//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2SharedBufferVtbl struct {
	IUnknownVtbl
	GetSize ComProc
	GetBuffer ComProc
	OpenStream ComProc
	GetFileMappingHandle ComProc
	Close ComProc
}

type ICoreWebView2SharedBuffer struct {
	Vtbl *ICoreWebView2SharedBufferVtbl
}

func (i *ICoreWebView2SharedBuffer) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2SharedBuffer) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2SharedBuffer) GetSize() (uint64, error) {

	var value uint64

	hr, _, err := i.Vtbl.GetSize.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2SharedBuffer) GetBuffer() (*uint8, error) {

	var value *uint8

	hr, _, err := i.Vtbl.GetBuffer.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2SharedBuffer) OpenStream() (*IStream, error) {

	var value *IStream

	hr, _, err := i.Vtbl.OpenStream.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2SharedBuffer) GetFileMappingHandle() (HANDLE, error) {

	var value HANDLE

	hr, _, err := i.Vtbl.GetFileMappingHandle.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2SharedBuffer) Close() error {


	hr, _, err := i.Vtbl.Close.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
