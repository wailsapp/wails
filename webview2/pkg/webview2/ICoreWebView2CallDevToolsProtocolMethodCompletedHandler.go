//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2CallDevToolsProtocolMethodCompletedHandler struct {
	Vtbl *ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerVtbl
	impl ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerImpl
}

func (i *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownAddRef(this *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownRelease(this *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerInvoke(this *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler, errorCode uintptr, result string) uintptr {
	return this.impl.CallDevToolsProtocolMethodCompleted(errorCode, result)
}

type ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerImpl interface {
	IUnknownImpl
	CallDevToolsProtocolMethodCompleted(errorCode uintptr, result string) uintptr
}

var ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerFn = ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerInvoke),
}

func NewICoreWebView2CallDevToolsProtocolMethodCompletedHandler(impl ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerImpl) *ICoreWebView2CallDevToolsProtocolMethodCompletedHandler {
	return &ICoreWebView2CallDevToolsProtocolMethodCompletedHandler{
		Vtbl: &ICoreWebView2CallDevToolsProtocolMethodCompletedHandlerFn,
		impl: impl,
	}
}
