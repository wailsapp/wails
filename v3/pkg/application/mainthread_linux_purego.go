//go:build linux && purego

package application

import "github.com/ebitengine/purego"

const (
	G_SOURCE_REMOVE = 0
)

func (m *linuxApp) dispatchOnMainThread(id uint) {
	var dispatch func(uintptr)
	purego.RegisterLibFunc(&dispatch, gtk, "g_idle_add")
	dispatch(purego.NewCallback(func(uintptr) int {
		dispatchOnMainThreadCallback(id)
		return G_SOURCE_REMOVE
	}))
}

func dispatchOnMainThreadCallback(callbackID uint) {
	mainThreadFunctionStoreLock.RLock()
	id := uint(callbackID)
	fn := mainThreadFunctionStore[id]
	if fn == nil {
		Fatal("dispatchCallback called with invalid id: %v", id)
	}
	delete(mainThreadFunctionStore, id)
	mainThreadFunctionStoreLock.RUnlock()
	fn()
}
