//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler struct {
	Vtbl *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVtbl
	impl ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerImpl
}

func (i *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerIUnknownAddRef(this *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerIUnknownRelease(this *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerInvoke(this *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler, errorCode uintptr, result string) uintptr {
	return this.impl.AddScriptToExecuteOnDocumentCreatedCompleted(errorCode, result)
}

type ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerImpl interface {
	IUnknownImpl
	AddScriptToExecuteOnDocumentCreatedCompleted(errorCode uintptr, result string) uintptr
}

var ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerFn = ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerInvoke),
}

func NewICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler(impl ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerImpl) *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler {
	return &ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler{
		Vtbl: &ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerFn,
		impl: impl,
	}
}
