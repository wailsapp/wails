//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2FrameCreatedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FrameCreatedEventHandler struct {
	Vtbl *ICoreWebView2FrameCreatedEventHandlerVtbl
	impl ICoreWebView2FrameCreatedEventHandlerImpl
}

func (i *ICoreWebView2FrameCreatedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2FrameCreatedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2FrameCreatedEventHandlerIUnknownQueryInterface(this *ICoreWebView2FrameCreatedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FrameCreatedEventHandlerIUnknownAddRef(this *ICoreWebView2FrameCreatedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2FrameCreatedEventHandlerIUnknownRelease(this *ICoreWebView2FrameCreatedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2FrameCreatedEventHandlerInvoke(this *ICoreWebView2FrameCreatedEventHandler, sender *ICoreWebView2, args *ICoreWebView2FrameCreatedEventArgs) uintptr {
	return this.impl.FrameCreated(sender, args)
}

type ICoreWebView2FrameCreatedEventHandlerImpl interface {
	IUnknownImpl
	FrameCreated(sender *ICoreWebView2, args *ICoreWebView2FrameCreatedEventArgs) uintptr
}

var ICoreWebView2FrameCreatedEventHandlerFn = ICoreWebView2FrameCreatedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2FrameCreatedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FrameCreatedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FrameCreatedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FrameCreatedEventHandlerInvoke),
}

func NewICoreWebView2FrameCreatedEventHandler(impl ICoreWebView2FrameCreatedEventHandlerImpl) *ICoreWebView2FrameCreatedEventHandler {
	return &ICoreWebView2FrameCreatedEventHandler{
		Vtbl: &ICoreWebView2FrameCreatedEventHandlerFn,
		impl: impl,
	}
}
