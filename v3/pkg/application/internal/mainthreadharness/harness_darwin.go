//go:build darwin && !ios && !server

package mainthreadharness

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>
#include <stdatomic.h>
#include <stdlib.h>

typedef struct {
	atomic_int running;
} wails_loop_t;

// A loop starts out running, so a Stop which lands before Run is honoured
// rather than overwritten by Run and lost.
static void* wailsLoopNew(void) {
	wails_loop_t* loop = malloc(sizeof(wails_loop_t));
	atomic_store(&loop->running, 1);
	return loop;
}

static void wailsLoopFree(void* handle) { free(handle); }

// Turns the run loop until the loop is stopped. Run in slices rather than with
// CFRunLoopRun, which returns immediately when the loop has no sources of its
// own.
static void wailsLoopRun(void* handle) {
	wails_loop_t* loop = (wails_loop_t*)handle;
	while (atomic_load(&loop->running)) {
		CFRunLoopRunInMode(kCFRunLoopDefaultMode, 0.05, false);
	}
}

static void wailsLoopStop(void* handle) {
	atomic_store(&((wails_loop_t*)handle)->running, 0);
}

static int wailsInitAppKit(void) {
	return (NSApplicationLoad() && NSApp != nil) ? 1 : 0;
}

static void wailsRunNested(double seconds) {
	CFRunLoopRunInMode(kCFRunLoopDefaultMode, seconds, false);
}
*/
import "C"

import (
	"time"
	"unsafe"
)

// InitAppKit brings AppKit up and reports whether it started. It matters for
// the dispatch tests: until AppKit is loaded the main dispatch queue is drained
// by CFRunLoop, which re-enters and delivers queued blocks inside a nested
// loop. Once AppKit owns the draining it does not re-enter, which is the
// behaviour a real application has.
func InitAppKit() bool { return C.wailsInitAppKit() != 0 }

// Loop turns a run loop until it is stopped.
type Loop struct {
	handle unsafe.Pointer
}

// NewLoop returns a loop which is already in the running state, so that Stop
// may be called before Run.
func NewLoop() *Loop { return &Loop{handle: C.wailsLoopNew()} }

// Run turns the run loop of the calling thread until Stop is called. It returns
// straight away if Stop was called first.
func (l *Loop) Run() { C.wailsLoopRun(l.handle) }

// Stop ends the loop. It is safe to call from any thread, and before Run.
func (l *Loop) Stop() { C.wailsLoopStop(l.handle) }

// Free releases the loop.
func (l *Loop) Free() { C.wailsLoopFree(l.handle) }

// RunNested runs a nested run loop for d without returning, which is what
// [NSAlert runModal] and menu tracking do while they are on screen.
func RunNested(d time.Duration) { C.wailsRunNested(C.double(d.Seconds())) }
