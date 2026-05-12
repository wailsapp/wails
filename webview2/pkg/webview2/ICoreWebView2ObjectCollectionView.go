//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2ObjectCollectionViewVtbl struct {
	IUnknownVtbl
	GetCount ComProc
	GetValueAtIndex ComProc
}

type ICoreWebView2ObjectCollectionView struct {
	Vtbl *ICoreWebView2ObjectCollectionViewVtbl
}

func (i *ICoreWebView2ObjectCollectionView) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ObjectCollectionView) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2ObjectCollectionView) GetCount() (uint32, error) {

	var value uint32

	hr, _, err := i.Vtbl.GetCount.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2ObjectCollectionView) GetValueAtIndex(index uint32) (*IUnknown, error) {

	var value *IUnknown

	hr, _, err := i.Vtbl.GetValueAtIndex.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(index),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}
