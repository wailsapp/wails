//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler struct {
	Vtbl *ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl
	impl ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl
}

func (i *ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownAddRef(this *ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownRelease(this *ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerInvoke(this *ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, errorCode uintptr, result *ICoreWebView2Environment) uintptr {
	return this.impl.CreateCoreWebView2EnvironmentCompleted(errorCode, result)
}

type ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl interface {
	IUnknownImpl
	CreateCoreWebView2EnvironmentCompleted(errorCode uintptr, result *ICoreWebView2Environment) uintptr
}

var ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerFn = ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerInvoke),
}

func NewICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler(impl ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerImpl) *ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler {
	return &ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler{
		Vtbl: &ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerFn,
		impl: impl,
	}
}
