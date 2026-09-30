---
title: "原始訊息"
description: "為效能至關重要的應用程式實作自訂的前端到後端通訊"
slug: "guides/raw-messages"
sourcePath: "guides/raw-messages.md"
---

原始訊息在前端與後端之間提供低階通訊通道，繞過標準繫結系統。這是以便利性換取速度。

## 何時使用原始訊息

原始訊息最適合極端的邊緣案例：

- **超高頻率更新**：每秒傳送數千則訊息，且每一微秒都至關重要
- **自訂訊息通訊協定**：需要完全掌控資料傳輸格式時

@note{type="tip"}
對幾乎所有使用案例，建議使用標準[服務繫結](/features/bindings/services/)，因為其能以可忽略不計的額外負擔提供型別安全、自動序列化及更佳的開發者體驗。

@end

## 後端設定

在應用程式選項中設定`RawMessageHandler`：

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

### 處理常式簽章

```go
RawMessageHandler func(window Window, message string, originInfo *application.OriginInfo)
```

| 參數 | 型別 | 說明 |
| --- | --- | --- |
| `window` | `Window` | 傳送訊息的視窗 |
| `message` | `string` | 原始訊息內容 |
| `originInfo` | `*application.OriginInfo` | 訊息來源的來源資訊 |

#### OriginInfo 結構

```go
type OriginInfo struct {
	Origin      string
	TopOrigin   string
	IsMainFrame bool
}
```

| 欄位 | 型別 | 說明 |
| --- | --- | --- |
| `Origin` | `string` | 傳送訊息之文件的來源 URL |
| `TopOrigin` | `string` | 頂層來源 URL（在 iframe 中可能與 Origin 不同） |
| `IsMainFrame` | `bool` | 訊息是否源自主框架 |

#### 各平台的可用資訊

- **macOS**：會提供`Origin`和`IsMainFrame`
- **Windows**：會提供`Origin`和`TopOrigin`
- **Linux**：只會提供`Origin`

### 來源驗證

@note{type="caution"}
絕不能因訊息已傳入處理常式，就假定訊息安全。在處理敏感作業或會修改狀態的作業之前，必須驗證來源資訊。

@end

**處理傳入訊息之前，一律先驗證其來源。**`originInfo`參數會提供重要的安全性資訊，必須加以驗證以防止未經授權的存取。 惡意內容、遭入侵的內容或非預期的指令碼可能會傳送原始訊息。如果未驗證來源，可能會處理來自不受信任來源的命令。請使用`originInfo`確保訊息來自預期的來源。

### 主要驗證要點

- **一律檢查`Origin`**：確認來源符合預期的受信任來源（本機資產通常是`wails://wails`或`http://wails.localhost`，也可能是應用程式的特定來源）
- **驗證`IsMainFrame`**（macOS）：留意訊息是否來自 iframe，因為這可能表示內嵌內容處於不同的安全性脈絡
- **使用`TopOrigin`**（Windows）：處理框架內容時，確認頂層來源
- **拒絕非預期的來源**：拒絕來自未明確允許之來源的訊息，以安全方式失敗

@note{type="info"}
以`wails:`為前置字串的訊息保留供 Wails 內部通訊使用，不會傳遞給處理常式。

@end

## 前端設定

使用`System.invoke()`傳送原始訊息：

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

### 使用預先建置的套件組合

如果不使用 npm，請透過全域`wails`物件存取`invoke`：

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

## 結構化訊息

對於複雜資料，請將其序列化為 JSON：

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

### 後端

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

## 效能比較

| 方式 | 額外負擔 | 型別安全性 | 使用情境 |
| --- | --- | --- | --- |
| 服務繫結 | 較高 | 完整 | 一般用途 |
| 原始訊息 | 極低 | 手動 | 高頻率、效能關鍵 |

### 基準測試範例

對於簡單的承載資料，原始訊息每秒可處理的訊息數量明顯高於服務繫結：

```go
// Raw message handler - minimal overhead
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Direct string processing, no reflection or marshaling
    counter++
}
```

## 完整範例

以下是實作簡易命令通訊協定的完整範例：

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

## 最佳實務

### 建議做法

- 在真正攸關效能的路徑中使用原始訊息
- 在處理常式中妥善處理錯誤
- 使用事件將回應傳回前端
- 對結構化資料，可考慮使用 JSON
- 保持快速的訊息處理速度，以免造成阻塞

### 避免的做法

- 在服務繫結已足以應付需求時使用原始訊息
- 忘記驗證傳入的訊息
- 在處理常式中執行耗時操作而造成阻塞（請使用 goroutine）
- 當回應需要傳送至特定視窗時，忽略 window 參數

## 多視窗注意事項

`window`參數可識別訊息是由哪個視窗傳送，讓你可以：

- 將回應傳送至正確的視窗
- 實作視窗專屬行為
- 追蹤訊息來源以進行偵錯

```go
RawMessageHandler: func(window application.Window, message string, originInfo *application.OriginInfo) {
    // Respond only to the sending window
    window.EmitEvent("response", result)

    // Or broadcast to all windows
    app.Event.Emit("broadcast", result)
}
```

## 後續步驟

- [服務繫結](/features/bindings/services/)－適用於大多數應用程式的標準方法
- [事件](/guides/events-reference/)－用於後端至前端通訊的事件系統
- [效能](/guides/performance/)－一般效能最佳化
