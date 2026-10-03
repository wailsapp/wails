//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2ProcessFailedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ProcessFailedEventHandler struct {
	Vtbl *ICoreWebView2ProcessFailedEventHandlerVtbl
	impl ICoreWebView2ProcessFailedEventHandlerImpl
}

func (i *ICoreWebView2ProcessFailedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ProcessFailedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2ProcessFailedEventHandlerIUnknownQueryInterface(this *ICoreWebView2ProcessFailedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ProcessFailedEventHandlerIUnknownAddRef(this *ICoreWebView2ProcessFailedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2ProcessFailedEventHandlerIUnknownRelease(this *ICoreWebView2ProcessFailedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2ProcessFailedEventHandlerInvoke(this *ICoreWebView2ProcessFailedEventHandler, sender *ICoreWebView2, args *ICoreWebView2ProcessFailedEventArgs) uintptr {
	return this.impl.ProcessFailed(sender, args)
}

type ICoreWebView2ProcessFailedEventHandlerImpl interface {
	IUnknownImpl
	ProcessFailed(sender *ICoreWebView2, args *ICoreWebView2ProcessFailedEventArgs) uintptr
}

var ICoreWebView2ProcessFailedEventHandlerFn = ICoreWebView2ProcessFailedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2ProcessFailedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ProcessFailedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ProcessFailedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ProcessFailedEventHandlerInvoke),
}

func NewICoreWebView2ProcessFailedEventHandler(impl ICoreWebView2ProcessFailedEventHandlerImpl) *ICoreWebView2ProcessFailedEventHandler {
	return &ICoreWebView2ProcessFailedEventHandler{
		Vtbl: &ICoreWebView2ProcessFailedEventHandlerFn,
		impl: impl,
	}
}
