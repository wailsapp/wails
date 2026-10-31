//go:build darwin && !ios && !server

package application

// runMainThreadFunction runs the Go function registered under id in
// mainThreadFunctionStore. Both backends call it from their dispatch callback.
func runMainThreadFunction(id uint) {
	mainThreadFunctionStoreLock.Lock()
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
