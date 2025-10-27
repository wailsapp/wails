# Transport API Improvements

This document outlines improvements to reduce developer complexity when implementing custom transports.

## Status Summary

| Improvement | Status | Lines Saved | Priority |
|-------------|--------|-------------|----------|
| 1. Pluggable Codec System | ✅ **Completed** | ~51 lines | Critical |
| 2. Callback Handler Helper | 🔄 Recommended | ~9 lines | High |
| 3. Message ID Generation | 🔄 Recommended | ~12 lines | Medium |
| 4. Message Builder Factory | 🔄 Recommended | ~10 lines | Medium |
| 5. Pending Request Manager | 🔄 Recommended | ~40 lines | High |
| 6. TransportMessage Export | 🔄 Recommended | ~7 lines | Low |
| 7. BaseHTTPTransport Helper | 🔄 Recommended | ~20 lines | Medium |
| 8. ConnectionManager Helper | 🔄 Recommended | ~30 lines | Medium |

**Total Potential Reduction**: ~179 additional lines (beyond the 51 already saved with codecs)

---

## ✅ 1. Pluggable Codec System (COMPLETED)

### Problem
Developers had to manually implement base64 decoding/encoding logic with 26 lines of boilerplate in 2 places (success and error paths).

### Solution Implemented
Created `TransportCodec` interface with built-in implementations:
- **Frontend**: `Base64JSONCodec`, `RawJSONCodec`, `RawStringCodec`
- **Backend**: `DefaultCodec`, `RawJSONCodec`, `RawStringCodec`

### Impact
**Before**: 52 lines of base64 decoding boilerplate
**After**: 1 line using codec
**Saved**: ~51 lines

### Files
- `v3/internal/runtime/desktop/@wailsio/runtime/src/transport-codec.ts`
- `v3/pkg/application/transport_codec.go`

---

## 🔄 2. Callback Handler Helper

### Problem
Developers must understand Wails' internal `window._wails.callResultHandler` mechanism and manually:
- Extract call-id from stored request args
- Check if window._wails exists
- Check if callResultHandler exists
- Manually invoke the handler

**Current complexity** (10 lines):
```javascript
// For binding calls (object=0, method=0), we need to call the Wails callback handler
// because Call.ByName expects window._wails.callResultHandler to be invoked
if (responseData && response.contentType?.includes('application/json')) {
    console.log('[WebSocket] Calling Wails result handler with JSON data');
    const callId = pending.request.args?.['call-id'];
    console.log('[WebSocket] Extracted call-id:', callId);
    console.log('[WebSocket] window._wails exists:', !!window._wails);
    console.log('[WebSocket] callResultHandler exists:', !!window._wails?.callResultHandler);
    if (callId && window._wails?.callResultHandler) {
        console.log('[WebSocket] Invoking callResultHandler');
        window._wails.callResultHandler(callId, responseData, true);
    }
    pending.resolve();
}
```

### Proposed Solution
Add helper function to `runtime.ts`:

```typescript
/**
 * Handles Wails callback invocation for binding calls
 * @param pending - Pending request object with stored args
 * @param responseData - Decoded response data
 * @param contentType - Response content type
 * @returns true if handled as binding call, false otherwise
 */
export function handleWailsCallback(pending: any, responseData: string, contentType: string): boolean {
    // Only for JSON responses (binding calls)
    if (!responseData || !contentType?.includes('application/json')) {
        return false;
    }

    // Extract call-id from stored request
    const callId = pending.request?.args?.['call-id'];

    // Invoke Wails callback handler if available
    if (callId && window._wails?.callResultHandler) {
        window._wails.callResultHandler(callId, responseData, true);
        return true;
    }

    return false;
}
```

**After**:
```javascript
if (handleWailsCallback(pending, responseData, response.contentType)) {
    pending.resolve(); // Handled by Wails
} else {
    pending.resolve(responseData); // Direct resolve
}
```

### Impact
**Before**: 10 lines with internal knowledge required
**After**: 1 line
**Saved**: ~9 lines
**Benefit**: Abstracts internal Wails callback mechanism from developers

---

## 🔄 3. Message ID Generation

### Problem
Developers must implement their own nanoid/UUID generator (12 lines):

```javascript
function nanoid(size = 21) {
    const alphabet = 'useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict';
    let id = '';
    let i = size;
    while (i--) {
        id += alphabet[(Math.random() * 64) | 0];
    }
    return id;
}
```

### Proposed Solution
Export existing ID generator from `runtime.ts`:

```typescript
// Already exists in runtime.ts as generateID()
export function generateMessageID(): string {
    return generateID();
}
```

**Usage**:
```javascript
import { generateMessageID } from '/wails/runtime.js';

const msgID = generateMessageID();
```

### Impact
**Before**: 12 lines of boilerplate
**After**: 1 import
**Saved**: ~12 lines

### Notes
- Wails already has `generateID()` internally
- Just needs to be exported for developer use

---

## 🔄 4. Message Builder Factory

### Problem
Developers must know exact message structure with `id`, `type`, `request` fields (11 lines):

```javascript
const message = {
    id: msgID,
    type: 'request',
    request: {
        object: objectID,
        method: method,
        args: args ? JSON.stringify(args) : undefined,
        windowName: windowName || undefined,
        clientId: clientId
    }
};
```

### Proposed Solution
Add factory function to `runtime.ts`:

```typescript
/**
 * Builds a transport request message
 * @param id - Unique message ID
 * @param objectID - Wails object ID (0=Call, 1=Clipboard, etc.)
 * @param method - Method ID within the object
 * @param windowName - Source window name (optional)
 * @param args - Method arguments (will be JSON stringified)
 * @returns Formatted transport request message
 */
export function buildTransportRequest(
    id: string,
    objectID: number,
    method: number,
    windowName: string,
    args: any
): any {
    return {
        id: id,
        type: 'request',
        request: {
            object: objectID,
            method: method,
            args: args ? JSON.stringify(args) : undefined,
            windowName: windowName || undefined,
            clientId: clientId
        }
    };
}
```

**Usage**:
```javascript
const message = buildTransportRequest(msgID, objectID, method, windowName, args);
```

### Impact
**Before**: 11 lines
**After**: 1 line
**Saved**: ~10 lines

---

## 🔄 5. Pending Request Manager

### Problem
Developers must implement Map-based request tracking with timeout cleanup (scattered across ~40 lines):

```javascript
// Field declaration
this.pendingRequests = new Map();

// Registration (lines 212-224)
const timeout = setTimeout(() => {
    if (this.pendingRequests.has(msgID)) {
        this.pendingRequests.delete(msgID);
        reject(new Error(`Request timeout (${this.requestTimeout}ms)`));
    }
}, this.requestTimeout);

this.pendingRequests.set(msgID, { resolve, reject, timeout, request: { object: objectID, method, args } });

// Retrieval and cleanup (lines 118-125)
const pending = this.pendingRequests.get(msg.id);
if (!pending) {
    console.warn('[WebSocket] No pending request for ID:', msg.id);
    return;
}
this.pendingRequests.delete(msg.id);
clearTimeout(pending.timeout);

// Cleanup on disconnect (lines 88-92)
this.pendingRequests.forEach(({ reject, timeout }) => {
    clearTimeout(timeout);
    reject(new Error('WebSocket connection closed'));
});
this.pendingRequests.clear();
```

### Proposed Solution
Create utility class in new file `transport-helpers.ts`:

```typescript
/**
 * Manages pending transport requests with automatic timeout handling
 */
export class PendingRequestManager {
    private pending = new Map<string, {
        resolve: Function;
        reject: Function;
        timeout: any;
        request: any;
    }>();

    /**
     * Add a pending request with timeout
     */
    add(
        id: string,
        resolve: Function,
        reject: Function,
        timeoutMs: number,
        request?: any
    ): void {
        const timeout = setTimeout(() => {
            if (this.pending.has(id)) {
                this.pending.delete(id);
                reject(new Error(`Request timeout (${timeoutMs}ms)`));
            }
        }, timeoutMs);

        this.pending.set(id, { resolve, reject, timeout, request });
    }

    /**
     * Get and remove a pending request
     */
    take(id: string): any {
        const item = this.pending.get(id);
        if (item) {
            this.pending.delete(id);
            clearTimeout(item.timeout);
        }
        return item;
    }

    /**
     * Check if request exists
     */
    has(id: string): boolean {
        return this.pending.has(id);
    }

    /**
     * Clear all pending requests with rejection
     */
    clear(reason?: string): void {
        this.pending.forEach(({ reject, timeout }) => {
            clearTimeout(timeout);
            reject(new Error(reason || 'Requests cleared'));
        });
        this.pending.clear();
    }

    /**
     * Get number of pending requests
     */
    get size(): number {
        return this.pending.size;
    }
}
```

**Usage**:
```javascript
import { PendingRequestManager } from '/wails/runtime.js';

class WebSocketTransport {
    constructor() {
        this.pendingRequests = new PendingRequestManager();
    }

    async call(objectID, method, windowName, args) {
        return new Promise((resolve, reject) => {
            const msgID = generateMessageID();

            // Add with automatic timeout handling
            this.pendingRequests.add(msgID, resolve, reject, this.requestTimeout,
                { object: objectID, method, args });

            // Send message...
        });
    }

    handleMessage(data) {
        const msg = JSON.parse(data);

        // Take and auto-cleanup
        const pending = this.pendingRequests.take(msg.id);
        if (!pending) return;

        // Process response...
    }

    onClose() {
        // Clear all pending with one call
        this.pendingRequests.clear('WebSocket connection closed');
    }
}
```

### Impact
**Before**: ~40 lines scattered across multiple methods
**After**: ~5 lines with class usage
**Saved**: ~35-40 lines
**Benefit**: Automatic timeout handling, thread-safe operations, cleaner code

---

## 🔄 6. TransportMessage Export

### Problem
Developers must define their own message struct (7 lines):

```go
type WebSocketMessage struct {
    ID      string           `json:"id"`
    Type    string           `json:"type"`
    Request *TransportRequest `json:"request,omitempty"`
    Response *TransportResponse `json:"response,omitempty"`
}
```

### Proposed Solution
Export standard message type from `transport.go`:

```go
// TransportMessage represents a message sent over custom transports.
// This is the standard message format for request/response matching.
type TransportMessage struct {
    // ID is a unique message identifier for request/response correlation
    ID string `json:"id"`

    // Type indicates the message kind: "request", "response", or "event"
    Type string `json:"type"`

    // Request contains the request data (for type="request")
    Request *TransportRequest `json:"request,omitempty"`

    // Response contains the response data (for type="response")
    Response *TransportResponse `json:"response,omitempty"`
}
```

**Usage**:
```go
import "github.com/wailsapp/wails/v3/pkg/application"

func (w *WebSocketTransport) handleWebSocket(conn *websocket.Conn) {
    var msg application.TransportMessage  // Use exported type
    err := conn.ReadJSON(&msg)
    // ...
}
```

### Impact
**Before**: 7 lines to define struct
**After**: 0 lines (import only)
**Saved**: ~7 lines

---

## 🔄 7. BaseHTTPTransport Helper

### Problem
Developers must manually handle:
- HTTP server creation and lifecycle
- Mux setup for WebSocket + assets
- ServeAssets implementation
- Graceful shutdown

**Current complexity** (24 lines in `transport_websocket_example.go`):

```go
// Lines 73-96: ServeAssets implementation
func (w *WebSocketTransport) ServeAssets(assetHandler http.Handler) error {
    mux := http.NewServeMux()

    // Mount WebSocket endpoint for IPC
    mux.HandleFunc("/wails/ws", w.handleWebSocket)

    // Mount asset server for all other requests
    mux.Handle("/", assetHandler)

    // Set the handler and start the server
    w.server.Handler = mux

    // Start server in background
    go func() {
        log.Printf("WebSocket transport serving assets and IPC on %s", w.addr)
        log.Printf("  - Assets: http://%s/", w.addr)
        log.Printf("  - WebSocket IPC: ws://%s/wails/ws", w.addr)
        if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Printf("WebSocket server error: %v", err)
        }
    }()

    return nil
}
```

### Proposed Solution
Create helper in `transport_helpers.go`:

```go
// BaseHTTPTransport provides HTTP server management for custom transports
type BaseHTTPTransport struct {
    addr   string
    server *http.Server
    mux    *http.ServeMux
}

// NewBaseHTTPTransport creates a new base HTTP transport helper
func NewBaseHTTPTransport(addr string) *BaseHTTPTransport {
    return &BaseHTTPTransport{
        addr:   addr,
        server: &http.Server{Addr: addr},
        mux:    http.NewServeMux(),
    }
}

// RegisterHandler registers a handler for a specific path
func (b *BaseHTTPTransport) RegisterHandler(pattern string, handler http.HandlerFunc) {
    b.mux.HandleFunc(pattern, handler)
}

// ServeAssets mounts the asset handler and starts the HTTP server
func (b *BaseHTTPTransport) ServeAssets(assetHandler http.Handler) error {
    // Mount asset server for all unhandled routes
    b.mux.Handle("/", assetHandler)

    // Set handler
    b.server.Handler = b.mux

    // Start server
    go func() {
        log.Printf("Transport serving on %s", b.addr)
        if err := b.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Printf("Server error: %v", err)
        }
    }()

    return nil
}

// Shutdown gracefully stops the HTTP server
func (b *BaseHTTPTransport) Shutdown(ctx context.Context) error {
    if b.server == nil {
        return nil
    }
    return b.server.Shutdown(ctx)
}
```

**Usage**:
```go
type WebSocketTransport struct {
    base    *application.BaseHTTPTransport
    // ... other fields
}

func NewWebSocketTransport(addr string) *WebSocketTransport {
    return &WebSocketTransport{
        base: application.NewBaseHTTPTransport(addr),
        // ...
    }
}

func (w *WebSocketTransport) ServeAssets(assetHandler http.Handler) error {
    // Register WebSocket endpoint
    w.base.RegisterHandler("/wails/ws", w.handleWebSocket)

    // Delegate to base helper
    return w.base.ServeAssets(assetHandler)
}

func (w *WebSocketTransport) Stop() error {
    return w.base.Shutdown(context.Background())
}
```

### Impact
**Before**: ~20-24 lines per transport
**After**: ~5 lines with helper
**Saved**: ~15-20 lines
**Benefit**: Consistent HTTP server handling, easier testing, less boilerplate

---

## 🔄 8. ConnectionManager Helper

### Problem
Developers must manually manage:
- Client connection map with mutex
- Thread-safe add/remove operations
- Broadcast to all clients
- Cleanup on shutdown

**Current complexity** (35 lines across multiple methods):

```go
// Fields
clients  map[*websocket.Conn]bool
mu       sync.RWMutex

// Add client (lines 123-125)
w.mu.Lock()
w.clients[conn] = true
w.mu.Unlock()

// Remove client (lines 127-131)
defer func() {
    w.mu.Lock()
    delete(w.clients, conn)
    w.mu.Unlock()
    conn.Close()
}()

// Broadcast (lines 195-202)
w.mu.RLock()
defer w.mu.RUnlock()

for conn := range w.clients {
    if err := conn.WriteJSON(msg); err != nil {
        log.Printf("Failed to broadcast event: %v", err)
    }
}

// Cleanup on stop (lines 104-110)
w.mu.Lock()
for conn := range w.clients {
    conn.Close()
}
w.clients = make(map[*websocket.Conn]bool)
w.mu.Unlock()
```

### Proposed Solution
Create helper in `transport_helpers.go`:

```go
// ConnectionManager manages WebSocket client connections thread-safely
type ConnectionManager struct {
    clients map[*websocket.Conn]bool
    mu      sync.RWMutex
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
    return &ConnectionManager{
        clients: make(map[*websocket.Conn]bool),
    }
}

// Add registers a new client connection
func (cm *ConnectionManager) Add(conn *websocket.Conn) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    cm.clients[conn] = true
}

// Remove unregisters and closes a client connection
func (cm *ConnectionManager) Remove(conn *websocket.Conn) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    delete(cm.clients, conn)
    conn.Close()
}

// Broadcast sends a message to all connected clients
func (cm *ConnectionManager) Broadcast(message interface{}) error {
    cm.mu.RLock()
    defer cm.mu.RUnlock()

    var lastErr error
    for conn := range cm.clients {
        if err := conn.WriteJSON(message); err != nil {
            lastErr = err
            log.Printf("Failed to broadcast to client: %v", err)
        }
    }
    return lastErr
}

// Send sends a message to a specific client
func (cm *ConnectionManager) Send(conn *websocket.Conn, message interface{}) error {
    cm.mu.RLock()
    defer cm.mu.RUnlock()

    if !cm.clients[conn] {
        return fmt.Errorf("client not found")
    }

    return conn.WriteJSON(message)
}

// CloseAll closes all client connections
func (cm *ConnectionManager) CloseAll() {
    cm.mu.Lock()
    defer cm.mu.Unlock()

    for conn := range cm.clients {
        conn.Close()
    }
    cm.clients = make(map[*websocket.Conn]bool)
}

// Count returns the number of connected clients
func (cm *ConnectionManager) Count() int {
    cm.mu.RLock()
    defer cm.mu.RUnlock()
    return len(cm.clients)
}
```

**Usage**:
```go
type WebSocketTransport struct {
    connections *application.ConnectionManager
    // ...
}

func NewWebSocketTransport(addr string) *WebSocketTransport {
    return &WebSocketTransport{
        connections: application.NewConnectionManager(),
        // ...
    }
}

func (w *WebSocketTransport) handleWebSocket(rw http.ResponseWriter, r *http.Request) {
    conn, err := w.upgrader.Upgrade(rw, r, nil)
    if err != nil {
        return
    }

    w.connections.Add(conn)
    defer w.connections.Remove(conn)

    // Handle messages...
}

func (w *WebSocketTransport) BroadcastEvent(event interface{}) error {
    return w.connections.Broadcast(event)
}

func (w *WebSocketTransport) Stop() error {
    w.connections.CloseAll()
    return w.base.Shutdown(context.Background())
}
```

### Impact
**Before**: ~30-35 lines of mutex management
**After**: ~5 lines with helper
**Saved**: ~25-30 lines
**Benefit**: Thread-safe by design, prevents common concurrency bugs, cleaner API

---

## Implementation Priority

### Phase 1: High Impact (Recommended Next)
1. ✅ **Codec System** - Already done, saves 51 lines
2. **Pending Request Manager** - Saves 40 lines, reduces complexity significantly
3. **Callback Handler Helper** - Saves 9 lines, abstracts internal mechanism

**Phase 1 Total**: ~100 lines saved

### Phase 2: Medium Impact
4. **BaseHTTPTransport Helper** - Saves 20 lines, improves consistency
5. **ConnectionManager Helper** - Saves 30 lines, prevents bugs
6. **Message Builder Factory** - Saves 10 lines, reduces errors
7. **Message ID Generation** - Saves 12 lines, simple improvement

**Phase 2 Total**: ~72 lines saved

### Phase 3: Low Impact (Nice to Have)
8. **TransportMessage Export** - Saves 7 lines, minimal effort

**Phase 3 Total**: ~7 lines saved

---

## Total Impact Summary

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| **Frontend Lines** | ~284 lines | ~110 lines | **61% reduction** |
| **Backend Lines** | ~206 lines | ~90 lines | **56% reduction** |
| **Developer Knowledge Required** | High (internal APIs, threading, encoding) | Low (public helpers) | Significant |
| **Bug Surface Area** | Large (manual memory management, timeouts, mutex) | Small (tested helpers) | Substantial |
| **Time to Implement** | ~2-3 days | ~4-6 hours | **~75% faster** |

---

## Example: Before vs After (Full Transport)

### Before (Current - 284 lines frontend, 206 lines backend)

**Frontend**: Implement nanoid, base64 decoding, pending requests, message building, callback handling
**Backend**: Implement message struct, HTTP server, mux setup, client tracking, mutex management, codec

### After (With All Helpers - ~110 lines frontend, ~90 lines backend)

**Frontend**:
```javascript
import {
    generateMessageID,
    buildTransportRequest,
    handleWailsCallback,
    PendingRequestManager,
    Base64JSONCodec
} from '/wails/runtime.js';

export class WebSocketTransport {
    constructor(url, options = {}) {
        this.ws = new WebSocket(url);
        this.requests = new PendingRequestManager();
        this.codec = options.codec || new Base64JSONCodec();
    }

    async call(objectID, method, windowName, args) {
        return new Promise((resolve, reject) => {
            const id = generateMessageID();
            this.requests.add(id, resolve, reject, 30000, { args });

            const msg = buildTransportRequest(id, objectID, method, windowName, args);
            this.ws.send(JSON.stringify(msg));
        });
    }

    handleMessage(data) {
        const msg = JSON.parse(data);
        const pending = this.requests.take(msg.id);
        if (!pending) return;

        const responseData = this.codec.decodeResponse(msg.response.data, msg.response.contentType);

        if (handleWailsCallback(pending, responseData, msg.response.contentType)) {
            pending.resolve();
        } else {
            pending.resolve(responseData);
        }
    }
}
```

**Backend**:
```go
import "github.com/wailsapp/wails/v3/pkg/application"

type WebSocketTransport struct {
    base        *application.BaseHTTPTransport
    connections *application.ConnectionManager
    codec       application.TransportCodec
    handler     application.TransportHandler
}

func NewWebSocketTransport(addr string, opts ...application.TransportOption) *WebSocketTransport {
    t := &WebSocketTransport{
        base:        application.NewBaseHTTPTransport(addr),
        connections: application.NewConnectionManager(),
        codec:       application.NewDefaultCodec(),
    }
    for _, opt := range opts {
        opt(t)
    }
    return t
}

func (w *WebSocketTransport) ServeAssets(assetHandler http.Handler) error {
    w.base.RegisterHandler("/wails/ws", w.handleWebSocket)
    return w.base.ServeAssets(assetHandler)
}

func (w *WebSocketTransport) handleWebSocket(rw http.ResponseWriter, r *http.Request) {
    conn, _ := upgrader.Upgrade(rw, r, nil)

    w.connections.Add(conn)
    defer w.connections.Remove(conn)

    for {
        var msg application.TransportMessage
        if err := conn.ReadJSON(&msg); err != nil {
            break
        }

        response := w.handler.HandleRuntimeCall(context.Background(), msg.Request)
        response.Data, _ = w.codec.EncodeResponse(response.Data.([]byte), response.ContentType)

        w.connections.Send(conn, application.TransportMessage{
            ID:       msg.ID,
            Type:     "response",
            Response: response,
        })
    }
}

func (w *WebSocketTransport) Stop() error {
    w.connections.CloseAll()
    return w.base.Shutdown(context.Background())
}
```

**Reduction**: From ~490 total lines to ~200 total lines (**~59% reduction**)

---

## Recommendations

### Must Have (Phase 1)
- ✅ Codec System (done)
- PendingRequestManager
- handleWailsCallback

These three provide the most value with minimal API surface.

### Should Have (Phase 2)
- BaseHTTPTransport
- ConnectionManager
- Message helpers

These significantly improve developer experience and prevent common bugs.

### Nice to Have (Phase 3)
- TransportMessage export

Low effort, small benefit, but completes the API.

---

## Next Steps

1. **Review this proposal** - Are these the right abstractions?
2. **Prioritize** - Which helpers provide the most value?
3. **Implement** - Start with Phase 1 (high impact)
4. **Document** - Update examples to use new helpers
5. **Test** - Ensure helpers work across different transport types

---

**Status**: 1 of 8 improvements completed (Codec System ✅)
**Remaining Work**: ~7 helpers across 3 phases
**Estimated Effort**: 2-3 days for all phases
**Developer Benefit**: 59% code reduction, significantly simplified implementation
