//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2FindStartCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FindStartCompletedHandler struct {
	Vtbl *ICoreWebView2FindStartCompletedHandlerVtbl
	impl ICoreWebView2FindStartCompletedHandlerImpl
}

func (i *ICoreWebView2FindStartCompletedHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2FindStartCompletedHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2FindStartCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2FindStartCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FindStartCompletedHandlerIUnknownAddRef(this *ICoreWebView2FindStartCompletedHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2FindStartCompletedHandlerIUnknownRelease(this *ICoreWebView2FindStartCompletedHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2FindStartCompletedHandlerInvoke(this *ICoreWebView2FindStartCompletedHandler, errorCode uintptr) uintptr {
	return this.impl.FindStartCompleted(errorCode)
}

type ICoreWebView2FindStartCompletedHandlerImpl interface {
	IUnknownImpl
	FindStartCompleted(errorCode uintptr) uintptr
}

var ICoreWebView2FindStartCompletedHandlerFn = ICoreWebView2FindStartCompletedHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2FindStartCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FindStartCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FindStartCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FindStartCompletedHandlerInvoke),
}

func NewICoreWebView2FindStartCompletedHandler(impl ICoreWebView2FindStartCompletedHandlerImpl) *ICoreWebView2FindStartCompletedHandler {
	return &ICoreWebView2FindStartCompletedHandler{
		Vtbl: &ICoreWebView2FindStartCompletedHandlerFn,
		impl: impl,
	}
}
