//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandler struct {
	Vtbl *ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerVtbl
	impl ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerImpl
}

func (i *ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerIUnknownAddRef(this *ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerIUnknownRelease(this *ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerInvoke(this *ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandler, sender *ICoreWebView2, args *IUnknown) uintptr {
	return this.impl.IsDefaultDownloadDialogOpenChanged(sender, args)
}

type ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerImpl interface {
	IUnknownImpl
	IsDefaultDownloadDialogOpenChanged(sender *ICoreWebView2, args *IUnknown) uintptr
}

var ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerFn = ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerInvoke),
}

func NewICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandler(impl ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerImpl) *ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandler {
	return &ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandler{
		Vtbl: &ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandlerFn,
		impl: impl,
	}
}
