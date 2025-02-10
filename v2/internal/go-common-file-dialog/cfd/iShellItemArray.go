//go:build windows
// +build windows

package cfd

import (
	"github.com/go-ole/go-ole"
	"syscall"
	"unsafe"
)

const (
	iidShellItemArrayGUID = "{b63ea76d-1f85-456f-a19c-48159efa858b}"
)

var (
	iidShellItemArray *ole.GUID
)

func init() {
	iidShellItemArray, _ = ole.IIDFromString(iidShellItemArrayGUID)
}

type iShellItemArray struct {
	vtbl *iShellItemArrayVtbl
}

type iShellItemArrayVtbl struct {
	iUnknownVtbl
	BindToHandler              uintptr
	GetPropertyStore           uintptr
	GetPropertyDescriptionList uintptr
	GetAttributes              uintptr
	GetCount                   uintptr // func (pdwNumItems *DWORD) HRESULT
	GetItemAt                  uintptr // func (dwIndex DWORD, ppsi **IShellItem) HRESULT
	EnumItems                  uintptr
}

func (vtbl *iShellItemArrayVtbl) getCount(objPtr unsafe.Pointer) (uintptr, error) {
	var count uintptr
	ret, _, _ := syscall.Syscall(vtbl.GetCount,
		1,
		uintptr(objPtr),
		uintptr(unsafe.Pointer(&count)),
		0)
	if err := hresultToError(ret); err != nil {
		return 0, err
	}
	return count, nil
}

func (vtbl *iShellItemArrayVtbl) getItemAt(objPtr unsafe.Pointer, index uintptr) (string, error) {
	var shellItem *iShellItem
	ret, _, _ := syscall.Syscall(vtbl.GetItemAt,
		2,
		uintptr(objPtr),
		index,
		uintptr(unsafe.Pointer(&shellItem)))
	if err := hresultToError(ret); err != nil {
		return "", err
	}
	if shellItem == nil {
		return "", ErrCancelled
	}
	defer shellItem.vtbl.release(unsafe.Pointer(shellItem))
	return shellItem.vtbl.getDisplayName(unsafe.Pointer(shellItem))
}
