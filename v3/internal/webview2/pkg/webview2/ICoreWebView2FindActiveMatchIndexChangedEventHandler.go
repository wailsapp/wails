//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2FindActiveMatchIndexChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FindActiveMatchIndexChangedEventHandler struct {
	Vtbl *ICoreWebView2FindActiveMatchIndexChangedEventHandlerVtbl
	impl ICoreWebView2FindActiveMatchIndexChangedEventHandlerImpl
}

func (i *ICoreWebView2FindActiveMatchIndexChangedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2FindActiveMatchIndexChangedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2FindActiveMatchIndexChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2FindActiveMatchIndexChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FindActiveMatchIndexChangedEventHandlerIUnknownAddRef(this *ICoreWebView2FindActiveMatchIndexChangedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2FindActiveMatchIndexChangedEventHandlerIUnknownRelease(this *ICoreWebView2FindActiveMatchIndexChangedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2FindActiveMatchIndexChangedEventHandlerInvoke(this *ICoreWebView2FindActiveMatchIndexChangedEventHandler, sender *ICoreWebView2Find, args *IUnknown) uintptr {
	return this.impl.FindActiveMatchIndexChanged(sender, args)
}

type ICoreWebView2FindActiveMatchIndexChangedEventHandlerImpl interface {
	IUnknownImpl
	FindActiveMatchIndexChanged(sender *ICoreWebView2Find, args *IUnknown) uintptr
}

var ICoreWebView2FindActiveMatchIndexChangedEventHandlerFn = ICoreWebView2FindActiveMatchIndexChangedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2FindActiveMatchIndexChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FindActiveMatchIndexChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FindActiveMatchIndexChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FindActiveMatchIndexChangedEventHandlerInvoke),
}

func NewICoreWebView2FindActiveMatchIndexChangedEventHandler(impl ICoreWebView2FindActiveMatchIndexChangedEventHandlerImpl) *ICoreWebView2FindActiveMatchIndexChangedEventHandler {
	return &ICoreWebView2FindActiveMatchIndexChangedEventHandler{
		Vtbl: &ICoreWebView2FindActiveMatchIndexChangedEventHandlerFn,
		impl: impl,
	}
}
