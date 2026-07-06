//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2ContentLoadingEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ContentLoadingEventHandler struct {
	Vtbl *ICoreWebView2ContentLoadingEventHandlerVtbl
	impl ICoreWebView2ContentLoadingEventHandlerImpl
}

func (i *ICoreWebView2ContentLoadingEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2ContentLoadingEventHandlerIUnknownQueryInterface(this *ICoreWebView2ContentLoadingEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ContentLoadingEventHandlerIUnknownAddRef(this *ICoreWebView2ContentLoadingEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2ContentLoadingEventHandlerIUnknownRelease(this *ICoreWebView2ContentLoadingEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2ContentLoadingEventHandlerInvoke(this *ICoreWebView2ContentLoadingEventHandler, sender *ICoreWebView2, args *ICoreWebView2ContentLoadingEventArgs) uintptr {
	return this.impl.ContentLoading(sender, args)
}

type ICoreWebView2ContentLoadingEventHandlerImpl interface {
	IUnknownImpl
	ContentLoading(sender *ICoreWebView2, args *ICoreWebView2ContentLoadingEventArgs) uintptr
}

var ICoreWebView2ContentLoadingEventHandlerFn = ICoreWebView2ContentLoadingEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2ContentLoadingEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ContentLoadingEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ContentLoadingEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ContentLoadingEventHandlerInvoke),
}

func NewICoreWebView2ContentLoadingEventHandler(impl ICoreWebView2ContentLoadingEventHandlerImpl) *ICoreWebView2ContentLoadingEventHandler {
	return &ICoreWebView2ContentLoadingEventHandler{
		Vtbl: &ICoreWebView2ContentLoadingEventHandlerFn,
		impl: impl,
	}
}
