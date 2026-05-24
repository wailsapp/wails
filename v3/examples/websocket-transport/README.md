# WebSocket Transport Example

This example demonstrates how to use a custom transport like WebSocket for Wails IPC instead of the default HTTP fetch transport. All Wails bindings and features work identically - only the underlying transport layer changes.

## What This Example Shows

- How to configure a WebSocket transport on the backend using `NewWebSocketTransport()`
- How to override the runtime transport on the frontend using `setTransport()` and `createWebSocketTransport()`
- Full compatibility with generated bindings - no code generation changes needed
- Real-time connection status monitoring
- Automatic reconnection handling

## Architecture

```text
┌─────────────────────────────────────────┐
│  Frontend (JavaScript)                  │
│  - setTransport(wsTransport)            │
│  - All bindings work unchanged          │
└───────────────┬─────────────────────────┘
                │
                │ WebSocket (ws://localhost:9099)
                │
┌───────────────▼─────────────────────────┐
│  Backend (Go)                           │
│  - WebSocketTransport on port 9099      │
│  - Standard MessageProcessor            │
│  - All Wails features available         │
└─────────────────────────────────────────┘
```

## How to Run

1. Navigate to this directory:
   ```bash
   cd v3/examples/websocket-transport
   ```

2. Run the example:
   ```bash
   go run .
   ```

3. The application will start with:
   - WebView window displaying the UI
   - WebSocket server listening on `ws://localhost:9099/wails/ws`
   - Real-time connection status indicator

## Backend Setup

The backend configuration is simple - just pass a `WebSocketTransport` to the application options:

```go
// Create WebSocket transport on port 9099
wsTransport := NewWebSocketTransport(":9099")

app := application.New(application.Options{
    Name: "WebSocket Transport Example",
    Services: []application.Service{
        application.NewService(&GreetService{}),
    },
    Assets: application.AssetOptions{
        Handler: application.BundledAssetFileServer(assets),
    },
    // Use WebSocket transport instead of default HTTP
    Transport: wsTransport,
})
```

## Frontend Setup

The frontend uses the WebSocket transport with **generated bindings**:

```typescript
import { setTransport } from "/wails/runtime.js";
import { createWebSocketTransport } from "/websocket-transport.js";
import { GreetService } from "/bindings/github.com/wailsapp/wails/v3/examples/websocket-transport/index.js";

// Create and configure WebSocket transport
const wsTransport = createWebSocketTransport('ws://localhost:9099/wails/ws', {
    reconnectDelay: 2000,      // Reconnect after 2 seconds if disconnected
    requestTimeout: 30000      // Request timeout of 30 seconds
});

// Set as the active transport
setTransport(wsTransport);

// Now all generated bindings use WebSocket instead of HTTP fetch!
const result = await GreetService.Greet("World");
```

**Key Point**: The generated bindings (`GreetService.Greet()`, `GreetService.Echo()`, etc.) automatically use whatever transport is configured via `setTransport()`. This proves the custom transport hijack works seamlessly with Wails code generation!

## Features Demonstrated

### 1. Generated Bindings with Custom Transport
All generated bindings work identically with WebSocket transport:
- `GreetService.Greet(name)` - Simple string parameter and return
- `GreetService.Echo(message)` - Echo back messages
- `GreetService.Add(a, b)` - Multiple parameters with numeric types
- `GreetService.GetTime()` - No parameters, string return

### 2. Connection Management
- Automatic connection establishment on startup
- Visual connection status indicator (green = connected, red = disconnected)
- Automatic reconnection with configurable delay
- Graceful handling of connection failures

### 3. Error Handling
- Request timeouts
- Connection errors
- Backend method errors
- All propagate correctly to the frontend
