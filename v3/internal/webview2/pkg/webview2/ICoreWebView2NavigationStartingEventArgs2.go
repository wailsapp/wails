//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2NavigationStartingEventArgs2Vtbl struct {
	IUnknownVtbl
	GetAdditionalAllowedFrameAncestors ComProc
	PutAdditionalAllowedFrameAncestors ComProc
}

type ICoreWebView2NavigationStartingEventArgs2 struct {
	Vtbl *ICoreWebView2NavigationStartingEventArgs2Vtbl
}

func (i *ICoreWebView2NavigationStartingEventArgs2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2NavigationStartingEventArgs2() *ICoreWebView2NavigationStartingEventArgs2 {
	var result *ICoreWebView2NavigationStartingEventArgs2

	iidICoreWebView2NavigationStartingEventArgs2 := NewGUID("{9086be93-91aa-472d-a7e0-579f2ba006ad}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2NavigationStartingEventArgs2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2NavigationStartingEventArgs2) GetAdditionalAllowedFrameAncestors() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetAdditionalAllowedFrameAncestors.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}

func (i *ICoreWebView2NavigationStartingEventArgs2) PutAdditionalAllowedFrameAncestors(value string) error {

	// Convert string 'value' to *uint16
	_value, err := UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutAdditionalAllowedFrameAncestors.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
