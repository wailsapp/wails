//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Frame6Vtbl struct {
	ICoreWebView2Frame5Vtbl
	AddScreenCaptureStarting ComProc
	RemoveScreenCaptureStarting ComProc
}

type ICoreWebView2Frame6 struct {
	Vtbl *ICoreWebView2Frame6Vtbl
}

func (i *ICoreWebView2Frame6) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Frame6) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2Frame6 queries the object for its ICoreWebView2Frame6 interface. The receiver
// is the root of ICoreWebView2Frame6's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2Frame) GetICoreWebView2Frame6() (*ICoreWebView2Frame6, error) {
	var result *ICoreWebView2Frame6

	iidICoreWebView2Frame6 := NewGUID("{0de611fd-31e9-5ddc-9d71-95eda26eff32}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Frame6)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Frame6) AddScreenCaptureStarting(eventHandler *ICoreWebView2FrameScreenCaptureStartingEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddScreenCaptureStarting.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2Frame6) RemoveScreenCaptureStarting(token EventRegistrationToken) error {

	// 8/16-byte by-value arguments encode differently per architecture; the
	// arch consts are compile-time constants so dead branches are eliminated.
	var hr uintptr
	switch {
	case archIs386:
		hr, _, _ = i.Vtbl.RemoveScreenCaptureStarting.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[0]),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[1]),
		)
	default:
		hr, _, _ = i.Vtbl.RemoveScreenCaptureStarting.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(*(*uint64)(unsafe.Pointer(&token))),
		)
	}
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
