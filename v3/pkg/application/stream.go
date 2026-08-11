package application

// GoStream gives Wails apps the WebSocket programming model without binding a
// TCP port. A WebSocket cannot be spoken over a custom URL scheme, so the only
// way to get one in a webview today is a real listener (see websocket_server.go,
// build tag "server"), which means an open local port reachable by any other
// process on the machine. Streams instead ride the asset server that already
// exists and is already origin-bound: Go→JS over a held poll, JS→Go over a
// normal POST.
//
// This is deliberately separate from the event system. It shares no code with
// Emit/events and does not change them. What it does take from the event work
// is that work's lessons, which are recorded at each rule below:
//
//   - Nothing in the Go→JS path touches the main thread. An event emitted from
//     the main thread used to run its eval inline while an earlier goroutine
//     emit was still queued, inverting 4.4% of events on all three platforms.
//     Send() appends under a mutex and returns; a single drainer preserves
//     order regardless of which goroutine produced the frame.
//   - Nothing in this path touches evaluateJavaScript, at any size. Splicing
//     payload into eval source retains host memory above a platform-specific
//     knee — 11.6 GB on macOS, 6.2 GB on WebKitGTK, at 100 events/sec of 1 MB.
//   - Every buffer is bounded, and every buffer is dropped when its window
//     goes away, following eventPayloadStore.
//
// A connection's lifetime is its handler goroutine's lifetime, exactly as with
// gorilla/coder WebSocket handlers. That single choice is what makes reload,
// shutdown and cleanup fall out rather than each needing its own policy.

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

var (
	// ErrStreamClosed is returned by Send and Receive once the connection is
	// gone: the page navigated or reloaded, the window closed, the frontend
	// called close(), or the app is shutting down.
	ErrStreamClosed = errors.New("wails: stream connection closed")

	// ErrStreamFull is returned by TrySend when the window's outbound buffer
	// has no room. Send blocks in the same situation.
	ErrStreamFull = errors.New("wails: stream send buffer full")
)

const (
	// streamOutQueueBytes bounds bytes buffered for one window awaiting its
	// next poll. This is the backstop that actually bounds host memory: the
	// depth limit below can be satisfied by 256 one-megabyte frames.
	streamOutQueueBytes = 8 << 20

	// streamOutQueueDepth bounds frames buffered for one window.
	//
	// Note this is deliberately NOT eventQueueCapacity (64). That constant was
	// measured for a queue drained one eval at a time, where depth bought
	// nothing but tail latency (capacity 16 → p99 179µs, capacity 1024 → 427µs
	// at identical throughput). A poll drains in batches, so depth here has to
	// cover one round trip of production instead: at 5000 frames/s and a 5 ms
	// round trip that is ~25 frames, and 256 leaves room for a burst without
	// stalling the producer.
	streamOutQueueDepth = 256

	// streamHoldTimeout is how long an empty poll is parked before returning
	// 204. It bounds idle cost and, until cancellation is plumbed through from
	// the platform layer, it also bounds how long a poll belonging to a page
	// that has navigated away can linger.
	streamHoldTimeout = 20 * time.Second

	// streamMaxResponseBytes caps one poll response. Necessary rather than
	// tidy: the Windows response writer accumulates the whole body in memory
	// and only hands it to WebView2 in Finish, so an unbounded response is an
	// unbounded allocation there. A truncated response sets the "more" flag and
	// the client re-polls without delay.
	streamMaxResponseBytes = 1 << 20

	// streamSessionTTL is how long a session survives without a poll before it
	// is treated as gone and its connections are closed. It must be comfortably
	// longer than streamHoldTimeout, or a session would be reaped while its own
	// poll is legitimately parked.
	streamSessionTTL = 3 * streamHoldTimeout

	streamSessionSweep = streamHoldTimeout

	// streamSessionGrace applies to a session that still has connections. It is
	// deliberately long: the cost of reaping a live session is losing an app's
	// streams mid-flight, while the cost of keeping a dead one is a little
	// memory until the window closes.
	streamSessionGrace = 10 * time.Minute

	// streamInQueueDepth and streamInQueueBytes bound frames received from the
	// frontend and not yet taken by Receive. Without them a handler slow to call
	// Receive lets the frontend grow host memory without limit, since the send
	// endpoint responds as soon as it has queued and the client is then free to
	// post again.
	//
	// Fullness is signalled, not waited out — the desktop endpoint answers 429
	// and the client retries the same frame. Blocking instead holds a request
	// open, and the held request is this transport's scarce resource: enough
	// stalled sends starve the window's single poll until its session TTL
	// expires. Measured at 25-32% of upload throughput plus outright failure of
	// multi-connection uploads.
	streamInQueueDepth = 256
	streamInQueueBytes = 8 << 20

	// A page normally has only a handful of Streams. The limit is deliberately
	// much larger than that, but finite: open acknowledgements and close
	// notifications must bypass data backpressure, and without a connection
	// admission bound a page that stopped polling could grow both protocol state
	// and handler goroutines without limit.
	streamMaxConnections = 256

	// Non-close control frames (open acknowledgements and refused opens) have a
	// separate allowance. Close notifications have their own reserved allowance
	// below, so a burst of bad opens cannot consume the space required to tell
	// already accepted peers that their connections ended.
	streamOutControlDepth = streamMaxConnections
)

// StreamHandler is invoked once per connection, on its own goroutine. The
// connection is live for as long as the handler runs; returning from it closes
// the connection, so the handler should block on Receive (or on c.Context())
// for as long as it wants the connection open.
type StreamHandler func(*StreamConn)

// streamSink is where a connection's outbound frames go. On the desktop it is
// the window's held-poll session; in server mode it is a real WebSocket. The
// handler code above the sink is identical either way, which is the whole point
// of mirroring the WebSocket model rather than inventing an API.
type streamSink interface {
	enqueue(c *StreamConn, connID uint32, kind uint8, data []byte, block bool) error
	removeConn(connID uint32)
	wake()
	// wakeProducers releases anything blocked waiting for queue space, so a
	// closing connection cannot leave a Send parked indefinitely.
	wakeProducers()
}

// StreamConn is one connection to a named stream. It is the moral equivalent of
// a *websocket.Conn: Send blocks like a socket write, Receive blocks like a
// socket read, and both fail once the peer is gone.
type StreamConn struct {
	id       uint32
	name     string
	windowID uint
	sink     streamSink

	ctx    context.Context
	cancel context.CancelFunc

	// Inbound frames from the frontend. Ordering is guaranteed by the client
	// serialising its sends and by handleSend appending before it responds, so
	// the next send cannot be issued until this one is queued.
	inMu    sync.Mutex
	inCond  *sync.Cond
	in      [][]byte
	inBytes int

	closeOnce sync.Once
}

// Name is the stream name this connection was opened against.
func (c *StreamConn) Name() string { return c.name }

// Context is cancelled when the connection closes, so anything the handler
// spawned can unwind with it.
func (c *StreamConn) Context() context.Context { return c.ctx }

// Window is the window that opened this connection. It is nil if the window has
// since been destroyed.
func (c *StreamConn) Window() Window {
	if globalApplication == nil {
		return nil
	}
	w, ok := globalApplication.Window.GetByID(c.windowID)
	if !ok {
		return nil
	}
	return w
}

// Send queues a frame for delivery to the frontend. It blocks while the
// window's outbound buffer is full, the way a socket write blocks on a full
// send buffer, and returns ErrStreamClosed once the connection is gone.
//
// The frame is queued, not copied. Do not mutate or reuse data after passing it
// to Send; hand over a slice the transport can own until it is delivered. This
// deviates from a real WebSocket, which copies on send — matching the API shape
// is the goal, and a memcpy per frame costs more than it is worth at the rates
// this sustains.
//
// It never touches the main thread, and never reaches evaluateJavaScript.
func (c *StreamConn) Send(data []byte) error {
	return c.sink.enqueue(c, c.id, frameData, data, true)
}

// TrySend is Send without the blocking: it returns ErrStreamFull rather than
// waiting for the frontend to catch up.
func (c *StreamConn) TrySend(data []byte) error {
	return c.sink.enqueue(c, c.id, frameData, data, false)
}

// SendJSON marshals v and sends it as a single frame. A convenience over Send
// for the common case; the wire carries the marshalled bytes and nothing else,
// so the frontend can read it with JSON.parse or with JSONStream.
func (c *StreamConn) SendJSON(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.Send(data)
}

// ReceiveJSON blocks for the next frame and unmarshals it into v. A malformed
// frame returns the unmarshal error and consumes the frame — it does not close
// the connection, so a handler can log and carry on.
func (c *StreamConn) ReceiveJSON(v any) error {
	frame, err := c.Receive()
	if err != nil {
		return err
	}
	return json.Unmarshal(frame, v)
}

// Receive returns the next frame from the frontend, blocking until one arrives
// or the connection closes.
func (c *StreamConn) Receive() ([]byte, error) {
	c.inMu.Lock()
	defer c.inMu.Unlock()
	for {
		if len(c.in) > 0 {
			frame := c.in[0]
			// Clear the slot so the backing array does not pin a delivered
			// frame for as long as the slice header lives.
			c.in[0] = nil
			c.in = c.in[1:]
			c.inBytes -= len(frame)
			// Wake a delivering goroutine waiting for room.
			c.inCond.Broadcast()
			return frame, nil
		}
		if c.ctx.Err() != nil {
			return nil, ErrStreamClosed
		}
		c.inCond.Wait()
	}
}

// Close ends the connection and tells the frontend, which sees an onclose. It
// is safe to call more than once and from any goroutine.
func (c *StreamConn) Close() error {
	c.closeOnce.Do(func() {
		// Mark the connection closed before queueing the notification. Otherwise
		// a concurrent Send can slip in behind the close frame, report success,
		// and then be discarded by a frontend that has already processed close.
		c.cancel()
		// Best-effort notification: if the buffer is full or the session is
		// already gone the frontend finds out via the session going away
		// instead, so this must not block.
		_ = c.sink.enqueue(c, c.id, frameClose, nil, false)
		c.shutdown()
	})
	return nil
}

// closedByPeer marks the connection as already closed from the other end, so a
// later Close() does not queue a close frame to a peer that has gone. It
// consumes the same sync.Once that Close() uses.
func (c *StreamConn) closedByPeer() {
	c.closeOnce.Do(func() {})
	c.shutdown()
}

// shutdown tears down the local end without notifying the frontend. Used when
// the frontend is the one that went away.
func (c *StreamConn) shutdown() {
	c.cancel()
	c.sink.removeConn(c.id)

	// Wake anyone blocked in Receive so it can observe the cancelled context,
	// and anyone blocked in Send waiting for queue space. The comment always
	// claimed the latter; only wake() was called, which nudges the poll and not
	// the producers, so a Send blocked on a full queue stayed blocked forever
	// after its connection closed.
	c.inMu.Lock()
	c.inCond.Broadcast()
	c.inMu.Unlock()
	c.sink.wakeProducers()
	c.sink.wake()
}

// deliver queues an inbound frame and wakes a waiting Receive. It returns
// ErrStreamFull when the inbox has no room, and ErrStreamClosed once the
// connection is gone.
//
// It does not block, deliberately. Blocking here bounds memory but holds the
// request open while it waits, and on this transport the held request is the
// scarce resource: enough stalled sends starve the window's single poll, which
// then misses its session TTL and the whole session is reaped. Measured, that
// cost 25-32% of upload throughput and broke multi-connection uploads outright.
// The caller signals fullness to the client instead, which retries — bounded
// memory without occupying a slot to achieve it.
//
// An empty inbox always accepts one frame however large, for the same reason
// the outbound queue does: the caller does not choose the frame size, so a
// frame bigger than the cap must not become impossible to deliver.
func (c *StreamConn) deliver(data []byte) error {
	c.inMu.Lock()
	defer c.inMu.Unlock()

	if c.ctx.Err() != nil {
		return ErrStreamClosed
	}
	if len(c.in) >= streamInQueueDepth ||
		(len(c.in) > 0 && c.inBytes+len(data) > streamInQueueBytes) {
		return ErrStreamFull
	}

	c.in = append(c.in, data)
	c.inBytes += len(data)
	c.inCond.Broadcast()
	return nil
}

// ---------------------------------------------------------------------------
// Registry
// ---------------------------------------------------------------------------

type streamManager struct {
	app *App

	mu       sync.Mutex
	handlers map[string]StreamHandler
	sessions map[string]*streamSession
	// retiredThrough is the highest page generation that has finished in each
	// window. Keeping a watermark rather than every retired session id prevents
	// reloads from growing manager memory for the lifetime of the application.
	retiredThrough map[uint]uint64
	closed         bool

	janitor sync.Once
	once    sync.Once
	stop    chan struct{}
}

func newStreamManager(app *App) *streamManager {
	return &streamManager{
		app:            app,
		handlers:       make(map[string]StreamHandler),
		sessions:       make(map[string]*streamSession),
		retiredThrough: make(map[uint]uint64),
		stop:           make(chan struct{}),
	}
}

// HandleStream registers a handler for a named stream. The frontend connects to
// it by that name, and the handler runs once per connection on its own
// goroutine:
//
//	app.HandleStream("telemetry", func(c *application.StreamConn) {
//		defer c.Close()
//		for {
//			frame, err := c.Receive()
//			if err != nil {
//				return // reload, close, or shutdown
//			}
//			c.Send(reply(frame))
//		}
//	})
//
// Registering the same name twice replaces the handler; connections already
// open keep the handler they started with.
func (a *App) HandleStream(name string, handler StreamHandler) {
	if a.streams == nil || name == "" || handler == nil {
		return
	}
	a.streams.mu.Lock()
	defer a.streams.mu.Unlock()
	a.streams.handlers[name] = handler
}

func (m *streamManager) handler(name string) (StreamHandler, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	h, ok := m.handlers[name]
	return h, ok
}

// session returns the session for id, creating it if this is the first request
// from that page. Sessions are created by whichever of poll or send arrives
// first; the frontend does not perform a separate handshake.
func (m *streamManager) session(id string, windowID uint, generation uint64, maySupersede bool) *streamSession {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil
	}
	s, ok := m.sessions[id]
	// Resolve an existing id only within the window that owns it. sessionFor
	// reports the eventual 403, but the manager must avoid mutating anything
	// before that check: a foreign poll used to borrow this session's newer
	// generation and retire legitimate sessions in the requesting window.
	if ok && windowID != 0 && s.windowID != windowID {
		m.mu.Unlock()
		return s
	}
	// The same rule applies to a generation mismatch. sessionFor reports 403,
	// but an invalid poll must not use the generation stored on the colliding
	// session id to retire other valid sessions before that rejection.
	if ok && s.generation != generation {
		m.mu.Unlock()
		return s
	}
	if !ok {
		// Once a page generation has finished, no delayed request from that page
		// may recreate it. Without this watermark, a late open POST could
		// allocate the old session again after its replacement retired it.
		if generation <= m.retiredThrough[windowID] {
			m.mu.Unlock()
			return nil
		}
		s = newStreamSession(id, windowID, m)
		s.generation = generation
		m.sessions[id] = s
		m.janitor.Do(func() { go m.reap() })
	}

	var superseded []*streamSession
	// The open POST and first poll start concurrently. Whichever request creates
	// the new session gives it a generation; when its poll arrives it retires
	// only older sessions. A late poll from the old page sees the newer
	// generation and leaves it alone.
	if maySupersede {
		// A poll proves that this is the current page generation for the window.
		// Retire the entire lower generation range, not just sessions that happen
		// to exist already: an older page may have dispatched its first request
		// before navigation but reach Go only after this poll completes.
		if s.generation > 1 && s.generation-1 > m.retiredThrough[windowID] {
			m.retiredThrough[windowID] = s.generation - 1
		}

		// A new session id in a window that already has one means that window
		// navigated or reloaded: the previous page is gone. Without this, the
		// old session survives until the TTL sweep — up to
		// streamSessionTTL + streamSessionSweep, 60-80s at the current values —
		// during which the app holds two live connections to the same stream
		// and a handler that owns a per-connection resource has two of them.
		// A WebSocket closes on unload; this is the closest portable equivalent
		// until the platform layer reports a cancelled request.
		//
		// Keyed on the window, not the client, so several runtimes under one
		// window id (an iframe loading /wails/runtime.js gets the same id and
		// its own session) would collide. Superseding is still correct there:
		// the reload that replaces the top-level document replaces its frames
		// too. What it does not survive is a frame reloading independently of
		// its parent, which would close the parent's session; if that ever
		// needs supporting it wants a per-document key, not a per-window one.
		for otherID, other := range m.sessions {
			if otherID != id && other.windowID == windowID && other.generation < s.generation {
				superseded = append(superseded, other)
				delete(m.sessions, otherID)
				m.retireLocked(other)
			}
		}
	}
	m.mu.Unlock()

	// Outside the lock: close() takes each session's own lock and shuts down
	// its connections.
	for _, old := range superseded {
		old.close()
	}

	s.touch()
	return s
}

// existingSession is session lookup without creation, for requests that must
// not be able to conjure state.
func (m *streamManager) existingSession(id string) *streamSession {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sessions[id]
}

// removeSession forgets closed transport state. It must not retire the page
// generation: the document that observed the close may still be loaded and
// create another stream later.
func (m *streamManager) removeSession(id string) {
	m.mu.Lock()
	delete(m.sessions, id)
	m.mu.Unlock()
}

func (m *streamManager) retireLocked(s *streamSession) {
	if s.generation > m.retiredThrough[s.windowID] {
		m.retiredThrough[s.windowID] = s.generation
	}
}

// dropWindow closes every session belonging to a window that is going away.
// Mirrors eventPayloadStore.dropWindow: the window is gone, so nothing will
// ever poll again and waiting for the TTL would just hold memory.
func (m *streamManager) dropWindow(windowID uint) {
	m.mu.Lock()
	var doomed []*streamSession
	for id, s := range m.sessions {
		if s.windowID == windowID {
			doomed = append(doomed, s)
			delete(m.sessions, id)
			m.retireLocked(s)
		}
	}
	m.mu.Unlock()

	for _, s := range doomed {
		s.close()
	}
}

// reap closes sessions whose page has stopped polling. Until the platform layer
// reports a cancelled request, this is how a reload or a crashed renderer is
// noticed: a page that has gone away stops polling, and after streamSessionTTL
// its connections are closed and their handlers unblock.
func (m *streamManager) reap() {
	ticker := time.NewTicker(streamSessionSweep)
	defer ticker.Stop()
	for {
		select {
		case <-m.stop:
			return
		case now := <-ticker.C:
			m.reapStale(now)
		}
	}
}

// reapStale performs one janitor sweep. Expiring a session is deliberately not
// the same as retiring its page generation: a still-loaded page stops polling
// when its last connection closes and must be able to create another stream
// after the idle session has expired. Only a newer page poll or destruction of
// the owning window proves that the page generation itself can never return.
func (m *streamManager) reapStale(now time.Time) {
	m.mu.Lock()
	var doomed []*streamSession
	for id, s := range m.sessions {
		// A session with live connections is owned by running handlers,
		// and under saturation the page can be busy dispatching a large
		// response for longer than the TTL before it issues the next
		// poll. Reaping on that basis killed streams mid-run at full
		// load — measured, after 230,077 frames. Idle time alone is only
		// trustworthy once nothing is connected; a page that really has
		// gone is caught by window destroy or by the next page's poll
		// superseding it, and the long grace below is the backstop for a
		// renderer that died without either.
		ttl := streamSessionTTL
		if s.connCount() > 0 {
			ttl = streamSessionGrace
		}
		if now.Sub(s.lastSeenAt()) > ttl {
			doomed = append(doomed, s)
			delete(m.sessions, id)
		}
	}
	m.mu.Unlock()

	for _, s := range doomed {
		s.close()
	}
}

func (m *streamManager) close() {
	m.mu.Lock()
	m.closed = true
	doomed := make([]*streamSession, 0, len(m.sessions))
	for _, s := range m.sessions {
		doomed = append(doomed, s)
	}
	m.sessions = map[string]*streamSession{}
	m.retiredThrough = map[uint]uint64{}
	m.mu.Unlock()

	for _, s := range doomed {
		s.close()
	}

	m.janitor.Do(func() {}) // consume the Once so reap can never start later
	m.once.Do(func() { close(m.stop) })
}
