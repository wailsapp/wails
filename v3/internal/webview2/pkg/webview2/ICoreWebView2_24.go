//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_24Vtbl struct {
	IUnknownVtbl
	AddNotificationReceived    ComProc
	RemoveNotificationReceived ComProc
}

type ICoreWebView2_24 struct {
	Vtbl *ICoreWebView2_24Vtbl
}

func (i *ICoreWebView2_24) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_24() *ICoreWebView2_24 {
	var result *ICoreWebView2_24

	iidICoreWebView2_24 := NewGUID("{39a7ad55-4287-5cc1-88a1-c6f458593824}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_24)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_24) AddNotificationReceived(eventHandler *ICoreWebView2NotificationReceivedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddNotificationReceived.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_24) RemoveNotificationReceived(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveNotificationReceived.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
