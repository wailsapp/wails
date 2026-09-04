//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2IsMutedChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2IsMutedChangedEventHandler struct {
	Vtbl *ICoreWebView2IsMutedChangedEventHandlerVtbl
	impl ICoreWebView2IsMutedChangedEventHandlerImpl
}

func (i *ICoreWebView2IsMutedChangedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2IsMutedChangedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2IsMutedChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2IsMutedChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2IsMutedChangedEventHandlerIUnknownAddRef(this *ICoreWebView2IsMutedChangedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2IsMutedChangedEventHandlerIUnknownRelease(this *ICoreWebView2IsMutedChangedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2IsMutedChangedEventHandlerInvoke(this *ICoreWebView2IsMutedChangedEventHandler, sender *ICoreWebView2, args *IUnknown) uintptr {
	return this.impl.IsMutedChanged(sender, args)
}

type ICoreWebView2IsMutedChangedEventHandlerImpl interface {
	IUnknownImpl
	IsMutedChanged(sender *ICoreWebView2, args *IUnknown) uintptr
}

var ICoreWebView2IsMutedChangedEventHandlerFn = ICoreWebView2IsMutedChangedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2IsMutedChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2IsMutedChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2IsMutedChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2IsMutedChangedEventHandlerInvoke),
}

func NewICoreWebView2IsMutedChangedEventHandler(impl ICoreWebView2IsMutedChangedEventHandlerImpl) *ICoreWebView2IsMutedChangedEventHandler {
	return &ICoreWebView2IsMutedChangedEventHandler{
		Vtbl: &ICoreWebView2IsMutedChangedEventHandlerFn,
		impl: impl,
	}
}
