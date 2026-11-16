---
title: "创建自定义传输层"
description: "了解如何创建和自定义您自己的 Wails v3 自定义 IPC 传输层"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Wails v3允许您提供自定义 IPC 传输层，同时保留所有生成的绑定和事件通信。因此，您可以用 WebSocket、自定义协议或任何其他传输机制取代默认的基于 HTTP fetch 的传输方式。

## 概述

默认情况下，Wails 使用来自前端的 HTTP fetch 请求，通过`/wails/runtime`与后端通信。自定义传输 API 可让您：

- 使用 WebSocket、gRPC 或任意自定义协议取代 HTTP 传输
- 保持与 Wails 代码生成的完全兼容
- 保留所有现有的绑定、事件、对话框及其他 Wails 功能
- 实现您自己的连接管理、身份验证和错误处理

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

## 运行完整的 WebSocket 示例

[WebSocket 传输示例](https://github.com/wailsapp/wails/tree/master/v3/examples/websocket-transport)包含协议的两端：请求 ID 与响应、生命周期管理、JavaScript 客户端和事件传递。请以这些文件为起点，而不是仅实现`Start`和`Stop`的传输。

在已安装原生平台依赖项的 Wails 仓库工作副本中：

```bash
cd v3/examples/websocket-transport
go run .
```

该示例会打开原生窗口，并在端口9099上启动 WebSocket 端点。请在可信的本地开发环境中使用它。部署可通过网络访问的变体之前，请限制监听地址、验证来源并对客户端进行身份验证。示例的来源策略较为宽松，不能作为生产环境的安全策略。

### 后端约定

传输必须实现全部三个必需方法。示例中的`WebSocketTransport`实现了`Start`、`JSClient`和`Stop`；`JSClient`提供 Wails 注入页面的 JavaScript 适配代码。请在实现中加入编译时接口检查：

```go
var _ application.Transport = (*WebSocketTransport)(nil)
```

使用示例中的构造函数配置应用程序：

```go
Transport: NewWebSocketTransport("127.0.0.1:9099"),
```

此字段应放在`application.Options`中，与常规资源和服务并列。除构造函数调用外，还需复制实现文件。`Start`接收`MessageProcessor`；`Stop`必须关闭监听器、活动连接和待处理工作。

### 前端约定

示例中的`websocket-transport.js`实现了运行时传输，包括请求与响应的关联以及对失败调用的拒绝处理。通过资源提供该文件后，请先初始化它，再调用生成的服务绑定：

```javascript
import { setTransport } from '/wails/runtime.js'
import { createWebSocketTransport } from '/websocket-transport.js'

setTransport(createWebSocketTransport('ws://127.0.0.1:9099/wails/ws'))
```

生成的服务绑定保持不变。如果您的传输替换了客户端，则必须承载应用程序使用的方法，并保留其错误语义。事件传递需要可选的事件传输约定；任意`Transport`实现不会自动提供此功能。完整的 WebSocket 示例对此进行了演示。

对于远程部署，请使用 TLS，并明确设置身份验证和来源策略。仅在浏览器中托管时，还需要按下文所述提供应用程序资源。

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

对于基于浏览器的部署，或者希望通过自定义传输同时提供资源和 IPC 时，请实现`AssetServerTransport`接口：

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
- 在使用自定义 IPC 传输的同时通过 HTTP 提供资源
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

调用`ServeAssets()`时，assetHandler 提供：

- 所有静态资源（HTML、CSS、JS、图像等）
- `/wails/runtime.js` - Wails 运行时库

## 另请参阅

- `transport.go` - 核心传输接口和类型
- `messageprocessor.go` - 处理所有 Wails IPC 的底层消息处理器
