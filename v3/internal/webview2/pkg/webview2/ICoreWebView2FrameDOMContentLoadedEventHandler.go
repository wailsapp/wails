//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2FrameDOMContentLoadedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FrameDOMContentLoadedEventHandler struct {
	Vtbl *ICoreWebView2FrameDOMContentLoadedEventHandlerVtbl
	impl ICoreWebView2FrameDOMContentLoadedEventHandlerImpl
}

func (i *ICoreWebView2FrameDOMContentLoadedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2FrameDOMContentLoadedEventHandlerIUnknownQueryInterface(this *ICoreWebView2FrameDOMContentLoadedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FrameDOMContentLoadedEventHandlerIUnknownAddRef(this *ICoreWebView2FrameDOMContentLoadedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2FrameDOMContentLoadedEventHandlerIUnknownRelease(this *ICoreWebView2FrameDOMContentLoadedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2FrameDOMContentLoadedEventHandlerInvoke(this *ICoreWebView2FrameDOMContentLoadedEventHandler, sender *ICoreWebView2Frame, args *ICoreWebView2DOMContentLoadedEventArgs) uintptr {
	return this.impl.FrameDOMContentLoaded(sender, args)
}

type ICoreWebView2FrameDOMContentLoadedEventHandlerImpl interface {
	IUnknownImpl
	FrameDOMContentLoaded(sender *ICoreWebView2Frame, args *ICoreWebView2DOMContentLoadedEventArgs) uintptr
}

var ICoreWebView2FrameDOMContentLoadedEventHandlerFn = ICoreWebView2FrameDOMContentLoadedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2FrameDOMContentLoadedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FrameDOMContentLoadedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FrameDOMContentLoadedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FrameDOMContentLoadedEventHandlerInvoke),
}

func NewICoreWebView2FrameDOMContentLoadedEventHandler(impl ICoreWebView2FrameDOMContentLoadedEventHandlerImpl) *ICoreWebView2FrameDOMContentLoadedEventHandler {
	return &ICoreWebView2FrameDOMContentLoadedEventHandler{
		Vtbl: &ICoreWebView2FrameDOMContentLoadedEventHandlerFn,
		impl: impl,
	}
}
