//go:build windows

/*
 * Copyright (C) 2019 Tad Vizbaras. All Rights Reserved.
 * Copyright (C) 2010-2012 The W32 Authors. All Rights Reserved.
 */
package w32

import (
	"github.com/wailsapp/go-webview2/pkg/combridge"
	"syscall"
	"unsafe"
)

var (
	modole32 = syscall.NewLazyDLL("ole32.dll")

	procCoInitializeEx        = modole32.NewProc("CoInitializeEx")
	procCoInitialize          = modole32.NewProc("CoInitialize")
	procOleInitialize         = modole32.NewProc("OleInitialize")
	procCoUninitialize        = modole32.NewProc("CoUninitialize")
	procCreateStreamOnHGlobal = modole32.NewProc("CreateStreamOnHGlobal")
	procRegisterDragDrop      = modole32.NewProc("RegisterDragDrop")
	procRevokeDragDrop        = modole32.NewProc("RevokeDragDrop")
)

func CoInitializeEx(coInit uintptr) HRESULT {
	ret, _, _ := procCoInitializeEx.Call(
		0,
		coInit)

	switch uint32(ret) {
	case E_INVALIDARG:
		panic("CoInitializeEx failed with E_INVALIDARG")
	case E_OUTOFMEMORY:
		panic("CoInitializeEx failed with E_OUTOFMEMORY")
	case E_UNEXPECTED:
		panic("CoInitializeEx failed with E_UNEXPECTED")
	}

	return HRESULT(ret)
}

func CoInitialize() {
	procCoInitialize.Call(0)
}

func CoUninitialize() {
	procCoUninitialize.Call()
}

func CreateStreamOnHGlobal(hGlobal HGLOBAL, fDeleteOnRelease bool) *IStream {
	stream := new(IStream)
	ret, _, _ := procCreateStreamOnHGlobal.Call(
		uintptr(hGlobal),
		uintptr(BoolToBOOL(fDeleteOnRelease)),
		uintptr(unsafe.Pointer(&stream)))

	switch uint32(ret) {
	case E_INVALIDARG:
		panic("CreateStreamOnHGlobal failed with E_INVALIDARG")
	case E_OUTOFMEMORY:
		panic("CreateStreamOnHGlobal failed with E_OUTOFMEMORY")
	case E_UNEXPECTED:
		panic("CreateStreamOnHGlobal failed with E_UNEXPECTED")
	}

	return stream
}
func OleInitialise() {
	procOleInitialize.Call()
}

func RegisterDragDrop(hwnd HWND, dropTarget *DropTarget) error {

	dt := combridge.New[iDropTarget](dropTarget)
	hr, _, _ := procRegisterDragDrop.Call(
		hwnd,
		dt.Ref(),
	)

	if hr != S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func RevokeDragDrop(hwnd HWND) error {
	hr, _, _ := procRevokeDragDrop.Call(
		hwnd,
	)

	if hr != S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
