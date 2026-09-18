//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2ProcessInfoCollectionVtbl struct {
	IUnknownVtbl
	GetCount        ComProc
	GetValueAtIndex ComProc
}

type ICoreWebView2ProcessInfoCollection struct {
	Vtbl *ICoreWebView2ProcessInfoCollectionVtbl
}

func (i *ICoreWebView2ProcessInfoCollection) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2ProcessInfoCollection) GetCount() (uint32, error) {

	var value uint32

	hr, _, _ := i.Vtbl.GetCount.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2ProcessInfoCollection) GetValueAtIndex(index uint32) (*ICoreWebView2ProcessInfo, error) {

	var value *ICoreWebView2ProcessInfo

	hr, _, _ := i.Vtbl.GetValueAtIndex.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&index)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}
