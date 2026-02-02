//go:build windows

package w32

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

// ComProc stores a COM procedure.
type ComProc uintptr

// NewComProc creates a new COM proc from a Go function.
func NewComProc(fn interface{}) ComProc {
	return ComProc(windows.NewCallback(fn))
}

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

func (i *IUnknownVtbl) CallRelease(this unsafe.Pointer) error {
	_, _, err := i.Release.Call(
		uintptr(this),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

type IUnknownImpl interface {
	QueryInterface(refiid, object uintptr) uintptr
	AddRef() uintptr
	Release() uintptr
}

// Call calls a COM procedure.
//
//go:uintptrescapes
func (p ComProc) Call(a ...uintptr) (r1, r2 uintptr, lastErr error) {
	return syscall.SyscallN(uintptr(p), a...)
}
