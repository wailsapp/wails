//go:build linux && cgo && !android

package webview

/*
#cgo linux pkg-config: glib-2.0

#include <glib.h>
#include <stdint.h>

extern void webviewMainThreadCallback(uintptr_t id);

// webview_dispatch_mu serializes the enabled-check plus scheduling in
// webview_invoke_on_main_sync against the flag clear in
// webview_disable_main_dispatch. Without it a worker could read
// enabled == TRUE, then have the flag flipped before it scheduled its source,
// and still queue work onto the now-dead main loop — blocking forever. A
// statically allocated GMutex needs no g_mutex_init.
static GMutex webview_dispatch_mu;

// webview_main_dispatch_enabled gates whether webview_invoke_on_main_sync may
// schedule work onto the GTK main loop. It starts enabled and is cleared once
// the loop has stopped (see webview_disable_main_dispatch). It is only ever read
// or written while holding webview_dispatch_mu.
static gboolean webview_main_dispatch_enabled = TRUE;

static void webview_disable_main_dispatch(void) {
	g_mutex_lock(&webview_dispatch_mu);
	webview_main_dispatch_enabled = FALSE;
	g_mutex_unlock(&webview_dispatch_mu);
}

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
// When run_if_disabled is false, it returns FALSE instead of running the
// callback inline after the GTK loop has stopped. That is required for native
// object releases that must never run on a worker goroutine during shutdown.
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
static gboolean webview_invoke_on_main_sync(uintptr_t id, gboolean run_if_disabled) {
	// The enabled-check and the g_main_context_invoke that acts on it must be
	// atomic with respect to webview_disable_main_dispatch. Holding
	// webview_dispatch_mu across both means a worker either schedules onto a live
	// loop or sees the loop already stopped — it can never schedule onto a loop
	// that stops in between (which would block it here forever). The trampoline
	// only touches the per-call mutex/cond, never webview_dispatch_mu, so when
	// g_main_context_invoke runs it inline (main-thread caller) there is no
	// self-deadlock.
	g_mutex_lock(&webview_dispatch_mu);
	if (!webview_main_dispatch_enabled) {
		g_mutex_unlock(&webview_dispatch_mu);
		if (!run_if_disabled) {
			return FALSE;
		}
		// The GTK main loop has stopped: a scheduled source would never run. The
		// loop is no longer iterating, so the cross-thread race that makes
		// main-thread confinement necessary is gone — running the callback inline
		// on the worker lets in-flight asset requests drain during shutdown
		// instead of wedging. See #5631 (review question 5).
		webviewMainThreadCallback(id);
		return TRUE;
	}

	webviewMainSyncCall call;
	call.id = id;
	call.done = FALSE;
	g_mutex_init(&call.mutex);
	g_cond_init(&call.cond);
	g_main_context_invoke(NULL, webview_main_sync_trampoline, &call);
	g_mutex_unlock(&webview_dispatch_mu);

	g_mutex_lock(&call.mutex);
	while (!call.done) {
		g_cond_wait(&call.cond, &call.mutex);
	}
	g_mutex_unlock(&call.mutex);

	g_mutex_clear(&call.mutex);
	g_cond_clear(&call.cond);
	return TRUE;
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

	C.webview_invoke_on_main_sync(C.uintptr_t(id), C.gboolean(1))
}

// invokeOnMainSyncIfEnabled runs fn on the GTK main thread while the main loop
// is active. It returns false after shutdown has disabled dispatch and leaves
// fn uncalled, allowing callers that release native objects to avoid touching
// them from a worker goroutine.
func invokeOnMainSyncIfEnabled(fn func()) bool {
	mainSyncMu.Lock()
	mainSyncNextID++
	id := mainSyncNextID
	mainSyncCallbacks[id] = fn
	mainSyncMu.Unlock()

	if C.webview_invoke_on_main_sync(C.uintptr_t(id), C.gboolean(0)) == 0 {
		mainSyncMu.Lock()
		delete(mainSyncCallbacks, id)
		mainSyncMu.Unlock()
		return false
	}
	return true
}

// DisableMainThreadDispatch marks the GTK main loop as stopped. After it is
// called, invokeOnMainSync runs callbacks inline on the calling goroutine
// instead of scheduling them onto the now-dead main loop, so asset-server
// workers that complete a request during shutdown cannot block forever waiting
// for a source that will never be serviced. The application layer calls this
// once g_application_run has returned. See issue #5631.
func DisableMainThreadDispatch() {
	C.webview_disable_main_dispatch()
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
