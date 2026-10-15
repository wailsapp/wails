---
title: "建立自訂傳輸層"
description: "瞭解如何建立及自訂自己的 Wails v3 IPC 傳輸層"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Wails v3 可讓您提供自訂 IPC 傳輸層，同時保留所有產生的繫結與事件通訊。因此，您可以使用 WebSocket、自訂通訊協定或任何其他傳輸機制，取代預設以 HTTP fetch 為基礎的傳輸方式。

## 概觀

Wails 預設會從前端透過 `/wails/runtime` 傳送 HTTP fetch 請求，以便與後端通訊。自訂傳輸 API 可讓您：

- 使用 WebSocket、gRPC 或任何自訂通訊協定取代 HTTP 傳輸
- 維持與 Wails 程式碼產生功能的完整相容性
- 保留所有現有的繫結、事件、對話方塊及其他 Wails 功能
- 實作自己的連線管理、驗證及錯誤處理

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

## 使用方式

### 1. 實作 Transport 介面

實作 `Transport` 介面以建立自訂傳輸：

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type MyCustomTransport struct {
    // Your fields
}

func (t *MyCustomTransport) Start(ctx context.Context, processor *application.MessageProcessor) error {
    // Initialize your transport (WebSocket server, gRPC server, etc.)
    // When you receive requests, call processor.HandleRuntimeCallWithIDs()
    return nil
}

func (t *MyCustomTransport) Stop() error {
    // Clean up your transport
    return nil
}
```

### 2. 設定應用程式

將自訂傳輸傳入應用程式選項：

```go
func main() {
    app := application.New(application.Options{
        Name: "My App",
        Transport: &MyCustomTransport{},
        // ... other options
    })

    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}
```

### 3. 修改前端執行階段

如果使用自訂傳輸，您需要修改前端執行階段，使其使用您的傳輸，而非 HTTP fetch。請實作將用於處理請求的 `RuntimeTransport` 介面：

```typescript
const { setTransport } = await import('/wails/runtime.js');

class MyRuntimeTransport {
  call(objectID: number, method: number, windowName: string, args: any): Promise<any> {
    // TODO: implement IPC call with your transport protocol

    return resp;
  }
}

const myTransport = new MyRuntimeTransport();
setTransport(myTransport);
```

## 注意事項

- 如果未指定自訂傳輸，預設 HTTP 傳輸仍會繼續運作
- 產生的繫結維持不變，只有傳輸層會變更
- 事件、對話方塊、剪貼簿及所有其他 Wails 功能皆可無縫運作
- 您必須負責自訂傳輸中的錯誤處理、重新連線邏輯及安全性
- 所提供的 WebSocket 範例僅供示範，若要用於正式環境，可能需要進一步強化

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

對於以瀏覽器為基礎的部署，或當您想要透過自訂傳輸同時提供資源與 IPC 時，請實作 `AssetServerTransport` 介面：

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

- 在瀏覽器而非 WebView 中執行應用程式
- 在自訂 IPC 傳輸之外，另透過 HTTP 提供資源
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

- 所有靜態資源（HTML、CSS、JS、圖片等）
- `/wails/runtime.js` - Wails 執行階段程式庫

## 另請參閱

- `transport.go` - 核心傳輸介面與型別
- `messageprocessor.go` - 處理所有 Wails IPC 的底層訊息處理器
