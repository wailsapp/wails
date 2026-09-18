//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2ObjectCollectionViewVtbl struct {
	IUnknownVtbl
	GetCount        ComProc
	GetValueAtIndex ComProc
}

type ICoreWebView2ObjectCollectionView struct {
	Vtbl *ICoreWebView2ObjectCollectionViewVtbl
}

func (i *ICoreWebView2ObjectCollectionView) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2ObjectCollectionView) GetCount() (uint32, error) {

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

func (i *ICoreWebView2ObjectCollectionView) GetValueAtIndex(index uint32) (*IUnknown, error) {

	var value *IUnknown

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
