//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2SourceChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2SourceChangedEventHandler struct {
	Vtbl *ICoreWebView2SourceChangedEventHandlerVtbl
	impl ICoreWebView2SourceChangedEventHandlerImpl
}

func (i *ICoreWebView2SourceChangedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2SourceChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2SourceChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2SourceChangedEventHandlerIUnknownAddRef(this *ICoreWebView2SourceChangedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2SourceChangedEventHandlerIUnknownRelease(this *ICoreWebView2SourceChangedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2SourceChangedEventHandlerInvoke(this *ICoreWebView2SourceChangedEventHandler, sender *ICoreWebView2, args *ICoreWebView2SourceChangedEventArgs) uintptr {
	return this.impl.SourceChanged(sender, args)
}

type ICoreWebView2SourceChangedEventHandlerImpl interface {
	IUnknownImpl
	SourceChanged(sender *ICoreWebView2, args *ICoreWebView2SourceChangedEventArgs) uintptr
}

var ICoreWebView2SourceChangedEventHandlerFn = ICoreWebView2SourceChangedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2SourceChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2SourceChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2SourceChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2SourceChangedEventHandlerInvoke),
}

func NewICoreWebView2SourceChangedEventHandler(impl ICoreWebView2SourceChangedEventHandlerImpl) *ICoreWebView2SourceChangedEventHandler {
	return &ICoreWebView2SourceChangedEventHandler{
		Vtbl: &ICoreWebView2SourceChangedEventHandlerFn,
		impl: impl,
	}
}
