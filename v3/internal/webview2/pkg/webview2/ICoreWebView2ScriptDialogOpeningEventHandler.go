//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2ScriptDialogOpeningEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ScriptDialogOpeningEventHandler struct {
	Vtbl *ICoreWebView2ScriptDialogOpeningEventHandlerVtbl
	impl ICoreWebView2ScriptDialogOpeningEventHandlerImpl
}

func (i *ICoreWebView2ScriptDialogOpeningEventHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2ScriptDialogOpeningEventHandlerIUnknownQueryInterface(this *ICoreWebView2ScriptDialogOpeningEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2ScriptDialogOpeningEventHandlerIUnknownAddRef(this *ICoreWebView2ScriptDialogOpeningEventHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2ScriptDialogOpeningEventHandlerIUnknownRelease(this *ICoreWebView2ScriptDialogOpeningEventHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2ScriptDialogOpeningEventHandlerInvoke(this *ICoreWebView2ScriptDialogOpeningEventHandler, sender *ICoreWebView2, args *ICoreWebView2ScriptDialogOpeningEventArgs) uintptr {
	return this.impl.ScriptDialogOpening(sender, args)
}

type ICoreWebView2ScriptDialogOpeningEventHandlerImpl interface {
	IUnknownImpl
	ScriptDialogOpening(sender *ICoreWebView2, args *ICoreWebView2ScriptDialogOpeningEventArgs) uintptr
}

var ICoreWebView2ScriptDialogOpeningEventHandlerFn = ICoreWebView2ScriptDialogOpeningEventHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2ScriptDialogOpeningEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2ScriptDialogOpeningEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2ScriptDialogOpeningEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2ScriptDialogOpeningEventHandlerInvoke),
}

func NewICoreWebView2ScriptDialogOpeningEventHandler(impl ICoreWebView2ScriptDialogOpeningEventHandlerImpl) *ICoreWebView2ScriptDialogOpeningEventHandler {
	return &ICoreWebView2ScriptDialogOpeningEventHandler{
		Vtbl: &ICoreWebView2ScriptDialogOpeningEventHandlerFn,
		impl: impl,
	}
}
