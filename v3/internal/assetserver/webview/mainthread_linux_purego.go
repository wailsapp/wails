//go:build linux && purego && !android

package webview

// CGO-free (purego) port of mainthread_linux.go.
//
// The C implementation used a statically allocated GMutex to protect the
// dispatch-enabled flag and a per-call GMutex/GCond pair to block the worker;
// here a Go sync.Mutex plays the former role and a per-call done channel the
// latter. The trampoline scheduled onto the GLib main context is a purego
// callback created exactly once (purego callback slots are a process-wide,
// never-freed resource, so a per-call NewCallback would leak them).

import (
	"fmt"
	"sync"

	"github.com/ebitengine/purego"
)

// ----------------------------------------------------------------------------
// Library loading helpers (shared with the other *_linux_purego.go files)
// ----------------------------------------------------------------------------

// dlopenWebviewLib loads the first of the given soname candidates. The
// versioned name is tried first; the unversioned name usually only exists when
// -dev packages are installed. In practice pkg/application has already loaded
// and validated these same libraries before any request can arrive, so a
// failure here is a safety net, not primary UX — panic with a clear message.
func dlopenWebviewLib(names ...string) uintptr {
	var lastErr error
	for _, name := range names {
		handle, err := purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err == nil && handle != 0 {
			return handle
		}
		if err != nil {
			lastErr = err
		}
	}
	panic(fmt.Sprintf("wails: assetserver/webview: failed to load required library %s: %v", names[0], lastErr))
}

// mustRegisterWebviewFunc binds fptr to the named symbol in lib, panicking if
// the symbol is missing (same rationale as dlopenWebviewLib: pkg/application
// validated the installed library versions before any request can arrive).
func mustRegisterWebviewFunc(fptr any, lib uintptr, name string) {
	sym, err := purego.Dlsym(lib, name)
	if err != nil || sym == 0 {
		panic(fmt.Sprintf("wails: assetserver/webview: failed to resolve required symbol %s: %v", name, err))
	}
	purego.RegisterFunc(fptr, sym)
}

var (
	webviewGLibOnce sync.Once
	webviewLibGLib  uintptr

	// void g_main_context_invoke(GMainContext *context, GSourceFunc function,
	// gpointer data). Pass 0 for the context to target the default (GTK) main
	// context.
	g_main_context_invoke func(context uintptr, function uintptr, data uintptr)
)

// ensureGLib lazily dlopens libglib-2.0 and binds the symbols this file needs.
// Lazy (rather than init-time) so that merely importing the package cannot
// panic on a system without GLib.
func ensureGLib() {
	webviewGLibOnce.Do(func() {
		webviewLibGLib = dlopenWebviewLib("libglib-2.0.so.0", "libglib-2.0.so")
		mustRegisterWebviewFunc(&g_main_context_invoke, webviewLibGLib, "g_main_context_invoke")
	})
}

// ----------------------------------------------------------------------------
// Main-thread dispatch
// ----------------------------------------------------------------------------

// webviewDispatchMu serializes the enabled-check plus scheduling in
// invokeOnMainSync against the flag clear in DisableMainThreadDispatch.
// Without it a worker could read webviewMainDispatchEnabled == true, then have
// the flag flipped before it scheduled its source, and still queue work onto
// the now-dead main loop — blocking forever.
var webviewDispatchMu sync.Mutex

// webviewMainDispatchEnabled gates whether invokeOnMainSync may schedule work
// onto the GTK main loop. It starts enabled and is cleared once the loop has
// stopped (see DisableMainThreadDispatch). It is only ever read or written
// while holding webviewDispatchMu.
var webviewMainDispatchEnabled = true

type mainSyncCall struct {
	fn   func()
	done chan struct{}
}

var (
	mainSyncMu     sync.Mutex
	mainSyncNextID uintptr
	mainSyncCalls  = map[uintptr]*mainSyncCall{}
)

// mainSyncTrampoline runs on the GTK main thread (scheduled via
// g_main_context_invoke). It invokes the Go callback identified by data, then
// signals the waiting worker: runMainSyncCallback closes the call's done
// channel only after the callback has finished, so the waiter cannot proceed
// (and drop the call record) until the work has completed — the Go analogue of
// the C version signalling the per-call GCond while holding its GMutex.
var mainSyncTrampoline = purego.NewCallback(func(data uintptr) uintptr {
	runMainSyncCallback(data)
	return 0 // G_SOURCE_REMOVE
})

// runMainSyncCallback is the Go twin of the cgo webviewMainThreadCallback
// export: it looks up the pending call by id, removes it from the registry,
// runs it, and releases the waiting worker.
func runMainSyncCallback(id uintptr) {
	mainSyncMu.Lock()
	call := mainSyncCalls[id]
	delete(mainSyncCalls, id)
	mainSyncMu.Unlock()

	if call == nil {
		return
	}
	if call.fn != nil {
		call.fn()
	}
	close(call.done)
}

// invokeOnMainSync runs fn on the GTK main thread and blocks until it returns.
// It is safe to call from any goroutine, including the main thread itself.
//
// WebKit2GTK objects may only be touched on the thread running the GTK main
// loop (g_application_run). Asset-server responses are produced on worker
// goroutines, so the WebKit calls that complete a request must hop here first.
// The wait is safe because webkit_uri_scheme_request_finish_with_response
// returns before the response stream is drained (WebKit reads it
// asynchronously), so the main loop never blocks waiting on the worker.
//
// If the caller is already the main thread, g_main_context_invoke runs the
// trampoline inline, so the done channel is closed before the wait begins and
// the receive completes immediately without deadlocking.
func invokeOnMainSync(fn func()) {
	ensureGLib()

	mainSyncMu.Lock()
	mainSyncNextID++
	id := mainSyncNextID
	call := &mainSyncCall{fn: fn, done: make(chan struct{})}
	mainSyncCalls[id] = call
	mainSyncMu.Unlock()

	// The enabled-check and the g_main_context_invoke that acts on it must be
	// atomic with respect to DisableMainThreadDispatch. Holding
	// webviewDispatchMu across both means a worker either schedules onto a live
	// loop or sees the loop already stopped — it can never schedule onto a loop
	// that stops in between (which would block it here forever). The trampoline
	// only touches the per-call registry and done channel, never
	// webviewDispatchMu, so when g_main_context_invoke runs it inline
	// (main-thread caller) there is no self-deadlock.
	webviewDispatchMu.Lock()
	if !webviewMainDispatchEnabled {
		// The GTK main loop has stopped: a scheduled source would never run. The
		// loop is no longer iterating, so the cross-thread race that makes
		// main-thread confinement necessary is gone — running the callback inline
		// on the worker lets in-flight asset requests drain during shutdown
		// instead of wedging. See #5631 (review question 5).
		webviewDispatchMu.Unlock()
		runMainSyncCallback(id)
		return
	}
	g_main_context_invoke(0, mainSyncTrampoline, id)
	webviewDispatchMu.Unlock()

	<-call.done
}

// DisableMainThreadDispatch marks the GTK main loop as stopped. After it is
// called, invokeOnMainSync runs callbacks inline on the calling goroutine
// instead of scheduling them onto the now-dead main loop, so asset-server
// workers that complete a request during shutdown cannot block forever waiting
// for a source that will never be serviced. The application layer calls this
// once g_application_run has returned. See issue #5631.
func DisableMainThreadDispatch() {
	webviewDispatchMu.Lock()
	webviewMainDispatchEnabled = false
	webviewDispatchMu.Unlock()
}
