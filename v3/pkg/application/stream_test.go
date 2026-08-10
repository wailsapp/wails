package application

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync"
	"testing"
	"time"
)

// newTestSession builds a session with no App behind it. Nothing exercised here
// needs one; StreamConn.Window is the only method that does, and it is not
// under test.
func newTestSession(t *testing.T) (*streamManager, *streamSession) {
	t.Helper()
	mgr := newStreamManager(nil)
	s := newStreamSession("sess", 1, mgr)
	mgr.sessions["sess"] = s
	t.Cleanup(func() { mgr.close() })
	return mgr, s
}

func newTestConn(s *streamSession, id uint32) *StreamConn {
	ctx, cancel := context.WithCancel(context.Background())
	c := &StreamConn{id: id, name: "test", sink: s, ctx: ctx, cancel: cancel}
	c.inCond = sync.NewCond(&c.inMu)
	s.mu.Lock()
	s.conns[id] = c
	s.mu.Unlock()
	return c
}

func TestStreamFrameRoundTrip(t *testing.T) {
	frames := []outFrame{
		{connID: 1, kind: frameOpen, data: nil},
		{connID: 1, kind: frameData, data: []byte("hello")},
		{connID: 7, kind: frameData, data: make([]byte, 70000)},
		{connID: 7, kind: frameClose, data: nil},
	}

	body := encodeStreamFrames(frames, true)

	if got := string(body[:4]); got != "WS1\x00" {
		t.Fatalf("magic = %q, want WS1\\0", got)
	}
	if body[4]&1 != 1 {
		t.Fatalf("more flag not set")
	}
	if n := binary.BigEndian.Uint32(body[5:9]); n != uint32(len(frames)) {
		t.Fatalf("count = %d, want %d", n, len(frames))
	}

	off := 9
	for i, want := range frames {
		connID := binary.BigEndian.Uint32(body[off : off+4])
		kind := body[off+4]
		length := binary.BigEndian.Uint32(body[off+5 : off+9])
		off += 9
		payload := body[off : off+int(length)]
		off += int(length)

		if connID != want.connID || kind != want.kind || int(length) != len(want.data) {
			t.Fatalf("frame %d header = (%d,%d,%d), want (%d,%d,%d)",
				i, connID, kind, length, want.connID, want.kind, len(want.data))
		}
		if string(payload) != string(want.data) {
			t.Fatalf("frame %d payload mismatch", i)
		}
	}
	if off != len(body) {
		t.Fatalf("trailing bytes: consumed %d of %d", off, len(body))
	}
}

// TestStreamOrderingUnderConcurrentSenders is the regression test for the
// failure that cost the event path 4.4% of its events: a send from one
// goroutine overtaking an earlier send from another. Frames must leave in the
// order Send accepted them, whichever goroutine called it.
func TestStreamOrderingUnderConcurrentSenders(t *testing.T) {
	_, s := newTestSession(t)
	c := newTestConn(s, 1)

	const senders = 8
	const perSender = 200

	// A single counter handed out under the session lock is the ground truth:
	// whatever order Send accepted, the queue must reproduce.
	var mu sync.Mutex
	var next int
	var accepted []int

	var wg sync.WaitGroup
	for i := 0; i < senders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perSender; j++ {
				mu.Lock()
				seq := next
				next++
				payload := []byte(fmt.Sprintf("%d", seq))
				if err := s.enqueue(c, 1, frameData, payload, true); err != nil {
					mu.Unlock()
					t.Errorf("enqueue: %v", err)
					return
				}
				accepted = append(accepted, seq)
				mu.Unlock()
			}
		}()
	}

	// Drain concurrently, the way a poll would.
	var drained []int
	done := make(chan struct{})
	go func() {
		defer close(done)
		for len(drained) < senders*perSender {
			s.mu.Lock()
			frames, _ := s.drainLocked(streamMaxResponseBytes)
			s.mu.Unlock()
			for _, f := range frames {
				var v int
				_, _ = fmt.Sscanf(string(f.data), "%d", &v)
				drained = append(drained, v)
			}
			if len(frames) == 0 {
				time.Sleep(time.Millisecond)
			}
		}
	}()

	wg.Wait()
	<-done

	if len(drained) != len(accepted) {
		t.Fatalf("drained %d frames, accepted %d", len(drained), len(accepted))
	}
	for i := range accepted {
		if drained[i] != accepted[i] {
			t.Fatalf("frame %d out of order: got %d, want %d", i, drained[i], accepted[i])
		}
	}
}

func TestStreamTrySendReportsFullAndSendBlocks(t *testing.T) {
	_, s := newTestSession(t)
	c := newTestConn(s, 1)

	for i := 0; i < streamOutQueueDepth; i++ {
		if err := s.enqueue(c, 1, frameData, []byte("x"), false); err != nil {
			t.Fatalf("fill %d: %v", i, err)
		}
	}

	if err := s.enqueue(c, 1, frameData, []byte("x"), false); err != ErrStreamFull {
		t.Fatalf("TrySend on a full queue = %v, want ErrStreamFull", err)
	}

	// A blocking send waits for room rather than failing.
	blocked := make(chan error, 1)
	go func() { blocked <- s.enqueue(c, 1, frameData, []byte("y"), true) }()

	select {
	case err := <-blocked:
		t.Fatalf("blocking send returned early: %v", err)
	case <-time.After(50 * time.Millisecond):
	}

	s.mu.Lock()
	s.drainLocked(streamMaxResponseBytes)
	s.mu.Unlock()

	select {
	case err := <-blocked:
		if err != nil {
			t.Fatalf("blocking send after drain: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("blocking send never woke after the queue drained")
	}
}

func TestStreamByteCapBindsBeforeDepth(t *testing.T) {
	_, s := newTestSession(t)
	c := newTestConn(s, 1)

	// Frames big enough that the byte cap is reached well before the depth cap.
	big := make([]byte, streamOutQueueBytes/4)
	for i := 0; i < 4; i++ {
		if err := s.enqueue(c, 1, frameData, big, false); err != nil {
			t.Fatalf("enqueue %d: %v", i, err)
		}
	}
	if err := s.enqueue(c, 1, frameData, big, false); err != ErrStreamFull {
		t.Fatalf("over byte cap = %v, want ErrStreamFull", err)
	}

	s.mu.Lock()
	depth := len(s.out)
	s.mu.Unlock()
	if depth >= streamOutQueueDepth {
		t.Fatalf("depth cap reached first (%d frames); byte cap did not bind", depth)
	}
}

// A frame larger than the per-window byte cap must still be sendable. Frame
// size is not always the caller's choice — a struct with a []byte field
// marshals to whatever it marshals to — so enforcing the cap unconditionally
// made a large frame impossible to send at all rather than merely slow: the
// wait condition could never come true.
func TestStreamSingleFrameLargerThanByteCap(t *testing.T) {
	_, s := newTestSession(t)
	c := newTestConn(s, 1)

	huge := make([]byte, streamOutQueueBytes*2) // 16 MB against an 8 MB cap

	done := make(chan error, 1)
	go func() { done <- s.enqueue(c, 1, frameData, huge, true) }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("blocking send of an oversized frame: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("blocking send of a frame larger than the byte cap never returned")
	}

	// It occupies the queue, so the next frame waits behind it as normal.
	if err := s.enqueue(c, 1, frameData, []byte("x"), false); err != ErrStreamFull {
		t.Fatalf("send behind an oversized frame = %v, want ErrStreamFull", err)
	}

	s.mu.Lock()
	frames, _ := s.drainLocked(streamMaxResponseBytes)
	s.mu.Unlock()
	if len(frames) != 1 || len(frames[0].data) != len(huge) {
		t.Fatalf("oversized frame not delivered intact")
	}

	// And a non-blocking send succeeds again once it has drained.
	if err := s.enqueue(c, 1, frameData, []byte("x"), false); err != nil {
		t.Fatalf("send after drain: %v", err)
	}
}

// A frame larger than the response cap must still be delivered, or it would
// wedge the queue forever.
func TestStreamDrainAlwaysTakesOneFrame(t *testing.T) {
	_, s := newTestSession(t)
	c := newTestConn(s, 1)

	oversized := make([]byte, streamMaxResponseBytes*2)
	if err := s.enqueue(c, 1, frameData, oversized, false); err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	s.mu.Lock()
	frames, more := s.drainLocked(streamMaxResponseBytes)
	s.mu.Unlock()

	if len(frames) != 1 {
		t.Fatalf("drained %d frames, want 1", len(frames))
	}
	if len(frames[0].data) != len(oversized) {
		t.Fatalf("payload truncated: %d of %d", len(frames[0].data), len(oversized))
	}
	if more {
		t.Fatalf("more set with an empty queue")
	}
}

func TestStreamDrainRespectsResponseCap(t *testing.T) {
	_, s := newTestSession(t)
	c := newTestConn(s, 1)

	// Sized so exactly three fit once each frame's 9-byte header is counted.
	chunk := make([]byte, streamMaxResponseBytes/3-streamFrameHeaderBytes)
	for i := 0; i < 5; i++ {
		if err := s.enqueue(c, 1, frameData, chunk, false); err != nil {
			t.Fatalf("enqueue %d: %v", i, err)
		}
	}

	s.mu.Lock()
	frames, more := s.drainLocked(streamMaxResponseBytes)
	s.mu.Unlock()

	if len(frames) != 3 {
		t.Fatalf("drained %d frames, want 3 to fit the cap", len(frames))
	}
	if !more {
		t.Fatalf("more not set with frames still queued")
	}
}

// A new poll must supersede the one belonging to the page that has navigated
// away, rather than both waiting and one stealing the other's frames.
func TestStreamPollSupersededByNewer(t *testing.T) {
	_, s := newTestSession(t)

	first := make(chan int, 1)
	go func() {
		frames, _, alive := s.awaitFrames(context.Background(), 5*time.Second, streamMaxResponseBytes)
		if !alive {
			first <- -1
			return
		}
		first <- len(frames)
	}()

	// Let the first poll park.
	time.Sleep(50 * time.Millisecond)

	go func() {
		_, _, _ = s.awaitFrames(context.Background(), 5*time.Second, streamMaxResponseBytes)
	}()

	select {
	case n := <-first:
		if n != 0 {
			t.Fatalf("superseded poll returned %d frames, want 0", n)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("superseded poll did not return; it held until timeout")
	}
}

func TestStreamPollReturnsOnHoldTimeout(t *testing.T) {
	_, s := newTestSession(t)

	start := time.Now()
	frames, _, alive := s.awaitFrames(context.Background(), 100*time.Millisecond, streamMaxResponseBytes)
	elapsed := time.Since(start)

	if !alive {
		t.Fatal("session reported dead")
	}
	if len(frames) != 0 {
		t.Fatalf("got %d frames from an empty session", len(frames))
	}
	if elapsed < 90*time.Millisecond {
		t.Fatalf("returned after %v, before the hold expired", elapsed)
	}
}

// Closing the session must unblock everything: a handler parked in Receive, and
// a producer parked on a full buffer.
func TestStreamSessionCloseUnblocksHandlers(t *testing.T) {
	_, s := newTestSession(t)
	c := newTestConn(s, 1)

	recvErr := make(chan error, 1)
	go func() {
		_, err := c.Receive()
		recvErr <- err
	}()

	for i := 0; i < streamOutQueueDepth; i++ {
		_ = s.enqueue(c, 1, frameData, []byte("x"), false)
	}
	sendErr := make(chan error, 1)
	go func() { sendErr <- s.enqueue(c, 1, frameData, []byte("y"), true) }()

	time.Sleep(50 * time.Millisecond)
	s.close()

	select {
	case err := <-recvErr:
		if err != ErrStreamClosed {
			t.Fatalf("Receive = %v, want ErrStreamClosed", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Receive never unblocked after close")
	}

	select {
	case err := <-sendErr:
		if err != ErrStreamClosed {
			t.Fatalf("blocked Send = %v, want ErrStreamClosed", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("blocked Send never unblocked after close")
	}
}

func TestStreamReceiveDeliversInOrder(t *testing.T) {
	_, s := newTestSession(t)
	c := newTestConn(s, 1)

	for i := 0; i < 100; i++ {
		c.deliver([]byte(fmt.Sprintf("%d", i)))
	}
	for i := 0; i < 100; i++ {
		frame, err := c.Receive()
		if err != nil {
			t.Fatalf("Receive %d: %v", i, err)
		}
		if string(frame) != fmt.Sprintf("%d", i) {
			t.Fatalf("frame %d = %q", i, frame)
		}
	}
}

func TestStreamDropWindowClosesOnlyThatWindow(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)

	s1 := newStreamSession("a", 1, mgr)
	s2 := newStreamSession("b", 2, mgr)
	mgr.sessions["a"] = s1
	mgr.sessions["b"] = s2

	c1 := newTestConn(s1, 1)
	c2 := newTestConn(s2, 1)

	mgr.dropWindow(1)

	if c1.ctx.Err() == nil {
		t.Fatal("connection in the dropped window is still live")
	}
	if c2.ctx.Err() != nil {
		t.Fatal("connection in another window was closed too")
	}
	if mgr.existingSession("a") != nil {
		t.Fatal("dropped session still registered")
	}
	if mgr.existingSession("b") == nil {
		t.Fatal("unrelated session was removed")
	}
}

func TestStreamOpenRefusedWithoutHandler(t *testing.T) {
	_, s := newTestSession(t)

	s.open(42, "nosuchstream")

	s.mu.Lock()
	frames, _ := s.drainLocked(streamMaxResponseBytes)
	_, registered := s.conns[42]
	s.mu.Unlock()

	if registered {
		t.Fatal("connection registered for a stream with no handler")
	}
	if len(frames) != 1 || frames[0].kind != frameError {
		t.Fatalf("frames = %+v, want a single error frame", frames)
	}
	if frames[0].connID != 42 {
		t.Fatalf("error frame carried connID %d, want 42", frames[0].connID)
	}
}

// The refusal frame must carry the right connection id even when a poll drains
// between it being queued and anything else happening. The id used to be
// patched onto the last queued frame under a second lock acquisition, so a
// drain in between either shipped the frame with id 0 — leaving the frontend in
// CONNECTING with no error — or retagged an unrelated connection's frame.
func TestStreamOpenRefusalTaggedUnderConcurrentDrain(t *testing.T) {
	for attempt := 0; attempt < 200; attempt++ {
		_, s := newTestSession(t)

		drained := make(chan []outFrame, 1)
		go func() {
			frames, _, _ := s.awaitFrames(context.Background(), 2*time.Second, streamMaxResponseBytes)
			drained <- frames
		}()

		s.open(77, "nosuchstream")

		select {
		case frames := <-drained:
			if len(frames) != 1 {
				t.Fatalf("attempt %d: drained %d frames, want 1", attempt, len(frames))
			}
			if frames[0].kind != frameError {
				t.Fatalf("attempt %d: kind = %d, want frameError", attempt, frames[0].kind)
			}
			if frames[0].connID != 77 {
				t.Fatalf("attempt %d: refusal addressed to conn %d, want 77", attempt, frames[0].connID)
			}
		case <-time.After(3 * time.Second):
			t.Fatalf("attempt %d: refusal never delivered", attempt)
		}
	}
}

// A reload gives the window a new session id. The previous page is gone, so its
// session must go with it rather than surviving to the TTL sweep - which would
// leave the app holding two live connections to the same stream for
// streamSessionTTL + streamSessionSweep.
func TestStreamNewSessionSupersedesPreviousInSameWindow(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	mgr.handlers["s"] = func(c *StreamConn) { <-c.Context().Done() }

	first := mgr.session("page-1", 1)
	if first == nil {
		t.Fatal("first session not created")
	}
	first.open(1, "s")
	c := first.conn(1)
	if c == nil {
		t.Fatal("connection not registered on the first session")
	}

	// Another window's session must be untouched by the reload below.
	other := mgr.session("other-window", 2)
	other.open(1, "s")
	otherConn := other.conn(1)

	second := mgr.session("page-2", 1)
	if second == first {
		t.Fatal("reload reused the previous session")
	}

	if c.ctx.Err() == nil {
		t.Fatal("previous page's connection is still live after reload")
	}
	if mgr.existingSession("page-1") != nil {
		t.Fatal("previous session still registered after reload")
	}
	if otherConn.ctx.Err() != nil {
		t.Fatal("reload in one window closed another window's connection")
	}
	if mgr.existingSession("other-window") == nil {
		t.Fatal("reload in one window removed another window's session")
	}
}

func TestStreamOpenAckPrecedesHandlerOutput(t *testing.T) {
	mgr, s := newTestSession(t)

	started := make(chan struct{})
	mgr.handlers["greet"] = func(c *StreamConn) {
		close(started)
		_ = c.Send([]byte("first"))
		<-c.Context().Done()
	}

	s.open(9, "greet")
	<-started

	// Give the handler's Send time to land.
	time.Sleep(50 * time.Millisecond)

	s.mu.Lock()
	frames, _ := s.drainLocked(streamMaxResponseBytes)
	s.mu.Unlock()

	if len(frames) < 2 {
		t.Fatalf("got %d frames, want the ack plus the handler's frame", len(frames))
	}
	if frames[0].kind != frameOpen {
		t.Fatalf("first frame kind = %d, want the open ack (%d)", frames[0].kind, frameOpen)
	}
	if frames[1].kind != frameData || string(frames[1].data) != "first" {
		t.Fatalf("second frame = %+v, want the handler's data", frames[1])
	}
}

func TestStreamDuplicateConnIDRefused(t *testing.T) {
	mgr, s := newTestSession(t)
	mgr.handlers["dup"] = func(c *StreamConn) { <-c.Context().Done() }

	s.open(5, "dup")
	first := s.conn(5)
	if first == nil {
		t.Fatal("first open did not register")
	}

	s.open(5, "dup")
	if s.conn(5) != first {
		t.Fatal("duplicate open replaced a live connection")
	}
}

func TestStreamChunkAssembly(t *testing.T) {
	store := &streamChunkStore{items: make(map[string]*streamChunkSet)}

	parts := [][]byte{[]byte("one "), []byte("two "), []byte("three")}
	// Deliver out of order: only the last arrival completes the frame, and the
	// result must still be in index order.
	if got, done := store.add("c1", 2, 3, parts[2]); done || got != nil {
		t.Fatal("completed early on the first chunk")
	}
	if _, done := store.add("c1", 0, 3, parts[0]); done {
		t.Fatal("completed early on the second chunk")
	}
	got, done := store.add("c1", 1, 3, parts[1])
	if !done {
		t.Fatal("did not complete on the final chunk")
	}
	if string(got) != "one two three" {
		t.Fatalf("assembled = %q", got)
	}
	if len(store.items) != 0 || store.bytes != 0 {
		t.Fatalf("store not drained: %d items, %d bytes", len(store.items), store.bytes)
	}
}

func TestStreamChunkRejectsInconsistentTotal(t *testing.T) {
	store := &streamChunkStore{items: make(map[string]*streamChunkSet)}

	if _, done := store.add("c1", 0, 3, []byte("a")); done {
		t.Fatal("completed on the first of three chunks")
	}
	if _, done := store.add("c1", 1, 4, []byte("b")); done {
		t.Fatal("accepted a chunk claiming a different total")
	}
	if len(store.items) != 0 {
		t.Fatal("inconsistent set was not discarded")
	}
}
