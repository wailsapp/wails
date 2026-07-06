//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2ClientCertificateRequestedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ClientCertificateRequestedEventHandler struct {
	Vtbl *ICoreWebView2ClientCertificateRequestedEventHandlerVtbl
	impl ICoreWebView2ClientCertificateRequestedEventHandlerImpl
}

func (i *ICoreWebView2ClientCertificateRequestedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2ClientCertificateRequestedEventHandlerIUnknownQueryInterface(this *ICoreWebView2ClientCertificateRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ClientCertificateRequestedEventHandlerIUnknownAddRef(this *ICoreWebView2ClientCertificateRequestedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2ClientCertificateRequestedEventHandlerIUnknownRelease(this *ICoreWebView2ClientCertificateRequestedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2ClientCertificateRequestedEventHandlerInvoke(this *ICoreWebView2ClientCertificateRequestedEventHandler, sender *ICoreWebView2, args *ICoreWebView2ClientCertificateRequestedEventArgs) uintptr {
	return this.impl.ClientCertificateRequested(sender, args)
}

type ICoreWebView2ClientCertificateRequestedEventHandlerImpl interface {
	IUnknownImpl
	ClientCertificateRequested(sender *ICoreWebView2, args *ICoreWebView2ClientCertificateRequestedEventArgs) uintptr
}

var ICoreWebView2ClientCertificateRequestedEventHandlerFn = ICoreWebView2ClientCertificateRequestedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2ClientCertificateRequestedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ClientCertificateRequestedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ClientCertificateRequestedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ClientCertificateRequestedEventHandlerInvoke),
}

func NewICoreWebView2ClientCertificateRequestedEventHandler(impl ICoreWebView2ClientCertificateRequestedEventHandlerImpl) *ICoreWebView2ClientCertificateRequestedEventHandler {
	return &ICoreWebView2ClientCertificateRequestedEventHandler{
		Vtbl: &ICoreWebView2ClientCertificateRequestedEventHandlerFn,
		impl: impl,
	}
}
