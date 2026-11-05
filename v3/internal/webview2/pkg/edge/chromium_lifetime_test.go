//go:build windows

package edge

import (
	"runtime"
	"testing"
	"unsafe"
	"weak"
)

// comRef is a handler reference as native code holds it: the address as an
// integer, which does not keep the Go object alive, and the vtable read from
// the object while it was reachable. The vtable is a package-level value shared
// by every handler of its type, so holding it keeps no handler alive. Nothing
// here turns the integer back into a Go pointer; checkptr, enabled by -race,
// rejects that conversion.
type comRef struct {
	addr uintptr
	vtbl *IUnknownVtbl
}

func newComRef(handler unsafe.Pointer) comRef {
	return comRef{addr: uintptr(handler), vtbl: (*IUnknown)(handler).Vtbl}
}

func (r comRef) addRef() uintptr {
	n, _, _ := r.vtbl.AddRef.Call(r.addr)
	return n
}

func (r comRef) release() uintptr {
	n, _, _ := r.vtbl.Release.Call(r.addr)
	return n
}

// refs reports the native references counted for the handler, matching it by
// address. It cannot call handlerRefCount, which takes the handler as a
// pointer: rebuilding that pointer from addr is the conversion checkptr rejects.
func (r comRef) refs() uint32 {
	handlerRefs.Lock()
	defer handlerRefs.Unlock()
	for handler, n := range handlerRefs.counts {
		if uintptr(handler) == r.addr {
			return n
		}
	}
	return 0
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
// Chromium. Only the native reference, the handler's Invoke slot and a weak
// pointer escape.
//
//go:noinline
func abandonWithNativeReference(t *testing.T) (comRef, ComProc, weak.Pointer[Chromium]) {
	e := NewChromium()
	handler := newComRef(unsafe.Pointer(e.controllerCompleted))
	if refs := handler.addRef(); refs == 0 {
		t.Fatalf("AddRef returned %d", refs)
	}
	invoke := e.controllerCompleted.vtbl.Invoke
	e.ShuttingDown()
	return handler, invoke, weak.Make(e)
}

func TestNativeReferenceKeepsAbandonedHandlerAlive(t *testing.T) {
	handler, invoke, instance := abandonWithNativeReference(t)
	collectGarbage()
	if instance.Value() == nil {
		// Invoking the handler now would read freed memory.
		t.Fatal("abandoned Chromium was collected while native code held its callback")
	}

	closed := false
	controller := &ICoreWebView2Controller{vtbl: &_ICoreWebView2ControllerVtbl{
		Close: NewComProc(func(uintptr) uintptr { closed = true; return 0 }),
	}}
	invoke.Call(handler.addr, 0, uintptr(unsafe.Pointer(controller)))
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
