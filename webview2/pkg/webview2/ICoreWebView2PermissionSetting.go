//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2PermissionSettingVtbl struct {
	IUnknownVtbl
	GetPermissionKind   ComProc
	GetPermissionOrigin ComProc
	GetPermissionState  ComProc
}

type ICoreWebView2PermissionSetting struct {
	Vtbl *ICoreWebView2PermissionSettingVtbl
}

func (i *ICoreWebView2PermissionSetting) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2PermissionSetting) GetPermissionKind() (COREWEBVIEW2_PERMISSION_KIND, error) {

	var value COREWEBVIEW2_PERMISSION_KIND

	hr, _, _ := i.Vtbl.GetPermissionKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2PermissionSetting) GetPermissionOrigin() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetPermissionOrigin.Call(
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

func (i *ICoreWebView2PermissionSetting) GetPermissionState() (COREWEBVIEW2_PERMISSION_STATE, error) {

	var value COREWEBVIEW2_PERMISSION_STATE

	hr, _, _ := i.Vtbl.GetPermissionState.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}
