//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2FrameInfoCollectionIteratorVtbl struct {
	IUnknownVtbl
	GetHasCurrent ComProc
	GetCurrent    ComProc
	MoveNext      ComProc
}

type ICoreWebView2FrameInfoCollectionIterator struct {
	Vtbl *ICoreWebView2FrameInfoCollectionIteratorVtbl
}

func (i *ICoreWebView2FrameInfoCollectionIterator) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2FrameInfoCollectionIterator) GetHasCurrent() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetHasCurrent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	value := _value != 0
	return value, nil
}

func (i *ICoreWebView2FrameInfoCollectionIterator) GetCurrent() (*ICoreWebView2FrameInfo, error) {

	var value *ICoreWebView2FrameInfo

	hr, _, _ := i.Vtbl.GetCurrent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2FrameInfoCollectionIterator) MoveNext() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.MoveNext.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	value := _value != 0
	return value, nil
}
