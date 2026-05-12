//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2CustomSchemeRegistrationVtbl struct {
	IUnknownVtbl
	GetSchemeName            ComProc
	GetTreatAsSecure         ComProc
	PutTreatAsSecure         ComProc
	GetAllowedOrigins        ComProc
	SetAllowedOrigins        ComProc
	GetHasAuthorityComponent ComProc
	PutHasAuthorityComponent ComProc
}

type ICoreWebView2CustomSchemeRegistration struct {
	Vtbl *ICoreWebView2CustomSchemeRegistrationVtbl
}

func (i *ICoreWebView2CustomSchemeRegistration) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2CustomSchemeRegistration) GetSchemeName() (string, error) {
	// Create *uint16 to hold result
	var _schemeName *uint16

	hr, _, _ := i.Vtbl.GetSchemeName.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_schemeName)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	schemeName := UTF16PtrToString(_schemeName)
	CoTaskMemFree(unsafe.Pointer(_schemeName))
	return schemeName, nil
}

func (i *ICoreWebView2CustomSchemeRegistration) GetTreatAsSecure() (bool, error) {
	// Create int32 to hold bool result
	var _treatAsSecure int32

	hr, _, _ := i.Vtbl.GetTreatAsSecure.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_treatAsSecure)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	treatAsSecure := _treatAsSecure != 0
	return treatAsSecure, nil
}

func (i *ICoreWebView2CustomSchemeRegistration) PutTreatAsSecure(value bool) error {

	hr, _, _ := i.Vtbl.PutTreatAsSecure.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2CustomSchemeRegistration) GetAllowedOrigins() (uint32, *string, error) {

	var allowedOriginsCount uint32
	var allowedOrigins *string

	hr, _, _ := i.Vtbl.GetAllowedOrigins.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&allowedOriginsCount)),
		uintptr(unsafe.Pointer(allowedOrigins)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, nil, syscall.Errno(hr)
	}
	return allowedOriginsCount, allowedOrigins, nil
}

func (i *ICoreWebView2CustomSchemeRegistration) SetAllowedOrigins(allowedOriginsCount uint32, allowedOrigins string) error {

	// Convert string 'allowedOrigins' to *uint16
	_allowedOrigins, err := UTF16PtrFromString(allowedOrigins)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.SetAllowedOrigins.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&allowedOriginsCount)),
		uintptr(unsafe.Pointer(_allowedOrigins)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2CustomSchemeRegistration) GetHasAuthorityComponent() (bool, error) {
	// Create int32 to hold bool result
	var _hasAuthorityComponent int32

	hr, _, _ := i.Vtbl.GetHasAuthorityComponent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_hasAuthorityComponent)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	hasAuthorityComponent := _hasAuthorityComponent != 0
	return hasAuthorityComponent, nil
}

func (i *ICoreWebView2CustomSchemeRegistration) PutHasAuthorityComponent(hasAuthorityComponent bool) error {

	hr, _, _ := i.Vtbl.PutHasAuthorityComponent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&hasAuthorityComponent)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
