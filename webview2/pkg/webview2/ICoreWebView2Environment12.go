//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Environment12Vtbl struct {
	IUnknownVtbl
	CreateSharedBuffer ComProc
}

type ICoreWebView2Environment12 struct {
	Vtbl *ICoreWebView2Environment12Vtbl
}

func (i *ICoreWebView2Environment12) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Environment12() *ICoreWebView2Environment12 {
	var result *ICoreWebView2Environment12

	iidICoreWebView2Environment12 := NewGUID("{f503db9b-739f-48dd-b151-fdfcf253f54e}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment12)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Environment12) CreateSharedBuffer(Size uint64) (*ICoreWebView2SharedBuffer, error) {

	var value *ICoreWebView2SharedBuffer

	hr, _, _ := i.Vtbl.CreateSharedBuffer.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&Size)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
