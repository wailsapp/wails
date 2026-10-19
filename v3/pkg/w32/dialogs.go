//go:build windows

package w32

import (
	"syscall"
	"unsafe"
)

func MessageBoxWithIcon(hwnd HWND, text *uint16, caption *uint16, iconID int, flags uint32) (int32, error) {

	params := MSGBOXPARAMS{
		cbSize:      uint32(unsafe.Sizeof(MSGBOXPARAMS{})),
		hwndOwner:   hwnd,
		hInstance:   GetApplicationHandle(),
		lpszText:    text,
		lpszCaption: caption,
		dwStyle:     flags,
		lpszIcon:    (*uint16)(unsafe.Pointer(uintptr(iconID))),
	}

	r, _, err := procMessageBoxIndirect.Call(
		uintptr(unsafe.Pointer(&params)),
	)
	if r == 0 {
		return 0, err
	}
	return int32(r), nil
}

var procTaskDialogIndirect = modcomctl32.NewProc("TaskDialogIndirect")

// TaskDialogAvailable reports whether TaskDialogIndirect can be called. It is
// only exported by comctl32.dll v6, which requires the application manifest to
// declare a dependency on Microsoft.Windows.Common-Controls 6.0.0.0.
func TaskDialogAvailable() bool {
	return procTaskDialogIndirect.Find() == nil
}

// TaskDialogIndirect shows the task dialog described by config, a serialised
// (1-byte packed) TASKDIALOGCONFIG, and returns the ID of the pressed button.
func TaskDialogIndirect(config []byte) (int32, error) {
	if err := procTaskDialogIndirect.Find(); err != nil {
		return 0, err
	}
	var button int32
	hr, _, _ := procTaskDialogIndirect.Call(
		uintptr(unsafe.Pointer(&config[0])),
		uintptr(unsafe.Pointer(&button)),
		0,
		0,
	)
	if int32(hr) < 0 {
		return 0, syscall.Errno(uint32(hr))
	}
	return button, nil
}
