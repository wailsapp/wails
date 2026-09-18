---
title: "Go－前端橋接器"
description: "深入瞭解 Wails 如何讓 Go 與 JavaScript 直接通訊"
slug: "concepts/bridge"
sourcePath: "concepts/bridge.md"
---

## Go 與 JavaScript 直接通訊

Wails 在 Go 與 JavaScript 之間提供<strong>直接的記憶體內橋接</strong>，無須承受 HTTP 額外負荷、跨越處理程序邊界，也不會遭遇序列化瓶頸，即可順暢通訊。

## 整體架構

```d2
direction: right

Frontend: 前端（JavaScript） {
  UI: React/Vue/原生 JavaScript {
    shape: rectangle
    style.fill: "#8B5CF6"
  }

  Bindings: 自動產生的繫結 {
    shape: rectangle
    style.fill: "#A78BFA"
  }
}

Bridge: Wails 橋接器 {
  Encoder: JSON 編碼器 {
    shape: rectangle
    style.fill: "#10B981"
  }

  Router: 方法路由器 {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: JSON 解碼器 {
    shape: rectangle
    style.fill: "#10B981"
  }

  TypeGen: 型別產生器 {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Backend: 後端（Go） {
  Services: 您的服務 {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Registry: 服務登錄庫 {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend.UI -> Frontend.Bindings: "import { Method }"
Frontend.Bindings -> Bridge.Encoder: "呼叫 Method('arg')"
Bridge.Encoder -> Bridge.Router: 編碼為 JSON
Bridge.Router -> Backend.Registry: 尋找服務
Backend.Registry -> Backend.Services: 叫用方法
Backend.Services -> Bridge.Decoder: 傳回結果
Bridge.Decoder -> Frontend.Bindings: 解碼為 JS
Frontend.Bindings -> Frontend.UI: Promise 完成
Bridge.TypeGen -> Frontend.Bindings: 產生型別
```

<strong>關鍵要點：</strong>沒有 HTTP、沒有 IPC，也沒有處理程序邊界。只有具備<strong>型別安全性</strong>的<strong>直接函式呼叫</strong>。

## 運作方式：逐步說明

### 1. 服務註冊（啟動時）

應用程式啟動時，Wails 會掃描您的服務：

```go
type GreetService struct {
    prefix string
}

func (g *GreetService) Greet(name string) string {
    return g.prefix + name + "!"
}

func (g *GreetService) Add(a, b int) int {
    return a + b
}

// Register service
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{prefix: "Hello, "}),
    },
})
```

**Wails 執行的作業：**

1. <strong>掃描結構體</strong>以尋找匯出的方法
2. **擷取型別資訊**（參數、傳回型別）
3. **建立登錄庫**，將方法名稱對應至函式
4. **產生 TypeScript 繫結**，其中包含完整的型別定義

### 2. 產生繫結（建置時）

Wails 會自動產生 TypeScript 繫結：

```typescript
// Auto-generated: frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function Add(a: number, b: number): Promise<number>
```

**型別對應：**

| Go 型別 | TypeScript 型別 |
| --- | --- |
| `string` | `string` |
| `int`、`int32`、`int64` | `number` |
| `float32`、`float64` | `number` |
| `bool` | `boolean` |
| `[]T` | `T[]` |
| `map[string]T` | `Record<string, T>` |
| `struct` | `interface` |
| `time.Time` | `Date` |
| `error` | 例外（擲出） |

### 3. 前端呼叫（執行階段）

開發人員從 JavaScript 呼叫 Go 方法：

```javascript
import { Greet, Add } from './bindings/GreetService'

// Call Go from JavaScript
const greeting = await Greet("World")
console.log(greeting)  // "Hello, World!"

const sum = await Add(5, 3)
console.log(sum)  // 8
```

**執行流程：**

1. **呼叫繫結函式**－`Greet("World")`
2. **建立訊息**－`{ service: "GreetService", method: "Greet", args: ["World"] }`
3. **傳送至橋接器**—透過 WebView 的 JavaScript 橋接器
4. **傳回 Promise**—等待回應

### 4. 橋接處理（執行階段）

橋接器接收訊息並進行處理：

```d2
direction: down

Receive: 接收訊息 {
  shape: rectangle
  style.fill: "#10B981"
}

Parse: 剖析 JSON {
  shape: rectangle
}

Validate: 驗證 {
  Check: 服務是否存在？ {
    shape: diamond
  }

  CheckMethod: 方法是否存在？ {
    shape: diamond
  }

  CheckTypes: 型別是否正確？ {
    shape: diamond
  }
}

Invoke: 叫用 Go 方法 {
  shape: rectangle
  style.fill: "#00ADD8"
}

Encode: 編碼結果 {
  shape: rectangle
}

Send: 傳送回應 {
  shape: rectangle
  style.fill: "#10B981"
}

Error: 傳送錯誤 {
  shape: rectangle
  style.fill: "#EF4444"
}

Receive -> Parse
Parse -> Validate.Check
Validate.Check -> Validate.CheckMethod: 是
Validate.Check -> Error: 否
Validate.CheckMethod -> Validate.CheckTypes: 是
Validate.CheckMethod -> Error: 否
Validate.CheckTypes -> Invoke: 是
Validate.CheckTypes -> Error: 否
Invoke -> Encode: 成功
Invoke -> Error: 錯誤
Encode -> Send
```

<strong>安全性：</strong>只能叫用已註冊的服務與已匯出的方法。

### 5. Go 執行（執行階段）

Go 方法會執行：

```go
func (g *GreetService) Greet(name string) string {
    // This runs in Go
    return g.prefix + name + "!"
}
```

**執行環境：**

- 在 **goroutine** 中執行（非阻塞）
- 可存取 **Go 的所有功能**（檔案系統、網路、資料庫）
- 可自由呼叫<strong>其他 Go 程式碼</strong>
- 傳回結果或錯誤

### 6. 回應（執行階段）

結果會傳回 JavaScript：

```javascript
// Promise resolves with result
const greeting = await Greet("World")
// greeting = "Hello, World!"
```

**錯誤處理：**

```go
func (g *GreetService) Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

```javascript
try {
    const result = await Divide(10, 0)
} catch (error) {
    console.error("Go error:", error)  // "division by zero"
}
```

## 效能特性

### 速度

**一般呼叫開銷：**&lt;1ms

```
Frontend Call → Bridge → Go Execution → Bridge → Frontend Response
     ↓            ↓           ↓            ↓            ↓
   &lt;0.1ms      &lt;0.1ms      [varies]     &lt;0.1ms      &lt;0.1ms
```

**與其他方案相比：**

- <strong>HTTP/REST：</strong>5-50ms（網路堆疊、序列化）
- <strong>IPC：</strong>1-10ms（跨處理程序、資料編組）
- **Wails 橋接器：**&lt;1ms（記憶體內直接呼叫）

### 記憶體

<strong>每次呼叫的開銷：</strong>約 1KB（訊息緩衝區）

<strong>零複製最佳化：</strong>大型資料（>1MB）會盡可能使用共用記憶體。

### 並行處理

**呼叫會並行執行：**

- 每次呼叫都在各自的 goroutine 中執行
- 多個呼叫可同時執行
- 呼叫之間不會互相阻塞

```javascript
// These run concurrently
const [result1, result2, result3] = await Promise.all([
    SlowOperation1(),
    SlowOperation2(),
    SlowOperation3(),
])
```

## 型別系統

### 支援的型別

#### 基本型別

```go
// Go
func Example(
    s string,
    i int,
    f float64,
    b bool,
) (string, int, float64, bool) {
    return s, i, f, b
}
```

```typescript
// TypeScript (auto-generated)
function Example(
    s: string,
    i: number,
    f: number,
    b: boolean,
): Promise<[string, number, number, boolean]>
```

#### 切片與陣列

```go
// Go
func Sum(numbers []int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}
```

```typescript
// TypeScript
function Sum(numbers: number[]): Promise<number>

// Usage
const total = await Sum([1, 2, 3, 4, 5])  // 15
```

#### 映射

```go
// Go
func GetConfig() map[string]interface{} {
    return map[string]interface{}{
        "theme": "dark",
        "fontSize": 14,
        "enabled": true,
    }
}
```

```typescript
// TypeScript
function GetConfig(): Promise<Record<string, any>>

// Usage
const config = await GetConfig()
console.log(config.theme)  // "dark"
```

#### 結構體

```go
// Go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func GetUser(id int) (*User, error) {
    return &User{
        ID:    id,
        Name:  "Alice",
        Email: "alice@example.com",
    }, nil
}
```

```typescript
// TypeScript (auto-generated)
interface User {
    id: number
    name: string
    email: string
}

function GetUser(id: number): Promise<User>

// Usage
const user = await GetUser(1)
console.log(user.name)  // "Alice"
```

<strong>JSON 標籤：</strong>使用 `json:` 標籤控制 TypeScript 中的欄位名稱。

#### 時間

```go
// Go
func GetTimestamp() time.Time {
    return time.Now()
}
```

```typescript
// TypeScript
function GetTimestamp(): Promise<Date>

// Usage
const timestamp = await GetTimestamp()
console.log(timestamp.toISOString())
```

#### 錯誤

```go
// Go
func Validate(input string) error {
    if input == "" {
        return errors.New("input cannot be empty")
    }
    return nil
}
```

```typescript
// TypeScript
function Validate(input: string): Promise<void>

// Usage
try {
    await Validate("")
} catch (error) {
    console.error(error)  // "input cannot be empty"
}
```

### 不支援的型別

以下型別<strong>無法</strong>透過橋接器傳遞：

- **通道**（`chan T`）
- **函式**（`func()`）
- **介面**（`interface{}`／`any`除外）
- **指標**（指向結構的指標除外）
- **未匯出的欄位**（小寫）

<strong>替代方案：</strong>使用 ID 或控制代碼：

```go
// ❌ Can't pass file handle
func OpenFile(path string) (*os.File, error) {
    return os.Open(path)
}

// ✅ Return file ID instead
var files = make(map[string]*os.File)

func OpenFile(path string) (string, error) {
    file, err := os.Open(path)
    if err != nil {
        return "", err
    }
    id := generateID()
    files[id] = file
    return id, nil
}

func ReadFile(id string) ([]byte, error) {
    file := files[id]
    return io.ReadAll(file)
}

func CloseFile(id string) error {
    file := files[id]
    delete(files, id)
    return file.Close()
}
```

## 進階模式

### 傳遞內容脈絡

服務可以存取呼叫的內容脈絡：

```go
type UserService struct{}

func (s *UserService) GetCurrentUser(ctx context.Context) (*User, error) {
    // Access the calling window via the context value
    window, _ := ctx.Value(application.WindowKey).(application.Window)
    _ = window

    // Access the application
    app := application.Get()
    _ = app

    // Your logic
    return getCurrentUser(), nil
}
```

**內容脈絡提供：**

- 發出呼叫的視窗
- 應用程式執行個體
- 請求中繼資料

### 串流資料

對於大量資料，請使用事件而非傳回值：

```go
func ProcessLargeFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lineNum := 0

    for scanner.Scan() {
        lineNum++
        // Emit progress events
        app.Event.Emit("file-progress", map[string]interface{}{
            "line": lineNum,
            "text": scanner.Text(),
        })
    }

    return scanner.Err()
}
```

```javascript
import { Events } from '@wailsio/runtime'
import { ProcessLargeFile } from './bindings/FileService'

// Listen for progress
Events.On('file-progress', (data) => {
    console.log(`Line ${data.line}: ${data.text}`)
})

// Start processing
await ProcessLargeFile('/path/to/large/file.txt')
```

### 取消

對於可取消的操作，請使用內容脈絡：

```go
func LongRunningTask(ctx context.Context) error {
    for i := 0; i < 1000; i++ {
        // Check if cancelled
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Continue work
            time.Sleep(100 * time.Millisecond)
        }
    }
    return nil
}
```

<strong>注意：</strong>前端中斷連線時，內容脈絡會自動取消。

### 批次操作

透過批次處理來降低橋接層的額外負擔：

```go
// ❌ Inefficient: N bridge calls
for _, item := range items {
    await ProcessItem(item)
}

// ✅ Efficient: 1 bridge call
await ProcessItems(items)
```

```go
func ProcessItems(items []Item) ([]Result, error) {
    results := make([]Result, len(items))
    for i, item := range items {
        results[i] = processItem(item)
    }
    return results, nil
}
```

## 偵錯橋接層

### 啟用偵錯記錄

```go
app := application.New(application.Options{
    Name:     "My App",
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // `Logger` is an optional *slog.Logger; the default-logger helper is
    // application.DefaultLogger(slog.Leveler) if you want to construct one explicitly.
})
```

**輸出會顯示：**

- 方法呼叫
- 參數
- 傳回值
- 錯誤
- 計時資訊

### 檢查產生的繫結

查看`frontend/bindings/`以檢視產生的 TypeScript：

```javascript
// frontend/bindings/<full-go-import-path>/myservice.js (real generated shape)
import { Call as $Call } from "/wails/runtime.js";

export function MyMethod($0) {
    return $Call.ByID(1234567890, $0); // numeric method ID assigned by the generator
}
```

### 直接測試服務

不透過前端直接測試 Go 服務：

```go
func TestGreetService(t *testing.T) {
    service := &GreetService{prefix: "Hello, "}
    result := service.Greet("Test")
    if result != "Hello, Test!" {
        t.Errorf("Expected 'Hello, Test!', got '%s'", result)
    }
}
```

## 效能提示

### ✅ 建議做法

- **批次執行操作**——減少橋接層呼叫
- **使用事件進行串流傳輸**——不要傳回大型陣列
- **讓方法快速完成**——理想情況下應&lt;100ms
- **使用 goroutine**——適用於耗時操作
- **在 Go 端快取**——避免重複計算

### ❌ 請勿這樣做

- **請勿過度呼叫**——可行時採用批次處理
- **請勿傳回巨量資料**——使用分頁或串流傳輸
- **請勿阻塞**——耗時操作請使用 goroutine
- **請勿傳遞複雜型別**——保持簡單
- **請勿忽略錯誤**——務必處理錯誤

## 安全性

橋接層預設即具安全性：

1. **僅限白名單**——只能呼叫已註冊的服務
2. **型別驗證**——依照 Go 型別檢查引數
3. **不使用 eval()**——前端無法執行任意 Go 程式碼
4. **無法濫用反射**——只能存取已匯出的方法

**最佳實務：**

- 在 Go 中<strong>驗證輸入</strong>（不要信任前端）
- <strong>使用執行脈絡</strong>進行身分驗證／授權
- 對耗費大量資源的操作進行<strong>速率限制</strong>
- <strong>清理</strong>檔案路徑和使用者輸入

## 後續步驟

**建置系統**——瞭解 Wails 如何建置並封裝您的應用程式 [深入瞭解 →](/concepts/build-system/)

**服務**——深入探討服務系統 [深入瞭解 →](/features/bindings/services/)

**事件** - 使用事件進行發布／訂閱通訊 [深入瞭解 →](/features/events/system/)

---

<strong>對橋接機制有疑問嗎？</strong>請到[Discord](https://discord.gg/JDdSxwjhGf)提問，或查看[繫結範例](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)。
