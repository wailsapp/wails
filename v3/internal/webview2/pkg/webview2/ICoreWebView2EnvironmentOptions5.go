//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2EnvironmentOptions5Vtbl struct {
	IUnknownVtbl
	GetEnableTrackingPrevention ComProc
	PutEnableTrackingPrevention ComProc
}

type ICoreWebView2EnvironmentOptions5 struct {
	Vtbl *ICoreWebView2EnvironmentOptions5Vtbl
}

func (i *ICoreWebView2EnvironmentOptions5) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2EnvironmentOptions5) GetEnableTrackingPrevention() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetEnableTrackingPrevention.Call(
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

func (i *ICoreWebView2EnvironmentOptions5) PutEnableTrackingPrevention(value bool) error {

	// BOOL is a 4-byte by-value parameter: pass the value, not a pointer
	// to a 1-byte Go bool.
	var _valueInt int32
	if value {
		_valueInt = 1
	}
	hr, _, _ := i.Vtbl.PutEnableTrackingPrevention.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_valueInt),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
