//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2ZoomFactorChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ZoomFactorChangedEventHandler struct {
	Vtbl *ICoreWebView2ZoomFactorChangedEventHandlerVtbl
	impl ICoreWebView2ZoomFactorChangedEventHandlerImpl
}

func (i *ICoreWebView2ZoomFactorChangedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ZoomFactorChangedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2ZoomFactorChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2ZoomFactorChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ZoomFactorChangedEventHandlerIUnknownAddRef(this *ICoreWebView2ZoomFactorChangedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2ZoomFactorChangedEventHandlerIUnknownRelease(this *ICoreWebView2ZoomFactorChangedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2ZoomFactorChangedEventHandlerInvoke(this *ICoreWebView2ZoomFactorChangedEventHandler, sender *ICoreWebView2Controller, args *IUnknown) uintptr {
	return this.impl.ZoomFactorChanged(sender, args)
}

type ICoreWebView2ZoomFactorChangedEventHandlerImpl interface {
	IUnknownImpl
	ZoomFactorChanged(sender *ICoreWebView2Controller, args *IUnknown) uintptr
}

var ICoreWebView2ZoomFactorChangedEventHandlerFn = ICoreWebView2ZoomFactorChangedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2ZoomFactorChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ZoomFactorChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ZoomFactorChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ZoomFactorChangedEventHandlerInvoke),
}

func NewICoreWebView2ZoomFactorChangedEventHandler(impl ICoreWebView2ZoomFactorChangedEventHandlerImpl) *ICoreWebView2ZoomFactorChangedEventHandler {
	return &ICoreWebView2ZoomFactorChangedEventHandler{
		Vtbl: &ICoreWebView2ZoomFactorChangedEventHandlerFn,
		impl: impl,
	}
}
