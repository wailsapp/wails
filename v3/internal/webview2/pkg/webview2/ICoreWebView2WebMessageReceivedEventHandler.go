//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2WebMessageReceivedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2WebMessageReceivedEventHandler struct {
	Vtbl *ICoreWebView2WebMessageReceivedEventHandlerVtbl
	impl ICoreWebView2WebMessageReceivedEventHandlerImpl
}

func (i *ICoreWebView2WebMessageReceivedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2WebMessageReceivedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2WebMessageReceivedEventHandlerIUnknownQueryInterface(this *ICoreWebView2WebMessageReceivedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2WebMessageReceivedEventHandlerIUnknownAddRef(this *ICoreWebView2WebMessageReceivedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2WebMessageReceivedEventHandlerIUnknownRelease(this *ICoreWebView2WebMessageReceivedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2WebMessageReceivedEventHandlerInvoke(this *ICoreWebView2WebMessageReceivedEventHandler, sender *ICoreWebView2, args *ICoreWebView2WebMessageReceivedEventArgs) uintptr {
	return this.impl.WebMessageReceived(sender, args)
}

type ICoreWebView2WebMessageReceivedEventHandlerImpl interface {
	IUnknownImpl
	WebMessageReceived(sender *ICoreWebView2, args *ICoreWebView2WebMessageReceivedEventArgs) uintptr
}

var ICoreWebView2WebMessageReceivedEventHandlerFn = ICoreWebView2WebMessageReceivedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2WebMessageReceivedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2WebMessageReceivedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2WebMessageReceivedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2WebMessageReceivedEventHandlerInvoke),
}

func NewICoreWebView2WebMessageReceivedEventHandler(impl ICoreWebView2WebMessageReceivedEventHandlerImpl) *ICoreWebView2WebMessageReceivedEventHandler {
	return &ICoreWebView2WebMessageReceivedEventHandler{
		Vtbl: &ICoreWebView2WebMessageReceivedEventHandlerFn,
		impl: impl,
	}
}
