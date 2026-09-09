//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2PermissionRequestedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2PermissionRequestedEventHandler struct {
	Vtbl *ICoreWebView2PermissionRequestedEventHandlerVtbl
	impl ICoreWebView2PermissionRequestedEventHandlerImpl
}

func (i *ICoreWebView2PermissionRequestedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2PermissionRequestedEventHandlerIUnknownQueryInterface(this *ICoreWebView2PermissionRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2PermissionRequestedEventHandlerIUnknownAddRef(this *ICoreWebView2PermissionRequestedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2PermissionRequestedEventHandlerIUnknownRelease(this *ICoreWebView2PermissionRequestedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2PermissionRequestedEventHandlerInvoke(this *ICoreWebView2PermissionRequestedEventHandler, sender *ICoreWebView2, args *ICoreWebView2PermissionRequestedEventArgs) uintptr {
	return this.impl.PermissionRequested(sender, args)
}

type ICoreWebView2PermissionRequestedEventHandlerImpl interface {
	IUnknownImpl
	PermissionRequested(sender *ICoreWebView2, args *ICoreWebView2PermissionRequestedEventArgs) uintptr
}

var ICoreWebView2PermissionRequestedEventHandlerFn = ICoreWebView2PermissionRequestedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2PermissionRequestedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2PermissionRequestedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2PermissionRequestedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2PermissionRequestedEventHandlerInvoke),
}

func NewICoreWebView2PermissionRequestedEventHandler(impl ICoreWebView2PermissionRequestedEventHandlerImpl) *ICoreWebView2PermissionRequestedEventHandler {
	return &ICoreWebView2PermissionRequestedEventHandler{
		Vtbl: &ICoreWebView2PermissionRequestedEventHandlerFn,
		impl: impl,
	}
}
