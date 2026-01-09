#!/bin/bash
# 构建脚本 - Linux/macOS

set -e

echo "构建 Go 引擎..."
cd engine/go
go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer ./cmd/main.go
echo "Go 引擎构建完成"

echo "构建 Tauri 应用..."
cd ../../apps/tauri
npm install
npm run tauri build

echo "构建完成！"
