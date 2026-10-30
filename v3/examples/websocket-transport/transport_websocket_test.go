package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func connectTransport(t testing.TB) (*WebSocketTransport, *websocket.Conn, chan *WebSocketMessage) {
	t.Helper()
	transport := NewWebSocketTransport("")
	done := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer close(done)
		transport.handleWebSocket(w, r)
	}))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, _, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(server.URL, "http"), &websocket.DialOptions{
		HTTPHeader: http.Header{"Origin": {"https://example.invalid"}},
	})
	if err != nil {
		server.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		conn.CloseNow()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Error("transport handler did not stop after disconnect")
		}
		server.Close()
	})
	for {
		transport.mu.RLock()
		for _, messages := range transport.clients {
			transport.mu.RUnlock()
			return transport, conn, messages
		}
		transport.mu.RUnlock()
		select {
		case <-ctx.Done():
			t.Fatal("connection was not registered")
		case <-time.After(time.Millisecond):
		}
	}
}

func TestTransportJSONResponses(t *testing.T) {
	transport, conn, messages := connectTransport(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, tc := range []struct {
		name   string
		err    error
		status int
		data   string
	}{
		{"success", nil, 200, "hello"},
		{"error", fmt.Errorf("failed"), 422, "failed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transport.sendResponse(ctx, messages, tc.name, "hello", tc.err)
			var response WebSocketMessage
			if err := wsjson.Read(ctx, conn, &response); err != nil {
				t.Fatal(err)
			}
			if response.ID != tc.name || response.Type != "response" || response.Response == nil {
				t.Fatalf("unexpected response envelope: %+v", response)
			}
			if response.Response.StatusCode != tc.status || response.Response.Data != tc.data {
				t.Fatalf("unexpected response body: %+v", response.Response)
			}
		})
	}
}

func TestTransportLargeMessageAndEvent(t *testing.T) {
	transport, conn, _ := connectTransport(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := make(chan error, 1)
	go func() {
		var msg WebSocketMessage
		err := wsjson.Read(ctx, conn, &msg)
		if err == nil && (msg.Type != "event" || msg.Event == nil || msg.Event.Name != "ready") {
			err = fmt.Errorf("unexpected event: %+v", msg)
		}
		result <- err
	}()
	// The old transport accepted messages larger than coder's default 32 KiB.
	if err := wsjson.Write(ctx, conn, WebSocketMessage{Type: "ignored", ID: strings.Repeat("x", 1<<20)}); err != nil {
		t.Fatal(err)
	}
	// A pong confirms the server read through the large message and stayed open.
	if err := conn.Ping(ctx); err != nil {
		t.Fatal(err)
	}
	transport.DispatchWailsEvent(&application.CustomEvent{Name: "ready"})
	if err := <-result; err != nil {
		t.Fatal(err)
	}
}

func TestTransportStopClosesConnections(t *testing.T) {
	transport, conn, _ := connectTransport(t)
	transport.server = &http.Server{}
	stopped := make(chan error, 1)
	go func() { stopped <- transport.Stop() }()
	select {
	case err := <-stopped:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Stop did not close active connections promptly")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, _, err := conn.Read(ctx); err == nil {
		t.Fatal("connection stayed open after Stop")
	}
}

func BenchmarkTransportResponse(b *testing.B) {
	transport, conn, messages := connectTransport(b)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		transport.sendResponse(ctx, messages, "benchmark", "hello", nil)
		var response WebSocketMessage
		if err := wsjson.Read(ctx, conn, &response); err != nil {
			b.Fatal(err)
		}
	}
}
