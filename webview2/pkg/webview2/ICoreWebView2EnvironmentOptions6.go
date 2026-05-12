//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2EnvironmentOptions6Vtbl struct {
	IUnknownVtbl
	GetAreBrowserExtensionsEnabled ComProc
	PutAreBrowserExtensionsEnabled ComProc
}

type ICoreWebView2EnvironmentOptions6 struct {
	Vtbl *ICoreWebView2EnvironmentOptions6Vtbl
}

func (i *ICoreWebView2EnvironmentOptions6) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2EnvironmentOptions6) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2EnvironmentOptions6) GetAreBrowserExtensionsEnabled() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, err := i.Vtbl.GetAreBrowserExtensionsEnabled.Call(
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

func (i *ICoreWebView2EnvironmentOptions6) PutAreBrowserExtensionsEnabled(value bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _value int32
	if value {
		_value = 1
	}

	hr, _, err := i.Vtbl.PutAreBrowserExtensionsEnabled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
