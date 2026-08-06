//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2FrameNavigationCompletedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FrameNavigationCompletedEventHandler struct {
	Vtbl *ICoreWebView2FrameNavigationCompletedEventHandlerVtbl
	impl ICoreWebView2FrameNavigationCompletedEventHandlerImpl
}

func (i *ICoreWebView2FrameNavigationCompletedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2FrameNavigationCompletedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownQueryInterface(this *ICoreWebView2FrameNavigationCompletedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownAddRef(this *ICoreWebView2FrameNavigationCompletedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownRelease(this *ICoreWebView2FrameNavigationCompletedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2FrameNavigationCompletedEventHandlerInvoke(this *ICoreWebView2FrameNavigationCompletedEventHandler, sender *ICoreWebView2Frame, args *ICoreWebView2NavigationCompletedEventArgs) uintptr {
	return this.impl.FrameNavigationCompleted(sender, args)
}

type ICoreWebView2FrameNavigationCompletedEventHandlerImpl interface {
	IUnknownImpl
	FrameNavigationCompleted(sender *ICoreWebView2Frame, args *ICoreWebView2NavigationCompletedEventArgs) uintptr
}

var ICoreWebView2FrameNavigationCompletedEventHandlerFn = ICoreWebView2FrameNavigationCompletedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FrameNavigationCompletedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FrameNavigationCompletedEventHandlerInvoke),
}

func NewICoreWebView2FrameNavigationCompletedEventHandler(impl ICoreWebView2FrameNavigationCompletedEventHandlerImpl) *ICoreWebView2FrameNavigationCompletedEventHandler {
	return &ICoreWebView2FrameNavigationCompletedEventHandler{
		Vtbl: &ICoreWebView2FrameNavigationCompletedEventHandlerFn,
		impl: impl,
	}
}
