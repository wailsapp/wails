//go:build windows

package webview2
import (
	"unsafe"
)

type ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2IsDocumentPlayingAudioChangedEventHandler struct {
	Vtbl *ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerVtbl
	impl ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerImpl
}

func (i *ICoreWebView2IsDocumentPlayingAudioChangedEventHandler) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2IsDocumentPlayingAudioChangedEventHandler) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerIUnknownQueryInterface(this *ICoreWebView2IsDocumentPlayingAudioChangedEventHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerIUnknownAddRef(this *ICoreWebView2IsDocumentPlayingAudioChangedEventHandler) uintptr {
	return uintptr(this.impl.AddRef())
}

func ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerIUnknownRelease(this *ICoreWebView2IsDocumentPlayingAudioChangedEventHandler) uintptr {
	return uintptr(this.impl.Release())
}

func ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerInvoke(this *ICoreWebView2IsDocumentPlayingAudioChangedEventHandler, sender *ICoreWebView2, args *IUnknown) uintptr {
	return this.impl.IsDocumentPlayingAudioChanged(sender, args)
}

type ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerImpl interface {
	IUnknownImpl
	IsDocumentPlayingAudioChanged(sender *ICoreWebView2, args *IUnknown) uintptr
}

var ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerFn = ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerVtbl{
	IUnknownVtbl {
		NewComProc(ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerInvoke),
}

func NewICoreWebView2IsDocumentPlayingAudioChangedEventHandler(impl ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerImpl) *ICoreWebView2IsDocumentPlayingAudioChangedEventHandler {
	return &ICoreWebView2IsDocumentPlayingAudioChangedEventHandler{
		Vtbl: &ICoreWebView2IsDocumentPlayingAudioChangedEventHandlerFn,
		impl: impl,
	}
}
