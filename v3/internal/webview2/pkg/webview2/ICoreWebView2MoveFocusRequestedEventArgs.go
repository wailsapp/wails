//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2MoveFocusRequestedEventArgsVtbl struct {
	IUnknownVtbl
	GetReason  ComProc
	GetHandled ComProc
	PutHandled ComProc
}

type ICoreWebView2MoveFocusRequestedEventArgs struct {
	Vtbl *ICoreWebView2MoveFocusRequestedEventArgsVtbl
}

func (i *ICoreWebView2MoveFocusRequestedEventArgs) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2MoveFocusRequestedEventArgs) GetReason() (COREWEBVIEW2_MOVE_FOCUS_REASON, error) {

	var reason COREWEBVIEW2_MOVE_FOCUS_REASON

	hr, _, _ := i.Vtbl.GetReason.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&reason)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return reason, nil
}

func (i *ICoreWebView2MoveFocusRequestedEventArgs) GetHandled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetHandled.Call(
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

func (i *ICoreWebView2MoveFocusRequestedEventArgs) PutHandled(value bool) error {

	hr, _, _ := i.Vtbl.PutHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
