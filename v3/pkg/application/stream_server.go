//go:build server

package application

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

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

// wsStreamSink writes a connection's frames to a socket through a bounded
// queue drained by a single writer goroutine.
//
// Writing inline would be simpler, but websocket.Conn.Write blocks until the
// peer makes progress, which would make TrySend block in server mode while it
// does not on the desktop — and the fan-out and broker-callback patterns in the
// migration guide pick TrySend precisely to avoid stalling their producer. The
// queue keeps one contract across both transports.
type wsStreamSink struct {
	conn *websocket.Conn
	ctx  context.Context

	mu     sync.Mutex
	cond   *sync.Cond
	out    []outFrame
	bytes  int
	closed bool
}

func newWSStreamSink(ctx context.Context, conn *websocket.Conn) *wsStreamSink {
	w := &wsStreamSink{conn: conn, ctx: ctx}
	w.cond = sync.NewCond(&w.mu)
	go w.pump()
	return w
}

func (w *wsStreamSink) enqueue(c *StreamConn, connID uint32, kind uint8, data []byte, block bool) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Control frames bypass the cap, as on the desktop: losing one is a
	// protocol failure rather than a slow-down.
	control := kind != frameData

	for !control {
		if w.closed {
			return ErrStreamClosed
		}
		if len(w.out) < streamOutQueueDepth &&
			(len(w.out) == 0 || w.bytes+len(data) <= streamOutQueueBytes) {
			break
		}
		if !block {
			return ErrStreamFull
		}
		w.cond.Wait()
	}
	if w.closed {
		return ErrStreamClosed
	}

	// No copy, as on the desktop: the caller owns the slice until the frame is
	// written. See StreamConn.Send.
	w.out = append(w.out, outFrame{connID: connID, kind: kind, data: data})
	w.bytes += len(data)
	w.cond.Broadcast()
	return nil
}

// pump is the single writer. websocket.Conn does not permit concurrent writes,
// and serialising here also preserves the order frames were accepted in.
func (w *wsStreamSink) pump() {
	for {
		w.mu.Lock()
		for len(w.out) == 0 && !w.closed {
			w.cond.Wait()
		}
		if w.closed && len(w.out) == 0 {
			w.mu.Unlock()
			return
		}
		frame := w.out[0]
		w.out[0] = outFrame{}
		w.out = w.out[1:]
		w.bytes -= len(frame.data)
		w.cond.Broadcast()
		w.mu.Unlock()

		switch frame.kind {
		case frameData:
			if err := w.conn.Write(w.ctx, websocket.MessageBinary, frame.data); err != nil {
				w.shut(websocket.StatusNormalClosure, "")
				return
			}
		case frameClose:
			w.shut(websocket.StatusNormalClosure, "")
			return
		case frameError:
			w.shut(websocket.StatusInternalError, string(frame.data))
			return
		case frameOpen:
			// The handshake is the acknowledgement; there is nothing to send.
		}
	}
}

func (w *wsStreamSink) shut(code websocket.StatusCode, reason string) {
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		return
	}
	w.closed = true
	w.cond.Broadcast()
	w.mu.Unlock()
	_ = w.conn.Close(code, reason)
}

// The connection table is a session concept; a socket is its own connection.
// removeConn releases the writer so a closing connection cannot leave it parked.
func (w *wsStreamSink) removeConn(uint32) {
	w.mu.Lock()
	w.closed = true
	w.cond.Broadcast()
	w.mu.Unlock()
}

func (w *wsStreamSink) wake() {}

func (w *wsStreamSink) wakeProducers() {
	w.mu.Lock()
	w.cond.Broadcast()
	w.mu.Unlock()
}

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
		a.error("stream websocket accept failed: %w", err)
		return
	}

	// coder/websocket defaults to a 32 KiB read limit, which would close the
	// connection on any frontend frame past that — well below the 64 MB this
	// transport documents.
	conn.SetReadLimit(streamMaxSendBytes)

	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	sink := newWSStreamSink(ctx, conn)
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
			// A socket has no retry channel, so here the read pump waits for
			// room and lets TCP apply the backpressure to the peer. Unlike the
			// desktop transport there is no request slot being held.
			for {
				err := c.deliver(data)
				if err == nil {
					break
				}
				if !errors.Is(err, ErrStreamFull) {
					return
				}
				select {
				case <-ctx.Done():
					return
				case <-time.After(time.Millisecond):
				}
			}
		}
	}()

	defer func() {
		_ = conn.CloseNow()
	}()
	defer handlePanic()
	handler(c)
}
