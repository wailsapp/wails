//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2EnvironmentOptionsVtbl struct {
	IUnknownVtbl
	GetAdditionalBrowserArguments             ComProc
	PutAdditionalBrowserArguments             ComProc
	GetLanguage                               ComProc
	PutLanguage                               ComProc
	GetTargetCompatibleBrowserVersion         ComProc
	PutTargetCompatibleBrowserVersion         ComProc
	GetAllowSingleSignOnUsingOSPrimaryAccount ComProc
	PutAllowSingleSignOnUsingOSPrimaryAccount ComProc
}

type ICoreWebView2EnvironmentOptions struct {
	Vtbl *ICoreWebView2EnvironmentOptionsVtbl
}

func (i *ICoreWebView2EnvironmentOptions) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2EnvironmentOptions) GetAdditionalBrowserArguments() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetAdditionalBrowserArguments.Call(
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

func (i *ICoreWebView2EnvironmentOptions) PutAdditionalBrowserArguments(value string) error {

	// Convert string 'value' to *uint16
	_value, err := UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutAdditionalBrowserArguments.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2EnvironmentOptions) GetLanguage() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetLanguage.Call(
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

func (i *ICoreWebView2EnvironmentOptions) PutLanguage(value string) error {

	// Convert string 'value' to *uint16
	_value, err := UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutLanguage.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2EnvironmentOptions) GetTargetCompatibleBrowserVersion() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetTargetCompatibleBrowserVersion.Call(
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

func (i *ICoreWebView2EnvironmentOptions) PutTargetCompatibleBrowserVersion(value string) error {

	// Convert string 'value' to *uint16
	_value, err := UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutTargetCompatibleBrowserVersion.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2EnvironmentOptions) GetAllowSingleSignOnUsingOSPrimaryAccount() (bool, error) {
	// Create int32 to hold bool result
	var _allow int32

	hr, _, _ := i.Vtbl.GetAllowSingleSignOnUsingOSPrimaryAccount.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_allow)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	allow := _allow != 0
	return allow, nil
}

func (i *ICoreWebView2EnvironmentOptions) PutAllowSingleSignOnUsingOSPrimaryAccount(allow bool) error {

	hr, _, _ := i.Vtbl.PutAllowSingleSignOnUsingOSPrimaryAccount.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&allow)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
