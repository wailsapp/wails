//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2DownloadStartingEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2DownloadStartingEventHandler struct {
	Vtbl *ICoreWebView2DownloadStartingEventHandlerVtbl
	impl ICoreWebView2DownloadStartingEventHandlerImpl
}

func (i *ICoreWebView2DownloadStartingEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2DownloadStartingEventHandlerIUnknownQueryInterface(this *ICoreWebView2DownloadStartingEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2DownloadStartingEventHandlerIUnknownAddRef(this *ICoreWebView2DownloadStartingEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2DownloadStartingEventHandlerIUnknownRelease(this *ICoreWebView2DownloadStartingEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2DownloadStartingEventHandlerInvoke(this *ICoreWebView2DownloadStartingEventHandler, sender *ICoreWebView2, args *ICoreWebView2DownloadStartingEventArgs) uintptr {
	return this.impl.DownloadStarting(sender, args)
}

type ICoreWebView2DownloadStartingEventHandlerImpl interface {
	IUnknownImpl
	DownloadStarting(sender *ICoreWebView2, args *ICoreWebView2DownloadStartingEventArgs) uintptr
}

var ICoreWebView2DownloadStartingEventHandlerFn = ICoreWebView2DownloadStartingEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2DownloadStartingEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2DownloadStartingEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2DownloadStartingEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2DownloadStartingEventHandlerInvoke),
}

func NewICoreWebView2DownloadStartingEventHandler(impl ICoreWebView2DownloadStartingEventHandlerImpl) *ICoreWebView2DownloadStartingEventHandler {
	return &ICoreWebView2DownloadStartingEventHandler{
		Vtbl: &ICoreWebView2DownloadStartingEventHandlerFn,
		impl: impl,
	}
}
