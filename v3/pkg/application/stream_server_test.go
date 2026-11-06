//go:build server

package application

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
)

// Server mode is only required to work, not to be fast: these check that a real
// WebSocket carries frames both ways through the same StreamConn the desktop
// build uses, and that the handler's lifetime tracks the socket.

func newServerTestApp(t *testing.T) *App {
	t.Helper()
	a := &App{Logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
	a.streams = newStreamManager(a)
	t.Cleanup(a.streams.close)
	return a
}

func TestServerModeStreamOriginPolicy(t *testing.T) {
	tests := []struct {
		name       string
		server     ServerOptions
		origin     string
		wantAccept bool
	}{
		{
			name:   "cross origin rejected by default",
			origin: "https://untrusted.example",
		},
		{
			name: "configured origin accepted",
			server: ServerOptions{
				WebSocketOriginPatterns: []string{"trusted.example"},
			},
			origin:     "https://trusted.example",
			wantAccept: true,
		},
		{
			name: "all origins require explicit opt in",
			server: ServerOptions{
				WebSocketAllowAllOrigins: true,
			},
			origin:     "https://untrusted.example",
			wantAccept: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newServerTestApp(t)
			a.options.Server = tt.server
			a.HandleStream("origin", func(c *StreamConn) {
				<-c.Context().Done()
			})

			srv := httptest.NewServer(http.HandlerFunc(a.serveStreamWS))
			defer srv.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			conn, resp, err := websocket.Dial(ctx, wsURL(srv.URL)+"?name=origin", &websocket.DialOptions{
				HTTPHeader: http.Header{"Origin": []string{tt.origin}},
			})
			if tt.wantAccept {
				if err != nil {
					t.Fatalf("dial from %s: %v", tt.origin, err)
				}
				conn.CloseNow()
				return
			}
			if err == nil {
				conn.CloseNow()
				t.Fatalf("dial from %s unexpectedly succeeded", tt.origin)
			}
			if resp == nil || resp.StatusCode != http.StatusForbidden {
				t.Fatalf("rejected response = %#v, want HTTP %d", resp, http.StatusForbidden)
			}
		})
	}
}

func TestServerModeStreamEchoesBothWays(t *testing.T) {
	a := newServerTestApp(t)

	a.HandleStream("echo", func(c *StreamConn) {
		defer c.Close()
		for {
			frame, err := c.Receive()
			if err != nil {
				return
			}
			if err := c.Send(append([]byte("re:"), frame...)); err != nil {
				return
			}
		}
	})

	srv := httptest.NewServer(http.HandlerFunc(a.serveStreamWS))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL(srv.URL)+"?name=echo", nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.CloseNow()

	for i := 0; i < 5; i++ {
		if err := conn.Write(ctx, websocket.MessageBinary, []byte("ping")); err != nil {
			t.Fatalf("write %d: %v", i, err)
		}
		typ, data, err := conn.Read(ctx)
		if err != nil {
			t.Fatalf("read %d: %v", i, err)
		}
		if typ != websocket.MessageBinary {
			t.Fatalf("message %d type = %v, want binary", i, typ)
		}
		if string(data) != "re:ping" {
			t.Fatalf("message %d = %q, want %q", i, data, "re:ping")
		}
	}
}

// A frame far larger than the desktop build's per-window buffer must still
// cross, since on a socket the send buffer is the backpressure and there is no
// queue cap to trip over.
func TestServerModeCarriesLargeFrames(t *testing.T) {
	a := newServerTestApp(t)

	payload := make([]byte, 12<<20) // 12 MB, past streamOutQueueBytes
	for i := range payload {
		payload[i] = byte(i)
	}

	a.HandleStream("big", func(c *StreamConn) {
		defer c.Close()
		if err := c.Send(payload); err != nil {
			t.Errorf("send: %v", err)
		}
		<-c.Context().Done()
	})

	srv := httptest.NewServer(http.HandlerFunc(a.serveStreamWS))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL(srv.URL)+"?name=big", nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.CloseNow()
	conn.SetReadLimit(int64(len(payload)) * 2)

	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(data) != len(payload) {
		t.Fatalf("received %d bytes, want %d", len(data), len(payload))
	}
	for i := range data {
		if data[i] != payload[i] {
			t.Fatalf("payload differs at byte %d", i)
		}
	}
}

// coder/websocket defaults to a 32 KiB inbound limit. The Stream endpoint
// promises streamMaxFrameBytes, so exercise the client-to-Go direction over a
// real socket rather than relying on the independent Go-to-client test above.
func TestServerModeAcceptsLargeInboundFrame(t *testing.T) {
	a := newServerTestApp(t)

	a.HandleStream("large-inbound", func(c *StreamConn) {
		frame, err := c.Receive()
		if err != nil {
			t.Errorf("receive: %v", err)
			return
		}
		if err := c.Send(frame); err != nil {
			t.Errorf("echo: %v", err)
		}
	})

	srv := httptest.NewServer(http.HandlerFunc(a.serveStreamWS))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	conn, _, err := websocket.Dial(ctx, wsURL(srv.URL)+"?name=large-inbound", nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.CloseNow()

	payload := make([]byte, 1<<20) // safely beyond coder/websocket's 32 KiB default
	for i := range payload {
		payload[i] = byte(i)
	}
	if err := conn.Write(ctx, websocket.MessageBinary, payload); err != nil {
		t.Fatalf("write: %v", err)
	}
	conn.SetReadLimit(int64(len(payload)) * 2)
	_, echoed, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read echo: %v", err)
	}
	if len(echoed) != len(payload) {
		t.Fatalf("received %d bytes, want %d", len(echoed), len(payload))
	}
	for i := range echoed {
		if echoed[i] != payload[i] {
			t.Fatalf("payload differs at byte %d", i)
		}
	}
}

func TestServerModeRejectsOversizedStreamName(t *testing.T) {
	a := &App{streams: newStreamManager(nil)}
	t.Cleanup(a.streams.close)
	req := httptest.NewRequest(http.MethodGet, "/wails/stream?name="+strings.Repeat("n", streamMaxNameLen+1), nil)
	rw := httptest.NewRecorder()

	a.serveStreamWS(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("oversized name status = %d, want 400", rw.Code)
	}
}

func TestServerModeAcceptsMaximumLengthStreamName(t *testing.T) {
	a := newServerTestApp(t)
	maximum := strings.Repeat("m", streamMaxNameLen)
	a.HandleStream(maximum, func(c *StreamConn) {})
	srv := httptest.NewServer(http.HandlerFunc(a.serveStreamWS))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, _, err := websocket.Dial(ctx, wsURL(srv.URL)+"?name="+maximum, nil)
	if err != nil {
		t.Fatalf("dial maximum-length stream name: %v", err)
	}
	defer conn.CloseNow()
}

func TestServerModeRejectsConnectionBeyondGlobalLifecycleBudget(t *testing.T) {
	a := newServerTestApp(t)
	a.HandleStream("hold", func(c *StreamConn) { <-c.Context().Done() })
	for i := 0; i < streamMaxConnectionsGlobal; i++ {
		if !a.streams.reserveLifecycle() {
			t.Fatalf("reserve lifecycle %d unexpectedly failed", i)
		}
	}
	defer func() {
		for i := 0; i < streamMaxConnectionsGlobal; i++ {
			a.streams.releaseLifecycle()
		}
	}()

	req := httptest.NewRequest(http.MethodGet, "/wails/stream?name=hold", nil)
	rw := httptest.NewRecorder()
	a.serveStreamWS(rw, req)
	if rw.Code != http.StatusServiceUnavailable {
		t.Fatalf("connection beyond global lifecycle budget status = %d, want 503", rw.Code)
	}
}

func TestServerModeSinkUsesApplicationWideOutboundBudget(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	if !reserveCounter(&mgr.outBytes, streamOutQueueBytesGlobal, streamOutQueueBytesGlobal) {
		t.Fatal("failed to fill global outbound byte budget")
	}
	defer releaseCounter(&mgr.outBytes, streamOutQueueBytesGlobal)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sink := &wsStreamSink{mgr: mgr}
	sink.cond = sync.NewCond(&sink.mu)
	c := &StreamConn{sink: sink, ctx: ctx, cancel: cancel, manager: mgr}
	c.inCond = sync.NewCond(&c.inMu)

	if err := c.TrySend([]byte("overflow")); err != ErrStreamFull {
		t.Fatalf("server TrySend beyond global byte budget = %v, want ErrStreamFull", err)
	}
}

func TestServerModeBlockedGlobalProducerWakesOnCapacity(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	if !reserveCounter(&mgr.outBytes, streamOutQueueBytesGlobal, streamOutQueueBytesGlobal) {
		t.Fatal("failed to fill global outbound byte budget")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sink := &wsStreamSink{mgr: mgr}
	sink.cond = sync.NewCond(&sink.mu)
	c := &StreamConn{sink: sink, ctx: ctx, cancel: cancel, manager: mgr}
	c.inCond = sync.NewCond(&c.inMu)

	done := make(chan error, 1)
	go func() { done <- c.Send([]byte("released")) }()
	select {
	case err := <-done:
		t.Fatalf("Send returned before capacity was released: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
	releaseCounter(&mgr.outBytes, streamOutQueueBytesGlobal)
	mgr.signalBudget()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Send after capacity release: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("blocked server Send did not wake after capacity release")
	}
	sink.releaseQueued()
	if got := mgr.outBytes.Load(); got != 0 {
		t.Fatalf("outbound bytes after queue cleanup = %d, want 0", got)
	}
}

func TestServerModeBlockedGlobalProducerWakesOnShutdown(t *testing.T) {
	mgr := newStreamManager(nil)
	t.Cleanup(mgr.close)
	if !reserveCounter(&mgr.outBytes, streamOutQueueBytesGlobal, streamOutQueueBytesGlobal) {
		t.Fatal("failed to fill global outbound byte budget")
	}
	ctx, cancel := context.WithCancel(context.Background())
	sink := &wsStreamSink{mgr: mgr}
	sink.cond = sync.NewCond(&sink.mu)
	c := &StreamConn{sink: sink, ctx: ctx, cancel: cancel, manager: mgr}
	c.inCond = sync.NewCond(&c.inMu)

	done := make(chan error, 1)
	go func() { done <- c.Send([]byte("blocked")) }()
	select {
	case err := <-done:
		t.Fatalf("Send returned before shutdown: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
	c.shutdown()
	select {
	case err := <-done:
		if err != ErrStreamClosed {
			t.Fatalf("Send after shutdown = %v, want ErrStreamClosed", err)
		}
	case <-time.After(time.Second):
		t.Fatal("blocked server Send did not wake after shutdown")
	}
	releaseCounter(&mgr.outBytes, streamOutQueueBytesGlobal)
}

// TrySend is documented as non-blocking in both transports. Pin the
// server-mode sink contract directly so a future inline websocket.Write does
// not make a broker callback or fan-out producer block again.
func TestServerModeTrySendReportsFullWithoutBlocking(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sink := &wsStreamSink{}
	sink.cond = sync.NewCond(&sink.mu)
	c := &StreamConn{sink: sink, ctx: ctx, cancel: cancel}
	c.inCond = sync.NewCond(&c.inMu)

	for i := 0; i < streamOutQueueDepth; i++ {
		if err := c.TrySend([]byte("x")); err != nil {
			t.Fatalf("fill %d: %v", i, err)
		}
	}

	done := make(chan error, 1)
	go func() { done <- c.TrySend([]byte("overflow")) }()
	select {
	case err := <-done:
		if err != ErrStreamFull {
			t.Fatalf("TrySend on full server queue = %v, want ErrStreamFull", err)
		}
	case <-time.After(time.Second):
		t.Fatal("TrySend blocked on a full server queue")
	}
}

func TestServerModeRefusesUnknownStream(t *testing.T) {
	a := newServerTestApp(t)

	srv := httptest.NewServer(http.HandlerFunc(a.serveStreamWS))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "?name=nosuchstream")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", resp.StatusCode)
	}
}

// Closing the socket must end the handler, the way a reload does on the desktop.
func TestServerModeSocketCloseEndsHandler(t *testing.T) {
	a := newServerTestApp(t)

	finished := make(chan error, 1)
	a.HandleStream("hold", func(c *StreamConn) {
		_, err := c.Receive()
		finished <- err
	})

	srv := httptest.NewServer(http.HandlerFunc(a.serveStreamWS))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL(srv.URL)+"?name=hold", nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	_ = conn.Close(websocket.StatusNormalClosure, "done")

	select {
	case err := <-finished:
		if err != ErrStreamClosed {
			t.Fatalf("Receive after socket close = %v, want ErrStreamClosed", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("handler never returned after the socket closed")
	}
}

func TestServerModeAppContextCancellationEndsHandler(t *testing.T) {
	a := newServerTestApp(t)
	appCtx, stopApp := context.WithCancel(context.Background())
	a.ctx = appCtx

	started := make(chan struct{})
	finished := make(chan error, 1)
	a.HandleStream("shutdown", func(c *StreamConn) {
		close(started)
		_, err := c.Receive()
		finished <- err
	})

	srv := httptest.NewServer(http.HandlerFunc(a.serveStreamWS))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, _, err := websocket.Dial(ctx, wsURL(srv.URL)+"?name=shutdown", nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.CloseNow()
	<-started

	stopApp()
	select {
	case err := <-finished:
		if err != ErrStreamClosed {
			t.Fatalf("Receive after app shutdown = %v, want ErrStreamClosed", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("handler never returned after application context cancellation")
	}
}

// Returning from a handler is a graceful close. Data accepted immediately
// before that return must reach the peer before the socket closes.
func TestServerModeFlushesFinalFrameBeforeHandlerReturn(t *testing.T) {
	a := newServerTestApp(t)

	a.HandleStream("final", func(c *StreamConn) {
		if err := c.Send([]byte("final frame")); err != nil {
			t.Errorf("send: %v", err)
		}
	})

	srv := httptest.NewServer(http.HandlerFunc(a.serveStreamWS))
	defer srv.Close()

	for i := 0; i < 20; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		conn, _, err := websocket.Dial(ctx, wsURL(srv.URL)+"?name=final", nil)
		if err != nil {
			cancel()
			t.Fatalf("iteration %d dial: %v", i, err)
		}

		typ, data, err := conn.Read(ctx)
		if err != nil {
			_ = conn.CloseNow()
			cancel()
			t.Fatalf("iteration %d read: %v", i, err)
		}
		if typ != websocket.MessageBinary || string(data) != "final frame" {
			_ = conn.CloseNow()
			cancel()
			t.Fatalf("iteration %d frame = (%v, %q), want binary final frame", i, typ, data)
		}

		_, _, err = conn.Read(ctx)
		if status := websocket.CloseStatus(err); status != websocket.StatusNormalClosure {
			_ = conn.CloseNow()
			cancel()
			t.Fatalf("iteration %d close status = %v (err %v), want normal closure", i, status, err)
		}

		_ = conn.CloseNow()
		cancel()
	}
}

func wsURL(httpURL string) string {
	return "ws" + strings.TrimPrefix(httpURL, "http")
}
