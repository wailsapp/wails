//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2ExecuteScriptCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ExecuteScriptCompletedHandler struct {
	Vtbl *ICoreWebView2ExecuteScriptCompletedHandlerVtbl
	impl ICoreWebView2ExecuteScriptCompletedHandlerImpl
}

func (i *ICoreWebView2ExecuteScriptCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2ExecuteScriptCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2ExecuteScriptCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ExecuteScriptCompletedHandlerIUnknownAddRef(this *ICoreWebView2ExecuteScriptCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2ExecuteScriptCompletedHandlerIUnknownRelease(this *ICoreWebView2ExecuteScriptCompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2ExecuteScriptCompletedHandlerInvoke(this *ICoreWebView2ExecuteScriptCompletedHandler, errorCode uintptr, result string) uintptr {
	return this.impl.ExecuteScriptCompleted(errorCode, result)
}

type ICoreWebView2ExecuteScriptCompletedHandlerImpl interface {
	IUnknownImpl
	ExecuteScriptCompleted(errorCode uintptr, result string) uintptr
}

var ICoreWebView2ExecuteScriptCompletedHandlerFn = ICoreWebView2ExecuteScriptCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2ExecuteScriptCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ExecuteScriptCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ExecuteScriptCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ExecuteScriptCompletedHandlerInvoke),
}

func NewICoreWebView2ExecuteScriptCompletedHandler(impl ICoreWebView2ExecuteScriptCompletedHandlerImpl) *ICoreWebView2ExecuteScriptCompletedHandler {
	return &ICoreWebView2ExecuteScriptCompletedHandler{
		Vtbl: &ICoreWebView2ExecuteScriptCompletedHandlerFn,
		impl: impl,
	}
}
