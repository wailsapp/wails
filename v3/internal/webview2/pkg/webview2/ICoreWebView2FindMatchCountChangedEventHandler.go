//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2FindMatchCountChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FindMatchCountChangedEventHandler struct {
	Vtbl *ICoreWebView2FindMatchCountChangedEventHandlerVtbl
	impl ICoreWebView2FindMatchCountChangedEventHandlerImpl
}

func (i *ICoreWebView2FindMatchCountChangedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2FindMatchCountChangedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2FindMatchCountChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2FindMatchCountChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FindMatchCountChangedEventHandlerIUnknownAddRef(this *ICoreWebView2FindMatchCountChangedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2FindMatchCountChangedEventHandlerIUnknownRelease(this *ICoreWebView2FindMatchCountChangedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2FindMatchCountChangedEventHandlerInvoke(this *ICoreWebView2FindMatchCountChangedEventHandler, sender *ICoreWebView2Find, args *IUnknown) uintptr {
	return this.impl.FindMatchCountChanged(sender, args)
}

type ICoreWebView2FindMatchCountChangedEventHandlerImpl interface {
	IUnknownImpl
	FindMatchCountChanged(sender *ICoreWebView2Find, args *IUnknown) uintptr
}

var ICoreWebView2FindMatchCountChangedEventHandlerFn = ICoreWebView2FindMatchCountChangedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2FindMatchCountChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FindMatchCountChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FindMatchCountChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FindMatchCountChangedEventHandlerInvoke),
}

func NewICoreWebView2FindMatchCountChangedEventHandler(impl ICoreWebView2FindMatchCountChangedEventHandlerImpl) *ICoreWebView2FindMatchCountChangedEventHandler {
	return &ICoreWebView2FindMatchCountChangedEventHandler{
		Vtbl: &ICoreWebView2FindMatchCountChangedEventHandlerFn,
		impl: impl,
	}
}
