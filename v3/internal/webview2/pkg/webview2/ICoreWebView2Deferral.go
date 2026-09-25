//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2DeferralVtbl struct {
	IUnknownVtbl
	Complete ComProc
}

type ICoreWebView2Deferral struct {
	Vtbl *ICoreWebView2DeferralVtbl
}

func (i *ICoreWebView2Deferral) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Deferral) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
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
