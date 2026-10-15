//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2BasicAuthenticationRequestedEventArgsVtbl struct {
	IUnknownVtbl
	GetUri ComProc
	GetChallenge ComProc
	GetResponse ComProc
	GetCancel ComProc
	PutCancel ComProc
	GetDeferral ComProc
}

type ICoreWebView2BasicAuthenticationRequestedEventArgs struct {
	Vtbl *ICoreWebView2BasicAuthenticationRequestedEventArgsVtbl
}

func (i *ICoreWebView2BasicAuthenticationRequestedEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2BasicAuthenticationRequestedEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2BasicAuthenticationRequestedEventArgs) GetUri() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, _ := i.Vtbl.GetUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}

func (i *ICoreWebView2BasicAuthenticationRequestedEventArgs) GetChallenge() (string, error) {
	// Create *uint16 to hold result
	var _challenge *uint16


	hr, _, _ := i.Vtbl.GetChallenge.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_challenge)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	challenge := UTF16PtrToString(_challenge)
	CoTaskMemFree(unsafe.Pointer(_challenge))
	return challenge, nil
}

func (i *ICoreWebView2BasicAuthenticationRequestedEventArgs) GetResponse() (*ICoreWebView2BasicAuthenticationResponse, error) {

	var response *ICoreWebView2BasicAuthenticationResponse

	hr, _, _ := i.Vtbl.GetResponse.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&response)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return response, nil
}

func (i *ICoreWebView2BasicAuthenticationRequestedEventArgs) GetCancel() (bool, error) {
	// Create int32 to hold bool result
	var _cancel int32

	hr, _, _ := i.Vtbl.GetCancel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_cancel)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    cancel := _cancel != 0
	return cancel, nil
}

func (i *ICoreWebView2BasicAuthenticationRequestedEventArgs) PutCancel(cancel bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _cancel int32
	if cancel {
		_cancel = 1
	}

	hr, _, _ := i.Vtbl.PutCancel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_cancel),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2BasicAuthenticationRequestedEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {

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
