//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2FramePermissionRequestedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FramePermissionRequestedEventHandler struct {
	Vtbl *ICoreWebView2FramePermissionRequestedEventHandlerVtbl
	impl ICoreWebView2FramePermissionRequestedEventHandlerImpl
}

func (i *ICoreWebView2FramePermissionRequestedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2FramePermissionRequestedEventHandlerIUnknownQueryInterface(this *ICoreWebView2FramePermissionRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FramePermissionRequestedEventHandlerIUnknownAddRef(this *ICoreWebView2FramePermissionRequestedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2FramePermissionRequestedEventHandlerIUnknownRelease(this *ICoreWebView2FramePermissionRequestedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2FramePermissionRequestedEventHandlerInvoke(this *ICoreWebView2FramePermissionRequestedEventHandler, sender *ICoreWebView2Frame, args *ICoreWebView2PermissionRequestedEventArgs2) uintptr {
	return this.impl.FramePermissionRequested(sender, args)
}

type ICoreWebView2FramePermissionRequestedEventHandlerImpl interface {
	IUnknownImpl
	FramePermissionRequested(sender *ICoreWebView2Frame, args *ICoreWebView2PermissionRequestedEventArgs2) uintptr
}

var ICoreWebView2FramePermissionRequestedEventHandlerFn = ICoreWebView2FramePermissionRequestedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2FramePermissionRequestedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FramePermissionRequestedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FramePermissionRequestedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FramePermissionRequestedEventHandlerInvoke),
}

func NewICoreWebView2FramePermissionRequestedEventHandler(impl ICoreWebView2FramePermissionRequestedEventHandlerImpl) *ICoreWebView2FramePermissionRequestedEventHandler {
	return &ICoreWebView2FramePermissionRequestedEventHandler{
		Vtbl: &ICoreWebView2FramePermissionRequestedEventHandlerFn,
		impl: impl,
	}
}
