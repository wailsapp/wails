//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2CompositionController5Vtbl struct {
	IUnknownVtbl
	AddDragStarting ComProc
	RemoveDragStarting ComProc
}

type ICoreWebView2CompositionController5 struct {
	Vtbl *ICoreWebView2CompositionController5Vtbl
}

func (i *ICoreWebView2CompositionController5) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2CompositionController5) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2CompositionController5() (*ICoreWebView2CompositionController5, error) {
	var result *ICoreWebView2CompositionController5

	iidICoreWebView2CompositionController5 := NewGUID("{8d0f82eb-7c33-5a4c-9108-84ca28ccc3b4}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2CompositionController5)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2CompositionController5) AddDragStarting(eventHandler *ICoreWebView2DragStartingEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddDragStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2CompositionController5) RemoveDragStarting(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveDragStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
