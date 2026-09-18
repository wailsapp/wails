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
	"strings"
	"sync"
	"sync/atomic"
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
	c := &StreamConn{id: id, name: "test", sink: s, ctx: ctx, cancel: cancel, manager: s.mgr}
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

func TestStreamBatchRejectsLengthBeyondRemainingBody(t *testing.T) {
	body := binary.BigEndian.AppendUint32(nil, 1)
	body = binary.BigEndian.AppendUint32(body, ^uint32(0))

	if _, err := decodeStreamBatch(body); err == nil {
		t.Fatal("batch with an overflowing declared length was accepted")
	}
}

func TestStreamRejectsOversizedOutboundFrame(t *testing.T) {
	_, s := newTestSession(t)
	c := newTestConn(s, 1)
	payload := make([]byte, streamMaxFrameBytes+1)

	if err := c.TrySend(payload); !errors.Is(err, ErrStreamTooLarge) {
		t.Fatalf("oversized TrySend = %v, want ErrStreamTooLarge", err)
	}
	if err := c.Send(payload); !errors.Is(err, ErrStreamTooLarge) {
		t.Fatalf("oversized Send = %v, want ErrStreamTooLarge", err)
	}
	if err := c.TrySend(payload[:streamMaxFrameBytes]); err != nil {
		t.Fatalf("maximum-size TrySend = %v", err)
	}

	s.mu.Lock()
	frames, _ := s.drainLocked(streamMaxResponseBytes)
	s.mu.Unlock()
	if len(frames) != 1 || len(frames[0].data) != streamMaxFrameBytes {
		t.Fatalf("queued %d frames with first length %d, want one %d-byte frame", len(frames), len(frames[0].data), streamMaxFrameBytes)
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

func TestStreamSendRejectsInvalidFramesBeforeRetainingBodies(t *testing.T) {
	t.Run("unsupported kind before session admission", func(t *testing.T) {
		mgr := newStreamManager(nil)
		t.Cleanup(mgr.close)
		a := &App{streams: mgr}
		req := httptest.NewRequest(http.MethodPost, streamPathSend, strings.NewReader("ignored"))
		req.Header.Set(streamHeaderSession, "invalid-kind")
		req.Header.Set(streamHeaderGeneration, "1")
		req.Header.Set(streamHeaderConn, "7")
		req.Header.Set(streamHeaderKind, strconv.Itoa(int(frameError)))
		rw := httptest.NewRecorder()

		a.serveStreamSend(rw, req)

		if rw.Code != http.StatusBadRequest {
			t.Fatalf("unsupported kind status = %d, want 400", rw.Code)
		}
		mgr.mu.Lock()
		defer mgr.mu.Unlock()
		if len(mgr.sessions) != 0 {
			t.Fatalf("unsupported kind admitted %d sessions", len(mgr.sessions))
		}
	})

	t.Run("unknown data connection before chunk storage", func(t *testing.T) {
		mgr, s := newTestSession(t)
		s.generation = 1
		a := &App{streams: mgr}
		req := httptest.NewRequest(http.MethodPost, streamPathSend, strings.NewReader("partial"))
		req.Header.Set(streamHeaderSession, s.id)
		req.Header.Set(streamHeaderGeneration, "1")
		req.Header.Set(streamHeaderConn, "404")
		req.Header.Set(streamHeaderKind, strconv.Itoa(int(frameData)))
		req.Header.Set(streamHeaderChunkID, "orphan")
		req.Header.Set(streamHeaderChunkIndex, "0")
		req.Header.Set(streamHeaderChunkTotal, "2")
		rw := httptest.NewRecorder()

		a.serveStreamSend(rw, req)

		if rw.Code != http.StatusGone {
			t.Fatalf("unknown connection status = %d, want 410", rw.Code)
		}
		s.mu.Lock()
		chunkStore := s.chunkStore
		s.mu.Unlock()
		if chunkStore != nil || mgr.chunkBytes != 0 || mgr.chunkParts != 0 {
			t.Fatalf("unknown connection retained chunk state: store=%v bytes=%d parts=%d", chunkStore != nil, mgr.chunkBytes, mgr.chunkParts)
		}
	})

	for _, tc := range []struct {
		name   string
		kind   uint8
		body   string
		chunks bool
	}{
		{name: "open body", kind: frameOpen, body: "unexpected"},
		{name: "open chunks", kind: frameOpen, chunks: true},
		{name: "close body", kind: frameClose, body: "unexpected"},
		{name: "close chunks", kind: frameClose, chunks: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mgr := newStreamManager(nil)
			t.Cleanup(mgr.close)
			a := &App{streams: mgr}
			req := httptest.NewRequest(http.MethodPost, streamPathSend, strings.NewReader(tc.body))
			req.Header.Set(streamHeaderSession, "control")
			req.Header.Set(streamHeaderGeneration, "1")
			req.Header.Set(streamHeaderConn, "7")
			req.Header.Set(streamHeaderKind, strconv.Itoa(int(tc.kind)))
			if tc.kind == frameOpen {
				req.Header.Set(streamHeaderName, "test")
			}
			if tc.chunks {
				req.Header.Set(streamHeaderChunkID, "control-chunk")
				req.Header.Set(streamHeaderChunkIndex, "0")
				req.Header.Set(streamHeaderChunkTotal, "2")
			}
			rw := httptest.NewRecorder()

			a.serveStreamSend(rw, req)

			if rw.Code != http.StatusBadRequest {
				t.Fatalf("control frame status = %d, want 400", rw.Code)
			}
			mgr.mu.Lock()
			defer mgr.mu.Unlock()
			if len(mgr.sessions) != 0 || mgr.chunkBytes != 0 || mgr.chunkParts != 0 {
				t.Fatalf("control frame retained sessions=%d bytes=%d parts=%d", len(mgr.sessions), mgr.chunkBytes, mgr.chunkParts)
			}
		})
	}
}

func TestStreamOpenRejectsOversizedNameBeforeSessionAdmission(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	a := &App{streams: mgr}
	req := httptest.NewRequest(http.MethodPost, streamPathSend, nil)
	req.Header.Set(streamHeaderSession, "unadmitted")
	req.Header.Set(streamHeaderGeneration, "1")
	req.Header.Set(streamHeaderConn, "1")
	req.Header.Set(streamHeaderKind, strconv.Itoa(int(frameOpen)))
	req.Header.Set(streamHeaderName, strings.Repeat("n", streamMaxNameLen+1))
	rw := httptest.NewRecorder()

	a.serveStreamSend(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("oversized name status = %d, want 400", rw.Code)
	}
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if len(mgr.sessions) != 0 {
		t.Fatalf("oversized name admitted %d sessions", len(mgr.sessions))
	}
}

func TestHandleStreamNameLimit(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	a := &App{streams: mgr}
	handler := func(c *StreamConn) { <-c.Context().Done() }
	maximum := strings.Repeat("m", streamMaxNameLen)
	oversized := maximum + "x"

	a.HandleStream(maximum, handler)
	a.HandleStream(oversized, handler)

	if _, ok := mgr.handler(maximum); !ok {
		t.Fatal("maximum-length stream name was not registered")
	}
	if _, ok := mgr.handler(oversized); ok {
		t.Fatal("oversized stream name was registered")
	}

	s := newStreamSession("maximum-name", 1, mgr)
	mgr.sessions[s.id] = s
	if err := s.open(1, maximum); err != nil {
		t.Fatalf("maximum-length stream name was not accepted: %v", err)
	}
	if c := s.conn(1); c == nil || c.Name() != maximum {
		t.Fatal("maximum-length stream name did not reach its connection")
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

func TestStreamManagerBoundsSessionsCreatedWithoutPoll(t *testing.T) {
	t.Run("per window", func(t *testing.T) {
		mgr := newStreamManager(nil)
		t.Cleanup(mgr.close)
		for generation := uint64(1); generation <= streamMaxSessionsPerWindow; generation++ {
			id := fmt.Sprintf("page-%d", generation)
			if got := mgr.session(id, 7, generation, false); got == nil {
				t.Fatalf("session %d refused before the per-window bound", generation)
			}
		}
		if got := mgr.session("one-too-many", 7, streamMaxSessionsPerWindow+1, false); got != nil {
			t.Fatal("open POST created a session past the per-window bound")
		}
	})

	t.Run("global windowless fallback", func(t *testing.T) {
		mgr := newStreamManager(nil)
		t.Cleanup(mgr.close)
		for generation := uint64(1); generation <= streamMaxSessions; generation++ {
			id := fmt.Sprintf("fallback-%d", generation)
			if got := mgr.session(id, 0, generation, false); got == nil {
				t.Fatalf("session %d refused before the global bound", generation)
			}
		}
		if got := mgr.session("one-too-many", 0, streamMaxSessions+1, false); got != nil {
			t.Fatal("windowless request created a session past the global bound")
		}
	})
}

func TestStreamWindowlessSessionsDoNotSupersedeEachOther(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	mgr.handlers["s"] = func(c *StreamConn) { <-c.Context().Done() }

	// A missing platform window id means the requests may belong to unrelated
	// browser clients. Their page generations come from independent browsing
	// contexts and therefore cannot be ordered against one another.
	first := mgr.session("client-a", 0, 1, true)
	if first == nil {
		t.Fatal("first windowless session was not created")
	}
	if err := first.open(1, "s"); err != nil {
		t.Fatalf("open first windowless stream: %v", err)
	}
	firstConn := first.conn(1)

	second := mgr.session("client-b", 0, 2, true)
	if second == nil {
		t.Fatal("second windowless session was not created")
	}
	if mgr.existingSession("client-a") != first || firstConn.ctx.Err() != nil {
		t.Fatal("one windowless client's poll superseded another client")
	}

	// A third independent context may legitimately start at the same generation
	// as the first. No shared retired watermark may reject it.
	third := mgr.session("client-c", 0, 1, true)
	if third == nil {
		t.Fatal("shared windowless generation watermark rejected an independent client")
	}
}

func TestStreamCurrentPollCanReplaceSessionsAtAdmissionBound(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	for generation := uint64(1); generation <= streamMaxSessionsPerWindow; generation++ {
		id := fmt.Sprintf("stale-%d", generation)
		if got := mgr.session(id, 7, generation, false); got == nil {
			t.Fatalf("create stale session %d", generation)
		}
	}
	other := mgr.session("other-window", 8, 1, false)
	if other == nil {
		t.Fatal("create other-window session")
	}

	replacement := mgr.session("current", 7, streamMaxSessionsPerWindow+1, true)
	if replacement == nil {
		t.Fatal("current poll could not replace stale sessions at the bound")
	}
	mgr.mu.Lock()
	windowCount := 0
	for _, session := range mgr.sessions {
		if session.windowID == 7 {
			windowCount++
		}
	}
	otherStillPresent := mgr.sessions["other-window"] == other
	mgr.mu.Unlock()
	if windowCount != 1 {
		t.Fatalf("replacement left %d sessions in its window, want 1", windowCount)
	}
	if !otherStillPresent {
		t.Fatal("replacement removed another window's session")
	}
}

func TestStreamSessionAdmissionBoundAtHTTPBoundary(t *testing.T) {
	request := func(method, id string, windowID uint, generation uint64) *http.Request {
		req := httptest.NewRequest(method, streamPathSend, nil)
		req.Header.Set(streamHeaderSession, id)
		req.Header.Set(streamHeaderGeneration, strconv.FormatUint(generation, 10))
		if windowID != 0 {
			req.Header.Set(webViewRequestHeaderWindowId, strconv.FormatUint(uint64(windowID), 10))
		}
		return req
	}

	t.Run("open retries until current poll supersedes stale pages", func(t *testing.T) {
		mgr := newStreamManager(nil)
		t.Cleanup(mgr.close)
		a := &App{streams: mgr}
		for generation := uint64(1); generation <= streamMaxSessionsPerWindow; generation++ {
			if got := mgr.session(fmt.Sprintf("stale-%d", generation), 7, generation, false); got == nil {
				t.Fatalf("create stale session %d", generation)
			}
		}
		other := mgr.session("other-window", 8, 1, false)
		if other == nil {
			t.Fatal("create other-window session")
		}

		generation := uint64(streamMaxSessionsPerWindow + 1)
		rw := httptest.NewRecorder()
		if got := a.sessionFor(rw, request(http.MethodPost, "current", 7, generation), true, false); got != nil {
			t.Fatal("open request bypassed the per-window admission bound")
		}
		if rw.Code != http.StatusTooManyRequests || rw.Header().Get("Retry-After") == "" {
			t.Fatalf("full-session response = %d Retry-After=%q, want retryable 429", rw.Code, rw.Header().Get("Retry-After"))
		}

		pollResponse := httptest.NewRecorder()
		current := a.sessionFor(pollResponse, request(http.MethodGet, "current", 7, generation), true, true)
		if current == nil {
			t.Fatalf("current poll could not replace stale pages: status %d body %q", pollResponse.Code, pollResponse.Body.String())
		}

		retryResponse := httptest.NewRecorder()
		if got := a.sessionFor(retryResponse, request(http.MethodPost, "current", 7, generation), true, false); got != current {
			t.Fatalf("retried open did not resolve replacement session: status %d", retryResponse.Code)
		}
		if mgr.existingSession("other-window") != other {
			t.Fatal("replacement disturbed another window")
		}
	})

	t.Run("windowless sessions obey the global bound", func(t *testing.T) {
		mgr := newStreamManager(nil)
		t.Cleanup(mgr.close)
		a := &App{streams: mgr}
		for generation := uint64(1); generation <= streamMaxSessions; generation++ {
			if got := mgr.session(fmt.Sprintf("fallback-%d", generation), 0, generation, false); got == nil {
				t.Fatalf("create windowless session %d", generation)
			}
		}

		rw := httptest.NewRecorder()
		if got := a.sessionFor(rw, request(http.MethodPost, "one-too-many", 0, streamMaxSessions+1), true, false); got != nil {
			t.Fatal("windowless request bypassed the global admission bound")
		}
		if rw.Code != http.StatusTooManyRequests || rw.Header().Get("Retry-After") == "" {
			t.Fatalf("global-bound response = %d Retry-After=%q, want retryable 429", rw.Code, rw.Header().Get("Retry-After"))
		}
	})
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
	if store.bytes != 3 || store.parts != 1 {
		t.Fatalf("store = %d bytes, %d parts; want 3 bytes, 1 part", store.bytes, store.parts)
	}

	_, done, err := store.add("c1", 1, 4, []byte("d"))
	if done || err == nil {
		t.Fatalf("conflicting total: done=%v err=%v, want an error", done, err)
	}
	if store.bytes != 0 || store.parts != 0 {
		t.Fatalf("store after discarding the set = %d bytes, %d parts; want zero", store.bytes, store.parts)
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
		if store := s.chunks(); store.bytes != 0 || store.parts != 0 || len(store.items) != 0 {
			t.Fatalf("duplicate retained %d bytes and %d parts in %d sets", store.bytes, store.parts, len(store.items))
		}
		if s.mgr.chunkBytes != 0 || s.mgr.chunkParts != 0 {
			t.Fatalf("duplicate retained %d shared bytes and %d parts", s.mgr.chunkBytes, s.mgr.chunkParts)
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
		if store := s.chunks(); store.bytes != 0 || store.parts != 0 || len(store.items) != 0 {
			t.Fatalf("conflict retained %d bytes and %d parts in %d sets", store.bytes, store.parts, len(store.items))
		}
		if s.mgr.chunkBytes != 0 || s.mgr.chunkParts != 0 {
			t.Fatalf("conflict retained %d shared bytes and %d parts", s.mgr.chunkBytes, s.mgr.chunkParts)
		}
	})

	t.Run("aggregate overflow", func(t *testing.T) {
		a, s := newEndpoint(t)
		store := s.chunks()
		store.items["overflow"] = &streamChunkSet{
			parts:     map[int][]byte{0: nil},
			total:     2,
			size:      streamMaxFrameBytes,
			partCount: 1,
			created:   time.Now(),
		}
		store.bytes = streamMaxFrameBytes
		store.parts = 1
		if got := sendChunk(t, a, s, "overflow", 1, 2, "x").Code; got != http.StatusBadRequest {
			t.Fatalf("overflow status = %d, want 400", got)
		}
		if store.bytes != 0 || store.parts != 0 || len(store.items) != 0 {
			t.Fatalf("overflow retained %d bytes and %d parts in %d sets", store.bytes, store.parts, len(store.items))
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
		if store := s.chunks(); store.bytes != 0 || store.parts != 0 || len(store.items) != 0 {
			t.Fatalf("invalid retry retained %d bytes and %d parts in %d sets", store.bytes, store.parts, len(store.items))
		}
		if s.mgr.chunkBytes != 0 || s.mgr.chunkParts != 0 {
			t.Fatalf("invalid retry retained %d shared bytes and %d parts", s.mgr.chunkBytes, s.mgr.chunkParts)
		}
	})
}

func TestStreamChunkStoreReapsExpiredIncompleteSets(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	store := &streamChunkStore{
		items: map[string]*streamChunkSet{
			"expired": {
				parts:     map[int][]byte{0: []byte("old")},
				total:     2,
				size:      3,
				partCount: 1,
				created:   time.Now().Add(-streamHoldTimeout - time.Second),
			},
		},
		bytes: 3,
		parts: 1,
		mgr:   mgr,
	}
	mgr.chunkBytes = 3
	mgr.chunkParts = 1

	if _, done, err := store.add("fresh", 0, 2, []byte("new")); done || err != nil {
		t.Fatalf("fresh chunk: done=%v err=%v", done, err)
	}
	if _, ok := store.items["expired"]; ok {
		t.Fatal("expired chunk set was not removed")
	}
	if store.bytes != 3 || store.parts != 1 {
		t.Fatalf("store after reap = %d bytes, %d parts; want only the fresh part", store.bytes, store.parts)
	}
	if mgr.chunkBytes != 3 || mgr.chunkParts != 1 {
		t.Fatalf("shared accounting after reap = %d bytes, %d parts; want only the fresh part", mgr.chunkBytes, mgr.chunkParts)
	}
}

func TestStreamChunkStoreBoundsZeroByteIncompleteSets(t *testing.T) {
	const desiredMaxSets = 256
	store := &streamChunkStore{items: make(map[string]*streamChunkSet)}

	for i := 0; i < desiredMaxSets; i++ {
		id := fmt.Sprintf("set-%d", i)
		if _, done, err := store.add(id, 0, 2, nil); done || err != nil {
			t.Fatalf("add %s: done=%v err=%v", id, done, err)
		}
	}
	if _, done, err := store.add("one-too-many", 0, 2, nil); done || !errors.Is(err, ErrStreamFull) {
		t.Fatalf("add past item bound: done=%v err=%v, want ErrStreamFull", done, err)
	}
	if len(store.items) != desiredMaxSets || store.bytes != 0 || store.parts != desiredMaxSets {
		t.Fatalf("store after overflow = %d sets, %d bytes, %d parts", len(store.items), store.bytes, store.parts)
	}
}

func TestStreamChunkStoresShareGlobalByteBudget(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)

	// Simulate other sessions retaining all but one byte of the manager-wide
	// allowance. This store's two-byte chunk must receive retryable
	// backpressure rather than multiplying the per-session allowance.
	mgr.chunkBytes = streamMaxChunkBytesGlobal - 1
	store := &streamChunkStore{items: make(map[string]*streamChunkSet), mgr: mgr}
	if _, done, err := store.add("over-global-budget", 0, 2, []byte("xx")); done || !errors.Is(err, ErrStreamFull) {
		t.Fatalf("add beyond global chunk budget: done=%v err=%v, want ErrStreamFull", done, err)
	}
	if len(store.items) != 0 || store.bytes != 0 {
		t.Fatalf("rejected chunk retained %d sets and %d bytes", len(store.items), store.bytes)
	}

	// Existing partial uploads survive transient global backpressure and can
	// resume without losing or duplicating their prefix once another store frees
	// capacity. Removal must return all accounting to the manager.
	mgr.chunkBytes = streamMaxChunkBytesGlobal - 2
	if _, done, err := store.add("resume", 0, 2, []byte("a")); done || err != nil {
		t.Fatalf("add prefix: done=%v err=%v", done, err)
	}
	if _, done, err := store.add("resume", 1, 2, []byte("bc")); done || !errors.Is(err, ErrStreamFull) {
		t.Fatalf("add blocked suffix: done=%v err=%v, want ErrStreamFull", done, err)
	}
	if got := len(store.items["resume"].parts); got != 1 {
		t.Fatalf("blocked suffix discarded the %d-part prefix", got)
	}
	mgr.releaseChunkResources(streamMaxChunkBytesGlobal-2, 0)
	assembled, done, err := store.add("resume", 1, 2, []byte("bc"))
	if !done || err != nil || string(assembled) != "abc" {
		t.Fatalf("retry after shared capacity freed: data=%q done=%v err=%v", assembled, done, err)
	}
	store.remove("resume")
	if mgr.chunkBytes != 0 || mgr.chunkParts != 0 {
		t.Fatalf("delivered chunk set retained %d shared bytes and %d parts", mgr.chunkBytes, mgr.chunkParts)
	}

	// Session shutdown is the other ordinary release path.
	if _, done, err := store.add("close", 0, 2, []byte("held")); done || err != nil {
		t.Fatalf("add before close: done=%v err=%v", done, err)
	}
	store.close()
	if mgr.chunkBytes != 0 || mgr.chunkParts != 0 {
		t.Fatalf("closed chunk store retained %d shared bytes and %d parts", mgr.chunkBytes, mgr.chunkParts)
	}
}

func TestStreamChunkGlobalBudgetIsAtomicAcrossSessions(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)

	const chunkSize = 1 << 20
	const storeCount = streamMaxChunkBytesGlobal/chunkSize + 32
	payload := make([]byte, chunkSize)
	stores := make([]*streamChunkStore, storeCount)
	var accepted int64
	var wg sync.WaitGroup

	for i := range stores {
		store := &streamChunkStore{items: make(map[string]*streamChunkSet), mgr: mgr}
		stores[i] = store
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, done, err := store.add("partial", 0, 2, payload)
			switch {
			case done:
				t.Errorf("incomplete chunk unexpectedly assembled")
			case err == nil:
				atomic.AddInt64(&accepted, 1)
			case !errors.Is(err, ErrStreamFull):
				t.Errorf("add = %v, want ErrStreamFull", err)
			}
		}()
	}
	wg.Wait()

	wantAccepted := int64(streamMaxChunkBytesGlobal / chunkSize)
	if accepted != wantAccepted {
		t.Fatalf("accepted %d concurrent stores, want %d", accepted, wantAccepted)
	}
	mgr.mu.Lock()
	gotBytes := mgr.chunkBytes
	gotParts := mgr.chunkParts
	mgr.mu.Unlock()
	if gotBytes != streamMaxChunkBytesGlobal {
		t.Fatalf("shared accounting = %d, want %d", gotBytes, streamMaxChunkBytesGlobal)
	}
	if gotParts != int(accepted) {
		t.Fatalf("shared part accounting = %d, want %d", gotParts, accepted)
	}

	for _, store := range stores {
		wg.Add(1)
		go func() {
			defer wg.Done()
			store.close()
		}()
	}
	wg.Wait()
	mgr.mu.Lock()
	gotBytes = mgr.chunkBytes
	gotParts = mgr.chunkParts
	mgr.mu.Unlock()
	if gotBytes != 0 || gotParts != 0 {
		t.Fatalf("concurrent shutdown retained %d shared bytes and %d parts", gotBytes, gotParts)
	}
}

func TestStreamChunkReassemblyCopyFitsEffectiveGlobalBudget(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)

	// Two maximum-size sets consume the admitted 128 MiB payload budget. Their
	// final chunks arrive concurrently, so each store temporarily retains both
	// its parts and its 64 MiB contiguous result. The conservative 2x admission
	// bound keeps that peak within the 256 MiB effective ceiling.
	half := make([]byte, streamMaxFrameBytes/2)
	stores := []*streamChunkStore{
		{items: make(map[string]*streamChunkSet), mgr: mgr},
		{items: make(map[string]*streamChunkSet), mgr: mgr},
	}
	for i, store := range stores {
		if _, done, err := store.add("maximum", 0, 2, half); done || err != nil {
			t.Fatalf("store %d first half: done=%v err=%v", i, done, err)
		}
	}

	var wg sync.WaitGroup
	errs := make(chan error, len(stores))
	for i, store := range stores {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assembled, done, err := store.add("maximum", 1, 2, half)
			if err != nil {
				errs <- fmt.Errorf("store %d final half: %w", i, err)
				return
			}
			if !done || len(assembled) != streamMaxFrameBytes {
				errs <- fmt.Errorf("store %d assembled %d bytes, done=%v", i, len(assembled), done)
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}

	if got := mgr.chunkBytes; got != streamMaxChunkBytesGlobal {
		t.Fatalf("shared payload accounting = %d, want %d", got, streamMaxChunkBytesGlobal)
	}
	overflow := &streamChunkStore{items: make(map[string]*streamChunkSet), mgr: mgr}
	if _, done, err := overflow.add("third", 0, 2, []byte{1}); done || !errors.Is(err, ErrStreamFull) {
		t.Fatalf("third upload at effective memory ceiling: done=%v err=%v, want ErrStreamFull", done, err)
	}

	for _, store := range stores {
		store.remove("maximum")
	}
	if mgr.chunkBytes != 0 || mgr.chunkParts != 0 {
		t.Fatalf("completed sets retained %d bytes and %d parts", mgr.chunkBytes, mgr.chunkParts)
	}
}

func TestStreamChunkStoresShareGlobalMetadataBudget(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)

	storeCount := streamMaxChunkPartsGlobal/streamMaxChunkSets + 1
	stores := make([]*streamChunkStore, storeCount)
	for i := range stores {
		stores[i] = &streamChunkStore{items: make(map[string]*streamChunkSet), mgr: mgr}
	}

	for i := 0; i < streamMaxChunkPartsGlobal; i++ {
		store := stores[i/streamMaxChunkSets]
		if _, done, err := store.add(fmt.Sprintf("partial-%d", i), 0, 2, nil); done || err != nil {
			t.Fatalf("add metadata part %d: done=%v err=%v", i, done, err)
		}
	}
	last := stores[len(stores)-1]
	if _, done, err := last.add("over-metadata-budget", 0, 2, nil); done || !errors.Is(err, ErrStreamFull) {
		t.Fatalf("add beyond metadata budget: done=%v err=%v, want ErrStreamFull", done, err)
	}
	if len(last.items) != 0 {
		t.Fatalf("rejected metadata chunk retained %d sets", len(last.items))
	}
	if mgr.chunkBytes != 0 || mgr.chunkParts != streamMaxChunkPartsGlobal {
		t.Fatalf("shared accounting at capacity = %d bytes, %d parts", mgr.chunkBytes, mgr.chunkParts)
	}

	stores[0].close()
	if _, done, err := last.add("after-release", 0, 2, nil); done || err != nil {
		t.Fatalf("add after metadata release: done=%v err=%v", done, err)
	}
	for _, store := range stores[1:] {
		store.close()
	}
	if mgr.chunkBytes != 0 || mgr.chunkParts != 0 {
		t.Fatalf("shutdown retained %d shared bytes and %d parts", mgr.chunkBytes, mgr.chunkParts)
	}
}

func TestStreamSessionCloseReleasesChunkAccounting(t *testing.T) {
	_, s := newTestSession(t)
	store := s.chunks()
	if _, done, err := store.add("partial", 0, 2, []byte("retained")); done || err != nil {
		t.Fatalf("add partial set: done=%v err=%v", done, err)
	}

	s.close()
	if len(store.items) != 0 || store.bytes != 0 || store.parts != 0 {
		t.Fatalf("closed session retained %d sets, %d bytes, and %d parts", len(store.items), store.bytes, store.parts)
	}
	if _, _, err := store.add("late", 0, 2, []byte("late")); !errors.Is(err, ErrStreamClosed) {
		t.Fatalf("add after session close = %v, want ErrStreamClosed", err)
	}
}

func TestStreamChunkCapacityReturnsRetryableHTTPStatus(t *testing.T) {
	mgr, s := newTestSession(t)
	s.generation = 1
	newTestConn(s, 7)
	a := &App{streams: mgr}
	store := s.chunks()
	for i := 0; i < streamMaxChunkSets; i++ {
		id := fmt.Sprintf("held-%d", i)
		if _, done, err := store.add(id, 0, 2, nil); done || err != nil {
			t.Fatalf("fill %s: done=%v err=%v", id, done, err)
		}
	}

	send := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, streamPathSend, nil)
		req.Header.Set(streamHeaderSession, s.id)
		req.Header.Set(streamHeaderGeneration, "1")
		req.Header.Set(streamHeaderConn, "7")
		req.Header.Set(streamHeaderKind, strconv.Itoa(int(frameData)))
		req.Header.Set(streamHeaderChunkID, "retry-after-capacity")
		req.Header.Set(streamHeaderChunkIndex, "0")
		req.Header.Set(streamHeaderChunkTotal, "2")
		rw := httptest.NewRecorder()
		a.serveStreamSend(rw, req)
		return rw
	}

	if rw := send(); rw.Code != http.StatusTooManyRequests || rw.Header().Get("Retry-After") == "" {
		t.Fatalf("capacity response = %d Retry-After=%q, want retryable 429", rw.Code, rw.Header().Get("Retry-After"))
	}
	store.remove("held-0")
	if rw := send(); rw.Code != http.StatusNoContent {
		t.Fatalf("retry after capacity freed = %d, want 204", rw.Code)
	}
}

func TestStreamChunkPermutationsAndExpiredIDReuse(t *testing.T) {
	parts := []string{"a", "b", "c"}
	orders := [][]int{
		{0, 1, 2}, {0, 2, 1}, {1, 0, 2},
		{1, 2, 0}, {2, 0, 1}, {2, 1, 0},
	}
	for _, order := range orders {
		name := fmt.Sprint(order)
		t.Run(name, func(t *testing.T) {
			store := &streamChunkStore{items: make(map[string]*streamChunkSet)}
			var assembled []byte
			for step, index := range order {
				got, done, err := store.add("permutation", index, len(parts), []byte(parts[index]))
				if err != nil {
					t.Fatalf("step %d: %v", step, err)
				}
				if done != (step == len(order)-1) {
					t.Fatalf("step %d done=%v", step, done)
				}
				if done {
					assembled = got
				}
			}
			if string(assembled) != "abc" {
				t.Fatalf("assembled %q, want abc", assembled)
			}
			store.remove("permutation")
			if len(store.items) != 0 || store.bytes != 0 {
				t.Fatalf("delivered permutation retained %d sets and %d bytes", len(store.items), store.bytes)
			}
		})
	}

	store := &streamChunkStore{items: make(map[string]*streamChunkSet)}
	if _, done, err := store.add("reused", 0, 2, []byte("old")); done || err != nil {
		t.Fatalf("old incomplete set: done=%v err=%v", done, err)
	}
	store.items["reused"].created = time.Now().Add(-streamHoldTimeout - time.Second)
	if _, done, err := store.add("reused", 1, 2, []byte("new-b")); done || err != nil {
		t.Fatalf("first chunk after expiry: done=%v err=%v", done, err)
	}
	got, done, err := store.add("reused", 0, 2, []byte("new-a"))
	if !done || err != nil || string(got) != "new-anew-b" {
		t.Fatalf("reused id assembled %q done=%v err=%v", got, done, err)
	}
	store.remove("reused")
	if len(store.items) != 0 || store.bytes != 0 {
		t.Fatalf("reused id retained %d sets and %d bytes", len(store.items), store.bytes)
	}
}

func TestStreamRejectsOversizedChunkIdentifier(t *testing.T) {
	mgr, s := newTestSession(t)
	s.generation = 1
	newTestConn(s, 7)
	req := httptest.NewRequest(http.MethodPost, streamPathSend, nil)
	req.Header.Set(streamHeaderSession, s.id)
	req.Header.Set(streamHeaderGeneration, "1")
	req.Header.Set(streamHeaderConn, "7")
	req.Header.Set(streamHeaderKind, strconv.Itoa(int(frameData)))
	req.Header.Set(streamHeaderChunkID, strings.Repeat("x", streamMaxChunkIDLen+1))
	req.Header.Set(streamHeaderChunkIndex, "0")
	req.Header.Set(streamHeaderChunkTotal, "2")
	rw := httptest.NewRecorder()
	(&App{streams: mgr}).serveStreamSend(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("oversized chunk id status = %d, want 400", rw.Code)
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

func TestStreamOutboundBytesAreBoundedAcrossSessions(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	payload := make([]byte, 1<<20)
	var sessions []*streamSession

	for sessionIndex := 0; sessionIndex < streamOutQueueBytesGlobal/(8<<20); sessionIndex++ {
		s := newStreamSession(fmt.Sprintf("session-%d", sessionIndex), uint(sessionIndex+1), mgr)
		mgr.sessions[s.id] = s
		sessions = append(sessions, s)
		c := newTestConn(s, 1)
		for frame := 0; frame < 8; frame++ {
			if err := c.TrySend(payload); err != nil {
				t.Fatalf("session %d frame %d: %v", sessionIndex, frame, err)
			}
		}
	}

	overflow := newStreamSession("overflow", 999, mgr)
	mgr.sessions[overflow.id] = overflow
	sessions = append(sessions, overflow)
	if err := newTestConn(overflow, 1).TrySend(payload); err != ErrStreamFull {
		t.Fatalf("TrySend beyond application byte budget = %v, want ErrStreamFull", err)
	}
	if got := mgr.outBytes.Load(); got != streamOutQueueBytesGlobal {
		t.Fatalf("outbound bytes = %d, want %d", got, streamOutQueueBytesGlobal)
	}

	for _, s := range sessions {
		s.close()
	}
	if got := mgr.outBytes.Load(); got != 0 {
		t.Fatalf("outbound bytes after cleanup = %d, want 0", got)
	}
	if got := mgr.outFrames.Load(); got != 0 {
		t.Fatalf("outbound frames after cleanup = %d, want 0", got)
	}
}

func TestStreamOutboundFramesAreBoundedAcrossSessions(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	var sessions []*streamSession

	for sessionIndex := 0; sessionIndex < streamOutQueueDepthGlobal/streamOutQueueDepth; sessionIndex++ {
		s := newStreamSession(fmt.Sprintf("session-%d", sessionIndex), uint(sessionIndex+1), mgr)
		mgr.sessions[s.id] = s
		sessions = append(sessions, s)
		c := newTestConn(s, 1)
		for frame := 0; frame < streamOutQueueDepth; frame++ {
			if err := c.TrySend(nil); err != nil {
				t.Fatalf("session %d frame %d: %v", sessionIndex, frame, err)
			}
		}
	}

	overflow := newStreamSession("overflow", 999, mgr)
	mgr.sessions[overflow.id] = overflow
	sessions = append(sessions, overflow)
	if err := newTestConn(overflow, 1).TrySend(nil); err != ErrStreamFull {
		t.Fatalf("TrySend beyond application frame budget = %v, want ErrStreamFull", err)
	}
	if got := mgr.outFrames.Load(); got != streamOutQueueDepthGlobal {
		t.Fatalf("outbound frames = %d, want %d", got, streamOutQueueDepthGlobal)
	}

	for _, s := range sessions {
		s.close()
	}
	if got := mgr.outFrames.Load(); got != 0 {
		t.Fatalf("outbound frames after cleanup = %d, want 0", got)
	}
}

func TestStreamInboundBytesAreBoundedAcrossConnections(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	payload := make([]byte, 1<<20)
	var conns []*StreamConn

	for sessionIndex := 0; sessionIndex < streamInQueueBytesGlobal/(8<<20); sessionIndex++ {
		s := newStreamSession(fmt.Sprintf("session-%d", sessionIndex), uint(sessionIndex+1), mgr)
		mgr.sessions[s.id] = s
		c := newTestConn(s, 1)
		conns = append(conns, c)
		for frame := 0; frame < 8; frame++ {
			if err := c.deliver(payload); err != nil {
				t.Fatalf("connection %d frame %d: %v", sessionIndex, frame, err)
			}
		}
	}

	overflowSession := newStreamSession("overflow", 999, mgr)
	mgr.sessions[overflowSession.id] = overflowSession
	overflow := newTestConn(overflowSession, 1)
	conns = append(conns, overflow)
	if err := overflow.deliver(payload); err != ErrStreamFull {
		t.Fatalf("deliver beyond application byte budget = %v, want ErrStreamFull", err)
	}
	if got := mgr.inBytes.Load(); got != streamInQueueBytesGlobal {
		t.Fatalf("inbound bytes = %d, want %d", got, streamInQueueBytesGlobal)
	}

	for _, c := range conns {
		c.shutdown()
	}
	if got := mgr.inBytes.Load(); got != 0 {
		t.Fatalf("inbound bytes after cleanup = %d, want 0", got)
	}
	if got := mgr.inFrames.Load(); got != 0 {
		t.Fatalf("inbound frames after cleanup = %d, want 0", got)
	}
}

func TestStreamInboundFramesAreBoundedAcrossConnections(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	var conns []*StreamConn

	for sessionIndex := 0; sessionIndex < streamInQueueDepthGlobal/streamInQueueDepth; sessionIndex++ {
		s := newStreamSession(fmt.Sprintf("session-%d", sessionIndex), uint(sessionIndex+1), mgr)
		mgr.sessions[s.id] = s
		c := newTestConn(s, 1)
		conns = append(conns, c)
		for frame := 0; frame < streamInQueueDepth; frame++ {
			if err := c.deliver(nil); err != nil {
				t.Fatalf("connection %d frame %d: %v", sessionIndex, frame, err)
			}
		}
	}

	overflowSession := newStreamSession("overflow", 999, mgr)
	mgr.sessions[overflowSession.id] = overflowSession
	overflow := newTestConn(overflowSession, 1)
	conns = append(conns, overflow)
	if err := overflow.deliver(nil); err != ErrStreamFull {
		t.Fatalf("deliver beyond application frame budget = %v, want ErrStreamFull", err)
	}
	if got := mgr.inFrames.Load(); got != streamInQueueDepthGlobal {
		t.Fatalf("inbound frames = %d, want %d", got, streamInQueueDepthGlobal)
	}

	for _, c := range conns {
		c.shutdown()
	}
	if got := mgr.inFrames.Load(); got != 0 {
		t.Fatalf("inbound frames after cleanup = %d, want 0", got)
	}
}

func TestStreamGlobalConnectionAdmissionCoversDesktopOpen(t *testing.T) {
	mgr, s := newTestSession(t)
	mgr.handlers["hold"] = func(c *StreamConn) { <-c.Context().Done() }

	for i := 0; i < streamMaxConnectionsGlobal-1; i++ {
		if !mgr.reserveOpen() {
			t.Fatalf("reserve open %d unexpectedly failed", i)
		}
	}
	if err := s.open(1, "hold"); err != nil {
		t.Fatalf("last globally admitted open: %v", err)
	}
	if err := s.open(2, "hold"); err != ErrStreamFull {
		t.Fatalf("open beyond global lifecycle budget = %v, want ErrStreamFull", err)
	}

	s.close()
	for i := 0; i < streamMaxConnectionsGlobal-1; i++ {
		mgr.releaseOpenReservation()
	}
	if got := mgr.lifecycles.Load(); got != 0 {
		t.Fatalf("lifecycles after cleanup = %d, want 0", got)
	}
	if got := mgr.outControls.Load(); got != 0 {
		t.Fatalf("controls after cleanup = %d, want 0", got)
	}
}

func TestStreamCloseAccountingIsIndependentOfLifecycle(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)

	if !mgr.reserveLifecycle() {
		t.Fatal("failed to reserve lifecycle slot")
	}
	if err := mgr.reserveOutbound(frameClose, 0, false, nil); err != nil {
		t.Fatalf("reserve close frame: %v", err)
	}
	if got := mgr.outCloses.Load(); got != 1 {
		t.Fatalf("close slots after reserve = %d, want 1", got)
	}

	// Releasing the close frame — whether it was drained or never queued at all
	// — returns only the close slot. The connection keeps its lifecycle slot
	// until shutdown, so neither path depends on the other's ordering.
	mgr.releaseOutbound(frameClose, 0)
	if got := mgr.outCloses.Load(); got != 0 {
		t.Fatalf("close slots after release = %d, want 0", got)
	}
	if got := mgr.lifecycles.Load(); got != 1 {
		t.Fatalf("lifecycle after close released = %d, want 1", got)
	}

	mgr.releaseLifecycle()
	if got := mgr.lifecycles.Load(); got != 0 {
		t.Fatalf("lifecycle after shutdown = %d, want 0", got)
	}
}

// A close racing its own session's teardown must still return the connection's
// lifecycle slot. These leak permanently, so a slow drip eventually refuses
// every new connection application-wide.
func TestStreamLifecycleSlotSurvivesCloseTeardownRace(t *testing.T) {
	const iterations = 3000
	for i := 0; i < iterations; i++ {
		mgr := newStreamManager(nil)
		s := newStreamSession("sess", 1, mgr)
		mgr.sessions["sess"] = s

		if !mgr.reserveOpen() {
			t.Fatal("failed to reserve open")
		}
		ctx, cancel := context.WithCancel(context.Background())
		c := &StreamConn{
			id: 1, name: "t", sink: s, ctx: ctx, cancel: cancel,
			manager: mgr, lifecycle: true,
		}
		c.inCond = sync.NewCond(&c.inMu)
		s.mu.Lock()
		s.conns[1] = c
		s.out = append(s.out, outFrame{connID: 1, kind: frameOpen})
		s.outControls++
		s.mu.Unlock()

		var wg sync.WaitGroup
		wg.Add(2)
		go func() { defer wg.Done(); _ = c.Close() }()
		go func() { defer wg.Done(); s.close() }()
		wg.Wait()

		if got := mgr.lifecycles.Load(); got != 0 {
			t.Fatalf("iteration %d: lifecycle slots outstanding = %d, want 0", i, got)
		}
		if got := mgr.outCloses.Load(); got != 0 {
			t.Fatalf("iteration %d: close slots outstanding = %d, want 0", i, got)
		}
	}
}

func TestStreamBlockedOnGlobalBudgetWakesWhenConnectionCloses(t *testing.T) {
	mgr, s := newTestSession(t)
	c := newTestConn(s, 1)
	if !reserveCounter(&mgr.outBytes, streamOutQueueBytesGlobal, streamOutQueueBytesGlobal) {
		t.Fatal("failed to fill global byte budget")
	}

	done := make(chan error, 1)
	go func() { done <- c.Send([]byte("blocked")) }()
	select {
	case err := <-done:
		t.Fatalf("Send returned before capacity or cancellation: %v", err)
	case <-time.After(50 * time.Millisecond):
	}

	c.shutdown()
	select {
	case err := <-done:
		if err != ErrStreamClosed {
			t.Fatalf("blocked Send after close = %v, want ErrStreamClosed", err)
		}
	case <-time.After(time.Second):
		t.Fatal("blocked Send did not wake after connection close")
	}
	releaseCounter(&mgr.outBytes, streamOutQueueBytesGlobal)
}

func TestStreamManagerCloseWakesEveryGlobalBudgetWaiter(t *testing.T) {
	mgr := newStreamManager(nil)
	if !reserveCounter(&mgr.outBytes, streamOutQueueBytesGlobal, streamOutQueueBytesGlobal) {
		t.Fatal("failed to fill global byte budget")
	}

	const waiters = 16
	done := make(chan error, waiters)
	for i := 0; i < waiters; i++ {
		s := newStreamSession(fmt.Sprintf("waiter-%d", i), uint(i+1), mgr)
		mgr.sessions[s.id] = s
		c := newTestConn(s, 1)
		go func() { done <- c.Send([]byte("blocked")) }()
	}
	select {
	case err := <-done:
		t.Fatalf("Send returned before manager close: %v", err)
	case <-time.After(50 * time.Millisecond):
	}

	mgr.close()
	for i := 0; i < waiters; i++ {
		select {
		case err := <-done:
			if err != ErrStreamClosed {
				t.Fatalf("waiter %d after manager close = %v, want ErrStreamClosed", i, err)
			}
		case <-time.After(time.Second):
			t.Fatalf("waiter %d did not wake after manager close", i)
		}
	}
}
