//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2ObjectCollectionVtbl struct {
	IUnknownVtbl
	RemoveValueAtIndex ComProc
	InsertValueAtIndex ComProc
}

type ICoreWebView2ObjectCollection struct {
	Vtbl *ICoreWebView2ObjectCollectionVtbl
}

func (i *ICoreWebView2ObjectCollection) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2ObjectCollection() *ICoreWebView2ObjectCollection {
	var result *ICoreWebView2ObjectCollection

	iidICoreWebView2ObjectCollection := NewGUID("{5cfec11c-25bd-4e8d-9e1a-7acdaeeec047}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2ObjectCollection)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2ObjectCollection) RemoveValueAtIndex(index uint32) error {

	hr, _, _ := i.Vtbl.RemoveValueAtIndex.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&index)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2ObjectCollection) InsertValueAtIndex(index uint32, value *IUnknown) error {

	hr, _, _ := i.Vtbl.InsertValueAtIndex.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&index)),
		uintptr(unsafe.Pointer(value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
