//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2BrowserProcessExitedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2BrowserProcessExitedEventHandler struct {
	Vtbl *ICoreWebView2BrowserProcessExitedEventHandlerVtbl
	impl ICoreWebView2BrowserProcessExitedEventHandlerImpl
}

func (i *ICoreWebView2BrowserProcessExitedEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2BrowserProcessExitedEventHandlerIUnknownQueryInterface(this *ICoreWebView2BrowserProcessExitedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2BrowserProcessExitedEventHandlerIUnknownAddRef(this *ICoreWebView2BrowserProcessExitedEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2BrowserProcessExitedEventHandlerIUnknownRelease(this *ICoreWebView2BrowserProcessExitedEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2BrowserProcessExitedEventHandlerInvoke(this *ICoreWebView2BrowserProcessExitedEventHandler, sender *ICoreWebView2Environment, args *ICoreWebView2BrowserProcessExitedEventArgs) uintptr {
	return this.impl.BrowserProcessExited(sender, args)
}

type ICoreWebView2BrowserProcessExitedEventHandlerImpl interface {
	IUnknownImpl
	BrowserProcessExited(sender *ICoreWebView2Environment, args *ICoreWebView2BrowserProcessExitedEventArgs) uintptr
}

var ICoreWebView2BrowserProcessExitedEventHandlerFn = ICoreWebView2BrowserProcessExitedEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2BrowserProcessExitedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2BrowserProcessExitedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2BrowserProcessExitedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2BrowserProcessExitedEventHandlerInvoke),
}

func NewICoreWebView2BrowserProcessExitedEventHandler(impl ICoreWebView2BrowserProcessExitedEventHandlerImpl) *ICoreWebView2BrowserProcessExitedEventHandler {
	return &ICoreWebView2BrowserProcessExitedEventHandler{
		Vtbl: &ICoreWebView2BrowserProcessExitedEventHandlerFn,
		impl: impl,
	}
}
