//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2FramePermissionRequestedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2FramePermissionRequestedEventHandler struct {
	Vtbl *ICoreWebView2FramePermissionRequestedEventHandlerVtbl
	impl ICoreWebView2FramePermissionRequestedEventHandlerImpl
}

func (i *ICoreWebView2FramePermissionRequestedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2FramePermissionRequestedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2FramePermissionRequestedEventHandlerIUnknownQueryInterface(this *ICoreWebView2FramePermissionRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2FramePermissionRequestedEventHandlerIUnknownAddRef(this *ICoreWebView2FramePermissionRequestedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2FramePermissionRequestedEventHandlerIUnknownRelease(this *ICoreWebView2FramePermissionRequestedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2FramePermissionRequestedEventHandlerInvoke(this *ICoreWebView2FramePermissionRequestedEventHandler, sender *ICoreWebView2Frame, args *ICoreWebView2PermissionRequestedEventArgs2) uintptr {
	return this.impl.FramePermissionRequested(sender, args)
}

type ICoreWebView2FramePermissionRequestedEventHandlerImpl interface {
	IUnknownImpl
	FramePermissionRequested(sender *ICoreWebView2Frame, args *ICoreWebView2PermissionRequestedEventArgs2) uintptr
}

var ICoreWebView2FramePermissionRequestedEventHandlerFn = ICoreWebView2FramePermissionRequestedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2FramePermissionRequestedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2FramePermissionRequestedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2FramePermissionRequestedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2FramePermissionRequestedEventHandlerInvoke),
}

func NewICoreWebView2FramePermissionRequestedEventHandler(impl ICoreWebView2FramePermissionRequestedEventHandlerImpl) *ICoreWebView2FramePermissionRequestedEventHandler {
	return &ICoreWebView2FramePermissionRequestedEventHandler{
		Vtbl: &ICoreWebView2FramePermissionRequestedEventHandlerFn,
		impl: impl,
	}
}
