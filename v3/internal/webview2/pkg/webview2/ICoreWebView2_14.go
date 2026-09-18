//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_14Vtbl struct {
	IUnknownVtbl
	AddServerCertificateErrorDetected    ComProc
	RemoveServerCertificateErrorDetected ComProc
	ClearServerCertificateErrorActions   ComProc
}

type ICoreWebView2_14 struct {
	Vtbl *ICoreWebView2_14Vtbl
}

func (i *ICoreWebView2_14) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_14() *ICoreWebView2_14 {
	var result *ICoreWebView2_14

	iidICoreWebView2_14 := NewGUID("{6daa4f10-4a90-4753-8898-77c5df534165}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_14)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_14) AddServerCertificateErrorDetected(eventHandler *ICoreWebView2ServerCertificateErrorDetectedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddServerCertificateErrorDetected.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_14) RemoveServerCertificateErrorDetected(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveServerCertificateErrorDetected.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_14) ClearServerCertificateErrorActions(handler *ICoreWebView2ClearServerCertificateErrorActionsCompletedHandler) error {

	hr, _, _ := i.Vtbl.ClearServerCertificateErrorActions.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
