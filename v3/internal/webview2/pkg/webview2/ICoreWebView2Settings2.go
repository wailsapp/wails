//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Settings2Vtbl struct {
	IUnknownVtbl
	GetUserAgent ComProc
	PutUserAgent ComProc
}

type ICoreWebView2Settings2 struct {
	Vtbl *ICoreWebView2Settings2Vtbl
}

func (i *ICoreWebView2Settings2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Settings2() *ICoreWebView2Settings2 {
	var result *ICoreWebView2Settings2

	iidICoreWebView2Settings2 := NewGUID("{ee9a0f68-f46c-4e32-ac23-ef8cac224d2a}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Settings2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Settings2) GetUserAgent() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetUserAgent.Call(
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

func (i *ICoreWebView2Settings2) PutUserAgent(value string) error {

	// Convert string 'value' to *uint16
	_value, err := UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PutUserAgent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
