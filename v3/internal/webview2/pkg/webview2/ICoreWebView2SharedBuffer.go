//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2SharedBufferVtbl struct {
	IUnknownVtbl
	GetSize              ComProc
	GetBuffer            ComProc
	OpenStream           ComProc
	GetFileMappingHandle ComProc
	Close                ComProc
}

type ICoreWebView2SharedBuffer struct {
	Vtbl *ICoreWebView2SharedBufferVtbl
}

func (i *ICoreWebView2SharedBuffer) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2SharedBuffer) GetSize() (uint64, error) {

	var value uint64

	hr, _, _ := i.Vtbl.GetSize.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2SharedBuffer) GetBuffer() (*uint8, error) {

	var value *uint8

	hr, _, _ := i.Vtbl.GetBuffer.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2SharedBuffer) OpenStream() (*IStream, error) {

	var value *IStream

	hr, _, _ := i.Vtbl.OpenStream.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2SharedBuffer) GetFileMappingHandle() (HANDLE, error) {

	var value HANDLE

	hr, _, _ := i.Vtbl.GetFileMappingHandle.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2SharedBuffer) Close() error {

	hr, _, _ := i.Vtbl.Close.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
