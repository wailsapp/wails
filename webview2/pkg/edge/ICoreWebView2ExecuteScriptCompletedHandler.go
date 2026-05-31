package edge

import (
	"unsafe"
)

type _ICoreWebView2ExecuteScriptCompletedHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2ExecuteScriptCompletedHandler struct {
	vtbl *_ICoreWebView2ExecuteScriptCompletedHandlerVtbl
	impl _ICoreWebView2ExecuteScriptCompletedHandlerImpl
}

func (i *iCoreWebView2ExecuteScriptCompletedHandler) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *iCoreWebView2ExecuteScriptCompletedHandler) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func _ICoreWebView2ExecuteScriptCompletedHandlerIUnknownQueryInterface(this *iCoreWebView2ExecuteScriptCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2ExecuteScriptCompletedHandlerIUnknownAddRef(this *iCoreWebView2ExecuteScriptCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2ExecuteScriptCompletedHandlerIUnknownRelease(this *iCoreWebView2ExecuteScriptCompletedHandler) uintptr {
	return this.impl.Release()
}

func iCoreWebView2ExecuteScriptCompletedHandlerInvoke(this *iCoreWebView2ExecuteScriptCompletedHandler, errorCode uintptr, executedScript *uint16) uintptr {
	return this.impl.ExecuteScriptCompleted(errorCode, executedScript)
}

type _ICoreWebView2ExecuteScriptCompletedHandlerImpl interface {
	_IUnknownImpl
	ExecuteScriptCompleted(errorCode uintptr, executedScript *uint16) uintptr
}

var _ICoreWebView2ExecuteScriptCompletedHandlerFn = _ICoreWebView2ExecuteScriptCompletedHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2ExecuteScriptCompletedHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2ExecuteScriptCompletedHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2ExecuteScriptCompletedHandlerIUnknownRelease),
	},
	NewComProc(iCoreWebView2ExecuteScriptCompletedHandlerInvoke),
}

func newICoreWebView2ExecuteScriptCompletedHandler(impl _ICoreWebView2ExecuteScriptCompletedHandlerImpl) *iCoreWebView2ExecuteScriptCompletedHandler {
	return &iCoreWebView2ExecuteScriptCompletedHandler{
		vtbl: &_ICoreWebView2ExecuteScriptCompletedHandlerFn,
		impl: impl,
	}
}
