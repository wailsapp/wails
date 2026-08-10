# Streams Example

The smallest useful demonstration of a [stream](https://v3.wails.io/guides/streams/): a
named, ordered, bidirectional byte channel between Go and the frontend, with the same
programming model as a WebSocket and **no listening socket**.

Go sends the time once a second, and echoes anything the frontend sends back.

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
        frame, err := c.Receive()   // blocks until a frame arrives
        if err != nil {
            return                  // page reloaded, window closed, or app shutting down
        }
        c.Send([]byte("echo: " + string(frame)))
    }
})
```

The handler runs once per connection on its own goroutine, and its lifetime *is* the
connection's lifetime — returning from it closes the connection. A separate goroutine does
the ticking, and exits when `c.Context()` is cancelled.

**`assets/index.html`** — the frontend object implements the `WebSocket` interface:

```js
const stream = Stream("hello");        // synchronous, like new WebSocket(url)
stream.onmessage = (ev) => log(decoder.decode(ev.data));
stream.send("anything");
```

Two things worth noticing, because they are the ones people trip over:

- **`ev.data` is always an `ArrayBuffer`, never a string.** Decode it before use. Calling
  `JSON.parse(ev.data)` will not throw — it will parse `"[object ArrayBuffer]"`.
- **Sending a string works unchanged**; it is encoded as UTF-8. Only the receive path
  needs the decoder.

## Try this

- **Reload the page.** Go logs a disconnect and then a fresh connect — a reload closes the
  connection the way it closes a socket, and the new page gets a new one.
- **Close the window.** The handler returns.

## See also

- [Streams](https://v3.wails.io/guides/streams/) — the API
- [Migrating a WebSocket to Streams](https://v3.wails.io/guides/streams-from-websockets/) —
  if you have an app that stands up a local HTTP server to push data to the frontend
