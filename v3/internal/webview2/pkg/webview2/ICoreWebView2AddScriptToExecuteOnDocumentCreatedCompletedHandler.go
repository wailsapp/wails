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

func (i *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerIUnknownAddRef(this *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerIUnknownRelease(this *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerInvoke(this *ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler, errorCode uintptr, result *uint16) uintptr {
	_result := UTF16PtrToString(result)
	return this.impl.AddScriptToExecuteOnDocumentCreatedCompleted(errorCode, _result)
}

type ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerImpl interface {
	IUnknownImpl
	AddScriptToExecuteOnDocumentCreatedCompleted(errorCode uintptr, result string) uintptr
}

var ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerFn = ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVtbl{
	IUnknownVtbl {
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
