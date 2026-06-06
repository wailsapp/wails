//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2CursorChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2CursorChangedEventHandler struct {
	Vtbl *ICoreWebView2CursorChangedEventHandlerVtbl
	impl ICoreWebView2CursorChangedEventHandlerImpl
}

func (i *ICoreWebView2CursorChangedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2CursorChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2CursorChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2CursorChangedEventHandlerIUnknownAddRef(this *ICoreWebView2CursorChangedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2CursorChangedEventHandlerIUnknownRelease(this *ICoreWebView2CursorChangedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2CursorChangedEventHandlerInvoke(this *ICoreWebView2CursorChangedEventHandler, sender *ICoreWebView2CompositionController, args *IUnknown) uintptr {
	return this.impl.CursorChanged(sender, args)
}

type ICoreWebView2CursorChangedEventHandlerImpl interface {
	IUnknownImpl
	CursorChanged(sender *ICoreWebView2CompositionController, args *IUnknown) uintptr
}

var ICoreWebView2CursorChangedEventHandlerFn = ICoreWebView2CursorChangedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2CursorChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2CursorChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2CursorChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2CursorChangedEventHandlerInvoke),
}

func NewICoreWebView2CursorChangedEventHandler(impl ICoreWebView2CursorChangedEventHandlerImpl) *ICoreWebView2CursorChangedEventHandler {
	return &ICoreWebView2CursorChangedEventHandler{
		Vtbl: &ICoreWebView2CursorChangedEventHandlerFn,
		impl: impl,
	}
}
