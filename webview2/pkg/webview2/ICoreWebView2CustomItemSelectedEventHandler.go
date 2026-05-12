//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2CustomItemSelectedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2CustomItemSelectedEventHandler struct {
	Vtbl *ICoreWebView2CustomItemSelectedEventHandlerVtbl
	impl ICoreWebView2CustomItemSelectedEventHandlerImpl
}

func (i *ICoreWebView2CustomItemSelectedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2CustomItemSelectedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2CustomItemSelectedEventHandlerIUnknownQueryInterface(this *ICoreWebView2CustomItemSelectedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2CustomItemSelectedEventHandlerIUnknownAddRef(this *ICoreWebView2CustomItemSelectedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2CustomItemSelectedEventHandlerIUnknownRelease(this *ICoreWebView2CustomItemSelectedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2CustomItemSelectedEventHandlerInvoke(this *ICoreWebView2CustomItemSelectedEventHandler, sender *ICoreWebView2ContextMenuItem, args *IUnknown) uintptr {
	return this.impl.CustomItemSelected(sender, args)
}

type ICoreWebView2CustomItemSelectedEventHandlerImpl interface {
	IUnknownImpl
	CustomItemSelected(sender *ICoreWebView2ContextMenuItem, args *IUnknown) uintptr
}

var ICoreWebView2CustomItemSelectedEventHandlerFn = ICoreWebView2CustomItemSelectedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2CustomItemSelectedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2CustomItemSelectedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2CustomItemSelectedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2CustomItemSelectedEventHandlerInvoke),
}

func NewICoreWebView2CustomItemSelectedEventHandler(impl ICoreWebView2CustomItemSelectedEventHandlerImpl) *ICoreWebView2CustomItemSelectedEventHandler {
	return &ICoreWebView2CustomItemSelectedEventHandler{
		Vtbl: &ICoreWebView2CustomItemSelectedEventHandlerFn,
		impl: impl,
	}
}
