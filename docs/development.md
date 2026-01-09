# 开发指南

## 项目结构

```
.
├── apps/
│   └── tauri/              # Tauri 应用
│       ├── src/            # 前端代码（TypeScript + React）
│       └── src-tauri/      # Rust 代码（Tauri 壳）
├── engine/
│   └── go/                 # Go 核心引擎
│       ├── cmd/            # 主程序入口
│       ├── internal/       # 内部包
│       │   ├── detector/   # 检测器
│       │   ├── sanitizer/  # 清洗器
│       │   └── engine/     # 引擎核心
│       ├── pkg/            # 公开包
│       │   └── types/       # 类型定义
│       └── testdata/        # 测试数据
└── docs/                   # 文档
```

## 开发环境设置

### 1. 安装工具链

- **Node.js**: 18+
- **Rust**: 1.70+ (通过 rustup 安装)
- **Go**: 1.21+
- **Tauri CLI**: `npm install -g @tauri-apps/cli`

### 2. 初始化项目

```bash
# 安装前端依赖
cd apps/tauri
npm install

# 初始化 Rust 项目（如果需要）
cd src-tauri
cargo build
```

### 3. 开发流程

#### Go 引擎开发

```bash
cd engine/go

# 运行测试
go test ./...

# 构建
go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer.exe ./cmd/main.go

# 手动测试（通过 stdin）
echo '{"text":"test","mode":"sanitize","strategy":"redact","level":"standard","enabled_categories":[],"allowlist":[],"semantic_mode":"off"}' | ./prompt-sanitizer.exe
```

#### 前端开发

```bash
cd apps/tauri

# 开发模式（只启动前端，不启动 Tauri）
npm run dev

# Tauri 开发模式（完整应用）
npm run tauri dev
```

#### Rust 开发

```bash
cd apps/tauri/src-tauri

# 构建
cargo build

# 运行测试
cargo test
```

## 添加新的检测类别

1. 在 `engine/go/internal/detector/detector.go` 中添加新的 `Category` 常量
2. 实现新的检测器（实现 `Detector` 接口）
3. 在 `engine/go/internal/engine/engine.go` 的 `NewEngine()` 中注册新检测器
4. 更新前端类型定义和 UI（`apps/tauri/src/types.ts` 和组件）

示例：

```go
// 1. 添加类别
const CategoryCreditCard Category = "credit_card"

// 2. 实现检测器
type CreditCardDetector struct {
    BaseDetector
}

func NewCreditCardDetector() *CreditCardDetector {
    return &CreditCardDetector{BaseDetector{category: CategoryCreditCard}}
}

func (d *CreditCardDetector) Detect(text string, level string) []Finding {
    // 实现检测逻辑
}

// 3. 注册检测器
detectors: []detector.Detector{
    // ... 其他检测器
    detector.NewCreditCardDetector(),
}
```

## 添加新的清洗策略

1. 在 `engine/go/internal/sanitizer/sanitizer.go` 中添加新的策略处理逻辑
2. 更新 `Request` 类型的 `strategy` 字段文档
3. 更新前端 UI 的选项

## 测试

### Go 单元测试

```bash
cd engine/go
go test ./...
go test -v ./internal/detector  # 详细输出
```

### 集成测试

手动测试完整流程：

1. 构建 Go 引擎
2. 启动 Tauri 应用
3. 输入测试文本
4. 验证结果

### 测试数据

测试样例位于 `engine/go/testdata/sample.txt`，包含各种敏感信息类型。

## 构建和打包

### 开发构建

```bash
# 构建 Go 引擎
cd engine/go
go build -o ../../apps/tauri/src-tauri/bin/prompt-sanitizer.exe ./cmd/main.go

# 构建 Tauri 应用
cd ../../apps/tauri
npm run tauri build
```

### 生产构建

确保：
1. Go 引擎已构建并放在正确位置
2. 所有依赖已安装
3. 运行完整测试

```bash
cd apps/tauri
npm run tauri build -- --release
```

## 代码规范

### Go

- 使用 `gofmt` 格式化代码
- 遵循 Go 官方代码规范
- 添加必要的注释和文档

### TypeScript

- 使用 ESLint 和 Prettier
- 遵循 React 最佳实践
- 使用 TypeScript 严格模式

### Rust

- 使用 `rustfmt` 格式化代码
- 遵循 Rust 官方代码规范
- 处理所有可能的错误情况

## 调试

### Go 引擎调试

```bash
# 使用 delve 调试器
dlv debug ./cmd/main.go

# 或添加日志输出
```

### Tauri 调试

- 前端：使用浏览器开发者工具（开发模式下）
- Rust：使用 `println!` 或日志库

## 性能优化

- Go 引擎：使用 `go test -bench` 进行基准测试
- 前端：使用 React DevTools Profiler
- 整体：使用性能分析工具

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 创建 Pull Request

请确保：
- 代码通过所有测试
- 添加必要的文档
- 遵循代码规范
