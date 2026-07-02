//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Controller2Vtbl struct {
	IUnknownVtbl
	GetDefaultBackgroundColor ComProc
	PutDefaultBackgroundColor ComProc
}

type ICoreWebView2Controller2 struct {
	Vtbl *ICoreWebView2Controller2Vtbl
}

func (i *ICoreWebView2Controller2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Controller2() *ICoreWebView2Controller2 {
	var result *ICoreWebView2Controller2

	iidICoreWebView2Controller2 := NewGUID("{c979903e-d4ca-4228-92eb-47ee3fa96eab}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Controller2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Controller2) GetDefaultBackgroundColor() (COREWEBVIEW2_COLOR, error) {

	var value COREWEBVIEW2_COLOR

	hr, _, _ := i.Vtbl.GetDefaultBackgroundColor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return COREWEBVIEW2_COLOR{}, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2Controller2) PutDefaultBackgroundColor(value COREWEBVIEW2_COLOR) error {

	hr, _, _ := i.Vtbl.PutDefaultBackgroundColor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
