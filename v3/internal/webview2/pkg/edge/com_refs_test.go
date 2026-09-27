//go:build windows

package edge

import (
	"testing"
	"unsafe"
	"weak"
)

func TestHandlerReferenceCounting(t *testing.T) {
	e := NewChromium()
	handler := comAddr(uintptr(unsafe.Pointer(e.processFailed)))
	unbalanced := unbalancedHandlerReleases.Load()
	if handler.addRef() != 1 || handler.addRef() != 2 {
		t.Fatal("AddRef did not count native references")
	}
	if handler.release() != 1 || handlerRefCount(unsafe.Pointer(e.processFailed)) != 1 {
		t.Fatal("Release did not drop one reference")
	}
	if handler.release() != 0 || handlerRefCount(unsafe.Pointer(e.processFailed)) != 0 {
		t.Fatal("final Release left the handler registered")
	}
	if handler.release() != 0 || unbalancedHandlerReleases.Load() != unbalanced+1 {
		t.Fatal("unbalanced Release was not reported")
	}
}

type hostWindow struct{ events int }

// abandonWithHostCallbacks returns a handler that native code still references
// and the window whose callbacks its abandoned Chromium used to hold.
//
//go:noinline
func abandonWithHostCallbacks() (comAddr, weak.Pointer[hostWindow]) {
	w := &hostWindow{}
	e := NewChromium()
	e.MessageCallback = func(string, *ICoreWebView2, *ICoreWebView2WebMessageReceivedEventArgs) { w.events++ }
	e.MessageWithAdditionalObjectsCallback = e.MessageCallback
	e.WebResourceRequestedCallback = func(*ICoreWebView2WebResourceRequest, *ICoreWebView2WebResourceRequestedEventArgs) { w.events++ }
	e.NavigationStartingCallback = func(*ICoreWebView2) { w.events++ }
	e.NavigationCompletedCallback = func(*ICoreWebView2, *ICoreWebView2NavigationCompletedEventArgs) { w.events++ }
	e.ProcessFailedCallback = func(*ICoreWebView2, *ICoreWebView2ProcessFailedEventArgs) { w.events++ }
	e.ContainsFullScreenElementChangedCallback = func(*ICoreWebView2, *ICoreWebView2ContainsFullScreenElementChangedEventArgs) { w.events++ }
	e.AcceleratorKeyCallback = func(uint) bool { w.events++; return false }
	e.CursorChangedCallback = func(HCURSOR, uint32) { w.events++ }
	e.SetErrorCallback(func(error) { w.events++ })
	handler := comAddr(uintptr(unsafe.Pointer(e.processFailed)))
	handler.addRef()
	e.ShuttingDown()
	return handler, weak.Make(w)
}

func TestAbandonedChromiumDoesNotRetainWindow(t *testing.T) {
	handler, window := abandonWithHostCallbacks()
	defer handler.release()
	if !waitCollected(window) {
		t.Fatal("a native reference to an abandoned Chromium retained its window")
	}
	if handlerRefCount(unsafe.Pointer(uintptr(handler))) != 1 {
		t.Fatal("handler lost its native reference")
	}
}
