//go:build server

package application

import (
	"context"
	"net/http"
	"sync"

	"github.com/coder/websocket"
)

// Server mode already has a listener, so there is nothing to emulate: streams
// there are real WebSockets. Only the sink differs — the handler, StreamConn,
// and the application code above them are identical to the desktop build, which
// is the payoff for mirroring the WebSocket model instead of inventing an API.
//
// The desktop poll endpoints stay reachable in server mode (the asset server is
// mounted at "/"), so a client that has not installed the WebSocket factory
// still works, just over the poll.

// wsStreamSink writes a connection's frames straight to a socket.
//
// There is no outbound buffer and the block flag is ignored, deliberately:
// websocket.Conn.Write blocks until the frame is written, so the socket's own
// send buffer already provides exactly the backpressure that the desktop
// build's bounded queue exists to imitate.
type wsStreamSink struct {
	conn *websocket.Conn
	ctx  context.Context

	// websocket.Conn does not permit concurrent writes.
	mu     sync.Mutex
	closed bool
}

func (w *wsStreamSink) enqueue(c *StreamConn, connID uint32, kind uint8, data []byte, block bool) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return ErrStreamClosed
	}

	switch kind {
	case frameData:
		if err := w.conn.Write(w.ctx, websocket.MessageBinary, data); err != nil {
			return ErrStreamClosed
		}
	case frameClose:
		w.closed = true
		_ = w.conn.Close(websocket.StatusNormalClosure, "")
	case frameError:
		w.closed = true
		_ = w.conn.Close(websocket.StatusInternalError, string(data))
	case frameOpen:
		// The handshake is the acknowledgement; there is nothing to send.
	}
	return nil
}

// The connection table and the poll wake-up are session concepts. A socket is
// its own connection, so both are no-ops here.
func (w *wsStreamSink) removeConn(uint32) {}
func (w *wsStreamSink) wake()             {}

// serveStreamWS upgrades a request and runs the registered handler for the
// stream named in the query string. The handler runs on this goroutine, so the
// request lives exactly as long as the connection — the same relationship the
// desktop build gets from the handler goroutine.
func (a *App) serveStreamWS(rw http.ResponseWriter, req *http.Request) {
	if a.streams == nil {
		http.Error(rw, "streams unavailable", http.StatusServiceUnavailable)
		return
	}

	name := req.URL.Query().Get("name")
	handler, ok := a.streams.handler(name)
	if !ok {
		http.Error(rw, "no handler registered for stream "+name, http.StatusNotFound)
		return
	}

	conn, err := websocket.Accept(rw, req, &websocket.AcceptOptions{
		// Matches the events broadcaster: server mode serves whatever origin
		// the operator points at it.
		InsecureSkipVerify: true,
	})
	if err != nil {
		a.error("stream websocket accept failed", "error", err)
		return
	}

	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	sink := &wsStreamSink{conn: conn, ctx: ctx}
	c := &StreamConn{
		name: name,
		sink: sink,
		ctx:  ctx,
		// windowID stays 0: server mode has no windows, and StreamConn.Window
		// already reports nil for it.
		cancel: cancel,
	}
	c.inCond = sync.NewCond(&c.inMu)

	// Read pump. A read error is the socket closing, which is what ends the
	// handler — the same signal a reload gives the desktop build.
	go func() {
		defer c.shutdown()
		for {
			_, data, err := conn.Read(ctx)
			if err != nil {
				return
			}
			c.deliver(data)
		}
	}()

	defer func() {
		_ = conn.CloseNow()
	}()
	defer handlePanic()
	handler(c)
}
