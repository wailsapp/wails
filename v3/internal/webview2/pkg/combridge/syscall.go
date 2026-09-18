//go:build windows

package combridge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modkernel32     = windows.NewLazySystemDLL("kernel32.dll")
	procGlobalAlloc = modkernel32.NewProc("GlobalAlloc")
	procGlobalFree  = modkernel32.NewProc("GlobalFree")

	uintptrSize = unsafe.Sizeof(uintptr(0))
)

func allocUintptrObject(size int) (uintptr, []uintptr) {
	v := globalAlloc(uintptr(size) * uintptrSize)
	slice := unsafe.Slice((*uintptr)(unsafe.Pointer(v)), size)
	return v, slice
}

func globalAlloc(dwBytes uintptr) uintptr {
	ret, _, _ := procGlobalAlloc.Call(uintptr(0), dwBytes)
	if ret == 0 {
		panic("globalAlloc failed")
	}

	return ret
}

func globalFree(data uintptr) {
	ret, _, _ := procGlobalFree.Call(data)
	if ret != 0 {
		panic("globalFree failed")
	}
}
