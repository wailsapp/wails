//go:build windows

package application

import (
	"sync"
	"syscall"
	"unsafe"
)

const (
	tdfAllowDialogCancellation = 0x0008
)

const (
	tdcbfOkButton     = 0x0001
	tdcbfCancelButton = 0x0008
)

const (
	tdiInformationIcon = 65534
	tdiErrorIcon       = 65531
	tdiWarningIcon     = 65533
)

type taskDialogButton struct {
	nButtonID     int32
	pszButtonText *uint16
}

type taskDialogConfig struct {
	cbSize               uint32
	hwndParent           uintptr
	hInstance            uintptr
	dwFlags              uint32
	dwCommonButtons      uint32
	pszWindowTitle       *uint16
	_                    uintptr
	pszMainIcon          *uint16
	pszMainInstruction   *uint16
	pszContent           *uint16
	cButtons             uint32
	pButtons             uintptr
	nDefaultButton       int32
	cRadioButtons        uint32
	pRadioButtons        uintptr
	nDefaultRadioButton  int32
	pszVerificationText  *uint16
	pszExpandedInfo      *uint16
	pszExpandedCtrlText  *uint16
	pszCollapsedCtrlText *uint16
	_                    uintptr
	pszFooterIcon        *uint16
	pszFooter            *uint16
	pfCallback           uintptr
	lpCallbackData       uintptr
	cxWidth              uint32
}

var (
	lazyComctl32             = syscall.NewLazyDLL("comctl32.dll")
	procTaskDialogIndirect   = lazyComctl32.NewProc("TaskDialogIndirect")
	taskDialogCallbackMutex  sync.Mutex
	taskDialogButtonCallback map[int32]func()
)

func init() {
	taskDialogButtonCallback = make(map[int32]func())
}

func taskDialogAvailable() bool {
	return procTaskDialogIndirect.Find() == nil
}

func showTaskDialog(dialog *MessageDialog) bool {
	if !taskDialogAvailable() {
		return false
	}

	var parentWindow uintptr
	if dialog.window != nil {
		if nativeWindow := dialog.window.NativeWindow(); nativeWindow != nil {
			parentWindow = uintptr(nativeWindow)
		}
	}

	cfg := taskDialogConfig{
		cbSize:     uint32(unsafe.Sizeof(taskDialogConfig{})),
		hwndParent: parentWindow,
		dwFlags:    tdfAllowDialogCancellation,
	}

	if dialog.Title != "" {
		cfg.pszWindowTitle = syscall.StringToUTF16Ptr(dialog.Title)
	}

	if dialog.Message != "" {
		cfg.pszMainInstruction = syscall.StringToUTF16Ptr(dialog.Message)
	}

	switch dialog.DialogType {
	case InfoDialogType:
		cfg.pszMainIcon = makeIntResource(tdiInformationIcon)
	case ErrorDialogType:
		cfg.pszMainIcon = makeIntResource(tdiErrorIcon)
	case WarningDialogType:
		cfg.pszMainIcon = makeIntResource(tdiWarningIcon)
	case QuestionDialogType:
		cfg.pszMainIcon = makeIntResource(tdiInformationIcon)
	}

	if len(dialog.Buttons) == 0 {
		cfg.dwCommonButtons = tdcbfOkButton
	}

	var buttons []taskDialogButton
	const customButtonBase = 100

	taskDialogCallbackMutex.Lock()
	for id := range taskDialogButtonCallback {
		delete(taskDialogButtonCallback, id)
	}

	for i, btn := range dialog.Buttons {
		id := int32(customButtonBase + i)
		buttons = append(buttons, taskDialogButton{
			nButtonID:     id,
			pszButtonText: syscall.StringToUTF16Ptr(btn.Label),
		})
		if btn.Callback != nil {
			taskDialogButtonCallback[id] = btn.Callback
		}
		if btn.IsDefault {
			cfg.nDefaultButton = id
		}
	}
	taskDialogCallbackMutex.Unlock()

	if len(buttons) > 0 {
		cfg.cButtons = uint32(len(buttons))
		cfg.pButtons = uintptr(unsafe.Pointer(&buttons[0]))
	}

	var buttonPressed int32
	ret, _, _ := procTaskDialogIndirect.Call(
		uintptr(unsafe.Pointer(&cfg)),
		uintptr(unsafe.Pointer(&buttonPressed)),
		0,
		0,
	)

	if ret != 0 {
		return false
	}

	if buttonPressed >= customButtonBase {
		taskDialogCallbackMutex.Lock()
		if cb, ok := taskDialogButtonCallback[buttonPressed]; ok && cb != nil {
			taskDialogCallbackMutex.Unlock()
			cb()
			return true
		}
		taskDialogCallbackMutex.Unlock()
	}

	return true
}

func makeIntResource(id uint16) *uint16 {
	return (*uint16)(unsafe.Pointer(uintptr(id)))
}
