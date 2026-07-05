//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Frame3Vtbl struct {
	ICoreWebView2Frame2Vtbl
	AddPermissionRequested ComProc
	RemovePermissionRequested ComProc
}

type ICoreWebView2Frame3 struct {
	Vtbl *ICoreWebView2Frame3Vtbl
}

func (i *ICoreWebView2Frame3) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Frame3) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2Frame3 queries the object for its ICoreWebView2Frame3 interface. The receiver
// is the root of ICoreWebView2Frame3's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2Frame) GetICoreWebView2Frame3() (*ICoreWebView2Frame3, error) {
	var result *ICoreWebView2Frame3

	iidICoreWebView2Frame3 := NewGUID("{b50d82cc-cc28-481d-9614-cb048895e6a0}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Frame3)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Frame3) AddPermissionRequested(eventHandler *ICoreWebView2FramePermissionRequestedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddPermissionRequested.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2Frame3) RemovePermissionRequested(token EventRegistrationToken) error {

	// 8/16-byte by-value arguments encode differently per architecture; the
	// arch consts are compile-time constants so dead branches are eliminated.
	var hr uintptr
	switch {
	case archIs386:
		hr, _, _ = i.Vtbl.RemovePermissionRequested.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[0]),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[1]),
		)
	default:
		hr, _, _ = i.Vtbl.RemovePermissionRequested.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(*(*uint64)(unsafe.Pointer(&token))),
		)
	}
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
