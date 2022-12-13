//go:build darwin

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa

#include "Cocoa/Cocoa.h"

extern void dispatchCallback(unsigned int);

static void dispatch(unsigned int id) {
	dispatch_async(dispatch_get_main_queue(), ^{
		dispatchCallback(id);
	});
}

*/
import "C"

func platformDispatch(id uint) {
	C.dispatch(C.uint(id))
}

//export dispatchCallback
func dispatchCallback(callbackID C.uint) {
	mainThreadFunctionStoreLock.RLock()
	id := uint(callbackID)
	fn := mainThreadFunctionStore[id]
	if fn == nil {
		Fatal("dispatchCallback called with invalid id: ", id)
	}
	delete(mainThreadFunctionStore, id)
	mainThreadFunctionStoreLock.RUnlock()
	fn()
}
