//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_27Vtbl struct {
	IUnknownVtbl
	AddScreenCaptureStarting    ComProc
	RemoveScreenCaptureStarting ComProc
}

type ICoreWebView2_27 struct {
	Vtbl *ICoreWebView2_27Vtbl
}

func (i *ICoreWebView2_27) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_27() *ICoreWebView2_27 {
	var result *ICoreWebView2_27

	iidICoreWebView2_27 := NewGUID("{00fbe33b-8c07-517c-aa23-0ddd4b5f6fa0}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_27)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_27) AddScreenCaptureStarting(eventHandler *ICoreWebView2ScreenCaptureStartingEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddScreenCaptureStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_27) RemoveScreenCaptureStarting(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveScreenCaptureStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
