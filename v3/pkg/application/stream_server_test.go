//go:build server

package application

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
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

func wsURL(httpURL string) string {
	return "ws" + strings.TrimPrefix(httpURL, "http")
}
