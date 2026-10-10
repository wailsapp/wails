//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2WebResourceResponseReceivedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2WebResourceResponseReceivedEventHandler struct {
	Vtbl *ICoreWebView2WebResourceResponseReceivedEventHandlerVtbl
	impl ICoreWebView2WebResourceResponseReceivedEventHandlerImpl
}

func (i *ICoreWebView2WebResourceResponseReceivedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2WebResourceResponseReceivedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2WebResourceResponseReceivedEventHandlerIUnknownQueryInterface(this *ICoreWebView2WebResourceResponseReceivedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2WebResourceResponseReceivedEventHandlerIUnknownAddRef(this *ICoreWebView2WebResourceResponseReceivedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2WebResourceResponseReceivedEventHandlerIUnknownRelease(this *ICoreWebView2WebResourceResponseReceivedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2WebResourceResponseReceivedEventHandlerInvoke(this *ICoreWebView2WebResourceResponseReceivedEventHandler, sender *ICoreWebView2, args *ICoreWebView2WebResourceResponseReceivedEventArgs) uintptr {
	return this.impl.WebResourceResponseReceived(sender, args)
}

type ICoreWebView2WebResourceResponseReceivedEventHandlerImpl interface {
	IUnknownImpl
	WebResourceResponseReceived(sender *ICoreWebView2, args *ICoreWebView2WebResourceResponseReceivedEventArgs) uintptr
}

var ICoreWebView2WebResourceResponseReceivedEventHandlerFn = ICoreWebView2WebResourceResponseReceivedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2WebResourceResponseReceivedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2WebResourceResponseReceivedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2WebResourceResponseReceivedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2WebResourceResponseReceivedEventHandlerInvoke),
}

func NewICoreWebView2WebResourceResponseReceivedEventHandler(impl ICoreWebView2WebResourceResponseReceivedEventHandlerImpl) *ICoreWebView2WebResourceResponseReceivedEventHandler {
	return &ICoreWebView2WebResourceResponseReceivedEventHandler{
		Vtbl: &ICoreWebView2WebResourceResponseReceivedEventHandlerFn,
		impl: impl,
	}
}
