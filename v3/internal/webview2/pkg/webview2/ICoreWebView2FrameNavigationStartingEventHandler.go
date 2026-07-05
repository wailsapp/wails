//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2FrameNavigationStartingEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FrameNavigationStartingEventHandler struct {
	Vtbl *ICoreWebView2FrameNavigationStartingEventHandlerVtbl
	impl ICoreWebView2FrameNavigationStartingEventHandlerImpl
}

func (i *ICoreWebView2FrameNavigationStartingEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2FrameNavigationStartingEventHandlerIUnknownQueryInterface(this *ICoreWebView2FrameNavigationStartingEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FrameNavigationStartingEventHandlerIUnknownAddRef(this *ICoreWebView2FrameNavigationStartingEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2FrameNavigationStartingEventHandlerIUnknownRelease(this *ICoreWebView2FrameNavigationStartingEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2FrameNavigationStartingEventHandlerInvoke(this *ICoreWebView2FrameNavigationStartingEventHandler, sender *ICoreWebView2Frame, args *ICoreWebView2NavigationStartingEventArgs) uintptr {
	return this.impl.FrameNavigationStarting(sender, args)
}

type ICoreWebView2FrameNavigationStartingEventHandlerImpl interface {
	IUnknownImpl
	FrameNavigationStarting(sender *ICoreWebView2Frame, args *ICoreWebView2NavigationStartingEventArgs) uintptr
}

var ICoreWebView2FrameNavigationStartingEventHandlerFn = ICoreWebView2FrameNavigationStartingEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2FrameNavigationStartingEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FrameNavigationStartingEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FrameNavigationStartingEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FrameNavigationStartingEventHandlerInvoke),
}

func NewICoreWebView2FrameNavigationStartingEventHandler(impl ICoreWebView2FrameNavigationStartingEventHandlerImpl) *ICoreWebView2FrameNavigationStartingEventHandler {
	return &ICoreWebView2FrameNavigationStartingEventHandler{
		Vtbl: &ICoreWebView2FrameNavigationStartingEventHandlerFn,
		impl: impl,
	}
}
