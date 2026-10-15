//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2BrowserExtensionRemoveCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2BrowserExtensionRemoveCompletedHandler struct {
	Vtbl *ICoreWebView2BrowserExtensionRemoveCompletedHandlerVtbl
	impl ICoreWebView2BrowserExtensionRemoveCompletedHandlerImpl
}

func (i *ICoreWebView2BrowserExtensionRemoveCompletedHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2BrowserExtensionRemoveCompletedHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2BrowserExtensionRemoveCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2BrowserExtensionRemoveCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2BrowserExtensionRemoveCompletedHandlerIUnknownAddRef(this *ICoreWebView2BrowserExtensionRemoveCompletedHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2BrowserExtensionRemoveCompletedHandlerIUnknownRelease(this *ICoreWebView2BrowserExtensionRemoveCompletedHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2BrowserExtensionRemoveCompletedHandlerInvoke(this *ICoreWebView2BrowserExtensionRemoveCompletedHandler, errorCode uintptr) uintptr {
	return this.impl.BrowserExtensionRemoveCompleted(errorCode)
}

type ICoreWebView2BrowserExtensionRemoveCompletedHandlerImpl interface {
	IUnknownImpl
	BrowserExtensionRemoveCompleted(errorCode uintptr) uintptr
}

var ICoreWebView2BrowserExtensionRemoveCompletedHandlerFn = ICoreWebView2BrowserExtensionRemoveCompletedHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2BrowserExtensionRemoveCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2BrowserExtensionRemoveCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2BrowserExtensionRemoveCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2BrowserExtensionRemoveCompletedHandlerInvoke),
}

func NewICoreWebView2BrowserExtensionRemoveCompletedHandler(impl ICoreWebView2BrowserExtensionRemoveCompletedHandlerImpl) *ICoreWebView2BrowserExtensionRemoveCompletedHandler {
	return &ICoreWebView2BrowserExtensionRemoveCompletedHandler{
		Vtbl: &ICoreWebView2BrowserExtensionRemoveCompletedHandlerFn,
		impl: impl,
	}
}
