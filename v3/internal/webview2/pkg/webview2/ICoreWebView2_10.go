//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_10Vtbl struct {
	ICoreWebView2_9Vtbl
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


// GetICoreWebView2_10 queries the object for its ICoreWebView2_10 interface. The receiver
// is the root of ICoreWebView2_10's inheritance chain — the object that actually
// implements it.
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

	hr, _, _ := i.Vtbl.AddBasicAuthenticationRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_10) RemoveBasicAuthenticationRequested(token EventRegistrationToken) error {

	// 8/16-byte by-value arguments encode differently per architecture; the
	// arch consts are compile-time constants so dead branches are eliminated.
	var hr uintptr
	switch {
	case archIs386:
		hr, _, _ = i.Vtbl.RemoveBasicAuthenticationRequested.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[0]),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[1]),
		)
	default:
		hr, _, _ = i.Vtbl.RemoveBasicAuthenticationRequested.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(*(*uint64)(unsafe.Pointer(&token))),
		)
	}
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
