//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Environment5Vtbl struct {
	IUnknownVtbl
	AddBrowserProcessExited    ComProc
	RemoveBrowserProcessExited ComProc
}

type ICoreWebView2Environment5 struct {
	Vtbl *ICoreWebView2Environment5Vtbl
}

func (i *ICoreWebView2Environment5) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Environment5() *ICoreWebView2Environment5 {
	var result *ICoreWebView2Environment5

	iidICoreWebView2Environment5 := NewGUID("{319e423d-e0d7-4b8d-9254-ae9475de9b17}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment5)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Environment5) AddBrowserProcessExited(eventHandler *ICoreWebView2BrowserProcessExitedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddBrowserProcessExited.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2Environment5) RemoveBrowserProcessExited(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveBrowserProcessExited.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
