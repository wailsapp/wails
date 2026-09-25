//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2WindowCloseRequestedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2WindowCloseRequestedEventHandler struct {
	Vtbl *ICoreWebView2WindowCloseRequestedEventHandlerVtbl
	impl ICoreWebView2WindowCloseRequestedEventHandlerImpl
}

func (i *ICoreWebView2WindowCloseRequestedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2WindowCloseRequestedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2WindowCloseRequestedEventHandlerIUnknownQueryInterface(this *ICoreWebView2WindowCloseRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2WindowCloseRequestedEventHandlerIUnknownAddRef(this *ICoreWebView2WindowCloseRequestedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2WindowCloseRequestedEventHandlerIUnknownRelease(this *ICoreWebView2WindowCloseRequestedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2WindowCloseRequestedEventHandlerInvoke(this *ICoreWebView2WindowCloseRequestedEventHandler, sender *ICoreWebView2, args *IUnknown) uintptr {
	return this.impl.WindowCloseRequested(sender, args)
}

type ICoreWebView2WindowCloseRequestedEventHandlerImpl interface {
	IUnknownImpl
	WindowCloseRequested(sender *ICoreWebView2, args *IUnknown) uintptr
}

var ICoreWebView2WindowCloseRequestedEventHandlerFn = ICoreWebView2WindowCloseRequestedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2WindowCloseRequestedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2WindowCloseRequestedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2WindowCloseRequestedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2WindowCloseRequestedEventHandlerInvoke),
}

func NewICoreWebView2WindowCloseRequestedEventHandler(impl ICoreWebView2WindowCloseRequestedEventHandlerImpl) *ICoreWebView2WindowCloseRequestedEventHandler {
	return &ICoreWebView2WindowCloseRequestedEventHandler{
		Vtbl: &ICoreWebView2WindowCloseRequestedEventHandlerFn,
		impl: impl,
	}
}
