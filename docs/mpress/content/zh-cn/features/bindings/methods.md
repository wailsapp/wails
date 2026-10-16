---
title: "方法绑定"
description: "以类型安全的方式从 JavaScript 调用 Go 方法"
slug: "features/bindings/methods"
sourcePath: "features/bindings/methods.md"
---

## 类型安全的 Go-JavaScript 绑定

Wails 会为你的 Go 方法<strong>自动生成类型安全的 JavaScript/TypeScript 绑定</strong>。编写 Go 代码并运行一条命令，即可获得类型完备的前端函数，无 HTTP 开销、无需手动操作，也没有任何样板代码。

## 快速开始

**1. 编写 Go 服务：**

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}
```

**2. 注册服务：**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
})
```

**3. 生成绑定：**

```bash
wails3 generate bindings
```

**4. 在 JavaScript 中使用：**

```javascript
import { Greet } from './bindings/changeme/greetservice'

const message = await Greet("World")
console.log(message)  // "Hello, World!"
```

<strong>就是这么简单！</strong>以类型安全的方式从 Go 调用 JavaScript。

## 创建服务

### 基本服务

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

**注册：**

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(&CalculatorService{}),
    },
})
```

**要点：**

- 仅绑定<strong>导出的方法</strong>（PascalCase）
- 方法可以返回值或`(value, error)`
- 服务是<strong>单例</strong>（每个应用程序一个实例）

### 有状态服务

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

<strong>重要：</strong>所有窗口共享服务。请使用互斥锁确保线程安全。

### 带依赖项的服务

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

**注册带依赖项的服务：**

```go
db, _ := sql.Open("sqlite3", "app.db")

app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewDatabaseService(db)),
    },
})
```

## 生成绑定

### 基本生成方式

```bash
wails3 generate bindings
```

**输出：**

```
INFO  347 Packages, 3 Services, 12 Methods, 0 Enums, 0 Models in 1.98s
INFO  Output directory: /myproject/frontend/bindings
```

**生成的结构：**

@filetree
- frontend/bindings
  - myapp
    - calculatorservice.js
    - counterservice.js
    - databaseservice.js
    - index.js
@end

### 生成 TypeScript 绑定

```bash
wails3 generate bindings -ts
```

**生成包含完整 TypeScript 类型的`.ts`文件**。

### 自定义输出目录

```bash
wails3 generate bindings -d ./src/bindings
```

### 监视模式（开发）

```bash
wails3 dev
```

Go 代码发生变化时，**自动重新生成绑定**。

## 使用绑定

### JavaScript

**生成的绑定：**

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

**生成的绑定：**

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

**优势：**

- 完整的类型检查
- IDE 自动补全
- 编译时错误检查
- 更易于重构

### 索引文件

**生成的索引：**

```javascript
// frontend/bindings/changeme/index.js

export * as CalculatorService from './calculatorservice.js'
export * as CounterService from './counterservice.js'
export * as DatabaseService from './databaseservice.js'
```

**简化的导入方式：**

```javascript
import { CalculatorService } from './bindings/myapp'

const sum = await CalculatorService.Add(5, 3)
```

## 类型映射

### 基本类型

| Go 类型 | JavaScript/TypeScript |
| --- | --- |
| `string` | `string` |
| `bool` | `boolean` |
| `int`、`int8`、`int16`、`int32`、`int64` | `number` |
| `uint`、`uint8`、`uint16`、`uint32`、`uint64` | `number` |
| `float32`、`float64` | `number` |
| `byte` | `number` |
| `rune` | `number` |

### 复杂类型

| Go 类型 | JavaScript/TypeScript | 说明 |
| --- | --- | --- |
| `[]T` | `T[]` | - |
| `[N]T` | `T[]` | - |
| `map[string]T` | `{ [_: string]: T }` | 以字符串为键的映射 |
| `map[K]V` | `{ [_ in K]?: V }` | 非字符串类型的`K`会呈现为映射类型，<strong>而不是</strong>JS `Map` |
| `[]byte` | `string` | Base64 编码 |
| `struct` | `class` / `interface` | 包含字段 |
| `time.Time` | `any` | 在运行时序列化为 RFC3339Nano 字符串 |
| `*T` | `T \| null` | 指针表示可为 null |
| `any` / `interface{}` | `any` | - |
| `error` | `any` / `Exception` | 作为返回值时为异常，否则为 any |

### 不支持的类型

以下类型<strong>不能</strong>通过桥接层传递：

- `chan T`（通道）
- `func()`（函数）
- 复杂接口（`interface{}`除外）
- 未导出的字段（小写）

<strong>解决方法：</strong>使用 ID 或句柄：

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

## 错误处理

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

绑定的方法失败时，返回的 Promise 会以 JavaScript `Error`对象拒绝：

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

运行时会根据具体的错误原因抛出不同类型的错误：

| 错误类型 | 抛出条件 |
| --- | --- |
| `TypeError` | 调用的参数数量错误，或某个参数无法转换为相应的 Go 类型 |
| `RuntimeError` | 方法返回了错误，或在运行时发生 panic |
| `Error` | 任何其他故障，例如调用了不存在的方法 |

每个错误都提供：

- `name`：上表中的错误类型
- `message`：Go 错误的消息
- `cause`：Go 错误的 JSON 序列化形式（如有）。如果方法返回了多个错误，`cause`将是一个数组，每个错误对应一个元素。

`RuntimeError`类由`@wailsio/runtime`包导出，因此可以区分 Go 代码返回的错误与其他故障：

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
在较早版本的运行时中，每次失败的调用都会以普通的`Error`拒绝，其消息包含底层错误的原始 JSON。

@end

### 结构化错误数据

从 Go 返回自定义错误类型后，可通过所抛出错误的`cause`属性访问其 JSON 形式：

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

错误使用标准`encoding/json`包进行序列化，因此只包含已导出的字段。使用`errors.New`或`fmt.Errorf`创建的错误没有已导出字段，会序列化为空对象。

如需完全控制错误的序列化方式，请在服务选项中提供`MarshalError`函数：

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

`MarshalError`必须返回有效的 JSON，或返回`nil`以回退到默认序列化方式。

## 性能

### 调用开销

**典型调用：**&lt;1ms

```
JavaScript → Bridge → Go → Bridge → JavaScript
    ↓           ↓       ↓       ↓           ↓
  &lt;0.1ms    &lt;0.1ms  [varies] &lt;0.1ms    &lt;0.1ms
```

**与其他方案相比：**

- HTTP/REST：5-50ms
- IPC：1-10ms
- Wails：&lt;1ms

### 优化技巧

**✅ 批量操作：**

```javascript
// ❌ Slow: N calls
for (const item of items) {
    await ProcessItem(item)
}

// ✅ Fast: 1 call
await ProcessItems(items)
```

**✅ 缓存结果：**

```javascript
// ❌ Repeated calls
const config1 = await GetConfig()
const config2 = await GetConfig()

// ✅ Cache
const config = await GetConfig()
// Use config multiple times
```

**✅ 使用事件传输流式数据：**

```go
func ProcessLargeFile(path string) error {
    // Emit progress events
    for line := range lines {
        app.Event.Emit("progress", line)
    }
    return nil
}
```

## 完整示例

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

## 最佳实践

### ✅ 应该做

- **保持方法简单**——遵循单一职责原则
- **返回错误**——不要触发 panic
- **使用线程安全的状态**——用互斥锁保护共享数据
- **批量执行操作**——减少跨桥调用
- **在 Go 端缓存**——避免重复工作
- **为方法编写文档**——注释会转换为 JSDoc

### ❌ 不应该做

- **不要阻塞**——长时间运行的操作应使用 goroutine
- **不要返回通道**——改用事件
- **不要返回函数**——不受支持
- **不要忽略错误**——始终处理错误
- **不要使用未导出的字段**——这些字段不会生成绑定

## 后续步骤

@cards{cols="2"}
◆ 服务
深入了解服务系统。

[了解更多 →](/features/bindings/services/)

---
📖 模型
绑定复杂的数据结构。

[了解更多 →](/features/bindings/models/)

---
🚀 Go-前端桥接
了解桥接机制。

[了解更多 →](/concepts/bridge/)

---
★ 事件
使用事件进行发布/订阅通信。

[了解更多 →](/features/events/system/)

@end

---

<strong>有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[绑定示例](https://github.com/wailsapp/wails/tree/master/v3/examples/binding)。
