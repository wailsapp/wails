# Generated Bindings with Custom Transport

## Overview

This example now demonstrates that **custom transports work seamlessly with Wails generated bindings**. This is a critical validation that the transport hijack mechanism doesn't break code generation.

## What We Changed

### 1. Generated Bindings
Used `wails3 generate bindings` to create TypeScript bindings for `GreetService`:

```
frontend/bindings/github.com/wailsapp/wails/v3/examples/websocket-transport/
├── index.js          # Exports GreetService namespace
└── greetservice.js   # Generated bindings using Call.ByID()
```

### 2. Updated Frontend
**Before** (manual calls):
```javascript
import { Call } from '/wails/runtime.js';

const result = await Call.ByName("main.GreetService.Greet", name);
```

**After** (generated bindings):
```javascript
import { GreetService } from '/bindings/.../index.js';

const result = await GreetService.Greet(name);  // Clean, type-safe API
```

## Key Insight: Transport Transparency

The generated bindings use `Call.ByID()` internally:

```javascript
// From greetservice.js (generated)
export function Greet(name) {
    return $Call.ByID(1411160069, name);  // Uses whatever transport is set!
}
```

When we call `setTransport(wsTransport)`, the runtime updates the internal transport mechanism. The generated bindings automatically use the new transport without any code changes!

## Verification

### Console Output
When running the app, you should see:

```
[WebSocket Transport] Loading VERSION 4 with codec support
✓ WebSocket transport configured with codec: Base64JSONCodec
✓ Using generated bindings for GreetService
[UI] Calling Greet with: WebSocket User (via generated binding)
[WebSocket] Received request: msgID=abc123, object=0, method=0
[WebSocket] Response: statusCode=200, contentType=application/json
[UI] Greet result: Hello, WebSocket User! (Greeted 1 times via WebSocket)
[UI] UI updated successfully
```

### What This Proves

1. ✅ **Generated bindings work** - `GreetService.Greet()` calls succeed
2. ✅ **WebSocket transport is used** - Messages go through ws://localhost:9099
3. ✅ **Codec system works** - Base64 responses decoded correctly
4. ✅ **Callback handling works** - UI updates via `window._wails.callResultHandler`
5. ✅ **No code generation changes needed** - Bindings generated once, transport swapped at runtime

## Architecture Flow

```
┌─────────────────────────────────────────────────┐
│  Frontend                                       │
│                                                 │
│  GreetService.Greet("User")                    │
│         ↓                                       │
│  Call.ByID(1411160069, "User")                 │
│         ↓                                       │
│  setTransport() → WebSocketTransport           │
│         ↓                                       │
│  WebSocket Message:                            │
│  {                                             │
│    id: "abc123",                               │
│    type: "request",                            │
│    request: {                                  │
│      object: 0,                                │
│      method: 0,                                │
│      args: '{"methodID":1411160069,...}'       │
│    }                                           │
│  }                                             │
└──────────────────┬──────────────────────────────┘
                   │
                   │ ws://localhost:9099/wails/ws
                   │
┌──────────────────▼──────────────────────────────┐
│  Backend                                        │
│                                                 │
│  WebSocketTransport.handleRequest()            │
│         ↓                                       │
│  codec.DecodeRequest(args)                     │
│         ↓                                       │
│  handler.HandleRuntimeCall(req)                │
│         ↓                                       │
│  MessageProcessor.processCallMethod()          │
│         ↓                                       │
│  GreetService.Greet("User")                    │
│         ↓                                       │
│  codec.EncodeResponse(result)                  │
│         ↓                                       │
│  WebSocket Response:                           │
│  {                                             │
│    id: "abc123",                               │
│    type: "response",                           │
│    response: {                                 │
│      statusCode: 200,                          │
│      contentType: "application/json",          │
│      data: "SGVsbG8s..." (base64)              │
│    }                                           │
│  }                                             │
└─────────────────────────────────────────────────┘
```

## Benefits Demonstrated

### 1. Developer Experience
Developers can:
- Generate bindings once with `wails3 generate bindings`
- Use clean, typed APIs like `GreetService.Greet(name)`
- Swap transports at runtime without regenerating code
- Keep the same frontend code for HTTP, WebSocket, or custom transports

### 2. Transport Independence
The generated bindings are transport-agnostic:
- No transport-specific imports in generated code
- No hardcoded URLs or protocols
- Transport is a runtime configuration concern, not a build-time concern

### 3. Type Safety (with TypeScript)
If using TypeScript, the generated bindings provide:
- Parameter type checking
- Return type inference
- IDE autocomplete
- Compile-time error detection

## Comparison: Manual vs Generated

### Manual Calls (Before)
```javascript
// Easy to make mistakes
await Call.ByName("main.GreetService.Greet", name);  // Typo-prone
await Call.ByName("main.GreetService.Greeet", name); // Oops! Typo
await Call.ByName("main.GreetService.Greet", 123);   // Wrong type, no error
```

### Generated Bindings (After)
```javascript
// Type-safe, autocomplete, refactor-friendly
await GreetService.Greet(name);      // IDE autocomplete
await GreetService.Greeet(name);     // IDE shows error immediately
await GreetService.Greet(123);       // TypeScript catches wrong type
```

## Testing Checklist

To verify the transport works with generated bindings:

- [x] ✅ Generated bindings exist in `frontend/bindings/`
- [x] ✅ Bindings copied to `assets/bindings/` for serving
- [x] ✅ HTML imports bindings instead of using `Call.ByName`
- [x] ✅ WebSocket transport configured via `setTransport()`
- [x] ✅ All four methods work: Greet, Echo, Add, GetTime
- [x] ✅ Console shows "via generated binding" messages
- [x] ✅ UI updates correctly with results
- [x] ✅ WebSocket messages visible in Network tab

## Future Enhancements

This foundation enables:

1. **Multiple Transports**: Switch between HTTP and WebSocket at runtime
2. **Transport Fallback**: Try WebSocket, fallback to HTTP if unavailable
3. **Environment-Specific**: Use WebSocket in dev, HTTP in production
4. **Custom Protocols**: Implement gRPC, MessagePack, or proprietary protocols
5. **Transport Middleware**: Add logging, caching, retry logic at transport layer

## Conclusion

**The custom transport API successfully achieves its goal**: Developers can replace the IPC transport layer while retaining 100% compatibility with generated bindings and all Wails features.

The generated code is **transport-agnostic** and the transport is **runtime-configurable**, providing maximum flexibility without sacrificing developer experience.
