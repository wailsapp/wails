package application

import (
	"context"
	"sync"
	"time"
)

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
	mgr      *streamManager

	mu       sync.Mutex
	space    *sync.Cond // producers waiting for room in out
	out      []outFrame
	outBytes int
	conns    map[uint32]*StreamConn
	closed   bool

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
	s.mu.Lock()
	defer s.mu.Unlock()

	for {
		if s.closed {
			return ErrStreamClosed
		}
		// A closed connection may still emit its own close frame — that is how
		// the frontend learns — but nothing else.
		if kind != frameClose && kind != frameError && c != nil && c.ctx.Err() != nil {
			return ErrStreamClosed
		}
		// An empty queue always accepts one frame, however large. Enforcing the
		// byte cap unconditionally would make a frame bigger than the cap
		// impossible to send at all: the condition could never come true, so a
		// blocking Send would wait forever and TrySend would report full
		// permanently. Frame size is not something the caller necessarily
		// controls — a struct with a []byte field marshals to whatever it
		// marshals to — so the cap governs how much may accumulate, not how
		// large any single frame may be. Same principle as drainLocked always
		// taking at least one frame.
		if len(s.out) < streamOutQueueDepth &&
			(len(s.out) == 0 || s.outBytes+len(data) <= streamOutQueueBytes) {
			break
		}
		if !block {
			return ErrStreamFull
		}
		s.space.Wait()
	}

	s.out = append(s.out, outFrame{connID: connID, kind: kind, data: data})
	s.outBytes += len(data)
	s.wake()
	return nil
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
		s.outBytes -= len(s.out[i].data)
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
func (s *streamSession) open(connID uint32, name string) {
	handler, ok := s.mgr.handler(name)
	if !ok {
		_ = s.enqueue(nil, connID, frameError, []byte("no handler registered for stream "+name), false)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	c := &StreamConn{
		id:       connID,
		name:     name,
		windowID: s.windowID,
		sink:     s,
		ctx:      ctx,
		cancel:   cancel,
	}
	c.inCond = sync.NewCond(&c.inMu)

	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		cancel()
		return
	}
	if _, exists := s.conns[connID]; exists {
		// The frontend allocates connection ids; a duplicate means a buggy or
		// hostile client. Refuse rather than replacing a live connection.
		s.mu.Unlock()
		cancel()
		return
	}
	s.conns[connID] = c
	s.mu.Unlock()

	_ = s.enqueue(c, connID, frameOpen, nil, false)

	go func() {
		defer handlePanic()
		defer c.shutdown()
		handler(c)
	}()
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
	// Release buffered payloads immediately; nothing will ever poll for them.
	s.out = nil
	s.outBytes = 0
	s.space.Broadcast()
	s.mu.Unlock()

	s.wake()

	for _, c := range conns {
		c.shutdown()
	}
}
