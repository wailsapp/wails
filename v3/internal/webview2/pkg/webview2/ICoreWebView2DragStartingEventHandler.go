//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2DragStartingEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2DragStartingEventHandler struct {
	Vtbl *ICoreWebView2DragStartingEventHandlerVtbl
	impl ICoreWebView2DragStartingEventHandlerImpl
}

func (i *ICoreWebView2DragStartingEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2DragStartingEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2DragStartingEventHandlerIUnknownQueryInterface(this *ICoreWebView2DragStartingEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2DragStartingEventHandlerIUnknownAddRef(this *ICoreWebView2DragStartingEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2DragStartingEventHandlerIUnknownRelease(this *ICoreWebView2DragStartingEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2DragStartingEventHandlerInvoke(this *ICoreWebView2DragStartingEventHandler, sender *ICoreWebView2CompositionController, args *ICoreWebView2DragStartingEventArgs) uintptr {
	return this.impl.DragStarting(sender, args)
}

type ICoreWebView2DragStartingEventHandlerImpl interface {
	IUnknownImpl
	DragStarting(sender *ICoreWebView2CompositionController, args *ICoreWebView2DragStartingEventArgs) uintptr
}

var ICoreWebView2DragStartingEventHandlerFn = ICoreWebView2DragStartingEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2DragStartingEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2DragStartingEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2DragStartingEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2DragStartingEventHandlerInvoke),
}

func NewICoreWebView2DragStartingEventHandler(impl ICoreWebView2DragStartingEventHandlerImpl) *ICoreWebView2DragStartingEventHandler {
	return &ICoreWebView2DragStartingEventHandler{
		Vtbl: &ICoreWebView2DragStartingEventHandlerFn,
		impl: impl,
	}
}
