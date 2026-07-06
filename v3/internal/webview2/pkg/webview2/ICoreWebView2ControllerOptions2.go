//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2ControllerOptions2Vtbl struct {
	IUnknownVtbl
	GetScriptLocale ComProc
	PutScriptLocale ComProc
}

type ICoreWebView2ControllerOptions2 struct {
	Vtbl *ICoreWebView2ControllerOptions2Vtbl
}

func (i *ICoreWebView2ControllerOptions2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2ControllerOptions2() *ICoreWebView2ControllerOptions2 {
	var result *ICoreWebView2ControllerOptions2

	iidICoreWebView2ControllerOptions2 := NewGUID("{06c991d8-9e7e-11ed-a8fc-0242ac120002}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2ControllerOptions2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2ControllerOptions2) GetScriptLocale() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetScriptLocale.Call(
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

func (i *ICoreWebView2ControllerOptions2) PutScriptLocale(value string) error {

	// Convert string 'value' to *uint16
	_value, err := UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutScriptLocale.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
