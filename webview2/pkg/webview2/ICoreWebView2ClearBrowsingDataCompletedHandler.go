//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2ClearBrowsingDataCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ClearBrowsingDataCompletedHandler struct {
	Vtbl *ICoreWebView2ClearBrowsingDataCompletedHandlerVtbl
	impl ICoreWebView2ClearBrowsingDataCompletedHandlerImpl
}

func (i *ICoreWebView2ClearBrowsingDataCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2ClearBrowsingDataCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2ClearBrowsingDataCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ClearBrowsingDataCompletedHandlerIUnknownAddRef(this *ICoreWebView2ClearBrowsingDataCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2ClearBrowsingDataCompletedHandlerIUnknownRelease(this *ICoreWebView2ClearBrowsingDataCompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2ClearBrowsingDataCompletedHandlerInvoke(this *ICoreWebView2ClearBrowsingDataCompletedHandler, errorCode uintptr) uintptr {
	return this.impl.ClearBrowsingDataCompleted(errorCode)
}

type ICoreWebView2ClearBrowsingDataCompletedHandlerImpl interface {
	IUnknownImpl
	ClearBrowsingDataCompleted(errorCode uintptr) uintptr
}

var ICoreWebView2ClearBrowsingDataCompletedHandlerFn = ICoreWebView2ClearBrowsingDataCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2ClearBrowsingDataCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ClearBrowsingDataCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ClearBrowsingDataCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ClearBrowsingDataCompletedHandlerInvoke),
}

func NewICoreWebView2ClearBrowsingDataCompletedHandler(impl ICoreWebView2ClearBrowsingDataCompletedHandlerImpl) *ICoreWebView2ClearBrowsingDataCompletedHandler {
	return &ICoreWebView2ClearBrowsingDataCompletedHandler{
		Vtbl: &ICoreWebView2ClearBrowsingDataCompletedHandlerFn,
		impl: impl,
	}
}
