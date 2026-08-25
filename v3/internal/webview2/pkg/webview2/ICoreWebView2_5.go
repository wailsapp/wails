//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_5Vtbl struct {
	IUnknownVtbl
	AddClientCertificateRequested    ComProc
	RemoveClientCertificateRequested ComProc
}

type ICoreWebView2_5 struct {
	Vtbl *ICoreWebView2_5Vtbl
}

func (i *ICoreWebView2_5) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_5() *ICoreWebView2_5 {
	var result *ICoreWebView2_5

	iidICoreWebView2_5 := NewGUID("{bedb11b8-d63c-11eb-b8bc-0242ac130003}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_5)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_5) AddClientCertificateRequested(eventHandler *ICoreWebView2ClientCertificateRequestedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddClientCertificateRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_5) RemoveClientCertificateRequested(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveClientCertificateRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
