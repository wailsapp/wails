package application

import (
	"sync"
	"testing"
	"time"
)

// lockProbeWindow is a minimal Window implementation that captures whether
// app.windowsLock could be exclusively acquired during the per-window
// dispatch. Embedding the Window interface gives nil-bodied placeholders
// for every other method; DispatchWailsEvent is the only one called by
// EventIPCTransport.DispatchWailsEvent, so the placeholders never run.
type lockProbeWindow struct {
	Window
	app                *App
	acquiredWriteLock  bool
}

func (w *lockProbeWindow) DispatchWailsEvent(_ *CustomEvent) {
	if w.app.windowsLock.TryLock() {
		w.acquiredWriteLock = true
		w.app.windowsLock.Unlock()
	}
}

// TestDispatchWailsEventReleasesLockBeforePerWindowDispatch is the regression
// test for the deadlock pattern from #5016 / #4424: holding windowsLock
// across the per-window DispatchWailsEvent call blocks any goroutine that
// needs the lock for write — including the main thread when ExecJS routes
// through it. The fix snapshots the windows under RLock then releases
// before iterating. A regression that re-acquires (or never releases)
// the lock for the iteration would prevent the probe window from taking
// the write lock and fail this test.
func TestDispatchWailsEventReleasesLockBeforePerWindowDispatch(t *testing.T) {
	app := &App{
		windows:     make(map[uint]Window),
		windowsLock: sync.RWMutex{},
	}
	probe := &lockProbeWindow{app: app}
	app.windows[1] = probe

	transport := &EventIPCTransport{app: app}
	transport.DispatchWailsEvent(&CustomEvent{Name: "test"})

	if !probe.acquiredWriteLock {
		t.Fatal("windowsLock was held during per-window DispatchWailsEvent; " +
			"the snapshot-and-release pattern in EventIPCTransport.DispatchWailsEvent has regressed")
	}

	// Sanity-check: the lock is also released after the whole transport call.
	acquired := make(chan struct{})
	go func() {
		app.windowsLock.Lock()
		close(acquired)
		app.windowsLock.Unlock()
	}()
	select {
	case <-acquired:
	case <-time.After(time.Second):
		t.Fatal("windowsLock should be released after DispatchWailsEvent completes")
	}
}

// recordingWindow counts the events delivered to it by EventIPCTransport.
type recordingWindow struct {
	Window
	id       uint
	received []string
}

func (w *recordingWindow) ID() uint { return w.id }

func (w *recordingWindow) DispatchWailsEvent(event *CustomEvent) {
	w.received = append(w.received, event.Name)
}

func newRecordingApp(windows ...*recordingWindow) *App {
	app := &App{windows: make(map[uint]Window)}
	for _, w := range windows {
		app.windows[w.id] = w
	}
	return app
}

func TestDispatchWailsEventBroadcastsUntargetedEvents(t *testing.T) {
	first := &recordingWindow{id: 1}
	second := &recordingWindow{id: 2}
	transport := &EventIPCTransport{app: newRecordingApp(first, second)}

	transport.DispatchWailsEvent(&CustomEvent{Name: "broadcast"})

	if len(first.received) != 1 || len(second.received) != 1 {
		t.Fatalf("untargeted event: first got %v, second got %v; want both to receive it",
			first.received, second.received)
	}
}

// TestDispatchWailsEventDeliversTargetedEventOnlyToTarget is the regression
// test for window.EmitEvent reaching every window instead of only its own.
func TestDispatchWailsEventDeliversTargetedEventOnlyToTarget(t *testing.T) {
	first := &recordingWindow{id: 1}
	second := &recordingWindow{id: 2}
	transport := &EventIPCTransport{app: newRecordingApp(first, second)}

	transport.DispatchWailsEvent(&CustomEvent{Name: "targeted", target: second})

	if len(first.received) != 0 {
		t.Fatalf("non-target window received %v", first.received)
	}
	if len(second.received) != 1 {
		t.Fatalf("target window received %v, want exactly one event", second.received)
	}
}

func TestDispatchWailsEventDropsEventForUnregisteredTarget(t *testing.T) {
	registered := &recordingWindow{id: 1}
	// Same ID as a registered window, but a different window instance: the
	// target must be matched by identity, not just by ID.
	stale := &recordingWindow{id: 1}
	closed := &recordingWindow{id: 2}
	transport := &EventIPCTransport{app: newRecordingApp(registered)}

	transport.DispatchWailsEvent(&CustomEvent{Name: "stale", target: stale})
	transport.DispatchWailsEvent(&CustomEvent{Name: "closed", target: closed})

	if len(registered.received) != 0 || len(stale.received) != 0 || len(closed.received) != 0 {
		t.Fatalf("event for an unregistered target was delivered: registered=%v stale=%v closed=%v",
			registered.received, stale.received, closed.received)
	}
}
