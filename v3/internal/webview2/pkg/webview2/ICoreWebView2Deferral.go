//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2DeferralVtbl struct {
	IUnknownVtbl
	Complete ComProc
}

type ICoreWebView2Deferral struct {
	Vtbl *ICoreWebView2DeferralVtbl
}

func (i *ICoreWebView2Deferral) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2Deferral) Complete() error {

	hr, _, _ := i.Vtbl.Complete.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
