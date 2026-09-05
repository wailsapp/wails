//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Environment11Vtbl struct {
	IUnknownVtbl
	GetFailureReportFolderPath ComProc
}

type ICoreWebView2Environment11 struct {
	Vtbl *ICoreWebView2Environment11Vtbl
}

func (i *ICoreWebView2Environment11) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Environment11() *ICoreWebView2Environment11 {
	var result *ICoreWebView2Environment11

	iidICoreWebView2Environment11 := NewGUID("{f0913dc6-a0ec-42ef-9805-91dff3a2966a}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment11)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Environment11) GetFailureReportFolderPath() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetFailureReportFolderPath.Call(
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
