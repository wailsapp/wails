//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2PermissionRequestedEventArgs3Vtbl struct {
	IUnknownVtbl
	GetSavesInProfile ComProc
	PutSavesInProfile ComProc
}

type ICoreWebView2PermissionRequestedEventArgs3 struct {
	Vtbl *ICoreWebView2PermissionRequestedEventArgs3Vtbl
}

func (i *ICoreWebView2PermissionRequestedEventArgs3) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2PermissionRequestedEventArgs3() *ICoreWebView2PermissionRequestedEventArgs3 {
	var result *ICoreWebView2PermissionRequestedEventArgs3

	iidICoreWebView2PermissionRequestedEventArgs3 := NewGUID("{e61670bc-3dce-4177-86d2-c629ae3cb6ac}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2PermissionRequestedEventArgs3)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2PermissionRequestedEventArgs3) GetSavesInProfile() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetSavesInProfile.Call(
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

func (i *ICoreWebView2PermissionRequestedEventArgs3) PutSavesInProfile(value bool) error {

	hr, _, _ := i.Vtbl.PutSavesInProfile.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
