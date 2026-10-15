---
title: "伺服器建置"
description: "將 Wails 應用程式作為 HTTP 伺服器執行，無需原生 GUI 視窗"
slug: "guides/server-build"
sourcePath: "guides/server-build.md"
---

Wails v3 支援伺服器模式，可讓您將應用程式作為純 HTTP 伺服器執行，而不建立原生視窗，也不需要 GUI 相依項目。因此，同一個 Wails 應用程式可部署至伺服器、容器及網頁瀏覽器。

伺服器模式適用於：

- **Docker／容器部署**：不依賴 X11／Wayland 即可執行
- **伺服器端應用程式**：部署為可透過瀏覽器存取的網頁伺服器
- **僅限網頁存取**：桌面版與網頁版共用同一套程式碼
- **CI/CD 測試**：無需顯示伺服器即可執行整合測試
- **微服務**：在無頭後端服務中使用 Wails 繫結

## 快速開始

伺服器模式透過 `server` 建置標籤啟用。應用程式碼無需變更，只要使用該標籤進行建置即可：

```bash
# Using Taskfile (recommended)
wails3 task build:server
wails3 task run:server

# Or build directly with Go
go build -tags server -o myapp-server .
```

以下是最精簡的範例：

```go
package main

import (
    "embed"
    "log"

    "github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

func main() {
    app := application.New(application.Options{
        Name: "My App",
        // Server options are used when built with -tags server
        Server: application.ServerOptions{
            Host: "localhost",
            Port: 8080,
        },
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.AssetFileServerFS(assets),
        },
    })

    log.Println("Starting application...")
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

同一份程式碼既可建置為桌面模式（不使用該標籤），也可建置為伺服器模式（使用 `-tags server`）。

## 設定

### ServerOptions

使用 `ServerOptions` 設定 HTTP 伺服器：

```go
Server: application.ServerOptions{
    // Host to bind to. Default: "localhost"
    // Use "0.0.0.0" to listen on all interfaces
    Host: "localhost",

    // Port to listen on. Default: 8080
    Port: 8080,

    // Request read timeout. Default: 30s
    ReadTimeout: 30 * time.Second,

    // Response write timeout. Default: 30s
    WriteTimeout: 30 * time.Second,

    // Idle connection timeout. Default: 120s
    IdleTimeout: 120 * time.Second,

    // Graceful shutdown timeout. Default: 30s
    ShutdownTimeout: 30 * time.Second,

    // Additional origins allowed to open WebSocket connections.
    // Same-origin connections are always allowed.
    WebSocketOriginPatterns: []string{"app.example.com"},

    // Disable WebSocket origin checks. Unsafe; default: false.
    WebSocketAllowAllOrigins: false,

    // TLS configuration (optional)
    TLS: &application.TLSOptions{
        CertFile: "/path/to/cert.pem",
        KeyFile:  "/path/to/key.pem",
    },
},
```

## 功能

### 健康狀態檢查端點

`/health` 會自動提供健康狀態檢查端點：

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

此端點適用於：

- Kubernetes 存活／就緒探針
- 負載平衡器健康狀態檢查
- 監控系統

### 服務繫結

所有服務繫結的運作方式都與桌面模式完全相同：

```go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}

// Register in options
Services: []application.Service{
    application.NewService(&GreetService{}),
},
```

前端可使用標準 Wails 執行階段呼叫這些繫結：

```javascript
const greeting = await wails.Call.ByName('main.GreetService.Greet', 'World');
```

### 事件

在伺服器模式中，事件可雙向傳遞：

- **前端到後端**：瀏覽器發出的事件會透過 HTTP 傳送，並由您的 Go 事件處理常式接收
- **後端到前端**：Go 發出的事件會透過 WebSocket 廣播至所有已連線的瀏覽器

每個瀏覽器分頁都會表示為具有唯一名稱（`browser-1`、`browser-2` 等）的「視窗」，並可透過 `event.Sender` 存取：

```go
// Listen for events from browsers
app.Event.On("user-action", func(event *application.CustomEvent) {
    log.Printf("Event from %s: %v", event.Sender, event.Data)
    // event.Sender will be "browser-1", "browser-2", etc.
})

// Emit events to all connected browsers
app.Event.Emit("server-update", data)
```

從前端：

```javascript
// Emit event to server (and all other browsers)
await wails.Events.Emit('user-action', { action: 'click' });

// Listen for events from server
wails.Events.On('server-update', (event) => {
    console.log('Update from server:', event.data);
});
```

### 順暢關閉

伺服器會妥善處理 `SIGINT` 和 `SIGTERM` 訊號：

1. 停止接受新連線
2. 等待作用中的要求完成（最多等待 `ShutdownTimeout`）
3. 執行 `OnShutdown` 掛鉤
4. 以相反順序關閉服務

## 與桌面模式的差異

| 功能 | 桌面模式 | 伺服器模式 |
| --- | --- | --- |
| 原生視窗 | 會建立 | 瀏覽器視窗（`browser-N`） |
| 系統匣 | 可用 | 不可用 |
| 原生對話方塊 | 可用 | 不可用 |
| 應用程式選單 | 可用 | 不可用 |
| 螢幕資訊 | 可用 | 傳回錯誤 |
| 服務繫結 | 可運作 | 可運作 |
| 事件 | 可運作 | 可運作（透過 WebSocket） |
| 資產 | 透過 WebView | 透過 HTTP |
| 需要 CGO | 是 | 否 |

### 視窗 API 行為

在伺服器模式中，與視窗相關的 API 會以安全的方式處理：

- `app.Window.NewWithOptions()`：記錄警告並傳回 nil
- `app.Hide()`／`app.Show()`：不執行任何操作
- `app.Screen.GetPrimary()`：傳回錯誤

如此一來，參照視窗的程式碼便可執行而不會當機，但視窗操作不會產生任何效果。

## 建置正式環境版本

### 使用 Task（建議）

使用`wails3 init`建立的專案包含一個`build:server`工作：

```bash
# Build for server mode
wails3 task build:server

# Build and run
wails3 task run:server
```

### 手動建置

```bash
# Build with server mode
go build -tags server -o myapp-server .
```

### Docker

Wails 專案包含可直接使用的 Docker 設定。若要在容器中建置並執行應用程式：

```bash
# Build the Docker image
wails3 task build:docker

# Run it
wails3 task run:docker
```

就這麼簡單！您可以在`http://localhost:8080`存取應用程式。

您可以使用幾個選項自訂建置：

```bash
# Use a custom image tag
wails3 task build:docker TAG=myapp:v1.0.0

# Run on a different port
wails3 task run:docker PORT=3000
```

產生的`Dockerfile.server`會建立以 distroless 為基礎的最小映像檔。它會自動處理網路繫結，因此可從容器外部存取您的應用程式。

### Docker Compose

對於較複雜的部署，以下是包含健康狀態檢查的 Docker Compose 設定：

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - WAILS_SERVER_HOST=0.0.0.0
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

@note{type="info"}
此健康狀態檢查範例使用`wget`。如果您使用 distroless 基礎映像檔，則需要在映像檔中加入健康狀態檢查二進位檔，或使用外部健康狀態檢查機制（例如 Docker 的`curl`選項或 Sidecar 容器）。

@end

### 自訂 Dockerfile

如果需要更多控制權，您可以自行建立 Dockerfile。請務必設定`WAILS_SERVER_HOST=0.0.0.0`，讓伺服器接受來自容器外部的連線：

```dockerfile
# Build stage
FROM golang:alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY . .
RUN go mod tidy
RUN go build -tags server -ldflags="-s -w" -o server .

# Runtime stage
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server
COPY --from=builder /app/frontend/dist /frontend/dist
EXPOSE 8080
ENV WAILS_SERVER_HOST=0.0.0.0
ENTRYPOINT ["/server"]
```

## 安全性考量

部署伺服器模式應用程式時：

1. **預設繫結至 localhost**－僅在必要時使用`0.0.0.0`
2. **在正式環境中使用 TLS**－設定`ServerOptions.TLS`
3. **置於反向 Proxy 後方**－使用 nginx/traefik 提供額外安全防護
4. **讓 WebSocket 保持同源**－僅使用`WebSocketOriginPatterns`新增受信任的來源；避免使用`WebSocketAllowAllOrigins`
5. **驗證所有輸入**－採用與任何 Web 應用程式相同的安全實務

## 範例

完整範例位於`v3/examples/server/`：

```bash
cd v3/examples/server

# Using Taskfile
task dev

# Or run directly
go run -tags server .

# Open http://localhost:8080 in browser
```

## 環境變數

若部署情境需要在不變更程式碼的情況下覆寫伺服器設定，Wails 可識別下列環境變數：

| 變數 | 說明 | 預設值 |
| --- | --- | --- |
| `WAILS_SERVER_HOST` | 要繫結的網路介面 | `localhost` |
| `WAILS_SERVER_PORT` | 要接聽的連接埠 | `8080` |

這些環境變數的優先順序高於程式碼中的`ServerOptions`，因此 Docker 範例會設定`WAILS_SERVER_HOST=0.0.0.0`；這可讓容器接受外部連線，而不必變更應用程式。

## 另請參閱

- [自訂傳輸](/guides/custom-transport/)－用於進階 IPC 自訂
- [服務](/features/bindings/services/)－服務繫結文件
- [事件](/guides/events-reference/)－事件系統文件

### 執行階段要求大小

傳送至`/wails/runtime`的要求在進行 JSON 處理前，大小上限為64 MiB。超過此限制的一般要求會收到 HTTP 413，未包含`Content-Length`的要求亦同。分塊執行階段上傳仍各自受每個分塊1 MiB 及組合後承載資料64 MiB 的大小限制。適當時，請使用應用程式中介軟體或反向 Proxy 設定更低的限制。
