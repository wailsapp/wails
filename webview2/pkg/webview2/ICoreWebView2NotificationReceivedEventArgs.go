//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2NotificationReceivedEventArgsVtbl struct {
	IUnknownVtbl
	GetSenderOrigin ComProc
	GetNotification ComProc
	PutHandled ComProc
	GetHandled ComProc
	GetDeferral ComProc
}

type ICoreWebView2NotificationReceivedEventArgs struct {
	Vtbl *ICoreWebView2NotificationReceivedEventArgsVtbl
}

func (i *ICoreWebView2NotificationReceivedEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2NotificationReceivedEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2NotificationReceivedEventArgs) GetSenderOrigin() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, err := i.Vtbl.GetSenderOrigin.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, err
}

func (i *ICoreWebView2NotificationReceivedEventArgs) GetNotification() (*ICoreWebView2Notification, error) {

	var value *ICoreWebView2Notification

	hr, _, err := i.Vtbl.GetNotification.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2NotificationReceivedEventArgs) PutHandled(value bool) error {

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

func (i *ICoreWebView2NotificationReceivedEventArgs) GetHandled() (bool, error) {
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

func (i *ICoreWebView2NotificationReceivedEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {

	var deferral *ICoreWebView2Deferral

	hr, _, err := i.Vtbl.GetDeferral.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&deferral)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return deferral, err
}
