# Custom Transport Layer

Wails v3 allows you to provide a custom IPC transport layer while retaining all generated bindings and event communication. This enables you to replace the default HTTP fetch-based transport with WebSockets, custom protocols, or any other transport mechanism.

## Overview

By default, Wails uses HTTP fetch requests from the frontend to communicate with the backend via `/wails/runtime`. The custom transport API allows you to:

- Replace the HTTP transport with WebSockets, gRPC, or any custom protocol
- Maintain full compatibility with Wails code generation
- Keep all existing bindings, events, dialogs, and other Wails features
- Implement your own connection management, authentication, and error handling

## Architecture

```
┌─────────────────────────────────────────────────┐
│  Frontend (TypeScript)                          │
│  - Generated bindings still work               │
│  - Your custom client transport               │
└──────────────────┬──────────────────────────────┘
                   │
                   │ Your Protocol (WebSocket/etc)
                   │
┌──────────────────▼──────────────────────────────┐
│  Backend (Go)                                   │
│  - Your Transport implementation               │
│  - Wails MessageProcessor                    │
│  - All existing Wails infrastructure           │
└─────────────────────────────────────────────────┘
```

## Usage

### 1. Implement the Transport Interface

Create a custom transport by implementing the `Transport` interface:

```go
package main

import (
    "context"
    "github.com/wailsapp/wails/v3/pkg/application"
)

type MyCustomTransport struct {
    // Your fields
}

func (t *MyCustomTransport) Start(ctx context.Context, handler application.TransportHandler) error {
    // Initialize your transport (WebSocket server, gRPC server, etc.)
    // When you receive requests, call handler.HandleRuntimeCall()
    return nil
}

func (t *MyCustomTransport) Stop() error {
    // Clean up your transport
    return nil
}
```

### 2. Configure Your Application

Pass your custom transport to the application options:

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

### 3. Modify Frontend Runtime

If using a custom transport, you'll need to modify the frontend runtime to use your transport instead of HTTP fetch. Replace the `runtimeCallWithID` function in your frontend:

```typescript
// Example: WebSocket-based transport
const ws = new WebSocket('ws://localhost:9998/wails/ws');

async function runtimeCallWithID(objectID: number, method: number, windowName: string, args: any): Promise<any> {
    const msgID = nanoid();

    return new Promise((resolve, reject) => {
        const handler = (event: MessageEvent) => {
            const response = JSON.parse(event.data);
            if (response.id === msgID) {
                ws.removeEventListener('message', handler);
                if (response.type === 'response') {
                    if (response.response.statusCode === 200) {
                        resolve(JSON.parse(new TextDecoder().decode(response.response.data)));
                    } else {
                        reject(new Error(new TextDecoder().decode(response.response.data)));
                    }
                }
            }
        };

        ws.addEventListener('message', handler);

        ws.send(JSON.stringify({
            id: msgID,
            type: 'request',
            request: {
                object: objectID,
                method: method,
                args: args ? JSON.stringify(args) : '',
                windowName: windowName
            }
        }));
    });
}
```

## Complete WebSocket Example

### Backend Implementation

```go
package main

import (
    "log"
    "github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
    // Create WebSocket transport on port 9998
    wsTransport := application.NewWebSocketTransport(":9998")

    app := application.New(application.Options{
        Name:        "WebSocket App",
        Description: "An app using WebSocket transport",
        Services: []application.Service{
            application.NewService(&GreetService{}),
        },
        Assets: application.AlphaAssets,
        Transport: wsTransport,
    })

    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello, " + name + "!"
}
```

### Frontend Runtime Modification

Replace the `runtimeCallWithID` function in `v3/internal/runtime/desktop/@wailsio/runtime/src/runtime.ts`:

```typescript
// WebSocket connection
let ws: WebSocket | null = null;
const pendingRequests = new Map<string, {resolve: Function, reject: Function}>();

function initWebSocket() {
    if (ws?.readyState === WebSocket.OPEN) return;

    ws = new WebSocket('ws://localhost:9998/wails/ws');

    ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        if (msg.type === 'response' && pendingRequests.has(msg.id)) {
            const {resolve, reject} = pendingRequests.get(msg.id)!;
            pendingRequests.delete(msg.id);

            if (msg.response.statusCode === 200) {
                const data = new TextDecoder().decode(new Uint8Array(msg.response.data));
                if (msg.response.contentType?.includes('application/json')) {
                    resolve(JSON.parse(data));
                } else {
                    resolve(data);
                }
            } else {
                reject(new Error(new TextDecoder().decode(new Uint8Array(msg.response.data))));
            }
        }
    };

    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
}

async function runtimeCallWithID(objectID: number, method: number, windowName: string, args: any): Promise<any> {
    initWebSocket();

    return new Promise((resolve, reject) => {
        const msgID = nanoid();
        pendingRequests.set(msgID, {resolve, reject});

        ws!.send(JSON.stringify({
            id: msgID,
            type: 'request',
            request: {
                object: objectID,
                method: method,
                args: args ? JSON.stringify(args) : '',
                windowName: windowName,
                clientId: clientId
            }
        }));
    });
}
```

## TransportRequest Structure

When implementing a custom transport, you'll receive `TransportRequest` objects that you need to pass to the handler:

```go
type TransportRequest struct {
    Object     int    // Which Wails subsystem (0=Call, 1=Clipboard, 2=Application, etc.)
    Method     int    // Which method within the object
    Args       string // JSON-encoded arguments
    WindowID   string // Source window ID (optional)
    WindowName string // Source window name (optional)
    ClientID   string // Frontend client ID (optional)
}
```

## TransportResponse Structure

The handler returns `TransportResponse` objects that you need to send back to the frontend:

```go
type TransportResponse struct {
    StatusCode  int    // HTTP status code equivalent (200=success, 422=error)
    ContentType string // "application/json" or "text/plain"
    Data        []byte // Response payload
    Error       error  // Error if call failed
}
```

## Benefits

1. **Better Performance**: WebSockets reduce overhead compared to HTTP fetch for frequent calls
2. **Server Push**: Easily implement server-to-client push notifications
3. **Custom Authentication**: Implement your own auth mechanism in the transport layer
4. **Protocol Flexibility**: Use any protocol (WebSocket, gRPC, custom binary protocol)
5. **Full Compatibility**: All existing Wails features continue to work

## Notes

- The default HTTP transport continues to work if no custom transport is specified
- Generated bindings remain unchanged - only the transport layer changes
- Events, dialogs, clipboard, and all other Wails features work transparently
- You're responsible for error handling, reconnection logic, and security in your custom transport
- The WebSocket example provided is for demonstration and may need hardening for production use

## API Reference

### Transport Interface

```go
type Transport interface {
    Start(ctx context.Context, handler TransportHandler) error
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
- `/wails/capabilities` - Capability information
- `/wails/flags` - Application flags

**Browser deployment:**
```bash
# Start your app
go run .

# Open in browser (if transport serves assets)
open http://localhost:8080/
```

The same code works for both webview and browser deployments - Wails automatically calls `ServeAssets()` if the interface is implemented.

### TransportHandler Interface

```go
type TransportHandler interface {
    HandleRuntimeCall(ctx context.Context, req *TransportRequest) *TransportResponse
}
```

### Factory Functions

- `NewHTTPTransport()` - Creates the default HTTP transport (used automatically if no transport specified)
- `NewWebSocketTransport(addr string)` - Example WebSocket transport implementation

## See Also

- `transport.go` - Core transport interfaces and types
- `transport_websocket_example.go` - Complete WebSocket transport implementation
- `messageprocessor.go` - The underlying message processor that handles all Wails IPC
