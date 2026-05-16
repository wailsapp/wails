package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// ICoreWebView2CookieManager vtable
type iCoreWebView2CookieManagerVtbl struct {
	_IUnknownVtbl
	CreateCookie                   ComProc
	CopyCookie                     ComProc
	GetCookies                     ComProc
	AddOrUpdateCookie              ComProc
	DeleteCookie                   ComProc
	DeleteCookies                  ComProc
	DeleteCookiesWithDomainAndPath ComProc
	DeleteAllCookies               ComProc
}

// ICoreWebView2CookieManager represents the cookie manager interface
type ICoreWebView2CookieManager struct {
	vtbl *iCoreWebView2CookieManagerVtbl
}

// AddRef increments the reference count of ICoreWebView2CookieManager interface
func (i *ICoreWebView2CookieManager) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

// Release decrements the reference count of ICoreWebView2CookieManager interface
func (i *ICoreWebView2CookieManager) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

// CreateCookie creates a new cookie with the given parameters
func (i *ICoreWebView2CookieManager) CreateCookie(name, value, domain, path string) (*ICoreWebView2Cookie, error) {
	var cookie *ICoreWebView2Cookie
	
	nameutf16, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}
	valueutf16, err := windows.UTF16PtrFromString(value)
	if err != nil {
		return nil, err
	}
	domainutf16, err := windows.UTF16PtrFromString(domain)
	if err != nil {
		return nil, err
	}
	pathutf16, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}

	hr, _, _ := i.vtbl.CreateCookie.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(nameutf16)),
		uintptr(unsafe.Pointer(valueutf16)),
		uintptr(unsafe.Pointer(domainutf16)),
		uintptr(unsafe.Pointer(pathutf16)),
		uintptr(unsafe.Pointer(&cookie)),
	)
	if hr != 0 {
		return nil, syscall.Errno(hr)
	}
	return cookie, nil
}

// CopyCookie creates a copy of the given cookie
func (i *ICoreWebView2CookieManager) CopyCookie(cookie *ICoreWebView2Cookie) (*ICoreWebView2Cookie, error) {
	var newCookie *ICoreWebView2Cookie
	hr, _, _ := i.vtbl.CopyCookie.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(cookie)),
		uintptr(unsafe.Pointer(&newCookie)),
	)
	if hr != 0 {
		return nil, syscall.Errno(hr)
	}
	return newCookie, nil
}

// GetCookies gets all cookies matching the URI
func (i *ICoreWebView2CookieManager) GetCookies(uri string) (*ICoreWebView2CookieList, error) {
	var list *ICoreWebView2CookieList
	uriutf16, err := windows.UTF16PtrFromString(uri)
	if err != nil {
		return nil, err
	}

	hr, _, _ := i.vtbl.GetCookies.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(uriutf16)),
		uintptr(unsafe.Pointer(&list)),
	)
	if hr != 0 {
		return nil, syscall.Errno(hr)
	}
	return list, nil
}

// DeleteCookies deletes all cookies with matching name and uri
func (i *ICoreWebView2CookieManager) DeleteCookies(name, uri string) error {
	nameutf16, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return err
	}
	uriutf16, err := windows.UTF16PtrFromString(uri)
	if err != nil {
		return err
	}

	hr, _, _ := i.vtbl.DeleteCookies.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(nameutf16)),
		uintptr(unsafe.Pointer(uriutf16)),
	)
	if hr != 0 {
		return syscall.Errno(hr)
	}
	return nil
}

// DeleteCookiesWithDomainAndPath deletes all cookies matching the domain and path
func (i *ICoreWebView2CookieManager) DeleteCookiesWithDomainAndPath(domain, path string) error {
	domainutf16, err := windows.UTF16PtrFromString(domain)
	if err != nil {
		return err
	}
	pathutf16, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	hr, _, _ := i.vtbl.DeleteCookiesWithDomainAndPath.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(domainutf16)),
		uintptr(unsafe.Pointer(pathutf16)),
	)
	if hr != 0 {
		return syscall.Errno(hr)
	}
	return nil
}

// AddOrUpdateCookie adds or updates a cookie
func (i *ICoreWebView2CookieManager) AddOrUpdateCookie(cookie *ICoreWebView2Cookie) error {
	hr, _, _ := i.vtbl.AddOrUpdateCookie.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(cookie)),
	)
	if hr != 0 {
		return syscall.Errno(hr)
	}
	return nil
}

// DeleteCookie deletes a specific cookie
func (i *ICoreWebView2CookieManager) DeleteCookie(cookie *ICoreWebView2Cookie) error {
	hr, _, _ := i.vtbl.DeleteCookie.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(cookie)),
	)
	if hr != 0 {
		return syscall.Errno(hr)
	}
	return nil
}

// DeleteAllCookies deletes all cookies
func (i *ICoreWebView2CookieManager) DeleteAllCookies() error {
	hr, _, _ := i.vtbl.DeleteAllCookies.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if hr != 0 {
		return syscall.Errno(hr)
	}
	return nil
}
