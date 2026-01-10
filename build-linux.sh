#!/bin/bash
# Linux 专用构建脚本

set -e

echo "=========================================="
echo "构建 PromptSanitizer - Linux"
echo "=========================================="

# 构建 Go 引擎 (Linux)
echo ""
echo "步骤 1/3: 构建 Go 引擎 (Linux)..."
cd engine/go
GOOS=linux GOARCH=amd64 go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer ./cmd/main.go
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

# 构建 Tauri 应用 (Linux AppImage)
echo ""
echo "步骤 3/3: 构建 Tauri 应用 (AppImage)..."
npm run tauri build -- --bundles appimage

echo ""
echo "=========================================="
echo "✓ Linux 构建完成！"
echo "=========================================="
echo ""
echo "构建产物位置:"
echo "  apps/tauri/src-tauri/target/release/bundle/appimage/"
echo ""
echo "可选格式:"
echo "  - AppImage: --bundles appimage (默认)"
echo "  - DEB:      --bundles deb"
echo "  - RPM:      --bundles rpm"
echo ""
