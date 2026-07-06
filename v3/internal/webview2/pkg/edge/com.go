//go:build windows

package edge

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

type EventRegistrationToken struct {
	value int64
}

// IUnknown
type IUnknown struct {
	Vtbl *IUnknownVtbl
}

type IUnknownVtbl struct {
	QueryInterface ComProc
	AddRef         ComProc
	Release        ComProc
}

func (i *IUnknownVtbl) CallRelease(this unsafe.Pointer) uint32 {
	ret, _, _ := i.Release.Call(
		uintptr(this),
	)

	return uint32(ret)
}

type IUnknownImpl interface {
	QueryInterface(refiid, object uintptr) uintptr
	AddRef() uintptr
	Release() uintptr
}

type POINT struct {
	X, Y int32
}
type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}
type HANDLE uintptr
type HBRUSH uintptr
type HCURSOR uintptr
type HICON uintptr
type HINSTANCE uintptr
type HMENU uintptr
type HMODULE uintptr
type HWND uintptr

// NOTE: For sure, this is wrong!
type VARIANT uintptr

type IDataObject struct {
	IUnknown
}

func ptr[T any](p T) *T {
	return &p
}

const ERROR_SUCCESS = windows.ERROR_SUCCESS

func UTF16PtrFromString(s string) (*uint16, error) {
	return windows.UTF16PtrFromString(s)
}

func UTF16PtrToString(s *uint16) string {
	return windows.UTF16PtrToString(s)
}

func CoTaskMemFree(pv unsafe.Pointer) {
	windows.CoTaskMemFree(pv)
}
