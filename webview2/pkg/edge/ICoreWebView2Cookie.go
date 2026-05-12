package edge

import (
	"unsafe"
	"math"

	"golang.org/x/sys/windows"
)

// ICoreWebView2Cookie vtable
type iCoreWebView2CookieVtbl struct {
	_IUnknownVtbl
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
}

// ICoreWebView2Cookie represents a cookie
type ICoreWebView2Cookie struct {
	vtbl *iCoreWebView2CookieVtbl
}

// Addref increments refernce count of the ICoreWebView2Cookie interface
func (i *ICoreWebView2Cookie) AddRef() uintptr {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return ret
}

// Release decrements reference count of the ICoreWebView2Cookie interface
func (i *ICoreWebView2Cookie) Release() uintptr {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return ret
}

// GetName gets the cookie name
func (i *ICoreWebView2Cookie) GetName() (string, error) {
	var name *uint16
	hr, _, _ := i.vtbl.GetName.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&name)),
	)
	if hr != 0 {
		return "", windows.Errno(hr)
	}
	return windows.UTF16PtrToString(name), nil
}

// GetValue gets the cookie value
func (i *ICoreWebView2Cookie) GetValue() (string, error) {
	var value *uint16
	hr, _, _ := i.vtbl.GetValue.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if hr != 0 {
		return "", windows.Errno(hr)
	}
	return windows.UTF16PtrToString(value), nil
}

// PutValue sets the cookie value
func (i *ICoreWebView2Cookie) PutValue(value string) error {
	ptr, err := windows.UTF16PtrFromString(value)
	if err != nil {
		return err
	}
	hr, _, _ := i.vtbl.PutValue.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(ptr)),
	)
	if hr != 0 {
		return windows.Errno(hr)
	}
	return nil
}

// GetDomain gets the cookie domain
func (i *ICoreWebView2Cookie) GetDomain() (string, error) {
	var domain *uint16
	hr, _, _ := i.vtbl.GetDomain.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&domain)),
	)
	if hr != 0 {
		return "", windows.Errno(hr)
	}
	return windows.UTF16PtrToString(domain), nil
}

// GetPath gets the cookie path
func (i *ICoreWebView2Cookie) GetPath() (string, error) {
	var path *uint16
	hr, _, _ := i.vtbl.GetPath.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&path)),
	)
	if hr != 0 {
		return "", windows.Errno(hr)
	}
	return windows.UTF16PtrToString(path), nil
}

// GetExpires gets the cookie expiration time
func (i *ICoreWebView2Cookie) GetExpires() (float64, error) {
	var expiresUint64 uint64
	hr, _, _ := i.vtbl.GetExpires.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&expiresUint64)),
	)
	if hr != 0 {
		return 0.0, windows.Errno(hr)
	}
	return math.Float64frombits(expiresUint64), nil
}

// PutExpires sets the cookie expiration time
func (i *ICoreWebView2Cookie) PutExpires(expires float64) error {
	hr, _, _ := i.vtbl.PutExpires.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(math.Float64bits(expires)),
	)
	if hr != 0 {
		return windows.Errno(hr)
	}
	return nil
}

// GetIsHttpOnly gets whether the cookie is HTTP-only
func (i *ICoreWebView2Cookie) GetIsHttpOnly() (bool, error) {
	var isHttpOnly int32
	hr, _, _ := i.vtbl.GetIsHttpOnly.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&isHttpOnly)),
	)
	if hr != 0 {
		return false, windows.Errno(hr)
	}
	return isHttpOnly != 0, nil
}

// PutIsHttpOnly sets whether the cookie is HTTP-only
func (i *ICoreWebView2Cookie) PutIsHttpOnly(isHttpOnly bool) error {
	value := int32(0)
	if isHttpOnly {
		value = 1
	}
	hr, _, _ := i.vtbl.PutIsHttpOnly.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if hr != 0 {
		return windows.Errno(hr)
	}
	return nil
}

// GetSameSite gets the cookie's SameSite attribute
func (i *ICoreWebView2Cookie) GetSameSite() (int32, error) {
	var sameSite int32
	hr, _, _ := i.vtbl.GetSameSite.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&sameSite)),
	)
	if hr != 0 {
		return 0, windows.Errno(hr)
	}
	return sameSite, nil
}

// PutSameSite sets the cookie's SameSite attribute
func (i *ICoreWebView2Cookie) PutSameSite(sameSite int32) error {
	hr, _, _ := i.vtbl.PutSameSite.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(sameSite),
	)
	if hr != 0 {
		return windows.Errno(hr)
	}
	return nil
}

// GetIsSecure gets whether the cookie is secure
func (i *ICoreWebView2Cookie) GetIsSecure() (bool, error) {
	var isSecure int32
	hr, _, _ := i.vtbl.GetIsSecure.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&isSecure)),
	)
	if hr != 0 {
		return false, windows.Errno(hr)
	}
	return isSecure != 0, nil
}

// PutIsSecure sets whether the cookie is secure
func (i *ICoreWebView2Cookie) PutIsSecure(isSecure bool) error {
	value := int32(0)
	if isSecure {
		value = 1
	}
	hr, _, _ := i.vtbl.PutIsSecure.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if hr != 0 {
		return windows.Errno(hr)
	}
	return nil
}

// QueryInterface queries for a specific interface
func (i *ICoreWebView2Cookie) QueryInterface(riid *windows.GUID, ppvObject *unsafe.Pointer) error {
	hr, _, _ := i.vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(ppvObject)),
	)
	if hr != 0 {
		return windows.Errno(hr)
	}
	return nil
}
