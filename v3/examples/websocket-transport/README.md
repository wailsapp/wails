# WebSocket Transport Example

This example demonstrates how to use a custom WebSocket transport for Wails IPC instead of the default HTTP fetch transport. All Wails bindings and features work identically - only the underlying transport layer changes.

## What This Example Shows

- How to configure a WebSocket transport on the backend using `application.NewWebSocketTransport()`
- How to override the runtime transport on the frontend using `setTransport()` and `createWebSocketTransport()`
- Full compatibility with generated bindings - no code generation changes needed
- Real-time connection status monitoring
- Automatic reconnection handling

## Architecture

```
┌─────────────────────────────────────────┐
│  Frontend (JavaScript)                  │
│  - setTransport(wsTransport)           │
│  - All bindings work unchanged         │
└───────────────┬─────────────────────────┘
                │
                │ WebSocket (ws://localhost:9998)
                │
┌───────────────▼─────────────────────────┐
│  Backend (Go)                           │
│  - WebSocketTransport on port 9998     │
│  - Standard MessageProcessor           │
│  - All Wails features available        │
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
   - WebSocket server listening on `ws://localhost:9998/wails/ws`
   - Real-time connection status indicator

## Backend Setup

The backend configuration is simple - just pass a `WebSocketTransport` to the application options:

```go
// Create WebSocket transport on port 9998
wsTransport := application.NewWebSocketTransport(":9998")

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

**The bindings are generated once, but the transport can be swapped at runtime!**

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

## Benefits of WebSocket Transport

1. **Better Performance**: Persistent connection reduces overhead for frequent calls
2. **Lower Latency**: No TCP/TLS handshake per request
3. **Server Push**: WebSocket enables server-to-client push notifications (future feature)
4. **Binary Support**: Can efficiently transfer binary data
5. **Full Compatibility**: All existing Wails features continue to work

## Files

- `main.go` - Application setup with WebSocket transport
- `GreetService.go` - Example service with bound methods
- `assets/index.html` - Frontend UI with WebSocket transport configuration

## Creating Custom Transports

You can create your own custom transport by implementing the `RuntimeTransport` interface:

### Backend (Go)

```go
type MyTransport struct {
    // Your fields
}

func (t *MyTransport) Start(ctx context.Context, handler application.TransportHandler) error {
    // Initialize your transport and call handler.HandleRuntimeCall() for requests
    return nil
}

func (t *MyTransport) Stop() error {
    // Clean up
    return nil
}
```

### Frontend (TypeScript)

```typescript
import { setTransport, type RuntimeTransport } from "/wails/runtime.js";

const myTransport: RuntimeTransport = {
    call: async (objectID, method, windowName, args) => {
        // Your transport implementation
        // Must return the response data
    }
};

setTransport(myTransport);
```

## Browser Support

The WebSocket transport implements the `AssetServerTransport` interface, which enables browser-based deployments. When this interface is implemented, Wails automatically configures the transport to serve both assets and IPC:

```go
// WebSocketTransport implements AssetServerTransport
func (w *WebSocketTransport) ServeAssets(assetHandler http.Handler) error {
    // Mount WebSocket IPC at /wails/ws
    // Mount assets at /
    // Start HTTP server
}
```

When configured this way, you can:

1. **Run in a webview** (default):
   - Assets served through webview
   - IPC via WebSocket on port 9099

2. **Run in a browser**:
   - Open `http://localhost:9099/` in any browser
   - Assets served via HTTP
   - IPC via WebSocket on the same port
   - Full application functionality maintained

The transport automatically detects whether `ServeAssets()` is called and configures itself accordingly. This means the same code works for both webview and browser deployments.

### Browser Deployment Example

```bash
# Start the app
go run .

# Open in browser
open http://localhost:9099/
```

All features work identically in both environments - the transport layer is completely transparent to the application code.

## Notes

- The WebSocket server runs on port 9099 (configurable)
- All Wails generated bindings remain unchanged
- Events, dialogs, clipboard, and all other Wails features work transparently
- The example WebSocket transport is production-ready but you may want to add authentication
- Connection status is monitored every second for demonstration purposes
- Browser support is automatic when the transport implements `AssetServerTransport`

## See Also

- `/v3/pkg/application/transport.go` - Transport interfaces and types
- `/v3/pkg/application/transport_websocket_example.go` - WebSocket transport implementation
- `/v3/pkg/application/CUSTOM_TRANSPORT.md` - Complete transport documentation
- `/v3/internal/runtime/desktop/@wailsio/runtime/src/transport-websocket.ts` - Frontend WebSocket transport
