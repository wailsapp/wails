//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2WebResourceResponseViewGetContentCompletedHandler struct {
	Vtbl *ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerVtbl
	impl ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerImpl
}

func (i *ICoreWebView2WebResourceResponseViewGetContentCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2WebResourceResponseViewGetContentCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerIUnknownAddRef(this *ICoreWebView2WebResourceResponseViewGetContentCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerIUnknownRelease(this *ICoreWebView2WebResourceResponseViewGetContentCompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerInvoke(this *ICoreWebView2WebResourceResponseViewGetContentCompletedHandler, errorCode uintptr, result *IStream) uintptr {
	return this.impl.WebResourceResponseViewGetContentCompleted(errorCode, result)
}

type ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerImpl interface {
	IUnknownImpl
	WebResourceResponseViewGetContentCompleted(errorCode uintptr, result *IStream) uintptr
}

var ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerFn = ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerInvoke),
}

func NewICoreWebView2WebResourceResponseViewGetContentCompletedHandler(impl ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerImpl) *ICoreWebView2WebResourceResponseViewGetContentCompletedHandler {
	return &ICoreWebView2WebResourceResponseViewGetContentCompletedHandler{
		Vtbl: &ICoreWebView2WebResourceResponseViewGetContentCompletedHandlerFn,
		impl: impl,
	}
}
