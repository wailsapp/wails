//go:build windows

package edge

import (
	"unsafe"
)

type _ICoreWebView2MoveFocusRequestedEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2MoveFocusRequestedEventHandler struct {
	vtbl *_ICoreWebView2MoveFocusRequestedEventHandlerVtbl
	impl _ICoreWebView2MoveFocusRequestedEventHandlerImpl
}

func (i *ICoreWebView2MoveFocusRequestedEventHandler) AddRef() uintptr {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return ret
}

func (i *ICoreWebView2MoveFocusRequestedEventHandler) Release() uintptr {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return ret
}

func _ICoreWebView2MoveFocusRequestedEventHandlerIUnknownQueryInterface(this *ICoreWebView2MoveFocusRequestedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2MoveFocusRequestedEventHandlerIUnknownAddRef(this *ICoreWebView2MoveFocusRequestedEventHandler) uintptr {
	return handlerAddRef(unsafe.Pointer(this))
}

func _ICoreWebView2MoveFocusRequestedEventHandlerIUnknownRelease(this *ICoreWebView2MoveFocusRequestedEventHandler) uintptr {
	return handlerRelease(unsafe.Pointer(this))
}

func _ICoreWebView2MoveFocusRequestedEventHandlerInvoke(this *ICoreWebView2MoveFocusRequestedEventHandler, sender *ICoreWebView2Controller, args *ICoreWebView2MoveFocusRequestedEventArgs) uintptr {
	return this.impl.MoveFocusRequested(sender, args)
}

type _ICoreWebView2MoveFocusRequestedEventHandlerImpl interface {
	_IUnknownImpl
	MoveFocusRequested(sender *ICoreWebView2Controller, args *ICoreWebView2MoveFocusRequestedEventArgs) uintptr
}

var _ICoreWebView2MoveFocusRequestedEventHandlerFn = _ICoreWebView2MoveFocusRequestedEventHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2MoveFocusRequestedEventHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2MoveFocusRequestedEventHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2MoveFocusRequestedEventHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2MoveFocusRequestedEventHandlerInvoke),
}

func newICoreWebView2MoveFocusRequestedEventHandler(impl _ICoreWebView2MoveFocusRequestedEventHandlerImpl) *ICoreWebView2MoveFocusRequestedEventHandler {
	return &ICoreWebView2MoveFocusRequestedEventHandler{
		vtbl: &_ICoreWebView2MoveFocusRequestedEventHandlerFn,
		impl: impl,
	}
}
