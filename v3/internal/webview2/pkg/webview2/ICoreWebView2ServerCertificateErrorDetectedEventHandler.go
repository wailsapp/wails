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

func (i *ICoreWebView2ServerCertificateErrorDetectedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ServerCertificateErrorDetectedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2ServerCertificateErrorDetectedEventHandlerIUnknownQueryInterface(this *ICoreWebView2ServerCertificateErrorDetectedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ServerCertificateErrorDetectedEventHandlerIUnknownAddRef(this *ICoreWebView2ServerCertificateErrorDetectedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2ServerCertificateErrorDetectedEventHandlerIUnknownRelease(this *ICoreWebView2ServerCertificateErrorDetectedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2ServerCertificateErrorDetectedEventHandlerInvoke(this *ICoreWebView2ServerCertificateErrorDetectedEventHandler, sender *ICoreWebView2, args *ICoreWebView2ServerCertificateErrorDetectedEventArgs) uintptr {
	return this.impl.ServerCertificateErrorDetected(sender, args)
}

type ICoreWebView2ServerCertificateErrorDetectedEventHandlerImpl interface {
	IUnknownImpl
	ServerCertificateErrorDetected(sender *ICoreWebView2, args *ICoreWebView2ServerCertificateErrorDetectedEventArgs) uintptr
}

var ICoreWebView2ServerCertificateErrorDetectedEventHandlerFn = ICoreWebView2ServerCertificateErrorDetectedEventHandlerVtbl{
	IUnknownVtbl {
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
