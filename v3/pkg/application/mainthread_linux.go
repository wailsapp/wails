//go:build linux && !android

package application

func (a *linuxApp) dispatchOnMainThread(id uint) {
	dispatchOnMainThread(id)
}

func executeOnMainThread(callbackID uint) {
	mainThreadFunctionStoreLock.Lock()
	fn := mainThreadFunctionStore[callbackID]
	if fn == nil {
		mainThreadFunctionStoreLock.Unlock()
		Fatal("dispatchCallback called with invalid id: %v", callbackID)
		return
	}
	delete(mainThreadFunctionStore, callbackID)
	mainThreadFunctionStoreLock.Unlock()
	fn()
}
