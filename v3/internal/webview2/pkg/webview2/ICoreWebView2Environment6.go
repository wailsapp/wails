//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Environment6Vtbl struct {
	IUnknownVtbl
	CreatePrintSettings ComProc
}

type ICoreWebView2Environment6 struct {
	Vtbl *ICoreWebView2Environment6Vtbl
}

func (i *ICoreWebView2Environment6) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Environment6() *ICoreWebView2Environment6 {
	var result *ICoreWebView2Environment6

	iidICoreWebView2Environment6 := NewGUID("{e59ee362-acbd-4857-9a8e-d3644d9459a9}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment6)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Environment6) CreatePrintSettings() (*ICoreWebView2PrintSettings, error) {

	var value *ICoreWebView2PrintSettings

	hr, _, _ := i.Vtbl.CreatePrintSettings.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
