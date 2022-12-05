//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2HttpHeadersCollectionIteratorVtbl struct {
	_IUnknownVtbl
	GetCurrentHeader    ComProc
	GetHasCurrentHeader ComProc
	MoveNext            ComProc
}

type ICoreWebView2HttpHeadersCollectionIterator struct {
	vtbl *_ICoreWebView2HttpHeadersCollectionIteratorVtbl
}

func (i *ICoreWebView2HttpHeadersCollectionIterator) Release() error {
	return i.vtbl.CallRelease(unsafe.Pointer(i))
}

func (i *ICoreWebView2HttpHeadersCollectionIterator) HasCurrentHeader() (bool, error) {
	var hasHeader int32
	res, _, err := i.vtbl.GetHasCurrentHeader.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&hasHeader)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	if windows.Handle(res) != windows.S_OK {
		return false, syscall.Errno(res)
	}
	return hasHeader != 0, nil
}

func (i *ICoreWebView2HttpHeadersCollectionIterator) GetCurrentHeader() (string, string, error) {
	// Create *uint16 to hold result
	var _name *uint16
	var _value *uint16
	res, _, err := i.vtbl.GetCurrentHeader.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_name)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if err != windows.ERROR_SUCCESS {
		return "", "", err
	}
	if windows.Handle(res) != windows.S_OK {
		return "", "", syscall.Errno(res)
	}
	// Get result and cleanup
	name := windows.UTF16PtrToString(_name)
	windows.CoTaskMemFree(unsafe.Pointer(_name))
	value := windows.UTF16PtrToString(_value)
	windows.CoTaskMemFree(unsafe.Pointer(_value))
	return name, value, nil
}

func (i *ICoreWebView2HttpHeadersCollectionIterator) MoveNext() (bool, error) {
	var next int32
	res, _, err := i.vtbl.MoveNext.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&next)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	if windows.Handle(res) != windows.S_OK {
		return false, syscall.Errno(res)
	}
	return next != 0, nil
}
