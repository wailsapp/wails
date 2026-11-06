//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	ERROR_ELEMENT_NOT_FOUND syscall.Errno = 0x80070490
)

type _ICoreWebView2HttpRequestHeadersVtbl struct {
	_IUnknownVtbl
	GetHeader    ComProc
	GetHeaders   ComProc
	Contains     ComProc
	SetHeader    ComProc
	RemoveHeader ComProc
	GetIterator  ComProc
}

type ICoreWebView2HttpRequestHeaders struct {
	vtbl *_ICoreWebView2HttpRequestHeadersVtbl
}

func (i *ICoreWebView2HttpRequestHeaders) AddRef() uint32 {
	// The return value is the new reference count, as documented for IUnknown::AddRef.
	// It is for diagnostics only. The Errno returned by Call is just the thread's
	// last-error and says nothing about success, so it is ignored.
	r, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(r)
}

func (i *ICoreWebView2HttpRequestHeaders) Release() uint32 {
	// The return value is the new reference count, as documented for IUnknown::Release.
	// It is for diagnostics only. The Errno returned by Call is just the thread's
	// last-error and says nothing about success, so it is ignored.
	r, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(r)
}

// GetHeader returns the value of the specified header. If the header is not found
// ERROR_ELEMENT_NOT_FOUND is returned as error.
func (i *ICoreWebView2HttpRequestHeaders) GetHeader(name string) (string, error) {
	_name, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return "", err
	}

	var _value *uint16
	hr, _, _ := i.vtbl.GetHeader.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_name)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", windows.Errno(hr)
	}

	value := windows.UTF16PtrToString(_value)
	windows.CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}

// SetHeader sets the specified header to the value.
func (i *ICoreWebView2HttpRequestHeaders) SetHeader(name, value string) error {
	_name, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return err
	}

	_value, err := windows.UTF16PtrFromString(value)
	if err != nil {
		return err
	}

	hr, _, _ := i.vtbl.SetHeader.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_name)),
		uintptr(unsafe.Pointer(_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}

	return nil
}

// GetIterator returns an iterator over the collection of request headers. Make sure to call
// Release on the returned Object after finished using it.
func (i *ICoreWebView2HttpRequestHeaders) GetIterator() (*ICoreWebView2HttpHeadersCollectionIterator, error) {
	var headers *ICoreWebView2HttpHeadersCollectionIterator
	hr, _, _ := i.vtbl.GetIterator.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&headers)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, windows.Errno(hr)
	}

	return headers, nil
}
