//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2ServerCertificateErrorDetectedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ServerCertificateErrorDetectedEventHandler struct {
	Vtbl *ICoreWebView2ServerCertificateErrorDetectedEventHandlerVtbl
	impl ICoreWebView2ServerCertificateErrorDetectedEventHandlerImpl
}

func (i *ICoreWebView2ServerCertificateErrorDetectedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2ServerCertificateErrorDetectedEventHandlerIUnknownQueryInterface(this *ICoreWebView2ServerCertificateErrorDetectedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ServerCertificateErrorDetectedEventHandlerIUnknownAddRef(this *ICoreWebView2ServerCertificateErrorDetectedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2ServerCertificateErrorDetectedEventHandlerIUnknownRelease(this *ICoreWebView2ServerCertificateErrorDetectedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2ServerCertificateErrorDetectedEventHandlerInvoke(this *ICoreWebView2ServerCertificateErrorDetectedEventHandler, sender *ICoreWebView2, args *ICoreWebView2ServerCertificateErrorDetectedEventArgs) uintptr {
	return this.impl.ServerCertificateErrorDetected(sender, args)
}

type ICoreWebView2ServerCertificateErrorDetectedEventHandlerImpl interface {
	IUnknownImpl
	ServerCertificateErrorDetected(sender *ICoreWebView2, args *ICoreWebView2ServerCertificateErrorDetectedEventArgs) uintptr
}

var ICoreWebView2ServerCertificateErrorDetectedEventHandlerFn = ICoreWebView2ServerCertificateErrorDetectedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2ServerCertificateErrorDetectedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ServerCertificateErrorDetectedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ServerCertificateErrorDetectedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ServerCertificateErrorDetectedEventHandlerInvoke),
}

func NewICoreWebView2ServerCertificateErrorDetectedEventHandler(impl ICoreWebView2ServerCertificateErrorDetectedEventHandlerImpl) *ICoreWebView2ServerCertificateErrorDetectedEventHandler {
	return &ICoreWebView2ServerCertificateErrorDetectedEventHandler{
		Vtbl: &ICoreWebView2ServerCertificateErrorDetectedEventHandlerFn,
		impl: impl,
	}
}
