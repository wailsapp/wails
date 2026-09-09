//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2GetCookiesCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2GetCookiesCompletedHandler struct {
	Vtbl *ICoreWebView2GetCookiesCompletedHandlerVtbl
	impl ICoreWebView2GetCookiesCompletedHandlerImpl
}

func (i *ICoreWebView2GetCookiesCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2GetCookiesCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2GetCookiesCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2GetCookiesCompletedHandlerIUnknownAddRef(this *ICoreWebView2GetCookiesCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2GetCookiesCompletedHandlerIUnknownRelease(this *ICoreWebView2GetCookiesCompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2GetCookiesCompletedHandlerInvoke(this *ICoreWebView2GetCookiesCompletedHandler, errorCode uintptr, result *ICoreWebView2CookieList) uintptr {
	return this.impl.GetCookiesCompleted(errorCode, result)
}

type ICoreWebView2GetCookiesCompletedHandlerImpl interface {
	IUnknownImpl
	GetCookiesCompleted(errorCode uintptr, result *ICoreWebView2CookieList) uintptr
}

var ICoreWebView2GetCookiesCompletedHandlerFn = ICoreWebView2GetCookiesCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2GetCookiesCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2GetCookiesCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2GetCookiesCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2GetCookiesCompletedHandlerInvoke),
}

func NewICoreWebView2GetCookiesCompletedHandler(impl ICoreWebView2GetCookiesCompletedHandlerImpl) *ICoreWebView2GetCookiesCompletedHandler {
	return &ICoreWebView2GetCookiesCompletedHandler{
		Vtbl: &ICoreWebView2GetCookiesCompletedHandlerFn,
		impl: impl,
	}
}
