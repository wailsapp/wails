//go:build darwin

package application

/*
extern void dispatch(unsigned int id);
*/
import "C"
import (
	"sync"
)

var mainThreadFuntionStore = make(map[uint]func())
var mainThreadFuntionStoreLock sync.RWMutex

func generateFunctionStoreID() uint {
	startID := 0
	for {
		if _, ok := mainThreadFuntionStore[uint(startID)]; !ok {
			return uint(startID)
		}
		startID++
		if startID == 0 {
			Fatal("Too many functions have been dispatched to the main thread")
		}
	}
}

func Dispatch(fn func()) {
	mainThreadFuntionStoreLock.Lock()
	id := generateFunctionStoreID()
	mainThreadFuntionStore[id] = fn
	mainThreadFuntionStoreLock.Unlock()
	C.dispatch(C.uint(id))
}

//export dispatchCallback
func dispatchCallback(callbackID C.uint) {
	mainThreadFuntionStoreLock.RLock()
	id := uint(callbackID)
	fn := mainThreadFuntionStore[id]
	if fn == nil {
		Fatal("dispatchCallback called with invalid id: ", id)
	}
	delete(mainThreadFuntionStore, id)
	mainThreadFuntionStoreLock.RUnlock()
	fn()
}
