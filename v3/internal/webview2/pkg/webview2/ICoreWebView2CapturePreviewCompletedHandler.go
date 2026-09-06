//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2CapturePreviewCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2CapturePreviewCompletedHandler struct {
	Vtbl *ICoreWebView2CapturePreviewCompletedHandlerVtbl
	impl ICoreWebView2CapturePreviewCompletedHandlerImpl
}

func (i *ICoreWebView2CapturePreviewCompletedHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2CapturePreviewCompletedHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2CapturePreviewCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2CapturePreviewCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2CapturePreviewCompletedHandlerIUnknownAddRef(this *ICoreWebView2CapturePreviewCompletedHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2CapturePreviewCompletedHandlerIUnknownRelease(this *ICoreWebView2CapturePreviewCompletedHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2CapturePreviewCompletedHandlerInvoke(this *ICoreWebView2CapturePreviewCompletedHandler, errorCode uintptr) uintptr {
	return this.impl.CapturePreviewCompleted(errorCode)
}

type ICoreWebView2CapturePreviewCompletedHandlerImpl interface {
	IUnknownImpl
	CapturePreviewCompleted(errorCode uintptr) uintptr
}

var ICoreWebView2CapturePreviewCompletedHandlerFn = ICoreWebView2CapturePreviewCompletedHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2CapturePreviewCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2CapturePreviewCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2CapturePreviewCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2CapturePreviewCompletedHandlerInvoke),
}

func NewICoreWebView2CapturePreviewCompletedHandler(impl ICoreWebView2CapturePreviewCompletedHandlerImpl) *ICoreWebView2CapturePreviewCompletedHandler {
	return &ICoreWebView2CapturePreviewCompletedHandler{
		Vtbl: &ICoreWebView2CapturePreviewCompletedHandlerFn,
		impl: impl,
	}
}
