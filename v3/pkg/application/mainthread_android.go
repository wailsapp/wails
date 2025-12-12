//go:build android

package application

// isOnMainThread returns whether the current goroutine is on the main thread
func (a *androidApp) isOnMainThread() bool {
	// On Android, Go runs in its own thread separate from the UI thread
	// UI operations need to be dispatched via JNI to the main thread
	return false
}

// dispatchOnMainThread executes a function on the Android main/UI thread
func (a *androidApp) dispatchOnMainThread(id uint) {
	// TODO: Implement via JNI callback to Activity.runOnUiThread()
	// For now, execute the callback directly
	mainThreadFunctionStoreLock.RLock()
	fn := mainThreadFunctionStore[id]
	if fn == nil {
		mainThreadFunctionStoreLock.RUnlock()
		return
	}
	delete(mainThreadFunctionStore, id)
	mainThreadFunctionStoreLock.RUnlock()
	fn()
}
