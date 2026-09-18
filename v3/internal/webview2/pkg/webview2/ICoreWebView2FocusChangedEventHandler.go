//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2FocusChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FocusChangedEventHandler struct {
	Vtbl *ICoreWebView2FocusChangedEventHandlerVtbl
	impl ICoreWebView2FocusChangedEventHandlerImpl
}

func (i *ICoreWebView2FocusChangedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2FocusChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2FocusChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FocusChangedEventHandlerIUnknownAddRef(this *ICoreWebView2FocusChangedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2FocusChangedEventHandlerIUnknownRelease(this *ICoreWebView2FocusChangedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2FocusChangedEventHandlerInvoke(this *ICoreWebView2FocusChangedEventHandler, sender *ICoreWebView2Controller, args *IUnknown) uintptr {
	return this.impl.FocusChanged(sender, args)
}

type ICoreWebView2FocusChangedEventHandlerImpl interface {
	IUnknownImpl
	FocusChanged(sender *ICoreWebView2Controller, args *IUnknown) uintptr
}

var ICoreWebView2FocusChangedEventHandlerFn = ICoreWebView2FocusChangedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2FocusChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FocusChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FocusChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FocusChangedEventHandlerInvoke),
}

func NewICoreWebView2FocusChangedEventHandler(impl ICoreWebView2FocusChangedEventHandlerImpl) *ICoreWebView2FocusChangedEventHandler {
	return &ICoreWebView2FocusChangedEventHandler{
		Vtbl: &ICoreWebView2FocusChangedEventHandlerFn,
		impl: impl,
	}
}
