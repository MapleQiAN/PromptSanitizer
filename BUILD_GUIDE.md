# 构建指南 - Linux 和 macOS

本文档说明如何为 Linux 和 macOS 平台打包 PromptSanitizer。

## 📋 前置要求

### Linux
安装系统依赖：
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install libwebkit2gtk-4.1-dev \
  build-essential \
  curl \
  wget \
  file \
  libssl-dev \
  libgtk-3-dev \
  libayatana-appindicator3-dev \
  librsvg2-dev

# Fedora
sudo dnf install webkit2gtk4.1-devel.x86_64 \
  openssl-devel \
  curl \
  wget \
  file \
  libappindicator-gtk3 \
  librsvg2-devel

# Arch Linux
sudo pacman -S webkit2gtk \
  base-devel \
  curl \
  wget \
  openssl \
  appmenu-gtk-module \
  gtk3 \
  libappindicator-gtk3 \
  librsvg \
  libvips
```

### macOS
安装 Xcode Command Line Tools：
```bash
xcode-select --install
```

## 🚀 快速开始

### 方法一：使用专用脚本（推荐）

#### Linux
```bash
chmod +x build-linux.sh
./build-linux.sh
```

#### macOS
```bash
chmod +x build-macos.sh
./build-macos.sh
```

### 方法二：使用通用脚本

```bash
# 自动检测平台并构建
chmod +x build.sh
./build.sh

# 或指定平台
./build.sh linux   # 构建 Linux 版本
./build.sh macos   # 构建 macOS 版本
```

## 📦 构建产物

### Linux
构建完成后，安装包位于：
```
apps/tauri/src-tauri/target/release/bundle/appimage/
```
- **AppImage**: `prompt-sanitizer_0.1.0_amd64.AppImage` (便携式，无需安装)

可选格式：
- **DEB**: 使用 `--bundles deb` (Debian/Ubuntu)
- **RPM**: 使用 `--bundles rpm` (Fedora/RHEL/CentOS)

### macOS
构建完成后，安装包位于：
```
apps/tauri/src-tauri/target/<arch>/release/bundle/dmg/
```
- **DMG**: `PromptSanitizer_0.1.0_<arch>.dmg` (标准 macOS 安装方式)
- **App Bundle**: `PromptSanitizer.app` (可直接运行的应用程序包)

## 🔧 高级选项

### 指定构建格式

#### Linux
```bash
cd apps/tauri

# AppImage (默认)
npm run tauri build -- --bundles appimage

# DEB 包
npm run tauri build -- --bundles deb

# RPM 包
npm run tauri build -- --bundles rpm

# 同时构建多种格式
npm run tauri build -- --bundles appimage,deb,rpm
```

#### macOS
```bash
cd apps/tauri

# DMG (默认)
npm run tauri build -- --bundles dmg

# App Bundle
npm run tauri build -- --bundles app

# 同时构建多种格式
npm run tauri build -- --bundles dmg,app
```

### 交叉编译

#### 在 Linux 上构建 macOS 版本
需要安装 macOS 交叉编译工具链（较复杂，建议在 macOS 上构建）

#### 在 macOS 上构建 Linux 版本
需要安装 Linux 交叉编译工具链（较复杂，建议在 Linux 上构建）

### 指定目标架构

#### macOS
```bash
# Apple Silicon (arm64)
npm run tauri build -- --target aarch64-apple-darwin --bundles dmg

# Intel (x86_64)
npm run tauri build -- --target x86_64-apple-darwin --bundles dmg
```

#### Linux
```bash
# x86_64 (默认)
npm run tauri build -- --target x86_64-unknown-linux-gnu --bundles appimage
```

## 📝 手动构建步骤

如果脚本无法使用，可以手动执行以下步骤：

### 1. 构建 Go 引擎

#### Linux
```bash
cd engine/go
GOOS=linux GOARCH=amd64 go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer ./cmd/main.go
```

#### macOS
```bash
cd engine/go

# Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer ./cmd/main.go

# Intel
GOOS=darwin GOARCH=amd64 go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer ./cmd/main.go
```

### 2. 安装前端依赖
```bash
cd apps/tauri
npm install
```

### 3. 构建 Tauri 应用
```bash
cd apps/tauri

# Linux
npm run tauri build -- --bundles appimage

# macOS
npm run tauri build -- --bundles dmg
```

## 🐛 常见问题

### Linux: 缺少系统依赖
**问题**: 构建时提示缺少库文件

**解决**: 根据你的 Linux 发行版安装相应的系统依赖（见前置要求）

### macOS: 代码签名警告
**问题**: 构建成功但运行时提示"无法验证开发者"

**解决**: 
- 开发版本可以忽略此警告
- 生产版本需要配置代码签名（见 `docs/building.md`）

### Go 引擎未找到
**问题**: 构建时提示找不到 `prompt-sanitizer` 二进制文件

**解决**: 
- 确保 Go 引擎已构建到 `apps/tauri/src-tauri/bin/prompt-sanitizer`
- 检查文件名是否正确（Linux/macOS 无 `.exe` 扩展名）

## 📚 更多信息

详细构建说明请参考: `docs/building.md`
