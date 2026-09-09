//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2ProcessInfoVtbl struct {
	IUnknownVtbl
	GetProcessId ComProc
	GetKind      ComProc
}

type ICoreWebView2ProcessInfo struct {
	Vtbl *ICoreWebView2ProcessInfoVtbl
}

func (i *ICoreWebView2ProcessInfo) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2ProcessInfo) GetProcessId() (int32, error) {

	var value int32

	hr, _, _ := i.Vtbl.GetProcessId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2ProcessInfo) GetKind() (COREWEBVIEW2_PROCESS_KIND, error) {

	var kind COREWEBVIEW2_PROCESS_KIND

	hr, _, _ := i.Vtbl.GetKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&kind)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return kind, nil
}
