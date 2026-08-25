//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2ProcessFailedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ProcessFailedEventHandler struct {
	Vtbl *ICoreWebView2ProcessFailedEventHandlerVtbl
	impl ICoreWebView2ProcessFailedEventHandlerImpl
}

func (i *ICoreWebView2ProcessFailedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2ProcessFailedEventHandlerIUnknownQueryInterface(this *ICoreWebView2ProcessFailedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ProcessFailedEventHandlerIUnknownAddRef(this *ICoreWebView2ProcessFailedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2ProcessFailedEventHandlerIUnknownRelease(this *ICoreWebView2ProcessFailedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2ProcessFailedEventHandlerInvoke(this *ICoreWebView2ProcessFailedEventHandler, sender *ICoreWebView2, args *ICoreWebView2ProcessFailedEventArgs) uintptr {
	return this.impl.ProcessFailed(sender, args)
}

type ICoreWebView2ProcessFailedEventHandlerImpl interface {
	IUnknownImpl
	ProcessFailed(sender *ICoreWebView2, args *ICoreWebView2ProcessFailedEventArgs) uintptr
}

var ICoreWebView2ProcessFailedEventHandlerFn = ICoreWebView2ProcessFailedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2ProcessFailedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ProcessFailedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ProcessFailedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ProcessFailedEventHandlerInvoke),
}

func NewICoreWebView2ProcessFailedEventHandler(impl ICoreWebView2ProcessFailedEventHandlerImpl) *ICoreWebView2ProcessFailedEventHandler {
	return &ICoreWebView2ProcessFailedEventHandler{
		Vtbl: &ICoreWebView2ProcessFailedEventHandlerFn,
		impl: impl,
	}
}
