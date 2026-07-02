//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2PermissionRequestedEventArgs2Vtbl struct {
	IUnknownVtbl
	GetHandled ComProc
	PutHandled ComProc
}

type ICoreWebView2PermissionRequestedEventArgs2 struct {
	Vtbl *ICoreWebView2PermissionRequestedEventArgs2Vtbl
}

func (i *ICoreWebView2PermissionRequestedEventArgs2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2PermissionRequestedEventArgs2() *ICoreWebView2PermissionRequestedEventArgs2 {
	var result *ICoreWebView2PermissionRequestedEventArgs2

	iidICoreWebView2PermissionRequestedEventArgs2 := NewGUID("{74d7127f-9de6-4200-8734-42d6fb4ff741}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2PermissionRequestedEventArgs2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2PermissionRequestedEventArgs2) GetHandled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetHandled.Call(
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

func (i *ICoreWebView2PermissionRequestedEventArgs2) PutHandled(value bool) error {

	hr, _, _ := i.Vtbl.PutHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
