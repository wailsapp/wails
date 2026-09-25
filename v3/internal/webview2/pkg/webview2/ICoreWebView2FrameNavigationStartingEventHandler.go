//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2FrameNavigationStartingEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FrameNavigationStartingEventHandler struct {
	Vtbl *ICoreWebView2FrameNavigationStartingEventHandlerVtbl
	impl ICoreWebView2FrameNavigationStartingEventHandlerImpl
}

func (i *ICoreWebView2FrameNavigationStartingEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2FrameNavigationStartingEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2FrameNavigationStartingEventHandlerIUnknownQueryInterface(this *ICoreWebView2FrameNavigationStartingEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FrameNavigationStartingEventHandlerIUnknownAddRef(this *ICoreWebView2FrameNavigationStartingEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2FrameNavigationStartingEventHandlerIUnknownRelease(this *ICoreWebView2FrameNavigationStartingEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2FrameNavigationStartingEventHandlerInvoke(this *ICoreWebView2FrameNavigationStartingEventHandler, sender *ICoreWebView2Frame, args *ICoreWebView2NavigationStartingEventArgs) uintptr {
	return this.impl.FrameNavigationStarting(sender, args)
}

type ICoreWebView2FrameNavigationStartingEventHandlerImpl interface {
	IUnknownImpl
	FrameNavigationStarting(sender *ICoreWebView2Frame, args *ICoreWebView2NavigationStartingEventArgs) uintptr
}

var ICoreWebView2FrameNavigationStartingEventHandlerFn = ICoreWebView2FrameNavigationStartingEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2FrameNavigationStartingEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FrameNavigationStartingEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FrameNavigationStartingEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FrameNavigationStartingEventHandlerInvoke),
}

func NewICoreWebView2FrameNavigationStartingEventHandler(impl ICoreWebView2FrameNavigationStartingEventHandlerImpl) *ICoreWebView2FrameNavigationStartingEventHandler {
	return &ICoreWebView2FrameNavigationStartingEventHandler{
		Vtbl: &ICoreWebView2FrameNavigationStartingEventHandlerFn,
		impl: impl,
	}
}
