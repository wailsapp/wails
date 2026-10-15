//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ClearServerCertificateErrorActionsCompletedHandler struct {
	Vtbl *ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerVtbl
	impl ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerImpl
}

func (i *ICoreWebView2ClearServerCertificateErrorActionsCompletedHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ClearServerCertificateErrorActionsCompletedHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2ClearServerCertificateErrorActionsCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerIUnknownAddRef(this *ICoreWebView2ClearServerCertificateErrorActionsCompletedHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerIUnknownRelease(this *ICoreWebView2ClearServerCertificateErrorActionsCompletedHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerInvoke(this *ICoreWebView2ClearServerCertificateErrorActionsCompletedHandler, errorCode uintptr) uintptr {
	return this.impl.ClearServerCertificateErrorActionsCompleted(errorCode)
}

type ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerImpl interface {
	IUnknownImpl
	ClearServerCertificateErrorActionsCompleted(errorCode uintptr) uintptr
}

var ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerFn = ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerInvoke),
}

func NewICoreWebView2ClearServerCertificateErrorActionsCompletedHandler(impl ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerImpl) *ICoreWebView2ClearServerCertificateErrorActionsCompletedHandler {
	return &ICoreWebView2ClearServerCertificateErrorActionsCompletedHandler{
		Vtbl: &ICoreWebView2ClearServerCertificateErrorActionsCompletedHandlerFn,
		impl: impl,
	}
}
