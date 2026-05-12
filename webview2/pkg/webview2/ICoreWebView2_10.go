//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_10Vtbl struct {
	IUnknownVtbl
	AddBasicAuthenticationRequested ComProc
	RemoveBasicAuthenticationRequested ComProc
}

type ICoreWebView2_10 struct {
	Vtbl *ICoreWebView2_10Vtbl
}

func (i *ICoreWebView2_10) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_10) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2_10() (*ICoreWebView2_10, error) {
	var result *ICoreWebView2_10

	iidICoreWebView2_10 := NewGUID("{b1690564-6f5a-4983-8e48-31d1143fecdb}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_10)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2_10) AddBasicAuthenticationRequested(eventHandler *ICoreWebView2BasicAuthenticationRequestedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, err := i.Vtbl.AddBasicAuthenticationRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, err
}

func (i *ICoreWebView2_10) RemoveBasicAuthenticationRequested(token EventRegistrationToken) error {


	hr, _, err := i.Vtbl.RemoveBasicAuthenticationRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
