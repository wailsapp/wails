//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2EstimatedEndTimeChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2EstimatedEndTimeChangedEventHandler struct {
	Vtbl *ICoreWebView2EstimatedEndTimeChangedEventHandlerVtbl
	impl ICoreWebView2EstimatedEndTimeChangedEventHandlerImpl
}

func (i *ICoreWebView2EstimatedEndTimeChangedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2EstimatedEndTimeChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2EstimatedEndTimeChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2EstimatedEndTimeChangedEventHandlerIUnknownAddRef(this *ICoreWebView2EstimatedEndTimeChangedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2EstimatedEndTimeChangedEventHandlerIUnknownRelease(this *ICoreWebView2EstimatedEndTimeChangedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2EstimatedEndTimeChangedEventHandlerInvoke(this *ICoreWebView2EstimatedEndTimeChangedEventHandler, sender *ICoreWebView2DownloadOperation, args *IUnknown) uintptr {
	return this.impl.EstimatedEndTimeChanged(sender, args)
}

type ICoreWebView2EstimatedEndTimeChangedEventHandlerImpl interface {
	IUnknownImpl
	EstimatedEndTimeChanged(sender *ICoreWebView2DownloadOperation, args *IUnknown) uintptr
}

var ICoreWebView2EstimatedEndTimeChangedEventHandlerFn = ICoreWebView2EstimatedEndTimeChangedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2EstimatedEndTimeChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2EstimatedEndTimeChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2EstimatedEndTimeChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2EstimatedEndTimeChangedEventHandlerInvoke),
}

func NewICoreWebView2EstimatedEndTimeChangedEventHandler(impl ICoreWebView2EstimatedEndTimeChangedEventHandlerImpl) *ICoreWebView2EstimatedEndTimeChangedEventHandler {
	return &ICoreWebView2EstimatedEndTimeChangedEventHandler{
		Vtbl: &ICoreWebView2EstimatedEndTimeChangedEventHandlerFn,
		impl: impl,
	}
}
