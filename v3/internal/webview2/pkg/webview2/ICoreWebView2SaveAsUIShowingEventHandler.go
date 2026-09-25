//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2SaveAsUIShowingEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2SaveAsUIShowingEventHandler struct {
	Vtbl *ICoreWebView2SaveAsUIShowingEventHandlerVtbl
	impl ICoreWebView2SaveAsUIShowingEventHandlerImpl
}

func (i *ICoreWebView2SaveAsUIShowingEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2SaveAsUIShowingEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2SaveAsUIShowingEventHandlerIUnknownQueryInterface(this *ICoreWebView2SaveAsUIShowingEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2SaveAsUIShowingEventHandlerIUnknownAddRef(this *ICoreWebView2SaveAsUIShowingEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2SaveAsUIShowingEventHandlerIUnknownRelease(this *ICoreWebView2SaveAsUIShowingEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2SaveAsUIShowingEventHandlerInvoke(this *ICoreWebView2SaveAsUIShowingEventHandler, sender *ICoreWebView2, args *ICoreWebView2SaveAsUIShowingEventArgs) uintptr {
	return this.impl.SaveAsUIShowing(sender, args)
}

type ICoreWebView2SaveAsUIShowingEventHandlerImpl interface {
	IUnknownImpl
	SaveAsUIShowing(sender *ICoreWebView2, args *ICoreWebView2SaveAsUIShowingEventArgs) uintptr
}

var ICoreWebView2SaveAsUIShowingEventHandlerFn = ICoreWebView2SaveAsUIShowingEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2SaveAsUIShowingEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2SaveAsUIShowingEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2SaveAsUIShowingEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2SaveAsUIShowingEventHandlerInvoke),
}

func NewICoreWebView2SaveAsUIShowingEventHandler(impl ICoreWebView2SaveAsUIShowingEventHandlerImpl) *ICoreWebView2SaveAsUIShowingEventHandler {
	return &ICoreWebView2SaveAsUIShowingEventHandler{
		Vtbl: &ICoreWebView2SaveAsUIShowingEventHandlerFn,
		impl: impl,
	}
}
