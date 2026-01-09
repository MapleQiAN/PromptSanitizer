# 快速开始指南

## 前置要求

- **Node.js** 18+ ([下载](https://nodejs.org/))
- **Rust** 1.70+ ([安装 rustup](https://rustup.rs/))
- **Go** 1.21+ ([下载](https://go.dev/dl/))
- **Tauri CLI**: `npm install -g @tauri-apps/cli`

## Windows 快速开始

### 1. 安装依赖

```powershell
# 安装 Node.js 依赖
npm install

# 安装 Tauri CLI（如果还没有）
npm install -g @tauri-apps/cli
```

### 2. 构建 Go 引擎

```powershell
cd engine\go
go mod download
go build -o ..\..\apps\tauri\src-tauri\bin\prompt-sanitizer.exe .\cmd\main.go
```

### 3. 运行应用（开发模式）

```powershell
cd ..\..\apps\tauri
npm install
npm run tauri dev
```

### 4. 打包应用

```powershell
cd apps\tauri
npm run tauri build
```

打包后的安装包位于：`apps/tauri/src-tauri/target/release/bundle/msi/`

## Linux/macOS 快速开始

### 1. 安装依赖

```bash
npm install
```

### 2. 构建 Go 引擎

```bash
cd engine/go
go mod download
go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer ./cmd/main.go
```

### 3. 运行应用（开发模式）

```bash
cd ../../apps/tauri
npm install
npm run tauri dev
```

### 4. 打包应用

```bash
npm run tauri build
```

## 使用一键构建脚本

### Windows

```powershell
.\build.bat
```

### Linux/macOS

```bash
chmod +x build.sh
./build.sh
```

## 测试 Go 引擎

在构建完成后，可以手动测试 Go 引擎：

### Windows

```powershell
cd engine\go
echo '{"text":"我的手机号是13812345678","mode":"sanitize","strategy":"redact","level":"standard","enabled_categories":[],"allowlist":[],"semantic_mode":"off"}' | .\prompt-sanitizer.exe
```

### Linux/macOS

```bash
cd engine/go
echo '{"text":"我的手机号是13812345678","mode":"sanitize","strategy":"redact","level":"standard","enabled_categories":[],"allowlist":[],"semantic_mode":"off"}' | ./prompt-sanitizer
```

## 常见问题

### Q: 找不到 Go 二进制文件

**A**: 确保已构建 Go 引擎并放在正确位置：
- Windows: `apps/tauri/src-tauri/bin/prompt-sanitizer.exe`
- Linux/macOS: `apps/tauri/src-tauri/bin/prompt-sanitizer`

### Q: Rust 编译错误

**A**: 确保已安装 Rust 工具链：
```bash
rustup update
```

### Q: Tauri 构建失败

**A**: 
1. 确保所有依赖已安装：`npm install`
2. 确保 Go 引擎已构建
3. 检查 Tauri 版本兼容性

### Q: 应用启动后无法调用 Go 引擎

**A**: 
1. 检查 Go 二进制文件是否存在
2. 检查文件权限（Linux/macOS）
3. 查看控制台错误信息

## 下一步

- 查看 [使用指南](docs/usage.md) 了解如何使用应用
- 查看 [开发文档](docs/development.md) 了解如何扩展功能
- 查看 [协议文档](docs/protocol.md) 了解通信协议
