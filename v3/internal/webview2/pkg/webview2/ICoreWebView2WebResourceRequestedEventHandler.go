//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2WebResourceRequestedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2WebResourceRequestedEventHandler struct {
	Vtbl *ICoreWebView2WebResourceRequestedEventHandlerVtbl
	impl ICoreWebView2WebResourceRequestedEventHandlerImpl
}

func (i *ICoreWebView2WebResourceRequestedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2WebResourceRequestedEventHandlerIUnknownQueryInterface(this *ICoreWebView2WebResourceRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2WebResourceRequestedEventHandlerIUnknownAddRef(this *ICoreWebView2WebResourceRequestedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2WebResourceRequestedEventHandlerIUnknownRelease(this *ICoreWebView2WebResourceRequestedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2WebResourceRequestedEventHandlerInvoke(this *ICoreWebView2WebResourceRequestedEventHandler, sender *ICoreWebView2, args *ICoreWebView2WebResourceRequestedEventArgs) uintptr {
	return this.impl.WebResourceRequested(sender, args)
}

type ICoreWebView2WebResourceRequestedEventHandlerImpl interface {
	IUnknownImpl
	WebResourceRequested(sender *ICoreWebView2, args *ICoreWebView2WebResourceRequestedEventArgs) uintptr
}

var ICoreWebView2WebResourceRequestedEventHandlerFn = ICoreWebView2WebResourceRequestedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2WebResourceRequestedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2WebResourceRequestedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2WebResourceRequestedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2WebResourceRequestedEventHandlerInvoke),
}

func NewICoreWebView2WebResourceRequestedEventHandler(impl ICoreWebView2WebResourceRequestedEventHandlerImpl) *ICoreWebView2WebResourceRequestedEventHandler {
	return &ICoreWebView2WebResourceRequestedEventHandler{
		Vtbl: &ICoreWebView2WebResourceRequestedEventHandlerFn,
		impl: impl,
	}
}
