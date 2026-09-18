//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2EnvironmentOptions8Vtbl struct {
	IUnknownVtbl
	GetScrollBarStyle ComProc
	PutScrollBarStyle ComProc
}

type ICoreWebView2EnvironmentOptions8 struct {
	Vtbl *ICoreWebView2EnvironmentOptions8Vtbl
}

func (i *ICoreWebView2EnvironmentOptions8) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2EnvironmentOptions8) GetScrollBarStyle() (COREWEBVIEW2_SCROLLBAR_STYLE, error) {

	var value COREWEBVIEW2_SCROLLBAR_STYLE

	hr, _, _ := i.Vtbl.GetScrollBarStyle.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2EnvironmentOptions8) PutScrollBarStyle(value COREWEBVIEW2_SCROLLBAR_STYLE) error {

	hr, _, _ := i.Vtbl.PutScrollBarStyle.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
