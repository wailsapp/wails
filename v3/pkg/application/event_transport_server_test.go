//go:build server

package application

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func TestServerEventDeliveryDoesNotQueueNativeJavaScript(t *testing.T) {
	a := newEventTransportTestApp(t, Options{})
	probe := &threadProbeApp{}
	probe.onMain.Store(true)
	a.impl = probe
	// Match hosts that still create a WebviewWindow in server mode. Its runtime
	// never becomes ready; it must not participate in server event delivery.
	window := &WebviewWindow{id: 1, impl: &stubWindowImpl{}}
	a.Window.Add(window)

	broadcaster := NewWebSocketBroadcaster(a)
	// This is the registration performed by serverApp.run, before accepting clients.
	a.wailsEventListeners = append(a.wailsEventListeners, broadcaster)
	server := httptest.NewServer(broadcaster)
	defer server.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	conn, _, err := websocket.Dial(ctx, wsURL(server.URL)+"?clientId=retention-test", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.CloseNow()
	conn.SetReadLimit(256 * 1024)

	// The upgrade completes before ServeHTTP necessarily finishes registration.
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()
	for {
		broadcaster.mu.RLock()
		registered := len(broadcaster.clients) == 1
		broadcaster.mu.RUnlock()
		if registered {
			break
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			t.Fatal("event WebSocket did not register")
		}
	}

	payload := strings.Repeat("x", 128*1024)
	for i := range 3 {
		name := fmt.Sprintf("event-%d", i)
		a.Event.dispatch(&CustomEvent{Name: name, Data: payload})
		var received CustomEvent
		if err := wsjson.Read(ctx, conn, &received); err != nil {
			t.Fatalf("reading %s: %v", name, err)
		}
		if received.Name != name || received.Data != payload {
			t.Fatalf("event %s: received name %q or changed payload", name, received.Name)
		}
	}
	window.pendingJSMutex.Lock()
	pending := len(window.pendingJS)
	window.pendingJSMutex.Unlock()
	if pending != 0 {
		t.Fatalf("retained %d native JavaScript strings", pending)
	}
	a.eventPayloads.mu.Lock()
	parked := a.eventPayloads.bytes
	a.eventPayloads.mu.Unlock()
	if parked != 0 {
		t.Fatalf("retained %d native payload bytes", parked)
	}

	conn.CloseNow()
	for {
		broadcaster.mu.RLock()
		remaining := len(broadcaster.clients)
		broadcaster.mu.RUnlock()
		if remaining == 0 {
			break
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			t.Fatal("event WebSocket did not unregister after close")
		}
	}
	// Events after disconnect must not enter native window queues.
	a.Event.dispatch(&CustomEvent{Name: "after-close", Data: payload})
	if len(window.pendingJS) != 0 || len(window.eventQueue) != 0 {
		t.Fatal("event after browser closure entered a native window queue")
	}
}
