//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Settings8Vtbl struct {
	IUnknownVtbl
	GetIsReputationCheckingRequired ComProc
	PutIsReputationCheckingRequired ComProc
}

type ICoreWebView2Settings8 struct {
	Vtbl *ICoreWebView2Settings8Vtbl
}

func (i *ICoreWebView2Settings8) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Settings8() *ICoreWebView2Settings8 {
	var result *ICoreWebView2Settings8

	iidICoreWebView2Settings8 := NewGUID("{9e6b0e8f-86ad-4e81-8147-a9b5edb68650}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Settings8)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Settings8) GetIsReputationCheckingRequired() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsReputationCheckingRequired.Call(
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

func (i *ICoreWebView2Settings8) PutIsReputationCheckingRequired(value bool) error {

	hr, _, _ := i.Vtbl.PutIsReputationCheckingRequired.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
