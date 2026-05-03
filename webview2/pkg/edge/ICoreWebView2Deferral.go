//go:build windows

package edge

import (
	"errors"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ICoreWebView2DeferralVtbl struct {
	_IUnknownVtbl
	Complete ComProc
}

type ICoreWebView2Deferral struct {
	Vtbl *ICoreWebView2DeferralVtbl
}

// AddRef increments the reference count of ICoreWebView2Deferral interface
func (i *ICoreWebView2Deferral) AddRef() error {
	_, _, err := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	if err != nil && !errors.Is(err, windows.ERROR_SUCCESS) {
		return err
	}

	return nil
}

// Release decrements the reference count of ICoreWebView2Deferral interface
func (i *ICoreWebView2Deferral) Release() error {
	_, _, err := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	if err != nil && !errors.Is(err, windows.ERROR_SUCCESS) {
		return err
	}

	return nil
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
