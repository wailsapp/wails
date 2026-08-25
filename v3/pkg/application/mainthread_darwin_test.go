//go:build darwin && !ios && !server

package application

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application/internal/mainthreadharness"
)

const (
	// How long the stand-in for a modal dialog holds the main thread.
	nestedLoopDuration = 2 * time.Second
	// What a main thread callback is allowed to take. In practice it arrives in
	// well under a millisecond; starved, it takes the whole nested loop. The
	// budget sits far from both, and is used as the timeout and as the
	// assertion, so the two can never disagree.
	deliveryBudget = 500 * time.Millisecond
)

var (
	appKitReady bool
	mainLoop    *mainthreadharness.Loop
)

// TestMain runs the package's tests on their own goroutine while the main
// thread turns the run loop, as an application's message loop does. Callbacks
// scheduled with dispatchOnMainThread are delivered to the main thread, so
// without a loop running there nothing would ever execute them.
func TestMain(m *testing.M) {
	appKitReady = mainthreadharness.InitAppKit()

	mainLoop = mainthreadharness.NewLoop()
	defer mainLoop.Free()

	var code int
	done := make(chan struct{})
	go func() {
		defer close(done)
		code = m.Run()
		mainLoop.Stop()
	}()
	mainLoop.Run()
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

	entered := make(chan struct{})
	delivered := make(chan struct{})

	// A callback which does not return for the duration, standing in for a
	// dialog or an open menu.
	scheduleOnMainForTest(func() {
		close(entered)
		mainthreadharness.RunNested(nestedLoopDuration)
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
		if elapsed := time.Since(start); elapsed > deliveryBudget {
			t.Fatalf("callback took %s to arrive, budget is %s", elapsed, deliveryBudget)
		}
	case <-time.After(deliveryBudget):
		t.Fatalf("main thread dispatch starved by a nested run loop: nothing delivered "+
			"within %s while a %s modal loop was up. An open dialog or menu would freeze "+
			"every window, menu and runtime call until it was dismissed",
			deliveryBudget, nestedLoopDuration)
	}
}

// TestLoopHonoursStopBeforeRun guards the harness itself. TestMain stops the
// loop from the goroutine running the suite, which can finish before the main
// thread reaches Run; if Run were to put the loop into the running state, that
// stop would be overwritten and the test binary would never exit.
func TestLoopHonoursStopBeforeRun(t *testing.T) {
	loop := mainthreadharness.NewLoop()
	defer loop.Free()

	loop.Stop()

	returned := make(chan struct{})
	go func() {
		loop.Run()
		close(returned)
	}()

	select {
	case <-returned:
	case <-time.After(3 * time.Second):
		t.Fatal("Run kept turning after Stop: if the suite finished first, the test binary would hang")
	}
}

// TestNoTestOnlyHelpersInShippedFiles keeps the run loop harness out of the
// package that ships. It is only needed by the tests, and it carries cgo.
func TestNoTestOnlyHelpersInShippedFiles(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		body, err := os.ReadFile(filepath.Join(".", name))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(body), "mainthreadharness") {
			t.Errorf("%s pulls the test harness into the shipped package", name)
		}
	}
}
