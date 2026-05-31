//go:build windows

package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2FileVtbl struct {
	_IUnknownVtbl
	GetPath ComProc
}

type ICoreWebView2File struct {
	vtbl *_ICoreWebView2FileVtbl
}

func (i *ICoreWebView2File) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2File) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))

	return uint32(ret)
}

func (i *ICoreWebView2File) GetPath() (string, error) {
	
	var _path *uint16
	hr, _, _ := i.vtbl.GetPath.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_path)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", windows.Errno(hr)
	}

	path := windows.UTF16PtrToString(_path)
	windows.CoTaskMemFree(unsafe.Pointer(_path))
	return path, nil
}
