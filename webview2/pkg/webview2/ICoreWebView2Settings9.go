//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Settings9Vtbl struct {
	IUnknownVtbl
	GetIsNonClientRegionSupportEnabled ComProc
	PutIsNonClientRegionSupportEnabled ComProc
}

type ICoreWebView2Settings9 struct {
	Vtbl *ICoreWebView2Settings9Vtbl
}

func (i *ICoreWebView2Settings9) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Settings9() *ICoreWebView2Settings9 {
	var result *ICoreWebView2Settings9

	iidICoreWebView2Settings9 := NewGUID("{0528a73b-e92d-49f4-927a-e547dddaa37d}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Settings9)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Settings9) GetIsNonClientRegionSupportEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsNonClientRegionSupportEnabled.Call(
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

func (i *ICoreWebView2Settings9) PutIsNonClientRegionSupportEnabled(value bool) error {

	hr, _, _ := i.Vtbl.PutIsNonClientRegionSupportEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
