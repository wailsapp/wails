//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler struct {
	Vtbl *ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerVtbl
	impl ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerImpl
}

func (i *ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerIUnknownAddRef(this *ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerIUnknownRelease(this *ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerInvoke(this *ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler, errorCode uintptr, result *ICoreWebView2CompositionController) uintptr {
	return this.impl.CreateCoreWebView2CompositionControllerCompleted(errorCode, result)
}

type ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerImpl interface {
	IUnknownImpl
	CreateCoreWebView2CompositionControllerCompleted(errorCode uintptr, result *ICoreWebView2CompositionController) uintptr
}

var ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerFn = ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerInvoke),
}

func NewICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler(impl ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerImpl) *ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler {
	return &ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandler{
		Vtbl: &ICoreWebView2CreateCoreWebView2CompositionControllerCompletedHandlerFn,
		impl: impl,
	}
}
