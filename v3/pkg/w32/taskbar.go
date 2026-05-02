//go:build windows

package w32

import (
	"syscall"
	"unsafe"
)

var (
	CLSID_TaskbarList = syscall.GUID{Data1: 0x56FDF344, Data2: 0xFD6D, Data3: 0x11D0, Data4: [8]byte{0x95, 0x8A, 0x00, 0x60, 0x97, 0xC9, 0xA0, 0x90}}
	IID_ITaskbarList3 = syscall.GUID{Data1: 0xEA1AFB91, Data2: 0x9E28, Data3: 0x4B86, Data4: [8]byte{0x90, 0xE9, 0x9E, 0x9F, 0x8A, 0x5E, 0xEF, 0xAF}}
)

// ITaskbarList3 interface for Windows taskbar functionality
type ITaskbarList3 struct {
	lpVtbl *taskbarList3Vtbl
}

type taskbarList3Vtbl struct {
	QueryInterface        uintptr
	AddRef                uintptr
	Release               uintptr
	HrInit                uintptr
	AddTab                uintptr
	DeleteTab             uintptr
	ActivateTab           uintptr
	SetActiveAlt          uintptr
	MarkFullscreenWindow  uintptr
	SetProgressValue      uintptr
	SetProgressState      uintptr
	RegisterTab           uintptr
	UnregisterTab         uintptr
	SetTabOrder           uintptr
	SetTabActive          uintptr
	ThumbBarAddButtons    uintptr
	ThumbBarUpdateButtons uintptr
	ThumbBarSetImageList  uintptr
	SetOverlayIcon        uintptr
	SetThumbnailTooltip   uintptr
	SetThumbnailClip      uintptr
}

// NewTaskbarList3 creates a new instance of ITaskbarList3
func NewTaskbarList3() (*ITaskbarList3, error) {
	const COINIT_APARTMENTTHREADED = 0x2

	if hrInit := CoInitializeEx(COINIT_APARTMENTTHREADED); hrInit != 0 && hrInit != 0x1 {
		return nil, syscall.Errno(hrInit)
	}

	var taskbar *ITaskbarList3
	hr := CoCreateInstance(
		&CLSID_TaskbarList,
		CLSCTX_INPROC_SERVER,
		&IID_ITaskbarList3,
		uintptr(unsafe.Pointer(&taskbar)),
	)

	if hr != 0 {
		CoUninitialize()
		return nil, syscall.Errno(hr)
	}

	if r, _, _ := syscall.SyscallN(taskbar.lpVtbl.HrInit, uintptr(unsafe.Pointer(taskbar))); r != 0 {
		syscall.SyscallN(taskbar.lpVtbl.Release, uintptr(unsafe.Pointer(taskbar)))
		CoUninitialize()
		return nil, syscall.Errno(r)
	}

	return taskbar, nil
}

// SetOverlayIcon sets an overlay icon on the taskbar
func (t *ITaskbarList3) SetOverlayIcon(hwnd HWND, hIcon HICON, description *uint16) error {
	ret, _, _ := syscall.SyscallN(
		t.lpVtbl.SetOverlayIcon,
		uintptr(unsafe.Pointer(t)),
		uintptr(hwnd),
		uintptr(hIcon),
		uintptr(unsafe.Pointer(description)),
	)
	if ret != 0 {
		return syscall.Errno(ret)
	}
	return nil
}

// Release releases the ITaskbarList3 interface
func (t *ITaskbarList3) Release() {
	if t != nil {
		syscall.SyscallN(t.lpVtbl.Release, uintptr(unsafe.Pointer(t)))
		CoUninitialize()
	}
}
