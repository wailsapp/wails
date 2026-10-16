---
title: "编码规范"
description: "Wails v3 的代码风格、约定和最佳实践"
slug: "contributing/standards"
sourcePath: "contributing/standards.md"
---

## 代码风格与约定

遵循一致的编码规范可使代码库更易于阅读、维护和贡献。

## Go 代码规范

### 代码格式化

使用标准的 Go 格式化工具：

```bash
# Format all code
gofmt -w .

# Use goimports for import organization
goimports -w .
```

<strong>要求：</strong>提交前，所有 Go 代码都必须通过`gofmt`和`goimports`。

### 命名约定

**包：**

- 使用小写字母，并尽可能采用单个单词
- `package application`、`package events`
- 避免使用下划线或大小写混合形式

**导出名称：**

- 类型、函数和常量使用 PascalCase
- `type WebviewWindow struct`、`func NewApplication()`

**未导出名称：**

- 内部类型、函数和变量使用 camelCase
- `type windowImpl struct`、`func createWindow()`

**接口：**

- 按行为命名：`Reader`、`Writer`、`Handler`
- 单方法接口：名称使用`-er`后缀

```go
// Good
type Closer interface {
    Close() error
}

// Avoid
type CloseInterface interface {
    Close() error
}
```

### 错误处理

**始终检查错误：**

```go
// Good
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad - ignoring errors
result, _ := doSomething()
```

**使用错误包装：**

```go
// Wrap errors to provide context
if err := validate(); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}
```

**必要时创建自定义错误类型：**

```go
type ValidationError struct {
    Field string
    Value string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid value %q for field %q", e.Value, e.Field)
}
```

### 注释与文档

**包注释：**

```go
// Package application provides the core Wails application runtime.
//
// It handles window management, event dispatching, and service lifecycle.
package application
```

**导出声明：**

```go
// NewApplication creates a new Wails application with the given options.
//
// The application must be started with Run() or RunWithContext().
func NewApplication(opts Options) *Application {
    // ...
}
```

**实现注释：**

```go
// processEvent handles incoming events from the runtime.
// It dispatches to registered handlers and manages event lifecycle.
func (a *Application) processEvent(event *Event) {
    // Validate event before processing
    if event == nil {
        return
    }

    // Find and invoke handlers
    // ...
}
```

### 函数与方法结构

**保持函数职责集中：**

```go
// Good - single responsibility
func (w *Window) setTitle(title string) {
    w.title = title
    w.updateNativeTitle()
}

// Bad - doing too much
func (w *Window) updateEverything() {
    w.setTitle(w.title)
    w.setSize(w.width, w.height)
    w.setPosition(w.x, w.y)
    // ... 20 more operations
}
```

**使用提前返回：**

```go
// Good
func validate(input string) error {
    if input == "" {
        return errors.New("empty input")
    }

    if len(input) > 100 {
        return errors.New("input too long")
    }

    return nil
}

// Avoid deep nesting
```

### 并发

**使用 context 进行取消：**

```go
func (a *Application) RunWithContext(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-a.done:
        return nil
    }
}
```

**使用互斥锁保护共享状态：**

```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}
```

**避免 goroutine 泄漏：**

```go
// Good - goroutine has exit condition
func (a *Application) startWorker(ctx context.Context) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                return  // Clean exit
            case work := <-a.workChan:
                a.process(work)
            }
        }
    }()
}
```

### 测试

**测试文件命名：**

```go
// Implementation: window.go
// Tests: window_test.go
```

**表驱动测试：**

```go
func TestValidate(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"empty input", "", true},
        {"valid input", "hello", false},
        {"too long", strings.Repeat("a", 101), true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## JavaScript/TypeScript 规范

### 代码格式化

使用 Prettier 统一代码格式：

```json
{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5"
}
```

### 命名约定

**变量和函数：**

- camelCase：`const userName = "John"`

**类和类型：**

- PascalCase：`class WindowManager`

**常量：**

- UPPER<em>SNAKE</em>CASE：`const MAX_RETRIES = 3`

### TypeScript

**使用显式类型：**

```typescript
// Good
function greet(name: string): string {
    return `Hello, ${name}`
}

// Avoid implicit any
function process(data) {  // Bad
    return data
}
```

**定义接口：**

```typescript
interface WindowOptions {
    title: string
    width: number
    height: number
}

function createWindow(options: WindowOptions): void {
    // ...
}
```

## 提交消息格式

使用[约定式提交](https://www.conventionalcommits.org/)：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型：**

- `feat`：新功能
- `fix`：错误修复
- `docs`：文档变更
- `refactor`：代码重构
- `test`：添加或更新测试
- `chore`：维护任务

**示例：**

```
feat(window): add SetAlwaysOnTop method

Implement SetAlwaysOnTop for keeping windows above others.
Adds platform implementations for macOS, Windows, and Linux.

Closes #123
```

```
fix(events): prevent event handler memory leak

Event listeners were not being properly cleaned up when
windows were closed. This adds explicit cleanup in the
window destructor.
```

## 拉取请求指南

### 提交前检查

- [ ] 代码通过`gofmt`和`goimports`检查
- [ ] 所有测试均通过（`go test ./...`）
- [ ] 新代码包含测试
- [ ] 已按需更新文档
- [ ] 提交消息遵循约定
- [ ] 与`master`不存在合并冲突

### PR 描述模板

```markdown
## Description
Brief description of what this PR does.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
How was this tested?

## Checklist
- [ ] Tests pass
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

## 代码审查流程

### 作为审查者

- 提供建设性意见并尊重他人
- 关注代码质量，而非个人偏好
- 说明建议更改的原因
- 确认满意后批准

### 作为作者

- 回复所有评论
- 如有需要，请请求澄清
- 按要求进行更改，或说明不更改的原因
- 以开放的态度接受反馈

## 最佳实践

### 性能

- 避免过早优化
- 优化前先进行性能分析
- 对性能关键型代码使用基准测试

```go
func BenchmarkProcess(b *testing.B) {
    for i := 0; i < b.N; i++ {
        process(testData)
    }
}
```

### 安全性

- 验证所有用户输入
- 显示数据前对其进行清理
- 使用`crypto/rand`生成随机数据
- 切勿记录敏感信息

### 文档

- 为导出的 API 编写文档
- 在文档中提供示例
- 更改 API 时同步更新文档
- 确保 README 文件始终为最新版本

## 平台特定代码

### 文件命名

```
window.go           // Common interface
window_darwin.go    // macOS implementation
window_windows.go   // Windows implementation
window_linux.go     // Linux implementation
```

### 构建标签

```go
//go:build darwin

package application

// macOS-specific code
```

## 代码检查

提交前运行代码检查工具：

```bash
# golangci-lint (recommended)
golangci-lint run

# Individual linters
go vet ./...
staticcheck ./...
```

## 有疑问？

如果不确定任何规范：

- 查看现有代码中的示例
- 在[Discord](https://discord.gg/JDdSxwjhGf)中提问
- 在 GitHub 上发起讨论
