//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2GetFaviconCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2GetFaviconCompletedHandler struct {
	Vtbl *ICoreWebView2GetFaviconCompletedHandlerVtbl
	impl ICoreWebView2GetFaviconCompletedHandlerImpl
}

func (i *ICoreWebView2GetFaviconCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2GetFaviconCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2GetFaviconCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2GetFaviconCompletedHandlerIUnknownAddRef(this *ICoreWebView2GetFaviconCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2GetFaviconCompletedHandlerIUnknownRelease(this *ICoreWebView2GetFaviconCompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2GetFaviconCompletedHandlerInvoke(this *ICoreWebView2GetFaviconCompletedHandler, errorCode uintptr, result *IStream) uintptr {
	return this.impl.GetFaviconCompleted(errorCode, result)
}

type ICoreWebView2GetFaviconCompletedHandlerImpl interface {
	IUnknownImpl
	GetFaviconCompleted(errorCode uintptr, result *IStream) uintptr
}

var ICoreWebView2GetFaviconCompletedHandlerFn = ICoreWebView2GetFaviconCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2GetFaviconCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2GetFaviconCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2GetFaviconCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2GetFaviconCompletedHandlerInvoke),
}

func NewICoreWebView2GetFaviconCompletedHandler(impl ICoreWebView2GetFaviconCompletedHandlerImpl) *ICoreWebView2GetFaviconCompletedHandler {
	return &ICoreWebView2GetFaviconCompletedHandler{
		Vtbl: &ICoreWebView2GetFaviconCompletedHandlerFn,
		impl: impl,
	}
}
