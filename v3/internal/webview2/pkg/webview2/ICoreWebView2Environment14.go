//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2Environment14Vtbl struct {
	IUnknownVtbl
	CreateWebFileSystemFileHandle      ComProc
	CreateWebFileSystemDirectoryHandle ComProc
	CreateObjectCollection             ComProc
}

type ICoreWebView2Environment14 struct {
	Vtbl *ICoreWebView2Environment14Vtbl
}

func (i *ICoreWebView2Environment14) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2Environment14() *ICoreWebView2Environment14 {
	var result *ICoreWebView2Environment14

	iidICoreWebView2Environment14 := NewGUID("{a5e9fad9-c875-59da-9bd7-473aa5ca1cef}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment14)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2Environment14) CreateWebFileSystemFileHandle(path string, permission COREWEBVIEW2_FILE_SYSTEM_HANDLE_PERMISSION) (*ICoreWebView2FileSystemHandle, error) {

	// Convert string 'path' to *uint16
	_path, err := UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	var value *ICoreWebView2FileSystemHandle

	hr, _, _ := i.Vtbl.CreateWebFileSystemFileHandle.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_path)),
		uintptr(permission),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2Environment14) CreateWebFileSystemDirectoryHandle(path string, permission COREWEBVIEW2_FILE_SYSTEM_HANDLE_PERMISSION) (*ICoreWebView2FileSystemHandle, error) {

	// Convert string 'path' to *uint16
	_path, err := UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	var value *ICoreWebView2FileSystemHandle

	hr, _, _ := i.Vtbl.CreateWebFileSystemDirectoryHandle.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_path)),
		uintptr(permission),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2Environment14) CreateObjectCollection(length uint32, items **IUnknown) (*ICoreWebView2ObjectCollection, error) {

	var objectCollection *ICoreWebView2ObjectCollection

	hr, _, _ := i.Vtbl.CreateObjectCollection.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&length)),
		uintptr(unsafe.Pointer(&items)),
		uintptr(unsafe.Pointer(&objectCollection)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return objectCollection, nil
}
