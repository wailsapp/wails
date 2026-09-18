//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2PrintToPdfStreamCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2PrintToPdfStreamCompletedHandler struct {
	Vtbl *ICoreWebView2PrintToPdfStreamCompletedHandlerVtbl
	impl ICoreWebView2PrintToPdfStreamCompletedHandlerImpl
}

func (i *ICoreWebView2PrintToPdfStreamCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2PrintToPdfStreamCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2PrintToPdfStreamCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2PrintToPdfStreamCompletedHandlerIUnknownAddRef(this *ICoreWebView2PrintToPdfStreamCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2PrintToPdfStreamCompletedHandlerIUnknownRelease(this *ICoreWebView2PrintToPdfStreamCompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2PrintToPdfStreamCompletedHandlerInvoke(this *ICoreWebView2PrintToPdfStreamCompletedHandler, errorCode uintptr, result *IStream) uintptr {
	return this.impl.PrintToPdfStreamCompleted(errorCode, result)
}

type ICoreWebView2PrintToPdfStreamCompletedHandlerImpl interface {
	IUnknownImpl
	PrintToPdfStreamCompleted(errorCode uintptr, result *IStream) uintptr
}

var ICoreWebView2PrintToPdfStreamCompletedHandlerFn = ICoreWebView2PrintToPdfStreamCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2PrintToPdfStreamCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2PrintToPdfStreamCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2PrintToPdfStreamCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2PrintToPdfStreamCompletedHandlerInvoke),
}

func NewICoreWebView2PrintToPdfStreamCompletedHandler(impl ICoreWebView2PrintToPdfStreamCompletedHandlerImpl) *ICoreWebView2PrintToPdfStreamCompletedHandler {
	return &ICoreWebView2PrintToPdfStreamCompletedHandler{
		Vtbl: &ICoreWebView2PrintToPdfStreamCompletedHandlerFn,
		impl: impl,
	}
}
