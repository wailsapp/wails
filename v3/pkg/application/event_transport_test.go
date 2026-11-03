package application

import "testing"

// newEventTransportTestApp exercises the same transport selection as a real
// application without starting a platform event loop or creating native windows.
func newEventTransportTestApp(t *testing.T, options Options) *App {
	t.Helper()
	previous := globalApplication
	globalApplication = nil
	a := New(options)
	t.Cleanup(func() {
		a.cancel()
		for _, shutdown := range a.shutdownTasks {
			shutdown()
		}
		globalApplication = previous
	})
	return a
}

func TestDefaultEventTransportMatchesBuildMode(t *testing.T) {
	a := newEventTransportTestApp(t, Options{})
	if isServerBuild {
		if len(a.wailsEventListeners) != 0 {
			t.Fatalf("server constructor installed %d listeners before its WebSocket broadcaster", len(a.wailsEventListeners))
		}
		return
	}
	if len(a.wailsEventListeners) != 1 {
		t.Fatalf("native constructor installed %d listeners, want one IPC listener", len(a.wailsEventListeners))
	}
	ipc, ok := a.wailsEventListeners[0].(*EventIPCTransport)
	if !ok || ipc.app != a {
		t.Fatalf("native listener = %#v, want IPC for this application", a.wailsEventListeners[0])
	}
}

type recordingEventTransport struct {
	*HTTPTransport
	events []*CustomEvent
}

func (transport *recordingEventTransport) DispatchWailsEvent(event *CustomEvent) {
	transport.events = append(transport.events, event)
}

func TestCustomEventTransportRemainsAuthoritative(t *testing.T) {
	transport := &recordingEventTransport{HTTPTransport: NewHTTPTransport()}
	a := newEventTransportTestApp(t, Options{Transport: transport})
	if len(a.wailsEventListeners) != 1 || a.wailsEventListeners[0] != transport {
		t.Fatalf("listeners = %#v, want the supplied event transport only", a.wailsEventListeners)
	}
	first := &CustomEvent{Name: "first"}
	second := &CustomEvent{Name: "second"}
	a.Event.dispatch(first)
	a.Event.dispatch(second)
	if len(transport.events) != 2 || transport.events[0] != first || transport.events[1] != second {
		t.Fatalf("custom transport received %#v, want both events once in order", transport.events)
	}
}
