//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2TrySuspendCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2TrySuspendCompletedHandler struct {
	Vtbl *ICoreWebView2TrySuspendCompletedHandlerVtbl
	impl ICoreWebView2TrySuspendCompletedHandlerImpl
}

func (i *ICoreWebView2TrySuspendCompletedHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2TrySuspendCompletedHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2TrySuspendCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2TrySuspendCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2TrySuspendCompletedHandlerIUnknownAddRef(this *ICoreWebView2TrySuspendCompletedHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2TrySuspendCompletedHandlerIUnknownRelease(this *ICoreWebView2TrySuspendCompletedHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2TrySuspendCompletedHandlerInvoke(this *ICoreWebView2TrySuspendCompletedHandler, errorCode uintptr, result bool) uintptr {
	return this.impl.TrySuspendCompleted(errorCode, result)
}

type ICoreWebView2TrySuspendCompletedHandlerImpl interface {
	IUnknownImpl
	TrySuspendCompleted(errorCode uintptr, result bool) uintptr
}

var ICoreWebView2TrySuspendCompletedHandlerFn = ICoreWebView2TrySuspendCompletedHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2TrySuspendCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2TrySuspendCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2TrySuspendCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2TrySuspendCompletedHandlerInvoke),
}

func NewICoreWebView2TrySuspendCompletedHandler(impl ICoreWebView2TrySuspendCompletedHandlerImpl) *ICoreWebView2TrySuspendCompletedHandler {
	return &ICoreWebView2TrySuspendCompletedHandler{
		Vtbl: &ICoreWebView2TrySuspendCompletedHandlerFn,
		impl: impl,
	}
}
