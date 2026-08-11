//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2NonClientRegionChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2NonClientRegionChangedEventHandler struct {
	Vtbl *ICoreWebView2NonClientRegionChangedEventHandlerVtbl
	impl ICoreWebView2NonClientRegionChangedEventHandlerImpl
}

func (i *ICoreWebView2NonClientRegionChangedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2NonClientRegionChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2NonClientRegionChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2NonClientRegionChangedEventHandlerIUnknownAddRef(this *ICoreWebView2NonClientRegionChangedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2NonClientRegionChangedEventHandlerIUnknownRelease(this *ICoreWebView2NonClientRegionChangedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2NonClientRegionChangedEventHandlerInvoke(this *ICoreWebView2NonClientRegionChangedEventHandler, sender *ICoreWebView2CompositionController, args *ICoreWebView2NonClientRegionChangedEventArgs) uintptr {
	return this.impl.NonClientRegionChanged(sender, args)
}

type ICoreWebView2NonClientRegionChangedEventHandlerImpl interface {
	IUnknownImpl
	NonClientRegionChanged(sender *ICoreWebView2CompositionController, args *ICoreWebView2NonClientRegionChangedEventArgs) uintptr
}

var ICoreWebView2NonClientRegionChangedEventHandlerFn = ICoreWebView2NonClientRegionChangedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2NonClientRegionChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2NonClientRegionChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2NonClientRegionChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2NonClientRegionChangedEventHandlerInvoke),
}

func NewICoreWebView2NonClientRegionChangedEventHandler(impl ICoreWebView2NonClientRegionChangedEventHandlerImpl) *ICoreWebView2NonClientRegionChangedEventHandler {
	return &ICoreWebView2NonClientRegionChangedEventHandler{
		Vtbl: &ICoreWebView2NonClientRegionChangedEventHandlerFn,
		impl: impl,
	}
}
