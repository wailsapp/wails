//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_12Vtbl struct {
	IUnknownVtbl
	AddStatusBarTextChanged    ComProc
	RemoveStatusBarTextChanged ComProc
	GetStatusBarText           ComProc
}

type ICoreWebView2_12 struct {
	Vtbl *ICoreWebView2_12Vtbl
}

func (i *ICoreWebView2_12) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_12() *ICoreWebView2_12 {
	var result *ICoreWebView2_12

	iidICoreWebView2_12 := NewGUID("{35D69927-BCFA-4566-9349-6B3E0D154CAC}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_12)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_12) AddStatusBarTextChanged(eventHandler *ICoreWebView2StatusBarTextChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddStatusBarTextChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_12) RemoveStatusBarTextChanged(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveStatusBarTextChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_12) GetStatusBarText() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetStatusBarText.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}
