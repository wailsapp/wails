//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2NewWindowRequestedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2NewWindowRequestedEventHandler struct {
	Vtbl *ICoreWebView2NewWindowRequestedEventHandlerVtbl
	impl ICoreWebView2NewWindowRequestedEventHandlerImpl
}

func (i *ICoreWebView2NewWindowRequestedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2NewWindowRequestedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2NewWindowRequestedEventHandlerIUnknownQueryInterface(this *ICoreWebView2NewWindowRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2NewWindowRequestedEventHandlerIUnknownAddRef(this *ICoreWebView2NewWindowRequestedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2NewWindowRequestedEventHandlerIUnknownRelease(this *ICoreWebView2NewWindowRequestedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2NewWindowRequestedEventHandlerInvoke(this *ICoreWebView2NewWindowRequestedEventHandler, sender *ICoreWebView2, args *ICoreWebView2NewWindowRequestedEventArgs) uintptr {
	return this.impl.NewWindowRequested(sender, args)
}

type ICoreWebView2NewWindowRequestedEventHandlerImpl interface {
	IUnknownImpl
	NewWindowRequested(sender *ICoreWebView2, args *ICoreWebView2NewWindowRequestedEventArgs) uintptr
}

var ICoreWebView2NewWindowRequestedEventHandlerFn = ICoreWebView2NewWindowRequestedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2NewWindowRequestedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2NewWindowRequestedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2NewWindowRequestedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2NewWindowRequestedEventHandlerInvoke),
}

func NewICoreWebView2NewWindowRequestedEventHandler(impl ICoreWebView2NewWindowRequestedEventHandlerImpl) *ICoreWebView2NewWindowRequestedEventHandler {
	return &ICoreWebView2NewWindowRequestedEventHandler{
		Vtbl: &ICoreWebView2NewWindowRequestedEventHandlerFn,
		impl: impl,
	}
}
