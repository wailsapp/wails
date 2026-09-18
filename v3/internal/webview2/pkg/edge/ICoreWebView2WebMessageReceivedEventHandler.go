//go:build windows

package edge

type iCoreWebView2WebMessageReceivedEventHandlerVtbl struct {
	_IUnknownVtbl
	Invoke ComProc
}

type iCoreWebView2WebMessageReceivedEventHandler struct {
	vtbl *iCoreWebView2WebMessageReceivedEventHandlerVtbl
	impl iCoreWebView2WebMessageReceivedEventHandlerImpl
}

func _ICoreWebView2WebMessageReceivedEventHandlerIUnknownQueryInterface(this *iCoreWebView2WebMessageReceivedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func _ICoreWebView2WebMessageReceivedEventHandlerIUnknownAddRef(this *iCoreWebView2WebMessageReceivedEventHandler) uintptr {
	return this.impl.AddRef()
}

func _ICoreWebView2WebMessageReceivedEventHandlerIUnknownRelease(this *iCoreWebView2WebMessageReceivedEventHandler) uintptr {
	return this.impl.Release()
}

func _ICoreWebView2WebMessageReceivedEventHandlerInvoke(this *iCoreWebView2WebMessageReceivedEventHandler, sender *ICoreWebView2, args *ICoreWebView2WebMessageReceivedEventArgs) uintptr {
	return this.impl.MessageReceived(sender, args)
}

type iCoreWebView2WebMessageReceivedEventHandlerImpl interface {
	_IUnknownImpl
	MessageReceived(sender *ICoreWebView2, args *ICoreWebView2WebMessageReceivedEventArgs) uintptr
}

var iCoreWebView2WebMessageReceivedEventHandlerFn = iCoreWebView2WebMessageReceivedEventHandlerVtbl{
	_IUnknownVtbl{
		NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerIUnknownQueryInterface),
		NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerIUnknownAddRef),
		NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerIUnknownRelease),
	},
	NewComProc(_ICoreWebView2WebMessageReceivedEventHandlerInvoke),
}

func newICoreWebView2WebMessageReceivedEventHandler(impl iCoreWebView2WebMessageReceivedEventHandlerImpl) *iCoreWebView2WebMessageReceivedEventHandler {
	return &iCoreWebView2WebMessageReceivedEventHandler{
		vtbl: &iCoreWebView2WebMessageReceivedEventHandlerFn,
		impl: impl,
	}
}
