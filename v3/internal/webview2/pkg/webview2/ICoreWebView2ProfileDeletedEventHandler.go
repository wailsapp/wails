//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2ProfileDeletedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ProfileDeletedEventHandler struct {
	Vtbl *ICoreWebView2ProfileDeletedEventHandlerVtbl
	impl ICoreWebView2ProfileDeletedEventHandlerImpl
}

func (i *ICoreWebView2ProfileDeletedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2ProfileDeletedEventHandlerIUnknownQueryInterface(this *ICoreWebView2ProfileDeletedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ProfileDeletedEventHandlerIUnknownAddRef(this *ICoreWebView2ProfileDeletedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2ProfileDeletedEventHandlerIUnknownRelease(this *ICoreWebView2ProfileDeletedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2ProfileDeletedEventHandlerInvoke(this *ICoreWebView2ProfileDeletedEventHandler, sender *ICoreWebView2Profile, args *IUnknown) uintptr {
	return this.impl.ProfileDeleted(sender, args)
}

type ICoreWebView2ProfileDeletedEventHandlerImpl interface {
	IUnknownImpl
	ProfileDeleted(sender *ICoreWebView2Profile, args *IUnknown) uintptr
}

var ICoreWebView2ProfileDeletedEventHandlerFn = ICoreWebView2ProfileDeletedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2ProfileDeletedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ProfileDeletedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ProfileDeletedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ProfileDeletedEventHandlerInvoke),
}

func NewICoreWebView2ProfileDeletedEventHandler(impl ICoreWebView2ProfileDeletedEventHandlerImpl) *ICoreWebView2ProfileDeletedEventHandler {
	return &ICoreWebView2ProfileDeletedEventHandler{
		Vtbl: &ICoreWebView2ProfileDeletedEventHandlerFn,
		impl: impl,
	}
}
