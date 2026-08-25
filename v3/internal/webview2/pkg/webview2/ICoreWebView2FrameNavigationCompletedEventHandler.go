//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2FrameNavigationCompletedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FrameNavigationCompletedEventHandler struct {
	Vtbl *ICoreWebView2FrameNavigationCompletedEventHandlerVtbl
	impl ICoreWebView2FrameNavigationCompletedEventHandlerImpl
}

func (i *ICoreWebView2FrameNavigationCompletedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownQueryInterface(this *ICoreWebView2FrameNavigationCompletedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownAddRef(this *ICoreWebView2FrameNavigationCompletedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownRelease(this *ICoreWebView2FrameNavigationCompletedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2FrameNavigationCompletedEventHandlerInvoke(this *ICoreWebView2FrameNavigationCompletedEventHandler, sender *ICoreWebView2Frame, args *ICoreWebView2NavigationCompletedEventArgs) uintptr {
	return this.impl.FrameNavigationCompleted(sender, args)
}

type ICoreWebView2FrameNavigationCompletedEventHandlerImpl interface {
	IUnknownImpl
	FrameNavigationCompleted(sender *ICoreWebView2Frame, args *ICoreWebView2NavigationCompletedEventArgs) uintptr
}

var ICoreWebView2FrameNavigationCompletedEventHandlerFn = ICoreWebView2FrameNavigationCompletedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FrameNavigationCompletedEventHandlerInvoke),
}

func NewICoreWebView2FrameNavigationCompletedEventHandler(impl ICoreWebView2FrameNavigationCompletedEventHandlerImpl) *ICoreWebView2FrameNavigationCompletedEventHandler {
	return &ICoreWebView2FrameNavigationCompletedEventHandler{
		Vtbl: &ICoreWebView2FrameNavigationCompletedEventHandlerFn,
		impl: impl,
	}
}
