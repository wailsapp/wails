//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2SetPermissionStateCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2SetPermissionStateCompletedHandler struct {
	Vtbl *ICoreWebView2SetPermissionStateCompletedHandlerVtbl
	impl ICoreWebView2SetPermissionStateCompletedHandlerImpl
}

func (i *ICoreWebView2SetPermissionStateCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2SetPermissionStateCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2SetPermissionStateCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2SetPermissionStateCompletedHandlerIUnknownAddRef(this *ICoreWebView2SetPermissionStateCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2SetPermissionStateCompletedHandlerIUnknownRelease(this *ICoreWebView2SetPermissionStateCompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2SetPermissionStateCompletedHandlerInvoke(this *ICoreWebView2SetPermissionStateCompletedHandler, errorCode uintptr) uintptr {
	return this.impl.SetPermissionStateCompleted(errorCode)
}

type ICoreWebView2SetPermissionStateCompletedHandlerImpl interface {
	IUnknownImpl
	SetPermissionStateCompleted(errorCode uintptr) uintptr
}

var ICoreWebView2SetPermissionStateCompletedHandlerFn = ICoreWebView2SetPermissionStateCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2SetPermissionStateCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2SetPermissionStateCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2SetPermissionStateCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2SetPermissionStateCompletedHandlerInvoke),
}

func NewICoreWebView2SetPermissionStateCompletedHandler(impl ICoreWebView2SetPermissionStateCompletedHandlerImpl) *ICoreWebView2SetPermissionStateCompletedHandler {
	return &ICoreWebView2SetPermissionStateCompletedHandler{
		Vtbl: &ICoreWebView2SetPermissionStateCompletedHandlerFn,
		impl: impl,
	}
}
