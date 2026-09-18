//go:build windows

package combridge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// IUnknownFromPointer cast a generic pointer into a IUnknownImpl pointer
func IUnknownFromPointer(ref unsafe.Pointer) *IUnknownImpl {
	return (*IUnknownImpl)(ref)
}

// IUnknownFromPointer cast native pointer into a IUnknownImpl pointer
func IUnknownFromUintptr(ref uintptr) *IUnknownImpl {
	return IUnknownFromPointer(unsafe.Pointer(ref))
}

type IUnknownVtbl struct {
	queryInterface uintptr
	addRef         uintptr
	release        uintptr
}

func (i *IUnknownVtbl) QueryInterface(this unsafe.Pointer, refiid *windows.GUID, ppvObject **IUnknownImpl) error {
	r, _, _ := syscall.SyscallN(
		i.queryInterface,
		uintptr(this),
		uintptr(unsafe.Pointer(refiid)),
		uintptr(unsafe.Pointer(ppvObject)),
	)

	if r != uintptr(windows.S_OK) {
		return syscall.Errno(r)
	}

	return nil
}

func (i *IUnknownVtbl) AddRef(this unsafe.Pointer) uint32 {
	r, _, _ := syscall.SyscallN(
		i.addRef,
		uintptr(this),
	)
	return uint32(r)
}

func (i *IUnknownVtbl) Release(this unsafe.Pointer) uint32 {
	r, _, _ := syscall.SyscallN(
		i.release,
		uintptr(this),
	)

	return uint32(r)
}

type IUnknownImpl struct {
	vtbl *IUnknownVtbl
}

func (i *IUnknownImpl) QueryInterface(refiid *windows.GUID, ppvObject **IUnknownImpl) error {
	return i.vtbl.QueryInterface(unsafe.Pointer(i), refiid, ppvObject)
}

func (i *IUnknownImpl) AddRef() uint32 {
	return i.vtbl.AddRef(unsafe.Pointer(i))
}

func (i *IUnknownImpl) Release() uint32 {
	return i.vtbl.Release(unsafe.Pointer(i))
}
