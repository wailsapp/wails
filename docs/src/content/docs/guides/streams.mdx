---
title: Streams
description: Bidirectional byte streams between Go and JavaScript, with the WebSocket programming model and no listening socket
sidebar:
    order: 53
---

Streams give you a named, ordered, bidirectional byte channel between Go and your
frontend, with the same programming model as a WebSocket — **without binding a TCP port**.

A WebSocket cannot be spoken over a custom URL scheme, so the only way to get one inside a
webview is to run a real HTTP server and listen on a port. On a desktop app that means an
open local port that any other process on the machine can reach, needing an origin check
and a token to be safe, and visible to every firewall and endpoint-security product your
users run. Streams avoid all of that: they ride the asset server your app already serves,
which is already origin-bound.

Migrating an existing WebSocket implementation? Follow
[Migrating a WebSocket to Streams](/guides/streams-from-websockets/) — it is written to be
applied mechanically, and it leads with the three differences that break silently.

## Quick start

Declare a stream in Go. The handler runs once per connection, on its own goroutine:

```go
app.HandleStream("telemetry", func(c *application.StreamConn) {
    defer c.Close()

    for {
        frame, err := c.Receive()   // blocks until a frame arrives
        if err != nil {
            return                  // page reloaded, window closed, or app shutting down
        }
        _ = c.Send(process(frame))  // blocks like a socket write
    }
})
```

Connect from the frontend by name. The object implements the `WebSocket` interface:

```js
import { Stream } from "@wailsio/runtime";

const s = Stream("telemetry");
s.onopen    = () => s.send(new TextEncoder().encode("hello"));
s.onmessage = (ev) => console.log(new Uint8Array(ev.data));
s.onclose   = (ev) => console.log("closed", ev.code);
```

`Stream(name)` returns **synchronously** with `readyState === CONNECTING`, exactly like
`new WebSocket(url)`, so you can create one at module scope:

```js
export const Telemetry = Stream("telemetry");
```

## Frames are bytes

Every frame is a `[]byte` in Go and an `ArrayBuffer` in JavaScript. There is no schema and
no encoding imposed on you — use JSON, protobuf, CBOR, or raw bytes as you prefer.

A frame is a **message, not a byte stream**: it arrives whole or not at all, and its length
travels with it. Neither side needs to know the size in advance, so a struct with a
`[]byte` field marshals to whatever it marshals to and is sent as one frame.

## Sending objects

Frames are bytes, but you rarely want to think in bytes. Both sides have a JSON
convenience that pairs up:

```go
type Reading struct {
    Sensor string  `json:"sensor"`
    Value  float64 `json:"value"`
}

app.HandleStream("telemetry", func(c *application.StreamConn) {
    defer c.Close()

    var cmd map[string]any
    if err := c.ReceiveJSON(&cmd); err != nil {
        return
    }

    _ = c.SendJSON(Reading{Sensor: "cpu", Value: 42.5})
})
```

```js
import { JSONStream } from "@wailsio/runtime";

const s = JSONStream("telemetry");
s.onopen    = () => s.send({ subscribe: "cpu" });   // stringified for you
s.onmessage = (ev) => console.log(ev.data.value);   // already an object
```

`JSONStream` is the same object as `Stream` with the encoding done at the boundary — there
is no separate protocol, and a Go handler cannot tell the difference. A frame that is not
valid JSON raises an `error` event and is dropped rather than taking the connection down.

Use plain `Stream` when you want the bytes: protobuf, CBOR, binary formats, or anything
where you would rather encode yourself.

## The Go API

```go
// Register a handler. Runs once per connection, on its own goroutine.
func (a *App) HandleStream(name string, handler StreamHandler)

// The connection.
func (c *StreamConn) Send(data []byte) error        // blocks when the buffer is full
func (c *StreamConn) TrySend(data []byte) error     // ErrStreamFull instead of blocking
func (c *StreamConn) Receive() ([]byte, error)      // blocks until a frame or close
func (c *StreamConn) SendJSON(v any) error          // marshal and send as one frame
func (c *StreamConn) ReceiveJSON(v any) error       // receive one frame and unmarshal
func (c *StreamConn) Context() context.Context      // cancelled on disconnect
func (c *StreamConn) Window() Window                // nil in server mode
func (c *StreamConn) Name() string
func (c *StreamConn) Close() error
```

**The handler goroutine's lifetime is the connection's lifetime.** Returning from the
handler closes the connection, so block on `Receive` (or on `c.Context()`) for as long as
you want it open. This is the same shape as a `gorilla`/`coder` WebSocket handler.

Errors are `ErrStreamClosed` (the peer is gone) and `ErrStreamFull` (from `TrySend` only).

## The JavaScript API

`Stream(name)` returns an object implementing the useful subset of `WebSocket`:

| supported | notes |
|---|---|
| `readyState` + `CONNECTING`/`OPEN`/`CLOSING`/`CLOSED` | |
| `onopen`, `onmessage`, `onclose`, `onerror` | plus `addEventListener` |
| `send(data)` | string, `ArrayBuffer`, typed array, or `Blob` — see ownership below |
| `JSONStream(name)` | same object, objects in and out |
| `close(code, reason)` | |
| `binaryType` | **defaults to `"arraybuffer"`**, not `"blob"` |
| `bufferedAmount` | bytes queued by `send` that have not reached Go |
| `protocol`, `extensions` | always `""` — not negotiated |

The `binaryType` default is the one deliberate divergence from the standard: frames are
always binary, and a `Blob` would force an extra async hop to read every message. Set it to
`"blob"` if you want the standard behaviour.

Many connections to the same stream name are allowed, from one window or several. Each
gets its own `StreamConn` and its own handler goroutine.

**Buffer ownership differs by direction.** JavaScript `send()` snapshots mutable binary
inputs synchronously, matching native WebSocket behaviour, so callers may reuse them as
soon as `send()` returns. Go `Send` transfers ownership of its slice to the transport and
does not copy it; do not mutate or reuse that slice after a successful call. Hand over a
fresh slice when the producer needs to reuse its storage.

## Lifecycle

A stream behaves like a socket, and the events that close a socket close a stream:

| event | what happens |
|---|---|
| Page reload or navigation | connection closes, handler's `Receive` returns an error, the new page connects fresh |
| `window.close()` / window destroyed | all that window's connections close |
| `s.close()` in JS | handler's `Receive` returns `ErrStreamClosed` |
| Handler returns | frontend sees `onclose` |
| App shutdown | every connection's context is cancelled |

There is **no automatic reconnect**, which matches `WebSocket`. If your app needs one, the
reconnect logic you already have for a WebSocket will work unchanged — recreate the stream
in `onclose`.

## Backpressure

`Send` blocks when the frontend has not kept up, the way a socket write blocks on a full
send buffer. Use `TrySend` if you would rather drop than wait:

```go
if err := c.TrySend(sample); errors.Is(err, application.ErrStreamFull) {
    // frontend is behind — skip this sample rather than stalling the producer
}
```

A paused frontend — a devtools breakpoint, a hidden window, App Nap — stops collecting, and
the buffer bound then blocks the producer. That is intentional: it bounds memory rather
than letting an unread stream grow without limit.

## Server mode

Building with `-tags server` swaps the transport for a **real WebSocket** at
`/wails/stream/ws`, since server mode already has a listener to upgrade on. The Go handler
and the frontend code are identical — nothing in your app changes. The runtime picks the
transport for you before any module code runs. WebSocket connections are same-origin by
default. A server that deliberately hosts its frontend on another trusted origin can add
that host with `ServerOptions.WebSocketOriginPatterns`.

## Performance

Measured with `v3/tests/stream-performance`, unthrottled, 0 drops and 0 reorders across
~41 million frames:

| | Go→JS peak | JS→Go peak |
|---|---:|---:|
| macOS / WebKit-Cocoa | **3,117 MB/s** | 2,793 MB/s |
| Linux / WebKitGTK | 226 MB/s | 727 MB/s |
| Windows / WebView2 | 100 MB/s | 99 MB/s |

Shape matters more than the peaks:

- **Go→JS is much faster for small frames** — 634,000 frames/s on macOS against ~6,200/s
  the other way. A single response coalesces up to 256 frames; JS→Go also batches frames
  that accumulate behind an in-flight request, but each connection still serialises its
  own POST chain. If you are sending many small messages, prefer Go→JS, or batch them at
  the application level before sending up.
- **On Windows, 512 KB is the sweet spot for uploads.** Frames larger than that are split
  into several requests, and a 4 MB frame measures *slower* than a 512 KB one.
- **Latency is low and stays low**: ~1–2 ms p99 on macOS, and it does not degrade with
  load — 20,000 frames/s measured a *lower* p99 than 100 frames/s.

Full per-platform tables and method are in the measurement record that accompanies this
feature.

## Limits

| | limit | what happens when you reach it |
|---|---|---|
| Buffered per window, awaiting collection | 8 MB or 256 frames, whichever first | `Send` blocks; `TrySend` returns `ErrStreamFull` |
| Buffered across the application, awaiting collection/write | 256 MB or 8,192 data frames | same |
| Received per connection, awaiting `Receive` | 8 MB or 256 frames | the frontend's `send()` is retried for you until the handler catches up |
| Received across the application, awaiting `Receive` | 256 MB or 8,192 frames | same |
| Connections per window | 256 | the open is retried for you until a slot frees |
| Live connections across the application | 4,096 | same |
| Sessions per window | 16 | a reload supersedes its own older session; otherwise the open is retried |
| A single frame, either direction | 64 MB | Go returns `ErrStreamTooLarge`; from JS the stream raises `error` and closes |
| A stream name | 256 UTF-8 bytes | the open is rejected and the stream raises `error` |
| Idle poll hold | 20 s | the poll returns empty and the runtime immediately reissues it |

Nothing above drops data silently. The two rows that say *retried for you* are ordinary
backpressure — the runtime holds the frame and retries with a short backoff, so your code
sees a slower stream rather than an error. The rows that raise `error` are programming
mistakes rather than load, and they surface rather than being papered over.

These are compile-time constants today, not options. See the internals guide if you need to
change them.

## When not to use a stream

- **For request/response, use bindings.** Streams are for continuous or unsolicited data;
  a call that returns a value is simpler as a bound method.
- **For application events, use `Emit`/`On`.** Events fan out to every listener and are a
  separate, well-worn system. Streams are point-to-point.
- **Not available for `InitialHTML` windows.** Those load with `origin === "null"`, so
  they cannot reach the asset server at all.
