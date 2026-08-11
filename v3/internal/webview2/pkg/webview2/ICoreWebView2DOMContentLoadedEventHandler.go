//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2DOMContentLoadedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2DOMContentLoadedEventHandler struct {
	Vtbl *ICoreWebView2DOMContentLoadedEventHandlerVtbl
	impl ICoreWebView2DOMContentLoadedEventHandlerImpl
}

func (i *ICoreWebView2DOMContentLoadedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2DOMContentLoadedEventHandlerIUnknownQueryInterface(this *ICoreWebView2DOMContentLoadedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2DOMContentLoadedEventHandlerIUnknownAddRef(this *ICoreWebView2DOMContentLoadedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2DOMContentLoadedEventHandlerIUnknownRelease(this *ICoreWebView2DOMContentLoadedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2DOMContentLoadedEventHandlerInvoke(this *ICoreWebView2DOMContentLoadedEventHandler, sender *ICoreWebView2, args *ICoreWebView2DOMContentLoadedEventArgs) uintptr {
	return this.impl.DOMContentLoaded(sender, args)
}

type ICoreWebView2DOMContentLoadedEventHandlerImpl interface {
	IUnknownImpl
	DOMContentLoaded(sender *ICoreWebView2, args *ICoreWebView2DOMContentLoadedEventArgs) uintptr
}

var ICoreWebView2DOMContentLoadedEventHandlerFn = ICoreWebView2DOMContentLoadedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2DOMContentLoadedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2DOMContentLoadedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2DOMContentLoadedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2DOMContentLoadedEventHandlerInvoke),
}

func NewICoreWebView2DOMContentLoadedEventHandler(impl ICoreWebView2DOMContentLoadedEventHandlerImpl) *ICoreWebView2DOMContentLoadedEventHandler {
	return &ICoreWebView2DOMContentLoadedEventHandler{
		Vtbl: &ICoreWebView2DOMContentLoadedEventHandlerFn,
		impl: impl,
	}
}
