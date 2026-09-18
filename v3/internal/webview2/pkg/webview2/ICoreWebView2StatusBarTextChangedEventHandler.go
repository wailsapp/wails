//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2StatusBarTextChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2StatusBarTextChangedEventHandler struct {
	Vtbl *ICoreWebView2StatusBarTextChangedEventHandlerVtbl
	impl ICoreWebView2StatusBarTextChangedEventHandlerImpl
}

func (i *ICoreWebView2StatusBarTextChangedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2StatusBarTextChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2StatusBarTextChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2StatusBarTextChangedEventHandlerIUnknownAddRef(this *ICoreWebView2StatusBarTextChangedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2StatusBarTextChangedEventHandlerIUnknownRelease(this *ICoreWebView2StatusBarTextChangedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2StatusBarTextChangedEventHandlerInvoke(this *ICoreWebView2StatusBarTextChangedEventHandler, sender *ICoreWebView2, args *IUnknown) uintptr {
	return this.impl.StatusBarTextChanged(sender, args)
}

type ICoreWebView2StatusBarTextChangedEventHandlerImpl interface {
	IUnknownImpl
	StatusBarTextChanged(sender *ICoreWebView2, args *IUnknown) uintptr
}

var ICoreWebView2StatusBarTextChangedEventHandlerFn = ICoreWebView2StatusBarTextChangedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2StatusBarTextChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2StatusBarTextChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2StatusBarTextChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2StatusBarTextChangedEventHandlerInvoke),
}

func NewICoreWebView2StatusBarTextChangedEventHandler(impl ICoreWebView2StatusBarTextChangedEventHandlerImpl) *ICoreWebView2StatusBarTextChangedEventHandler {
	return &ICoreWebView2StatusBarTextChangedEventHandler{
		Vtbl: &ICoreWebView2StatusBarTextChangedEventHandlerFn,
		impl: impl,
	}
}
