//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_27Vtbl struct {
	IUnknownVtbl
	AddScreenCaptureStarting ComProc
	RemoveScreenCaptureStarting ComProc
}

type ICoreWebView2_27 struct {
	Vtbl *ICoreWebView2_27Vtbl
}

func (i *ICoreWebView2_27) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_27) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2_27() (*ICoreWebView2_27, error) {
	var result *ICoreWebView2_27

	iidICoreWebView2_27 := NewGUID("{00fbe33b-8c07-517c-aa23-0ddd4b5f6fa0}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_27)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2_27) AddScreenCaptureStarting(eventHandler *ICoreWebView2ScreenCaptureStartingEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddScreenCaptureStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2_27) RemoveScreenCaptureStarting(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveScreenCaptureStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
