//go:build windows

package edge

type iCoreWebView2CursorChangedEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2CursorChangedEventHandler struct {
	vtbl *iCoreWebView2CursorChangedEventHandlerVtbl
	impl iCoreWebView2CursorChangedEventHandlerImpl
}

func iCoreWebView2CursorChangedEventHandlerIUnknownQueryInterface(this *iCoreWebView2CursorChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func iCoreWebView2CursorChangedEventHandlerIUnknownAddRef(this *iCoreWebView2CursorChangedEventHandler) uintptr {
	return this.impl.AddRef()
}

func iCoreWebView2CursorChangedEventHandlerIUnknownRelease(this *iCoreWebView2CursorChangedEventHandler) uintptr {
	return this.impl.Release()
}

func iCoreWebView2CursorChangedEventHandlerInvoke(this *iCoreWebView2CursorChangedEventHandler, sender *ICoreWebView2CompositionController, args *IUnknown) uintptr {
	return this.impl.CursorChanged(sender, args)
}

type iCoreWebView2CursorChangedEventHandlerImpl interface {
	_IUnknownImpl
	CursorChanged(sender *ICoreWebView2CompositionController, args *IUnknown) uintptr
}

var iCoreWebView2CursorChangedEventHandlerFn = iCoreWebView2CursorChangedEventHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(iCoreWebView2CursorChangedEventHandlerIUnknownQueryInterface),
		NewComProc(iCoreWebView2CursorChangedEventHandlerIUnknownAddRef),
		NewComProc(iCoreWebView2CursorChangedEventHandlerIUnknownRelease),
	},
	NewComProc(iCoreWebView2CursorChangedEventHandlerInvoke),
}

func newICoreWebView2CursorChangedEventHandler(impl iCoreWebView2CursorChangedEventHandlerImpl) *iCoreWebView2CursorChangedEventHandler {
	return &iCoreWebView2CursorChangedEventHandler{
		vtbl: &iCoreWebView2CursorChangedEventHandlerFn,
		impl: impl,
	}
}
