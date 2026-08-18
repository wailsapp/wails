//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2DocumentTitleChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2DocumentTitleChangedEventHandler struct {
	Vtbl *ICoreWebView2DocumentTitleChangedEventHandlerVtbl
	impl ICoreWebView2DocumentTitleChangedEventHandlerImpl
}

func (i *ICoreWebView2DocumentTitleChangedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2DocumentTitleChangedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2DocumentTitleChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2DocumentTitleChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2DocumentTitleChangedEventHandlerIUnknownAddRef(this *ICoreWebView2DocumentTitleChangedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2DocumentTitleChangedEventHandlerIUnknownRelease(this *ICoreWebView2DocumentTitleChangedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2DocumentTitleChangedEventHandlerInvoke(this *ICoreWebView2DocumentTitleChangedEventHandler, sender *ICoreWebView2, args *IUnknown) uintptr {
	return this.impl.DocumentTitleChanged(sender, args)
}

type ICoreWebView2DocumentTitleChangedEventHandlerImpl interface {
	IUnknownImpl
	DocumentTitleChanged(sender *ICoreWebView2, args *IUnknown) uintptr
}

var ICoreWebView2DocumentTitleChangedEventHandlerFn = ICoreWebView2DocumentTitleChangedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2DocumentTitleChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2DocumentTitleChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2DocumentTitleChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2DocumentTitleChangedEventHandlerInvoke),
}

func NewICoreWebView2DocumentTitleChangedEventHandler(impl ICoreWebView2DocumentTitleChangedEventHandlerImpl) *ICoreWebView2DocumentTitleChangedEventHandler {
	return &ICoreWebView2DocumentTitleChangedEventHandler{
		Vtbl: &ICoreWebView2DocumentTitleChangedEventHandlerFn,
		impl: impl,
	}
}
