//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2GetProcessExtendedInfosCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2GetProcessExtendedInfosCompletedHandler struct {
	Vtbl *ICoreWebView2GetProcessExtendedInfosCompletedHandlerVtbl
	impl ICoreWebView2GetProcessExtendedInfosCompletedHandlerImpl
}

func (i *ICoreWebView2GetProcessExtendedInfosCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2GetProcessExtendedInfosCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2GetProcessExtendedInfosCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2GetProcessExtendedInfosCompletedHandlerIUnknownAddRef(this *ICoreWebView2GetProcessExtendedInfosCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2GetProcessExtendedInfosCompletedHandlerIUnknownRelease(this *ICoreWebView2GetProcessExtendedInfosCompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2GetProcessExtendedInfosCompletedHandlerInvoke(this *ICoreWebView2GetProcessExtendedInfosCompletedHandler, errorCode uintptr, result *ICoreWebView2ProcessExtendedInfoCollection) uintptr {
	return this.impl.GetProcessExtendedInfosCompleted(errorCode, result)
}

type ICoreWebView2GetProcessExtendedInfosCompletedHandlerImpl interface {
	IUnknownImpl
	GetProcessExtendedInfosCompleted(errorCode uintptr, result *ICoreWebView2ProcessExtendedInfoCollection) uintptr
}

var ICoreWebView2GetProcessExtendedInfosCompletedHandlerFn = ICoreWebView2GetProcessExtendedInfosCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2GetProcessExtendedInfosCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2GetProcessExtendedInfosCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2GetProcessExtendedInfosCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2GetProcessExtendedInfosCompletedHandlerInvoke),
}

func NewICoreWebView2GetProcessExtendedInfosCompletedHandler(impl ICoreWebView2GetProcessExtendedInfosCompletedHandlerImpl) *ICoreWebView2GetProcessExtendedInfosCompletedHandler {
	return &ICoreWebView2GetProcessExtendedInfosCompletedHandler{
		Vtbl: &ICoreWebView2GetProcessExtendedInfosCompletedHandlerFn,
		impl: impl,
	}
}
