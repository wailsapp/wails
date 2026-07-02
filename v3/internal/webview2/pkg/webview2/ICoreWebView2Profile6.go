//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Profile6Vtbl struct {
	IUnknownVtbl
	GetIsPasswordAutosaveEnabled ComProc
	PutIsPasswordAutosaveEnabled ComProc
	GetIsGeneralAutofillEnabled  ComProc
	PutIsGeneralAutofillEnabled  ComProc
}

type ICoreWebView2Profile6 struct {
	Vtbl *ICoreWebView2Profile6Vtbl
}

func (i *ICoreWebView2Profile6) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Profile6() *ICoreWebView2Profile6 {
	var result *ICoreWebView2Profile6

	iidICoreWebView2Profile6 := NewGUID("{BD82FA6A-1D65-4C33-B2B4-0393020CC61B}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Profile6)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Profile6) GetIsPasswordAutosaveEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsPasswordAutosaveEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	value := _value != 0
	return value, nil
}

func (i *ICoreWebView2Profile6) PutIsPasswordAutosaveEnabled(value bool) error {

	hr, _, _ := i.Vtbl.PutIsPasswordAutosaveEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Profile6) GetIsGeneralAutofillEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsGeneralAutofillEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	value := _value != 0
	return value, nil
}

func (i *ICoreWebView2Profile6) PutIsGeneralAutofillEnabled(value bool) error {

	hr, _, _ := i.Vtbl.PutIsGeneralAutofillEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
