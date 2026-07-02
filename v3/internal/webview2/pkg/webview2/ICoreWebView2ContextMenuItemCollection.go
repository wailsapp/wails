//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2ContextMenuItemCollectionVtbl struct {
	IUnknownVtbl
	GetCount           ComProc
	GetValueAtIndex    ComProc
	RemoveValueAtIndex ComProc
	InsertValueAtIndex ComProc
}

type ICoreWebView2ContextMenuItemCollection struct {
	Vtbl *ICoreWebView2ContextMenuItemCollectionVtbl
}

func (i *ICoreWebView2ContextMenuItemCollection) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2ContextMenuItemCollection) GetCount() (uint32, error) {

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

func (i *ICoreWebView2ContextMenuItemCollection) GetValueAtIndex(index uint32) (*ICoreWebView2ContextMenuItem, error) {

	var value *ICoreWebView2ContextMenuItem

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

func (i *ICoreWebView2ContextMenuItemCollection) RemoveValueAtIndex(index uint32) error {

	hr, _, _ := i.Vtbl.RemoveValueAtIndex.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&index)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2ContextMenuItemCollection) InsertValueAtIndex(index uint32, value *ICoreWebView2ContextMenuItem) error {

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
