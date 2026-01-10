@echo off
REM 构建脚本 - Windows

echo 构建 Go 引擎...
cd engine\go
go build -ldflags="-H windowsgui" -o ..\..\apps\tauri\src-tauri\bin\prompt-sanitizer.exe .\cmd\main.go
if errorlevel 1 (
    echo Go 引擎构建失败
    exit /b 1
)
echo Go 引擎构建完成

echo 复制 Go 引擎为平台特定文件名...
cd ..\..\apps\tauri\src-tauri\bin
copy prompt-sanitizer.exe prompt-sanitizer-x86_64-pc-windows-msvc.exe
if errorlevel 1 (
    echo 复制失败
    exit /b 1
)
echo 平台特定文件创建完成

echo 构建 Tauri 应用...
cd ..\..\apps\tauri
call npm install
if errorlevel 1 (
    echo npm install 失败
    exit /b 1
)
call npm run tauri build
if errorlevel 1 (
    echo Tauri 构建失败
    exit /b 1
)

echo 构建完成！
