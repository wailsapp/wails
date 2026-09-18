//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2BytesReceivedChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2BytesReceivedChangedEventHandler struct {
	Vtbl *ICoreWebView2BytesReceivedChangedEventHandlerVtbl
	impl ICoreWebView2BytesReceivedChangedEventHandlerImpl
}

func (i *ICoreWebView2BytesReceivedChangedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2BytesReceivedChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2BytesReceivedChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2BytesReceivedChangedEventHandlerIUnknownAddRef(this *ICoreWebView2BytesReceivedChangedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2BytesReceivedChangedEventHandlerIUnknownRelease(this *ICoreWebView2BytesReceivedChangedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2BytesReceivedChangedEventHandlerInvoke(this *ICoreWebView2BytesReceivedChangedEventHandler, sender *ICoreWebView2DownloadOperation, args *IUnknown) uintptr {
	return this.impl.BytesReceivedChanged(sender, args)
}

type ICoreWebView2BytesReceivedChangedEventHandlerImpl interface {
	IUnknownImpl
	BytesReceivedChanged(sender *ICoreWebView2DownloadOperation, args *IUnknown) uintptr
}

var ICoreWebView2BytesReceivedChangedEventHandlerFn = ICoreWebView2BytesReceivedChangedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2BytesReceivedChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2BytesReceivedChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2BytesReceivedChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2BytesReceivedChangedEventHandlerInvoke),
}

func NewICoreWebView2BytesReceivedChangedEventHandler(impl ICoreWebView2BytesReceivedChangedEventHandlerImpl) *ICoreWebView2BytesReceivedChangedEventHandler {
	return &ICoreWebView2BytesReceivedChangedEventHandler{
		Vtbl: &ICoreWebView2BytesReceivedChangedEventHandlerFn,
		impl: impl,
	}
}
