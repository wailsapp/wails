//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_15Vtbl struct {
	ICoreWebView2_14Vtbl
	AddFaviconChanged ComProc
	RemoveFaviconChanged ComProc
	GetFaviconUri ComProc
	GetFavicon ComProc
}

type ICoreWebView2_15 struct {
	Vtbl *ICoreWebView2_15Vtbl
}

func (i *ICoreWebView2_15) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_15) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2_15 queries the object for its ICoreWebView2_15 interface. The receiver
// is the root of ICoreWebView2_15's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2) GetICoreWebView2_15() (*ICoreWebView2_15, error) {
	var result *ICoreWebView2_15

	iidICoreWebView2_15 := NewGUID("{517B2D1D-7DAE-4A66-A4F4-10352FFB9518}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_15)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2_15) AddFaviconChanged(eventHandler *ICoreWebView2FaviconChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddFaviconChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_15) RemoveFaviconChanged(token EventRegistrationToken) error {

	// 8/16-byte by-value arguments encode differently per architecture; the
	// arch consts are compile-time constants so dead branches are eliminated.
	var hr uintptr
	switch {
	case archIs386:
		hr, _, _ = i.Vtbl.RemoveFaviconChanged.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[0]),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[1]),
		)
	default:
		hr, _, _ = i.Vtbl.RemoveFaviconChanged.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(*(*uint64)(unsafe.Pointer(&token))),
		)
	}
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_15) GetFaviconUri() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, _ := i.Vtbl.GetFaviconUri.Call(
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

func (i *ICoreWebView2_15) GetFavicon(format COREWEBVIEW2_FAVICON_IMAGE_FORMAT, completedHandler *ICoreWebView2GetFaviconCompletedHandler) error {


	hr, _, _ := i.Vtbl.GetFavicon.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(format),
		uintptr(unsafe.Pointer(completedHandler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
