---
title: "Creating Custom Transport Layer"
description: "Learn how to create and customize your own Wails v3 custom IPC transport layer"
slug: "guides/custom-transport"
sourcePath: "guides/custom-transport.md"
---

Wails v3 allows you to provide a custom IPC transport layer while retaining all generated bindings and event communication. This enables you to replace the default HTTP fetch-based transport with WebSockets, custom protocols, or any other transport mechanism.

## Overview

By default, Wails uses HTTP fetch requests from the frontend to communicate with the backend via `/wails/runtime`. The custom transport API allows you to:

- Replace the HTTP transport with WebSockets, gRPC, or any custom protocol
- Maintain full compatibility with Wails code generation
- Keep all existing bindings, events, dialogs, and other Wails features
- Implement your own connection management, authentication, and error handling

## Architecture

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

## Run the complete WebSocket example

The [WebSocket transport example](https://github.com/wailsapp/wails/tree/master/v3/examples/websocket-transport) includes both halves of the protocol: request IDs and responses, lifecycle management, the JavaScript client, and event delivery. Start from those files rather than a transport that only implements `Start` and `Stop`.

From a Wails repository checkout with the native platform dependencies installed:

```bash
cd v3/examples/websocket-transport
go run .
```

The example opens a native window and starts a WebSocket endpoint at port 9099. Use it in a trusted local development environment. Before deploying a network-accessible variant, restrict the listener, validate origins and authenticate clients. The sample origin policy is permissive and is not a production security policy.

### Backend contract

A transport must implement all three required methods. The example's `WebSocketTransport` implements `Start`, `JSClient` and `Stop`; `JSClient` supplies the JavaScript shim Wails injects into the page. Include a compile-time interface check in your implementation:

```go
var _ application.Transport = (*WebSocketTransport)(nil)
```

Configure the application using the constructor from the example:

```go
Transport: NewWebSocketTransport("127.0.0.1:9099"),
```

This field belongs in `application.Options`, alongside the ordinary assets and services. Copy the implementation files as well as the constructor call. `Start` receives the `MessageProcessor`; `Stop` must close listeners, active connections and pending work.

### Frontend contract

The example's `websocket-transport.js` implements the runtime transport, including request/response correlation and rejection of failed calls. After serving that file from your assets, initialise it before invoking generated service bindings:

```javascript
import { setTransport } from '/wails/runtime.js'
import { createWebSocketTransport } from '/websocket-transport.js'

setTransport(createWebSocketTransport('ws://127.0.0.1:9099/wails/ws'))
```

Generated service bindings remain unchanged. If your transport replaces the client, it must carry the methods your application uses and preserve their error semantics. Event delivery requires the optional event transport contract; it is not supplied automatically by an arbitrary `Transport` implementation. The complete WebSocket example demonstrates it.

For remote deployments, use TLS and an explicit authentication and origin policy. Browser-only hosting additionally requires serving the application assets, as described below.

## API Reference

### Transport Interface

```go
type Transport interface {
    Start(ctx context.Context, messageProcessor *application.MessageProcessor) error
    // JSClient returns the JavaScript shim that the runtime injects into the
    // window so that frontend code can call into the transport.
    JSClient() []byte
    Stop() error
}
```

### AssetServerTransport Interface (Optional)

For browser-based deployments or when you want to serve both assets and IPC through your custom transport, implement the `AssetServerTransport` interface:

```go
type AssetServerTransport interface {
    Transport

    // ServeAssets configures the transport to serve assets alongside IPC.
    // The assetHandler is Wails' internal asset server that handles all assets,
    // runtime.js, capabilities, flags, etc.
    ServeAssets(assetHandler http.Handler) error
}
```

**When to implement this interface:**

- Running the app in a browser instead of a webview
- Serving assets over HTTP alongside your custom IPC transport
- Building network-accessible applications

**Example implementation:**

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

When `ServeAssets()` is called, the assetHandler provides:

- All static assets (HTML, CSS, JS, images, etc.)
- `/wails/runtime.js` - The Wails runtime library

## See Also

- `transport.go` - Core transport interfaces and types
- `messageprocessor.go` - The underlying message processor that handles all Wails IPC
