//go:build linux && !wailscgo && !android

package webview

import (
	"context"
	"github.com/ebitengine/purego"
	"sync"
	"unsafe"
)

// requestLifetime retains the exact native request without hiding WebKit's last
// reference. WebKit drops that reference when the browser cancels the request,
// or when its response finishes. Unlike URI matching, this also distinguishes
// simultaneous requests for the same URL.
//
// The callback handle holds only a CancelFunc, not the Go request, so it cannot
// keep the request finalizer alive. All native reference operations must run on
// the GTK main thread; WebKit invokes NewRequest there, and close dispatches back.
// This mechanism is shared by GTK4/WebKitGTK 6.0 and GTK3/WebKit2GTK 4.1.
type requestLifetime struct {
	context.Context
	cancel context.CancelFunc
	object unsafe.Pointer
	handle uintptr
	once   sync.Once
	state  struct {
		sync.Mutex
		nativeReleased bool
		handleDeleted  bool
	}
}

var activeRequestLifetimes = struct {
	sync.Mutex
	values map[*requestLifetime]struct{}
}{
	values: make(map[*requestLifetime]struct{}),
}

func retainRequest(object unsafe.Pointer) *requestLifetime {
	ctx, cancel := context.WithCancel(context.Background())
	lifetime := &requestLifetime{
		Context: ctx,
		cancel:  cancel,
		object:  object,
		handle:  newRequestCancelHandle(cancel),
	}
	activeRequestLifetimes.Lock()
	activeRequestLifetimes.values[lifetime] = struct{}{}
	activeRequestLifetimes.Unlock()
	g_object_add_toggle_ref(uintptr(object), requestToggleNotify, lifetime.handle)
	return lifetime
}

func (l *requestLifetime) close() {
	l.once.Do(func() {
		if !invokeOnMainSyncIfEnabled(l.releaseNative) {
			// The GTK loop is gone, so no native callback can arrive after this
			// point. Do not release the GObject from the worker; process shutdown
			// will reclaim it with the rest of the WebKit instance.
			l.abandon()
		}
	})
}

// CloseActiveRequests releases all native request references while the GTK
// thread is still active. The application calls this immediately before it
// quits the main loop, covering handlers that outlive window destruction.
func CloseActiveRequests() {
	activeRequestLifetimes.Lock()
	lifetimes := make([]*requestLifetime, 0, len(activeRequestLifetimes.values))
	for lifetime := range activeRequestLifetimes.values {
		lifetimes = append(lifetimes, lifetime)
	}
	activeRequestLifetimes.Unlock()

	for _, lifetime := range lifetimes {
		lifetime.releaseNative()
	}
}

func (l *requestLifetime) releaseNative() {
	l.state.Lock()
	if l.state.nativeReleased {
		l.state.Unlock()
		return
	}
	l.state.nativeReleased = true
	l.state.Unlock()

	// Remove the native callback before deleting its Go handle. Doing both on
	// the GTK thread serializes this with WebKit's ref changes.
	g_object_remove_toggle_ref(uintptr(l.object), requestToggleNotify, l.handle)
	l.cancel()
	l.deleteHandle()
	l.removeActive()
}

func (l *requestLifetime) abandon() {
	l.state.Lock()
	if l.state.nativeReleased {
		l.state.Unlock()
		return
	}
	l.state.nativeReleased = true
	l.state.Unlock()
	l.cancel()
	// Keep the cancellation handle alive because the native toggle reference remains
	// attached until WebKit tears down the object with the application. There is
	// no live GTK loop left on which it could be removed safely.
	l.removeActive()
}

func (l *requestLifetime) deleteHandle() {
	l.state.Lock()
	if l.state.handleDeleted {
		l.state.Unlock()
		return
	}
	l.state.handleDeleted = true
	l.state.Unlock()
	requestCancelHandles.Lock()
	delete(requestCancelHandles.values, l.handle)
	requestCancelHandles.Unlock()
}

func (l *requestLifetime) removeActive() {
	activeRequestLifetimes.Lock()
	delete(activeRequestLifetimes.values, l)
	activeRequestLifetimes.Unlock()
}

// Only the cancellation function is retained, so handles cannot keep the Go
// request finalizer alive. A single trampoline serves every native request.
var requestCancelHandles = struct {
	sync.Mutex
	next   uintptr
	values map[uintptr]context.CancelFunc
}{values: make(map[uintptr]context.CancelFunc)}

func newRequestCancelHandle(cancel context.CancelFunc) uintptr {
	requestCancelHandles.Lock()
	defer requestCancelHandles.Unlock()
	requestCancelHandles.next++
	id := requestCancelHandles.next
	requestCancelHandles.values[id] = cancel
	return id
}

var requestToggleNotify = purego.NewCallback(func(data, object uintptr, last int32) {
	if last == 0 {
		return
	}
	requestCancelHandles.Lock()
	cancel := requestCancelHandles.values[data]
	requestCancelHandles.Unlock()
	if cancel != nil {
		cancel()
	}
})
