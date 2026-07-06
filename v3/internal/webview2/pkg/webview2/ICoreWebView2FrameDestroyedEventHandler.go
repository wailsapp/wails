//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2FrameDestroyedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FrameDestroyedEventHandler struct {
	Vtbl *ICoreWebView2FrameDestroyedEventHandlerVtbl
	impl ICoreWebView2FrameDestroyedEventHandlerImpl
}

func (i *ICoreWebView2FrameDestroyedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2FrameDestroyedEventHandlerIUnknownQueryInterface(this *ICoreWebView2FrameDestroyedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FrameDestroyedEventHandlerIUnknownAddRef(this *ICoreWebView2FrameDestroyedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2FrameDestroyedEventHandlerIUnknownRelease(this *ICoreWebView2FrameDestroyedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2FrameDestroyedEventHandlerInvoke(this *ICoreWebView2FrameDestroyedEventHandler, sender *ICoreWebView2Frame, args *IUnknown) uintptr {
	return this.impl.FrameDestroyed(sender, args)
}

type ICoreWebView2FrameDestroyedEventHandlerImpl interface {
	IUnknownImpl
	FrameDestroyed(sender *ICoreWebView2Frame, args *IUnknown) uintptr
}

var ICoreWebView2FrameDestroyedEventHandlerFn = ICoreWebView2FrameDestroyedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2FrameDestroyedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FrameDestroyedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FrameDestroyedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FrameDestroyedEventHandlerInvoke),
}

func NewICoreWebView2FrameDestroyedEventHandler(impl ICoreWebView2FrameDestroyedEventHandlerImpl) *ICoreWebView2FrameDestroyedEventHandler {
	return &ICoreWebView2FrameDestroyedEventHandler{
		Vtbl: &ICoreWebView2FrameDestroyedEventHandlerFn,
		impl: impl,
	}
}
