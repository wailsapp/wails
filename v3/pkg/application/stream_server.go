//go:build server

package application

import (
	"context"
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
	conn     *websocket.Conn
	mgr      *streamManager
	ioCtx    context.Context
	cancelIO context.CancelFunc
	done     chan struct{}
	shutOnce sync.Once

	mu     sync.Mutex
	cond   *sync.Cond
	out    []outFrame
	bytes  int
	closed bool
}

func newWSStreamSink(conn *websocket.Conn, parent context.Context, mgr *streamManager) *wsStreamSink {
	// Accepted writes must outlive the handler context long enough to drain.
	// StreamConn.Close cancels that context immediately after enqueueing the
	// close frame, so using it for writes loses any data queued just before a
	// handler returns. It is still rooted in the application context so app
	// shutdown closes server streams and releases handlers blocked in Receive.
	if parent == nil {
		parent = context.Background()
	}
	ioCtx, cancelIO := context.WithCancel(parent)
	w := &wsStreamSink{
		conn:     conn,
		mgr:      mgr,
		ioCtx:    ioCtx,
		cancelIO: cancelIO,
		done:     make(chan struct{}),
	}
	w.cond = sync.NewCond(&w.mu)
	go w.pump()
	return w
}

func (w *wsStreamSink) enqueue(c *StreamConn, connID uint32, kind uint8, data []byte, block bool) error {
	// Control frames bypass the cap, as on the desktop: losing one is a
	// protocol failure rather than a slow-down.
	control := kind != frameData

	for {
		w.mu.Lock()
		if w.closed {
			w.mu.Unlock()
			return ErrStreamClosed
		}
		if c != nil && c.ctx.Err() != nil {
			w.mu.Unlock()
			return ErrStreamClosed
		}
		localRoom := control || (len(w.out) < streamOutQueueDepth &&
			(len(w.out) == 0 || w.bytes+len(data) <= streamOutQueueBytes))
		if localRoom {
			w.mu.Unlock()
		} else {
			if !block {
				w.mu.Unlock()
				return ErrStreamFull
			}
			w.cond.Wait()
			w.mu.Unlock()
			continue
		}

		if !control && w.mgr != nil {
			var ctx context.Context
			if c != nil {
				ctx = c.ctx
			}
			if err := w.mgr.reserveOutbound(frameData, len(data), block, ctx); err != nil {
				return err
			}
		}

		w.mu.Lock()
		localRoom = control || (len(w.out) < streamOutQueueDepth &&
			(len(w.out) == 0 || w.bytes+len(data) <= streamOutQueueBytes))
		if w.closed || (c != nil && c.ctx.Err() != nil) {
			w.mu.Unlock()
			if !control && w.mgr != nil {
				w.mgr.releaseOutbound(frameData, len(data))
			}
			return ErrStreamClosed
		}
		if !localRoom {
			w.mu.Unlock()
			if !control && w.mgr != nil {
				w.mgr.releaseOutbound(frameData, len(data))
			}
			if !block {
				return ErrStreamFull
			}
			continue
		}

		// No copy, as on the desktop: the caller owns the slice until the frame is
		// written. See StreamConn.Send.
		w.out = append(w.out, outFrame{connID: connID, kind: kind, data: data})
		w.bytes += len(data)
		w.cond.Broadcast()
		w.mu.Unlock()
		return nil
	}
}

// pump is the single writer. websocket.Conn does not permit concurrent writes,
// and serialising here also preserves the order frames were accepted in.
func (w *wsStreamSink) pump() {
	defer close(w.done)
	for {
		w.mu.Lock()
		for len(w.out) == 0 && !w.closed {
			w.cond.Wait()
		}
		if w.closed && len(w.out) == 0 {
			w.mu.Unlock()
			w.shut(websocket.StatusNormalClosure, "")
			return
		}
		frame := w.out[0]
		w.out[0] = outFrame{}
		w.out = w.out[1:]
		w.bytes -= len(frame.data)
		w.cond.Broadcast()
		w.mu.Unlock()

		var writeErr error
		switch frame.kind {
		case frameData:
			writeErr = w.conn.Write(w.ioCtx, websocket.MessageBinary, frame.data)
		case frameClose:
			w.shut(websocket.StatusNormalClosure, "")
			return
		case frameError:
			w.shut(websocket.StatusInternalError, string(frame.data))
			return
		case frameOpen:
			// The handshake is the acknowledgement; there is nothing to send.
		}
		if frame.kind == frameData && w.mgr != nil {
			w.mgr.releaseOutbound(frameData, len(frame.data))
		}
		if writeErr != nil {
			w.shut(websocket.StatusNormalClosure, "")
			return
		}
	}
}

func (w *wsStreamSink) releaseQueued() {
	w.mu.Lock()
	queued := w.out
	w.out = nil
	w.bytes = 0
	w.cond.Broadcast()
	w.mu.Unlock()
	if w.mgr != nil {
		for i := range queued {
			if queued[i].kind == frameData {
				w.mgr.releaseOutbound(frameData, len(queued[i].data))
			}
		}
	}
}

func (w *wsStreamSink) shut(code websocket.StatusCode, reason string) {
	w.shutOnce.Do(func() {
		w.mu.Lock()
		w.closed = true
		w.cond.Broadcast()
		w.mu.Unlock()
		_ = w.conn.Close(code, reason)
		w.cancelIO()
		w.releaseQueued()
	})
}

func (w *wsStreamSink) abort() {
	w.cancelIO()
	_ = w.conn.CloseNow()
	w.mu.Lock()
	w.closed = true
	w.cond.Broadcast()
	w.mu.Unlock()
	w.releaseQueued()
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
	if name == "" || len(name) > streamMaxNameLen {
		http.Error(rw, "invalid stream name", http.StatusBadRequest)
		return
	}
	handler, ok := a.streams.handler(name)
	if !ok {
		http.Error(rw, "no handler registered for stream "+name, http.StatusNotFound)
		return
	}
	if !a.streams.reserveLifecycle() {
		http.Error(rw, "stream connection capacity reached", http.StatusServiceUnavailable)
		return
	}
	defer a.streams.releaseLifecycle()

	conn, err := websocket.Accept(rw, req, a.websocketAcceptOptions())
	if err != nil {
		a.error("stream websocket accept failed: %w", err)
		return
	}

	// coder/websocket defaults to a 32 KiB read limit, which would close the
	// connection on any frontend frame past that — well below the 64 MB this
	// transport documents.
	conn.SetReadLimit(streamMaxFrameBytes)

	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	sink := newWSStreamSink(conn, a.ctx, a.streams)
	c := &StreamConn{
		name:    name,
		sink:    sink,
		ctx:     ctx,
		manager: a.streams,
		// windowID stays 0: server mode has no windows, and StreamConn.Window
		// already reports nil for it.
		cancel: cancel,
	}
	c.inCond = sync.NewCond(&c.inMu)

	// Read pump. A read error is the socket closing, which is what ends the
	// handler — the same signal a reload gives the desktop build.
	go func() {
		defer c.closedByPeer()
		for {
			// Socket I/O has a separate lifetime from the StreamConn context.
			// A local Close cancels the latter to unblock the handler, but the
			// reader must stay alive while the writer flushes accepted data and
			// completes the WebSocket close handshake.
			_, data, err := conn.Read(sink.ioCtx)
			if err != nil {
				return
			}
			// A socket has no retry channel, so the read pump waits on the
			// connection's inbox condition and lets TCP apply backpressure to the
			// peer. Receive and shutdown wake it without polling.
			if err := c.deliverBlocking(data); err != nil {
				return
			}
		}
	}()

	defer func() { _ = conn.CloseNow() }()

	// Match the desktop lifecycle: returning from a handler closes the stream.
	// Keep panic recovery and Close in an inner scope so the accepted write
	// queue is still drained below after either a normal return or a panic.
	func() {
		defer handlePanic()
		defer c.Close()
		handler(c)
	}()

	// Close queues a control frame behind every accepted data frame. Do not
	// tear down the socket until the writer has drained that queue. A bounded
	// escape hatch is still required for a peer that has stopped reading.
	select {
	case <-sink.done:
	case <-time.After(5 * time.Second):
		sink.abort()
		<-sink.done
	}
}
