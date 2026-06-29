//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2CookieVtbl struct {
	IUnknownVtbl
	GetName       ComProc
	GetValue      ComProc
	PutValue      ComProc
	GetDomain     ComProc
	GetPath       ComProc
	GetExpires    ComProc
	PutExpires    ComProc
	GetIsHttpOnly ComProc
	PutIsHttpOnly ComProc
	GetSameSite   ComProc
	PutSameSite   ComProc
	GetIsSecure   ComProc
	PutIsSecure   ComProc
	GetIsSession  ComProc
}

type ICoreWebView2Cookie struct {
	Vtbl *ICoreWebView2CookieVtbl
}

func (i *ICoreWebView2Cookie) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2Cookie) GetName() (string, error) {
	// Create *uint16 to hold result
	var _name *uint16

	hr, _, _ := i.Vtbl.GetName.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_name)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	name := UTF16PtrToString(_name)
	CoTaskMemFree(unsafe.Pointer(_name))
	return name, nil
}

func (i *ICoreWebView2Cookie) GetValue() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetValue.Call(
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

func (i *ICoreWebView2Cookie) PutValue(value string) error {

	// Convert string 'value' to *uint16
	_value, err := UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutValue.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Cookie) GetDomain() (string, error) {
	// Create *uint16 to hold result
	var _domain *uint16

	hr, _, _ := i.Vtbl.GetDomain.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_domain)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	domain := UTF16PtrToString(_domain)
	CoTaskMemFree(unsafe.Pointer(_domain))
	return domain, nil
}

func (i *ICoreWebView2Cookie) GetPath() (string, error) {
	// Create *uint16 to hold result
	var _path *uint16

	hr, _, _ := i.Vtbl.GetPath.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_path)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	path := UTF16PtrToString(_path)
	CoTaskMemFree(unsafe.Pointer(_path))
	return path, nil
}

func (i *ICoreWebView2Cookie) GetExpires() (float64, error) {

	var expires float64

	hr, _, _ := i.Vtbl.GetExpires.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&expires)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0.0, syscall.Errno(hr)
	}
	return expires, nil
}

func (i *ICoreWebView2Cookie) PutExpires(expires float64) error {

	hr, _, _ := i.Vtbl.PutExpires.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&expires)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Cookie) GetIsHttpOnly() (bool, error) {
	// Create int32 to hold bool result
	var _isHttpOnly int32

	hr, _, _ := i.Vtbl.GetIsHttpOnly.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_isHttpOnly)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	isHttpOnly := _isHttpOnly != 0
	return isHttpOnly, nil
}

func (i *ICoreWebView2Cookie) PutIsHttpOnly(isHttpOnly bool) error {

	hr, _, _ := i.Vtbl.PutIsHttpOnly.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&isHttpOnly)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Cookie) GetSameSite() (COREWEBVIEW2_COOKIE_SAME_SITE_KIND, error) {

	var sameSite COREWEBVIEW2_COOKIE_SAME_SITE_KIND

	hr, _, _ := i.Vtbl.GetSameSite.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&sameSite)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return sameSite, nil
}

func (i *ICoreWebView2Cookie) PutSameSite(sameSite COREWEBVIEW2_COOKIE_SAME_SITE_KIND) error {

	hr, _, _ := i.Vtbl.PutSameSite.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(sameSite),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Cookie) GetIsSecure() (bool, error) {
	// Create int32 to hold bool result
	var _isSecure int32

	hr, _, _ := i.Vtbl.GetIsSecure.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_isSecure)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	isSecure := _isSecure != 0
	return isSecure, nil
}

func (i *ICoreWebView2Cookie) PutIsSecure(isSecure bool) error {

	hr, _, _ := i.Vtbl.PutIsSecure.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&isSecure)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Cookie) GetIsSession() (bool, error) {
	// Create int32 to hold bool result
	var _isSession int32

	hr, _, _ := i.Vtbl.GetIsSession.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_isSession)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	isSession := _isSession != 0
	return isSession, nil
}
