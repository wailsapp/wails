//go:build windows

package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2MoveFocusRequestedEventArgsVtbl struct {
	_IUnknownVtbl
	GetReason  ComProc
	GetHandled ComProc
	PutHandled ComProc
}

type ICoreWebView2MoveFocusRequestedEventArgs struct {
	vtbl *_ICoreWebView2MoveFocusRequestedEventArgsVtbl
}

func (i *ICoreWebView2MoveFocusRequestedEventArgs) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2MoveFocusRequestedEventArgs) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2MoveFocusRequestedEventArgs) GetReason() (COREWEBVIEW2_MOVE_FOCUS_REASON, error) {
	var reason COREWEBVIEW2_MOVE_FOCUS_REASON
	hr, _, _ := i.vtbl.GetReason.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&reason)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, windows.Errno(hr)
	}
	return reason, nil
}

func (i *ICoreWebView2MoveFocusRequestedEventArgs) PutHandled(handled bool) error {
	hr, _, _ := i.vtbl.PutHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(handled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
}
