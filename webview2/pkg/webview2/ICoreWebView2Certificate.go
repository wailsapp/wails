//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2CertificateVtbl struct {
	IUnknownVtbl
	GetSubject ComProc
	GetIssuer ComProc
	GetValidFrom ComProc
	GetValidTo ComProc
	GetDerEncodedSerialNumber ComProc
	GetDisplayName ComProc
	ToPemEncoding ComProc
	GetPemEncodedIssuerCertificateChain ComProc
}

type ICoreWebView2Certificate struct {
	Vtbl *ICoreWebView2CertificateVtbl
}

func (i *ICoreWebView2Certificate) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Certificate) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2Certificate) GetSubject() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, err := i.Vtbl.GetSubject.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, err
}

func (i *ICoreWebView2Certificate) GetIssuer() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, err := i.Vtbl.GetIssuer.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, err
}

func (i *ICoreWebView2Certificate) GetValidFrom() (float64, error) {

	var value float64

	hr, _, err := i.Vtbl.GetValidFrom.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0.0, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2Certificate) GetValidTo() (float64, error) {

	var value float64

	hr, _, err := i.Vtbl.GetValidTo.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0.0, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2Certificate) GetDerEncodedSerialNumber() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, err := i.Vtbl.GetDerEncodedSerialNumber.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, err
}

func (i *ICoreWebView2Certificate) GetDisplayName() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, err := i.Vtbl.GetDisplayName.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, err
}

func (i *ICoreWebView2Certificate) ToPemEncoding() (string, error) {
	// Create *uint16 to hold result
	var _pemEncodedData *uint16


	hr, _, err := i.Vtbl.ToPemEncoding.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_pemEncodedData)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	pemEncodedData := UTF16PtrToString(_pemEncodedData)
	CoTaskMemFree(unsafe.Pointer(_pemEncodedData))
	return pemEncodedData, err
}

func (i *ICoreWebView2Certificate) GetPemEncodedIssuerCertificateChain() (*ICoreWebView2StringCollection, error) {

	var value *ICoreWebView2StringCollection

	hr, _, err := i.Vtbl.GetPemEncodedIssuerCertificateChain.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}
