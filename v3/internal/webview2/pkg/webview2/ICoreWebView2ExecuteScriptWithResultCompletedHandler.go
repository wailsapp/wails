//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2ExecuteScriptWithResultCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ExecuteScriptWithResultCompletedHandler struct {
	Vtbl *ICoreWebView2ExecuteScriptWithResultCompletedHandlerVtbl
	impl ICoreWebView2ExecuteScriptWithResultCompletedHandlerImpl
}

func (i *ICoreWebView2ExecuteScriptWithResultCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2ExecuteScriptWithResultCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2ExecuteScriptWithResultCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ExecuteScriptWithResultCompletedHandlerIUnknownAddRef(this *ICoreWebView2ExecuteScriptWithResultCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2ExecuteScriptWithResultCompletedHandlerIUnknownRelease(this *ICoreWebView2ExecuteScriptWithResultCompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2ExecuteScriptWithResultCompletedHandlerInvoke(this *ICoreWebView2ExecuteScriptWithResultCompletedHandler, errorCode uintptr, result *ICoreWebView2ExecuteScriptResult) uintptr {
	return this.impl.ExecuteScriptWithResultCompleted(errorCode, result)
}

type ICoreWebView2ExecuteScriptWithResultCompletedHandlerImpl interface {
	IUnknownImpl
	ExecuteScriptWithResultCompleted(errorCode uintptr, result *ICoreWebView2ExecuteScriptResult) uintptr
}

var ICoreWebView2ExecuteScriptWithResultCompletedHandlerFn = ICoreWebView2ExecuteScriptWithResultCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2ExecuteScriptWithResultCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ExecuteScriptWithResultCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ExecuteScriptWithResultCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ExecuteScriptWithResultCompletedHandlerInvoke),
}

func NewICoreWebView2ExecuteScriptWithResultCompletedHandler(impl ICoreWebView2ExecuteScriptWithResultCompletedHandlerImpl) *ICoreWebView2ExecuteScriptWithResultCompletedHandler {
	return &ICoreWebView2ExecuteScriptWithResultCompletedHandler{
		Vtbl: &ICoreWebView2ExecuteScriptWithResultCompletedHandlerFn,
		impl: impl,
	}
}
