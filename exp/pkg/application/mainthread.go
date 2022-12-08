//go:build darwin

package application

/*
extern void dispatch(unsigned int id);
*/
import "C"
import "strconv"

var mainThreadFuntionStore = make(map[uint]func())

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
	id := generateFunctionStoreID()
	mainThreadFuntionStore[id] = fn
	C.dispatch(C.uint(id))
}

//export dispatchCallback
func dispatchCallback(id C.uint) {

	fn := mainThreadFuntionStore[uint(id)]
	if fn == nil {
		panic("dispatchCallback called with invalid id " + strconv.Itoa(int(id)))
	}
	fn()
	delete(mainThreadFuntionStore, uint(id))
}
