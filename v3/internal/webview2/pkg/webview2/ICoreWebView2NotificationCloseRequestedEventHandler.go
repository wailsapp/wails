//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2NotificationCloseRequestedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2NotificationCloseRequestedEventHandler struct {
	Vtbl *ICoreWebView2NotificationCloseRequestedEventHandlerVtbl
	impl ICoreWebView2NotificationCloseRequestedEventHandlerImpl
}

func (i *ICoreWebView2NotificationCloseRequestedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2NotificationCloseRequestedEventHandlerIUnknownQueryInterface(this *ICoreWebView2NotificationCloseRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2NotificationCloseRequestedEventHandlerIUnknownAddRef(this *ICoreWebView2NotificationCloseRequestedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2NotificationCloseRequestedEventHandlerIUnknownRelease(this *ICoreWebView2NotificationCloseRequestedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2NotificationCloseRequestedEventHandlerInvoke(this *ICoreWebView2NotificationCloseRequestedEventHandler, sender *ICoreWebView2Notification, args *IUnknown) uintptr {
	return this.impl.NotificationCloseRequested(sender, args)
}

type ICoreWebView2NotificationCloseRequestedEventHandlerImpl interface {
	IUnknownImpl
	NotificationCloseRequested(sender *ICoreWebView2Notification, args *IUnknown) uintptr
}

var ICoreWebView2NotificationCloseRequestedEventHandlerFn = ICoreWebView2NotificationCloseRequestedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2NotificationCloseRequestedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2NotificationCloseRequestedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2NotificationCloseRequestedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2NotificationCloseRequestedEventHandlerInvoke),
}

func NewICoreWebView2NotificationCloseRequestedEventHandler(impl ICoreWebView2NotificationCloseRequestedEventHandlerImpl) *ICoreWebView2NotificationCloseRequestedEventHandler {
	return &ICoreWebView2NotificationCloseRequestedEventHandler{
		Vtbl: &ICoreWebView2NotificationCloseRequestedEventHandlerFn,
		impl: impl,
	}
}
