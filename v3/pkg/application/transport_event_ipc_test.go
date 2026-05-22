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
