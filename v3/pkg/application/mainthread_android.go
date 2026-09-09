//go:build android

package application

// Main-thread dispatch is routed through the WailsBridge: runOnMainThread
// posts to the Android main looper, which calls back into Go via
// nativeMainThreadCallback (see application_android.go).

func (a *androidApp) isOnMainThread() bool {
	return androidBridgeBool("isMainThread")
}

func (a *androidApp) dispatchOnMainThread(id uint) {
	androidBridgeVoidInt("runOnMainThread", int(id))
}

// androidMainThreadCallback is invoked from JNI on the Android main thread.
func androidMainThreadCallback(callbackID uint) {
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
