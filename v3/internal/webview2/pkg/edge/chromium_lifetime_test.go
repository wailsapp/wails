//go:build windows

package edge

import (
	"runtime"
	"testing"
	"unsafe"
	"weak"
)

// comAddr is a handler address as native code holds it: an integer that does
// not keep the Go object alive.
type comAddr uintptr

func (a comAddr) unknown() *IUnknown { return (*IUnknown)(unsafe.Pointer(a)) }

func (a comAddr) addRef() uintptr {
	r, _, _ := a.unknown().Vtbl.AddRef.Call(uintptr(a))
	return r
}

func (a comAddr) release() uintptr {
	r, _, _ := a.unknown().Vtbl.Release.Call(uintptr(a))
	return r
}

// collectGarbage runs several cycles with allocation churn so that freed
// handler memory is likely to be reused.
func collectGarbage() {
	for i := 0; i < 8; i++ {
		runtime.GC()
		churn := make([][]byte, 0, 256)
		for j := 0; j < cap(churn); j++ {
			churn = append(churn, make([]byte, 128))
		}
		runtime.KeepAlive(churn)
	}
}

func waitCollected[T any](p weak.Pointer[T]) bool {
	for i := 0; i < 10; i++ {
		collectGarbage()
		if p.Value() == nil {
			return true
		}
	}
	return false
}

// abandonWithNativeReference mirrors a timed-out initialization: native code
// takes a reference to the controller callback, then the window drops the
// Chromium. Only the handler address and a weak pointer escape.
//
//go:noinline
func abandonWithNativeReference(t *testing.T) (comAddr, weak.Pointer[Chromium]) {
	e := NewChromium()
	handler := comAddr(uintptr(unsafe.Pointer(e.controllerCompleted)))
	if refs := handler.addRef(); refs == 0 {
		t.Fatalf("AddRef returned %d", refs)
	}
	e.ShuttingDown()
	return handler, weak.Make(e)
}

func TestNativeReferenceKeepsAbandonedHandlerAlive(t *testing.T) {
	handler, instance := abandonWithNativeReference(t)
	collectGarbage()
	if instance.Value() == nil {
		// Invoking the handler now would read freed memory.
		t.Fatal("abandoned Chromium was collected while native code held its callback")
	}

	closed := false
	controller := &ICoreWebView2Controller{vtbl: &_ICoreWebView2ControllerVtbl{
		Close: NewComProc(func(uintptr) uintptr { closed = true; return 0 }),
	}}
	invoke := (*iCoreWebView2CreateCoreWebView2ControllerCompletedHandler)(unsafe.Pointer(handler)).vtbl.Invoke
	invoke.Call(uintptr(handler), 0, uintptr(unsafe.Pointer(controller)))
	if !closed {
		t.Fatal("late controller callback did not close the controller")
	}

	if refs := handler.release(); refs != 0 {
		t.Fatalf("final Release returned %d", refs)
	}
	if !waitCollected(instance) {
		t.Fatal("abandoned Chromium was retained after the final native Release")
	}
}
