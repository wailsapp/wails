//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa

#include "Cocoa/Cocoa.h"

extern void dispatchOnMainThreadCallback(unsigned int);

static void dispatchOnMainThread(unsigned int id) {
	// Not dispatch_async to the main queue: AppKit drains that queue and does
	// not re-enter, so once a callback shows a modal dialog or a menu - calls
	// which do not return until they are dismissed - everything queued behind it
	// waits that long too. The run loop keeps turning throughout, so schedule
	// the callback on it instead and it is delivered either way.
	CFRunLoopRef mainLoop = CFRunLoopGetMain();
	CFRunLoopPerformBlock(mainLoop, kCFRunLoopCommonModes, ^{
		dispatchOnMainThreadCallback(id);
	});
	CFRunLoopWakeUp(mainLoop);
}

static bool onMainThread() {
	return [NSThread isMainThread];
}

*/
import "C"

func (m *macosApp) isOnMainThread() bool {
	return bool(C.onMainThread())
}

func (m *macosApp) dispatchOnMainThread(id uint) {
	C.dispatchOnMainThread(C.uint(id))
}

//export dispatchOnMainThreadCallback
func dispatchOnMainThreadCallback(callbackID C.uint) {
	mainThreadFunctionStoreLock.Lock()
	id := uint(callbackID)
	fn := mainThreadFunctionStore[id]
	if fn == nil {
		mainThreadFunctionStoreLock.Unlock()
		Fatal("dispatchCallback called with invalid id: %v", id)
		return
	}
	delete(mainThreadFunctionStore, id)
	mainThreadFunctionStoreLock.Unlock()
	fn()
}
