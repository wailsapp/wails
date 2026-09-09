//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2SaveFileSecurityCheckStartingEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2SaveFileSecurityCheckStartingEventHandler struct {
	Vtbl *ICoreWebView2SaveFileSecurityCheckStartingEventHandlerVtbl
	impl ICoreWebView2SaveFileSecurityCheckStartingEventHandlerImpl
}

func (i *ICoreWebView2SaveFileSecurityCheckStartingEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2SaveFileSecurityCheckStartingEventHandlerIUnknownQueryInterface(this *ICoreWebView2SaveFileSecurityCheckStartingEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2SaveFileSecurityCheckStartingEventHandlerIUnknownAddRef(this *ICoreWebView2SaveFileSecurityCheckStartingEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2SaveFileSecurityCheckStartingEventHandlerIUnknownRelease(this *ICoreWebView2SaveFileSecurityCheckStartingEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2SaveFileSecurityCheckStartingEventHandlerInvoke(this *ICoreWebView2SaveFileSecurityCheckStartingEventHandler, sender *ICoreWebView2, args *ICoreWebView2SaveFileSecurityCheckStartingEventArgs) uintptr {
	return this.impl.SaveFileSecurityCheckStarting(sender, args)
}

type ICoreWebView2SaveFileSecurityCheckStartingEventHandlerImpl interface {
	IUnknownImpl
	SaveFileSecurityCheckStarting(sender *ICoreWebView2, args *ICoreWebView2SaveFileSecurityCheckStartingEventArgs) uintptr
}

var ICoreWebView2SaveFileSecurityCheckStartingEventHandlerFn = ICoreWebView2SaveFileSecurityCheckStartingEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2SaveFileSecurityCheckStartingEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2SaveFileSecurityCheckStartingEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2SaveFileSecurityCheckStartingEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2SaveFileSecurityCheckStartingEventHandlerInvoke),
}

func NewICoreWebView2SaveFileSecurityCheckStartingEventHandler(impl ICoreWebView2SaveFileSecurityCheckStartingEventHandlerImpl) *ICoreWebView2SaveFileSecurityCheckStartingEventHandler {
	return &ICoreWebView2SaveFileSecurityCheckStartingEventHandler{
		Vtbl: &ICoreWebView2SaveFileSecurityCheckStartingEventHandlerFn,
		impl: impl,
	}
}
