//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Frame6Vtbl struct {
	IUnknownVtbl
	AddScreenCaptureStarting    ComProc
	RemoveScreenCaptureStarting ComProc
}

type ICoreWebView2Frame6 struct {
	Vtbl *ICoreWebView2Frame6Vtbl
}

func (i *ICoreWebView2Frame6) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Frame6() *ICoreWebView2Frame6 {
	var result *ICoreWebView2Frame6

	iidICoreWebView2Frame6 := NewGUID("{0de611fd-31e9-5ddc-9d71-95eda26eff32}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Frame6)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Frame6) AddScreenCaptureStarting(eventHandler *ICoreWebView2FrameScreenCaptureStartingEventHandler) (EventRegistrationToken, error) {

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

func (i *ICoreWebView2Frame6) RemoveScreenCaptureStarting(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveScreenCaptureStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
