//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2ClientCertificateRequestedEventArgsVtbl struct {
	IUnknownVtbl
	GetHost                          ComProc
	GetPort                          ComProc
	GetIsProxy                       ComProc
	GetAllowedCertificateAuthorities ComProc
	GetMutuallyTrustedCertificates   ComProc
	GetSelectedCertificate           ComProc
	PutSelectedCertificate           ComProc
	GetCancel                        ComProc
	PutCancel                        ComProc
	GetHandled                       ComProc
	PutHandled                       ComProc
	GetDeferral                      ComProc
}

type ICoreWebView2ClientCertificateRequestedEventArgs struct {
	Vtbl *ICoreWebView2ClientCertificateRequestedEventArgsVtbl
}

func (i *ICoreWebView2ClientCertificateRequestedEventArgs) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2ClientCertificateRequestedEventArgs) GetHost() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetHost.Call(
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

func (i *ICoreWebView2ClientCertificateRequestedEventArgs) GetPort() (int, error) {

	var value int

	hr, _, _ := i.Vtbl.GetPort.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2ClientCertificateRequestedEventArgs) GetIsProxy() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsProxy.Call(
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

func (i *ICoreWebView2ClientCertificateRequestedEventArgs) GetAllowedCertificateAuthorities() (*ICoreWebView2StringCollection, error) {

	var value *ICoreWebView2StringCollection

	hr, _, _ := i.Vtbl.GetAllowedCertificateAuthorities.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2ClientCertificateRequestedEventArgs) GetMutuallyTrustedCertificates() (*ICoreWebView2ClientCertificateCollection, error) {

	var value *ICoreWebView2ClientCertificateCollection

	hr, _, _ := i.Vtbl.GetMutuallyTrustedCertificates.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2ClientCertificateRequestedEventArgs) GetSelectedCertificate() (*ICoreWebView2ClientCertificate, error) {

	var value *ICoreWebView2ClientCertificate

	hr, _, _ := i.Vtbl.GetSelectedCertificate.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2ClientCertificateRequestedEventArgs) PutSelectedCertificate(value *ICoreWebView2ClientCertificate) error {

	hr, _, _ := i.Vtbl.PutSelectedCertificate.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2ClientCertificateRequestedEventArgs) GetCancel() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetCancel.Call(
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

func (i *ICoreWebView2ClientCertificateRequestedEventArgs) PutCancel(value bool) error {

	hr, _, _ := i.Vtbl.PutCancel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2ClientCertificateRequestedEventArgs) GetHandled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetHandled.Call(
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

func (i *ICoreWebView2ClientCertificateRequestedEventArgs) PutHandled(value bool) error {

	hr, _, _ := i.Vtbl.PutHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2ClientCertificateRequestedEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {

	var deferral *ICoreWebView2Deferral

	hr, _, _ := i.Vtbl.GetDeferral.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&deferral)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return deferral, nil
}
