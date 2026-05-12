//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2AcceleratorKeyPressedEventArgs2Vtbl struct {
	IUnknownVtbl
	GetIsBrowserAcceleratorKeyEnabled ComProc
	PutIsBrowserAcceleratorKeyEnabled ComProc
}

type ICoreWebView2AcceleratorKeyPressedEventArgs2 struct {
	Vtbl *ICoreWebView2AcceleratorKeyPressedEventArgs2Vtbl
}

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs2) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs2) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2AcceleratorKeyPressedEventArgs2() (*ICoreWebView2AcceleratorKeyPressedEventArgs2, error) {
	var result *ICoreWebView2AcceleratorKeyPressedEventArgs2

	iidICoreWebView2AcceleratorKeyPressedEventArgs2 := NewGUID("{03b2c8c8-7799-4e34-bd66-ed26aa85f2bf}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2AcceleratorKeyPressedEventArgs2)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2AcceleratorKeyPressedEventArgs2) GetIsBrowserAcceleratorKeyEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, err := i.Vtbl.GetIsBrowserAcceleratorKeyEnabled.Call(
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

func (i *ICoreWebView2AcceleratorKeyPressedEventArgs2) PutIsBrowserAcceleratorKeyEnabled(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, err := i.Vtbl.PutIsBrowserAcceleratorKeyEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
