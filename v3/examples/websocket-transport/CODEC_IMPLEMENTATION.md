# Transport Codec Implementation

## What Was Implemented

We added a pluggable codec system that allows developers to customize how data is encoded/decoded between the frontend and backend transport layers.

### Files Created/Modified

#### Frontend (TypeScript/JavaScript)
1. **Created**: `v3/internal/runtime/desktop/@wailsio/runtime/src/transport-codec.ts`
   - `TransportCodec` interface
   - `Base64JSONCodec` (default - matches Go's []byte base64 encoding)
   - `RawStringCodec` (plain strings)
   - `RawJSONCodec` (direct JSON objects)

2. **Modified**: `v3/internal/runtime/desktop/@wailsio/runtime/src/index.ts`
   - Exported codec types from runtime

3. **Modified**: `v3/examples/websocket-transport/assets/websocket-transport.js`
   - Added codec support with `options.codec` parameter
   - Uses codec to decode responses and errors

4. **Modified**: `v3/examples/websocket-transport/assets/index.html`
   - Added example of how to configure codec

#### Backend (Go)
1. **Created**: `v3/pkg/application/transport_codec.go`
   - `TransportCodec` interface
   - `DefaultCodec` (base64/JSON - default)
   - `RawJSONCodec` (direct JSON)
   - `RawStringCodec` (plain strings)

2. **Modified**: `v3/pkg/application/transport.go`
   - Changed `TransportResponse.Data` from `[]byte` to `interface{}`

3. **Modified**: `v3/pkg/application/transport_websocket_example.go`
   - Added `codec` field and `WithCodec()` option
   - Uses codec in `handleRequest()` to encode/decode data

4. **Modified**: `v3/examples/websocket-transport/main.go`
   - Added documentation about codec usage

## How to Verify It Works

### Method 1: Check Console Logs

Run the example and check browser console:

```bash
cd v3/examples/websocket-transport
go run .
```

**In browser console, you should see:**
```
[WebSocket Transport] Loading VERSION 4 with codec support
✓ WebSocket transport configured with codec: Base64JSONCodec
```

### Method 2: Test with Default Codec (Base64JSON)

1. Run the example
2. Click "Greet" button
3. **Expected behavior**: UI updates with greeting message
4. **Console should show**:
   ```
   [WebSocket] Response data: "Hello, WebSocket User! (Greeted X times via WebSocket)"
   [WebSocket] Content type: application/json
   [WebSocket] Calling Wails result handler with JSON data
   ```

### Method 3: Test with RawJSONCodec

**Frontend** - Edit `assets/index.html` line 220:
```javascript
const wsTransport = createWebSocketTransport('ws://localhost:9099/wails/ws', {
    reconnectDelay: 2000,
    requestTimeout: 30000,
    codec: new RawJSONCodec()  // Use raw JSON instead of base64
});
```

**Backend** - Edit `main.go` line 17:
```go
wsTransport := application.NewWebSocketTransport(":9099",
    application.WithCodec(application.NewRawJSONCodec()))
```

**Rebuild and run:**
```bash
go run .
```

**Expected difference**:
- Response data will be direct JSON object instead of base64 string
- No base64 decoding step in browser console

### Method 4: Verify Codec Interface

**Test custom codec** - Add to `assets/index.html`:

```javascript
// Custom codec that logs everything
class LoggingCodec {
    decodeResponse(data, contentType) {
        console.log('[CustomCodec] Decoding response:', data, 'type:', contentType);
        return new Base64JSONCodec().decodeResponse(data, contentType);
    }

    decodeError(data) {
        console.log('[CustomCodec] Decoding error:', data);
        return new Base64JSONCodec().decodeError(data);
    }
}

const wsTransport = createWebSocketTransport('ws://localhost:9099/wails/ws', {
    codec: new LoggingCodec()
});
```

**Expected**: You'll see `[CustomCodec]` logs showing all decode operations.

### Method 5: Network Inspection

**Open browser DevTools → Network → WS tab**

1. Click "Greet" button
2. Select the WebSocket connection
3. Click on a message frame

**With DefaultCodec, you'll see**:
```json
{
  "id": "abc123...",
  "type": "response",
  "response": {
    "statusCode": 200,
    "contentType": "application/json",
    "data": "SGVsbG8sIFdlYlNvY2tldCBVc2VyISAoR3JlZXRlZCAx..."  // base64 string
  }
}
```

**With RawJSONCodec, you'll see**:
```json
{
  "id": "abc123...",
  "type": "response",
  "response": {
    "statusCode": 200,
    "contentType": "application/json",
    "data": "Hello, WebSocket User! (Greeted 1 times via WebSocket)"  // direct string
  }
}
```

## Success Criteria

✅ **Working correctly if:**
1. Browser console shows "VERSION 4 with codec support"
2. Console shows "codec: Base64JSONCodec" (or your chosen codec)
3. Greet button updates UI with greeting message
4. No decode errors in console
5. Backend logs show successful response encoding

❌ **Not working if:**
1. Console shows codec errors like "Failed to decode response"
2. UI doesn't update after clicking buttons
3. Browser shows "VERSION 3" or earlier
4. Response data is empty or garbled

## Codec Selection Guide

| Use Case | Frontend Codec | Backend Codec | Why |
|----------|---------------|---------------|-----|
| Default (recommended) | `Base64JSONCodec` | `DefaultCodec` | Matches Go's JSON behavior, works out of box |
| Performance optimization | `RawJSONCodec` | `RawJSONCodec` | Skips base64 encoding/decoding |
| Text-only transport | `RawStringCodec` | `RawStringCodec` | Minimal overhead for plain text |
| Binary protocol | Custom implementation | Custom implementation | Full control over format |

## Common Issues

### Issue: "Failed to decode response"
**Cause**: Frontend/backend codec mismatch
**Fix**: Ensure both sides use compatible codecs (e.g., Base64JSONCodec ↔ DefaultCodec)

### Issue: UI doesn't update
**Cause**: Runtime not rebuilt after codec changes
**Fix**: Rebuild runtime with `npm run build` in runtime directory

### Issue: Response data is object instead of string
**Cause**: Using RawJSONCodec on frontend but DefaultCodec on backend
**Fix**: Match codecs on both sides

## Benefits Achieved

**Before (manual base64 decoding):**
```javascript
// 26 lines of boilerplate × 2 places = 52 lines
const binaryString = atob(response.data);
const bytes = new Uint8Array(binaryString.length);
for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i);
}
responseData = new TextDecoder().decode(bytes);
```

**After (with codec):**
```javascript
// 1 line
const responseData = this.codec.decodeResponse(response.data, response.contentType);
```

**Reduction**: ~51 lines of boilerplate eliminated, plus flexibility to use any encoding scheme.
