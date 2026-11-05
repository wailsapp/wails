---
title: "原始消息"
description: "为性能关键型应用实现从前端到后端的自定义通信"
slug: "guides/raw-messages"
sourcePath: "guides/raw-messages.md"
---

原始消息绕过标准绑定系统，在前端与后端之间提供底层通信通道。这以便利性换取速度。

## 何时使用原始消息

原始消息最适合以下极端边缘场景：

- **超高频更新**——每秒传输数千条消息，且每一微秒都至关重要
- **自定义消息协议**——需要完全控制消息传输格式时

@note{type="tip"}
对于几乎所有使用场景，建议使用标准[服务绑定](/features/bindings/services/)，因为它们只产生可忽略不计的开销，同时提供类型安全、自动序列化和更好的开发体验。

@end

## 后端设置

在应用选项中配置`RawMessageHandler`：

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            fmt.Printf("Raw message from window '%s': %s (origin: %+v)\n", window.Name(), message, originInfo.Origin)

            // Process the message and respond via events
            response := processMessage(message)
            window.EmitEvent("raw-response", response)
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "My App",
        Name:  "main",
    })

    app.Run()
}

func processMessage(message string) map[string]any {
    // Your custom message processing logic
    return map[string]any{
        "received": message,
        "status":   "processed",
    }
}
```

### 处理程序签名

```go
RawMessageHandler func(window Window, message string, originInfo *application.OriginInfo)
```

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| `window` | `Window` | 发送消息的窗口 |
| `message` | `string` | 原始消息内容 |
| `originInfo` | `*application.OriginInfo` | 关于消息来源的源信息 |

#### OriginInfo 结构

```go
type OriginInfo struct {
	Origin      string
	TopOrigin   string
	IsMainFrame bool
}
```

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `Origin` | `string` | 发送消息的文档的源 URL |
| `TopOrigin` | `string` | 顶层源 URL（在 iframe 中可能与 Origin 不同） |
| `IsMainFrame` | `bool` | 消息是否来自主框架 |

#### 特定平台的可用性

- **macOS**：提供`Origin`和`IsMainFrame`
- **Windows**：提供`Origin`和`TopOrigin`
- **Linux**：仅提供`Origin`

### 源验证

@note{type="caution"}
绝不能因为消息到达了处理程序，就认定它是安全的。在处理敏感操作或会修改状态的操作之前，必须验证源信息。

@end

**处理传入消息前，始终验证其来源。**`originInfo`参数提供关键安全信息，必须对其进行验证，以防止未经授权的访问。 恶意内容、遭到入侵的内容或非预期脚本都可能发送原始消息。如果不验证来源，您可能会处理来自不可信来源的命令。使用`originInfo`确保消息来自预期来源。

### 关键验证要点

- **始终检查`Origin`**——验证来源是否与预期的可信来源相符（对于本地资源，通常为`wails://wails`或`http://wails.localhost`；也可以是应用的特定来源）
- **验证`IsMainFrame`**（macOS）——注意消息是否来自 iframe，因为这可能表示其中嵌入了具有不同安全上下文的内容
- **使用`TopOrigin`**（Windows）——处理框架中的内容时，验证顶层来源
- **拒绝非预期来源**——拒绝来自未明确允许的来源的消息，以安全方式失败

@note{type="info"}
以`wails:`为前缀的消息保留供 Wails 内部通信使用，不会传递给您的处理程序。

@end

## 前端设置

使用`System.invoke()`发送原始消息：

```html
<!DOCTYPE html>
<html>
<head>
    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        // Send raw message
        document.getElementById('send').addEventListener('click', () => {
            const message = document.getElementById('input').value
            System.invoke(message)
        })

        // Listen for response
        Events.On('raw-response', (event) => {
            console.log('Response:', event.data)
        })
    </script>
</head>
<body>
    <input type="text" id="input" placeholder="Enter message" />
    <button id="send">Send</button>
</body>
</html>
```

### 使用预构建包

如果不使用 npm，请通过全局`wails`对象访问`invoke`：

```html
<script type="module" src="/wails/runtime.js"></script>
<script>
    window.onload = function() {
        document.getElementById('send').onclick = function() {
            wails.System.invoke('my-message')
        }
    }
</script>
```

## 结构化消息

对于复杂数据，请将其序列化为 JSON：

### 前端

```javascript
import { System } from '@wailsio/runtime'

const command = {
    action: 'update',
    payload: {
        id: 123,
        value: 'new value'
    }
}

System.invoke(JSON.stringify(command))
```

### 后端

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    var cmd struct {
        Action  string `json:"action"`
        Payload struct {
            ID    int    `json:"id"`
            Value string `json:"value"`
        } `json:"payload"`
    }

    if err := json.Unmarshal([]byte(message), &cmd); err != nil {
        window.EmitEvent("error", err.Error())
        return
    }

    switch cmd.Action {
    case "update":
        // Handle update
        result := handleUpdate(cmd.Payload.ID, cmd.Payload.Value)
        window.EmitEvent("update-complete", result)
    default:
        window.EmitEvent("error", "unknown action")
    }
}
```

## 性能比较

| 方式 | 开销 | 类型安全性 | 用例 |
| --- | --- | --- | --- |
| 服务绑定 | 较高 | 完整 | 通用 |
| 原始消息 | 极低 | 手动 | 高频、性能关键型 |

### 基准测试示例

对于简单的有效载荷，原始消息每秒可处理的消息数量远高于服务绑定：

```go
// Raw message handler - minimal overhead
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Direct string processing, no reflection or marshaling
    counter++
}
```

## 完整示例

以下是实现简单命令协议的完整示例：

### main.go

```go
package main

import (
    "embed"
    "encoding/json"
    "fmt"
    "time"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

type Command struct {
    Type string          `json:"type"`
    Data json.RawMessage `json:"data"`
}

func main() {
    app := application.New(application.Options{
        Name: "Raw Message Demo",
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        Mac: application.MacOptions{
            ApplicationShouldTerminateAfterLastWindowClosed: true,
        },
        RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
            var cmd Command
            if err := json.Unmarshal([]byte(message), &cmd); err != nil {
                window.EmitEvent("error", map[string]string{"error": err.Error()})
                return
            }

            switch cmd.Type {
            case "ping":
                window.EmitEvent("pong", map[string]any{
                    "time":   time.Now().UnixMilli(),
                    "window": window.Name(),
                })
            case "echo":
                var text string
                json.Unmarshal(cmd.Data, &text)
                window.EmitEvent("echo", text)
            default:
                window.EmitEvent("error", map[string]string{
                    "error": fmt.Sprintf("unknown command: %s", cmd.Type),
                })
            }
        },
    })

    app.Window.NewWithOptions(application.WebviewWindowOptions{
        Title: "Raw Message Demo",
        Name:  "main",
        Width: 400,
        Height: 300,
    })

    app.Run()
}
```

### assets/index.html

```html
<!DOCTYPE html>
<html>
<head>
    <title>Raw Message Demo</title>
    <style>
        body { font-family: sans-serif; padding: 20px; }
        button { margin: 5px; padding: 10px 20px; }
        #output { margin-top: 20px; padding: 10px; background: #f0f0f0; }
    </style>
</head>
<body>
    <h1>Raw Message Demo</h1>

    <button id="ping">Ping</button>
    <button id="echo">Echo "Hello"</button>

    <div id="output">Waiting for response...</div>

    <script type="module">
        import { System, Events } from '@wailsio/runtime'

        const output = document.getElementById('output')

        function send(type, data) {
            System.invoke(JSON.stringify({ type, data }))
        }

        document.getElementById('ping').onclick = () => send('ping')
        document.getElementById('echo').onclick = () => send('echo', 'Hello')

        Events.On('pong', (e) => {
            output.textContent = `Pong from ${e.data.window} at ${e.data.time}`
        })

        Events.On('echo', (e) => {
            output.textContent = `Echo: ${e.data}`
        })

        Events.On('error', (e) => {
            output.textContent = `Error: ${e.data.error}`
        })
    </script>
</body>
</html>
```

## 最佳实践

### 应当

- 将原始消息用于真正的性能关键路径
- 在处理程序中妥善处理错误
- 使用事件将响应发回前端
- 对于结构化数据，可考虑使用 JSON
- 保持消息处理快速，以免造成阻塞

### 不应

- 在服务绑定足以满足需求时使用原始消息
- 忘记验证传入的消息
- 在处理程序中执行长时间运行的操作并造成阻塞（请使用 goroutine）
- 当响应需要发送到特定窗口时忽略 window 参数

## 多窗口注意事项

`window` 参数用于标识消息由哪个窗口发送，因此你可以：

- 将响应发送到正确的窗口
- 实现窗口特定的行为
- 跟踪消息来源以便调试

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Respond only to the sending window
    window.EmitEvent("response", result)

    // Or broadcast to all windows
    app.Event.Emit("broadcast", result)
}
```

## 后续步骤

- [服务绑定](/features/bindings/services/)——适用于大多数应用程序的标准方法
- [事件](/guides/events-reference/)——用于从后端向前端通信的事件系统
- [性能](/guides/performance/)——常规性能优化
