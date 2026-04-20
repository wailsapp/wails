//go:build darwin && !ios

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa

#include "Cocoa/Cocoa.h"

extern void dispatchOnMainThreadCallback(unsigned int);

static void dispatchOnMainThread(unsigned int id) {
	dispatch_async(dispatch_get_main_queue(), ^{
		dispatchOnMainThreadCallback(id);
	});
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
