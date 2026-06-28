//go:build linux && cgo && !android

package webview

// This file provides cgo helpers for mainthread_linux_test.go. Go forbids cgo
// (an `import "C"` preamble) inside _test.go files, so the small GLib main-loop
// scaffolding the test needs lives here instead. None of it is referenced by
// production code paths.

/*
#cgo linux pkg-config: glib-2.0

#include <glib.h>

static GMainLoop *webview_test_loop;
static gpointer   webview_test_loop_thread;

// webview_test_run_loop records the running thread and drives the default GLib
// main context, mimicking the GTK main loop that invokeOnMainSync dispatches to.
static void webview_test_run_loop(void) {
	webview_test_loop_thread = (gpointer)g_thread_self();
	webview_test_loop = g_main_loop_new(NULL, FALSE);
	g_main_loop_run(webview_test_loop);
	g_main_loop_unref(webview_test_loop);
	webview_test_loop = NULL;
}

static void webview_test_wait_running(void) {
	while (webview_test_loop == NULL || !g_main_loop_is_running(webview_test_loop)) {
		g_usleep(1000);
	}
}

static void webview_test_quit_loop(void) {
	webview_test_wait_running();
	g_main_loop_quit(webview_test_loop);
}

// webview_test_on_loop_thread reports whether the caller is the thread that runs
// the main loop.
static int webview_test_on_loop_thread(void) {
	return g_thread_self() == (GThread *)webview_test_loop_thread ? 1 : 0;
}
*/
import "C"

func testRunMainLoop()     { C.webview_test_run_loop() }
func testWaitLoopRunning() { C.webview_test_wait_running() }
func testQuitMainLoop()    { C.webview_test_quit_loop() }
func testOnLoopThread() bool {
	return C.webview_test_on_loop_thread() != 0
}
