//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2CookieManagerVtbl struct {
	IUnknownVtbl
	CreateCookie                   ComProc
	CopyCookie                     ComProc
	GetCookies                     ComProc
	AddOrUpdateCookie              ComProc
	DeleteCookie                   ComProc
	DeleteCookies                  ComProc
	DeleteCookiesWithDomainAndPath ComProc
	DeleteAllCookies               ComProc
}

type ICoreWebView2CookieManager struct {
	Vtbl *ICoreWebView2CookieManagerVtbl
}

func (i *ICoreWebView2CookieManager) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2CookieManager) CreateCookie(name string, value string, domain string, path string) (*ICoreWebView2Cookie, error) {

	// Convert string 'name' to *uint16
	_name, err := UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}
	// Convert string 'value' to *uint16
	_value, err := UTF16PtrFromString(value)
	if err != nil {
		return nil, err
	}
	// Convert string 'domain' to *uint16
	_domain, err := UTF16PtrFromString(domain)
	if err != nil {
		return nil, err
	}
	// Convert string 'path' to *uint16
	_path, err := UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	var cookie *ICoreWebView2Cookie

	hr, _, _ := i.Vtbl.CreateCookie.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_name)),
		uintptr(unsafe.Pointer(_value)),
		uintptr(unsafe.Pointer(_domain)),
		uintptr(unsafe.Pointer(_path)),
		uintptr(unsafe.Pointer(&cookie)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return cookie, nil
}

func (i *ICoreWebView2CookieManager) CopyCookie(cookieParam *ICoreWebView2Cookie) (*ICoreWebView2Cookie, error) {

	var cookie *ICoreWebView2Cookie

	hr, _, _ := i.Vtbl.CopyCookie.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(cookieParam)),
		uintptr(unsafe.Pointer(&cookie)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return cookie, nil
}

func (i *ICoreWebView2CookieManager) GetCookies(uri string, handler *ICoreWebView2GetCookiesCompletedHandler) error {

	// Convert string 'uri' to *uint16
	_uri, err := UTF16PtrFromString(uri)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.GetCookies.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2CookieManager) AddOrUpdateCookie(cookie *ICoreWebView2Cookie) error {

	hr, _, _ := i.Vtbl.AddOrUpdateCookie.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(cookie)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2CookieManager) DeleteCookie(cookie *ICoreWebView2Cookie) error {

	hr, _, _ := i.Vtbl.DeleteCookie.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(cookie)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2CookieManager) DeleteCookies(name string, uri string) error {

	// Convert string 'name' to *uint16
	_name, err := UTF16PtrFromString(name)
	if err != nil {
		return err
	}
	// Convert string 'uri' to *uint16
	_uri, err := UTF16PtrFromString(uri)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.DeleteCookies.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_name)),
		uintptr(unsafe.Pointer(_uri)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2CookieManager) DeleteCookiesWithDomainAndPath(name string, domain string, path string) error {

	// Convert string 'name' to *uint16
	_name, err := UTF16PtrFromString(name)
	if err != nil {
		return err
	}
	// Convert string 'domain' to *uint16
	_domain, err := UTF16PtrFromString(domain)
	if err != nil {
		return err
	}
	// Convert string 'path' to *uint16
	_path, err := UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.DeleteCookiesWithDomainAndPath.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_name)),
		uintptr(unsafe.Pointer(_domain)),
		uintptr(unsafe.Pointer(_path)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2CookieManager) DeleteAllCookies() error {

	hr, _, _ := i.Vtbl.DeleteAllCookies.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
