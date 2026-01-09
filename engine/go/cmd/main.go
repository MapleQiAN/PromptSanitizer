package main

import (
	"encoding/json"
	"fmt"
	"github.com/prompt-sanitizer/engine/internal/engine"
	"github.com/prompt-sanitizer/engine/pkg/types"
	"io"
	"os"
	"strings"
)

func main() {
	// 从 stdin 读取 JSON 请求
	// 使用 io.ReadAll 一次性读取所有内容
	requestJSONBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		respondError(fmt.Sprintf("read stdin error: %v", err))
		return
	}

	// 移除末尾的空白字符（包括换行符）
	requestJSON := strings.TrimSpace(string(requestJSONBytes))

	if requestJSON == "" {
		respondError("no input received")
		return
	}

	// 解析请求
	var req types.Request
	if err := json.Unmarshal([]byte(requestJSON), &req); err != nil {
		respondError(fmt.Sprintf("invalid JSON: %v", err))
		return
	}

	// 验证请求
	if req.Text == "" {
		respondError("text field is required")
		return
	}

	// 设置默认值
	if req.Mode == "" {
		req.Mode = "sanitize"
	}
	if req.Strategy == "" {
		req.Strategy = "redact"
	}
	if req.Level == "" {
		req.Level = "standard"
	}
	if req.SemanticMode == "" {
		req.SemanticMode = "off"
	}

	// 创建引擎并处理
	eng := engine.NewEngine()
	resp, err := eng.Process(&req)
	if err != nil {
		respondError(fmt.Sprintf("processing error: %v", err))
		return
	}

	// 输出响应
	responseJSON, err := json.Marshal(resp)
	if err != nil {
		respondError(fmt.Sprintf("marshal error: %v", err))
		return
	}

	fmt.Println(string(responseJSON))
}

func respondError(message string) {
	errorResp := map[string]interface{}{
		"error": message,
	}
	jsonBytes, _ := json.Marshal(errorResp)
	fmt.Println(string(jsonBytes))
	os.Exit(1)
}
