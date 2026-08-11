//go:build windows

package webviewloader

import (
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modkernel32     = windows.NewLazySystemDLL("kernel32.dll")
	procGlobalAlloc = modkernel32.NewProc("GlobalAlloc")
	procGlobalFree  = modkernel32.NewProc("GlobalFree")

	modversion                 = windows.NewLazySystemDLL("version.dll")
	procGetFileVersionInfoSize = modversion.NewProc("GetFileVersionInfoSizeW")
	procGetFileVersionInfo     = modversion.NewProc("GetFileVersionInfoW")
	procVerQueryValue          = modversion.NewProc("VerQueryValueW")

	modole32           = windows.NewLazySystemDLL("ole32.dll")
	procCoTaskMemAlloc = modole32.NewProc("CoTaskMemAlloc")
)

func getFileVersionInfo(path string) ([]byte, error) {
	lptstrFilename, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}

	size, _, err := procGetFileVersionInfoSize.Call(
		uintptr(unsafe.Pointer(lptstrFilename)),
		0,
	)

	err = maskErrorSuccess(err)
	if size == 0 && err == nil {
		err = fmt.Errorf("GetFileVersionInfoSize failed")
	}

	if err != nil {
		return nil, err
	}

	data := make([]byte, size)
	ret, _, err := procGetFileVersionInfo.Call(
		uintptr(unsafe.Pointer(lptstrFilename)),
		0,
		uintptr(size),
		uintptr(unsafe.Pointer(&data[0])),
	)

	err = maskErrorSuccess(err)
	if ret == 0 && err == nil {
		err = fmt.Errorf("GetFileVersionInfo failed")
	}

	if err != nil {
		return nil, err
	}
	return data, nil
}

func verQueryValueString(block []byte, subBlock string) (string, error) {
	// Allocate memory from native side to make sure the block doesn't get moved
	// because we get a pointer into that memory block from the native verQueryValue
	// call back.
	pBlock := globalAlloc(0, uint32(len(block)))
	defer globalFree(unsafe.Pointer(pBlock))

	// Copy the memory region into native side memory
	copy(unsafe.Slice((*byte)(pBlock), len(block)), block)

	lpSubBlock, err := syscall.UTF16PtrFromString(subBlock)
	if err != nil {
		return "", err
	}

	var lplpBuffer unsafe.Pointer
	var puLen uint
	ret, _, err := procVerQueryValue.Call(
		uintptr(pBlock),
		uintptr(unsafe.Pointer(lpSubBlock)),
		uintptr(unsafe.Pointer(&lplpBuffer)),
		uintptr(unsafe.Pointer(&puLen)),
	)

	err = maskErrorSuccess(err)
	if ret == 0 && err == nil {
		err = fmt.Errorf("VerQueryValue failed")
	}

	if err != nil {
		return "", err
	}

	if puLen <= 1 {
		return "", nil
	}
	puLen -= 1 // Remove Null-Terminator

	wchar := unsafe.Slice((*uint16)(lplpBuffer), puLen)
	return string(utf16.Decode(wchar)), nil
}

func globalAlloc(uFlags uint, dwBytes uint32) unsafe.Pointer {
	ret, _, _ := procGlobalAlloc.Call(
		uintptr(uFlags),
		uintptr(dwBytes))

	if ret == 0 {
		panic("globalAlloc failed")
	}

	return unsafe.Pointer(ret)
}

func globalFree(data unsafe.Pointer) {
	ret, _, _ := procGlobalFree.Call(uintptr(data))
	if ret != 0 {
		panic("globalFree failed")
	}
}

func maskErrorSuccess(err error) error {
	if err == windows.ERROR_SUCCESS {
		return nil
	}
	return err
}

func coTaskMemAlloc(size int) unsafe.Pointer {
	ret, _, _ := procCoTaskMemAlloc.Call(
		uintptr(size))

	if ret == 0 {
		panic("coTaskMemAlloc failed")
	}
	return unsafe.Pointer(ret)
}
