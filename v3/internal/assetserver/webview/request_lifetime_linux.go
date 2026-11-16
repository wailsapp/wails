//go:build linux && cgo && !purego && !android

package webview

/*
#cgo linux pkg-config: gobject-2.0
#include <glib-object.h>
#include <stdint.h>

extern void webviewRequestReleased(uintptr_t handle);

static void request_toggle_notify(gpointer data, GObject *object, gboolean last) {
	if (last) {
		webviewRequestReleased((uintptr_t)data);
	}
}

static void retain_request(GObject *request, uintptr_t handle) {
	g_object_add_toggle_ref(request, request_toggle_notify, (gpointer)handle);
}

static void release_request(GObject *request, uintptr_t handle) {
	g_object_remove_toggle_ref(request, request_toggle_notify, (gpointer)handle);
}
*/
import "C"

import (
	"context"
	"runtime/cgo"
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
	handle cgo.Handle
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
		handle:  cgo.NewHandle(cancel),
	}
	activeRequestLifetimes.Lock()
	activeRequestLifetimes.values[lifetime] = struct{}{}
	activeRequestLifetimes.Unlock()
	C.retain_request((*C.GObject)(object), C.uintptr_t(lifetime.handle))
	return lifetime
}

// releaseRequestRef removes the toggle reference that retainRequest added.
func releaseRequestRef(object unsafe.Pointer, handle uintptr) {
	C.release_request((*C.GObject)(object), C.uintptr_t(handle))
}

func (l *requestLifetime) deleteHandle() {
	l.state.Lock()
	if l.state.handleDeleted {
		l.state.Unlock()
		return
	}
	l.state.handleDeleted = true
	l.state.Unlock()
	l.handle.Delete()
}

//export webviewRequestReleased
func webviewRequestReleased(handle C.uintptr_t) {
	cgo.Handle(handle).Value().(context.CancelFunc)()
}
