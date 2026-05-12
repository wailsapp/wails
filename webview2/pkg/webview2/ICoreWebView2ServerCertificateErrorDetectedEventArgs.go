//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2ServerCertificateErrorDetectedEventArgsVtbl struct {
	IUnknownVtbl
	GetErrorStatus ComProc
	GetRequestUri ComProc
	GetServerCertificate ComProc
	GetAction ComProc
	PutAction ComProc
	GetDeferral ComProc
}

type ICoreWebView2ServerCertificateErrorDetectedEventArgs struct {
	Vtbl *ICoreWebView2ServerCertificateErrorDetectedEventArgsVtbl
}

func (i *ICoreWebView2ServerCertificateErrorDetectedEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ServerCertificateErrorDetectedEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2ServerCertificateErrorDetectedEventArgs) GetErrorStatus() (COREWEBVIEW2_WEB_ERROR_STATUS, error) {

	var value COREWEBVIEW2_WEB_ERROR_STATUS

	hr, _, err := i.Vtbl.GetErrorStatus.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2ServerCertificateErrorDetectedEventArgs) GetRequestUri() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, err := i.Vtbl.GetRequestUri.Call(
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

func (i *ICoreWebView2ServerCertificateErrorDetectedEventArgs) GetServerCertificate() (*ICoreWebView2Certificate, error) {

	var value *ICoreWebView2Certificate

	hr, _, err := i.Vtbl.GetServerCertificate.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2ServerCertificateErrorDetectedEventArgs) GetAction() (COREWEBVIEW2_SERVER_CERTIFICATE_ERROR_ACTION, error) {

	var value COREWEBVIEW2_SERVER_CERTIFICATE_ERROR_ACTION

	hr, _, err := i.Vtbl.GetAction.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2ServerCertificateErrorDetectedEventArgs) PutAction(value COREWEBVIEW2_SERVER_CERTIFICATE_ERROR_ACTION) error {


	hr, _, err := i.Vtbl.PutAction.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2ServerCertificateErrorDetectedEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {

	var deferral *ICoreWebView2Deferral

	hr, _, err := i.Vtbl.GetDeferral.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&deferral)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return deferral, err
}
