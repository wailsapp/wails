//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_26Vtbl struct {
	IUnknownVtbl
	AddSaveFileSecurityCheckStarting    ComProc
	RemoveSaveFileSecurityCheckStarting ComProc
}

type ICoreWebView2_26 struct {
	Vtbl *ICoreWebView2_26Vtbl
}

func (i *ICoreWebView2_26) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_26() *ICoreWebView2_26 {
	var result *ICoreWebView2_26

	iidICoreWebView2_26 := NewGUID("{806268b8-f897-5685-88e5-c45fca0b1a48}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_26)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_26) AddSaveFileSecurityCheckStarting(eventHandler *ICoreWebView2SaveFileSecurityCheckStartingEventHandler) (EventRegistrationToken, error) {
	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddSaveFileSecurityCheckStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_26) RemoveSaveFileSecurityCheckStarting(token EventRegistrationToken) error {
	hr, _, _ := i.Vtbl.RemoveSaveFileSecurityCheckStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
