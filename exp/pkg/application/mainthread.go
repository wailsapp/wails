//go:build darwin

package application

/*
extern void dispatch(unsigned int id);
*/
import "C"
import (
	"os"
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
			panic("Too many functions stored")
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
		println("***** dispatchCallback called with invalid id: ", id)
		os.Exit(1)
	}
	delete(mainThreadFuntionStore, id)
	mainThreadFuntionStoreLock.RUnlock()
	fn()
}
