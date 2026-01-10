#!/bin/bash
# 构建脚本 - Linux/macOS
# 用法: ./build.sh [linux|macos|auto]
# 如果不指定平台，将自动检测当前系统

set -e

# 检测平台
detect_platform() {
    case "$(uname -s)" in
        Linux*)
            echo "linux"
            ;;
        Darwin*)
            echo "macos"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

# 获取目标平台
TARGET_PLATFORM="${1:-auto}"

if [ "$TARGET_PLATFORM" = "auto" ]; then
    TARGET_PLATFORM=$(detect_platform)
fi

echo "=========================================="
echo "构建 PromptSanitizer - $TARGET_PLATFORM"
echo "=========================================="

# 构建 Go 引擎
echo ""
echo "步骤 1/3: 构建 Go 引擎..."
cd engine/go

if [ "$TARGET_PLATFORM" = "linux" ]; then
    echo "构建 Linux 版本..."
    GOOS=linux GOARCH=amd64 go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer ./cmd/main.go
elif [ "$TARGET_PLATFORM" = "macos" ]; then
    # 检测 macOS 架构
    MACOS_ARCH=$(uname -m)
    if [ "$MACOS_ARCH" = "arm64" ]; then
        echo "构建 macOS Apple Silicon (arm64) 版本..."
        GOOS=darwin GOARCH=arm64 go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer ./cmd/main.go
    else
        echo "构建 macOS Intel (x86_64) 版本..."
        GOOS=darwin GOARCH=amd64 go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer ./cmd/main.go
    fi
else
    echo "构建当前平台版本..."
    go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer ./cmd/main.go
fi

echo "✓ Go 引擎构建完成"

# 安装前端依赖
echo ""
echo "步骤 2/3: 安装前端依赖..."
cd ../../apps/tauri
if [ ! -d "node_modules" ]; then
    npm install
else
    echo "依赖已存在，跳过安装"
fi

# 构建 Tauri 应用
echo ""
echo "步骤 3/3: 构建 Tauri 应用..."

if [ "$TARGET_PLATFORM" = "linux" ]; then
    echo "构建 Linux AppImage..."
    npm run tauri build -- --bundles appimage
elif [ "$TARGET_PLATFORM" = "macos" ]; then
    echo "构建 macOS DMG..."
    npm run tauri build -- --bundles dmg
else
    echo "构建当前平台默认格式..."
    npm run tauri build
fi

echo ""
echo "=========================================="
echo "✓ 构建完成！"
echo "=========================================="
echo ""
echo "构建产物位置:"
echo "  apps/tauri/src-tauri/target/release/bundle/"
echo ""
