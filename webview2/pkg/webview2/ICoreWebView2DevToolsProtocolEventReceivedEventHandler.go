//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2DevToolsProtocolEventReceivedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2DevToolsProtocolEventReceivedEventHandler struct {
	Vtbl *ICoreWebView2DevToolsProtocolEventReceivedEventHandlerVtbl
	impl ICoreWebView2DevToolsProtocolEventReceivedEventHandlerImpl
}

func (i *ICoreWebView2DevToolsProtocolEventReceivedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2DevToolsProtocolEventReceivedEventHandlerIUnknownQueryInterface(this *ICoreWebView2DevToolsProtocolEventReceivedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2DevToolsProtocolEventReceivedEventHandlerIUnknownAddRef(this *ICoreWebView2DevToolsProtocolEventReceivedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2DevToolsProtocolEventReceivedEventHandlerIUnknownRelease(this *ICoreWebView2DevToolsProtocolEventReceivedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2DevToolsProtocolEventReceivedEventHandlerInvoke(this *ICoreWebView2DevToolsProtocolEventReceivedEventHandler, sender *ICoreWebView2, args *ICoreWebView2DevToolsProtocolEventReceivedEventArgs) uintptr {
	return this.impl.DevToolsProtocolEventReceived(sender, args)
}

type ICoreWebView2DevToolsProtocolEventReceivedEventHandlerImpl interface {
	IUnknownImpl
	DevToolsProtocolEventReceived(sender *ICoreWebView2, args *ICoreWebView2DevToolsProtocolEventReceivedEventArgs) uintptr
}

var ICoreWebView2DevToolsProtocolEventReceivedEventHandlerFn = ICoreWebView2DevToolsProtocolEventReceivedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2DevToolsProtocolEventReceivedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2DevToolsProtocolEventReceivedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2DevToolsProtocolEventReceivedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2DevToolsProtocolEventReceivedEventHandlerInvoke),
}

func NewICoreWebView2DevToolsProtocolEventReceivedEventHandler(impl ICoreWebView2DevToolsProtocolEventReceivedEventHandlerImpl) *ICoreWebView2DevToolsProtocolEventReceivedEventHandler {
	return &ICoreWebView2DevToolsProtocolEventReceivedEventHandler{
		Vtbl: &ICoreWebView2DevToolsProtocolEventReceivedEventHandlerFn,
		impl: impl,
	}
}
