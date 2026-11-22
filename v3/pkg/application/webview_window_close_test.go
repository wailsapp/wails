package application

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/events"
)

type closeTestApp struct{ platformApp }

func (*closeTestApp) isOnMainThread() bool { return true }

type closeTestWindowImpl struct {
	webviewWindowImpl
	closes atomic.Int32
}

func (w *closeTestWindowImpl) close() {
	w.closes.Add(1)
}

func newCloseTestWindow(t *testing.T) (*WebviewWindow, *closeTestWindowImpl) {
	t.Helper()

	previous := globalApplication
	app := &App{impl: &closeTestApp{}, windows: make(map[uint]Window)}
	app.Window = newWindowManager(app)
	globalApplication = app
	t.Cleanup(func() { globalApplication = previous })

	window := NewWindow(WebviewWindowOptions{})
	impl := &closeTestWindowImpl{}
	window.impl = impl
	app.Window.Add(window)
	return window, impl
}

func TestOverlappingWindowClosingEventsCloseNativeWindowOnce(t *testing.T) {
	window, impl := newCloseTestWindow(t)
	listener := window.eventListeners[uint(events.Common.WindowClosing)][0]

	var done sync.WaitGroup
	for range 2 {
		done.Add(1)
		go func() {
			defer done.Done()
			listener.callback(NewWindowEvent())
		}()
	}
	done.Wait()

	if got := impl.closes.Load(); got != 1 {
		t.Fatalf("native close called %d times, want 1", got)
	}
	if !window.isDestroyed() {
		t.Fatal("window was not marked destroyed")
	}
}

func TestCancelledWindowClosingAllowsLaterClose(t *testing.T) {
	window, impl := newCloseTestWindow(t)
	listener := window.eventListeners[uint(events.Common.WindowClosing)][0]
	completed := make(chan struct{})
	original := listener.callback
	listener.callback = func(event *WindowEvent) {
		defer close(completed)
		original(event)
	}
	cancelHook := window.RegisterHook(events.Common.WindowClosing, func(event *WindowEvent) {
		event.Cancel()
	})

	window.HandleWindowEvent(uint(events.Common.WindowClosing))
	if window.isDestroyed() || impl.closes.Load() != 0 {
		t.Fatal("cancelled close changed window state")
	}

	cancelHook()
	window.HandleWindowEvent(uint(events.Common.WindowClosing))
	<-completed
	if got := impl.closes.Load(); got != 1 {
		t.Fatalf("native close called %d times, want 1", got)
	}
}

// newCleanupTestApp is newCloseTestWindow's app, with the fields cleanup reads.
func newCleanupTestApp(t *testing.T, windows int) (*App, []*closeTestWindowImpl, []*WebviewWindow) {
	t.Helper()
	first, firstImpl := newCloseTestWindow(t)
	app := globalApplication
	app.cancel = func() {}
	impls := []*closeTestWindowImpl{firstImpl}
	wins := []*WebviewWindow{first}
	for i := 1; i < windows; i++ {
		w := NewWindow(WebviewWindowOptions{})
		impl := &closeTestWindowImpl{}
		w.impl = impl
		app.Window.Add(w)
		impls = append(impls, impl)
		wins = append(wins, w)
	}
	return app, impls, wins
}

// App.cleanup runs on quit. It must destroy every window itself, on the spot:
// a Close() only queues a WindowClosing event, and cleanup empties the window
// map straight after, so the queued event finds no window and nothing is
// destroyed (no WebView2 shutdown on Windows).
func TestCleanupDestroysEveryWindowSynchronously(t *testing.T) {
	app, impls, wins := newCleanupTestApp(t, 2)

	app.cleanup()

	for i, impl := range impls {
		if got := impl.closes.Load(); got != 1 {
			t.Errorf("window %d: native close called %d times, want 1", i, got)
		}
		if !wins[i].isDestroyed() {
			t.Errorf("window %d was not marked destroyed", i)
		}
	}
}

// A hook that hides to tray and cancels (the usual close-to-tray hook) must not
// keep a window alive through app teardown: the app is quitting, not the user
// closing one window.
func TestCleanupBypassesWindowClosingHooks(t *testing.T) {
	app, impls, wins := newCleanupTestApp(t, 1)
	wins[0].RegisterHook(events.Common.WindowClosing, func(event *WindowEvent) {
		event.Cancel()
	})

	app.cleanup()

	if got := impls[0].closes.Load(); got != 1 {
		t.Fatalf("native close called %d times, want 1 - a cancelling hook kept the window alive", got)
	}
}
