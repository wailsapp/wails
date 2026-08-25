//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2NotificationReceivedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2NotificationReceivedEventHandler struct {
	Vtbl *ICoreWebView2NotificationReceivedEventHandlerVtbl
	impl ICoreWebView2NotificationReceivedEventHandlerImpl
}

func (i *ICoreWebView2NotificationReceivedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2NotificationReceivedEventHandlerIUnknownQueryInterface(this *ICoreWebView2NotificationReceivedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2NotificationReceivedEventHandlerIUnknownAddRef(this *ICoreWebView2NotificationReceivedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2NotificationReceivedEventHandlerIUnknownRelease(this *ICoreWebView2NotificationReceivedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2NotificationReceivedEventHandlerInvoke(this *ICoreWebView2NotificationReceivedEventHandler, sender *ICoreWebView2, args *ICoreWebView2NotificationReceivedEventArgs) uintptr {
	return this.impl.NotificationReceived(sender, args)
}

type ICoreWebView2NotificationReceivedEventHandlerImpl interface {
	IUnknownImpl
	NotificationReceived(sender *ICoreWebView2, args *ICoreWebView2NotificationReceivedEventArgs) uintptr
}

var ICoreWebView2NotificationReceivedEventHandlerFn = ICoreWebView2NotificationReceivedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2NotificationReceivedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2NotificationReceivedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2NotificationReceivedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2NotificationReceivedEventHandlerInvoke),
}

func NewICoreWebView2NotificationReceivedEventHandler(impl ICoreWebView2NotificationReceivedEventHandlerImpl) *ICoreWebView2NotificationReceivedEventHandler {
	return &ICoreWebView2NotificationReceivedEventHandler{
		Vtbl: &ICoreWebView2NotificationReceivedEventHandlerFn,
		impl: impl,
	}
}
