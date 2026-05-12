//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2MoveFocusRequestedEventArgsVtbl struct {
	IUnknownVtbl
	GetReason ComProc
	GetHandled ComProc
	PutHandled ComProc
}

type ICoreWebView2MoveFocusRequestedEventArgs struct {
	Vtbl *ICoreWebView2MoveFocusRequestedEventArgsVtbl
}

func (i *ICoreWebView2MoveFocusRequestedEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2MoveFocusRequestedEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2MoveFocusRequestedEventArgs) GetReason() (COREWEBVIEW2_MOVE_FOCUS_REASON, error) {

	var reason COREWEBVIEW2_MOVE_FOCUS_REASON

	hr, _, err := i.Vtbl.GetReason.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&reason)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return reason, err
}

func (i *ICoreWebView2MoveFocusRequestedEventArgs) GetHandled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, err := i.Vtbl.GetHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    value := _value != 0
	return value, err
}

func (i *ICoreWebView2MoveFocusRequestedEventArgs) PutHandled(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, err := i.Vtbl.PutHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
