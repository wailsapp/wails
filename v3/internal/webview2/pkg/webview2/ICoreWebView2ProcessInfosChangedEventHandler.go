//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2ProcessInfosChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ProcessInfosChangedEventHandler struct {
	Vtbl *ICoreWebView2ProcessInfosChangedEventHandlerVtbl
	impl ICoreWebView2ProcessInfosChangedEventHandlerImpl
}

func (i *ICoreWebView2ProcessInfosChangedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ProcessInfosChangedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2ProcessInfosChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2ProcessInfosChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ProcessInfosChangedEventHandlerIUnknownAddRef(this *ICoreWebView2ProcessInfosChangedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2ProcessInfosChangedEventHandlerIUnknownRelease(this *ICoreWebView2ProcessInfosChangedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2ProcessInfosChangedEventHandlerInvoke(this *ICoreWebView2ProcessInfosChangedEventHandler, sender *ICoreWebView2Environment, args *IUnknown) uintptr {
	return this.impl.ProcessInfosChanged(sender, args)
}

type ICoreWebView2ProcessInfosChangedEventHandlerImpl interface {
	IUnknownImpl
	ProcessInfosChanged(sender *ICoreWebView2Environment, args *IUnknown) uintptr
}

var ICoreWebView2ProcessInfosChangedEventHandlerFn = ICoreWebView2ProcessInfosChangedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2ProcessInfosChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ProcessInfosChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ProcessInfosChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ProcessInfosChangedEventHandlerInvoke),
}

func NewICoreWebView2ProcessInfosChangedEventHandler(impl ICoreWebView2ProcessInfosChangedEventHandlerImpl) *ICoreWebView2ProcessInfosChangedEventHandler {
	return &ICoreWebView2ProcessInfosChangedEventHandler{
		Vtbl: &ICoreWebView2ProcessInfosChangedEventHandlerFn,
		impl: impl,
	}
}
