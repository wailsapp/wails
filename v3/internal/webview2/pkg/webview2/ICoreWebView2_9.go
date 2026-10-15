//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_9Vtbl struct {
	ICoreWebView2_8Vtbl
	AddIsDefaultDownloadDialogOpenChanged ComProc
	RemoveIsDefaultDownloadDialogOpenChanged ComProc
	GetIsDefaultDownloadDialogOpen ComProc
	OpenDefaultDownloadDialog ComProc
	CloseDefaultDownloadDialog ComProc
	GetDefaultDownloadDialogCornerAlignment ComProc
	PutDefaultDownloadDialogCornerAlignment ComProc
	GetDefaultDownloadDialogMargin ComProc
	PutDefaultDownloadDialogMargin ComProc
}

type ICoreWebView2_9 struct {
	Vtbl *ICoreWebView2_9Vtbl
}

func (i *ICoreWebView2_9) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_9) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2_9 queries the object for its ICoreWebView2_9 interface. The receiver
// is the root of ICoreWebView2_9's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2) GetICoreWebView2_9() (*ICoreWebView2_9, error) {
	var result *ICoreWebView2_9

	iidICoreWebView2_9 := NewGUID("{4d7b2eab-9fdc-468d-b998-a9260b5ed651}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_9)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2_9) AddIsDefaultDownloadDialogOpenChanged(handler *ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddIsDefaultDownloadDialogOpenChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_9) RemoveIsDefaultDownloadDialogOpenChanged(token EventRegistrationToken) error {

	// 8/16-byte by-value arguments encode differently per architecture; the
	// arch consts are compile-time constants so dead branches are eliminated.
	var hr uintptr
	switch {
	case archIs386:
		hr, _, _ = i.Vtbl.RemoveIsDefaultDownloadDialogOpenChanged.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[0]),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&token)))[1]),
		)
	default:
		hr, _, _ = i.Vtbl.RemoveIsDefaultDownloadDialogOpenChanged.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(*(*uint64)(unsafe.Pointer(&token))),
		)
	}
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_9) GetIsDefaultDownloadDialogOpen() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsDefaultDownloadDialogOpen.Call(
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

func (i *ICoreWebView2_9) OpenDefaultDownloadDialog() error {


	hr, _, _ := i.Vtbl.OpenDefaultDownloadDialog.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_9) CloseDefaultDownloadDialog() error {


	hr, _, _ := i.Vtbl.CloseDefaultDownloadDialog.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_9) GetDefaultDownloadDialogCornerAlignment() (COREWEBVIEW2_DEFAULT_DOWNLOAD_DIALOG_CORNER_ALIGNMENT, error) {

	var value COREWEBVIEW2_DEFAULT_DOWNLOAD_DIALOG_CORNER_ALIGNMENT

	hr, _, _ := i.Vtbl.GetDefaultDownloadDialogCornerAlignment.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2_9) PutDefaultDownloadDialogCornerAlignment(value COREWEBVIEW2_DEFAULT_DOWNLOAD_DIALOG_CORNER_ALIGNMENT) error {


	hr, _, _ := i.Vtbl.PutDefaultDownloadDialogCornerAlignment.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_9) GetDefaultDownloadDialogMargin() (POINT, error) {

	var value POINT

	hr, _, _ := i.Vtbl.GetDefaultDownloadDialogMargin.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return POINT{}, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2_9) PutDefaultDownloadDialogMargin(value POINT) error {

	// 8/16-byte by-value arguments encode differently per architecture; the
	// arch consts are compile-time constants so dead branches are eliminated.
	var hr uintptr
	switch {
	case archIs386:
		hr, _, _ = i.Vtbl.PutDefaultDownloadDialogMargin.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&value)))[0]),
			uintptr((*(*[2]uint32)(unsafe.Pointer(&value)))[1]),
		)
	default:
		hr, _, _ = i.Vtbl.PutDefaultDownloadDialogMargin.Call(
			uintptr(unsafe.Pointer(i)),
			uintptr(*(*uint64)(unsafe.Pointer(&value))),
		)
	}
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
