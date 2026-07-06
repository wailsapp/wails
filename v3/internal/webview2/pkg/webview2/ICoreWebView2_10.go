//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_10Vtbl struct {
	IUnknownVtbl
	AddBasicAuthenticationRequested    ComProc
	RemoveBasicAuthenticationRequested ComProc
}

type ICoreWebView2_10 struct {
	Vtbl *ICoreWebView2_10Vtbl
}

func (i *ICoreWebView2_10) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_10() *ICoreWebView2_10 {
	var result *ICoreWebView2_10

	iidICoreWebView2_10 := NewGUID("{b1690564-6f5a-4983-8e48-31d1143fecdb}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_10)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_10) AddBasicAuthenticationRequested(eventHandler *ICoreWebView2BasicAuthenticationRequestedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddBasicAuthenticationRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_10) RemoveBasicAuthenticationRequested(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveBasicAuthenticationRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
