//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2ContextMenuRequestedEventArgsVtbl struct {
	IUnknownVtbl
	GetMenuItems ComProc
	GetContextMenuTarget ComProc
	GetLocation ComProc
	PutSelectedCommandId ComProc
	GetSelectedCommandId ComProc
	PutHandled ComProc
	GetHandled ComProc
	GetDeferral ComProc
}

type ICoreWebView2ContextMenuRequestedEventArgs struct {
	Vtbl *ICoreWebView2ContextMenuRequestedEventArgsVtbl
}

func (i *ICoreWebView2ContextMenuRequestedEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ContextMenuRequestedEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2ContextMenuRequestedEventArgs) GetMenuItems() (*ICoreWebView2ContextMenuItemCollection, error) {

	var value *ICoreWebView2ContextMenuItemCollection

	hr, _, err := i.Vtbl.GetMenuItems.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2ContextMenuRequestedEventArgs) GetContextMenuTarget() (*ICoreWebView2ContextMenuTarget, error) {

	var value *ICoreWebView2ContextMenuTarget

	hr, _, err := i.Vtbl.GetContextMenuTarget.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2ContextMenuRequestedEventArgs) GetLocation() (POINT, error) {

	var value POINT

	hr, _, err := i.Vtbl.GetLocation.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return POINT{}, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2ContextMenuRequestedEventArgs) PutSelectedCommandId(value int32) error {


	hr, _, err := i.Vtbl.PutSelectedCommandId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2ContextMenuRequestedEventArgs) GetSelectedCommandId() (int32, error) {

	var value int32

	hr, _, err := i.Vtbl.GetSelectedCommandId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2ContextMenuRequestedEventArgs) PutHandled(value bool) error {

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

func (i *ICoreWebView2ContextMenuRequestedEventArgs) GetHandled() (bool, error) {
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

func (i *ICoreWebView2ContextMenuRequestedEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {

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
