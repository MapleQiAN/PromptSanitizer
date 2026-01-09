# 使用指南

## 快速开始

### 1. 安装依赖

```bash
# 安装 Node.js 依赖
npm install

# 安装 Rust 工具链（如果还没有）
# Windows: 下载并安装 https://rustup.rs/
# macOS/Linux: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# 安装 Go（如果还没有）
# 下载: https://go.dev/dl/
```

### 2. 构建 Go 引擎

```bash
# Windows
cd engine/go
go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer.exe ./cmd/main.go

# Linux/macOS
cd engine/go
go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer ./cmd/main.go
```

### 3. 开发模式运行

```bash
cd apps/tauri
npm install
npm run tauri dev
```

### 4. 打包应用

```bash
cd apps/tauri
npm run tauri build
```

打包后的安装包位于 `apps/tauri/src-tauri/target/release/bundle/`

## 使用说明

### 基本操作

1. **输入文本**: 在左侧文本框中粘贴或输入要清洗的文本
2. **导入文件**: 点击"导入文件"按钮选择文本文件（.txt, .md, .json, .yaml, .yml）
3. **一键清洗**: 点击"一键清洗"按钮执行清洗
4. **查看结果**: 右侧显示清洗后的文本，底部显示命中列表

### 配置选项

点击"显示配置"打开配置面板：

- **处理模式**: 
  - 清洗模式：执行替换
  - 标注模式：只识别不替换，用于预览
  
- **清洗策略**:
  - 占位符替换：替换为 `[REDACTED:TYPE]`
  - 部分打码：保留前后缀，中间打码
  - 一致化替换：同一实体使用相同代号

- **清洗强度**:
  - 宽松：只识别明显高风险内容
  - 标准：平衡误报和漏报（推荐）
  - 严格：尽可能识别所有可能的敏感信息

- **检测类别**: 选择要检测的敏感信息类型

- **白名单**: 输入要排除的字符串，每行一个（精确匹配）

### 视图模式

- **分栏视图**: 左右分栏显示原文和清洗结果
- **对比视图**: 上下对比显示原文和清洗结果
- **报告视图**: 显示详细的统计信息和报告

### 导出报告

在报告视图中点击"导出 JSON"按钮，可以导出清洗报告（不包含完整敏感内容）。

## 支持的敏感信息类型

- **手机号**: 中国手机号（11位，1开头）
- **邮箱**: 标准邮箱格式
- **身份证**: 中国居民身份证号（18位）
- **IP地址**: IPv4 地址
- **域名/URL**: 域名和 URL 地址
- **Token/Key**: API Key、Bearer Token、JWT、Cookie 等
- **密码**: 密码字段（如 `password=xxx`）
- **私钥**: PEM 格式私钥

## 注意事项

1. **隐私保护**: 所有处理完全在本地进行，不会上传任何内容
2. **误报处理**: 如果发现误报，可以使用白名单功能排除
3. **性能**: 对于大文本（>50k字符），处理可能需要几秒钟
4. **备份**: 建议在清洗前备份原始文本

## 常见问题

### Q: Go 引擎找不到？
A: 确保已构建 Go 二进制文件，并放在正确的位置（`apps/tauri/src-tauri/bin/`）

### Q: 清洗结果不准确？
A: 可以调整清洗强度或使用白名单功能排除误报

### Q: 如何添加自定义规则？
A: 当前版本不支持自定义规则，未来版本会添加此功能
