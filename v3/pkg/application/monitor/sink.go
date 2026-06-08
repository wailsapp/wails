package monitor

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// sink is the active global tap. nil means the monitor is disabled, in which
// case Emit is a cheap no-op.
var sink atomic.Pointer[Sink]

// Enabled reports whether the monitor is currently running. Callers should gate
// any expensive work (e.g. JSON marshaling of arguments) behind this so there
// is zero cost when the monitor is off.
func Enabled() bool { return sink.Load() != nil }

// Emit records a trace. It stamps the sequence number and time, stores the
// record in the replay ring, and broadcasts it to all connected clients wrapped
// in an Envelope. It never blocks the caller: a slow client simply drops records.
func Emit(t Trace) {
	s := sink.Load()
	if s == nil {
		return
	}
	t.Seq = s.seq.Add(1)
	if t.Time.IsZero() {
		t.Time = time.Now()
	}
	line, err := json.Marshal(Envelope{Type: MsgTrace, Trace: &t})
	if err != nil {
		return
	}
	line = append(line, '\n')

	s.mu.Lock()
	// store in ring for replay-on-connect
	if len(s.ring) < s.ringSize {
		s.ring = append(s.ring, line)
	} else {
		s.ring[s.ringHead] = line
		s.ringHead = (s.ringHead + 1) % s.ringSize
	}
	// broadcast (non-blocking)
	for c := range s.clients {
		select {
		case c.ch <- line:
		default:
			c.dropped.Add(1)
		}
	}
	s.mu.Unlock()
}

// Config configures a Sink.
type Config struct {
	// SocketPath is the unix socket to listen on. Required.
	SocketPath string
	// RingSize is the number of recent records replayed to a newly connected
	// client. Defaults to 4096 when <= 0.
	RingSize int
	// ClientBuffer is the per-client send-queue depth. Defaults to 1024.
	ClientBuffer int
}

type client struct {
	ch      chan []byte
	dropped atomic.Uint64
}

// Sink owns the unix listener, the replay ring, and the set of connected
// clients.
type Sink struct {
	socketPath   string
	ringSize     int
	clientBuffer int

	listener net.Listener

	mu       sync.Mutex
	ring     [][]byte
	ringHead int
	clients  map[*client]struct{}

	seq     atomic.Uint64
	closed  atomic.Bool
	closeWg sync.WaitGroup
	sampler *sampler
}

// broadcast marshals an envelope and fans it out to all connected clients
// without touching the replay ring. Used for transient messages (samples) that
// need not be replayed to late joiners. Never blocks: slow clients drop.
func (s *Sink) broadcast(env Envelope) {
	line, err := json.Marshal(env)
	if err != nil {
		return
	}
	line = append(line, '\n')
	s.mu.Lock()
	for c := range s.clients {
		select {
		case c.ch <- line:
		default:
			c.dropped.Add(1)
		}
	}
	s.mu.Unlock()
}

// Start creates the unix listener, installs the global sink, and begins
// accepting clients. Only one sink may be active at a time; starting a second
// replaces (but does not stop) the first, so callers should Stop the previous
// one first.
func Start(cfg Config) (*Sink, error) {
	if cfg.RingSize <= 0 {
		cfg.RingSize = 4096
	}
	if cfg.ClientBuffer <= 0 {
		cfg.ClientBuffer = 1024
	}

	// Remove any stale socket file from a previous crashed run.
	_ = os.Remove(cfg.SocketPath)

	ln, err := net.Listen("unix", cfg.SocketPath)
	if err != nil {
		return nil, err
	}

	s := &Sink{
		socketPath:   cfg.SocketPath,
		ringSize:     cfg.RingSize,
		clientBuffer: cfg.ClientBuffer,
		listener:     ln,
		clients:      make(map[*client]struct{}),
	}

	s.closeWg.Add(1)
	go s.acceptLoop()

	sink.Store(s)
	s.sampler = startSampler(s, time.Second)
	return s, nil
}

func (s *Sink) acceptLoop() {
	defer s.closeWg.Done()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			continue
		}
		s.handleConn(conn)
	}
}

func (s *Sink) handleConn(conn net.Conn) {
	c := &client{ch: make(chan []byte, s.clientBuffer)}

	// Snapshot the ring for replay, then register the client so we don't miss
	// records that arrive between snapshot and registration.
	s.mu.Lock()
	replay := s.snapshotRingLocked()
	s.clients[c] = struct{}{}
	s.mu.Unlock()

	var once sync.Once
	cleanup := func() {
		once.Do(func() {
			s.mu.Lock()
			delete(s.clients, c)
			s.mu.Unlock()
			_ = conn.Close()
		})
	}

	// Reader: handle client requests (e.g. describe).
	go func() {
		defer cleanup()
		sc := bufio.NewScanner(conn)
		sc.Buffer(make([]byte, 0, 4096), 1<<20)
		for sc.Scan() {
			b := sc.Bytes()
			if len(b) == 0 {
				continue
			}
			var req Request
			if err := json.Unmarshal(b, &req); err != nil {
				continue
			}
			if req.Type == ReqDescribe {
				s.sendSnapshot(c)
			}
		}
	}()

	// Writer: replay history, then stream live.
	go func() {
		defer cleanup()
		for _, line := range replay {
			if _, err := conn.Write(line); err != nil {
				return
			}
		}
		for line := range c.ch {
			if _, err := conn.Write(line); err != nil {
				return
			}
		}
	}()
}

// sendSnapshot enqueues a snapshot envelope to a single client.
func (s *Sink) sendSnapshot(c *client) {
	snap := describe()
	if snap == nil {
		snap = &Snapshot{}
	}
	line, err := json.Marshal(Envelope{Type: MsgSnapshot, Snapshot: snap})
	if err != nil {
		return
	}
	line = append(line, '\n')
	select {
	case c.ch <- line:
	default:
		c.dropped.Add(1)
	}
}

func (s *Sink) snapshotRingLocked() [][]byte {
	out := make([][]byte, 0, len(s.ring))
	if len(s.ring) < s.ringSize {
		out = append(out, s.ring...)
		return out
	}
	// full ring: oldest is at ringHead
	for i := 0; i < s.ringSize; i++ {
		out = append(out, s.ring[(s.ringHead+i)%s.ringSize])
	}
	return out
}

// Stop closes the listener, disconnects clients, clears the global sink, and
// removes the socket file. It is idempotent.
func (s *Sink) Stop() {
	if s.closed.Swap(true) {
		return
	}
	sink.CompareAndSwap(s, nil)
	s.sampler.stop()
	_ = s.listener.Close()

	s.mu.Lock()
	for c := range s.clients {
		close(c.ch)
		delete(s.clients, c)
	}
	s.mu.Unlock()

	s.closeWg.Wait()
	_ = os.Remove(s.socketPath)
}

// SocketPath returns the path this sink is listening on.
func (s *Sink) SocketPath() string { return s.socketPath }
