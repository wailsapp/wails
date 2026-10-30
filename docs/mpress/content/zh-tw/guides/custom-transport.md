---
title: "建立自訂傳輸層"
description: "瞭解如何建立及自訂自己的 Wails v3 自訂 IPC 傳輸層"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Wails v3 可讓你提供自訂 IPC 傳輸層，同時保留所有產生的繫結與事件通訊。你可以藉此將預設以 HTTP fetch 為基礎的傳輸方式，替換為 WebSocket、自訂通訊協定或任何其他傳輸機制。

## 概觀

Wails 預設會讓前端透過 `/wails/runtime` 使用 HTTP fetch 要求與後端通訊。自訂傳輸 API 可讓你：

- 以 WebSocket、gRPC 或任何自訂通訊協定取代 HTTP 傳輸
- 維持與 Wails 程式碼產生功能的完整相容性
- 保留所有現有的繫結、事件、對話方塊及其他 Wails 功能
- 實作自己的連線管理、驗證與錯誤處理

## 架構

```text
┌─────────────────────────────────────────────────┐
│  Frontend (TypeScript)                          │
│  - Generated bindings still work                │
│  - Your custom client transport                 │
└──────────────────┬──────────────────────────────┘
                   │
                   │ Your Protocol (WebSocket/etc)
                   │
┌──────────────────▼──────────────────────────────┐
│  Backend (Go)                                   │
│  - Your Transport implementation                │
│  - Wails MessageProcessor                       │
│  - All existing Wails infrastructure            │
└─────────────────────────────────────────────────┘
```

## 執行完整的 WebSocket 範例

[WebSocket 傳輸範例](https://github.com/wailsapp/wails/tree/master/v3/examples/websocket-transport)包含通訊協定的兩端：要求 ID 與回應、生命週期管理、JavaScript 用戶端及事件傳遞。請以這些檔案為起點，不要從僅實作 `Start` 和 `Stop` 的傳輸實作開始。

在已安裝原生平台相依項目的 Wails 儲存庫工作目錄中：

```bash
cd v3/examples/websocket-transport
go run .
```

此範例會開啟原生視窗，並在連接埠 9099 啟動 WebSocket 端點。請在受信任的本機開發環境中使用。部署可透過網路存取的變體之前，請限制監聽位址、驗證來源並驗證用戶端身分。範例的來源原則較為寬鬆，並非正式環境的安全性原則。

### 後端契約

傳輸實作必須實作全部三個必要方法。範例中的 `WebSocketTransport` 實作了 `Start`、`JSClient` 和 `Stop`；`JSClient` 則提供 Wails 注入頁面的 JavaScript 轉接程式。請在你的實作中加入編譯期介面檢查：

```go
var _ application.Transport = (*WebSocketTransport)(nil)
```

使用範例中的建構函式設定應用程式：

```go
Transport: NewWebSocketTransport("127.0.0.1:9099"),
```

此欄位應放在 `application.Options` 中，與一般資產及服務並列。除了建構函式呼叫，也請複製實作檔案。`Start` 會接收 `MessageProcessor`；`Stop` 必須關閉監聽器、作用中的連線及待處理工作。

### 前端契約

範例中的 `websocket-transport.js` 實作了執行階段傳輸，包括要求與回應的關聯，以及拒絕失敗的呼叫。從資產提供該檔案後，請先將它初始化，再叫用產生的服務繫結：

```javascript
import { setTransport } from '/wails/runtime.js'
import { createWebSocketTransport } from '/websocket-transport.js'

setTransport(createWebSocketTransport('ws://127.0.0.1:9099/wails/ws'))
```

產生的服務繫結維持不變。如果你的傳輸實作取代了用戶端，就必須承載應用程式使用的方法，並保留其錯誤語意。事件傳遞需要選用的事件傳輸契約；任意的 `Transport` 實作不會自動提供此契約。完整的 WebSocket 範例示範了這項功能。

遠端部署時，請使用 TLS，並明確設定驗證與來源原則。若僅在瀏覽器中託管，還必須依照下文說明提供應用程式資產。

## API 參考

### Transport 介面

```go
type Transport interface {
    Start(ctx context.Context, messageProcessor *application.MessageProcessor) error
    // JSClient returns the JavaScript shim that the runtime injects into the
    // window so that frontend code can call into the transport.
    JSClient() []byte
    Stop() error
}
```

### AssetServerTransport 介面（選用）

若要進行瀏覽器型部署，或想透過自訂傳輸同時提供資產與 IPC，請實作 `AssetServerTransport` 介面：

```go
type AssetServerTransport interface {
    Transport

    // ServeAssets configures the transport to serve assets alongside IPC.
    // The assetHandler is Wails' internal asset server that handles all assets,
    // runtime.js, capabilities, flags, etc.
    ServeAssets(assetHandler http.Handler) error
}
```

**何時應實作此介面：**

- 在瀏覽器而非 webview 中執行應用程式
- 透過 HTTP 提供資產，並同時使用自訂 IPC 傳輸
- 建置可透過網路存取的應用程式

**實作範例：**

```go
func (t *MyTransport) ServeAssets(assetHandler http.Handler) error {
    mux := http.NewServeMux()

    // Mount your IPC endpoint
    mux.HandleFunc("/my/ipc/endpoint", t.handleIPC)

    // Mount Wails asset server for everything else
    mux.Handle("/", assetHandler)

    // Start HTTP server
    t.httpServer.Handler = mux
    go t.httpServer.ListenAndServe()

    return nil
}
```

呼叫 `ServeAssets()` 時，assetHandler 會提供：

- 所有靜態資產（HTML、CSS、JS、圖片等）
- `/wails/runtime.js` - Wails 執行階段程式庫

## 另請參閱

- `transport.go` - 核心傳輸介面與型別
- `messageprocessor.go` - 處理所有 Wails IPC 的底層訊息處理器
