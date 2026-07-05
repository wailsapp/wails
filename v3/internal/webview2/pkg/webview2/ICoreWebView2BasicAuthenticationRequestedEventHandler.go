//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2BasicAuthenticationRequestedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2BasicAuthenticationRequestedEventHandler struct {
	Vtbl *ICoreWebView2BasicAuthenticationRequestedEventHandlerVtbl
	impl ICoreWebView2BasicAuthenticationRequestedEventHandlerImpl
}

func (i *ICoreWebView2BasicAuthenticationRequestedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2BasicAuthenticationRequestedEventHandlerIUnknownQueryInterface(this *ICoreWebView2BasicAuthenticationRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2BasicAuthenticationRequestedEventHandlerIUnknownAddRef(this *ICoreWebView2BasicAuthenticationRequestedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2BasicAuthenticationRequestedEventHandlerIUnknownRelease(this *ICoreWebView2BasicAuthenticationRequestedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2BasicAuthenticationRequestedEventHandlerInvoke(this *ICoreWebView2BasicAuthenticationRequestedEventHandler, sender *ICoreWebView2, args *ICoreWebView2BasicAuthenticationRequestedEventArgs) uintptr {
	return this.impl.BasicAuthenticationRequested(sender, args)
}

type ICoreWebView2BasicAuthenticationRequestedEventHandlerImpl interface {
	IUnknownImpl
	BasicAuthenticationRequested(sender *ICoreWebView2, args *ICoreWebView2BasicAuthenticationRequestedEventArgs) uintptr
}

var ICoreWebView2BasicAuthenticationRequestedEventHandlerFn = ICoreWebView2BasicAuthenticationRequestedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2BasicAuthenticationRequestedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2BasicAuthenticationRequestedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2BasicAuthenticationRequestedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2BasicAuthenticationRequestedEventHandlerInvoke),
}

func NewICoreWebView2BasicAuthenticationRequestedEventHandler(impl ICoreWebView2BasicAuthenticationRequestedEventHandlerImpl) *ICoreWebView2BasicAuthenticationRequestedEventHandler {
	return &ICoreWebView2BasicAuthenticationRequestedEventHandler{
		Vtbl: &ICoreWebView2BasicAuthenticationRequestedEventHandlerFn,
		impl: impl,
	}
}
