package edge

import (
	"unsafe"
)

type _ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler struct {
	vtbl *_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVtbl
	impl _ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerImpl
}

func (i *iCoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *iCoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func _ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerIUnknownQueryInterface(this *iCoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func iCoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerIUnknownAddRef(this *iCoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerIUnknownRelease(this *iCoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerInvoke(this *iCoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler, errorCode uintptr, addedScript *ICoreWebView2Controller) uintptr {
	return this.impl.AddScriptToExecuteOnDocumentCreatedCompleted(errorCode, addedScript)
}

type _ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerImpl interface {
	_IUnknownImpl
	AddScriptToExecuteOnDocumentCreatedCompleted(errorCode uintptr, addedScript *ICoreWebView2Controller) uintptr
}

var _ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerFn = _ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerIUnknownQueryInterface),
		NewComProc(iCoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerInvoke),
}

func newICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler(impl _ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerImpl) *iCoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler {
	return &iCoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler{
		vtbl: &_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerFn,
		impl: impl,
	}
}
