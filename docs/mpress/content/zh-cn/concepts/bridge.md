---
title: "Go-前端桥接"
description: "深入了解 Wails 如何实现 Go 与 JavaScript 之间的直接通信"
slug: "concepts/bridge"
sourcePath: "concepts/bridge.md"
---

## Go-JavaScript 直接通信

Wails 在 Go 与 JavaScript 之间提供**直接的内存桥接**，无需承担 HTTP 开销、跨越进程边界或遭遇序列化瓶颈，即可实现无缝通信。

## 整体架构

```d2
direction: right

Frontend: 前端（JavaScript） {
  UI: React/Vue/原生 JavaScript {
    shape: rectangle
    style.fill: "#8B5CF6"
  }

  Bindings: 自动生成的绑定 {
    shape: rectangle
    style.fill: "#A78BFA"
  }
}

Bridge: Wails 桥接 {
  Encoder: JSON 编码器 {
    shape: rectangle
    style.fill: "#10B981"
  }

  Router: 方法路由器 {
    shape: diamond
    style.fill: "#10B981"
  }

  Decoder: JSON 解码器 {
    shape: rectangle
    style.fill: "#10B981"
  }

  TypeGen: 类型生成器 {
    shape: rectangle
    style.fill: "#10B981"
  }
}

Backend: 后端（Go） {
  Services: 你的服务 {
    shape: rectangle
    style.fill: "#00ADD8"
  }

  Registry: 服务注册表 {
    shape: rectangle
    style.fill: "#00ADD8"
  }
}

Frontend.UI -> Frontend.Bindings: "import { Method }"
Frontend.Bindings -> Bridge.Encoder: "调用 Method('arg')"
Bridge.Encoder -> Bridge.Router: 编码为 JSON
Bridge.Router -> Backend.Registry: 查找服务
Backend.Registry -> Backend.Services: 调用方法
Backend.Services -> Bridge.Decoder: 返回结果
Bridge.Decoder -> Frontend.Bindings: 解码为 JS
Frontend.Bindings -> Frontend.UI: Promise 完成
Bridge.TypeGen -> Frontend.Bindings: 生成类型
```

**关键要点：**没有 HTTP、IPC，也没有进程边界，只有具备**类型安全**的**直接函数调用**。

## 工作原理：分步说明

### 1. 服务注册（启动时）

应用程序启动时，Wails 会扫描你的服务：

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

**Wails 执行的操作：**

1. **扫描结构体**中的已导出方法
2. **提取类型信息**（参数和返回类型）
3. **构建注册表**，将方法名称映射到函数
4. **生成 TypeScript 绑定**，其中包含完整的类型定义

### 2. 生成绑定（构建时）

Wails 会自动生成 TypeScript 绑定：

```typescript
// Auto-generated: frontend/bindings/GreetService.ts
export function Greet(name: string): Promise<string>
export function Add(a: number, b: number): Promise<number>
```

**类型映射：**

| Go 类型 | TypeScript 类型 |
| --- | --- |
| `string` | `string` |
| `int`, `int32`, `int64` | `number` |
| `float32`, `float64` | `number` |
| `bool` | `boolean` |
| `[]T` | `T[]` |
| `map[string]T` | `Record<string, T>` |
| `struct` | `interface` |
| `time.Time` | 默认为`string`；搭配`-time-type=Date`时为`Date` |
| `error` | 异常（抛出） |

### 3. 前端调用（运行时）

开发者从 JavaScript 调用 Go 方法：

```javascript
import { Greet, Add } from './bindings/GreetService'

// Call Go from JavaScript
const greeting = await Greet("World")
console.log(greeting)  // "Hello, World!"

const sum = await Add(5, 3)
console.log(sum)  // 8
```

**执行过程：**

1. **调用绑定函数** — `Greet("World")`
2. **创建消息** — `{ service: "GreetService", method: "Greet", args: ["World"] }`
3. **发送至桥接层** — 通过 WebView 的 JavaScript 桥接
4. **返回 Promise** — 等待响应

### 4. 桥接处理（运行时）

桥接层接收并处理消息：

```d2
direction: down

Receive: 接收消息 {
  shape: rectangle
  style.fill: "#10B981"
}

Parse: 解析 JSON {
  shape: rectangle
}

Validate: 验证 {
  Check: 服务是否存在？ {
    shape: diamond
  }

  CheckMethod: 方法是否存在？ {
    shape: diamond
  }

  CheckTypes: 类型是否正确？ {
    shape: diamond
  }
}

Invoke: 调用 Go 方法 {
  shape: rectangle
  style.fill: "#00ADD8"
}

Encode: 编码结果 {
  shape: rectangle
}

Send: 发送响应 {
  shape: rectangle
  style.fill: "#10B981"
}

Error: 发送错误 {
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
Invoke -> Error: 错误
Encode -> Send
```

<strong>安全性：</strong>只能调用已注册的服务和已导出的方法。

### 5. Go 执行（运行时）

Go 方法开始执行：

```go
func (g *GreetService) Greet(name string) string {
    // This runs in Go
    return g.prefix + name + "!"
}
```

**执行上下文：**

- 在**goroutine**中运行（非阻塞）
- 可以使用**Go 的全部功能**（文件系统、网络和数据库）
- 可以自由调用**其他 Go 代码**
- 返回结果或错误

### 6. 响应（运行时）

结果将发回 JavaScript：

```javascript
// Promise resolves with result
const greeting = await Greet("World")
// greeting = "Hello, World!"
```

**错误处理：**

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

## 性能特征

### 速度

**典型调用开销：**&lt;1ms

```
Frontend Call → Bridge → Go Execution → Bridge → Frontend Response
     ↓            ↓           ↓            ↓            ↓
   &lt;0.1ms      &lt;0.1ms      [varies]     &lt;0.1ms      &lt;0.1ms
```

**与其他方案相比：**

- <strong>HTTP/REST：</strong>5-50ms（网络栈、序列化）
- <strong>IPC：</strong>1-10毫秒（跨进程、编组）
- **Wails 桥接：**&lt;1毫秒（内存中直接调用）

### 内存

<strong>每次调用的开销：</strong>约1KB（消息缓冲区）

<strong>零拷贝优化：</strong>对于大型数据（>1MB），会尽可能使用共享内存。

### 并发

**调用可并发执行：**

- 每次调用都在各自的 goroutine 中运行
- 多个调用可以同时执行
- 调用之间不会相互阻塞

```javascript
// These run concurrently
const [result1, result2, result3] = await Promise.all([
    SlowOperation1(),
    SlowOperation2(),
    SlowOperation3(),
])
```

## 类型系统

### 支持的类型

#### 基本类型

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

#### 切片和数组

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

#### 结构体

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

<strong>JSON 标签：</strong>使用`json:`标签控制 TypeScript 中的字段名称。

#### 时间

```go
// Go
func GetTimestamp() time.Time {
    return time.Now()
}
```

```typescript
// TypeScript
function GetTimestamp(): Promise<string> // default -time-type=string

// Usage
const timestamp = await GetTimestamp()
console.log(new Date(timestamp).toISOString())
```

#### 错误

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

### 不支持的类型

以下类型**无法**通过桥接传递：

- **通道**（`chan T`）
- **函数**（`func()`）
- **接口**（`interface{}`/`any`除外）
- **指针**（指向结构体的指针除外）
- **未导出字段**（小写字母开头）

<strong>解决方法：</strong>使用 ID 或句柄：

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

## 高级模式

### 上下文传递

服务可以访问调用上下文：

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

**上下文提供：**

- 发起调用的窗口
- 应用实例
- 请求元数据

### 流式数据

对于大型数据，请使用事件而非返回值：

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

### 取消操作

对于可取消的操作，请使用上下文：

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

<strong>注意：</strong>前端断开连接时会自动取消上下文。

### 批量操作

通过批量处理减少桥接开销：

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

## 调试桥接

### 启用调试日志

```go
app := application.New(application.Options{
    Name:     "My App",
    LogLevel: slog.LevelDebug, // requires `import "log/slog"`
    // `Logger` is an optional *slog.Logger; the default-logger helper is
    // application.DefaultLogger(slog.Leveler) if you want to construct one explicitly.
})
```

**输出内容包括：**

- 方法调用
- 参数
- 返回值
- 错误
- 耗时信息

### 检查生成的绑定

查看`frontend/bindings/`中的已生成 TypeScript：

```javascript
// frontend/bindings/<full-go-import-path>/myservice.js (real generated shape)
import { Call as $Call } from "/wails/runtime.js";

export function MyMethod($0) {
    return $Call.ByID(1234567890, $0); // numeric method ID assigned by the generator
}
```

### 直接测试服务

在不使用前端的情况下测试 Go 服务：

```go
func TestGreetService(t *testing.T) {
    service := &GreetService{prefix: "Hello, "}
    result := service.Greet("Test")
    if result != "Hello, Test!" {
        t.Errorf("Expected 'Hello, Test!', got '%s'", result)
    }
}
```

## 性能建议

### ✅ 推荐做法

- **批量处理操作**——减少桥接调用
- **使用事件传输流式数据**——不要返回大型数组
- **确保方法快速执行**——理想耗时为&lt;100毫秒
- **使用 goroutine**——用于耗时较长的操作
- **在 Go 端缓存**——避免重复计算

### ❌ 请勿这样做

- **不要进行过多调用**——尽可能批量处理
- **不要返回海量数据**——使用分页或流式传输
- **不要阻塞**——对耗时较长的操作使用 goroutine
- **不要传递复杂类型**——保持简单
- **不要忽略错误**——务必处理错误

## 安全性

桥接默认安全：

1. **仅限白名单**——只能调用已注册的服务
2. **类型验证**——根据 Go 类型检查参数
3. **不使用 eval()**——前端无法执行任意 Go 代码
4. **不会滥用反射**——只能访问已导出的方法

**最佳实践：**

- 在 Go 中**验证输入**（不要信任前端）
- 使用**上下文**进行身份验证和授权
- 对高开销操作进行**速率限制**
- 对文件路径和用户输入进行**清理**

## 后续步骤

**构建系统**——了解 Wails 如何构建应用并将其打包[了解更多 →](/concepts/build-system/)

**服务**——深入了解服务系统[了解更多 →](/features/bindings/services/)

**事件**——使用事件进行发布/订阅通信[了解更多 →](/features/events/system/)

---

<strong>对桥接有疑问？</strong>请在[Discord](https://discord.gg/JDdSxwjhGf)中提问，或查看[绑定示例](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples/binding)。
