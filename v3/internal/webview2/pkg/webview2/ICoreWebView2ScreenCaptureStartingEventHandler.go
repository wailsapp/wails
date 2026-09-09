//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2ScreenCaptureStartingEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ScreenCaptureStartingEventHandler struct {
	Vtbl *ICoreWebView2ScreenCaptureStartingEventHandlerVtbl
	impl ICoreWebView2ScreenCaptureStartingEventHandlerImpl
}

func (i *ICoreWebView2ScreenCaptureStartingEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2ScreenCaptureStartingEventHandlerIUnknownQueryInterface(this *ICoreWebView2ScreenCaptureStartingEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ScreenCaptureStartingEventHandlerIUnknownAddRef(this *ICoreWebView2ScreenCaptureStartingEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2ScreenCaptureStartingEventHandlerIUnknownRelease(this *ICoreWebView2ScreenCaptureStartingEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2ScreenCaptureStartingEventHandlerInvoke(this *ICoreWebView2ScreenCaptureStartingEventHandler, sender *ICoreWebView2, args *ICoreWebView2ScreenCaptureStartingEventArgs) uintptr {
	return this.impl.ScreenCaptureStarting(sender, args)
}

type ICoreWebView2ScreenCaptureStartingEventHandlerImpl interface {
	IUnknownImpl
	ScreenCaptureStarting(sender *ICoreWebView2, args *ICoreWebView2ScreenCaptureStartingEventArgs) uintptr
}

var ICoreWebView2ScreenCaptureStartingEventHandlerFn = ICoreWebView2ScreenCaptureStartingEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2ScreenCaptureStartingEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ScreenCaptureStartingEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ScreenCaptureStartingEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ScreenCaptureStartingEventHandlerInvoke),
}

func NewICoreWebView2ScreenCaptureStartingEventHandler(impl ICoreWebView2ScreenCaptureStartingEventHandlerImpl) *ICoreWebView2ScreenCaptureStartingEventHandler {
	return &ICoreWebView2ScreenCaptureStartingEventHandler{
		Vtbl: &ICoreWebView2ScreenCaptureStartingEventHandlerFn,
		impl: impl,
	}
}
