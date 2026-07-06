//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2FaviconChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FaviconChangedEventHandler struct {
	Vtbl *ICoreWebView2FaviconChangedEventHandlerVtbl
	impl ICoreWebView2FaviconChangedEventHandlerImpl
}

func (i *ICoreWebView2FaviconChangedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2FaviconChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2FaviconChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FaviconChangedEventHandlerIUnknownAddRef(this *ICoreWebView2FaviconChangedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2FaviconChangedEventHandlerIUnknownRelease(this *ICoreWebView2FaviconChangedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2FaviconChangedEventHandlerInvoke(this *ICoreWebView2FaviconChangedEventHandler, sender *ICoreWebView2, args *IUnknown) uintptr {
	return this.impl.FaviconChanged(sender, args)
}

type ICoreWebView2FaviconChangedEventHandlerImpl interface {
	IUnknownImpl
	FaviconChanged(sender *ICoreWebView2, args *IUnknown) uintptr
}

var ICoreWebView2FaviconChangedEventHandlerFn = ICoreWebView2FaviconChangedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2FaviconChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FaviconChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FaviconChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FaviconChangedEventHandlerInvoke),
}

func NewICoreWebView2FaviconChangedEventHandler(impl ICoreWebView2FaviconChangedEventHandlerImpl) *ICoreWebView2FaviconChangedEventHandler {
	return &ICoreWebView2FaviconChangedEventHandler{
		Vtbl: &ICoreWebView2FaviconChangedEventHandlerFn,
		impl: impl,
	}
}
