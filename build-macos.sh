#!/bin/bash
# macOS 专用构建脚本

set -e

echo "=========================================="
echo "构建 PromptSanitizer - macOS"
echo "=========================================="

# 检测 macOS 架构
MACOS_ARCH=$(uname -m)
if [ "$MACOS_ARCH" = "arm64" ]; then
    echo "检测到 Apple Silicon (arm64)"
    GOARCH=arm64
    TARGET_ARCH="aarch64-apple-darwin"
else
    echo "检测到 Intel (x86_64)"
    GOARCH=amd64
    TARGET_ARCH="x86_64-apple-darwin"
fi

# 构建 Go 引擎 (macOS)
echo ""
echo "步骤 1/3: 构建 Go 引擎 (macOS $MACOS_ARCH)..."
cd engine/go
GOOS=darwin GOARCH=$GOARCH go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer ./cmd/main.go
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

# 构建 Tauri 应用 (macOS DMG)
echo ""
echo "步骤 3/3: 构建 Tauri 应用 (DMG)..."
npm run tauri build -- --bundles dmg --target $TARGET_ARCH

echo ""
echo "=========================================="
echo "✓ macOS 构建完成！"
echo "=========================================="
echo ""
echo "构建产物位置:"
echo "  apps/tauri/src-tauri/target/$TARGET_ARCH/release/bundle/dmg/"
echo ""
echo "可选格式:"
echo "  - DMG: --bundles dmg (默认)"
echo "  - App: --bundles app"
echo ""
