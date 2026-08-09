//go:build server

package application

import (
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/wailsapp/wails/v3/internal/mailbox"
)

func TestWebSocketBroadcasterPreservesOrderPerClient(t *testing.T) {
	firstStarted := make(chan struct{})
	releaseFirst := make(chan struct{})
	firstClientEvents := make(chan int, 2)
	secondClientEvents := make(chan int, 2)

	firstClient := &clientInfo{
		events: mailbox.New(func(event *CustomEvent) {
			value := event.Data.(int)
			if value == 1 {
				close(firstStarted)
				<-releaseFirst
			}
			firstClientEvents <- value
		}),
	}
	secondClient := &clientInfo{
		events: mailbox.New(func(event *CustomEvent) {
			secondClientEvents <- event.Data.(int)
		}),
	}

	broadcaster := &WebSocketBroadcaster{
		clients: map[*websocket.Conn]*clientInfo{
			new(websocket.Conn): firstClient,
			new(websocket.Conn): secondClient,
		},
	}

	broadcaster.DispatchWailsEvent(&CustomEvent{Data: 1})
	<-firstStarted
	broadcaster.DispatchWailsEvent(&CustomEvent{Data: 2})

	if got := <-secondClientEvents; got != 1 {
		t.Fatalf("second client received %d first, want 1", got)
	}
	if got := <-secondClientEvents; got != 2 {
		t.Fatalf("second client received %d second, want 2", got)
	}

	select {
	case got := <-firstClientEvents:
		t.Fatalf("first client received %d while its first event was blocked", got)
	case <-time.After(100 * time.Millisecond):
	}

	close(releaseFirst)
	if got := <-firstClientEvents; got != 1 {
		t.Fatalf("first client received %d first, want 1", got)
	}
	if got := <-firstClientEvents; got != 2 {
		t.Fatalf("first client received %d second, want 2", got)
	}
}

func TestWebSocketBroadcasterPreservesBurstOrder(t *testing.T) {
	const eventCount = 1000

	events := make(chan int, eventCount)
	client := &clientInfo{
		events: mailbox.New(func(event *CustomEvent) {
			events <- event.Data.(int)
		}),
	}
	broadcaster := &WebSocketBroadcaster{
		clients: map[*websocket.Conn]*clientInfo{new(websocket.Conn): client},
	}

	for value := range eventCount {
		broadcaster.DispatchWailsEvent(&CustomEvent{Data: value})
	}
	for want := range eventCount {
		if got := <-events; got != want {
			t.Fatalf("event %d received as %d", want, got)
		}
	}
}

func TestWebSocketBroadcasterDoesNotDispatchAfterClientRemoval(t *testing.T) {
	conn := new(websocket.Conn)
	events := make(chan *CustomEvent, 1)
	window := NewBrowserWindow(1, "client")
	client := &clientInfo{
		window: window,
		events: mailbox.New(func(event *CustomEvent) {
			events <- event
		}),
	}
	broadcaster := &WebSocketBroadcaster{
		clients: map[*websocket.Conn]*clientInfo{conn: client},
		windows: map[string]*BrowserWindow{"client": window},
	}

	removed, remaining := broadcaster.removeClient(conn, "client")
	if removed != client {
		t.Fatal("removeClient returned the wrong client")
	}
	if remaining != 0 {
		t.Fatalf("%d clients remain, want 0", remaining)
	}
	if broadcaster.GetBrowserWindow("client") != nil {
		t.Fatal("browser window was not removed")
	}

	broadcaster.DispatchWailsEvent(&CustomEvent{Name: "test"})
	select {
	case <-events:
		t.Fatal("removed client received an event")
	case <-time.After(100 * time.Millisecond):
	}
}

func TestWebSocketBroadcasterKeepsReplacementWindowOnStaleClientRemoval(t *testing.T) {
	conn := new(websocket.Conn)
	staleWindow := NewBrowserWindow(1, "client")
	replacementWindow := NewBrowserWindow(2, "client")
	broadcaster := &WebSocketBroadcaster{
		clients: map[*websocket.Conn]*clientInfo{
			conn: {window: staleWindow},
		},
		windows: map[string]*BrowserWindow{"client": replacementWindow},
	}

	broadcaster.removeClient(conn, "client")

	if got := broadcaster.GetBrowserWindow("client"); got != replacementWindow {
		t.Fatal("stale client removal removed the replacement browser window")
	}
}
