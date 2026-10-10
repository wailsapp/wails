//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2NavigationCompletedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2NavigationCompletedEventHandler struct {
	Vtbl *ICoreWebView2NavigationCompletedEventHandlerVtbl
	impl ICoreWebView2NavigationCompletedEventHandlerImpl
}

func (i *ICoreWebView2NavigationCompletedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2NavigationCompletedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2NavigationCompletedEventHandlerIUnknownQueryInterface(this *ICoreWebView2NavigationCompletedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2NavigationCompletedEventHandlerIUnknownAddRef(this *ICoreWebView2NavigationCompletedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2NavigationCompletedEventHandlerIUnknownRelease(this *ICoreWebView2NavigationCompletedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2NavigationCompletedEventHandlerInvoke(this *ICoreWebView2NavigationCompletedEventHandler, sender *ICoreWebView2, args *ICoreWebView2NavigationCompletedEventArgs) uintptr {
	return this.impl.NavigationCompleted(sender, args)
}

type ICoreWebView2NavigationCompletedEventHandlerImpl interface {
	IUnknownImpl
	NavigationCompleted(sender *ICoreWebView2, args *ICoreWebView2NavigationCompletedEventArgs) uintptr
}

var ICoreWebView2NavigationCompletedEventHandlerFn = ICoreWebView2NavigationCompletedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2NavigationCompletedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2NavigationCompletedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2NavigationCompletedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2NavigationCompletedEventHandlerInvoke),
}

func NewICoreWebView2NavigationCompletedEventHandler(impl ICoreWebView2NavigationCompletedEventHandlerImpl) *ICoreWebView2NavigationCompletedEventHandler {
	return &ICoreWebView2NavigationCompletedEventHandler{
		Vtbl: &ICoreWebView2NavigationCompletedEventHandlerFn,
		impl: impl,
	}
}
