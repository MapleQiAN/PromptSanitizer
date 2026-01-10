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

**安装包格式要求**：
- **NSIS 安装包（.exe）**：✅ 无需额外配置，推荐使用
- **MSI 安装包**：⚠️ 需要启用 Windows VBScript 功能（设置 → 应用 → 可选功能 → 更多 Windows 功能 → VBSCRIPT）

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
go build -ldflags="-H windowsgui" -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer.exe ./cmd/main.go
```

> 💡 **隐藏控制台窗口**：使用 `-ldflags="-H windowsgui"` 将 Go 程序编译为 Windows GUI 程序，这样运行时不会显示黑色命令行窗口。

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

#### NSIS 安装包（推荐，无需 VBScript）

- **NSIS 安装包**：`bundle/nsis/PromptSanitizer_0.1.0_x64-setup.exe`
  - ✅ **无需 VBScript**：不需要启用 Windows VBScript 功能
  - ✅ **无需 WiX Toolset**：使用 NSIS，更简单
  - ✅ **自定义安装程序**：更灵活的安装选项
  - ✅ **包含卸载程序**：自动生成卸载功能
  - 📦 **安装方式**：双击 `.exe` 文件运行安装向导

> 💡 **推荐使用 NSIS**：如果遇到 WiX/VBScript 相关问题，使用 NSIS 安装包更简单，无需额外配置。

#### MSI 安装包（需要 VBScript）

- **MSI 安装包**：`bundle/msi/PromptSanitizer_0.1.0_x64_en-US.msi`
  - 标准 Windows 安装程序
  - 支持自动更新
  - 包含卸载程序
  - ⚠️ **需要 VBScript**：需要启用 Windows VBScript 功能
  - ⚠️ **需要 WiX Toolset**：依赖 WiX 工具集

#### Windows 安装后的文件结构

安装后，应用文件通常位于：
```
C:\Users\<用户名>\AppData\Local\Programs\PromptSanitizer\
```

**安装后的文件结构**：

- `PromptSanitizer.exe` - 主应用程序（Tauri 前端 + Rust 后端，已嵌入图标）
- `resources/` - 资源目录
  - `prompt-sanitizer.exe` - Go 核心引擎（从 `resources` 配置打包）
  - `icon.ico` - 应用图标
- 其他 Tauri 运行时文件和依赖

> ⚠️ **重要**：如果安装后只有 `PromptSanitizer.exe` 一个文件，或缺少 `resources/` 目录，说明 Go sidecar 没有正确打包。请确保：
> 1. 在构建前，Go 二进制文件已构建到 `apps/tauri/src-tauri/bin/prompt-sanitizer.exe`
> 2. `tauri.conf.json` 中的 `externalBin` 配置包含 `"bin/prompt-sanitizer"`（不含扩展名）
> 3. `tauri.conf.json` 中的 `resources` 配置包含 `"bin/prompt-sanitizer.exe"`（含扩展名）
> 4. `Cargo.toml` 中包含 `[[bin]]` 配置，指定二进制名称为 `PromptSanitizer`（与 `productName` 一致）
> 5. 重新构建应用

> 💡 **应用图标**：应用图标通过 `build.rs` 嵌入到 exe 文件中，使用 `winres` crate。如果图标未显示，请确保：
> 1. `icons/icon.ico` 文件存在且格式正确
> 2. `Cargo.toml` 中包含 `winres` 依赖
> 3. `build.rs` 正确配置了图标路径

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
    "targets": ["nsis"],     // Windows: 仅 NSIS (.exe，推荐，无需 VBScript)
    // "targets": ["msi"],   // Windows: 仅 MSI (需要 VBScript)
    // "targets": ["dmg"],    // macOS: 仅 DMG
    // "targets": ["appimage"] // Linux: 仅 AppImage
  }
}
```

**Windows 安装包格式对比**：

| 格式 | 文件扩展名 | 需要 VBScript | 需要 WiX | 推荐度 |
|------|-----------|--------------|---------|--------|
| **NSIS** | `.exe` | ❌ 不需要 | ❌ 不需要 | ⭐⭐⭐⭐⭐ 推荐 |
| **MSI** | `.msi` | ✅ 需要 | ✅ 需要 | ⭐⭐⭐ |

> 💡 **建议**：如果遇到 VBScript 相关问题，使用 NSIS 格式（`.exe`）更简单。

#### 自定义 NSIS 安装程序

> ⚠️ **注意**：Tauri 2.0 对 NSIS 的自定义配置支持可能有限。如果遇到配置错误，建议使用默认配置。

**Tauri 2.0 默认的 NSIS 安装程序已经包含**：
- ✅ 完整的安装向导（不是一键安装）
- ✅ 用户可以选择安装位置
- ✅ 用户可以选择是否创建快捷方式
- ✅ 显示安装进度
- ✅ 提供卸载功能
- ✅ 使用应用图标（从 `bundle.icon` 配置中获取）

**确保安装程序使用正确的图标**：

1. ✅ **图标文件位置**：确保 `icons/icon.ico` 文件存在且格式正确
2. ✅ **图标配置顺序**：在 `bundle.icon` 中将 `icon.ico` 放在第一位，确保优先使用
3. ✅ **图标格式**：确保 `.ico` 文件包含多种尺寸（16x16, 32x32, 48x48, 256x256）

**如果安装程序图标不正确，尝试以下步骤**：

1. **重新生成图标**：
   ```bash
   # 在项目根目录准备 app-icon.png (512x512 或更大)
   npm run tauri icon
   ```

2. **清理并重新构建**：
   ```bash
   cd apps/tauri
   npm run tauri clean
   npm run tauri build
   ```

3. **检查图标文件**：确保 `apps/tauri/src-tauri/icons/icon.ico` 文件存在且不是损坏的

**让安装程序更专业的建议**（无需额外配置）：

1. ✅ **填写版权信息**：在 `bundle.copyright` 中填写版权声明（已配置）
2. ✅ **清晰的描述**：在 `bundle.shortDescription` 和 `bundle.longDescription` 中提供清晰的说明（已配置）
3. ✅ **应用图标**：确保 `bundle.icon` 中包含 `.ico` 文件，并将其放在列表第一位（已配置）

> 💡 **提示**：当前配置已经足够专业。Tauri 2.0 的默认 NSIS 安装程序行为已经比很多"流氓软件"更透明和用户友好。如果安装程序图标仍然不正确，可能需要重新生成图标文件。

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

### 1. WiX light.exe 执行失败（MSI 安装包）

**问题**：构建 MSI 安装包时提示 `failed to run light.exe`

> 💡 **快速解决方案**：如果不需要 MSI 格式，可以改用 NSIS 格式（`.exe`），无需 VBScript 和 WiX。在 `tauri.conf.json` 中设置 `"targets": ["nsis"]`。

**如果必须使用 MSI，可能原因和解决方案**：

#### 原因 1：VBSCRIPT 功能未启用（最常见）

WiX 工具集需要 Windows 的 VBSCRIPT 功能才能正常工作。

**解决步骤**：
1. 打开"设置" → "应用" → "可选功能" → "更多 Windows 功能"
2. 在列表中找到并勾选 **"VBSCRIPT"**
3. 点击"确定"并等待安装完成
4. 如果提示需要重启，请重启计算机
5. 重新运行构建命令

#### 原因 2：WiX 工具集下载不完整或损坏

**解决步骤**：
1. 删除现有的 WiX 工具集目录：
   ```powershell
   Remove-Item -Recurse -Force "$env:LOCALAPPDATA\tauri\WixTools314"
   ```
2. 重新运行构建，Tauri 会自动重新下载
3. 如果自动下载失败，可以手动下载：
   - 从 [WiX 官方发布页面](https://github.com/wixtoolset/wix3/releases/download/wix3141rtm/wix314-binaries.zip) 下载
   - 解压到 `C:\Users\<用户名>\AppData\Local\tauri\WixTools314\`
   - 确保 `light.exe` 文件存在于该目录

#### 原因 3：权限问题

**解决步骤**：
- 以管理员身份运行 PowerShell 或命令提示符
- 重新执行构建命令

#### 原因 4：productName 包含特殊字符

虽然当前配置看起来正常，但如果修改了 `productName`，确保：
- 只包含字母、数字和连字符
- 避免使用句点 (`.`)、空格等特殊字符

### 2. 安装后应用无法运行（找不到 Go 引擎）

**问题**：安装后运行应用，提示"找不到 Go sidecar 二进制文件"

**原因**：这通常是因为 `externalBin` 中的文件没有正确打包到安装包中，或路径解析不正确。

**解决**：
- 确保 Go 引擎已构建到 `apps/tauri/src-tauri/bin/` 目录
- 检查 `tauri.conf.json` 中的 `externalBin` 配置：`["bin/prompt-sanitizer"]`（不含扩展名）
- Windows 上确保源文件名为 `prompt-sanitizer.exe`
- 检查安装目录的 `resources/` 子目录，确认 `prompt-sanitizer.exe` 是否存在
- 如果不存在，重新构建并确保构建过程中没有错误

### 3. Go 二进制文件未找到（构建时）

**问题**：构建时提示找不到 `prompt-sanitizer` 二进制文件

**解决**：
- 确保 Go 引擎已构建到 `apps/tauri/src-tauri/bin/` 目录
- 检查文件名是否与 `tauri.conf.json` 中的 `externalBin` 配置一致
- Windows 上确保文件名包含 `.exe` 扩展名

### 4. Rust 编译错误

**问题**：`cargo build` 失败

**解决**：
- 检查 Rust 工具链：`rustup show`
- 更新 Rust：`rustup update`
- 清理构建缓存：`cargo clean`

### 5. 前端构建失败

**问题**：`npm run build` 失败

**解决**：
- 删除 `node_modules` 和 `package-lock.json`，重新安装：`npm install`
- 检查 Node.js 版本：`node --version`（需要 18+）

### 6. Windows 构建工具缺失

**问题**：提示找不到 C++ 编译器

**解决**：
- 安装 Visual Studio Build Tools
- 或安装 Visual Studio Community（包含 C++ 工作负载）

### 7. macOS 代码签名失败

**问题**：代码签名或公证失败

**解决**：
- 检查证书是否有效：`security find-identity -v -p codesigning`
- 确保环境变量配置正确
- 检查 Apple ID 是否启用了 App 专用密码

### 8. Linux 依赖缺失

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
