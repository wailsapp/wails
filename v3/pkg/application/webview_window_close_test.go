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
