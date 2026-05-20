//go:build windows

package edge

import (
	"unsafe"
)

type ICoreWebView2RasterizationScaleChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2RasterizationScaleChangedEventHandler struct {
	Vtbl *ICoreWebView2RasterizationScaleChangedEventHandlerVtbl
	impl ICoreWebView2RasterizationScaleChangedEventHandlerImpl
}

func (i *ICoreWebView2RasterizationScaleChangedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2RasterizationScaleChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2RasterizationScaleChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2RasterizationScaleChangedEventHandlerIUnknownAddRef(this *ICoreWebView2RasterizationScaleChangedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2RasterizationScaleChangedEventHandlerIUnknownRelease(this *ICoreWebView2RasterizationScaleChangedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2RasterizationScaleChangedEventHandlerInvoke(this *ICoreWebView2RasterizationScaleChangedEventHandler, sender *ICoreWebView2Controller, args *IUnknown) uintptr {
	return this.impl.RasterizationScaleChanged(sender, args)
}

type ICoreWebView2RasterizationScaleChangedEventHandlerImpl interface {
	IUnknownImpl
	RasterizationScaleChanged(sender *ICoreWebView2Controller, args *IUnknown) uintptr
}

var ICoreWebView2RasterizationScaleChangedEventHandlerFn = ICoreWebView2RasterizationScaleChangedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2RasterizationScaleChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2RasterizationScaleChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2RasterizationScaleChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2RasterizationScaleChangedEventHandlerInvoke),
}

func NewICoreWebView2RasterizationScaleChangedEventHandler(impl ICoreWebView2RasterizationScaleChangedEventHandlerImpl) *ICoreWebView2RasterizationScaleChangedEventHandler {
	return &ICoreWebView2RasterizationScaleChangedEventHandler{
		Vtbl: &ICoreWebView2RasterizationScaleChangedEventHandlerFn,
		impl: impl,
	}
}
