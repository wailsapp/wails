//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Profile8Vtbl struct {
	IUnknownVtbl
	Delete        ComProc
	AddDeleted    ComProc
	RemoveDeleted ComProc
}

type ICoreWebView2Profile8 struct {
	Vtbl *ICoreWebView2Profile8Vtbl
}

func (i *ICoreWebView2Profile8) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Profile8() *ICoreWebView2Profile8 {
	var result *ICoreWebView2Profile8

	iidICoreWebView2Profile8 := NewGUID("{fbf70c2f-eb1f-4383-85a0-163e92044011}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Profile8)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Profile8) Delete() error {

	hr, _, _ := i.Vtbl.Delete.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Profile8) AddDeleted(eventHandler *ICoreWebView2ProfileDeletedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddDeleted.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2Profile8) RemoveDeleted(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveDeleted.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
