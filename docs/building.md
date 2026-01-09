# 打包指南

本文档说明如何将 PromptSanitizer 打包为各平台的安装包。

## 📋 前置要求

### 开发环境

| 工具 | 版本要求 | 说明 |
|------|---------|------|
| **Node.js** | 18+ | JavaScript 运行时 |
| **Rust** | 1.70+ | 系统编程语言（通过 rustup 安装） |
| **Go** | 1.21+ | 核心引擎语言 |
| **Tauri CLI** | Latest | `npm install -g @tauri-apps/cli` |

### 平台特定要求

#### Windows
- **Visual Studio Build Tools** 或 **Visual Studio Community**（包含 C++ 构建工具）
- **Windows SDK**
- 安装命令：`winget install Microsoft.VisualStudio.2022.BuildTools --override "--wait --add Microsoft.VisualStudio.Workload.VCTools"`

#### macOS
- **Xcode Command Line Tools**：`xcode-select --install`
- **Apple Developer Certificate**（用于代码签名，可选但推荐）

#### Linux
- **系统依赖**：
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

## 🚀 打包步骤

### 方法一：使用构建脚本（推荐）

#### Windows

```bash
# 在项目根目录执行
.\build.bat
```

构建脚本会自动：
1. 构建 Go 引擎到 `apps/tauri/src-tauri/bin/prompt-sanitizer.exe`
2. 安装前端依赖
3. 构建 Tauri 应用并生成安装包

#### Linux/macOS

```bash
# 在项目根目录执行
chmod +x build.sh
./build.sh
```

### 方法二：手动构建

#### 1. 构建 Go 引擎

首先需要构建 Go sidecar 二进制文件：

**Windows:**
```bash
cd engine/go
go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer.exe ./cmd/main.go
```

**Linux/macOS:**
```bash
cd engine/go
go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer ./cmd/main.go
```

> ⚠️ **重要**：Go 二进制文件名必须与 `tauri.conf.json` 中 `externalBin` 配置一致（不含扩展名）。

#### 2. 安装前端依赖

```bash
cd apps/tauri
npm install
```

#### 3. 构建 Tauri 应用

```bash
cd apps/tauri
npm run tauri build
```

或者使用 Tauri CLI 直接构建：

```bash
cd apps/tauri
tauri build
```

## 📦 构建产物

构建完成后，安装包会生成在以下目录：

```
apps/tauri/src-tauri/target/release/bundle/
```

### Windows 构建产物

- **MSI 安装包**：`bundle/msi/PromptSanitizer_0.1.0_x64_en-US.msi`
  - 标准 Windows 安装程序
  - 支持自动更新
  - 包含卸载程序

- **NSIS 安装包**（如果配置了）：`bundle/nsis/PromptSanitizer_0.1.0_x64-setup.exe`
  - 自定义安装程序
  - 更灵活的安装选项

### macOS 构建产物

- **DMG 磁盘镜像**：`bundle/dmg/PromptSanitizer_0.1.0_x64.dmg`
  - 标准 macOS 安装方式
  - 拖拽安装

- **App Bundle**：`bundle/macos/PromptSanitizer.app`
  - 可直接运行的应用程序包

### Linux 构建产物

根据配置，可能生成以下格式之一：

- **AppImage**：`bundle/appimage/prompt-sanitizer_0.1.0_amd64.AppImage`
  - 便携式应用，无需安装
  - 所有依赖打包在一个文件中

- **DEB 包**：`bundle/deb/prompt-sanitizer_0.1.0_amd64.deb`
  - Debian/Ubuntu 系统安装包

- **RPM 包**：`bundle/rpm/prompt-sanitizer-0.1.0-1.x86_64.rpm`
  - Fedora/RHEL/CentOS 系统安装包

## ⚙️ 构建配置

### 自定义构建选项

#### 指定目标平台

```bash
# 构建特定平台（需要交叉编译工具链）
tauri build --target x86_64-pc-windows-msvc    # Windows 64位
tauri build --target x86_64-apple-darwin       # macOS Intel
tauri build --target aarch64-apple-darwin      # macOS Apple Silicon
tauri build --target x86_64-unknown-linux-gnu  # Linux 64位
```

#### 调试构建

```bash
# 开发模式构建（未优化，包含调试信息）
tauri build --debug
```

#### 仅构建特定格式

在 `tauri.conf.json` 中配置 `bundle.targets`：

```json
{
  "bundle": {
    "targets": "msi",        // Windows: 仅 MSI
    // "targets": "dmg",      // macOS: 仅 DMG
    // "targets": "appimage" // Linux: 仅 AppImage
  }
}
```

### 代码签名（生产环境）

#### Windows

1. 获取代码签名证书（.pfx 文件）
2. 在 `tauri.conf.json` 中配置：

```json
{
  "bundle": {
    "windows": {
      "certificateThumbprint": "YOUR_CERTIFICATE_THUMBPRINT",
      "digestAlgorithm": "sha256",
      "timestampUrl": "http://timestamp.digicert.com"
    }
  }
}
```

#### macOS

1. 在 Apple Developer 中创建证书
2. 配置环境变量：

```bash
export APPLE_CERTIFICATE="Developer ID Application: Your Name"
export APPLE_CERTIFICATE_PASSWORD="your_password"
export APPLE_SIGNING_IDENTITY="Developer ID Application: Your Name"
export APPLE_ID="your@email.com"
export APPLE_PASSWORD="app-specific-password"
export APPLE_TEAM_ID="YOUR_TEAM_ID"
```

3. 构建时会自动进行代码签名和公证

#### Linux

Linux 通常不需要代码签名，但可以为 AppImage 添加 GPG 签名：

```bash
gpg --detach-sign --armor PromptSanitizer_0.1.0_amd64.AppImage
```

## 🔍 验证构建

### 检查二进制文件

```bash
# Windows
file apps/tauri/src-tauri/target/release/prompt-sanitizer.exe

# macOS/Linux
file apps/tauri/src-tauri/target/release/prompt-sanitizer
```

### 测试安装包

#### Windows

```bash
# 安装测试
.\bundle\msi\PromptSanitizer_0.1.0_x64_en-US.msi

# 卸载测试
# 通过控制面板或使用 msiexec
msiexec /x {PRODUCT_GUID}
```

#### macOS

```bash
# 挂载 DMG
open bundle/dmg/PromptSanitizer_0.1.0_x64.dmg

# 检查签名
codesign --verify --deep --strict PromptSanitizer.app
spctl --assess --verbose PromptSanitizer.app
```

#### Linux

```bash
# AppImage
chmod +x prompt-sanitizer_0.1.0_amd64.AppImage
./prompt-sanitizer_0.1.0_amd64.AppImage

# DEB
sudo dpkg -i prompt-sanitizer_0.1.0_amd64.deb
sudo apt-get install -f  # 修复依赖

# RPM
sudo rpm -i prompt-sanitizer-0.1.0-1.x86_64.rpm
```

## 🐛 常见问题

### 1. Go 二进制文件未找到

**问题**：构建时提示找不到 `prompt-sanitizer` 二进制文件

**解决**：
- 确保 Go 引擎已构建到 `apps/tauri/src-tauri/bin/` 目录
- 检查文件名是否与 `tauri.conf.json` 中的 `externalBin` 配置一致
- Windows 上确保文件名包含 `.exe` 扩展名

### 2. Rust 编译错误

**问题**：`cargo build` 失败

**解决**：
- 检查 Rust 工具链：`rustup show`
- 更新 Rust：`rustup update`
- 清理构建缓存：`cargo clean`

### 3. 前端构建失败

**问题**：`npm run build` 失败

**解决**：
- 删除 `node_modules` 和 `package-lock.json`，重新安装：`npm install`
- 检查 Node.js 版本：`node --version`（需要 18+）

### 4. Windows 构建工具缺失

**问题**：提示找不到 C++ 编译器

**解决**：
- 安装 Visual Studio Build Tools
- 或安装 Visual Studio Community（包含 C++ 工作负载）

### 5. macOS 代码签名失败

**问题**：代码签名或公证失败

**解决**：
- 检查证书是否有效：`security find-identity -v -p codesigning`
- 确保环境变量配置正确
- 检查 Apple ID 是否启用了 App 专用密码

### 6. Linux 依赖缺失

**问题**：构建时提示缺少系统库

**解决**：
- 根据发行版安装相应的系统依赖（见前置要求）
- 使用包管理器搜索缺失的库：`apt search <package-name>`

## 📝 发布清单

发布前请确认：

- [ ] Go 引擎已构建并包含在 bundle 中
- [ ] 版本号已更新（`tauri.conf.json` 和 `package.json`）
- [ ] 应用图标已正确配置
- [ ] 代码签名已配置（生产环境）
- [ ] 所有测试通过
- [ ] 构建产物已测试安装和运行
- [ ] 更新日志已更新
- [ ] 文档已更新

## 🔗 相关资源

- [Tauri 官方文档 - 构建](https://tauri.app/v2/guides/building/)
- [Tauri 官方文档 - 打包](https://tauri.app/v2/guides/bundling/)
- [Tauri 官方文档 - 代码签名](https://tauri.app/v2/guides/distribution/sign-windows/)
