//go:build linux && cgo && !android

package webview

// This file provides cgo helpers for mainthread_linux_test.go. Go forbids cgo
// (an `import "C"` preamble) inside _test.go files, so the small GLib main-loop
// scaffolding the test needs lives here instead. None of it is referenced by
// production code paths.

/*
#cgo linux pkg-config: glib-2.0

#include <glib.h>

// All shared loop state is published and read under webview_test_mu so the
// worker goroutines that call the helpers below never observe it across threads
// without synchronization. webview_test_ready is signalled from inside the loop
// (via an idle source) so waiters block on the cond instead of spinning on an
// unsynchronized flag.
static GMutex     webview_test_mu;
static GCond      webview_test_cond;
static GMainLoop *webview_test_loop;        // guarded by webview_test_mu
static GThread   *webview_test_loop_thread; // guarded by webview_test_mu
static gboolean   webview_test_ready;       // guarded by webview_test_mu

// webview_test_mark_ready runs on the loop thread once the loop starts
// iterating, publishing that it is live.
static gboolean webview_test_mark_ready(gpointer data) {
	g_mutex_lock(&webview_test_mu);
	webview_test_ready = TRUE;
	g_cond_signal(&webview_test_cond);
	g_mutex_unlock(&webview_test_mu);
	return G_SOURCE_REMOVE;
}

// webview_test_run_loop records the running thread and drives the default GLib
// main context, mimicking the GTK main loop that invokeOnMainSync dispatches to.
static void webview_test_run_loop(void) {
	GMainLoop *loop = g_main_loop_new(NULL, FALSE);

	g_mutex_lock(&webview_test_mu);
	webview_test_loop = loop;
	webview_test_loop_thread = g_thread_self();
	g_mutex_unlock(&webview_test_mu);

	g_idle_add(webview_test_mark_ready, NULL);
	g_main_loop_run(loop);

	g_mutex_lock(&webview_test_mu);
	webview_test_loop = NULL;
	webview_test_ready = FALSE;
	g_mutex_unlock(&webview_test_mu);

	g_main_loop_unref(loop);
}

static void webview_test_wait_running(void) {
	g_mutex_lock(&webview_test_mu);
	while (!webview_test_ready) {
		g_cond_wait(&webview_test_cond, &webview_test_mu);
	}
	g_mutex_unlock(&webview_test_mu);
}

static void webview_test_quit_loop(void) {
	g_mutex_lock(&webview_test_mu);
	GMainLoop *loop = webview_test_loop;
	g_mutex_unlock(&webview_test_mu);
	if (loop != NULL) {
		g_main_loop_quit(loop);
	}
}

// webview_test_on_loop_thread reports whether the caller is the thread that runs
// the main loop.
static int webview_test_on_loop_thread(void) {
	g_mutex_lock(&webview_test_mu);
	GThread *loopThread = webview_test_loop_thread;
	g_mutex_unlock(&webview_test_mu);
	return g_thread_self() == loopThread ? 1 : 0;
}
*/
import "C"

func testRunMainLoop()     { C.webview_test_run_loop() }
func testWaitLoopRunning() { C.webview_test_wait_running() }
func testQuitMainLoop()    { C.webview_test_quit_loop() }
func testOnLoopThread() bool {
	return C.webview_test_on_loop_thread() != 0
}
