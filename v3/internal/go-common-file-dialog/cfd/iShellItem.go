//go:build windows
// +build windows

package cfd

import (
	"github.com/go-ole/go-ole"
	"syscall"
	"unsafe"
)

var (
	procSHCreateItemFromParsingName = syscall.NewLazyDLL("Shell32.dll").NewProc("SHCreateItemFromParsingName")
	iidShellItem                    = ole.NewGUID("43826d1e-e718-42ee-bc55-a1e261c37bfe")
)

type iShellItem struct {
	vtbl *iShellItemVtbl
}

type iShellItemVtbl struct {
	iUnknownVtbl
	BindToHandler  uintptr
	GetParent      uintptr
	GetDisplayName uintptr // func (sigdnName SIGDN, ppszName *LPWSTR) HRESULT
	GetAttributes  uintptr
	Compare        uintptr
}

func newIShellItem(path string) (*iShellItem, error) {
	var shellItem *iShellItem
	pathPtr := ole.SysAllocString(path)
	ret, _, _ := procSHCreateItemFromParsingName.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		0,
		uintptr(unsafe.Pointer(iidShellItem)),
		uintptr(unsafe.Pointer(&shellItem)))
	return shellItem, hresultToError(ret)
}

func (vtbl *iShellItemVtbl) getDisplayName(objPtr unsafe.Pointer) (string, error) {
	var ptr *uint16
	ret, _, _ := syscall.Syscall(vtbl.GetDisplayName,
		2,
		uintptr(objPtr),
		0x80058000, // SIGDN_FILESYSPATH
		uintptr(unsafe.Pointer(&ptr)))
	if err := hresultToError(ret); err != nil {
		return "", err
	}
	defer ole.CoTaskMemFree(uintptr(unsafe.Pointer(ptr)))
	return ole.LpOleStrToString(ptr), nil
}
