//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2NewBrowserVersionAvailableEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2NewBrowserVersionAvailableEventHandler struct {
	Vtbl *ICoreWebView2NewBrowserVersionAvailableEventHandlerVtbl
	impl ICoreWebView2NewBrowserVersionAvailableEventHandlerImpl
}

func (i *ICoreWebView2NewBrowserVersionAvailableEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2NewBrowserVersionAvailableEventHandlerIUnknownQueryInterface(this *ICoreWebView2NewBrowserVersionAvailableEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2NewBrowserVersionAvailableEventHandlerIUnknownAddRef(this *ICoreWebView2NewBrowserVersionAvailableEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2NewBrowserVersionAvailableEventHandlerIUnknownRelease(this *ICoreWebView2NewBrowserVersionAvailableEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2NewBrowserVersionAvailableEventHandlerInvoke(this *ICoreWebView2NewBrowserVersionAvailableEventHandler, sender *ICoreWebView2Environment, args *IUnknown) uintptr {
	return this.impl.NewBrowserVersionAvailable(sender, args)
}

type ICoreWebView2NewBrowserVersionAvailableEventHandlerImpl interface {
	IUnknownImpl
	NewBrowserVersionAvailable(sender *ICoreWebView2Environment, args *IUnknown) uintptr
}

var ICoreWebView2NewBrowserVersionAvailableEventHandlerFn = ICoreWebView2NewBrowserVersionAvailableEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2NewBrowserVersionAvailableEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2NewBrowserVersionAvailableEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2NewBrowserVersionAvailableEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2NewBrowserVersionAvailableEventHandlerInvoke),
}

func NewICoreWebView2NewBrowserVersionAvailableEventHandler(impl ICoreWebView2NewBrowserVersionAvailableEventHandlerImpl) *ICoreWebView2NewBrowserVersionAvailableEventHandler {
	return &ICoreWebView2NewBrowserVersionAvailableEventHandler{
		Vtbl: &ICoreWebView2NewBrowserVersionAvailableEventHandlerFn,
		impl: impl,
	}
}
