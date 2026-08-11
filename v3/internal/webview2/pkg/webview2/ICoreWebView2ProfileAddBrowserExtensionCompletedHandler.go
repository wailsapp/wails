//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ProfileAddBrowserExtensionCompletedHandler struct {
	Vtbl *ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerVtbl
	impl ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerImpl
}

func (i *ICoreWebView2ProfileAddBrowserExtensionCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2ProfileAddBrowserExtensionCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerIUnknownAddRef(this *ICoreWebView2ProfileAddBrowserExtensionCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerIUnknownRelease(this *ICoreWebView2ProfileAddBrowserExtensionCompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerInvoke(this *ICoreWebView2ProfileAddBrowserExtensionCompletedHandler, errorCode uintptr, result *ICoreWebView2BrowserExtension) uintptr {
	return this.impl.ProfileAddBrowserExtensionCompleted(errorCode, result)
}

type ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerImpl interface {
	IUnknownImpl
	ProfileAddBrowserExtensionCompleted(errorCode uintptr, result *ICoreWebView2BrowserExtension) uintptr
}

var ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerFn = ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerInvoke),
}

func NewICoreWebView2ProfileAddBrowserExtensionCompletedHandler(impl ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerImpl) *ICoreWebView2ProfileAddBrowserExtensionCompletedHandler {
	return &ICoreWebView2ProfileAddBrowserExtensionCompletedHandler{
		Vtbl: &ICoreWebView2ProfileAddBrowserExtensionCompletedHandlerFn,
		impl: impl,
	}
}
