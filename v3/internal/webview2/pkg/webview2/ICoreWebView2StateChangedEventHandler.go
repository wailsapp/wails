//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2StateChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2StateChangedEventHandler struct {
	Vtbl *ICoreWebView2StateChangedEventHandlerVtbl
	impl ICoreWebView2StateChangedEventHandlerImpl
}

func (i *ICoreWebView2StateChangedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2StateChangedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2StateChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2StateChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2StateChangedEventHandlerIUnknownAddRef(this *ICoreWebView2StateChangedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2StateChangedEventHandlerIUnknownRelease(this *ICoreWebView2StateChangedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2StateChangedEventHandlerInvoke(this *ICoreWebView2StateChangedEventHandler, sender *ICoreWebView2DownloadOperation, args *IUnknown) uintptr {
	return this.impl.StateChanged(sender, args)
}

type ICoreWebView2StateChangedEventHandlerImpl interface {
	IUnknownImpl
	StateChanged(sender *ICoreWebView2DownloadOperation, args *IUnknown) uintptr
}

var ICoreWebView2StateChangedEventHandlerFn = ICoreWebView2StateChangedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2StateChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2StateChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2StateChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2StateChangedEventHandlerInvoke),
}

func NewICoreWebView2StateChangedEventHandler(impl ICoreWebView2StateChangedEventHandlerImpl) *ICoreWebView2StateChangedEventHandler {
	return &ICoreWebView2StateChangedEventHandler{
		Vtbl: &ICoreWebView2StateChangedEventHandlerFn,
		impl: impl,
	}
}
