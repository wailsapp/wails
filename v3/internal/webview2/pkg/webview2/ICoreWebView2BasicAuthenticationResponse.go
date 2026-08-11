//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2BasicAuthenticationResponseVtbl struct {
	IUnknownVtbl
	GetUserName ComProc
	PutUserName ComProc
	GetPassword ComProc
	PutPassword ComProc
}

type ICoreWebView2BasicAuthenticationResponse struct {
	Vtbl *ICoreWebView2BasicAuthenticationResponseVtbl
}

func (i *ICoreWebView2BasicAuthenticationResponse) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2BasicAuthenticationResponse) GetUserName() (string, error) {
	// Create *uint16 to hold result
	var _userName *uint16

	hr, _, _ := i.Vtbl.GetUserName.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_userName)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	userName := UTF16PtrToString(_userName)
	CoTaskMemFree(unsafe.Pointer(_userName))
	return userName, nil
}

func (i *ICoreWebView2BasicAuthenticationResponse) PutUserName(userName string) error {

	// Convert string 'userName' to *uint16
	_userName, err := UTF16PtrFromString(userName)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutUserName.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_userName)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2BasicAuthenticationResponse) GetPassword() (string, error) {
	// Create *uint16 to hold result
	var _password *uint16

	hr, _, _ := i.Vtbl.GetPassword.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_password)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	password := UTF16PtrToString(_password)
	CoTaskMemFree(unsafe.Pointer(_password))
	return password, nil
}

func (i *ICoreWebView2BasicAuthenticationResponse) PutPassword(password string) error {

	// Convert string 'password' to *uint16
	_password, err := UTF16PtrFromString(password)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutPassword.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_password)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
