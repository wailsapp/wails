# Event Support with Custom Transports

## Overview

**Yes, events work with custom transports!** This example demonstrates that Wails events can be broadcast over custom transport layers like WebSocket, maintaining full compatibility with the Wails event system.

## How It Works

### Architecture

```
┌─────────────────────────────────────────────────┐
│  Backend (Go)                                   │
│                                                 │
│  GreetService.Greet() calls:                   │
│    app.Events.Emit("greet:count", count)       │
│         ↓                                       │
│  main.go subscribes and broadcasts:            │
│    app.Events.On("greet:count", ...)           │
│    wsTransport.BroadcastEvent(event)           │
│         ↓                                       │
│  WebSocket Message:                            │
│  {                                             │
│    type: "event",                              │
│    event: {                                    │
│      name: "greet:count",                      │
│      data: 5                                   │
│    }                                           │
│  }                                             │
└──────────────────┬──────────────────────────────┘
                   │
                   │ ws://localhost:9099/wails/ws
                   │
┌──────────────────▼──────────────────────────────┐
│  Frontend (JavaScript)                          │
│                                                 │
│  WebSocket Transport receives:                 │
│    handleMessage() sees type="event"           │
│         ↓                                       │
│  window._wails.dispatchWailsEvent(event)       │
│         ↓                                       │
│  Wails Event System dispatches to listeners:   │
│    Events.On("greet:count", callback)          │
│         ↓                                       │
│  UI updates:                                   │
│    eventCounter.textContent = event.data       │
└─────────────────────────────────────────────────┘
```

### Key Components

#### 1. Backend Event Emission (GreetService.go)
```go
func (g *GreetService) Greet(name string) string {
    g.greetCount++

    // Emit event
    if g.app != nil {
        g.app.Events.Emit(&application.WailsEvent{
            Name: "greet:count",
            Data: g.greetCount,
        })
    }

    return result
}
```

#### 2. Event Forwarding to Transport (main.go)
```go
// Subscribe to application events and broadcast them over WebSocket
app.Events.On("greet:count", func(event *application.WailsEvent) {
    log.Printf("[Events] Broadcasting greet:count event: %v", event.Data)
    wsTransport.BroadcastEvent(event)
})
```

#### 3. WebSocket Transport Broadcasting (transport_websocket_example.go)
```go
func (w *WebSocketTransport) BroadcastEvent(event interface{}) error {
    msg := WebSocketMessage{
        Type:  "event",
        Event: event,
    }

    // Send to all connected clients
    for conn := range w.clients {
        conn.WriteJSON(msg)
    }

    return nil
}
```

#### 4. Frontend Event Reception (websocket-transport.js)
```javascript
handleMessage(data) {
    const msg = JSON.parse(data);

    if (msg.type === 'event') {
        // Dispatch to Wails event system
        if (msg.event && window._wails?.dispatchWailsEvent) {
            window._wails.dispatchWailsEvent(msg.event);
        }
    }
}
```

#### 5. Frontend Event Subscription (index.html)
```javascript
import { Events } from '/wails/runtime.js';

// Subscribe to events just like normal!
Events.On('greet:count', (event) => {
    console.log('Received event:', event.data);
    document.getElementById('eventCounter').textContent = event.data;
});
```

## What This Demonstrates

✅ **Events work over custom transports** - Server-to-client push notifications
✅ **Full Wails event API support** - `Events.On()`, `Events.Once()`, `Events.Off()`
✅ **Broadcast to all clients** - WebSocket naturally supports this
✅ **Type-safe event data** - Events carry structured data
✅ **Seamless integration** - No special code needed in service methods

## Differences from Default Transport

### Default (Webview) Transport
- Events sent via `window._wails.dispatchWailsEvent()` through webview's `ExecJS()`
- Only works within webview context
- Per-window delivery

### WebSocket Transport
- Events sent as WebSocket messages
- Works in browser and webview
- Broadcast to all connected clients
- Requires manual forwarding from app events to transport

## Implementation Pattern

For any custom transport to support events:

### Backend
1. **Subscribe to app events**:
   ```go
   app.Events.On("event:name", func(event *application.WailsEvent) {
       transport.BroadcastEvent(event)
   })
   ```

2. **Implement broadcast method**:
   ```go
   func (t *Transport) BroadcastEvent(event interface{}) error {
       // Send event to all clients via your protocol
   }
   ```

### Frontend
1. **Detect event messages**:
   ```javascript
   if (msg.type === 'event') {
       window._wails.dispatchWailsEvent(msg.event);
   }
   ```

2. **Use standard Wails event API**:
   ```javascript
   Events.On('event:name', callback);
   ```

## Testing

Run the example and click "Greet" multiple times:

### Console Output
```
[Events] Broadcasting greet:count event: 1
[WebSocket] Received server event: {type: "event", event: {...}}
[WebSocket] Event dispatched to Wails: greet:count
[Events] Received greet:count event: {name: "greet:count", data: 1}
```

### UI Behavior
- **Event Counter** updates in real-time
- Counter increments with each greet
- All clients see the same count (broadcast)

## Benefits

### 1. Real-Time Updates
Server can push updates to clients without polling:
- Status changes
- Progress updates
- Notifications
- Live data

### 2. Multi-Client Sync
WebSocket broadcast means all clients stay synchronized:
- Collaborative features
- Shared state
- Live dashboards

### 3. Bidirectional Communication
Events can flow both ways:
- Client → Server: `Events.Emit()`
- Server → Client: `app.Events.Emit()` + broadcast

### 4. Clean API
No transport-specific code in business logic:
```go
// This works the same regardless of transport!
app.Events.Emit(&application.WailsEvent{
    Name: "data:update",
    Data: newData,
})
```

## Limitations & Considerations

### Current Implementation
- **Manual forwarding required**: Must subscribe to events and forward to transport
- **All clients receive all events**: No per-client filtering (can be added)
- **No guaranteed delivery**: WebSocket disconnect = missed events

### Production Enhancements
1. **Event filtering**: Only send relevant events to each client
2. **Event buffering**: Queue events during disconnect
3. **Selective broadcast**: Target specific clients/windows
4. **Event acknowledgment**: Confirm delivery

## Comparison with HTTP Transport

| Feature | HTTP Transport | WebSocket Transport |
|---------|---------------|---------------------|
| **Binding Calls** | ✅ Request/Response | ✅ Request/Response |
| **Events** | ✅ Via webview ExecJS | ✅ Via WebSocket message |
| **Server Push** | ❌ Not possible | ✅ Native support |
| **Multi-Client** | N/A (single webview) | ✅ Broadcast |
| **Browser Support** | ❌ Webview only | ✅ Works in browser |

## Example Use Cases

### 1. Progress Tracking
```go
for i := 0; i <= 100; i++ {
    app.Events.Emit(&application.WailsEvent{
        Name: "progress:update",
        Data: i,
    })
    time.Sleep(100 * time.Millisecond)
}
```

### 2. Live Notifications
```go
app.Events.Emit(&application.WailsEvent{
    Name: "notification:new",
    Data: map[string]interface{}{
        "title": "New Message",
        "body": "You have a new message",
    },
})
```

### 3. Multi-User Sync
```go
// When user A makes a change
app.Events.Emit(&application.WailsEvent{
    Name: "document:updated",
    Data: document,
}) // All connected users see the update
```

## Conclusion

**Custom transports fully support Wails events!** The implementation requires:
- Backend: Forward app events to transport
- Transport: Broadcast events to clients
- Frontend: Dispatch to `window._wails.dispatchWailsEvent()`

The result is seamless integration where:
- Service methods emit events normally
- Frontend subscribes with `Events.On()`
- No transport-specific code in business logic
- Full compatibility with Wails event system

Events over WebSocket enable real-time, bidirectional communication while maintaining the clean Wails event API! 🎉
