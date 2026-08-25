//go:build darwin && !ios && !server

package application

// Support for the dispatch regression test in mainthread_darwin_test.go. It
// lives outside the test file because the go tool does not allow cgo there, and
// the test needs to drive the main run loop and nest a loop inside a dispatched
// callback the way an AppKit modal loop does.

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

static volatile int wailsTestLoopRunning = 0;

// Brings AppKit up and reports whether it started. This is what makes the test
// meaningful: until AppKit is loaded the main dispatch queue is drained by
// CFRunLoop, which re-enters and delivers queued blocks inside a nested loop.
// Once AppKit owns the draining it does not re-enter, which is the behaviour a
// real application has.
static int wailsTestInitAppKit(void) {
	return (NSApplicationLoad() && NSApp != nil) ? 1 : 0;
}

// Turns the main run loop until wailsTestStopMainLoop is called. Run in slices
// rather than with CFRunLoopRun, which returns immediately when the loop has no
// sources of its own.
static void wailsTestRunMainLoop(void) {
	wailsTestLoopRunning = 1;
	while (wailsTestLoopRunning) {
		CFRunLoopRunInMode(kCFRunLoopDefaultMode, 0.05, false);
	}
}

static void wailsTestStopMainLoop(void) { wailsTestLoopRunning = 0; }

// Runs a nested run loop for the given number of seconds without returning,
// which is what [NSAlert runModal] and menu tracking do while they are on
// screen.
static void wailsTestRunNestedLoop(double seconds) {
	CFRunLoopRunInMode(kCFRunLoopDefaultMode, seconds, false);
}
*/
import "C"

import "time"

// initAppKitForTest reports whether AppKit came up.
func initAppKitForTest() bool { return C.wailsTestInitAppKit() != 0 }
func runMainRunLoopForTest()  { C.wailsTestRunMainLoop() }
func stopMainRunLoopForTest() { C.wailsTestStopMainLoop() }

func runNestedRunLoopForTest(d time.Duration) {
	C.wailsTestRunNestedLoop(C.double(d.Seconds()))
}
