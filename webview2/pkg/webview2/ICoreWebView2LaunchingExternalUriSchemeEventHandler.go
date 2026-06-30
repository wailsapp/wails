//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2LaunchingExternalUriSchemeEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2LaunchingExternalUriSchemeEventHandler struct {
	Vtbl *ICoreWebView2LaunchingExternalUriSchemeEventHandlerVtbl
	impl ICoreWebView2LaunchingExternalUriSchemeEventHandlerImpl
}

func (i *ICoreWebView2LaunchingExternalUriSchemeEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2LaunchingExternalUriSchemeEventHandlerIUnknownQueryInterface(this *ICoreWebView2LaunchingExternalUriSchemeEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2LaunchingExternalUriSchemeEventHandlerIUnknownAddRef(this *ICoreWebView2LaunchingExternalUriSchemeEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2LaunchingExternalUriSchemeEventHandlerIUnknownRelease(this *ICoreWebView2LaunchingExternalUriSchemeEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2LaunchingExternalUriSchemeEventHandlerInvoke(this *ICoreWebView2LaunchingExternalUriSchemeEventHandler, sender *ICoreWebView2, args *ICoreWebView2LaunchingExternalUriSchemeEventArgs) uintptr {
	return this.impl.LaunchingExternalUriScheme(sender, args)
}

type ICoreWebView2LaunchingExternalUriSchemeEventHandlerImpl interface {
	IUnknownImpl
	LaunchingExternalUriScheme(sender *ICoreWebView2, args *ICoreWebView2LaunchingExternalUriSchemeEventArgs) uintptr
}

var ICoreWebView2LaunchingExternalUriSchemeEventHandlerFn = ICoreWebView2LaunchingExternalUriSchemeEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2LaunchingExternalUriSchemeEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2LaunchingExternalUriSchemeEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2LaunchingExternalUriSchemeEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2LaunchingExternalUriSchemeEventHandlerInvoke),
}

func NewICoreWebView2LaunchingExternalUriSchemeEventHandler(impl ICoreWebView2LaunchingExternalUriSchemeEventHandlerImpl) *ICoreWebView2LaunchingExternalUriSchemeEventHandler {
	return &ICoreWebView2LaunchingExternalUriSchemeEventHandler{
		Vtbl: &ICoreWebView2LaunchingExternalUriSchemeEventHandlerFn,
		impl: impl,
	}
}
