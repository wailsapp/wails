//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2FrameChildFrameCreatedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FrameChildFrameCreatedEventHandler struct {
	Vtbl *ICoreWebView2FrameChildFrameCreatedEventHandlerVtbl
	impl ICoreWebView2FrameChildFrameCreatedEventHandlerImpl
}

func (i *ICoreWebView2FrameChildFrameCreatedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2FrameChildFrameCreatedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2FrameChildFrameCreatedEventHandlerIUnknownQueryInterface(this *ICoreWebView2FrameChildFrameCreatedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FrameChildFrameCreatedEventHandlerIUnknownAddRef(this *ICoreWebView2FrameChildFrameCreatedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2FrameChildFrameCreatedEventHandlerIUnknownRelease(this *ICoreWebView2FrameChildFrameCreatedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2FrameChildFrameCreatedEventHandlerInvoke(this *ICoreWebView2FrameChildFrameCreatedEventHandler, sender *ICoreWebView2Frame, args *ICoreWebView2FrameCreatedEventArgs) uintptr {
	return this.impl.FrameChildFrameCreated(sender, args)
}

type ICoreWebView2FrameChildFrameCreatedEventHandlerImpl interface {
	IUnknownImpl
	FrameChildFrameCreated(sender *ICoreWebView2Frame, args *ICoreWebView2FrameCreatedEventArgs) uintptr
}

var ICoreWebView2FrameChildFrameCreatedEventHandlerFn = ICoreWebView2FrameChildFrameCreatedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2FrameChildFrameCreatedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FrameChildFrameCreatedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FrameChildFrameCreatedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FrameChildFrameCreatedEventHandlerInvoke),
}

func NewICoreWebView2FrameChildFrameCreatedEventHandler(impl ICoreWebView2FrameChildFrameCreatedEventHandlerImpl) *ICoreWebView2FrameChildFrameCreatedEventHandler {
	return &ICoreWebView2FrameChildFrameCreatedEventHandler{
		Vtbl: &ICoreWebView2FrameChildFrameCreatedEventHandlerFn,
		impl: impl,
	}
}
