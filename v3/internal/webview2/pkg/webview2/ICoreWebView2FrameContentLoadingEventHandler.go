//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2FrameContentLoadingEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FrameContentLoadingEventHandler struct {
	Vtbl *ICoreWebView2FrameContentLoadingEventHandlerVtbl
	impl ICoreWebView2FrameContentLoadingEventHandlerImpl
}

func (i *ICoreWebView2FrameContentLoadingEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2FrameContentLoadingEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2FrameContentLoadingEventHandlerIUnknownQueryInterface(this *ICoreWebView2FrameContentLoadingEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FrameContentLoadingEventHandlerIUnknownAddRef(this *ICoreWebView2FrameContentLoadingEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2FrameContentLoadingEventHandlerIUnknownRelease(this *ICoreWebView2FrameContentLoadingEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2FrameContentLoadingEventHandlerInvoke(this *ICoreWebView2FrameContentLoadingEventHandler, sender *ICoreWebView2Frame, args *ICoreWebView2ContentLoadingEventArgs) uintptr {
	return this.impl.FrameContentLoading(sender, args)
}

type ICoreWebView2FrameContentLoadingEventHandlerImpl interface {
	IUnknownImpl
	FrameContentLoading(sender *ICoreWebView2Frame, args *ICoreWebView2ContentLoadingEventArgs) uintptr
}

var ICoreWebView2FrameContentLoadingEventHandlerFn = ICoreWebView2FrameContentLoadingEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2FrameContentLoadingEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FrameContentLoadingEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FrameContentLoadingEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FrameContentLoadingEventHandlerInvoke),
}

func NewICoreWebView2FrameContentLoadingEventHandler(impl ICoreWebView2FrameContentLoadingEventHandlerImpl) *ICoreWebView2FrameContentLoadingEventHandler {
	return &ICoreWebView2FrameContentLoadingEventHandler{
		Vtbl: &ICoreWebView2FrameContentLoadingEventHandlerFn,
		impl: impl,
	}
}
