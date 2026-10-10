//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2ContentLoadingEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ContentLoadingEventHandler struct {
	Vtbl *ICoreWebView2ContentLoadingEventHandlerVtbl
	impl ICoreWebView2ContentLoadingEventHandlerImpl
}

func (i *ICoreWebView2ContentLoadingEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ContentLoadingEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2ContentLoadingEventHandlerIUnknownQueryInterface(this *ICoreWebView2ContentLoadingEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ContentLoadingEventHandlerIUnknownAddRef(this *ICoreWebView2ContentLoadingEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2ContentLoadingEventHandlerIUnknownRelease(this *ICoreWebView2ContentLoadingEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2ContentLoadingEventHandlerInvoke(this *ICoreWebView2ContentLoadingEventHandler, sender *ICoreWebView2, args *ICoreWebView2ContentLoadingEventArgs) uintptr {
	return this.impl.ContentLoading(sender, args)
}

type ICoreWebView2ContentLoadingEventHandlerImpl interface {
	IUnknownImpl
	ContentLoading(sender *ICoreWebView2, args *ICoreWebView2ContentLoadingEventArgs) uintptr
}

var ICoreWebView2ContentLoadingEventHandlerFn = ICoreWebView2ContentLoadingEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2ContentLoadingEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ContentLoadingEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ContentLoadingEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ContentLoadingEventHandlerInvoke),
}

func NewICoreWebView2ContentLoadingEventHandler(impl ICoreWebView2ContentLoadingEventHandlerImpl) *ICoreWebView2ContentLoadingEventHandler {
	return &ICoreWebView2ContentLoadingEventHandler{
		Vtbl: &ICoreWebView2ContentLoadingEventHandlerFn,
		impl: impl,
	}
}
