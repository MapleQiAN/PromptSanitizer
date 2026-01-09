// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use serde::{Deserialize, Serialize};
use std::path::PathBuf;
use std::process::{Command, Stdio};
use tauri::Manager;

#[derive(Debug, Serialize, Deserialize)]
struct Request {
    text: String,
    mode: String,
    strategy: String,
    level: String,
    enabled_categories: Vec<String>,
    allowlist: Vec<String>,
    semantic_mode: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct Response {
    sanitized_text: String,
    findings: Vec<Finding>,
    stats: Stats,
    risk_score: i32,
    version: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct Finding {
    #[serde(rename = "type")]
    type_: String,
    start: usize,
    end: usize,
    confidence: f64,
    risk: i32,
    replacement: String,
    replacement_preview: String,
    reason: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct Stats {
    total_findings: i32,
    by_category: std::collections::HashMap<String, i32>,
    high_risk_count: i32,
    medium_risk_count: i32,
    low_risk_count: i32,
}

#[tauri::command]
async fn sanitize(app: tauri::AppHandle, request: Request) -> Result<Response, String> {
    // 获取 Go sidecar 路径
    let exe_path = get_sidecar_path(&app)?;

    // 准备请求 JSON
    let request_json = serde_json::to_string(&request)
        .map_err(|e| format!("序列化请求失败: {}", e))?;

    // 调用 Go sidecar
    let mut child = Command::new(&exe_path)
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()
        .map_err(|e| format!("启动 Go 引擎失败: {}. 请确保已构建 Go 二进制文件。", e))?;

    // 写入请求到 stdin
    if let Some(mut stdin) = child.stdin.take() {
        use std::io::Write;
        stdin
            .write_all(request_json.as_bytes())
            .map_err(|e| format!("写入请求失败: {}", e))?;
        stdin.flush().map_err(|e| format!("刷新stdin失败: {}", e))?;
    }

    let output = child
        .wait_with_output()
        .map_err(|e| format!("执行 Go 引擎失败: {}", e))?;

    if !output.status.success() {
        let stderr = String::from_utf8_lossy(&output.stderr);
        return Err(format!("Go 引擎执行失败: {}", stderr));
    }

    // 解析响应
    let stdout = String::from_utf8_lossy(&output.stdout);
    let response: Response = serde_json::from_str(&stdout)
        .map_err(|e| format!("解析响应失败: {}. 输出: {}", e, stdout))?;

    Ok(response)
}

#[tauri::command]
async fn open_file_dialog() -> Result<Option<String>, String> {
    // Tauri 2 使用不同的文件对话框 API
    // 这里简化实现，实际应该使用 tauri-plugin-dialog
    // 暂时返回错误，提示用户手动输入路径
    Err("文件对话框功能需要配置 tauri-plugin-dialog。请手动输入文件路径或使用粘贴功能。".to_string())
}

#[tauri::command]
async fn read_file(path: String) -> Result<String, String> {
    std::fs::read_to_string(&path).map_err(|e| format!("读取文件失败: {}", e))
}

fn get_sidecar_path(app: &tauri::AppHandle) -> Result<PathBuf, String> {
    // 生产环境：使用 Tauri 的路径解析器获取 externalBin
    if let Ok(resource_dir) = app.path().resource_dir() {
        #[cfg(target_os = "windows")]
        let exe_name = "prompt-sanitizer.exe";
        #[cfg(not(target_os = "windows"))]
        let exe_name = "prompt-sanitizer";
        
        let exe_path = resource_dir.join(exe_name);
        if exe_path.exists() {
            return Ok(exe_path);
        }
    }

    // 开发环境：尝试从多个可能的位置查找
    let possible_paths: Vec<PathBuf> = vec![
        // 开发环境：在项目根目录
        PathBuf::from("bin/prompt-sanitizer.exe"),
        PathBuf::from("bin/prompt-sanitizer"),
        PathBuf::from("src-tauri/bin/prompt-sanitizer.exe"),
        PathBuf::from("src-tauri/bin/prompt-sanitizer"),
        PathBuf::from("../../engine/go/cmd/main.exe"),
        PathBuf::from("../../engine/go/cmd/prompt-sanitizer.exe"),
    ];

    for path in possible_paths {
        if path.exists() {
            return Ok(path);
        }
    }

    Err("找不到 Go sidecar 二进制文件。请先构建 Go 引擎到 bin/ 目录。".to_string())
}

fn main() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![sanitize, open_file_dialog, read_file])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
