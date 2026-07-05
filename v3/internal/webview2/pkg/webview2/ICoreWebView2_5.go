//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_5Vtbl struct {
	ICoreWebView2_4Vtbl
	AddClientCertificateRequested ComProc
	RemoveClientCertificateRequested ComProc
}

type ICoreWebView2_5 struct {
	Vtbl *ICoreWebView2_5Vtbl
}

func (i *ICoreWebView2_5) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_5) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2_5 queries the object for its ICoreWebView2_5 interface. The receiver
// is the root of ICoreWebView2_5's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2) GetICoreWebView2_5() (*ICoreWebView2_5, error) {
	var result *ICoreWebView2_5

	iidICoreWebView2_5 := NewGUID("{bedb11b8-d63c-11eb-b8bc-0242ac130003}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_5)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2_5) AddClientCertificateRequested(eventHandler *ICoreWebView2ClientCertificateRequestedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddClientCertificateRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_5) RemoveClientCertificateRequested(token EventRegistrationToken) error {

	// 8/16-byte by-value arguments encode differently per architecture; the
	// arch consts are compile-time constants so dead branches are eliminated.
	var hr uintptr
	switch {
	case archIs386:
		hr, _, _ = i.Vtbl.RemoveClientCertificateRequested.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[0]),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[1]),
		)
	default:
		hr, _, _ = i.Vtbl.RemoveClientCertificateRequested.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(*(*uint64)(unsafe.Pointer(&token))),
		)
	}
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
