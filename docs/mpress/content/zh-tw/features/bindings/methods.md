---
title: "方法繫結"
description: "以型別安全的方式從 JavaScript 呼叫 Go 方法"
slug: "features/bindings/methods"
sourcePath: "features/bindings/methods.md"
---

## 型別安全的 Go-JavaScript 繫結

Wails 會<strong>自動為您的 Go 方法產生型別安全的 JavaScript/TypeScript 繫結</strong>。編寫 Go 程式碼、執行一個命令，即可取得具有完整型別的前端函式，無 HTTP 額外負擔、無須手動處理，也不需要任何樣板程式碼。

## 快速開始

**1. 編寫 Go 服務：**

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}
```

**2. 註冊服務：**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

**3. 產生繫結：**

```bash
wails3 generate bindings
```

**4. 在 JavaScript 中使用：**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const message = await Greet("World")
console.log(message)  // "Hello, World!"
```

<strong>就這樣！</strong>以型別安全的方式從 Go 呼叫 JavaScript。

## 建立服務

### 基本服務

```go
package main

import "github.com/wailsapp/wails/v3/pkg/application"

type CalculatorService struct{}

func (c *CalculatorService) Add(a, b int) int {
    return a + b
}

func (c *CalculatorService) Subtract(a, b int) int {
    return a - b
}

func (c *CalculatorService) Multiply(a, b int) int {
    return a * b
}

func (c *CalculatorService) Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

**註冊：**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&CalculatorService{}),
    },
})
```

**要點：**

- 僅繫結<strong>匯出的方法</strong>（PascalCase）
- 方法可以傳回值或`(value, error)`
- 服務是<strong>單例</strong>（每個應用程式一個執行個體）

### 具狀態的服務

```go
type CounterService struct {
    count int
    mu    sync.Mutex
}

func (c *CounterService) Increment() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
    return c.count
}

func (c *CounterService) Decrement() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count--
    return c.count
}

func (c *CounterService) GetCount() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.count
}

func (c *CounterService) Reset() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count = 0
}
```

<strong>重要：</strong>所有視窗會共用服務。請使用互斥鎖來確保執行緒安全。

### 具有相依性的服務

```go
type DatabaseService struct {
    db *sql.DB
}

func NewDatabaseService(db *sql.DB) *DatabaseService {
    return &DatabaseService{db: db}
}

func (d *DatabaseService) GetUser(id int) (*User, error) {
    var user User
    err := d.db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user)
    return &user, err
}
```

**連同相依性一起註冊：**

```go
db, _ := sql.Open("sqlite3", "app.db")

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewDatabaseService(db)),
    },
})
```

## 產生繫結

### 基本產生方式

```bash
wails3 generate bindings
```

**輸出：**

```
INFO  347 Packages, 3 Services, 12 Methods, 0 Enums, 0 Models in 1.98s
INFO  Output directory: /myproject/frontend/bindings
```

**產生的結構：**

@filetree
- frontend/bindings
  - myapp
    - calculatorservice.js
    - counterservice.js
    - databaseservice.js
    - index.js
@end

### 產生 TypeScript 繫結

```bash
wails3 generate bindings -ts
```

**產生具有完整 TypeScript 型別的`.ts`檔案**。

### 自訂輸出目錄

```bash
wails3 generate bindings -d ./src/bindings
```

### 監看模式（開發）

```bash
wails3 dev
```

Go 程式碼變更時，**自動重新產生繫結**。

## 使用繫結

### JavaScript

**產生的繫結：**

```javascript
// frontend/bindings/<full-go-import-path>/calculatorservice.js
// (Real generated output — imports $Call from /wails/runtime.js and calls $Call.ByID
// with a numeric method ID. Generate with `wails3 generate bindings -names` to get
// $Call.ByName("<package>.<Struct>.<Method>", ...) instead.)
import { Call as $Call, Create as $Create } from "/wails/runtime.js";

/**
 * @param {number} $0
 * @param {number} $1
 * @returns {Promise<number>}
 */
export function Add($0, $1) {
    return $Call.ByID(1234567890, $0, $1); // numeric ID assigned by the generator
}
```

**用法：**

```javascript
import { Add, Subtract, Multiply, Divide } from './bindings/changeme/calculatorservice'

// Simple calls
const sum = await Add(5, 3)        // 8
const diff = await Subtract(10, 4)  // 6
const product = await Multiply(7, 6) // 42

// Error handling
try {
    const result = await Divide(10, 0)
} catch (error) {
    console.error("Error:", error)  // "division by zero"
}
```

### TypeScript

**產生的繫結：**

```typescript
// frontend/bindings/changeme/calculatorservice.ts

export function Add(a: number, b: number): Promise<number>
export function Subtract(a: number, b: number): Promise<number>
export function Multiply(a: number, b: number): Promise<number>
export function Divide(a: number, b: number): Promise<number>
```

**用法：**

```typescript
import { Add, Divide } from './bindings/changeme/calculatorservice'

const sum: number = await Add(5, 3)

try {
    const result = await Divide(10, 0)
} catch (error: unknown) {
    if (error instanceof Error) {
        console.error(error.message)
    }
}
```

**優點：**

- 完整的型別檢查
- IDE 自動完成
- 編譯時期錯誤
- 更易於重構

### 索引檔案

**產生的索引：**

```javascript
// frontend/bindings/changeme/index.js

export * as CalculatorService from './calculatorservice.js'
export * as CounterService from './counterservice.js'
export * as DatabaseService from './databaseservice.js'
```

**簡化的匯入方式：**

```javascript
import { CalculatorService } from './bindings/myapp'

const sum = await CalculatorService.Add(5, 3)
```

## 型別對應

### 基本型別

| Go 型別 | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`、`int8`、`int16`、`int32`、`int64` | `number` |
| `uint`、`uint8`、`uint16`、`uint32`、`uint64` | `number` |
| `float32`、`float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### 複合型別

| Go 型別 | JavaScript/TypeScript | 備註 |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | 以字串為鍵的映射 |
| `map[K]V` | `{ [_ in K]?: V }` | 非字串的`K`會轉譯為映射型別，<strong>而不是</strong>JS 的`Map` |
| `[]byte` | `string` | 以 base64 編碼 |
| `struct` | `class` / `interface` | 包含欄位 |
| `time.Time` | `any` | 執行階段會序列化為 RFC3339Nano 字串 |
| `*T` | `T \| null` | 指標表示可為 null |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | 若作為回傳值則為例外，否則為 any |

### 不支援的型別

以下型別<strong>無法</strong>透過橋接傳遞：

- `chan T`（通道）
- `func()`（函式）
- 複雜介面（`interface{}`除外）
- 未匯出的欄位（小寫）

<strong>因應方法：</strong>使用 ID 或控制代碼：

```go
// ❌ Can't return file handle
func OpenFile(path string) (*os.File, error)

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

## 錯誤處理

### Go 端

```go
func (d *DatabaseService) GetUser(id int) (*User, error) {
    if id <= 0 {
        return nil, errors.New("invalid user ID")
    }
    
    var user User
    err := d.db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user %d not found", id)
    }
    if err != nil {
        return nil, fmt.Errorf("database error: %w", err)
    }
    
    return &user, nil
}
```

### JavaScript 端

繫結的方法失敗時，回傳的 Promise 會以 JavaScript `Error`物件拒絕：

```javascript
import { GetUser } from './bindings/changeme/databaseservice'

try {
    const user = await GetUser(123)
    console.log("User:", user)
} catch (error) {
    console.error(error.message)
    // "user 123 not found"
}
```

執行階段會根據發生的問題擲回不同的錯誤型別：

| 錯誤型別 | 擲回時機 |
| --- | --- |
| `TypeError` | 呼叫所傳入的引數數量不正確，或某個引數無法轉換為 Go 型別 |
| `RuntimeError` | 方法回傳錯誤，或在執行期間發生 panic |
| `Error` | 任何其他失敗，例如呼叫不存在的方法 |

每個錯誤都會提供：

- `name`：上表中的錯誤型別
- `message`：Go 錯誤的訊息
- `cause`：Go 錯誤序列化後的 JSON（若可取得）。若方法回傳多個錯誤，`cause`會是陣列，每個錯誤各有一個項目。

`RuntimeError`類別由`@wailsio/runtime`套件匯出，因此你可以區分 Go 程式碼回傳的錯誤與其他失敗：

```javascript
import { Call } from '@wailsio/runtime'
import { GetUser } from './bindings/changeme/databaseservice'

try {
    const user = await GetUser(123)
} catch (error) {
    if (error instanceof Call.RuntimeError) {
        // GetUser returned an error
    } else {
        // The call itself failed
    }
}
```

@note{type="info"}
在較舊版本的執行階段中，每次失敗的呼叫都會以一般的`Error`拒絕，其訊息包含底層錯誤的原始 JSON。

@end

### 結構化錯誤資料

從 Go 回傳自訂錯誤型別後，該錯誤的 JSON 形式可透過所擲回錯誤的`cause`屬性取得：

```go
type ValidationError struct {
    Field  string `json:"field"`
    Reason string `json:"reason"`
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Reason)
}

func (s *UserService) UpdateEmail(id int, email string) error {
    if !strings.Contains(email, "@") {
        return &ValidationError{Field: "email", Reason: "invalid email address"}
    }
    // ...
    return nil
}
```

```javascript
import { UpdateEmail } from './bindings/changeme/userservice'

try {
    await UpdateEmail(1, "not-an-email")
} catch (error) {
    console.log(error.message)  // "email: invalid email address"
    console.log(error.cause)    // { field: "email", reason: "invalid email address" }
}
```

錯誤會使用標準`encoding/json`套件序列化，因此只會包含已匯出的欄位。使用`errors.New`或`fmt.Errorf`建立的錯誤沒有已匯出的欄位，會序列化為空物件。

若要完全控制錯誤的序列化方式，請在服務選項中提供`MarshalError`函式：

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewServiceWithOptions(&UserService{}, application.ServiceOptions{
            MarshalError: func(err error) []byte {
                var validationErr *ValidationError
                if errors.As(err, &validationErr) {
                    data, _ := json.Marshal(map[string]string{
                        "type":  "validation",
                        "field": validationErr.Field,
                    })
                    return data
                }
                return nil // fall back to the default serialisation
            },
        }),
    },
})
```

`MarshalError`必須回傳有效的 JSON，或回傳`nil`以改用預設序列化。

## 效能

### 呼叫額外負荷

**一般呼叫：**&lt;1ms

```
JavaScript → Bridge → Go → Bridge → JavaScript
    ↓           ↓       ↓       ↓           ↓
  &lt;0.1ms    &lt;0.1ms  [varies] &lt;0.1ms    &lt;0.1ms
```

**與其他方式比較：**

- HTTP/REST：5-50ms
- IPC：1-10ms
- Wails：&lt;1ms

### 最佳化技巧

**✅ 批次執行操作：**

```javascript
// ❌ Slow: N calls
for (const item of items) {
    await ProcessItem(item)
}

// ✅ Fast: 1 call
await ProcessItems(items)
```

**✅ 快取結果：**

```javascript
// ❌ Repeated calls
const config1 = await GetConfig()
const config2 = await GetConfig()

// ✅ Cache
const config = await GetConfig()
// Use config multiple times
```

**✅ 使用事件串流傳輸資料：**

```go
func ProcessLargeFile(path string) error {
    // Emit progress events
    for line := range lines {
        app.Event.Emit("progress", line)
    }
    return nil
}
```

## 完整範例

**Go：**

```go
package main

import (
    "fmt"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type TodoService struct {
    todos []Todo
}

type Todo struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

func (t *TodoService) GetAll() []Todo {
    return t.todos
}

func (t *TodoService) Add(title string) Todo {
    todo := Todo{
        ID:        len(t.todos) + 1,
        Title:     title,
        Completed: false,
    }
    t.todos = append(t.todos, todo)
    return todo
}

func (t *TodoService) Toggle(id int) error {
    for i := range t.todos {
        if t.todos[i].ID == id {
            t.todos[i].Completed = !t.todos[i].Completed
            return nil
        }
    }
    return fmt.Errorf("todo %d not found", id)
}

func (t *TodoService) Delete(id int) error {
    for i := range t.todos {
        if t.todos[i].ID == id {
            t.todos = append(t.todos[:i], t.todos[i+1:]...)
            return nil
        }
    }
    return fmt.Errorf("todo %d not found", id)
}

func main() {
    app := application.New(application.Options{
        Services: []application.Service{
            application.NewService(&TodoService{}),
        },
    })
    
    app.Window.New()
    app.Run()
}
```

**JavaScript：**

```javascript
import { GetAll, Add, Toggle, Delete } from './bindings/changeme/todoservice'

class TodoApp {
    async loadTodos() {
        const todos = await GetAll()
        this.renderTodos(todos)
    }
    
    async addTodo(title) {
        try {
            const todo = await Add(title)
            this.loadTodos()
        } catch (error) {
            console.error("Failed to add todo:", error)
        }
    }
    
    async toggleTodo(id) {
        try {
            await Toggle(id)
            this.loadTodos()
        } catch (error) {
            console.error("Failed to toggle todo:", error)
        }
    }
    
    async deleteTodo(id) {
        try {
            await Delete(id)
            this.loadTodos()
        } catch (error) {
            console.error("Failed to delete todo:", error)
        }
    }
    
    renderTodos(todos) {
        const list = document.getElementById('todo-list')
        list.innerHTML = todos.map(todo => `
            <div class="todo ${todo.Completed ? 'completed' : ''}">
                <input type="checkbox" 
                       ${todo.Completed ? 'checked' : ''}
                       onchange="app.toggleTodo(${todo.ID})">
                <span>${todo.Title}</span>
                <button onclick="app.deleteTodo(${todo.ID})">Delete</button>
            </div>
        `).join('')
    }
}

const app = new TodoApp()
app.loadTodos()
```

## 最佳實務

### ✅ 建議做法

- **保持方法簡潔**——每個方法只負責一項工作
- **傳回錯誤**——不要引發 panic
- **使用執行緒安全的狀態**——以互斥鎖保護共用資料
- **批次執行操作**——減少跨橋接層的呼叫
- **在 Go 端快取**——避免重複處理
- **為方法撰寫說明**——註解會轉換為 JSDoc

### ❌ 避免事項

- **不要阻塞**——長時間執行的操作應使用 goroutine
- **不要傳回 channel**——改用事件
- **不要傳回函式**——不支援此做法
- **不要忽略錯誤**——一律處理錯誤
- **不要使用未匯出的欄位**——這些欄位不會產生繫結

## 後續步驟

@cards{cols="2"}
◆ 服務
深入瞭解服務系統。

[深入瞭解 →](/features/bindings/services/)

---
📖 模型
繫結複雜的資料結構。

[深入瞭解 →](/features/bindings/models/)

---
🚀 Go 與前端橋接
瞭解橋接機制。

[深入瞭解 →](/concepts/bridge/)

---
★ 事件
使用事件進行發布／訂閱通訊。

[深入瞭解 →](/features/events/system/)

@end

---

<strong>有問題嗎？</strong>請到[Discord](https://discord.gg/JDdSxwjhGf)提問，或查看[繫結範例](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)。
