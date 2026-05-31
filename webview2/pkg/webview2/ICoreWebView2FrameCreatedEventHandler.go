//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2FrameCreatedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FrameCreatedEventHandler struct {
	Vtbl *ICoreWebView2FrameCreatedEventHandlerVtbl
	impl ICoreWebView2FrameCreatedEventHandlerImpl
}

func (i *ICoreWebView2FrameCreatedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2FrameCreatedEventHandlerIUnknownQueryInterface(this *ICoreWebView2FrameCreatedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FrameCreatedEventHandlerIUnknownAddRef(this *ICoreWebView2FrameCreatedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2FrameCreatedEventHandlerIUnknownRelease(this *ICoreWebView2FrameCreatedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2FrameCreatedEventHandlerInvoke(this *ICoreWebView2FrameCreatedEventHandler, sender *ICoreWebView2, args *ICoreWebView2FrameCreatedEventArgs) uintptr {
	return this.impl.FrameCreated(sender, args)
}

type ICoreWebView2FrameCreatedEventHandlerImpl interface {
	IUnknownImpl
	FrameCreated(sender *ICoreWebView2, args *ICoreWebView2FrameCreatedEventArgs) uintptr
}

var ICoreWebView2FrameCreatedEventHandlerFn = ICoreWebView2FrameCreatedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2FrameCreatedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FrameCreatedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FrameCreatedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FrameCreatedEventHandlerInvoke),
}

func NewICoreWebView2FrameCreatedEventHandler(impl ICoreWebView2FrameCreatedEventHandlerImpl) *ICoreWebView2FrameCreatedEventHandler {
	return &ICoreWebView2FrameCreatedEventHandler{
		Vtbl: &ICoreWebView2FrameCreatedEventHandlerFn,
		impl: impl,
	}
}
