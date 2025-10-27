# Developer Experience Improvements for Custom Transports

## Current Pain Points

### 1. Binding Import Path Issues ❌
**Problem**: Generated bindings use `@wailsio/runtime` which doesn't resolve in browser without build tools.

**Current Workaround** (manual):
```bash
# Generate bindings
wails3 generate bindings

# Manually copy to assets
cp -r frontend/bindings assets/

# Manually edit import path in greetservice.js
# Change: from "@wailsio/runtime"
# To:     from "/wails/runtime.js"
```

**DX-Friendly Solution**:
```go
// Option 1: Automatic binding serving
app := application.New(application.Options{
    BindingsPath: "frontend/bindings",  // Wails automatically serves at /bindings/
})

// Option 2: Generate with correct import path
wails3 generate bindings --runtime-import="/wails/runtime.js"  // Or auto-detect from config
```

### 2. No Built-in Transport Override API ❌
**Problem**: Developers must create entire transport implementations from scratch.

**Current** (200+ lines of boilerplate):
```javascript
// Developer must implement:
class WebSocketTransport {
    constructor() { /* 40 lines */ }
    connect() { /* 50 lines */ }
    handleMessage() { /* 60 lines */ }
    call() { /* 30 lines */ }
    // ... more methods
}
```

**DX-Friendly Solution**:
```javascript
import { createWebSocketTransport } from '/wails/runtime.js';

// Wails provides ready-to-use transports
const transport = createWebSocketTransport('ws://localhost:9099/wails/ws', {
    codec: 'base64',  // or 'json', 'msgpack'
    reconnect: true,
    timeout: 30000
});

setTransport(transport);
```

### 3. Codec Selection Not User-Friendly ❌
**Problem**: Developers must understand base64 encoding, JSON marshaling, etc.

**Current**:
```javascript
// Frontend
import { Base64JSONCodec } from '/wails/runtime.js';
const transport = new WebSocketTransport(url, { codec: new Base64JSONCodec() });

// Backend
wsTransport := application.NewWebSocketTransport(":9099",
    application.WithCodec(application.NewDefaultCodec()))
```

**DX-Friendly Solution**:
```javascript
// Frontend - simple string option
const transport = createWebSocketTransport(url, { codec: 'base64' });

// Backend - matched automatically or explicit
wsTransport := application.NewWebSocketTransport(":9099",
    application.WithCodec("base64"))  // String option
```

### 4. No Transport Discovery/Registration ❌
**Problem**: No way to list available transports or register custom ones globally.

**DX-Friendly Solution**:
```javascript
// List available transports
const transports = Wails.Transports.list();
// ['http', 'websocket', 'custom-plugin']

// Register custom transport
Wails.Transports.register('my-transport', MyTransportClass);

// Use by name
setTransport('websocket', { url: '...' });
```

### 5. No Transport Fallback/Retry ❌
**Problem**: If WebSocket fails, app is unusable. No automatic fallback to HTTP.

**DX-Friendly Solution**:
```javascript
// Automatic fallback
setTransport(['websocket', 'http'], {
    websocket: { url: 'ws://...' },
    http: { /* defaults */ }
});

// Or explicit fallback strategy
setTransport('websocket', {
    url: 'ws://...',
    fallback: 'http',
    retries: 3
});
```

---

## Proposed DX-Friendly API

### Frontend API (`/wails/runtime.js`)

```typescript
// ============================================
// 1. Simple Transport Creation
// ============================================

// Built-in transports with sane defaults
export function createWebSocketTransport(url: string, options?: {
    codec?: 'base64' | 'json' | 'raw' | TransportCodec;
    reconnect?: boolean;
    reconnectDelay?: number;
    timeout?: number;
}): RuntimeTransport;

export function createHTTPTransport(options?: {
    baseURL?: string;
    timeout?: number;
}): RuntimeTransport;

// ============================================
// 2. Transport Management
// ============================================

// Set active transport (existing, but enhanced)
export function setTransport(
    transport: RuntimeTransport | string | RuntimeTransport[],
    options?: any
): void;

// Get current transport
export function getTransport(): RuntimeTransport;

// Transport registry
export namespace Transports {
    export function register(name: string, factory: TransportFactory): void;
    export function list(): string[];
    export function create(name: string, options?: any): RuntimeTransport;
}

// ============================================
// 3. Codec Helpers
// ============================================

// Simple codec creation by name
export function createCodec(name: 'base64' | 'json' | 'raw' | 'msgpack'): TransportCodec;

// Codec registry
export namespace Codecs {
    export function register(name: string, codec: TransportCodec): void;
    export function list(): string[];
}

// ============================================
// 4. Transport Utilities
// ============================================

// Generate message IDs (already planned)
export function generateMessageID(): string;

// Build transport requests (already planned)
export function buildTransportRequest(
    id: string,
    objectID: number,
    method: number,
    windowName: string,
    args: any
): any;

// Handle Wails callbacks (already planned)
export function handleWailsCallback(
    pending: any,
    responseData: string,
    contentType: string
): boolean;

// Pending request manager (already planned)
export class PendingRequestManager {
    add(id: string, resolve: Function, reject: Function, timeoutMs: number, request?: any): void;
    take(id: string): any;
    has(id: string): boolean;
    clear(reason?: string): void;
    get size(): number;
}
```

### Backend API (`pkg/application`)

```go
// ============================================
// 1. Transport Factory Functions
// ============================================

// NewWebSocketTransport with codec string option
func NewWebSocketTransport(addr string, options ...TransportOption) *WebSocketTransport

// TransportOption functional options
func WithCodec(codec interface{}) TransportOption  // Accepts string or TransportCodec
func WithReconnect(enabled bool) TransportOption
func WithTimeout(ms int) TransportOption

// ============================================
// 2. Built-in Transport Helpers
// ============================================

// BaseHTTPTransport helper (already planned)
type BaseHTTPTransport struct { /* ... */ }
func NewBaseHTTPTransport(addr string) *BaseHTTPTransport

// ConnectionManager helper (already planned)
type ConnectionManager struct { /* ... */ }
func NewConnectionManager() *ConnectionManager

// ============================================
// 3. Codec Helpers
// ============================================

// Codec factory by name
func NewCodec(name string) (TransportCodec, error)
// "base64", "json", "raw", "msgpack"

// Or keep existing typed factories
func NewDefaultCodec() TransportCodec     // base64/JSON
func NewRawJSONCodec() TransportCodec     // direct JSON
func NewRawStringCodec() TransportCodec   // plain strings

// ============================================
// 4. Message Types (already planned)
// ============================================

// Export standard message type
type TransportMessage struct {
    ID       string              `json:"id"`
    Type     string              `json:"type"`
    Request  *TransportRequest   `json:"request,omitempty"`
    Response *TransportResponse  `json:"response,omitempty"`
}
```

### Build Tool Enhancements

```bash
# ============================================
# 1. Bindings Generation
# ============================================

# Generate with correct runtime import
wails3 generate bindings --runtime-import=/wails/runtime.js

# Or auto-detect from wails.json config
# wails.json:
{
  "bindings": {
    "output": "frontend/bindings",
    "runtimeImport": "/wails/runtime.js",  // Auto-use bundled runtime
    "serveFrom": "/bindings/"              // Auto-configure serving
  }
}

# ============================================
# 2. Transport Scaffolding
# ============================================

# Generate custom transport boilerplate
wails3 generate transport websocket --name=MyCustomTransport

# Creates:
# - backend/transport_mycustom.go (with interfaces implemented)
# - frontend/transports/mycustom.ts (with RuntimeTransport interface)
# - Example usage in README

# ============================================
# 3. Dev Server
# ============================================

# Dev server automatically serves bindings
wails3 dev --serve-bindings

# Or configured in wails.json
{
  "dev": {
    "serveBindings": true,
    "bindingsPath": "/bindings/"
  }
}
```

---

## Recommended Implementation Priority

### Phase 1: Core DX Improvements (High Priority)
1. ✅ **Codec System** - Done
2. **Binding Auto-Serving** - Wails automatically serves `frontend/bindings` at `/bindings/`
3. **Runtime Import Flag** - `wails3 generate bindings --runtime-import=/wails/runtime.js`
4. **Built-in WebSocket Transport** - Ready-to-use `createWebSocketTransport()`

### Phase 2: Helper APIs (Medium Priority)
5. **Helper Utilities** - Export `generateMessageID()`, `buildTransportRequest()`, etc.
6. **PendingRequestManager** - Export helper class
7. **BaseHTTPTransport** - Backend helper for HTTP server management
8. **ConnectionManager** - Backend helper for WebSocket connections

### Phase 3: Advanced Features (Lower Priority)
9. **Transport Registry** - `Transports.register()`, `.list()`, `.create()`
10. **Codec Registry** - String-based codec selection
11. **Transport Fallback** - Automatic HTTP fallback if WebSocket fails
12. **Transport Scaffolding** - `wails3 generate transport` command

---

## Example: Ideal Developer Experience

### Backend (3 lines to add custom transport)
```go
func main() {
    // That's it! Transport with sensible defaults
    wsTransport := application.NewWebSocketTransport(":9099")

    app := application.New(application.Options{
        Name: "My App",
        Services: []application.Service{
            application.NewService(&MyService{}),
        },
        Assets: application.AssetOptions{
            Handler: application.BundledAssetFileServer(assets),
        },
        Transport: wsTransport,  // Just set it
    })

    app.Run()
}
```

### Frontend (2 lines to override transport)
```javascript
import { createWebSocketTransport, setTransport } from '/wails/runtime.js';

// Bindings automatically served from /bindings/ by Wails
import { MyService } from '/bindings/myapp/index.js';

// One-liner transport override
setTransport(createWebSocketTransport('ws://localhost:9099/wails/ws'));

// Use generated bindings - they just work!
await MyService.DoSomething();
```

### Configuration (wails.json)
```json
{
  "name": "myapp",
  "bindings": {
    "output": "frontend/bindings",
    "runtimeImport": "/wails/runtime.js",
    "autoServe": true
  },
  "dev": {
    "serveBindings": true
  }
}
```

### Commands
```bash
# Generate bindings (auto-configured)
wails3 generate bindings

# Run (bindings auto-served)
wails3 dev

# No manual copying, no import path editing, no build steps!
```

---

## Current State vs Ideal State

| Feature | Current | Ideal | Benefit |
|---------|---------|-------|---------|
| **Binding Generation** | Manual import path fix | Auto-configured | No manual edits |
| **Binding Serving** | Manual copy to assets | Auto-served from `/bindings/` | No build step needed |
| **Transport Creation** | 200+ lines boilerplate | 1 line factory function | 99% less code |
| **Codec Selection** | Manual class instantiation | String option `codec: 'base64'` | Simpler API |
| **Error Handling** | Manual try/catch everywhere | Built-in reconnect & fallback | Robust by default |
| **Helper Utilities** | None - DIY | Exported helpers | Reusable components |
| **Transport Fallback** | Not possible | `['websocket', 'http']` | Better reliability |
| **Type Safety** | Manual typing | Generated TypeScript types | IDE autocomplete |

---

## Breaking Changes Assessment

### Phase 1 (Non-Breaking)
- Add `--runtime-import` flag (optional, defaults to current behavior)
- Add `createWebSocketTransport()` export (new export, doesn't break existing)
- Add binding auto-serving (opt-in via config)
- Export helper utilities (new exports, doesn't break existing)

### Phase 2 (Potentially Breaking)
- Change default binding import from `@wailsio/runtime` to `/wails/runtime.js`
  - **Migration**: Add `--runtime-import=@wailsio/runtime` flag for legacy projects
  - **Timeline**: 1 major version notice, migrate in next major

### Phase 3 (Breaking)
- Remove manual transport implementation requirement
  - **Migration**: Built-in transports cover 95% of use cases
  - **Custom transports**: Still supported, but simpler interface

---

## Conclusion

**Current State**: Custom transports work, but require significant manual setup and boilerplate.

**Ideal State**: Custom transports are:
1. ✅ **Automatic** - Bindings auto-served, no manual copying
2. ✅ **Simple** - One-line transport override
3. ✅ **Robust** - Built-in reconnect, fallback, error handling
4. ✅ **Flexible** - Plugin system for custom transports
5. ✅ **Type-Safe** - Full TypeScript support with autocomplete

**Next Steps**:
1. Implement binding auto-serving (biggest DX win)
2. Add `--runtime-import` flag to `wails3 generate bindings`
3. Export built-in `createWebSocketTransport()`
4. Export helper utilities (generateMessageID, etc.)
5. Add transport registry system
6. Implement transport fallback mechanism

This will make custom transports a **first-class feature** that's as easy to use as any other Wails API.
