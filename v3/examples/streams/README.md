# Streams Example

The smallest useful demonstration of a [stream](https://v3.wails.io/guides/streams/): a
named, ordered, bidirectional byte channel between Go and the frontend, with the same
programming model as a WebSocket and **no listening socket**.

Go sends the time once a second, and echoes anything the frontend sends back — as objects,
using the JSON convenience on both sides.

## Running the example

```bash
go run .
```

Type into the box and press Enter. Messages from Go appear in the log, and what you send
comes back prefixed with `echo:`. The Go side logs what it receives.

## What to look at

**`main.go`** — one handler is the whole registration:

```go
app.HandleStream("hello", func(c *application.StreamConn) {
    defer c.Close()
    for {
        var msg map[string]any
        if err := c.ReceiveJSON(&msg); err != nil {
            return                  // page reloaded, window closed, or app shutting down
        }
        msg["echoed"] = true
        c.SendJSON(msg)
    }
})
```

The handler runs once per connection on its own goroutine, and its lifetime *is* the
connection's lifetime — returning from it closes the connection. A separate goroutine does
the ticking, and exits when `c.Context()` is cancelled.

**`assets/index.html`** — the frontend object implements the `WebSocket` interface:

```js
const stream = JSONStream("hello");    // synchronous, like new WebSocket(url)
stream.onmessage = (ev) => log(JSON.stringify(ev.data));   // already an object
stream.send({ text: "anything" });                         // stringified for you
```

## If you want bytes instead

`JSONStream` is `Stream` with the encoding done at the boundary — the wire carries bytes
either way, and a Go handler cannot tell which the frontend used. Swap in `Stream` and
`Send`/`Receive` when you want to control the encoding yourself: protobuf, CBOR, or any
binary format.

One thing to know if you do: **`ev.data` is then an `ArrayBuffer`, never a string.**
`JSON.parse(ev.data)` receives the string `"[object ArrayBuffer]"` and fails. Decode it
first:

```js
const decoder = new TextDecoder();
stream.onmessage = (ev) => JSON.parse(decoder.decode(ev.data));
```

Sending a string needs no change; it is encoded as UTF-8 for you.

## Try this

- **Reload the page.** Go logs a disconnect and then a fresh connect — a reload closes the
  connection the way it closes a socket, and the new page gets a new one.
- **Close the window.** The handler returns.

## See also

- [Streams](https://v3.wails.io/guides/streams/) — the API
- [Migrating a WebSocket to Streams](https://v3.wails.io/guides/streams-from-websockets/) —
  if you have an app that stands up a local HTTP server to push data to the frontend
