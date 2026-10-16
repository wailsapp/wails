package application

import (
	"context"
	"errors"
	"sync"
	"time"
)

var errStreamDuplicateConnection = errors.New("wails: duplicate stream connection id")

// Frame kinds, on the wire in both directions.
const (
	frameData  uint8 = 0 // payload
	frameOpen  uint8 = 1 // JS→Go: connect to a named stream. Go→JS: accepted.
	frameClose uint8 = 2 // either direction: this connection is finished
	frameError uint8 = 3 // Go→JS: connect refused, payload is the reason
)

type outFrame struct {
	connID uint32
	kind   uint8
	data   []byte
}

// streamSession is one page load in one window. Every connection that page
// opens is multiplexed over the session's single held poll.
//
// One poll at a time, one queue, one drainer. That is what makes ordering
// correct by construction rather than by argument: frames leave in the order
// they were appended, whichever goroutine appended them, and there is never a
// second delivery mechanism that could overtake the first.
type streamSession struct {
	id       string
	windowID uint
	// generation orders page sessions within a window. It is assigned by the
	// client when the session is first created and never changes. A poll may
	// retire only an older generation, so a delayed request from the previous
	// page cannot tear down its replacement.
	generation uint64
	mgr        *streamManager

	mu       sync.Mutex
	space    *sync.Cond // producers waiting for room in out
	out      []outFrame
	outBytes int
	// Data and protocol frames have independent bounds. Control frames cannot
	// consume data capacity, and close frames have a reserved allowance so an
	// accepted connection can always report its end.
	outDataFrames int
	outControls   int
	outCloses     int
	conns         map[uint32]*StreamConn
	closed        bool

	// pollGen supersedes an earlier poll when a new one arrives. A page that
	// reloads leaves its predecessor's poll parked in the platform layer with
	// no way to cancel it; bumping the generation makes that request return
	// 204 as soon as the new page polls, instead of holding until timeout.
	pollGen uint64

	// notify carries "something changed" to the parked poll. Buffered at 1: a
	// burst of appends collapses into one wake-up, and the poll re-reads state
	// under the lock anyway.
	notify chan struct{}

	lastSeen time.Time

	// chunkStore reassembles oversized frames from the frontend. Created on
	// first use, since most sessions never send one.
	chunkStore *streamChunkStore
}

func newStreamSession(id string, windowID uint, mgr *streamManager) *streamSession {
	s := &streamSession{
		id:       id,
		windowID: windowID,
		mgr:      mgr,
		conns:    make(map[uint32]*StreamConn),
		notify:   make(chan struct{}, 1),
		lastSeen: time.Now(),
	}
	s.space = sync.NewCond(&s.mu)
	return s
}

func (s *streamSession) touch() {
	s.mu.Lock()
	s.lastSeen = time.Now()
	s.mu.Unlock()
}

func (s *streamSession) lastSeenAt() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastSeen
}

// wake nudges the parked poll without blocking, whatever its state.
func (s *streamSession) wake() {
	select {
	case s.notify <- struct{}{}:
	default:
	}
}

// enqueue appends a frame for delivery. When block is true it waits for room
// the way a socket write waits on a full send buffer; otherwise it reports
// ErrStreamFull.
//
// The bound is per session rather than per connection, which is the honest
// consequence of multiplexing: connections in one window share one pipe, so a
// connection that floods it does slow its neighbours. Per-connection fairness
// would need a scheduler, and there is no evidence yet that it is needed.
// The connection id is passed explicitly rather than read from c, so a frame
// can be addressed to a connection that was never created — the refusal path,
// where there is no StreamConn to carry the id. Deriving it from c and patching
// the frame afterwards would need a second acquisition of s.mu, and a poll
// draining in between would either send the frame with id 0 (leaving the
// frontend stuck in CONNECTING with no error) or retag an unrelated frame.
func (s *streamSession) enqueue(c *StreamConn, connID uint32, kind uint8, data []byte, block bool) error {
	// Control frames bypass the data caps, but not their own bounds. Dropping one
	// is a protocol failure: a lost open ack leaves the frontend in CONNECTING
	// forever, and a lost close leaves it believing an ended connection is live.
	// Separate close capacity plus the admission invariant in open guarantees
	// one reliable close slot for every accepted connection.
	control := kind != frameData

	for {
		s.mu.Lock()
		if s.closed {
			s.mu.Unlock()
			return ErrStreamClosed
		}
		// A closed connection may still emit its own close frame — that is how
		// the frontend learns — but nothing else.
		if !control && c != nil && c.ctx.Err() != nil {
			s.mu.Unlock()
			return ErrStreamClosed
		}

		localRoom := false
		if control {
			if kind == frameClose {
				localRoom = s.outCloses < streamMaxConnections
			} else {
				localRoom = s.outControls < streamOutControlDepth
			}
		} else if s.outDataFrames < streamOutQueueDepth &&
			(s.outDataFrames == 0 || s.outBytes+len(data) <= streamOutQueueBytes) {
			localRoom = true
		}

		if !localRoom {
			if !block {
				s.mu.Unlock()
				return ErrStreamFull
			}
			s.space.Wait()
			s.mu.Unlock()
			continue
		}
		s.mu.Unlock()

		var ctx context.Context
		if !control && c != nil {
			ctx = c.ctx
		}
		if err := s.mgr.reserveOutbound(kind, len(data), block, ctx); err != nil {
			return err
		}

		// Local and global capacity are protected independently. Recheck the
		// session after reserving the shared allowance: another producer may have
		// filled it while this producer was contending for global room.
		s.mu.Lock()
		if s.closed || (!control && c != nil && c.ctx.Err() != nil) {
			s.mu.Unlock()
			s.mgr.releaseOutbound(kind, len(data))
			return ErrStreamClosed
		}
		if control {
			if kind == frameClose {
				localRoom = s.outCloses < streamMaxConnections
			} else {
				localRoom = s.outControls < streamOutControlDepth
			}
		} else {
			localRoom = s.outDataFrames < streamOutQueueDepth &&
				(s.outDataFrames == 0 || s.outBytes+len(data) <= streamOutQueueBytes)
		}
		if !localRoom {
			s.mu.Unlock()
			s.mgr.releaseOutbound(kind, len(data))
			if !block {
				return ErrStreamFull
			}
			continue
		}

		// The frame retains the caller's slice. Copying here would be a full memcpy
		// per frame, and at gigabytes per second that dominates the transport for
		// no benefit that callers cannot get themselves: almost every caller hands
		// over freshly marshalled bytes. The ownership rule is documented on Send
		// instead — do not mutate a slice after passing it.
		s.out = append(s.out, outFrame{connID: connID, kind: kind, data: data})
		s.outBytes += len(data)
		s.countQueued(kind, 1)
		s.mu.Unlock()
		s.wake()
		return nil
	}
}

func (s *streamSession) countQueued(kind uint8, delta int) {
	switch kind {
	case frameData:
		s.outDataFrames += delta
	case frameClose:
		s.outCloses += delta
	default:
		s.outControls += delta
	}
}

// awaitFrames parks until there is something to send, the hold expires, this
// poll is superseded, or the session dies. alive is false once the session is
// gone, which the handler reports so the frontend can stop polling.
func (s *streamSession) awaitFrames(ctx context.Context, hold time.Duration, maxBytes int) (frames []outFrame, more bool, alive bool) {
	timer := time.NewTimer(hold)
	defer timer.Stop()

	s.mu.Lock()
	s.pollGen++
	gen := s.pollGen
	s.lastSeen = time.Now()
	s.mu.Unlock()

	// Any poll already parked belongs to a page that has moved on.
	s.wake()

	for {
		s.mu.Lock()
		if s.closed {
			s.mu.Unlock()
			return nil, false, false
		}
		if gen != s.pollGen {
			// Superseded by a newer poll: return empty rather than stealing
			// frames the newer poll should carry.
			s.mu.Unlock()
			return nil, false, true
		}
		if len(s.out) > 0 {
			frames, more = s.drainLocked(maxBytes)
			s.mu.Unlock()
			return frames, more, true
		}
		s.mu.Unlock()

		select {
		case <-s.notify:
		case <-timer.C:
			return nil, false, true
		case <-ctx.Done():
			return nil, false, true
		}
	}
}

// drainLocked removes as many frames as fit in maxBytes. It always takes at
// least one frame, so a single frame larger than the cap is still delivered
// rather than wedging the queue forever.
func (s *streamSession) drainLocked(maxBytes int) (frames []outFrame, more bool) {
	total := 0
	n := 0
	for n < len(s.out) {
		size := len(s.out[n].data) + streamFrameHeaderBytes
		if n > 0 && total+size > maxBytes {
			break
		}
		total += size
		n++
	}

	frames = make([]outFrame, n)
	copy(frames, s.out[:n])
	for i := 0; i < n; i++ {
		s.mgr.releaseOutbound(s.out[i].kind, len(s.out[i].data))
		s.outBytes -= len(s.out[i].data)
		s.countQueued(s.out[i].kind, -1)
		s.out[i] = outFrame{}
	}
	s.out = s.out[n:]
	more = len(s.out) > 0

	// Whatever we just removed is room for a blocked producer.
	s.space.Broadcast()
	return frames, more
}

// open accepts a connection request from the frontend. The ack is queued before
// the handler starts, so the frontend cannot see data from a stream it has not
// yet been told is open.
func (s *streamSession) open(connID uint32, name string) error {
	handler, ok := s.mgr.handler(name)
	if !ok {
		return s.enqueue(nil, connID, frameError, []byte("no handler registered for stream "+name), false)
	}
	if !s.mgr.reserveOpen() {
		return ErrStreamFull
	}

	ctx, cancel := context.WithCancel(context.Background())
	c := &StreamConn{
		id:        connID,
		name:      name,
		windowID:  s.windowID,
		sink:      s,
		ctx:       ctx,
		cancel:    cancel,
		manager:   s.mgr,
		lifecycle: true,
	}
	c.inCond = sync.NewCond(&c.inMu)

	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		cancel()
		s.mgr.releaseOpenReservation()
		return ErrStreamClosed
	}
	if _, exists := s.conns[connID]; exists {
		// The frontend allocates connection ids; a duplicate means a buggy or
		// hostile client. Refuse rather than replacing a live connection.
		s.mu.Unlock()
		cancel()
		s.mgr.releaseOpenReservation()
		return errStreamDuplicateConnection
	}
	// A queued close still owns its connection slot until the frontend drains
	// that notification. This invariant reserves enough close-frame capacity
	// even when short-lived connection ids are recycled without any polling.
	if len(s.conns)+s.outCloses >= streamMaxConnections || s.outControls >= streamOutControlDepth {
		s.mu.Unlock()
		cancel()
		s.mgr.releaseOpenReservation()
		return ErrStreamFull
	}
	s.conns[connID] = c
	// Queue the acknowledgement in the same critical section as registration.
	// A poll can therefore never deliver the ack before data requests can find
	// the connection, and no competing open can consume its reserved slot.
	s.out = append(s.out, outFrame{connID: connID, kind: frameOpen})
	s.outControls++
	s.mu.Unlock()
	s.wake()

	go func() {
		defer handlePanic()
		// Close, not shutdown: the API says returning from the handler closes
		// the connection, and the frontend has to be told or its socket stays
		// open with onclose never firing. Close is idempotent, so a handler
		// that closed explicitly pays nothing here, and a connection the peer
		// already closed has had the once consumed by closedByPeer.
		defer c.Close()
		handler(c)
	}()
	return nil
}

// connCount reports how many connections the session still has.
func (s *streamSession) connCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.conns)
}

// wakeProducers releases anything parked in enqueue waiting for room.
func (s *streamSession) wakeProducers() {
	s.mu.Lock()
	s.space.Broadcast()
	s.mu.Unlock()
}

func (s *streamSession) conn(connID uint32) *StreamConn {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conns[connID]
}

func (s *streamSession) removeConn(connID uint32) {
	s.mu.Lock()
	delete(s.conns, connID)
	s.mu.Unlock()
}

// close ends the session and every connection on it. Handlers blocked in
// Receive or Send unblock with ErrStreamClosed.
func (s *streamSession) close() {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return
	}
	s.closed = true
	conns := make([]*StreamConn, 0, len(s.conns))
	for _, c := range s.conns {
		conns = append(conns, c)
	}
	// Release buffered payloads and their application-wide reservations
	// immediately; nothing will ever poll for them.
	for i := range s.out {
		s.mgr.releaseOutbound(s.out[i].kind, len(s.out[i].data))
		s.out[i] = outFrame{}
	}
	s.out = nil
	s.outBytes = 0
	s.outDataFrames = 0
	s.outControls = 0
	s.outCloses = 0
	chunks := s.chunkStore
	s.space.Broadcast()
	s.mu.Unlock()
	if chunks != nil {
		chunks.close()
	}

	s.wake()

	for _, c := range conns {
		c.shutdown()
	}
}
