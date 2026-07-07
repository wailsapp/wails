//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var procDCompositionCreateDevice2 = windows.NewLazySystemDLL("dcomp.dll").NewProc("DCompositionCreateDevice2")

type iDCompositionDeviceVtbl struct {
	_IUnknownVtbl
	Commit                  ComProc
	WaitForCommitCompletion ComProc
	GetFrameStatistics      ComProc
	CreateTargetForHwnd     ComProc
	CreateVisual            ComProc
}

type iDCompositionDevice struct {
	vtbl *iDCompositionDeviceVtbl
}

func (d *iDCompositionDevice) AddRef() uintptr {
	ret, _, _ := d.vtbl.AddRef.Call(uintptr(unsafe.Pointer(d)))

	return ret
}

func (d *iDCompositionDevice) Release() uintptr {
	ret, _, _ := d.vtbl.Release.Call(uintptr(unsafe.Pointer(d)))

	return ret
}

func dCompositionCreateDevice2() (*iDCompositionDevice, error) {
	var device *iDCompositionDevice
	iidIDCompositionDevice := NewGUID("{C37EA93A-E7AA-450D-B16F-9746CB0407F3}")

	hr, _, _ := procDCompositionCreateDevice2.Call(
		0,
		uintptr(unsafe.Pointer(iidIDCompositionDevice)),
		uintptr(unsafe.Pointer(&device)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return device, nil
}

func (d *iDCompositionDevice) CreateTargetForHwnd(hwnd uintptr, topmost bool) (*iDCompositionTarget, error) {
	var target *iDCompositionTarget
	hr, _, _ := d.vtbl.CreateTargetForHwnd.Call(
		uintptr(unsafe.Pointer(d)),
		hwnd,
		uintptr(boolToInt(topmost)),
		uintptr(unsafe.Pointer(&target)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return target, nil
}

func (d *iDCompositionDevice) CreateVisual() (*iDCompositionVisual, error) {
	var visual *iDCompositionVisual
	hr, _, _ := d.vtbl.CreateVisual.Call(
		uintptr(unsafe.Pointer(d)),
		uintptr(unsafe.Pointer(&visual)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return visual, nil
}

func (d *iDCompositionDevice) Commit() error {
	hr, _, _ := d.vtbl.Commit.Call(uintptr(unsafe.Pointer(d)))
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
