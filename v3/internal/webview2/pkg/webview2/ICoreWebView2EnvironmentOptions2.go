//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2EnvironmentOptions2Vtbl struct {
	IUnknownVtbl
	GetExclusiveUserDataFolderAccess ComProc
	PutExclusiveUserDataFolderAccess ComProc
}

type ICoreWebView2EnvironmentOptions2 struct {
	Vtbl *ICoreWebView2EnvironmentOptions2Vtbl
}

func (i *ICoreWebView2EnvironmentOptions2) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2EnvironmentOptions2) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2EnvironmentOptions2) GetExclusiveUserDataFolderAccess() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetExclusiveUserDataFolderAccess.Call(
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

func (i *ICoreWebView2EnvironmentOptions2) PutExclusiveUserDataFolderAccess(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, _ := i.Vtbl.PutExclusiveUserDataFolderAccess.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
