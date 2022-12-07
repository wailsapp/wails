//go:build windows

package edge

type _ICoreWebView2AcceleratorKeyPressedEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2AcceleratorKeyPressedEventHandler struct {
	vtbl *_ICoreWebView2AcceleratorKeyPressedEventHandlerVtbl
	impl _ICoreWebView2AcceleratorKeyPressedEventHandlerImpl
}

func (i *ICoreWebView2AcceleratorKeyPressedEventHandler) AddRef() uintptr {
	return i.AddRef()
}
func _ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownQueryInterface(this *ICoreWebView2AcceleratorKeyPressedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownAddRef(this *ICoreWebView2AcceleratorKeyPressedEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownRelease(this *ICoreWebView2AcceleratorKeyPressedEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2AcceleratorKeyPressedEventHandlerInvoke(this *ICoreWebView2AcceleratorKeyPressedEventHandler, sender *ICoreWebView2Controller, args *ICoreWebView2AcceleratorKeyPressedEventArgs) uintptr {
	return this.impl.AcceleratorKeyPressed(sender, args)
}

type _ICoreWebView2AcceleratorKeyPressedEventHandlerImpl interface {
	_IUnknownImpl
	AcceleratorKeyPressed(sender *ICoreWebView2Controller, args *ICoreWebView2AcceleratorKeyPressedEventArgs) uintptr
}

var _ICoreWebView2AcceleratorKeyPressedEventHandlerFn = _ICoreWebView2AcceleratorKeyPressedEventHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2AcceleratorKeyPressedEventHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2AcceleratorKeyPressedEventHandlerInvoke),
}

func newICoreWebView2AcceleratorKeyPressedEventHandler(impl _ICoreWebView2AcceleratorKeyPressedEventHandlerImpl) *ICoreWebView2AcceleratorKeyPressedEventHandler {
	return &ICoreWebView2AcceleratorKeyPressedEventHandler{
		vtbl: &_ICoreWebView2AcceleratorKeyPressedEventHandlerFn,
		impl: impl,
	}
}
