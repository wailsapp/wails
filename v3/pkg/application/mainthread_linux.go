//go:build linux

package application

func (m *linuxApp) dispatchOnMainThread(id uint) {
	dispatchOnMainThread(id)
}

func executeOnMainThread(callbackID uint) {
	mainThreadFunctionStoreLock.RLock()
	fn := mainThreadFunctionStore[callbackID]
	if fn == nil {
		Fatal("dispatchCallback called with invalid id: %v", callbackID)
	}
	delete(mainThreadFunctionStore, callbackID)
	mainThreadFunctionStoreLock.RUnlock()
	fn()
}
