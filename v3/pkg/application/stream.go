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
	"sync/atomic"
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

	// ErrStreamTooLarge is returned when one frame exceeds the transport's
	// maximum wire size. The same limit applies in both directions and in both
	// desktop and server transports.
	ErrStreamTooLarge = errors.New("wails: stream frame too large")
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

	// Open POSTs may arrive before a page's first poll. Bound the number of
	// independent session ids they can create so varying that header cannot
	// multiply every per-session allowance without limit. The per-window limit
	// is intentionally well above the expected top-level page plus embedded
	// frames; the global cap also covers transports where the platform cannot
	// tag a request with a window id.
	streamMaxSessionsPerWindow = 16
	streamMaxSessions          = 1024

	// Per-session and per-connection queue limits still multiply unless the
	// manager applies an application-wide ceiling. These allowances support a
	// substantial burst across independent windows while preventing admitted
	// sessions from turning nominal 8 MiB bounds into multi-gigabyte retention.
	streamOutQueueBytesGlobal = 256 << 20
	streamOutQueueDepthGlobal = 8192
	streamInQueueBytesGlobal  = 256 << 20
	streamInQueueDepthGlobal  = 8192

	// A live desktop connection holds one lifecycle slot from open until its
	// StreamConn is torn down.
	streamMaxConnectionsGlobal = 4096

	// Queued close notifications are counted separately from other controls, so
	// refused opens cannot consume the capacity an accepted connection needs to
	// report that it ended. Each connection queues at most one close, and a
	// close carries no payload, so matching the connection allowance bounds this
	// at a few kilobytes of bookkeeping.
	streamOutCloseDepthGlobal   = streamMaxConnectionsGlobal
	streamOutControlDepthGlobal = streamMaxConnectionsGlobal
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

	closeOnce    sync.Once
	shutdownOnce sync.Once
	manager      *streamManager
	lifecycle    bool
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
	if len(data) > streamMaxFrameBytes {
		return ErrStreamTooLarge
	}
	return c.sink.enqueue(c, c.id, frameData, data, true)
}

// TrySend is Send without the blocking: it returns ErrStreamFull rather than
// waiting for the frontend to catch up.
func (c *StreamConn) TrySend(data []byte) error {
	if len(data) > streamMaxFrameBytes {
		return ErrStreamTooLarge
	}
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
	for {
		if len(c.in) > 0 {
			frame := c.in[0]
			// Clear the slot so the backing array does not pin a delivered
			// frame for as long as the slice header lives.
			c.in[0] = nil
			c.in = c.in[1:]
			c.inBytes -= len(frame)
			c.inCond.Broadcast()
			c.inMu.Unlock()
			if c.manager != nil {
				c.manager.releaseInbound(len(frame), 1)
			}
			return frame, nil
		}
		if c.ctx.Err() != nil {
			c.inMu.Unlock()
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
	c.shutdownOnce.Do(func() {
		c.cancel()
		c.sink.removeConn(c.id)

		c.inMu.Lock()
		releasedBytes := c.inBytes
		releasedFrames := len(c.in)
		for i := range c.in {
			c.in[i] = nil
		}
		c.in = nil
		c.inBytes = 0
		c.inCond.Broadcast()
		c.inMu.Unlock()
		if c.manager != nil {
			c.manager.releaseInbound(releasedBytes, releasedFrames)
			// The lifecycle slot belongs to this connection for exactly as long
			// as the connection exists, and shutdown runs exactly once. A queued
			// close frame owns a separate close slot, released when that frame is
			// drained or discarded, so nothing here has to reason about whether
			// the notification made it out.
			if c.lifecycle {
				c.manager.releaseLifecycle()
			}
			c.manager.signalBudget()
		}
		c.sink.wakeProducers()
		c.sink.wake()
	})
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
	return c.deliverWithBackpressure(data, false)
}

// deliverBlocking is the socket-transport form of deliver. A real WebSocket
// has no HTTP response on which to signal retry, so its single read pump waits
// here and lets TCP backpressure reach the peer. Receive and shutdown broadcast
// on inCond, making this event-driven rather than a millisecond retry loop.
func (c *StreamConn) deliverBlocking(data []byte) error {
	return c.deliverWithBackpressure(data, true)
}

func (c *StreamConn) deliverWithBackpressure(data []byte, block bool) error {
	for {
		c.inMu.Lock()
		if c.ctx.Err() != nil {
			c.inMu.Unlock()
			return ErrStreamClosed
		}
		if len(c.in) < streamInQueueDepth &&
			(len(c.in) == 0 || c.inBytes+len(data) <= streamInQueueBytes) {
			c.inMu.Unlock()
		} else {
			if !block {
				c.inMu.Unlock()
				return ErrStreamFull
			}
			c.inCond.Wait()
			c.inMu.Unlock()
			continue
		}

		if c.manager != nil {
			if err := c.manager.reserveInbound(len(data), block, c.ctx); err != nil {
				return err
			}
		}

		c.inMu.Lock()
		if c.ctx.Err() == nil && len(c.in) < streamInQueueDepth &&
			(len(c.in) == 0 || c.inBytes+len(data) <= streamInQueueBytes) {
			c.in = append(c.in, data)
			c.inBytes += len(data)
			c.inCond.Broadcast()
			c.inMu.Unlock()
			return nil
		}
		c.inMu.Unlock()
		if c.manager != nil {
			c.manager.releaseInbound(len(data), 1)
		}
	}
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
	// chunkBytes and chunkParts account for incomplete and retryable chunk sets
	// across every session. Per-session accounting alone can be multiplied by
	// the global session allowance into an impractical host-memory bound.
	chunkBytes int
	chunkParts int
	closed     bool
	closedFlag atomic.Bool

	outBytes    atomic.Int64
	outFrames   atomic.Int64
	outControls atomic.Int64
	outCloses   atomic.Int64
	lifecycles  atomic.Int64
	inBytes     atomic.Int64
	inFrames    atomic.Int64
	budgetMu    sync.Mutex
	budgetCond  *sync.Cond

	janitor sync.Once
	once    sync.Once
	stop    chan struct{}
}

func (m *streamManager) reserveChunkResources(bytes, parts int) bool {
	if bytes < 0 || parts < 0 {
		return false
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed || bytes > streamMaxChunkBytesGlobal-m.chunkBytes ||
		parts > streamMaxChunkPartsGlobal-m.chunkParts {
		return false
	}
	m.chunkBytes += bytes
	m.chunkParts += parts
	return true
}

func (m *streamManager) releaseChunkResources(bytes, parts int) {
	if bytes <= 0 && parts <= 0 {
		return
	}
	m.mu.Lock()
	// All mutations are internal and balanced, but do not let a defensive
	// cleanup path turn an accounting defect into a negative allowance.
	if bytes >= m.chunkBytes {
		m.chunkBytes = 0
	} else {
		m.chunkBytes -= bytes
	}
	if parts >= m.chunkParts {
		m.chunkParts = 0
	} else {
		m.chunkParts -= parts
	}
	m.mu.Unlock()
}

func newStreamManager(app *App) *streamManager {
	m := &streamManager{
		app:            app,
		handlers:       make(map[string]StreamHandler),
		sessions:       make(map[string]*streamSession),
		retiredThrough: make(map[uint]uint64),
		stop:           make(chan struct{}),
	}
	m.budgetCond = sync.NewCond(&m.budgetMu)
	return m
}

func reserveCounter(counter *atomic.Int64, amount, limit int64) bool {
	if amount < 0 || amount > limit {
		return false
	}
	for {
		current := counter.Load()
		if amount > limit-current {
			return false
		}
		if counter.CompareAndSwap(current, current+amount) {
			return true
		}
	}
}

func releaseCounter(counter *atomic.Int64, amount int64) {
	if amount <= 0 {
		return
	}
	for {
		current := counter.Load()
		next := current - amount
		if next < 0 {
			next = 0
		}
		if counter.CompareAndSwap(current, next) {
			return
		}
	}
}

func (m *streamManager) signalBudget() {
	m.budgetMu.Lock()
	m.budgetCond.Broadcast()
	m.budgetMu.Unlock()
}

// waitForBudget waits with budgetMu held. Every cancellation path calls
// signalBudget, so the condition can be checked without a polling timer.
func (m *streamManager) waitForBudget(ctx context.Context) error {
	if m.closedFlag.Load() {
		return ErrStreamClosed
	}
	if ctx != nil && ctx.Err() != nil {
		return ErrStreamClosed
	}
	m.budgetCond.Wait()
	return nil
}

func (m *streamManager) reserveOutbound(kind uint8, bytes int, block bool, ctx context.Context) error {
	m.budgetMu.Lock()
	defer m.budgetMu.Unlock()
	for {
		if m.closedFlag.Load() || (ctx != nil && ctx.Err() != nil) {
			return ErrStreamClosed
		}
		var reserved bool
		switch kind {
		case frameData:
			reserved = reserveCounter(&m.outBytes, int64(bytes), streamOutQueueBytesGlobal)
			if reserved && !reserveCounter(&m.outFrames, 1, streamOutQueueDepthGlobal) {
				releaseCounter(&m.outBytes, int64(bytes))
				reserved = false
			}
		case frameClose:
			// Closes draw on their own allowance rather than inheriting the
			// connection's lifecycle slot. Ownership that moves between two
			// parties has to be handed over atomically; a separate counter has
			// no handover at all, so a close that fails to queue and a shutdown
			// that runs concurrently cannot disagree about who releases what.
			reserved = reserveCounter(&m.outCloses, 1, streamOutCloseDepthGlobal)
		default:
			reserved = reserveCounter(&m.outControls, 1, streamOutControlDepthGlobal)
		}
		if reserved {
			if m.closedFlag.Load() || (ctx != nil && ctx.Err() != nil) {
				switch kind {
				case frameData:
					releaseCounter(&m.outBytes, int64(bytes))
					releaseCounter(&m.outFrames, 1)
				case frameClose:
					releaseCounter(&m.outCloses, 1)
				default:
					releaseCounter(&m.outControls, 1)
				}
				return ErrStreamClosed
			}
			return nil
		}
		if !block {
			return ErrStreamFull
		}
		if err := m.waitForBudget(ctx); err != nil {
			return err
		}
	}
}

// releaseOutbound returns the capacity reserveOutbound acquired for a frame,
// whether that frame was drained, discarded with its session, or never queued
// at all. Every kind releases exactly what it reserved, so the rollback and the
// drain paths are the same call.
func (m *streamManager) releaseOutbound(kind uint8, bytes int) {
	switch kind {
	case frameData:
		releaseCounter(&m.outBytes, int64(bytes))
		releaseCounter(&m.outFrames, 1)
	case frameClose:
		releaseCounter(&m.outCloses, 1)
	default:
		releaseCounter(&m.outControls, 1)
	}
	m.signalBudget()
}

func (m *streamManager) reserveOpen() bool {
	m.budgetMu.Lock()
	defer m.budgetMu.Unlock()
	if m.closedFlag.Load() ||
		!reserveCounter(&m.lifecycles, 1, streamMaxConnectionsGlobal) {
		return false
	}
	if m.closedFlag.Load() ||
		!reserveCounter(&m.outControls, 1, streamOutControlDepthGlobal) {
		releaseCounter(&m.lifecycles, 1)
		return false
	}
	return true
}

func (m *streamManager) releaseOpenReservation() {
	releaseCounter(&m.lifecycles, 1)
	releaseCounter(&m.outControls, 1)
	m.signalBudget()
}

func (m *streamManager) releaseLifecycle() {
	releaseCounter(&m.lifecycles, 1)
	m.signalBudget()
}

func (m *streamManager) reserveLifecycle() bool {
	m.budgetMu.Lock()
	defer m.budgetMu.Unlock()
	if m.closedFlag.Load() ||
		!reserveCounter(&m.lifecycles, 1, streamMaxConnectionsGlobal) {
		return false
	}
	if m.closedFlag.Load() {
		releaseCounter(&m.lifecycles, 1)
		return false
	}
	return true
}

func (m *streamManager) reserveInbound(bytes int, block bool, ctx context.Context) error {
	m.budgetMu.Lock()
	defer m.budgetMu.Unlock()
	for {
		if m.closedFlag.Load() || (ctx != nil && ctx.Err() != nil) {
			return ErrStreamClosed
		}
		reserved := reserveCounter(&m.inBytes, int64(bytes), streamInQueueBytesGlobal)
		if reserved && !reserveCounter(&m.inFrames, 1, streamInQueueDepthGlobal) {
			releaseCounter(&m.inBytes, int64(bytes))
			reserved = false
		}
		if reserved {
			if m.closedFlag.Load() || (ctx != nil && ctx.Err() != nil) {
				releaseCounter(&m.inBytes, int64(bytes))
				releaseCounter(&m.inFrames, 1)
				return ErrStreamClosed
			}
			return nil
		}
		if !block {
			return ErrStreamFull
		}
		if err := m.waitForBudget(ctx); err != nil {
			return err
		}
	}
}

func (m *streamManager) releaseInbound(bytes, frames int) {
	releaseCounter(&m.inBytes, int64(bytes))
	releaseCounter(&m.inFrames, int64(frames))
	m.signalBudget()
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
	if a.streams == nil || name == "" || len(name) > streamMaxNameLen || handler == nil {
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
	s, _ := m.sessionWithAdmission(id, windowID, generation, maySupersede)
	return s
}

// sessionWithAdmission is session plus the one distinction the HTTP boundary
// needs: a fresh id rejected by a resource bound is retryable, whereas a
// retired generation or closed manager is terminal. Keeping that distinction
// out of session preserves the small manager API used by non-HTTP tests.
func (m *streamManager) sessionWithAdmission(id string, windowID uint, generation uint64, maySupersede bool) (*streamSession, bool) {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil, false
	}
	s, ok := m.sessions[id]
	// Resolve an existing id only within the window that owns it. sessionFor
	// reports the eventual 403, but the manager must avoid mutating anything
	// before that check: a foreign poll used to borrow this session's newer
	// generation and retire legitimate sessions in the requesting window.
	if ok && windowID != 0 && s.windowID != windowID {
		m.mu.Unlock()
		return s, false
	}
	// The same rule applies to a generation mismatch. sessionFor reports 403,
	// but an invalid poll must not use the generation stored on the colliding
	// session id to retire other valid sessions before that rejection.
	if ok && s.generation != generation {
		m.mu.Unlock()
		return s, false
	}
	// Zero means the platform could not identify a window. Such requests can
	// come from independent browser clients, whose generation counters have no
	// ordering relationship. Keep them globally bounded, but never let one
	// retire another or advance a shared generation watermark.
	scopedSupersede := maySupersede && windowID != 0
	if !ok {
		// Once a page generation has finished, no delayed request from that page
		// may recreate it. Without this watermark, a late open POST could
		// allocate the old session again after its replacement retired it.
		if windowID != 0 && generation <= m.retiredThrough[windowID] {
			m.mu.Unlock()
			return nil, false
		}
		// A current page poll may replace older sessions in its own window even
		// when a buggy predecessor filled the allowance. The supersede loop below
		// removes those older sessions before this lock is released, so allowing
		// that one replacement does not relax the steady-state bounds.
		replacesOlder := false
		windowSessions := 0
		for _, existing := range m.sessions {
			if existing.windowID == windowID {
				windowSessions++
				if scopedSupersede && existing.generation < generation {
					replacesOlder = true
				}
			}
		}
		if (windowID != 0 && windowSessions >= streamMaxSessionsPerWindow && !replacesOlder) ||
			(len(m.sessions) >= streamMaxSessions && !replacesOlder) {
			m.mu.Unlock()
			return nil, true
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
	if scopedSupersede {
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
	return s, false
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
	m.closedFlag.Store(true)
	doomed := make([]*streamSession, 0, len(m.sessions))
	for _, s := range m.sessions {
		doomed = append(doomed, s)
	}
	m.sessions = map[string]*streamSession{}
	m.retiredThrough = map[uint]uint64{}
	m.mu.Unlock()
	m.signalBudget()

	for _, s := range doomed {
		s.close()
	}

	m.janitor.Do(func() {}) // consume the Once so reap can never start later
	m.once.Do(func() { close(m.stop) })
}
