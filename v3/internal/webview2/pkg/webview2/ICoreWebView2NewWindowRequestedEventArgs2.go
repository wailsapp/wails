//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2NewWindowRequestedEventArgs2Vtbl struct {
	IUnknownVtbl
	GetName ComProc
}

type ICoreWebView2NewWindowRequestedEventArgs2 struct {
	Vtbl *ICoreWebView2NewWindowRequestedEventArgs2Vtbl
}

func (i *ICoreWebView2NewWindowRequestedEventArgs2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2NewWindowRequestedEventArgs2() *ICoreWebView2NewWindowRequestedEventArgs2 {
	var result *ICoreWebView2NewWindowRequestedEventArgs2

	iidICoreWebView2NewWindowRequestedEventArgs2 := NewGUID("{bbc7baed-74c6-4c92-b63a-7f5aeae03de3}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2NewWindowRequestedEventArgs2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2NewWindowRequestedEventArgs2) GetName() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetName.Call(
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
