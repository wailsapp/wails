//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ProfileGetBrowserExtensionsCompletedHandler struct {
	Vtbl *ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerVtbl
	impl ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerImpl
}

func (i *ICoreWebView2ProfileGetBrowserExtensionsCompletedHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ProfileGetBrowserExtensionsCompletedHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2ProfileGetBrowserExtensionsCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerIUnknownAddRef(this *ICoreWebView2ProfileGetBrowserExtensionsCompletedHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerIUnknownRelease(this *ICoreWebView2ProfileGetBrowserExtensionsCompletedHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerInvoke(this *ICoreWebView2ProfileGetBrowserExtensionsCompletedHandler, errorCode uintptr, result *ICoreWebView2BrowserExtensionList) uintptr {
	return this.impl.ProfileGetBrowserExtensionsCompleted(errorCode, result)
}

type ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerImpl interface {
	IUnknownImpl
	ProfileGetBrowserExtensionsCompleted(errorCode uintptr, result *ICoreWebView2BrowserExtensionList) uintptr
}

var ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerFn = ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerInvoke),
}

func NewICoreWebView2ProfileGetBrowserExtensionsCompletedHandler(impl ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerImpl) *ICoreWebView2ProfileGetBrowserExtensionsCompletedHandler {
	return &ICoreWebView2ProfileGetBrowserExtensionsCompletedHandler{
		Vtbl: &ICoreWebView2ProfileGetBrowserExtensionsCompletedHandlerFn,
		impl: impl,
	}
}
