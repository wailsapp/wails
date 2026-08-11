//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_20Vtbl struct {
	IUnknownVtbl
	GetFrameId ComProc
}

type ICoreWebView2_20 struct {
	Vtbl *ICoreWebView2_20Vtbl
}

func (i *ICoreWebView2_20) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_20() *ICoreWebView2_20 {
	var result *ICoreWebView2_20

	iidICoreWebView2_20 := NewGUID("{b4bc1926-7305-11ee-b962-0242ac120002}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_20)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_20) GetFrameId() (uint32, error) {

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
