//go:build darwin && purego && !ios && !server

package application

import (
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

// CoreFoundation run loop bindings. As in the cgo backend, callbacks are not
// sent with dispatch_async to the main queue: AppKit drains that queue and does
// not re-enter, so once a callback shows a modal dialog or a menu - calls which
// do not return until they are dismissed - everything queued behind it waits
// that long too. The run loop keeps turning throughout, so schedule the
// callback on it in the common modes instead and it is delivered either way.
var (
	cfRunLoopGetMain      func() uintptr
	cfRunLoopPerformBlock func(loop, mode uintptr, block objc.Block)
	cfRunLoopWakeUp       func(loop uintptr)
	cfRunLoopCommonModes  uintptr

	// dispatch_async to the main queue, for the places where the cgo backend
	// uses it (the tray click handler). dispatch_get_main_queue() is the
	// address of the global _dispatch_main_q symbol.
	dispatchMainQueue uintptr
	dispatchAsync     func(queue uintptr, block objc.Block)
)

func init() {
	sys, err := purego.Dlopen("/usr/lib/libSystem.B.dylib", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic("wails/purego: failed to load libSystem: " + err.Error())
	}
	if dispatchMainQueue, err = purego.Dlsym(sys, "_dispatch_main_q"); err != nil {
		panic("wails/purego: failed to resolve _dispatch_main_q: " + err.Error())
	}
	purego.RegisterLibFunc(&dispatchAsync, sys, "dispatch_async")

	lib, err := purego.Dlopen("/System/Library/Frameworks/CoreFoundation.framework/CoreFoundation", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic("wails/purego: failed to load CoreFoundation: " + err.Error())
	}
	modes, err := purego.Dlsym(lib, "kCFRunLoopCommonModes")
	if err != nil {
		panic("wails/purego: failed to resolve kCFRunLoopCommonModes: " + err.Error())
	}
	cfRunLoopCommonModes = **(**uintptr)(unsafe.Pointer(&modes))
	purego.RegisterLibFunc(&cfRunLoopGetMain, lib, "CFRunLoopGetMain")
	purego.RegisterLibFunc(&cfRunLoopPerformBlock, lib, "CFRunLoopPerformBlock")
	purego.RegisterLibFunc(&cfRunLoopWakeUp, lib, "CFRunLoopWakeUp")
}

// performOnMainRunLoop schedules fn on the main run loop in the common modes
// and wakes it, so fn also runs inside nested (modal, menu tracking) loops.
func performOnMainRunLoop(fn func()) {
	block := objc.NewBlock(func(objc.Block) { fn() })
	mainLoop := cfRunLoopGetMain()
	cfRunLoopPerformBlock(mainLoop, cfRunLoopCommonModes, block)
	cfRunLoopWakeUp(mainLoop)
	// CFRunLoopPerformBlock copies the block; drop our +1 so the block (and its
	// pinned Go closure) is freed after execution.
	block.Release()
}

func (m *macosApp) isOnMainThread() bool {
	return get[bool](class("NSThread"), "isMainThread")
}

func (m *macosApp) dispatchOnMainThread(id uint) {
	performOnMainRunLoop(func() { dispatchOnMainThreadCallback(id) })
}

// runOnMain runs fn synchronously on the main thread if we are not already
// there. Used internally by helpers that must touch AppKit off the main thread.
// It goes through the run loop too, so it does not wait for a modal to close.
func runOnMain(fn func()) {
	if get[bool](class("NSThread"), "isMainThread") {
		fn()
		return
	}
	done := make(chan struct{})
	performOnMainRunLoop(func() {
		fn()
		close(done)
	})
	<-done
}

// dispatchOnMainThreadCallback is the pure-Go analog of the cgo //export: it
// runs the function registered under id.
func dispatchOnMainThreadCallback(id uint) {
	runMainThreadFunction(id)
}
