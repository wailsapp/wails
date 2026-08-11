package application

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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

type blockingCloseSink struct {
	closeStarted chan struct{}
	releaseClose chan struct{}
}

func (s *blockingCloseSink) enqueue(c *StreamConn, _ uint32, kind uint8, _ []byte, _ bool) error {
	if kind == frameClose {
		close(s.closeStarted)
		<-s.releaseClose
		return nil
	}
	if c.ctx.Err() != nil {
		return ErrStreamClosed
	}
	return nil
}

func (*blockingCloseSink) removeConn(uint32) {}
func (*blockingCloseSink) wake()             {}
func (*blockingCloseSink) wakeProducers()    {}

func TestStreamFrameRoundTrip(t *testing.T) {
	frames := []outFrame{
		{connID: 1, kind: frameOpen, data: nil},
		{connID: 1, kind: frameData, data: []byte("hello")},
		{connID: 7, kind: frameData, data: make([]byte, 70000)},
		{connID: 7, kind: frameClose, data: nil},
	}

	body := encodeStreamFrames(nil, frames, true)

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

func TestStreamCloseRejectsSendRacingBehindCloseFrame(t *testing.T) {
	sink := &blockingCloseSink{
		closeStarted: make(chan struct{}),
		releaseClose: make(chan struct{}),
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := &StreamConn{id: 1, sink: sink, ctx: ctx, cancel: cancel}
	c.inCond = sync.NewCond(&c.inMu)

	closed := make(chan struct{})
	go func() {
		_ = c.Close()
		close(closed)
	}()
	<-sink.closeStarted

	if err := c.Send([]byte("too late")); err != ErrStreamClosed {
		t.Fatalf("Send racing behind close frame = %v, want ErrStreamClosed", err)
	}
	close(sink.releaseClose)
	<-closed
}

func TestStreamBatchRejectsTrailingBytes(t *testing.T) {
	body := binary.BigEndian.AppendUint32(nil, 1)
	body = binary.BigEndian.AppendUint32(body, 3)
	body = append(body, []byte("one")...)
	body = append(body, 0xff)

	if _, err := decodeStreamBatch(body); err == nil {
		t.Fatal("batch with trailing bytes was accepted")
	}
}

func TestStreamSendRejectsChunkedBatch(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, streamPathSend, nil)
	req.Header.Set(streamHeaderConn, "1")
	req.Header.Set(streamHeaderKind, strconv.Itoa(int(frameData)))
	req.Header.Set(streamHeaderBatch, "1")
	req.Header.Set(streamHeaderChunkID, "chunked-batch")
	rw := httptest.NewRecorder()

	(&App{}).serveStreamSend(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rw.Code, http.StatusBadRequest)
	}
}

func TestStreamMissingGenerationStopsStaleRuntime(t *testing.T) {
	a := &App{streams: newStreamManager(nil)}
	t.Cleanup(a.streams.close)
	req := httptest.NewRequest(http.MethodGet, streamPathPoll, nil)
	req.Header.Set(streamHeaderSession, "old-runtime")
	rw := httptest.NewRecorder()

	if got := a.sessionFor(rw, req, true, true); got != nil {
		t.Fatal("request without a page generation created a session")
	}
	if rw.Code != http.StatusGone {
		t.Fatalf("status = %d, want %d so an old poll client stops retrying", rw.Code, http.StatusGone)
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

// HEAD used to take the normal poll path: net/http suppressed the body after
// the handler had already drained it, silently losing every queued frame.
func TestStreamHeadPollDoesNotDrainFrames(t *testing.T) {
	mgr, s := newTestSession(t)
	c := newTestConn(s, 1)
	if err := s.enqueue(c, 1, frameData, []byte("still queued"), false); err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	a := &App{streams: mgr}
	req := httptest.NewRequest(http.MethodHead, streamPathPoll, nil)
	rw := httptest.NewRecorder()
	a.serveStreamPoll(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("HEAD status = %d, want 405", rw.Code)
	}
	if allow := rw.Header().Get("Allow"); allow != "GET" {
		t.Fatalf("Allow = %q, want GET", allow)
	}

	s.mu.Lock()
	frames, _ := s.drainLocked(streamMaxResponseBytes)
	s.mu.Unlock()
	if len(frames) != 1 || string(frames[0].data) != "still queued" {
		t.Fatalf("frames after HEAD = %+v, want original queued frame", frames)
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

	s1 := mgr.session("a", 1, 1, true)
	s2 := mgr.session("b", 2, 1, true)
	if s1 == nil || s2 == nil {
		t.Fatal("failed to create test sessions")
	}

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
	if got := mgr.session("a", 1, 1, false); got != nil {
		t.Fatal("destroyed window generation was recreated by a late request")
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

	first := mgr.session("page-1", 1, 1, true)
	if first == nil {
		t.Fatal("first session not created")
	}
	first.open(1, "s")
	c := first.conn(1)
	if c == nil {
		t.Fatal("connection not registered on the first session")
	}

	// Another window's session must be untouched by the reload below.
	other := mgr.session("other-window", 2, 1, true)
	other.open(1, "s")
	otherConn := other.conn(1)

	second := mgr.session("page-2", 1, 2, true)
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
	if got, done, _ := store.add("c1", 2, 3, parts[2]); done || got != nil {
		t.Fatal("completed early on the first chunk")
	}
	if _, done, _ := store.add("c1", 0, 3, parts[0]); done {
		t.Fatal("completed early on the second chunk")
	}
	got, done, _ := store.add("c1", 1, 3, parts[1])
	if !done {
		t.Fatal("did not complete on the final chunk")
	}
	if string(got) != "one two three" {
		t.Fatalf("assembled = %q", got)
	}
	// Completion remains retained until the endpoint acknowledges downstream
	// delivery, so a final-chunk retry after 429 can reproduce the same frame.
	store.remove("c1")
	if len(store.items) != 0 || store.bytes != 0 {
		t.Fatalf("store not drained: %d items, %d bytes", len(store.items), store.bytes)
	}
}

func TestStreamChunkRetryAfterBackpressureKeepsAssembledFrame(t *testing.T) {
	mgr, s := newTestSession(t)
	s.generation = 1
	c := newTestConn(s, 7)
	a := &App{streams: mgr}

	for i := 0; i < streamInQueueDepth; i++ {
		if err := c.deliver([]byte("queued")); err != nil {
			t.Fatalf("fill inbox %d: %v", i, err)
		}
	}

	sendChunk := func(index int, data string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(http.MethodPost, streamPathSend, bytes.NewReader([]byte(data)))
		req.Header.Set(streamHeaderSession, s.id)
		req.Header.Set(streamHeaderGeneration, "1")
		req.Header.Set(streamHeaderConn, "7")
		req.Header.Set(streamHeaderKind, strconv.Itoa(int(frameData)))
		req.Header.Set(streamHeaderChunkID, "large-frame")
		req.Header.Set(streamHeaderChunkIndex, strconv.Itoa(index))
		req.Header.Set(streamHeaderChunkTotal, "2")
		rw := httptest.NewRecorder()
		a.serveStreamSend(rw, req)
		return rw
	}

	if rw := sendChunk(0, "hello "); rw.Code != http.StatusNoContent {
		t.Fatalf("first chunk status = %d, want 204", rw.Code)
	}
	if rw := sendChunk(1, "world"); rw.Code != http.StatusTooManyRequests {
		t.Fatalf("completed frame against full inbox status = %d, want 429", rw.Code)
	}

	if _, err := c.Receive(); err != nil {
		t.Fatalf("free inbox slot: %v", err)
	}
	if rw := sendChunk(1, "world"); rw.Code != http.StatusNoContent {
		t.Fatalf("retried final chunk status = %d, want 204", rw.Code)
	}

	c.inMu.Lock()
	if got := string(c.in[len(c.in)-1]); got != "hello world" {
		c.inMu.Unlock()
		t.Fatalf("retried assembled frame = %q, want %q", got, "hello world")
	}
	c.inMu.Unlock()
	if len(s.chunks().items) != 0 || s.chunks().bytes != 0 {
		t.Fatal("acknowledged chunk set was not released")
	}
}

func TestStreamChunkRejectsInconsistentTotal(t *testing.T) {
	store := &streamChunkStore{items: make(map[string]*streamChunkSet)}

	if _, done, _ := store.add("c1", 0, 3, []byte("a")); done {
		t.Fatal("completed on the first of three chunks")
	}
	if _, done, err := store.add("c1", 1, 4, []byte("b")); done || err == nil {
		t.Fatal("accepted a chunk claiming a different total")
	}
	if len(store.items) != 0 {
		t.Fatal("inconsistent set was not discarded")
	}
}

// Copilot review: a handler that simply returns must tell the frontend. The
// wrapper used to defer shutdown(), which cancels the connection but queues no
// close frame, so the browser socket stayed open and onclose never fired.
func TestStreamHandlerReturnNotifiesFrontend(t *testing.T) {
	mgr, s := newTestSession(t)
	done := make(chan struct{})
	mgr.handlers["brief"] = func(c *StreamConn) { close(done) }

	s.open(3, "brief")
	<-done
	time.Sleep(50 * time.Millisecond)

	s.mu.Lock()
	frames, _ := s.drainLocked(streamMaxResponseBytes)
	s.mu.Unlock()

	var sawClose bool
	for _, f := range frames {
		if f.kind == frameClose && f.connID == 3 {
			sawClose = true
		}
	}
	if !sawClose {
		t.Fatalf("no close frame after the handler returned; frames=%+v", frames)
	}
}

// Copilot review: control frames must not be dropped by data backpressure. A
// lost open ack leaves the frontend in CONNECTING forever.
func TestStreamControlFramesBypassFullQueue(t *testing.T) {
	mgr, s := newTestSession(t)
	mgr.handlers["late"] = func(c *StreamConn) { <-c.Context().Done() }

	filler := newTestConn(s, 99)
	for i := 0; i < streamOutQueueDepth; i++ {
		if err := s.enqueue(filler, 99, frameData, []byte("x"), false); err != nil {
			t.Fatalf("fill %d: %v", i, err)
		}
	}
	if err := s.enqueue(filler, 99, frameData, []byte("x"), false); err != ErrStreamFull {
		t.Fatalf("queue not full: %v", err)
	}

	s.open(7, "late")

	s.mu.Lock()
	var sawOpen bool
	for _, f := range s.out {
		if f.kind == frameOpen && f.connID == 7 {
			sawOpen = true
		}
	}
	s.mu.Unlock()
	if !sawOpen {
		t.Fatal("open ack dropped because data frames had filled the queue")
	}
}

// A page can issue open requests faster than it polls their acknowledgements.
// Protocol frames must remain reliable, but that must not turn the control
// queue or connection lifecycle into an unbounded allocation. In particular,
// handlers that return immediately leave an open acknowledgement and a close
// notification queued for every request.
func TestStreamControlStateBoundedWithoutPolling(t *testing.T) {
	const maxConnections = 256

	t.Run("accepted connections", func(t *testing.T) {
		mgr, s := newTestSession(t)
		handled := make(chan *StreamConn)
		mgr.handlers["short"] = func(c *StreamConn) { handled <- c }

		for id := uint32(1); id <= maxConnections; id++ {
			if err := s.open(id, "short"); err != nil {
				t.Fatalf("open %d: %v", id, err)
			}
			c := <-handled
			select {
			case <-c.Context().Done():
			case <-time.After(time.Second):
				t.Fatalf("connection %d did not close after its handler returned", id)
			}
		}
		if err := s.open(maxConnections+1, "short"); !errors.Is(err, ErrStreamFull) {
			t.Fatalf("open past the bound = %v, want ErrStreamFull", err)
		}

		s.mu.Lock()
		queued := len(s.out)
		s.mu.Unlock()
		if queued > 2*maxConnections {
			t.Fatalf("queued %d control frames without a poll, want at most %d", queued, 2*maxConnections)
		}
	})

	t.Run("refused connections", func(t *testing.T) {
		_, s := newTestSession(t)
		for id := uint32(1); id <= maxConnections; id++ {
			if err := s.open(id, "missing"); err != nil {
				t.Fatalf("refuse %d: %v", id, err)
			}
		}
		if err := s.open(maxConnections+1, "missing"); !errors.Is(err, ErrStreamFull) {
			t.Fatalf("refusal past the bound = %v, want ErrStreamFull", err)
		}

		s.mu.Lock()
		queued := len(s.out)
		s.mu.Unlock()
		if queued > maxConnections {
			t.Fatalf("queued %d refusal frames without a poll, want at most %d", queued, maxConnections)
		}
	})
}

func TestStreamPendingClosesReserveConnectionCapacity(t *testing.T) {
	mgr, s := newTestSession(t)
	release := make(chan struct{})
	started := make(chan struct{}, streamMaxConnections)
	mgr.handlers["held"] = func(c *StreamConn) {
		started <- struct{}{}
		<-release
	}

	for id := uint32(1); id <= streamMaxConnections; id++ {
		if err := s.open(id, "held"); err != nil {
			t.Fatalf("open %d: %v", id, err)
		}
	}
	for range streamMaxConnections {
		<-started
	}

	// Drain every acknowledgement, then let all handlers return. Their close
	// notifications now occupy the reserved capacity while no live connection
	// remains in the map.
	s.mu.Lock()
	frames, _ := s.drainLocked(streamMaxResponseBytes)
	s.mu.Unlock()
	if len(frames) != streamMaxConnections {
		t.Fatalf("drained %d acknowledgements, want %d", len(frames), streamMaxConnections)
	}
	close(release)

	deadline := time.Now().Add(2 * time.Second)
	for {
		s.mu.Lock()
		conns, closes := len(s.conns), s.outCloses
		s.mu.Unlock()
		if conns == 0 && closes == streamMaxConnections {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("after handlers returned: %d connections, %d queued closes", conns, closes)
		}
		time.Sleep(time.Millisecond)
	}

	if err := s.open(streamMaxConnections+1, "held"); !errors.Is(err, ErrStreamFull) {
		t.Fatalf("open while close notifications are pending = %v, want ErrStreamFull", err)
	}

	s.mu.Lock()
	closes, _ := s.drainLocked(streamMaxResponseBytes)
	s.mu.Unlock()
	if len(closes) != streamMaxConnections {
		t.Fatalf("drained %d closes, want %d", len(closes), streamMaxConnections)
	}
	if err := s.open(streamMaxConnections+1, "held"); err != nil {
		t.Fatalf("open after close notifications drained: %v", err)
	}
}

func TestStreamOpenBackpressureIsRetryableAtHTTPBoundary(t *testing.T) {
	mgr, s := newTestSession(t)
	s.generation = 1
	started := make(chan struct{}, 1)
	mgr.handlers["held"] = func(c *StreamConn) {
		started <- struct{}{}
		<-c.Context().Done()
	}
	a := &App{streams: mgr}

	// Fill only the non-close control allowance with refusal notifications.
	// This must not consume the slots reserved for existing connections to
	// report close, but it must backpressure another open until a poll drains it.
	for id := uint32(1); id <= streamOutControlDepth; id++ {
		if err := s.open(id, "missing"); err != nil {
			t.Fatalf("queue refusal %d: %v", id, err)
		}
	}

	sendOpen := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, streamPathSend, nil)
		req.Header.Set(streamHeaderSession, s.id)
		req.Header.Set(streamHeaderGeneration, "1")
		req.Header.Set(streamHeaderConn, "1000")
		req.Header.Set(streamHeaderKind, strconv.Itoa(int(frameOpen)))
		req.Header.Set(streamHeaderName, "held")
		rw := httptest.NewRecorder()
		a.serveStreamSend(rw, req)
		return rw
	}

	if rw := sendOpen(); rw.Code != http.StatusTooManyRequests || rw.Header().Get("Retry-After") == "" {
		t.Fatalf("full control queue response = %d Retry-After=%q, want retryable 429", rw.Code, rw.Header().Get("Retry-After"))
	}
	if s.conn(1000) != nil {
		t.Fatal("backpressured open registered a connection")
	}

	s.mu.Lock()
	_, _ = s.drainLocked(streamMaxResponseBytes)
	s.mu.Unlock()
	if rw := sendOpen(); rw.Code != http.StatusNoContent {
		t.Fatalf("retried open status = %d, want 204", rw.Code)
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("retried open did not start its handler")
	}
}

// Send queues the caller's slice rather than copying it. The ownership rule is
// documented on Send; this pins the behaviour so a copy is not reintroduced for
// tidiness, since at these rates a memcpy per frame is the dominant cost.
func TestStreamSendDoesNotCopyPayload(t *testing.T) {
	_, s := newTestSession(t)
	c := newTestConn(s, 1)

	buf := []byte("original")
	if err := c.Send(buf); err != nil {
		t.Fatalf("Send: %v", err)
	}

	s.mu.Lock()
	frames, _ := s.drainLocked(streamMaxResponseBytes)
	s.mu.Unlock()
	if len(frames) != 1 {
		t.Fatalf("got %d frames", len(frames))
	}
	if &frames[0].data[0] != &buf[0] {
		t.Fatal("frame does not alias the caller's slice; a copy was reintroduced")
	}
}

// Copilot review: the inbound queue was unbounded, so a frontend could grow
// host memory without limit against a handler that is slow to Receive.
func TestStreamInboundQueueBounded(t *testing.T) {
	_, s := newTestSession(t)
	c := newTestConn(s, 1)

	for i := 0; i < streamInQueueDepth; i++ {
		if err := c.deliver([]byte("x")); err != nil {
			t.Fatalf("fill %d: %v", i, err)
		}
	}

	// Signalled, not waited out: blocking here would hold a request slot and
	// starve the window's poll.
	if err := c.deliver([]byte("y")); err != ErrStreamFull {
		t.Fatalf("deliver past the bound = %v, want ErrStreamFull", err)
	}

	if _, err := c.Receive(); err != nil {
		t.Fatalf("Receive: %v", err)
	}

	if err := c.deliver([]byte("y")); err != nil {
		t.Fatalf("deliver after Receive freed space: %v", err)
	}
}

func TestStreamInboundRejectsOnceClosed(t *testing.T) {
	_, s := newTestSession(t)
	c := newTestConn(s, 1)

	c.shutdown()
	if err := c.deliver([]byte("x")); err != ErrStreamClosed {
		t.Fatalf("deliver on a closed connection = %v, want ErrStreamClosed", err)
	}
}

func TestStreamBlockingInboundDeliveryWakesOnReceive(t *testing.T) {
	_, s := newTestSession(t)
	c := newTestConn(s, 1)

	for i := 0; i < streamInQueueDepth; i++ {
		if err := c.deliver([]byte{byte(i)}); err != nil {
			t.Fatalf("fill %d: %v", i, err)
		}
	}

	delivered := make(chan error, 1)
	go func() { delivered <- c.deliverBlocking([]byte("after-space")) }()
	select {
	case err := <-delivered:
		t.Fatalf("blocking delivery returned before Receive freed space: %v", err)
	case <-time.After(20 * time.Millisecond):
	}

	if _, err := c.Receive(); err != nil {
		t.Fatalf("Receive: %v", err)
	}
	select {
	case err := <-delivered:
		if err != nil {
			t.Fatalf("blocking delivery after Receive: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("blocking delivery did not wake after Receive freed space")
	}
}

func TestStreamBlockingInboundDeliveryWakesOnClose(t *testing.T) {
	_, s := newTestSession(t)
	c := newTestConn(s, 1)

	for i := 0; i < streamInQueueDepth; i++ {
		if err := c.deliver([]byte("full")); err != nil {
			t.Fatalf("fill %d: %v", i, err)
		}
	}

	delivered := make(chan error, 1)
	go func() { delivered <- c.deliverBlocking([]byte("blocked")) }()
	select {
	case err := <-delivered:
		t.Fatalf("blocking delivery returned before close: %v", err)
	case <-time.After(20 * time.Millisecond):
	}

	c.shutdown()
	select {
	case err := <-delivered:
		if !errors.Is(err, ErrStreamClosed) {
			t.Fatalf("blocking delivery after close = %v, want ErrStreamClosed", err)
		}
	case <-time.After(time.Second):
		t.Fatal("blocking delivery did not wake after close")
	}
}

func TestStreamSustainedBlockingInboundRemainsBoundedAndOrdered(t *testing.T) {
	_, s := newTestSession(t)
	c := newTestConn(s, 1)
	const frames = streamInQueueDepth * 4

	producer := make(chan error, 1)
	go func() {
		for i := 0; i < frames; i++ {
			payload := binary.BigEndian.AppendUint32(nil, uint32(i))
			if err := c.deliverBlocking(payload); err != nil {
				producer <- err
				return
			}
			c.inMu.Lock()
			if len(c.in) > streamInQueueDepth || c.inBytes > streamInQueueBytes {
				c.inMu.Unlock()
				producer <- fmt.Errorf("inbox exceeded bound: %d frames, %d bytes", len(c.in), c.inBytes)
				return
			}
			c.inMu.Unlock()
		}
		producer <- nil
	}()

	for want := 0; want < frames; want++ {
		payload, err := c.Receive()
		if err != nil {
			t.Fatalf("Receive %d: %v", want, err)
		}
		if got := int(binary.BigEndian.Uint32(payload)); got != want {
			t.Fatalf("frame %d carried sequence %d", want, got)
		}
		if want%32 == 0 {
			time.Sleep(time.Millisecond)
		}
	}
	if err := <-producer; err != nil {
		t.Fatal(err)
	}
}

// Copilot review: a rejected chunk set answered 204, so the client believed the
// frame had been sent, and an inconsistent total leaked its bytes forever.
func TestStreamChunkRejectionIsReportedAndAccounted(t *testing.T) {
	store := &streamChunkStore{items: make(map[string]*streamChunkSet)}

	if _, done, err := store.add("c1", 0, 3, []byte("abc")); done || err != nil {
		t.Fatalf("first chunk: done=%v err=%v", done, err)
	}
	if store.bytes != 3 {
		t.Fatalf("bytes = %d, want 3", store.bytes)
	}

	_, done, err := store.add("c1", 1, 4, []byte("d"))
	if done || err == nil {
		t.Fatalf("conflicting total: done=%v err=%v, want an error", done, err)
	}
	if store.bytes != 0 {
		t.Fatalf("bytes = %d after discarding the set, want 0", store.bytes)
	}
	if len(store.items) != 0 {
		t.Fatal("inconsistent set was not discarded")
	}
}

func TestStreamChunkRejectionsReturnBadRequestAndReleaseAccounting(t *testing.T) {
	newEndpoint := func(t *testing.T) (*App, *streamSession) {
		t.Helper()
		mgr, s := newTestSession(t)
		s.generation = 1
		newTestConn(s, 7)
		return &App{streams: mgr}, s
	}
	sendChunk := func(t *testing.T, a *App, s *streamSession, id string, index, total int, data string) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(http.MethodPost, streamPathSend, bytes.NewReader([]byte(data)))
		req.Header.Set(streamHeaderSession, s.id)
		req.Header.Set(streamHeaderGeneration, "1")
		req.Header.Set(streamHeaderConn, "7")
		req.Header.Set(streamHeaderKind, strconv.Itoa(int(frameData)))
		req.Header.Set(streamHeaderChunkID, id)
		req.Header.Set(streamHeaderChunkIndex, strconv.Itoa(index))
		req.Header.Set(streamHeaderChunkTotal, strconv.Itoa(total))
		rw := httptest.NewRecorder()
		a.serveStreamSend(rw, req)
		return rw
	}

	t.Run("duplicate", func(t *testing.T) {
		a, s := newEndpoint(t)
		if got := sendChunk(t, a, s, "duplicate", 0, 2, "a").Code; got != http.StatusNoContent {
			t.Fatalf("first chunk status = %d, want 204", got)
		}
		if got := sendChunk(t, a, s, "duplicate", 0, 2, "a").Code; got != http.StatusBadRequest {
			t.Fatalf("duplicate status = %d, want 400", got)
		}
		if store := s.chunks(); store.bytes != 0 || len(store.items) != 0 {
			t.Fatalf("duplicate retained %d bytes in %d sets", store.bytes, len(store.items))
		}
	})

	t.Run("conflicting total", func(t *testing.T) {
		a, s := newEndpoint(t)
		if got := sendChunk(t, a, s, "conflict", 0, 2, "a").Code; got != http.StatusNoContent {
			t.Fatalf("first chunk status = %d, want 204", got)
		}
		if got := sendChunk(t, a, s, "conflict", 1, 3, "b").Code; got != http.StatusBadRequest {
			t.Fatalf("conflict status = %d, want 400", got)
		}
		if store := s.chunks(); store.bytes != 0 || len(store.items) != 0 {
			t.Fatalf("conflict retained %d bytes in %d sets", store.bytes, len(store.items))
		}
	})

	t.Run("aggregate overflow", func(t *testing.T) {
		a, s := newEndpoint(t)
		store := s.chunks()
		store.items["overflow"] = &streamChunkSet{
			parts:   map[int][]byte{0: nil},
			total:   2,
			size:    streamMaxSendBytes,
			created: time.Now(),
		}
		store.bytes = streamMaxSendBytes
		if got := sendChunk(t, a, s, "overflow", 1, 2, "x").Code; got != http.StatusBadRequest {
			t.Fatalf("overflow status = %d, want 400", got)
		}
		if store.bytes != 0 || len(store.items) != 0 {
			t.Fatalf("overflow retained %d bytes in %d sets", store.bytes, len(store.items))
		}
	})

	t.Run("invalid completed retry", func(t *testing.T) {
		a, s := newEndpoint(t)
		c := s.conn(7)
		for i := 0; i < streamInQueueDepth; i++ {
			if err := c.deliver([]byte("queued")); err != nil {
				t.Fatalf("fill inbox %d: %v", i, err)
			}
		}
		if got := sendChunk(t, a, s, "completed", 0, 2, "hello ").Code; got != http.StatusNoContent {
			t.Fatalf("first chunk status = %d, want 204", got)
		}
		if got := sendChunk(t, a, s, "completed", 1, 2, "world").Code; got != http.StatusTooManyRequests {
			t.Fatalf("completed chunk status = %d, want 429", got)
		}
		if got := sendChunk(t, a, s, "completed", 1, 2, "changed").Code; got != http.StatusBadRequest {
			t.Fatalf("changed retry status = %d, want 400", got)
		}
		if store := s.chunks(); store.bytes != 0 || len(store.items) != 0 {
			t.Fatalf("invalid retry retained %d bytes in %d sets", store.bytes, len(store.items))
		}
	})
}

func TestStreamChunkStoreReapsExpiredIncompleteSets(t *testing.T) {
	store := &streamChunkStore{
		items: map[string]*streamChunkSet{
			"expired": {
				parts:   map[int][]byte{0: []byte("old")},
				total:   2,
				size:    3,
				created: time.Now().Add(-streamHoldTimeout - time.Second),
			},
		},
		bytes: 3,
	}

	if _, done, err := store.add("fresh", 0, 2, []byte("new")); done || err != nil {
		t.Fatalf("fresh chunk: done=%v err=%v", done, err)
	}
	if _, ok := store.items["expired"]; ok {
		t.Fatal("expired chunk set was not removed")
	}
	if store.bytes != 3 {
		t.Fatalf("bytes after reap = %d, want only the 3 fresh bytes", store.bytes)
	}
}

// Codex review: the client starts its open POST and its first poll
// concurrently, so the POST can create the session first. Supersede used to run
// only when the poll created the session, so in that ordering the previous
// page's session survived the whole live-connection grace.
func TestStreamPollSupersedesWhenOpenWonTheRace(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	mgr.handlers["s"] = func(c *StreamConn) { <-c.Context().Done() }

	old := mgr.session("page-1", 1, 1, true)
	old.open(1, "s")
	victim := old.conn(1)
	if victim == nil {
		t.Fatal("first session has no connection")
	}

	// The new page's open POST lands first: creates the session, no supersede.
	fresh := mgr.session("page-2", 1, 2, false)
	if fresh == nil {
		t.Fatal("send did not create the session")
	}
	if victim.ctx.Err() != nil {
		t.Fatal("a send superseded the previous session; only a poll should")
	}

	// Its poll follows, and must retire the previous page even though the
	// session already exists.
	if got := mgr.session("page-2", 1, 2, true); got != fresh {
		t.Fatal("poll did not resolve to the session the send created")
	}
	if victim.ctx.Err() == nil {
		t.Fatal("previous page still live after the new page polled")
	}
	if mgr.existingSession("page-1") != nil {
		t.Fatal("previous session still registered")
	}
	if mgr.existingSession("page-2") == nil {
		t.Fatal("poll superseded the session doing the polling")
	}
}

// A poll issued by the previous page can reach the manager after the new
// page's open POST. It must not retire the newer session merely because it was
// processed later.
func TestStreamLateOldPollDoesNotSupersedeFreshOpen(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	mgr.handlers["s"] = func(c *StreamConn) { <-c.Context().Done() }

	old := mgr.session("page-1", 1, 1, true)
	old.open(1, "s")
	oldConn := old.conn(1)
	if oldConn == nil {
		t.Fatal("old session has no connection")
	}

	// The replacement page's open POST wins its race with its first poll.
	fresh := mgr.session("page-2", 1, 2, false)
	fresh.open(1, "s")
	freshConn := fresh.conn(1)
	if freshConn == nil {
		t.Fatal("fresh session has no connection")
	}

	// A delayed poll from page 1 must not close page 2.
	if got := mgr.session("page-1", 1, 1, true); got != old {
		t.Fatal("late old poll did not resolve to the old session")
	}
	if freshConn.ctx.Err() != nil {
		t.Fatal("late old poll closed the fresh page connection")
	}

	// The first poll from page 2 now retires page 1.
	if got := mgr.session("page-2", 1, 2, true); got != fresh {
		t.Fatal("fresh poll did not resolve to the session created by its open")
	}
	if oldConn.ctx.Err() == nil {
		t.Fatal("fresh poll did not close the old page connection")
	}
	if freshConn.ctx.Err() != nil {
		t.Fatal("fresh poll closed its own connection")
	}
}

// Page order must not be inferred from server arrival order. An old page can
// start a request just before navigation and have that first request dispatched
// only after the replacement page is already connected.
func TestStreamPreviouslyUnseenOldPageCannotBecomeNewerByArrivingLate(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	mgr.handlers["s"] = func(c *StreamConn) { <-c.Context().Done() }

	fresh := mgr.session("page-2", 1, 2, false)
	fresh.open(1, "s")
	freshConn := fresh.conn(1)

	// The old page's first request reaches Go second, but carries its older page
	// generation. Its handler may start briefly; the next fresh poll retires it.
	old := mgr.session("page-1", 1, 1, false)
	old.open(1, "s")
	oldConn := old.conn(1)
	if freshConn == nil || oldConn == nil {
		t.Fatal("failed to establish test connections")
	}

	if got := mgr.session("page-2", 1, 2, true); got != fresh {
		t.Fatal("fresh poll did not retain the current page session")
	}
	if freshConn.ctx.Err() != nil {
		t.Fatal("late old request closed the current page connection")
	}
	if oldConn.ctx.Err() == nil {
		t.Fatal("current page poll did not retire the late old session")
	}
	if got := mgr.session("page-1", 1, 1, true); got != nil {
		t.Fatal("retired late page was recreated by its delayed poll")
	}
}

// Once the current page has polled, even a previously unseen request from an
// older page generation must be rejected immediately. There may be no older
// session for the supersede sweep to discover when the current poll arrives.
func TestStreamCurrentPollRejectsPreviouslyUnseenOlderPage(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)

	fresh := mgr.session("page-2", 1, 2, false)
	if fresh == nil {
		t.Fatal("failed to create current page session")
	}
	if got := mgr.session("page-2", 1, 2, true); got != fresh {
		t.Fatal("current page poll did not retain its session")
	}

	if got := mgr.session("page-1", 1, 1, false); got != nil {
		t.Fatal("late open created a previously unseen older page session")
	}
	if got := mgr.session("page-1", 1, 1, true); got != nil {
		t.Fatal("late poll created a previously unseen older page session")
	}
	if got := mgr.existingSession("page-2"); got != fresh {
		t.Fatal("late older request disturbed the current page session")
	}
}

// Once a session is superseded, a request that was already in flight for its
// id must receive Gone rather than recreating the old page as a newer session.
func TestStreamRetiredSessionCannotBeRecreated(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)

	old := mgr.session("page-1", 1, 1, true)
	fresh := mgr.session("page-2", 1, 2, true)
	if old == nil || fresh == nil {
		t.Fatal("failed to create test sessions")
	}

	if got := mgr.session("page-1", 1, 1, false); got != nil {
		t.Fatal("late open recreated a retired session")
	}
	if got := mgr.session("page-1", 1, 1, true); got != nil {
		t.Fatal("late poll recreated a retired session")
	}
	if got := mgr.existingSession("page-2"); got != fresh {
		t.Fatal("late retired request disturbed the active session")
	}
}

// Expiring an idle session cleans up transport state, not the page that owned
// it. A page stops polling when it has no connections, so it must be able to
// open another stream after the janitor has removed that idle session.
func TestStreamReapedIdleSessionCanBeRecreated(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)

	now := time.Now()
	old := mgr.session("page-1", 1, 1, true)
	if old == nil {
		t.Fatal("failed to create page session")
	}
	old.mu.Lock()
	old.lastSeen = now.Add(-streamSessionTTL - time.Second)
	old.mu.Unlock()

	mgr.reapStale(now)
	if got := mgr.existingSession("page-1"); got != nil {
		t.Fatal("idle session was not reaped")
	}
	old.mu.Lock()
	closed := old.closed
	old.mu.Unlock()
	if !closed {
		t.Fatal("reaped session was not closed")
	}

	reopened := mgr.session("page-1", 1, 1, false)
	if reopened == nil {
		t.Fatal("live page could not recreate its session after idle reap")
	}
	if reopened == old {
		t.Fatal("idle reap reused the closed session")
	}
	if got := mgr.session("page-1", 1, 1, true); got != reopened {
		t.Fatal("poll did not retain the live page's recreated session")
	}
}

// The poll cleanup path observes an already-closed session. Removing that
// transport state must likewise leave the still-loaded page generation usable.
func TestStreamRemovedClosedSessionCanBeRecreated(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)

	old := mgr.session("page-1", 1, 1, true)
	if old == nil {
		t.Fatal("failed to create page session")
	}
	old.close()
	mgr.removeSession(old.id)

	if got := mgr.session("page-1", 1, 1, false); got == nil {
		t.Fatal("live page could not recreate its session after closed-session cleanup")
	}
}

// Reaping the current page's idle transport must not weaken the watermark
// established when that page superseded its predecessor. The current page may
// reopen; a delayed request from the previous page must remain gone.
func TestStreamReapPreservesSupersededGenerationWatermark(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)

	old := mgr.session("page-1", 1, 1, true)
	current := mgr.session("page-2", 1, 2, true)
	if old == nil || current == nil {
		t.Fatal("failed to create page generations")
	}

	now := time.Now()
	current.mu.Lock()
	current.lastSeen = now.Add(-streamSessionTTL - time.Second)
	current.mu.Unlock()
	mgr.reapStale(now)

	if got := mgr.session("page-1", 1, 1, false); got != nil {
		t.Fatal("idle reap allowed the superseded page generation to return")
	}
	if got := mgr.session("page-2", 1, 2, false); got == nil {
		t.Fatal("idle reap prevented the current page generation from reopening")
	}
}

func TestStreamForeignSessionIDCannotSupersedeRequestingWindow(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	mgr.handlers["s"] = func(c *StreamConn) { <-c.Context().Done() }

	// Create the legitimate session first so the foreign session has the newer
	// generation that triggered the old, destructive sweep.
	legitimate := mgr.session("window-2", 2, 1, true)
	legitimate.open(1, "s")
	conn := legitimate.conn(1)
	foreign := mgr.session("window-1", 1, 2, true)
	if conn == nil || foreign == nil {
		t.Fatal("failed to establish test sessions")
	}

	// The request claims window 2 while carrying window 1's id. sessionFor will
	// reject it as forbidden; manager lookup must have no lifecycle side effect.
	if got := mgr.session("window-1", 2, 2, true); got != foreign {
		t.Fatal("foreign lookup did not return the existing session for ownership validation")
	}
	if conn.ctx.Err() != nil {
		t.Fatal("foreign session id closed the requesting window's legitimate connection")
	}
	if mgr.existingSession("window-2") != legitimate {
		t.Fatal("foreign session id removed the requesting window's session")
	}
}

func TestStreamGenerationMismatchCannotSupersedeOtherSessions(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)

	current := mgr.session("current", 1, 3, false)
	other := mgr.session("other", 1, 2, false)
	if current == nil || other == nil {
		t.Fatal("failed to create test sessions")
	}

	// This request claims the current id but carries an older generation. The
	// endpoint will reject it; the lookup must be side-effect free first.
	if got := mgr.session("current", 1, 1, true); got != current {
		t.Fatal("generation mismatch did not resolve for endpoint validation")
	}
	if got := mgr.existingSession("other"); got != other {
		t.Fatal("generation-mismatched poll retired another valid session")
	}
}
