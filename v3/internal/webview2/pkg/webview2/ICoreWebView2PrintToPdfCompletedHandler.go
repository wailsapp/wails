//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2PrintToPdfCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2PrintToPdfCompletedHandler struct {
	Vtbl *ICoreWebView2PrintToPdfCompletedHandlerVtbl
	impl ICoreWebView2PrintToPdfCompletedHandlerImpl
}

func (i *ICoreWebView2PrintToPdfCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2PrintToPdfCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2PrintToPdfCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2PrintToPdfCompletedHandlerIUnknownAddRef(this *ICoreWebView2PrintToPdfCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2PrintToPdfCompletedHandlerIUnknownRelease(this *ICoreWebView2PrintToPdfCompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2PrintToPdfCompletedHandlerInvoke(this *ICoreWebView2PrintToPdfCompletedHandler, errorCode uintptr, result bool) uintptr {
	return this.impl.PrintToPdfCompleted(errorCode, result)
}

type ICoreWebView2PrintToPdfCompletedHandlerImpl interface {
	IUnknownImpl
	PrintToPdfCompleted(errorCode uintptr, result bool) uintptr
}

var ICoreWebView2PrintToPdfCompletedHandlerFn = ICoreWebView2PrintToPdfCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2PrintToPdfCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2PrintToPdfCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2PrintToPdfCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2PrintToPdfCompletedHandlerInvoke),
}

func NewICoreWebView2PrintToPdfCompletedHandler(impl ICoreWebView2PrintToPdfCompletedHandlerImpl) *ICoreWebView2PrintToPdfCompletedHandler {
	return &ICoreWebView2PrintToPdfCompletedHandler{
		Vtbl: &ICoreWebView2PrintToPdfCompletedHandlerFn,
		impl: impl,
	}
}
