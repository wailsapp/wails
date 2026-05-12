//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2EnvironmentOptions2Vtbl struct {
	IUnknownVtbl
	GetExclusiveUserDataFolderAccess ComProc
	PutExclusiveUserDataFolderAccess ComProc
}

type ICoreWebView2EnvironmentOptions2 struct {
	Vtbl *ICoreWebView2EnvironmentOptions2Vtbl
}

func (i *ICoreWebView2EnvironmentOptions2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2EnvironmentOptions2) GetExclusiveUserDataFolderAccess() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetExclusiveUserDataFolderAccess.Call(
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

func (i *ICoreWebView2EnvironmentOptions2) PutExclusiveUserDataFolderAccess(value bool) error {

	hr, _, _ := i.Vtbl.PutExclusiveUserDataFolderAccess.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
