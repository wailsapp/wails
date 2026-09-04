//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2PrintToPdfStreamCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2PrintToPdfStreamCompletedHandler struct {
	Vtbl *ICoreWebView2PrintToPdfStreamCompletedHandlerVtbl
	impl ICoreWebView2PrintToPdfStreamCompletedHandlerImpl
}

func (i *ICoreWebView2PrintToPdfStreamCompletedHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2PrintToPdfStreamCompletedHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2PrintToPdfStreamCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2PrintToPdfStreamCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2PrintToPdfStreamCompletedHandlerIUnknownAddRef(this *ICoreWebView2PrintToPdfStreamCompletedHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2PrintToPdfStreamCompletedHandlerIUnknownRelease(this *ICoreWebView2PrintToPdfStreamCompletedHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2PrintToPdfStreamCompletedHandlerInvoke(this *ICoreWebView2PrintToPdfStreamCompletedHandler, errorCode uintptr, result *IStream) uintptr {
	return this.impl.PrintToPdfStreamCompleted(errorCode, result)
}

type ICoreWebView2PrintToPdfStreamCompletedHandlerImpl interface {
	IUnknownImpl
	PrintToPdfStreamCompleted(errorCode uintptr, result *IStream) uintptr
}

var ICoreWebView2PrintToPdfStreamCompletedHandlerFn = ICoreWebView2PrintToPdfStreamCompletedHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2PrintToPdfStreamCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2PrintToPdfStreamCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2PrintToPdfStreamCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2PrintToPdfStreamCompletedHandlerInvoke),
}

func NewICoreWebView2PrintToPdfStreamCompletedHandler(impl ICoreWebView2PrintToPdfStreamCompletedHandlerImpl) *ICoreWebView2PrintToPdfStreamCompletedHandler {
	return &ICoreWebView2PrintToPdfStreamCompletedHandler{
		Vtbl: &ICoreWebView2PrintToPdfStreamCompletedHandlerFn,
		impl: impl,
	}
}
