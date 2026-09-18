//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2MoveFocusRequestedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2MoveFocusRequestedEventHandler struct {
	Vtbl *ICoreWebView2MoveFocusRequestedEventHandlerVtbl
	impl ICoreWebView2MoveFocusRequestedEventHandlerImpl
}

func (i *ICoreWebView2MoveFocusRequestedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2MoveFocusRequestedEventHandlerIUnknownQueryInterface(this *ICoreWebView2MoveFocusRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2MoveFocusRequestedEventHandlerIUnknownAddRef(this *ICoreWebView2MoveFocusRequestedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2MoveFocusRequestedEventHandlerIUnknownRelease(this *ICoreWebView2MoveFocusRequestedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2MoveFocusRequestedEventHandlerInvoke(this *ICoreWebView2MoveFocusRequestedEventHandler, sender *ICoreWebView2Controller, args *ICoreWebView2MoveFocusRequestedEventArgs) uintptr {
	return this.impl.MoveFocusRequested(sender, args)
}

type ICoreWebView2MoveFocusRequestedEventHandlerImpl interface {
	IUnknownImpl
	MoveFocusRequested(sender *ICoreWebView2Controller, args *ICoreWebView2MoveFocusRequestedEventArgs) uintptr
}

var ICoreWebView2MoveFocusRequestedEventHandlerFn = ICoreWebView2MoveFocusRequestedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2MoveFocusRequestedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2MoveFocusRequestedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2MoveFocusRequestedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2MoveFocusRequestedEventHandlerInvoke),
}

func NewICoreWebView2MoveFocusRequestedEventHandler(impl ICoreWebView2MoveFocusRequestedEventHandlerImpl) *ICoreWebView2MoveFocusRequestedEventHandler {
	return &ICoreWebView2MoveFocusRequestedEventHandler{
		Vtbl: &ICoreWebView2MoveFocusRequestedEventHandlerFn,
		impl: impl,
	}
}
