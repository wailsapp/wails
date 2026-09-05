//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Frame5Vtbl struct {
	IUnknownVtbl
	GetFrameId ComProc
}

type ICoreWebView2Frame5 struct {
	Vtbl *ICoreWebView2Frame5Vtbl
}

func (i *ICoreWebView2Frame5) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Frame5() *ICoreWebView2Frame5 {
	var result *ICoreWebView2Frame5

	iidICoreWebView2Frame5 := NewGUID("{99d199c4-7305-11ee-b962-0242ac120002}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Frame5)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Frame5) GetFrameId() (uint32, error) {

	var value uint32

	hr, _, _ := i.Vtbl.GetFrameId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}
