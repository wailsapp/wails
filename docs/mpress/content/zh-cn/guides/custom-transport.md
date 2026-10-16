---
title: "创建自定义传输层"
description: "了解如何创建和定制自己的 Wails v3 自定义 IPC 传输层"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Wails v3 允许你提供自定义 IPC 传输层，同时保留所有生成的绑定和事件通信。这样，你可以使用 WebSocket、自定义协议或任何其他传输机制替换默认的基于 HTTP fetch 的传输。

## 概述

默认情况下，Wails 通过 `/wails/runtime` 使用前端发出的 HTTP fetch 请求与后端通信。你可以使用自定义传输 API：

- 使用 WebSocket、gRPC 或任意自定义协议替换 HTTP 传输
- 保持与 Wails 代码生成机制完全兼容
- 保留所有现有的绑定、事件、对话框及其他 Wails 功能
- 实现自己的连接管理、身份验证和错误处理

## 架构

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

## 用法

### 1. 实现 Transport 接口

实现 `Transport` 接口以创建自定义传输：

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

### 2. 配置应用程序

将自定义传输传入应用程序选项：

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

### 3. 修改前端运行时

如果使用自定义传输，则需要修改前端运行时，使其使用你的传输而不是 HTTP fetch。实现用于处理请求的 `RuntimeTransport` 接口：

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

## 注意事项

- 如果未指定自定义传输，默认 HTTP 传输仍可继续工作
- 生成的绑定保持不变，只有传输层会发生变化
- 事件、对话框、剪贴板及其他所有 Wails 功能均可透明运行
- 你需要负责自定义传输中的错误处理、重连逻辑和安全性
- 所提供的 WebSocket 示例仅用于演示，投入生产环境前可能需要进一步加固

## API 参考

### Transport 接口

```go
type Transport interface {
    Start(ctx context.Context, messageProcessor *application.MessageProcessor) error
    // JSClient returns the JavaScript shim that the runtime injects into the
    // window so that frontend code can call into the transport.
    JSClient() []byte
    Stop() error
}
```

### AssetServerTransport 接口（可选）

对于基于浏览器的部署，或者当你希望通过自定义传输同时提供静态资源和 IPC 时，请实现 `AssetServerTransport` 接口：

```go
type AssetServerTransport interface {
    Transport

    // ServeAssets configures the transport to serve assets alongside IPC.
    // The assetHandler is Wails' internal asset server that handles all assets,
    // runtime.js, capabilities, flags, etc.
    ServeAssets(assetHandler http.Handler) error
}
```

**何时实现此接口：**

- 在浏览器而非 WebView 中运行应用程序
- 在使用自定义 IPC 传输的同时通过 HTTP 提供静态资源
- 构建可通过网络访问的应用程序

**实现示例：**

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

调用 `ServeAssets()` 时，assetHandler 会提供：

- 所有静态资源（HTML、CSS、JS、图像等）
- `/wails/runtime.js` - Wails 运行时库

## 另请参阅

- `transport.go` - 核心传输接口和类型
- `messageprocessor.go` - 处理所有 Wails IPC 的底层消息处理器
