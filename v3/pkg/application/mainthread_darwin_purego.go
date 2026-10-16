//go:build darwin && purego && !ios && !server

package application

import (
	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

// libdispatch bindings. dispatch_get_main_queue() is defined as the address of
// the global _dispatch_main_q symbol, so we resolve that symbol once and treat
// its address as the main queue handle. Blocks submitted here run on the
// process main thread while [NSApp run] pumps the main run loop.
var (
	dispatchMainQueue uintptr
	dispatchAsync     func(queue uintptr, block objc.Block)
)

func init() {
	lib, err := purego.Dlopen("/usr/lib/libSystem.B.dylib", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic("wails/purego: failed to load libSystem: " + err.Error())
	}
	mainQ, err := purego.Dlsym(lib, "_dispatch_main_q")
	if err != nil {
		panic("wails/purego: failed to resolve _dispatch_main_q: " + err.Error())
	}
	dispatchMainQueue = mainQ
	purego.RegisterLibFunc(&dispatchAsync, lib, "dispatch_async")
}

func (m *macosApp) isOnMainThread() bool {
	return get[bool](class("NSThread"), "isMainThread")
}

func (m *macosApp) dispatchOnMainThread(id uint) {
	// Mirror the cgo path: hop to the main queue, then run the stored callback.
	block := objc.NewBlock(func(objc.Block) {
		dispatchOnMainThreadCallback(id)
	})
	dispatchAsync(dispatchMainQueue, block)
	// dispatch_async takes its own copy of the block synchronously; drop our
	// +1 so the block (and its pinned Go closure) is freed after execution.
	block.Release()
}

// dispatchOnMainThreadCallback runs the Go function previously registered under
// id in the shared mainThreadFunctionStore. This is a straight port of the cgo
// //export callback (same store, same semantics) minus the C signature.
func dispatchOnMainThreadCallback(id uint) {
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

// runOnMain runs fn synchronously on the main thread if we are not already
// there. Used internally by helpers that must touch AppKit off the main thread.
func runOnMain(fn func()) {
	if get[bool](class("NSThread"), "isMainThread") {
		fn()
		return
	}
	done := make(chan struct{})
	block := objc.NewBlock(func(objc.Block) {
		fn()
		close(done)
	})
	dispatchAsync(dispatchMainQueue, block)
	// dispatch_async copied the block; release our reference (see above).
	block.Release()
	<-done
}
