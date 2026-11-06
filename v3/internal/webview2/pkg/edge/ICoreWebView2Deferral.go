//go:build windows

package edge

import (
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
func (i *ICoreWebView2Deferral) AddRef() uint32 {
	// The return value is the new reference count, as documented for IUnknown::AddRef.
	// It is for diagnostics only. The Errno returned by Call is just the thread's
	// last-error and says nothing about success, so it is ignored.
	r, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(r)
}

// Release decrements the reference count of ICoreWebView2Deferral interface
func (i *ICoreWebView2Deferral) Release() uint32 {
	// The return value is the new reference count, as documented for IUnknown::Release.
	// It is for diagnostics only. The Errno returned by Call is just the thread's
	// last-error and says nothing about success, so it is ignored.
	r, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(r)
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
