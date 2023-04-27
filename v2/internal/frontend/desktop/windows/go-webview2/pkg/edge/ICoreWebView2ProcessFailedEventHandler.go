//go:build windows

package edge

type _ICoreWebView2ProcessFailedEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2ProcessFailedEventHandler struct {
	vtbl *_ICoreWebView2ProcessFailedEventHandlerVtbl
	impl _ICoreWebView2ProcessFailedEventHandlerImpl
}

func (i *ICoreWebView2ProcessFailedEventHandler) AddRef() uintptr {
	return i.AddRef()
}
func _ICoreWebView2ProcessFailedEventHandlerIUnknownQueryInterface(this *ICoreWebView2ProcessFailedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2ProcessFailedEventHandlerIUnknownAddRef(this *ICoreWebView2ProcessFailedEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2ProcessFailedEventHandlerIUnknownRelease(this *ICoreWebView2ProcessFailedEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2ProcessFailedEventHandlerInvoke(this *ICoreWebView2ProcessFailedEventHandler, sender *ICoreWebView2, args *ICoreWebView2ProcessFailedEventArgs) uintptr {
	return this.impl.ProcessFailed(sender, args)
}

type _ICoreWebView2ProcessFailedEventHandlerImpl interface {
	_IUnknownImpl
	ProcessFailed(sender *ICoreWebView2, args *ICoreWebView2ProcessFailedEventArgs) uintptr
}

var _ICoreWebView2ProcessFailedEventHandlerFn = _ICoreWebView2ProcessFailedEventHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2ProcessFailedEventHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2ProcessFailedEventHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2ProcessFailedEventHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2ProcessFailedEventHandlerInvoke),
}

func newICoreWebView2ProcessFailedEventHandler(impl _ICoreWebView2ProcessFailedEventHandlerImpl) *ICoreWebView2ProcessFailedEventHandler {
	return &ICoreWebView2ProcessFailedEventHandler{
		vtbl: &_ICoreWebView2ProcessFailedEventHandlerFn,
		impl: impl,
	}
}
