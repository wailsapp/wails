//go:build openharmony && cgo && !server

package application

// Main-thread dispatch is routed through the WailsBridge: runOnMainThread posts
// to the OHOS JS thread, which calls back into Go via
// WailsOhNativeMainThreadCallback (mirrors mainthread_android.go).

func (a *openharmonyApp) isOnMainThread() bool {
	return ohBridgeBool("isMainThread", "{}")
}

func (a *openharmonyApp) dispatchOnMainThread(id uint) {
	ohBridgeVoidInt("runOnMainThread", int(id))
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
