---
title: "伺服器建置"
description: "以 HTTP 伺服器執行 Wails 應用程式，不建立原生 GUI 視窗"
slug: "guides/server-build"
sourcePath: "guides/server-build.md"
---

Wails v3 支援伺服器模式，可讓您將應用程式當作純 HTTP 伺服器執行，無須建立原生視窗或使用 GUI 相依項目。因此，同一個 Wails 應用程式可部署至伺服器、容器和網頁瀏覽器。

伺服器模式適用於：

- **Docker／容器部署**：無須 X11／Wayland 相依項目即可執行
- **伺服器端應用程式**：部署為可透過瀏覽器存取的網頁伺服器
- **僅限網頁存取**：桌面版與網頁版共用同一套程式碼
- **CI/CD 測試**：無須顯示伺服器即可執行整合測試
- **微服務**：在無頭後端服務中使用 Wails 繫結

## 快速入門

伺服器模式透過 `server` 建置標籤啟用。應用程式碼不需要變更，只要使用此標籤建置即可：

```bash
# Using Taskfile (recommended)
wails3 task build:server
wails3 task run:server

# Or build directly with Go
go build -tags server -o myapp-server .
```

以下是最小範例：

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

同一套程式碼既可建置為桌面模式（不使用標籤），也可建置為伺服器模式（使用 `-tags server`）。

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

    // Additional origins allowed to open WebSocket connections and call
    // /wails/runtime. Same-origin requests are always allowed.
    WebSocketOriginPatterns: []string{"app.example.com"},

    // Disable WebSocket and runtime origin checks. Unsafe; default: false.
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

系統會自動在 `/health` 提供健康狀態檢查端點：

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

此端點適用於：

- Kubernetes 存活／就緒探針
- 負載平衡器健康狀態檢查
- 監控系統

### 服務繫結

所有服務繫結的運作方式都與桌面模式相同：

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

前端可以使用標準 Wails 執行階段呼叫這些繫結：

```javascript
const greeting = await wails.Call.ByName('main.GreetService.Greet', 'World');
```

### 事件

在伺服器模式下，事件可雙向傳遞：

- **前端至後端**：瀏覽器發出的事件會透過 HTTP 傳送，並由 Go 事件處理常式接收
- **後端至前端**：Go 發出的事件會透過 WebSocket 廣播至所有已連線的瀏覽器

每個瀏覽器分頁都會以具有唯一名稱（`browser-1`、`browser-2` 等）的「視窗」表示，並可透過 `event.Sender` 存取：

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

### 正常關閉

伺服器會妥善處理 `SIGINT` 和 `SIGTERM` 訊號：

1. 停止接受新連線
2. 等待進行中的要求完成（最長等待 `ShutdownTimeout`）
3. 執行 `OnShutdown` 掛鉤
4. 以相反順序關閉服務

## 與桌面模式的差異

| 功能 | 桌面模式 | 伺服器模式 |
| --- | --- | --- |
| 原生視窗 | 會建立 | 瀏覽器視窗（`browser-N`） |
| 系統匣 | 可用 | 不可用 |
| 原生對話方塊 | 可用 | 不可用 |
| 應用程式選單 | 可用 | 不可用 |
| 螢幕資訊 | 可用 | 沒有原生螢幕；`GetPrimary()` 會傳回 `nil` |
| 服務繫結 | 可運作 | 可運作 |
| 事件 | 可運作 | 可運作（透過 WebSocket） |
| 資源 | 透過 webview | 透過 HTTP |
| 需要 CGO | 是 | 否 |

### 視窗 API 行為

在伺服器模式下，系統會安全地處理視窗相關 API：

- `app.Window.NewWithOptions(options)`：傳回非 nil 的 Go 視窗物件；這不會建立原生桌面視窗
- `app.Hide()`／`app.Show()`：不執行任何操作
- `app.Screen.GetPrimary()`：傳回 `*application.Screen`；若沒有可用的主要螢幕，其值為 `nil`。解除參照前請先檢查

請區分 Go 視窗物件與原生視窗的可用性。建置伺服器時，請保護僅限桌面的行為；非 nil 的 Go 視窗物件並不表示原生視窗或螢幕存在。

### 支援的平台

伺服器模式以純 Go 實作，無須 CGO 即可針對 Wails 支援的每個 `GOOS`（包括 `android` 和 `ios`）進行編譯。在行動裝置目標上，只會建置無頭伺服器：

- 原生行動裝置 API `application.Android` 和 `application.IOS` 無法使用。使用這些 API 的程式碼必須透過 `!server` 建置條件約束，排除在伺服器建置之外。
- 預設記錄器為桌面版記錄器；在開發建置中，它會寫入 stderr，而不是 logcat 或 WebView 主控台。

## 建置正式版本

### 使用 Task（建議）

使用 `wails3 init` 建立的專案包含 `build:server` 工作：

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

就是這麼簡單！您的應用程式將可透過 `http://localhost:8080` 存取。

您可以使用幾個選項自訂建置：

```bash
# Use a custom image tag
wails3 task build:docker TAG=myapp:v1.0.0

# Run on a different port
wails3 task run:docker PORT=3000
```

產生的 `Dockerfile.server` 會建立以 distroless 為基礎的最小映像檔。它會自動處理網路繫結，因此您可以從容器外部存取應用程式。

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
健康檢查範例使用`wget`。若使用 distroless 基礎映像檔，您必須在映像檔中包含健康檢查二進位檔，或使用外部健康檢查機制（例如 Docker 的`curl`選項或 sidecar 容器）。

@end

### 自訂 Dockerfile

如果需要更完整的控制，可以自行建立 Dockerfile。請務必設定`WAILS_SERVER_HOST=0.0.0.0`，讓伺服器接受來自容器外部的連線：

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

1. **預設繫結至 localhost** - 僅在必要時使用`0.0.0.0`
2. **在正式環境中使用 TLS** - 設定`ServerOptions.TLS`
3. **置於反向 Proxy 後方** - 使用 nginx/traefik 提供額外安全防護
4. **讓 WebSocket 與執行階段呼叫保持同源** - 僅透過`WebSocketOriginPatterns`新增受信任的來源；避免使用`WebSocketAllowAllOrigins`
5. **驗證所有輸入** - 採用與任何 Web 應用程式相同的安全性實務

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

在需要覆寫伺服器設定而不變更程式碼的部署情境中，Wails 可辨識以下環境變數：

| 變數 | 說明 | 預設值 |
| --- | --- | --- |
| `WAILS_SERVER_HOST` | 要繫結的網路介面 | `localhost` |
| `WAILS_SERVER_PORT` | 要監聽的連接埠 | `8080` |

這些環境變數的優先順序高於程式碼中的`ServerOptions`，因此 Docker 範例會設定`WAILS_SERVER_HOST=0.0.0.0`；這讓容器無須變更應用程式即可接受外部連線。

## 另請參閱

- [自訂傳輸](/guides/custom-transport/) - 用於進階 IPC 自訂
- [服務](/features/bindings/services/) - 服務繫結文件
- [事件](/guides/events-reference/) - 事件系統文件

### 執行階段請求大小

在進行 JSON 處理之前，傳送至`/wails/runtime`的請求大小上限為64 MiB。超過此限制的一般請求會收到 HTTP 413，未包含`Content-Length`的請求也包括在內。分塊執行階段上傳仍分別採用每個區塊1 MiB 與組合後承載資料64 MiB 的限制。適當時，請使用應用程式中介軟體或反向 Proxy 設定更低的限制。

### 執行階段請求來源

`/wails/runtime`僅接受`POST`請求；其他方法會收到 HTTP 405。若瀏覽器請求的`Origin`既不是請求主機，也不符合`WebSocketOriginPatterns`，則會收到 HTTP 403；未包含`Origin`且其`Sec-Fetch-Site`為`cross-site`或`same-site`的請求也是如此。未攜帶這兩個標頭的請求（例如來自`curl`或其他非瀏覽器用戶端的請求）會獲得接受。若反向 Proxy 重寫`Host`標頭，請將公開來源新增至`WebSocketOriginPatterns`。

`WebSocketOriginPatterns`用於公開來源與應用程式所見`Host`不同的反向 Proxy 部署。不支援由不同來源提供、並直接呼叫`/wails/runtime`的前端：該端點不會發出 CORS 標頭，且會以 HTTP 405回應預檢`OPTIONS`請求。請由 Wails 伺服器提供前端，或將兩者置於相同的 Proxy 來源後方。

### 迴送 Host 檢查

當伺服器繫結至迴送位址（預設的`localhost`、任何`127.0.0.0/8`位址或`::1`，無論是在`ServerOptions.Host`或`WAILS_SERVER_HOST`中設定）時，每個請求都必須攜帶`Host`標頭，並在其中指定`localhost`、IP 常值、繫結主機，或與`WebSocketOriginPatterns`項目之主機部分相符的主機；任何其他`Host`都會收到 HTTP 403。這可阻止 DNS 重新繫結攻擊；在這類攻擊中，惡意網頁會將自己的主機名稱指向`127.0.0.1`，使其請求通過上述同源檢查。系統不會比較連接埠，且會拒絕`localhost.`（結尾有句點）或`127.1`等寫法。即使已設定`WebSocketAllowAllOrigins`，此檢查仍會套用。若本機反向 Proxy 轉送其公開`Host`，請將該主機新增至`WebSocketOriginPatterns`。繫結至其他位址（例如容器中的`0.0.0.0`）不受影響。
