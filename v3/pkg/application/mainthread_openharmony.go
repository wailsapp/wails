//go:build openharmony && cgo && !server

package application

import "strconv"

// Main-thread dispatch is routed through the WailsBridge: runOnMainThread posts
// to the OHOS JS thread, which calls back into Go via
// WailsOhNativeMainThreadCallback (mirrors mainthread_android.go).

func (a *openharmonyApp) isOnMainThread() bool { return ohOnArkTSThread() }

func (a *openharmonyApp) dispatchOnMainThread(id uint) {
	if !ohBridgePost(0, "runOnMainThread", `{"id":`+strconv.FormatUint(uint64(id), 10)+`}`) {
		mainThreadFunctionStoreLock.Lock()
		delete(mainThreadFunctionStore, id)
		mainThreadFunctionStoreLock.Unlock()
		panic("OpenHarmony main-thread dispatcher unavailable")
	}
}

// openharmonyMainThreadCallback is invoked from NAPI on the OHOS main thread.
func openharmonyMainThreadCallback(callbackID uint) {
	mainThreadFunctionStoreLock.Lock()
	fn := mainThreadFunctionStore[callbackID]
	if fn == nil {
		mainThreadFunctionStoreLock.Unlock()
		Fatal("dispatchOnMainThread called with invalid id: %v", callbackID)
		return
	}
	delete(mainThreadFunctionStore, callbackID)
	mainThreadFunctionStoreLock.Unlock()
	fn()
}
