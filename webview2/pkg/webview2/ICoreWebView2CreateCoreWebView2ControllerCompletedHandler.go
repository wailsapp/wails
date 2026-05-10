//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2CreateCoreWebView2ControllerCompletedHandler struct {
	Vtbl *ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerVtbl
	impl ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerImpl
}

func (i *ICoreWebView2CreateCoreWebView2ControllerCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2CreateCoreWebView2ControllerCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerIUnknownAddRef(this *ICoreWebView2CreateCoreWebView2ControllerCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerIUnknownRelease(this *ICoreWebView2CreateCoreWebView2ControllerCompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerInvoke(this *ICoreWebView2CreateCoreWebView2ControllerCompletedHandler, errorCode uintptr, result *ICoreWebView2Controller) uintptr {
	return this.impl.CreateCoreWebView2ControllerCompleted(errorCode, result)
}

type ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerImpl interface {
	IUnknownImpl
	CreateCoreWebView2ControllerCompleted(errorCode uintptr, result *ICoreWebView2Controller) uintptr
}

var ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerFn = ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerInvoke),
}

func NewICoreWebView2CreateCoreWebView2ControllerCompletedHandler(impl ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerImpl) *ICoreWebView2CreateCoreWebView2ControllerCompletedHandler {
	return &ICoreWebView2CreateCoreWebView2ControllerCompletedHandler{
		Vtbl: &ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerFn,
		impl: impl,
	}
}
