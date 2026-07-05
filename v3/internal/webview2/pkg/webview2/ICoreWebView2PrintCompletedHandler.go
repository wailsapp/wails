//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2PrintCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2PrintCompletedHandler struct {
	Vtbl *ICoreWebView2PrintCompletedHandlerVtbl
	impl ICoreWebView2PrintCompletedHandlerImpl
}

func (i *ICoreWebView2PrintCompletedHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2PrintCompletedHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2PrintCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2PrintCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2PrintCompletedHandlerIUnknownAddRef(this *ICoreWebView2PrintCompletedHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2PrintCompletedHandlerIUnknownRelease(this *ICoreWebView2PrintCompletedHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2PrintCompletedHandlerInvoke(this *ICoreWebView2PrintCompletedHandler, errorCode uintptr, result COREWEBVIEW2_PRINT_STATUS) uintptr {
	return this.impl.PrintCompleted(errorCode, result)
}

type ICoreWebView2PrintCompletedHandlerImpl interface {
	IUnknownImpl
	PrintCompleted(errorCode uintptr, result COREWEBVIEW2_PRINT_STATUS) uintptr
}

var ICoreWebView2PrintCompletedHandlerFn = ICoreWebView2PrintCompletedHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2PrintCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2PrintCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2PrintCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2PrintCompletedHandlerInvoke),
}

func NewICoreWebView2PrintCompletedHandler(impl ICoreWebView2PrintCompletedHandlerImpl) *ICoreWebView2PrintCompletedHandler {
	return &ICoreWebView2PrintCompletedHandler{
		Vtbl: &ICoreWebView2PrintCompletedHandlerFn,
		impl: impl,
	}
}
