//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Frame3Vtbl struct {
	IUnknownVtbl
	AddPermissionRequested    ComProc
	RemovePermissionRequested ComProc
}

type ICoreWebView2Frame3 struct {
	Vtbl *ICoreWebView2Frame3Vtbl
}

func (i *ICoreWebView2Frame3) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Frame3() *ICoreWebView2Frame3 {
	var result *ICoreWebView2Frame3

	iidICoreWebView2Frame3 := NewGUID("{b50d82cc-cc28-481d-9614-cb048895e6a0}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Frame3)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Frame3) AddPermissionRequested(eventHandler *ICoreWebView2FramePermissionRequestedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddPermissionRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2Frame3) RemovePermissionRequested(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemovePermissionRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
