//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2AcceleratorKeyPressedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2AcceleratorKeyPressedEventHandler struct {
	Vtbl *ICoreWebView2AcceleratorKeyPressedEventHandlerVtbl
	impl ICoreWebView2AcceleratorKeyPressedEventHandlerImpl
}

func (i *ICoreWebView2AcceleratorKeyPressedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2AcceleratorKeyPressedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownQueryInterface(this *ICoreWebView2AcceleratorKeyPressedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownAddRef(this *ICoreWebView2AcceleratorKeyPressedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownRelease(this *ICoreWebView2AcceleratorKeyPressedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2AcceleratorKeyPressedEventHandlerInvoke(this *ICoreWebView2AcceleratorKeyPressedEventHandler, sender *ICoreWebView2Controller, args *ICoreWebView2AcceleratorKeyPressedEventArgs) uintptr {
	return this.impl.AcceleratorKeyPressed(sender, args)
}

type ICoreWebView2AcceleratorKeyPressedEventHandlerImpl interface {
	IUnknownImpl
	AcceleratorKeyPressed(sender *ICoreWebView2Controller, args *ICoreWebView2AcceleratorKeyPressedEventArgs) uintptr
}

var ICoreWebView2AcceleratorKeyPressedEventHandlerFn = ICoreWebView2AcceleratorKeyPressedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2AcceleratorKeyPressedEventHandlerInvoke),
}

func NewICoreWebView2AcceleratorKeyPressedEventHandler(impl ICoreWebView2AcceleratorKeyPressedEventHandlerImpl) *ICoreWebView2AcceleratorKeyPressedEventHandler {
	return &ICoreWebView2AcceleratorKeyPressedEventHandler{
		Vtbl: &ICoreWebView2AcceleratorKeyPressedEventHandlerFn,
		impl: impl,
	}
}
