//go:build windows

package w32

import (
	"fmt"
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

// CustomTaskDialog displays a task dialog with custom buttons and icons
func CustomTaskDialog(hwnd HWND, title, instruction, content *uint16, buttons []TASKDIALOG_BUTTON, icon uintptr, defaultButton int32) (int32, error) {
	if len(buttons) == 0 {
		return 0, fmt.Errorf("no buttons specified")
	}

	config := TASKDIALOGCONFIG{
		CbSize:               uint32(unsafe.Sizeof(TASKDIALOGCONFIG{})),
		HwndParent:           hwnd,
		DwFlags:              TDF_ALLOW_DIALOG_CANCELLATION | TDF_USE_HICON_MAIN,
		PszWindowTitle:       title,
		PszMainInstruction:   instruction,
		PszContent:          content,
		CButtons:            uint32(len(buttons)),
		PButtons:            &buttons[0],
		NDefaultButton:      defaultButton,
		HMainIcon:           icon,
	}

	var buttonPressed int32
	err := TaskDialogIndirect(&config, &buttonPressed, nil, nil)
	if err != nil {
		return 0, err
	}

	return buttonPressed, nil
}
