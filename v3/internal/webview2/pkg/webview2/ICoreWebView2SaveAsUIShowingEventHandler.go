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

func (i *ICoreWebView2SaveAsUIShowingEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2SaveAsUIShowingEventHandlerIUnknownQueryInterface(this *ICoreWebView2SaveAsUIShowingEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2SaveAsUIShowingEventHandlerIUnknownAddRef(this *ICoreWebView2SaveAsUIShowingEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2SaveAsUIShowingEventHandlerIUnknownRelease(this *ICoreWebView2SaveAsUIShowingEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2SaveAsUIShowingEventHandlerInvoke(this *ICoreWebView2SaveAsUIShowingEventHandler, sender *ICoreWebView2, args *ICoreWebView2SaveAsUIShowingEventArgs) uintptr {
	return this.impl.SaveAsUIShowing(sender, args)
}

type ICoreWebView2SaveAsUIShowingEventHandlerImpl interface {
	IUnknownImpl
	SaveAsUIShowing(sender *ICoreWebView2, args *ICoreWebView2SaveAsUIShowingEventArgs) uintptr
}

var ICoreWebView2SaveAsUIShowingEventHandlerFn = ICoreWebView2SaveAsUIShowingEventHandlerVtbl{
	IUnknownVtbl{
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
