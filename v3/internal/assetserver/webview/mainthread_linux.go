//go:build linux && cgo && !android

package webview

/*
#cgo linux pkg-config: glib-2.0

#include <glib.h>
#include <stdint.h>

extern void webviewMainThreadCallback(uintptr_t id);

typedef struct {
	uintptr_t id;
	GMutex    mutex;
	GCond     cond;
	gboolean  done;
} webviewMainSyncCall;

// webview_main_sync_trampoline runs on the GTK main thread (scheduled via
// g_main_context_invoke). It invokes the Go callback, then signals the waiting
// worker. g_cond_signal happens while the mutex is held and before the unlock,
// so the waiter cannot re-acquire the mutex (and destroy the primitives) until
// the signal has completed — making the stack-allocated GMutex/GCond safe.
static gboolean webview_main_sync_trampoline(gpointer data) {
	webviewMainSyncCall *call = (webviewMainSyncCall *)data;
	webviewMainThreadCallback(call->id);
	g_mutex_lock(&call->mutex);
	call->done = TRUE;
	g_cond_signal(&call->cond);
	g_mutex_unlock(&call->mutex);
	return G_SOURCE_REMOVE;
}

// webview_invoke_on_main_sync schedules the Go callback identified by id on the
// default GTK main context and blocks the calling thread until it has finished.
//
// WebKit2GTK objects may only be touched on the thread running the GTK main loop
// (g_application_run). Asset-server responses are produced on worker goroutines,
// so the WebKit calls that complete a request must hop here first. The wait is
// safe because webkit_uri_scheme_request_finish_with_response returns before the
// response stream is drained (WebKit reads it asynchronously), so the main loop
// never blocks waiting on the worker.
//
// If the caller is already the main thread, g_main_context_invoke runs the
// trampoline inline, so the wait completes immediately without deadlocking.
static void webview_invoke_on_main_sync(uintptr_t id) {
	webviewMainSyncCall call;
	call.id = id;
	call.done = FALSE;
	g_mutex_init(&call.mutex);
	g_cond_init(&call.cond);

	g_main_context_invoke(NULL, webview_main_sync_trampoline, &call);

	g_mutex_lock(&call.mutex);
	while (!call.done) {
		g_cond_wait(&call.cond, &call.mutex);
	}
	g_mutex_unlock(&call.mutex);

	g_mutex_clear(&call.mutex);
	g_cond_clear(&call.cond);
}
*/
import "C"

import (
	"sync"
)

var (
	mainSyncMu        sync.Mutex
	mainSyncNextID    uintptr
	mainSyncCallbacks = map[uintptr]func(){}
)

// invokeOnMainSync runs fn on the GTK main thread and blocks until it returns.
// It is safe to call from any goroutine, including the main thread itself.
func invokeOnMainSync(fn func()) {
	mainSyncMu.Lock()
	mainSyncNextID++
	id := mainSyncNextID
	mainSyncCallbacks[id] = fn
	mainSyncMu.Unlock()

	C.webview_invoke_on_main_sync(C.uintptr_t(id))
}

//export webviewMainThreadCallback
func webviewMainThreadCallback(id C.uintptr_t) {
	mainSyncMu.Lock()
	fn := mainSyncCallbacks[uintptr(id)]
	delete(mainSyncCallbacks, uintptr(id))
	mainSyncMu.Unlock()

	if fn != nil {
		fn()
	}
}
