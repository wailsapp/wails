//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_25Vtbl struct {
	IUnknownVtbl
	AddSaveAsUIShowing    ComProc
	RemoveSaveAsUIShowing ComProc
	ShowSaveAsUI          ComProc
}

type ICoreWebView2_25 struct {
	Vtbl *ICoreWebView2_25Vtbl
}

func (i *ICoreWebView2_25) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_25() *ICoreWebView2_25 {
	var result *ICoreWebView2_25

	iidICoreWebView2_25 := NewGUID("{b5a86092-df50-5b4f-a17b-6c8f8b40b771}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_25)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_25) AddSaveAsUIShowing(eventHandler *ICoreWebView2SaveAsUIShowingEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddSaveAsUIShowing.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_25) RemoveSaveAsUIShowing(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveSaveAsUIShowing.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_25) ShowSaveAsUI(handler *ICoreWebView2ShowSaveAsUICompletedHandler) error {

	hr, _, _ := i.Vtbl.ShowSaveAsUI.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
