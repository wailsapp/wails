//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2CustomSchemeRegistrationVtbl struct {
	IUnknownVtbl
	GetSchemeName ComProc
	GetTreatAsSecure ComProc
	PutTreatAsSecure ComProc
	GetAllowedOrigins ComProc
	SetAllowedOrigins ComProc
	GetHasAuthorityComponent ComProc
	PutHasAuthorityComponent ComProc
}

type ICoreWebView2CustomSchemeRegistration struct {
	Vtbl *ICoreWebView2CustomSchemeRegistrationVtbl
}

func (i *ICoreWebView2CustomSchemeRegistration) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2CustomSchemeRegistration) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2CustomSchemeRegistration) GetSchemeName() (string, error) {
	// Create *uint16 to hold result
	var _schemeName *uint16


	hr, _, _ := i.Vtbl.GetSchemeName.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_schemeName)),
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

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, _ := i.Vtbl.PutTreatAsSecure.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
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
		uintptr(unsafe.Pointer(&allowedOrigins)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, nil, syscall.Errno(hr)
	}
	return allowedOriginsCount, allowedOrigins, nil
}

func (i *ICoreWebView2CustomSchemeRegistration) SetAllowedOrigins(allowedOriginsCount uint32, allowedOrigins []string) error {
	// Convert []string 'allowedOrigins' to **uint16 (LPCWSTR* / LPWSTR*)
	_allowedOriginsptrs := make([]*uint16, len(allowedOrigins))
	for _i, _s := range allowedOrigins {
		_p, err := UTF16PtrFromString(_s)
		if err != nil {
			return err
		}
		_allowedOriginsptrs[_i] = _p
	}
	var _allowedOrigins **uint16
	if len(_allowedOriginsptrs) > 0 {
		_allowedOrigins = &_allowedOriginsptrs[0]
	}


	hr, _, _ := i.Vtbl.SetAllowedOrigins.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(allowedOriginsCount),
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

	// Convert Go bool to COM BOOL (int32)
	var _hasAuthorityComponent int32
	if hasAuthorityComponent {
		_hasAuthorityComponent = 1
	}

	hr, _, _ := i.Vtbl.PutHasAuthorityComponent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_hasAuthorityComponent),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
