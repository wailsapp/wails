//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type iDCompositionTargetVtbl struct {
	_IUnknownVtbl
	SetRoot ComProc
}

type iDCompositionTarget struct {
	vtbl *iDCompositionTargetVtbl
}

func (t *iDCompositionTarget) AddRef() uintptr {
	ret, _, _ := t.vtbl.AddRef.Call(uintptr(unsafe.Pointer(t)))

	return ret
}

func (t *iDCompositionTarget) Release() uintptr {
	ret, _, _ := t.vtbl.Release.Call(uintptr(unsafe.Pointer(t)))

	return ret
}

func (t *iDCompositionTarget) SetRoot(visual *iDCompositionVisual) error {
	hr, _, _ := t.vtbl.SetRoot.Call(
		uintptr(unsafe.Pointer(t)),
		uintptr(unsafe.Pointer(visual)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
