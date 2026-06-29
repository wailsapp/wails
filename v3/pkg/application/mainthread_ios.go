//go:build ios

package application

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UIKit
#import <Foundation/Foundation.h>
#import <dispatch/dispatch.h>

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

func (a *iosApp) isOnMainThread() bool {
	return bool(C.onMainThread())
}

func (a *iosApp) dispatchOnMainThread(id uint) {
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
