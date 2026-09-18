//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2ContextMenuRequestedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ContextMenuRequestedEventHandler struct {
	Vtbl *ICoreWebView2ContextMenuRequestedEventHandlerVtbl
	impl ICoreWebView2ContextMenuRequestedEventHandlerImpl
}

func (i *ICoreWebView2ContextMenuRequestedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2ContextMenuRequestedEventHandlerIUnknownQueryInterface(this *ICoreWebView2ContextMenuRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ContextMenuRequestedEventHandlerIUnknownAddRef(this *ICoreWebView2ContextMenuRequestedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2ContextMenuRequestedEventHandlerIUnknownRelease(this *ICoreWebView2ContextMenuRequestedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2ContextMenuRequestedEventHandlerInvoke(this *ICoreWebView2ContextMenuRequestedEventHandler, sender *ICoreWebView2, args *ICoreWebView2ContextMenuRequestedEventArgs) uintptr {
	return this.impl.ContextMenuRequested(sender, args)
}

type ICoreWebView2ContextMenuRequestedEventHandlerImpl interface {
	IUnknownImpl
	ContextMenuRequested(sender *ICoreWebView2, args *ICoreWebView2ContextMenuRequestedEventArgs) uintptr
}

var ICoreWebView2ContextMenuRequestedEventHandlerFn = ICoreWebView2ContextMenuRequestedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2ContextMenuRequestedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ContextMenuRequestedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ContextMenuRequestedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ContextMenuRequestedEventHandlerInvoke),
}

func NewICoreWebView2ContextMenuRequestedEventHandler(impl ICoreWebView2ContextMenuRequestedEventHandlerImpl) *ICoreWebView2ContextMenuRequestedEventHandler {
	return &ICoreWebView2ContextMenuRequestedEventHandler{
		Vtbl: &ICoreWebView2ContextMenuRequestedEventHandlerFn,
		impl: impl,
	}
}
