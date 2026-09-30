---
title: "程式碼撰寫標準"
description: "Wails v3 的程式碼風格、慣例與最佳實務"
slug: "contributing/standards"
sourcePath: "contributing/standards.md"
---

## 程式碼風格與慣例

遵循一致的程式碼撰寫標準，可讓程式碼庫更容易閱讀、維護及參與貢獻。

## Go 程式碼標準

### 程式碼格式

使用標準 Go 格式化工具：

```bash
# Format all code
gofmt -w .

# Use goimports for import organization
goimports -w .
```

<strong>必要條件：</strong>所有 Go 程式碼都必須先通過`gofmt`和`goimports`，才能提交。

### 命名慣例

**套件：**

- 使用小寫，並儘可能採用單一單字
- `package application`、`package events`
- 避免使用底線或混合大小寫

**匯出的名稱：**

- 型別、函式和常數使用 PascalCase
- `type WebviewWindow struct`、`func NewApplication()`

**未匯出的名稱：**

- 內部型別、函式和變數使用 camelCase
- `type windowImpl struct`、`func createWindow()`

**介面：**

- 依行為命名：`Reader`、`Writer`、`Handler`
- 單一方法介面：名稱使用`-er`後綴

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

### 錯誤處理

**一律檢查錯誤：**

```go
// Good
result, err := doSomething()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Bad - ignoring errors
result, _ := doSomething()
```

**使用錯誤包裝：**

```go
// Wrap errors to provide context
if err := validate(); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}
```

**需要時建立自訂錯誤型別：**

```go
type ValidationError struct {
    Field string
    Value string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid value %q for field %q", e.Value, e.Field)
}
```

### 註解與文件

**套件註解：**

```go
// Package application provides the core Wails application runtime.
//
// It handles window management, event dispatching, and service lifecycle.
package application
```

**匯出的宣告：**

```go
// NewApplication creates a new Wails application with the given options.
//
// The application must be started with Run() or RunWithContext().
func NewApplication(opts Options) *Application {
    // ...
}
```

**實作註解：**

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

### 函式與方法的結構

**讓函式聚焦於單一職責：**

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

**使用提前回傳：**

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

### 並行處理

**使用 context 進行取消：**

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

**使用互斥鎖保護共用狀態：**

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

**避免 goroutine 洩漏：**

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

### 測試

**測試檔案命名：**

```go
// Implementation: window.go
// Tests: window_test.go
```

**表格驅動測試：**

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

## JavaScript/TypeScript 標準

### 程式碼格式

使用 Prettier 確保格式一致：

```json
{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5"
}
```

### 命名慣例

**變數與函式：**

- camelCase：`const userName = "John"`

**類別與型別：**

- PascalCase：`class WindowManager`

**常數：**

- UPPER<em>SNAKE</em>CASE：`const MAX_RETRIES = 3`

### TypeScript

**使用明確型別：**

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

**定義介面：**

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

## 提交訊息格式

使用[約定式提交](https://www.conventionalcommits.org/)：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**類型：**

- `feat`：新功能
- `fix`：錯誤修正
- `docs`：文件變更
- `refactor`：程式碼重構
- `test`：新增或更新測試
- `chore`：維護工作

**範例：**

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

## Pull Request 指南

### 提交前

- [ ] 程式碼通過`gofmt`和`goimports`檢查
- [ ] 所有測試均通過（`go test ./...`）
- [ ] 新程式碼附有測試
- [ ] 視需要更新文件
- [ ] 提交訊息遵循慣例
- [ ] 與`master`沒有合併衝突

### PR 說明範本

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

## 程式碼審查流程

### 審查者須知

- 提供具建設性且尊重他人的意見
- 著重於程式碼品質，而非個人偏好
- 說明提出變更建議的原因
- 確認滿意後予以核准

### 作者須知

- 回覆所有意見
- 如有需要，請求進一步說明
- 進行要求的變更，或說明不變更的原因
- 虛心接受意見

## 最佳實務

### 效能

- 避免過早最佳化
- 先進行效能分析，再做最佳化
- 對效能關鍵程式碼使用基準測試

```go
func BenchmarkProcess(b *testing.B) {
    for i := 0; i < b.N; i++ {
        process(testData)
    }
}
```

### 安全性

- 驗證所有使用者輸入
- 顯示資料前先進行清理
- 使用`crypto/rand`產生隨機資料
- 絕不記錄敏感資訊

### 文件

- 為匯出的 API 撰寫文件
- 在文件中加入範例
- 變更 API 時一併更新文件
- 確保 README 檔案保持最新

## 平台特定程式碼

### 檔案命名

```
window.go           // Common interface
window_darwin.go    // macOS implementation
window_windows.go   // Windows implementation
window_linux.go     // Linux implementation
```

### 建置標籤

```go
//go:build darwin

package application

// macOS-specific code
```

## 程式碼檢查

提交前執行程式碼檢查工具：

```bash
# golangci-lint (recommended)
golangci-lint run

# Individual linters
go vet ./...
staticcheck ./...
```

## 有疑問？

如果不確定任何規範：

- 查看現有程式碼中的範例
- 在[Discord](https://discord.gg/JDdSxwjhGf)中詢問
- 在 GitHub 上發起討論
