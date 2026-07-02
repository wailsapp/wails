//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2AcceleratorKeyPressedEventArgs2Vtbl struct {
	IUnknownVtbl
	GetIsBrowserAcceleratorKeyEnabled ComProc
	PutIsBrowserAcceleratorKeyEnabled ComProc
}

type ICoreWebView2AcceleratorKeyPressedEventArgs2 struct {
	Vtbl *ICoreWebView2AcceleratorKeyPressedEventArgs2Vtbl
}

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2AcceleratorKeyPressedEventArgs2() *ICoreWebView2AcceleratorKeyPressedEventArgs2 {
	var result *ICoreWebView2AcceleratorKeyPressedEventArgs2

	iidICoreWebView2AcceleratorKeyPressedEventArgs2 := NewGUID("{03b2c8c8-7799-4e34-bd66-ed26aa85f2bf}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2AcceleratorKeyPressedEventArgs2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs2) GetIsBrowserAcceleratorKeyEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsBrowserAcceleratorKeyEnabled.Call(
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

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs2) PutIsBrowserAcceleratorKeyEnabled(value bool) error {

	hr, _, _ := i.Vtbl.PutIsBrowserAcceleratorKeyEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
