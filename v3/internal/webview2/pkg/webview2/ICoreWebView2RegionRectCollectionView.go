//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2RegionRectCollectionViewVtbl struct {
	IUnknownVtbl
	GetCount        ComProc
	GetValueAtIndex ComProc
}

type ICoreWebView2RegionRectCollectionView struct {
	Vtbl *ICoreWebView2RegionRectCollectionViewVtbl
}

func (i *ICoreWebView2RegionRectCollectionView) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2RegionRectCollectionView) GetCount() (uint32, error) {

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

func (i *ICoreWebView2RegionRectCollectionView) GetValueAtIndex(index uint32) (RECT, error) {

	var value RECT

	hr, _, _ := i.Vtbl.GetValueAtIndex.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&index)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return RECT{}, syscall.Errno(hr)
	}
	return value, nil
}
