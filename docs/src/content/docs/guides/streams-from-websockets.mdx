---
title: Migrating a WebSocket to Streams
description: Step-by-step conversion of an existing WebSocket implementation to Wails streams, including the differences that break silently
sidebar:
    order: 55
---

A mechanical guide for converting an app that currently runs a WebSocket server to
[Streams](/guides/streams/). Written to be followed literally, including by an agent.

The payoff is deleting the listener: no bound TCP port, no origin check, no token, nothing
for a firewall or endpoint-security product to object to. The API is close enough that most
frontend code is untouched — but **three differences break silently**, and they are listed
first because they are the ones that will cost you an afternoon.

## The common case: a local HTTP server as a workaround

The usual reason a Wails app has a WebSocket is that there was no other way to push a
continuous feed to the frontend, so the app stands up its own `http.Server` on a local
port and the frontend connects back to it. If that is your shape, this migration deletes
the server outright — and several things you built *around* it go with it.

**The port discovery mechanism goes.** Something has to tell the frontend which port to
connect to: a bound `GetServerPort()`, a fixed port with a fallback when it is taken, an
injected global, or a value stashed in `localStorage`. All of it disappears — a stream is
addressed by name, and the name is a compile-time constant on both sides.

```go
// BEFORE
ln, _ := net.Listen("tcp", "127.0.0.1:0")
go http.Serve(ln, mux)
port := ln.Addr().(*net.TCPAddr).Port     // ...and a binding to hand `port` to the frontend

// AFTER
app.HandleStream("feed", handler)         // that is the entire replacement
```

**CORS configuration goes.** The webview's origin is `wails://` or
`http://wails.localhost` depending on platform, so a local server needs `CheckOrigin`, an
`Access-Control-Allow-Origin` header, or both. Streams ride the asset server the page was
loaded from, so there is no cross-origin request to permit.

**Any auth token you invented goes.** A port bound on localhost is reachable by every
process on the machine, so a careful implementation adds a token or a nonce to stop other
software connecting. There is no port to reach any more.

**Non-WebSocket endpoints move to asset server middleware.** These servers rarely stay
pure — a file download, an image endpoint, a health check tend to accumulate alongside the
socket. Streams do not replace those, but you do not need a second server for them either.
Mount the same handlers on the asset server:

```go
app := application.New(application.Options{
    Assets: application.AssetOptions{
        Handler: yourFrontendAssets,
        Middleware: func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                if strings.HasPrefix(r.URL.Path, "/api/") {
                    yourExistingMux.ServeHTTP(w, r)   // the handlers you already wrote
                    return
                }
                next.ServeHTTP(w, r)
            })
        },
    },
})
```

The frontend then calls `/api/...` as a same-origin relative URL — no host, no port, no
CORS. Between streams and this, the local server has nothing left to do.

## Read this before you start

### 1. `ev.data` is an `ArrayBuffer`, never a string

This is the big one. A WebSocket delivers text messages as strings; a stream delivers every
message as bytes. Code like this **compiles, runs, and misbehaves**:

```js
// BEFORE — works with WebSocket
ws.onmessage = (ev) => { const msg = JSON.parse(ev.data); ... };

// AFTER — ev.data is an ArrayBuffer, JSON.parse gets "[object ArrayBuffer]"
```

Fix it at the boundary rather than at every call site:

```js
const dec = new TextDecoder();
s.onmessage = (ev) => { const msg = JSON.parse(dec.decode(ev.data)); ... };
```

Or, if your traffic is JSON — which it usually is — use `JSONStream` instead of `Stream`
and skip the problem entirely:

```js
import { JSONStream } from "@wailsio/runtime";

const s = JSONStream("feed");
s.onmessage = (ev) => dispatch(ev.data);   // already an object
s.send({ hello: true });                   // stringified for you
```

That is the shortest migration for a JSON WebSocket: swap the constructor, drop the
`JSON.parse` and `JSON.stringify` calls, and the rest of the handler is unchanged. For
non-JSON traffic, wrap once and leave every handler untouched — see
[the compatibility shim](#compatibility-shim).

### 2. Sending is compatible; receiving is not

`send()` accepts a string and encodes it as UTF-8, so `s.send(JSON.stringify(x))` works
unchanged. Only the receive path needs editing. The asymmetry is easy to miss because half
your code keeps working.

### 3. There is no URL

A WebSocket carries connection parameters in its URL — path, query string, subprotocol,
auth token. A stream has only a name. Anything you passed in the URL has to move into the
first frame, or into a bound method called before connecting.

```js
// BEFORE
const ws = new WebSocket(`wss://host/feed?topic=${topic}&token=${token}`);

// AFTER — no token needed at all; the app is the only possible caller
const s = Stream("feed");
s.onopen = () => s.send(JSON.stringify({ subscribe: topic }));
```

## Go side

Delete the HTTP server, the upgrader, and the connection registry. Each becomes a handler.

```go
// BEFORE — gorilla/coder websocket
func (a *App) serveWS(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    clients.add(conn)
    defer clients.remove(conn)

    for {
        _, data, err := conn.ReadMessage()
        if err != nil {
            return
        }
        handle(data)
    }
}

// ...plus http.ListenAndServe, a mux entry, an origin checker, and a token check
```

```go
// AFTER
app.HandleStream("feed", func(c *application.StreamConn) {
    defer c.Close()

    for {
        frame, err := c.Receive()
        if err != nil {
            return                  // reload, close, or shutdown
        }
        handle(frame)
    }
})
```

| WebSocket | Stream |
|---|---|
| `upgrader.Upgrade` / `websocket.Accept` | *(nothing — `HandleStream` is the whole registration)* |
| `conn.ReadMessage()` / `conn.Read(ctx)` | `c.Receive()` |
| `conn.WriteMessage(TextMessage, b)` | `c.Send(b)` |
| `conn.Close()` | `c.Close()`, or just return from the handler |
| connection registry for broadcast | keep your own — see [Broadcasting](#broadcasting) |
| `http.ListenAndServe`, mux, origin check, token | **delete** |
| ping/pong keepalive | **delete** — there is no idle socket to keep alive |
| `r.Context()` | `c.Context()` |

The handler goroutine's lifetime is the connection's lifetime, exactly like a gorilla
handler, so the structure of your existing loop carries over unchanged.

## Frontend side

```js
// BEFORE
const ws = new WebSocket(url);
ws.onopen    = () => ws.send(JSON.stringify(hello));
ws.onmessage = (ev) => dispatch(JSON.parse(ev.data));
ws.onclose   = () => scheduleReconnect();

// AFTER
import { Stream } from "@wailsio/runtime";
const dec = new TextDecoder();

const s = Stream("feed");
s.onopen    = () => s.send(JSON.stringify(hello));   // unchanged
s.onmessage = (ev) => dispatch(JSON.parse(dec.decode(ev.data)));
s.onclose   = () => scheduleReconnect();             // unchanged
```

`readyState`, the four state constants, `addEventListener`, `close(code, reason)` and
`bufferedAmount` all behave as they do on a `WebSocket`.

### Compatibility shim

If you would rather not touch handlers at all, wrap the constructor once. Existing code
that expects string messages then works verbatim:

```js
import { Stream } from "@wailsio/runtime";

/** A Stream that delivers text messages as strings, like a WebSocket. */
export function TextStream(name) {
    const s = Stream(name);
    const dec = new TextDecoder();
    s.binaryType = "arraybuffer";

    const add = s.addEventListener.bind(s);
    const remove = s.removeEventListener.bind(s);
    const wrappers = new WeakMap();
    const decoded = new WeakMap();

    const decodeEvent = (ev) => {
        if (decoded.has(ev)) return decoded.get(ev);
        const data = typeof ev.data === "string" ? ev.data : dec.decode(ev.data);
        const textEvent = new MessageEvent("message", { data });
        decoded.set(ev, textEvent);
        return textEvent;
    };

    const wrap = (listener) => {
        let wrapper = wrappers.get(listener);
        if (wrapper) return wrapper;
        wrapper = (ev) => {
            const textEvent = decodeEvent(ev);
            if (typeof listener === "function") listener.call(s, textEvent);
            else listener.handleEvent(textEvent);
        };
        wrappers.set(listener, wrapper);
        return wrapper;
    };

    s.addEventListener = (type, listener, options) =>
        add(type, type === "message" && listener ? wrap(listener) : listener, options);
    s.removeEventListener = (type, listener, options) =>
        remove(type, type === "message" && listener ? wrappers.get(listener) ?? listener : listener, options);

    // A WailsSocket implements onmessage through addEventListener, but a native
    // WebSocket uses an internal event-handler slot. Define the property on the
    // instance so both transports pass property handlers through the same
    // decoding wrapper as addEventListener listeners.
    let onmessage = null;
    Object.defineProperty(s, "onmessage", {
        get: () => onmessage,
        set(listener) {
            if (onmessage) s.removeEventListener("message", onmessage);
            onmessage = typeof listener === "function" ? listener : null;
            if (onmessage) s.addEventListener("message", onmessage);
        },
        configurable: true,
        enumerable: true,
    });
    return s;
}
```

Then `const ws = TextStream("feed")` is a drop-in for `new WebSocket(url)`.

## Broadcasting

A WebSocket server usually keeps a registry so it can fan out. Streams have no built-in
broadcast — keep the registry, just store `*StreamConn` instead of `*websocket.Conn`:

```go
type hub struct {
    mu    sync.Mutex
    conns map[*application.StreamConn]struct{}
}

func (h *hub) add(c *application.StreamConn)    { h.mu.Lock(); h.conns[c] = struct{}{}; h.mu.Unlock() }
func (h *hub) remove(c *application.StreamConn) { h.mu.Lock(); delete(h.conns, c); h.mu.Unlock() }

func (h *hub) broadcast(msg []byte) {
    h.mu.Lock()
    conns := make([]*application.StreamConn, 0, len(h.conns))
    for c := range h.conns {
        conns = append(conns, c)
    }
    h.mu.Unlock()                         // never hold the lock across Send

    for _, c := range conns {
        // TrySend, not Send: one stalled frontend must not block the fan-out.
        _ = c.TrySend(msg)
    }
}

app.HandleStream("feed", func(c *application.StreamConn) {
    h.add(c)
    defer h.remove(c)
    defer c.Close()
    <-c.Context().Done()
})
```

Two rules worth keeping: release the lock before sending, and prefer `TrySend` in a
fan-out so a single slow consumer cannot stall every other client.

## The other case: the frontend talks to a broker directly

Less common in a Wails app, but worth knowing. If the frontend opens a WebSocket **to a
broker rather than to your app** — `nats.ws` to a NATS server, MQTT over WebSocket — a
stream is not a drop-in, because a stream connects the frontend to *your Go code*, not to a
third party.

The migration is an architecture change, and usually a good one:

```
BEFORE   frontend ──ws──► NATS server            (credentials in the frontend)
AFTER    frontend ──stream──► Go ──►  NATS       (credentials stay in Go)
```

Move the broker client into Go, where the native library is better than the browser one,
and expose the parts the frontend needs over a stream:

```go
nc, _ := nats.Connect(url, nats.UserCredentials(credsPath))  // creds never reach the frontend

app.HandleStream("nats", func(c *application.StreamConn) {
    defer c.Close()

    var subs []*nats.Subscription
    defer func() {
        for _, s := range subs {
            _ = s.Unsubscribe()
        }
    }()

    for {
        frame, err := c.Receive()
        if err != nil {
            return
        }

        var cmd struct {
            Op      string          `json:"op"`       // "sub" | "pub"
            Subject string          `json:"subject"`
            Data    json.RawMessage `json:"data"`
        }
        if json.Unmarshal(frame, &cmd) != nil {
            continue
        }

        switch cmd.Op {
        case "sub":
            sub, err := nc.Subscribe(cmd.Subject, func(m *nats.Msg) {
                out, _ := json.Marshal(map[string]any{"subject": m.Subject, "data": m.Data})
                // TrySend: a slow frontend must not block the NATS callback.
                _ = c.TrySend(out)
            })
            if err == nil {
                subs = append(subs, sub)
            }
        case "pub":
            _ = nc.Publish(cmd.Subject, cmd.Data)
        }
    }
})
```

Note `TrySend` inside the subscription callback: that callback runs on the broker client's
goroutine, and blocking it would stall delivery for every subscription on the connection.

What you gain: broker credentials never reach the frontend, no WebSocket port exposed to
the machine, and reconnection/backoff handled by the mature Go client instead of the
browser one.

## Migration checklist

- [ ] `HandleStream` registered for each WebSocket endpoint you had
- [ ] Read loop converted: `ReadMessage`/`Read` → `c.Receive()`
- [ ] Writes converted: `WriteMessage` → `c.Send()`, or `TrySend` in any fan-out or broker callback
- [ ] HTTP server, mux entry, upgrader, origin check and auth token **deleted**
- [ ] Ping/pong keepalive **deleted**
- [ ] URL parameters moved into a first frame or a bound method
- [ ] **Every `ev.data` read decoded** — `new TextDecoder().decode(ev.data)` — or the shim adopted
- [ ] Reconnect logic kept as-is (there is no built-in reconnect, by design)
- [ ] Broker clients moved into Go if the frontend was talking to one directly
- [ ] Checked for `InitialHTML` windows — they cannot use streams at all

Things worth grepping for during conversion:

```
new WebSocket(     ev.data            .onmessage
websocket.Accept   upgrader.Upgrade   ReadMessage
WriteMessage       ListenAndServe     CheckOrigin
```

## Behaviour differences to expect

| | WebSocket | Stream |
|---|---|---|
| Message type | text or binary | bytes only |
| `ev.data` | string or `Blob`/`ArrayBuffer` | always `ArrayBuffer` (unless `binaryType = "blob"`) |
| `binaryType` default | `"blob"` | `"arraybuffer"` |
| Subprotocols, `extensions` | negotiated | not supported; always `""` |
| Connection parameters | URL and query | first frame, or a bound call |
| Auth | token or cookie | none needed — the app is the only caller |
| Keepalive | ping/pong | none required |
| Automatic reconnect | none | none (same) |
| Close codes | full range | `1000` normal, `1001` session closed, `1002` framing mismatch, `1006` error |
| Backpressure | kernel socket buffer | 8 MB / 256 frames per window, then `Send` blocks |
| Many connections, one endpoint | yes | yes |

## After migrating

Sanity checks that catch the common mistakes:

1. Reload the page repeatedly — the handler should exit and a fresh one start each time,
   never accumulating.
2. Send a message larger than 512 KB in each direction.
3. Leave it idle for over a minute; traffic should resume without a reconnect.
4. Open devtools and pause on a breakpoint under load, then resume — the producer should
   block and recover rather than lose data or grow without bound.
