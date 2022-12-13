package application

import (
	"sync"
)

var mainThreadFunctionStore = make(map[uint]func())
var mainThreadFunctionStoreLock sync.RWMutex

func generateFunctionStoreID() uint {
	startID := 0
	for {
		if _, ok := mainThreadFunctionStore[uint(startID)]; !ok {
			return uint(startID)
		}
		startID++
		if startID == 0 {
			Fatal("Too many functions have been dispatched to the main thread")
		}
	}
}

func DispatchOnMainThread(fn func()) {
	mainThreadFunctionStoreLock.Lock()
	id := generateFunctionStoreID()
	mainThreadFunctionStore[id] = fn
	mainThreadFunctionStoreLock.Unlock()
	// Call platform specific dispatch function
	platformDispatch(id)
}
