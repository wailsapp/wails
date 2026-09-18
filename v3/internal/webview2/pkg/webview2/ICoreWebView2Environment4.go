//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Environment4Vtbl struct {
	IUnknownVtbl
	GetAutomationProviderForWindow ComProc
}

type ICoreWebView2Environment4 struct {
	Vtbl *ICoreWebView2Environment4Vtbl
}

func (i *ICoreWebView2Environment4) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Environment4() *ICoreWebView2Environment4 {
	var result *ICoreWebView2Environment4

	iidICoreWebView2Environment4 := NewGUID("{20944379-6dcf-41d6-a0a0-abc0fc50de0d}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment4)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Environment4) GetAutomationProviderForWindow(hwnd HWND) (*IUnknown, error) {

	var value *IUnknown

	hr, _, _ := i.Vtbl.GetAutomationProviderForWindow.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&hwnd)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
