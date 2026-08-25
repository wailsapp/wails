//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2FrameScreenCaptureStartingEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FrameScreenCaptureStartingEventHandler struct {
	Vtbl *ICoreWebView2FrameScreenCaptureStartingEventHandlerVtbl
	impl ICoreWebView2FrameScreenCaptureStartingEventHandlerImpl
}

func (i *ICoreWebView2FrameScreenCaptureStartingEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2FrameScreenCaptureStartingEventHandlerIUnknownQueryInterface(this *ICoreWebView2FrameScreenCaptureStartingEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FrameScreenCaptureStartingEventHandlerIUnknownAddRef(this *ICoreWebView2FrameScreenCaptureStartingEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2FrameScreenCaptureStartingEventHandlerIUnknownRelease(this *ICoreWebView2FrameScreenCaptureStartingEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2FrameScreenCaptureStartingEventHandlerInvoke(this *ICoreWebView2FrameScreenCaptureStartingEventHandler, sender *ICoreWebView2Frame, args *ICoreWebView2ScreenCaptureStartingEventArgs) uintptr {
	return this.impl.FrameScreenCaptureStarting(sender, args)
}

type ICoreWebView2FrameScreenCaptureStartingEventHandlerImpl interface {
	IUnknownImpl
	FrameScreenCaptureStarting(sender *ICoreWebView2Frame, args *ICoreWebView2ScreenCaptureStartingEventArgs) uintptr
}

var ICoreWebView2FrameScreenCaptureStartingEventHandlerFn = ICoreWebView2FrameScreenCaptureStartingEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2FrameScreenCaptureStartingEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FrameScreenCaptureStartingEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FrameScreenCaptureStartingEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FrameScreenCaptureStartingEventHandlerInvoke),
}

func NewICoreWebView2FrameScreenCaptureStartingEventHandler(impl ICoreWebView2FrameScreenCaptureStartingEventHandlerImpl) *ICoreWebView2FrameScreenCaptureStartingEventHandler {
	return &ICoreWebView2FrameScreenCaptureStartingEventHandler{
		Vtbl: &ICoreWebView2FrameScreenCaptureStartingEventHandlerFn,
		impl: impl,
	}
}
