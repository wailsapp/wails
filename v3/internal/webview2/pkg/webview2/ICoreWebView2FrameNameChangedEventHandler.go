//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2FrameNameChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FrameNameChangedEventHandler struct {
	Vtbl *ICoreWebView2FrameNameChangedEventHandlerVtbl
	impl ICoreWebView2FrameNameChangedEventHandlerImpl
}

func (i *ICoreWebView2FrameNameChangedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2FrameNameChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2FrameNameChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FrameNameChangedEventHandlerIUnknownAddRef(this *ICoreWebView2FrameNameChangedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2FrameNameChangedEventHandlerIUnknownRelease(this *ICoreWebView2FrameNameChangedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2FrameNameChangedEventHandlerInvoke(this *ICoreWebView2FrameNameChangedEventHandler, sender *ICoreWebView2Frame, args *IUnknown) uintptr {
	return this.impl.FrameNameChanged(sender, args)
}

type ICoreWebView2FrameNameChangedEventHandlerImpl interface {
	IUnknownImpl
	FrameNameChanged(sender *ICoreWebView2Frame, args *IUnknown) uintptr
}

var ICoreWebView2FrameNameChangedEventHandlerFn = ICoreWebView2FrameNameChangedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2FrameNameChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FrameNameChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FrameNameChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FrameNameChangedEventHandlerInvoke),
}

func NewICoreWebView2FrameNameChangedEventHandler(impl ICoreWebView2FrameNameChangedEventHandlerImpl) *ICoreWebView2FrameNameChangedEventHandler {
	return &ICoreWebView2FrameNameChangedEventHandler{
		Vtbl: &ICoreWebView2FrameNameChangedEventHandlerFn,
		impl: impl,
	}
}
