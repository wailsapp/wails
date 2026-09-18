//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2FrameWebMessageReceivedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FrameWebMessageReceivedEventHandler struct {
	Vtbl *ICoreWebView2FrameWebMessageReceivedEventHandlerVtbl
	impl ICoreWebView2FrameWebMessageReceivedEventHandlerImpl
}

func (i *ICoreWebView2FrameWebMessageReceivedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2FrameWebMessageReceivedEventHandlerIUnknownQueryInterface(this *ICoreWebView2FrameWebMessageReceivedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FrameWebMessageReceivedEventHandlerIUnknownAddRef(this *ICoreWebView2FrameWebMessageReceivedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2FrameWebMessageReceivedEventHandlerIUnknownRelease(this *ICoreWebView2FrameWebMessageReceivedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2FrameWebMessageReceivedEventHandlerInvoke(this *ICoreWebView2FrameWebMessageReceivedEventHandler, sender *ICoreWebView2Frame, args *ICoreWebView2WebMessageReceivedEventArgs) uintptr {
	return this.impl.FrameWebMessageReceived(sender, args)
}

type ICoreWebView2FrameWebMessageReceivedEventHandlerImpl interface {
	IUnknownImpl
	FrameWebMessageReceived(sender *ICoreWebView2Frame, args *ICoreWebView2WebMessageReceivedEventArgs) uintptr
}

var ICoreWebView2FrameWebMessageReceivedEventHandlerFn = ICoreWebView2FrameWebMessageReceivedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2FrameWebMessageReceivedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FrameWebMessageReceivedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FrameWebMessageReceivedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FrameWebMessageReceivedEventHandlerInvoke),
}

func NewICoreWebView2FrameWebMessageReceivedEventHandler(impl ICoreWebView2FrameWebMessageReceivedEventHandlerImpl) *ICoreWebView2FrameWebMessageReceivedEventHandler {
	return &ICoreWebView2FrameWebMessageReceivedEventHandler{
		Vtbl: &ICoreWebView2FrameWebMessageReceivedEventHandlerFn,
		impl: impl,
	}
}
