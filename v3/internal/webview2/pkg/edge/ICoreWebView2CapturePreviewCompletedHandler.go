//go:build windows

package edge

import "unsafe"

type _ICoreWebView2CapturePreviewCompletedHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2CapturePreviewCompletedHandler struct {
	vtbl *_ICoreWebView2CapturePreviewCompletedHandlerVtbl
	impl _ICoreWebView2CapturePreviewCompletedHandlerImpl
}

func (i *iCoreWebView2CapturePreviewCompletedHandler) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(ret)
}

func (i *iCoreWebView2CapturePreviewCompletedHandler) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(ret)
}

func iCoreWebView2CapturePreviewCompletedHandlerIUnknownQueryInterface(this *iCoreWebView2CapturePreviewCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func iCoreWebView2CapturePreviewCompletedHandlerIUnknownAddRef(this *iCoreWebView2CapturePreviewCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func iCoreWebView2CapturePreviewCompletedHandlerIUnknownRelease(this *iCoreWebView2CapturePreviewCompletedHandler) uintptr {
	return this.impl.Release()
}

func iCoreWebView2CapturePreviewCompletedHandlerInvoke(this *iCoreWebView2CapturePreviewCompletedHandler, errorCode uintptr) uintptr {
	return this.impl.CapturePreviewCompleted(errorCode)
}

type _ICoreWebView2CapturePreviewCompletedHandlerImpl interface {
	_IUnknownImpl
	CapturePreviewCompleted(errorCode uintptr) uintptr
}

var _ICoreWebView2CapturePreviewCompletedHandlerFn = _ICoreWebView2CapturePreviewCompletedHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(iCoreWebView2CapturePreviewCompletedHandlerIUnknownQueryInterface),
		NewComProc(iCoreWebView2CapturePreviewCompletedHandlerIUnknownAddRef),
		NewComProc(iCoreWebView2CapturePreviewCompletedHandlerIUnknownRelease),
	},
	NewComProc(iCoreWebView2CapturePreviewCompletedHandlerInvoke),
}

func newICoreWebView2CapturePreviewCompletedHandler(impl _ICoreWebView2CapturePreviewCompletedHandlerImpl) *iCoreWebView2CapturePreviewCompletedHandler {
	return &iCoreWebView2CapturePreviewCompletedHandler{
		vtbl: &_ICoreWebView2CapturePreviewCompletedHandlerFn,
		impl: impl,
	}
}
