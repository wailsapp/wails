//go:build darwin && !ios && !server

package application

import (
	"os"
	"testing"
	"time"
)

var appKitReady bool

// TestMain runs the package's tests on their own goroutine while the main
// thread turns the run loop, as an application's message loop does. Callbacks
// scheduled with dispatchOnMainThread are delivered to the main thread, so
// without a loop running there nothing would ever execute them.
func TestMain(m *testing.M) {
	appKitReady = initAppKitForTest()

	var code int
	done := make(chan struct{})
	go func() {
		defer close(done)
		code = m.Run()
		stopMainRunLoopForTest()
	}()
	runMainRunLoopForTest()
	<-done
	os.Exit(code)
}

// scheduleOnMainForTest queues fn for the main thread exactly as
// App.dispatchOnMainThread does, without needing a running application.
func scheduleOnMainForTest(fn func()) {
	mainThreadFunctionStoreLock.Lock()
	id := generateFunctionStoreID()
	mainThreadFunctionStore[id] = fn
	mainThreadFunctionStoreLock.Unlock()

	var app macosApp
	app.dispatchOnMainThread(id)
}

// TestDispatchIsDeliveredDuringNestedRunLoop covers the main thread dispatch
// used by every InvokeSync and InvokeAsync call, and so by every window, menu
// and runtime operation.
//
// Dialogs and menus reach AppKit through that dispatch, and the AppKit call
// which shows them - [NSAlert runModal], or menu tracking - does not return
// until the user dismisses them; it turns a nested run loop instead. AppKit
// drains the main dispatch queue and does not re-enter, so while a callback
// posted to that queue sits in such a call, everything queued behind it waits
// for the dialog or menu to close. Scheduling on the run loop keeps the queue
// free and the callbacks flowing.
func TestDispatchIsDeliveredDuringNestedRunLoop(t *testing.T) {
	if !appKitReady {
		// Without AppKit the main queue is drained by CFRunLoop, which re-enters
		// and delivers the callback regardless, so the test would pass whether or
		// not the bug is present. Skip rather than report a false green.
		t.Skip("AppKit is not available, so the behaviour under test does not arise")
	}

	const nested = 2 * time.Second

	entered := make(chan struct{})
	delivered := make(chan struct{})

	// A callback that does not return for `nested`, standing in for a dialog.
	scheduleOnMainForTest(func() {
		close(entered)
		runNestedRunLoopForTest(nested)
	})

	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		t.Fatal("the first callback never ran: the main run loop is not turning")
	}

	start := time.Now()
	scheduleOnMainForTest(func() { close(delivered) })

	select {
	case <-delivered:
		if elapsed := time.Since(start); elapsed >= nested {
			t.Fatalf("callback ran only once the nested loop had finished, after %s", elapsed)
		}
	case <-time.After(nested - 500*time.Millisecond):
		t.Fatal("main thread dispatch starved by a nested run loop: an open dialog " +
			"or menu would freeze every window, menu and runtime call until it was dismissed")
	}
}
