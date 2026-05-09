//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2BrowserExtensionEnableCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2BrowserExtensionEnableCompletedHandler struct {
	Vtbl *ICoreWebView2BrowserExtensionEnableCompletedHandlerVtbl
	impl ICoreWebView2BrowserExtensionEnableCompletedHandlerImpl
}

func (i *ICoreWebView2BrowserExtensionEnableCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2BrowserExtensionEnableCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2BrowserExtensionEnableCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2BrowserExtensionEnableCompletedHandlerIUnknownAddRef(this *ICoreWebView2BrowserExtensionEnableCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2BrowserExtensionEnableCompletedHandlerIUnknownRelease(this *ICoreWebView2BrowserExtensionEnableCompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2BrowserExtensionEnableCompletedHandlerInvoke(this *ICoreWebView2BrowserExtensionEnableCompletedHandler, errorCode uintptr) uintptr {
	return this.impl.BrowserExtensionEnableCompleted(errorCode)
}

type ICoreWebView2BrowserExtensionEnableCompletedHandlerImpl interface {
	IUnknownImpl
	BrowserExtensionEnableCompleted(errorCode uintptr) uintptr
}

var ICoreWebView2BrowserExtensionEnableCompletedHandlerFn = ICoreWebView2BrowserExtensionEnableCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2BrowserExtensionEnableCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2BrowserExtensionEnableCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2BrowserExtensionEnableCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2BrowserExtensionEnableCompletedHandlerInvoke),
}

func NewICoreWebView2BrowserExtensionEnableCompletedHandler(impl ICoreWebView2BrowserExtensionEnableCompletedHandlerImpl) *ICoreWebView2BrowserExtensionEnableCompletedHandler {
	return &ICoreWebView2BrowserExtensionEnableCompletedHandler{
		Vtbl: &ICoreWebView2BrowserExtensionEnableCompletedHandlerFn,
		impl: impl,
	}
}
