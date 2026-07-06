//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_6Vtbl struct {
	IUnknownVtbl
	OpenTaskManagerWindow ComProc
}

type ICoreWebView2_6 struct {
	Vtbl *ICoreWebView2_6Vtbl
}

func (i *ICoreWebView2_6) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_6() *ICoreWebView2_6 {
	var result *ICoreWebView2_6

	iidICoreWebView2_6 := NewGUID("{499aadac-d92c-4589-8a75-111bfc167795}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_6)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_6) OpenTaskManagerWindow() error {

	hr, _, _ := i.Vtbl.OpenTaskManagerWindow.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
