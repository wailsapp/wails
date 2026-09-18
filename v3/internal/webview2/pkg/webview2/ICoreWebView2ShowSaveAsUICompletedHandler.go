//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2ShowSaveAsUICompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ShowSaveAsUICompletedHandler struct {
	Vtbl *ICoreWebView2ShowSaveAsUICompletedHandlerVtbl
	impl ICoreWebView2ShowSaveAsUICompletedHandlerImpl
}

func (i *ICoreWebView2ShowSaveAsUICompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2ShowSaveAsUICompletedHandlerIUnknownQueryInterface(this *ICoreWebView2ShowSaveAsUICompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ShowSaveAsUICompletedHandlerIUnknownAddRef(this *ICoreWebView2ShowSaveAsUICompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2ShowSaveAsUICompletedHandlerIUnknownRelease(this *ICoreWebView2ShowSaveAsUICompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2ShowSaveAsUICompletedHandlerInvoke(this *ICoreWebView2ShowSaveAsUICompletedHandler, errorCode uintptr, result COREWEBVIEW2_SAVE_AS_UI_RESULT) uintptr {
	return this.impl.ShowSaveAsUICompleted(errorCode, result)
}

type ICoreWebView2ShowSaveAsUICompletedHandlerImpl interface {
	IUnknownImpl
	ShowSaveAsUICompleted(errorCode uintptr, result COREWEBVIEW2_SAVE_AS_UI_RESULT) uintptr
}

var ICoreWebView2ShowSaveAsUICompletedHandlerFn = ICoreWebView2ShowSaveAsUICompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2ShowSaveAsUICompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ShowSaveAsUICompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ShowSaveAsUICompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ShowSaveAsUICompletedHandlerInvoke),
}

func NewICoreWebView2ShowSaveAsUICompletedHandler(impl ICoreWebView2ShowSaveAsUICompletedHandlerImpl) *ICoreWebView2ShowSaveAsUICompletedHandler {
	return &ICoreWebView2ShowSaveAsUICompletedHandler{
		Vtbl: &ICoreWebView2ShowSaveAsUICompletedHandlerFn,
		impl: impl,
	}
}
