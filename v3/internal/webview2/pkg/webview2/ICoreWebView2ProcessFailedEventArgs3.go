//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2ProcessFailedEventArgs3Vtbl struct {
	IUnknownVtbl
	GetFailureSourceModulePath ComProc
}

type ICoreWebView2ProcessFailedEventArgs3 struct {
	Vtbl *ICoreWebView2ProcessFailedEventArgs3Vtbl
}

func (i *ICoreWebView2ProcessFailedEventArgs3) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2ProcessFailedEventArgs3() *ICoreWebView2ProcessFailedEventArgs3 {
	var result *ICoreWebView2ProcessFailedEventArgs3

	iidICoreWebView2ProcessFailedEventArgs3 := NewGUID("{ab667428-094d-5fd1-b480-8b4c0fdbdf2f}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2ProcessFailedEventArgs3)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2ProcessFailedEventArgs3) GetFailureSourceModulePath() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetFailureSourceModulePath.Call(
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
