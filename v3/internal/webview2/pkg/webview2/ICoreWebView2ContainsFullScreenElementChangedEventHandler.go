//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2ContainsFullScreenElementChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ContainsFullScreenElementChangedEventHandler struct {
	Vtbl *ICoreWebView2ContainsFullScreenElementChangedEventHandlerVtbl
	impl ICoreWebView2ContainsFullScreenElementChangedEventHandlerImpl
}

func (i *ICoreWebView2ContainsFullScreenElementChangedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2ContainsFullScreenElementChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2ContainsFullScreenElementChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ContainsFullScreenElementChangedEventHandlerIUnknownAddRef(this *ICoreWebView2ContainsFullScreenElementChangedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2ContainsFullScreenElementChangedEventHandlerIUnknownRelease(this *ICoreWebView2ContainsFullScreenElementChangedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2ContainsFullScreenElementChangedEventHandlerInvoke(this *ICoreWebView2ContainsFullScreenElementChangedEventHandler, sender *ICoreWebView2, args *IUnknown) uintptr {
	return this.impl.ContainsFullScreenElementChanged(sender, args)
}

type ICoreWebView2ContainsFullScreenElementChangedEventHandlerImpl interface {
	IUnknownImpl
	ContainsFullScreenElementChanged(sender *ICoreWebView2, args *IUnknown) uintptr
}

var ICoreWebView2ContainsFullScreenElementChangedEventHandlerFn = ICoreWebView2ContainsFullScreenElementChangedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2ContainsFullScreenElementChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ContainsFullScreenElementChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ContainsFullScreenElementChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ContainsFullScreenElementChangedEventHandlerInvoke),
}

func NewICoreWebView2ContainsFullScreenElementChangedEventHandler(impl ICoreWebView2ContainsFullScreenElementChangedEventHandlerImpl) *ICoreWebView2ContainsFullScreenElementChangedEventHandler {
	return &ICoreWebView2ContainsFullScreenElementChangedEventHandler{
		Vtbl: &ICoreWebView2ContainsFullScreenElementChangedEventHandlerFn,
		impl: impl,
	}
}
