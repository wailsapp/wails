//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2HttpHeadersCollectionIteratorVtbl struct {
	IUnknownVtbl
	GetCurrentHeader    ComProc
	GetHasCurrentHeader ComProc
	MoveNext            ComProc
}

type ICoreWebView2HttpHeadersCollectionIterator struct {
	Vtbl *ICoreWebView2HttpHeadersCollectionIteratorVtbl
}

func (i *ICoreWebView2HttpHeadersCollectionIterator) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2HttpHeadersCollectionIterator) GetCurrentHeader() (string, string, error) {
	// Create *uint16 to hold result
	var _name *uint16
	// Create *uint16 to hold result
	var _value *uint16

	hr, _, _ := i.Vtbl.GetCurrentHeader.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_name)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", "", syscall.Errno(hr)
	}
	// Get result and cleanup
	name := UTF16PtrToString(_name)
	CoTaskMemFree(unsafe.Pointer(_name))
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return name, value, nil
}

func (i *ICoreWebView2HttpHeadersCollectionIterator) GetHasCurrentHeader() (bool, error) {
	// Create int32 to hold bool result
	var _hasCurrent int32

	hr, _, _ := i.Vtbl.GetHasCurrentHeader.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_hasCurrent)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	hasCurrent := _hasCurrent != 0
	return hasCurrent, nil
}

func (i *ICoreWebView2HttpHeadersCollectionIterator) MoveNext() (bool, error) {
	// Create int32 to hold bool result
	var _hasNext int32

	hr, _, _ := i.Vtbl.MoveNext.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_hasNext)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	hasNext := _hasNext != 0
	return hasNext, nil
}
