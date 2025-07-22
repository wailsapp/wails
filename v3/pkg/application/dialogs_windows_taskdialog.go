//go:build windows

package application

import (
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// TaskDialog constants
const (
	TDF_ENABLE_HYPERLINKS          = 0x0001
	TDF_USE_HICON_MAIN             = 0x0002
	TDF_USE_HICON_FOOTER           = 0x0004
	TDF_ALLOW_DIALOG_CANCELLATION  = 0x0008
	TDF_USE_COMMAND_LINKS          = 0x0010
	TDF_USE_COMMAND_LINKS_NO_ICON  = 0x0020
	TDF_EXPAND_FOOTER_AREA         = 0x0040
	TDF_EXPANDED_BY_DEFAULT        = 0x0080
	TDF_VERIFICATION_FLAG_CHECKED  = 0x0100
	TDF_SHOW_PROGRESS_BAR          = 0x0200
	TDF_SHOW_MARQUEE_PROGRESS_BAR  = 0x0400
	TDF_CALLBACK_TIMER             = 0x0800
	TDF_POSITION_RELATIVE_TO_WINDOW = 0x1000
	TDF_RTL_LAYOUT                 = 0x2000
	TDF_NO_DEFAULT_RADIO_BUTTON    = 0x4000
	TDF_CAN_BE_MINIMIZED           = 0x8000
	TDF_SIZE_TO_CONTENT            = 0x01000000
)

// TaskDialog button IDs
const (
	IDOK       = 1
	IDCANCEL   = 2
	IDABORT    = 3
	IDRETRY    = 4
	IDIGNORE   = 5
	IDYES      = 6
	IDNO       = 7
	IDCLOSE    = 8
	IDHELP     = 9
	IDCUSTOM   = 100 // Custom button IDs start at 100
)

// TaskDialogConfig structure
type TASKDIALOGCONFIG struct {
	cbSize                     uint32
	hwndParent                 uintptr
	hInstance                  uintptr
	dwFlags                    uint32
	dwCommonButtons            uint32
	pszWindowTitle             *uint16
	hMainIcon                  uintptr
	pszMainIcon                *uint16
	pszMainInstruction         *uint16
	pszContent                 *uint16
	cButtons                   uint32
	pButtons                   uintptr
	nDefaultButton             int32
	cRadioButtons              uint32
	pRadioButtons              uintptr
	nDefaultRadioButton        int32
	pszVerificationText        *uint16
	pszExpandedInformation     *uint16
	pszExpandedControlText     *uint16
	pszCollapsedControlText    *uint16
	hFooterIcon                uintptr
	pszFooterIcon              *uint16
	pszFooter                  *uint16
	pfCallback                 uintptr
	lpCallbackData             uintptr
	cxWidth                    uint32
}

// TASKDIALOG_BUTTON structure
type TASKDIALOG_BUTTON struct {
	nButtonID     int32
	pszButtonText *uint16
}

var (
	comctl32             = windows.NewLazySystemDLL("comctl32.dll")
	procTaskDialogIndirect = comctl32.NewProc("TaskDialogIndirect")
)

// TaskDialog callback storage
type taskDialogCallback struct {
	buttons   []*Button
	callbacks map[int32]func()
}

var (
	taskDialogCallbacks = make(map[uintptr]*taskDialogCallback)
	taskDialogMutex     sync.Mutex
)

// windowsTaskDialog implements custom dialog using TaskDialog API
type windowsTaskDialog struct {
	dialog *MessageDialog
}

// newTaskDialogImpl creates a new TaskDialog implementation
func newTaskDialogImpl(d *MessageDialog) messageDialogImpl {
	// Check if TaskDialog is available (Windows Vista+)
	if err := procTaskDialogIndirect.Find(); err != nil {
		// Fall back to standard MessageBox
		return newDialogImpl(d)
	}
	return &windowsTaskDialog{
		dialog: d,
	}
}

func (m *windowsTaskDialog) show() {
	var parentWindow uintptr
	var err error
	if m.dialog.window != nil {
		parentWindow, err = m.dialog.window.NativeWindowHandle()
		if err != nil {
			globalApplication.handleFatalError(err)
		}
	}

	// Create TaskDialog config
	config := TASKDIALOGCONFIG{
		cbSize:     uint32(unsafe.Sizeof(TASKDIALOGCONFIG{})),
		hwndParent: parentWindow,
		dwFlags:    TDF_ALLOW_DIALOG_CANCELLATION | TDF_SIZE_TO_CONTENT,
	}

	// Set title
	if m.dialog.Title != "" {
		config.pszWindowTitle = syscall.StringToUTF16Ptr(m.dialog.Title)
	} else {
		config.pszWindowTitle = syscall.StringToUTF16Ptr(defaultTitles[m.dialog.DialogType])
	}

	// Set message
	if m.dialog.Message != "" {
		config.pszContent = syscall.StringToUTF16Ptr(m.dialog.Message)
	}

	// Set icon
	switch m.dialog.DialogType {
	case InfoDialogType:
		config.pszMainIcon = makeIntResource(65534) // TD_INFORMATION_ICON
	case ErrorDialogType:
		config.pszMainIcon = makeIntResource(65531) // TD_ERROR_ICON
	case WarningDialogType:
		config.pszMainIcon = makeIntResource(65533) // TD_WARNING_ICON
	case QuestionDialogType:
		config.pszMainIcon = makeIntResource(65534) // TD_INFORMATION_ICON
	}

	// Handle custom icon
	if m.dialog.Icon != nil {
		// For custom icons, we would need to create an HICON from the byte data
		// This is complex and would require additional Windows API calls
		// For now, we'll use the default icons
	}

	// Prepare buttons
	var buttons []TASKDIALOG_BUTTON
	var defaultButtonID int32 = IDOK
	var cancelButtonID int32 = -1
	callbackData := &taskDialogCallback{
		buttons:   m.dialog.Buttons,
		callbacks: make(map[int32]func()),
	}

	if len(m.dialog.Buttons) > 0 {
		for i, btn := range m.dialog.Buttons {
			buttonID := int32(IDCUSTOM + i)
			buttons = append(buttons, TASKDIALOG_BUTTON{
				nButtonID:     buttonID,
				pszButtonText: syscall.StringToUTF16Ptr(btn.Label),
			})
			
			if btn.Callback != nil {
				callbackData.callbacks[buttonID] = btn.Callback
			}
			
			if btn.IsDefault {
				defaultButtonID = buttonID
			}
			if btn.IsCancel {
				cancelButtonID = buttonID
				config.dwFlags |= TDF_ALLOW_DIALOG_CANCELLATION
			}
		}

		config.cButtons = uint32(len(buttons))
		config.pButtons = uintptr(unsafe.Pointer(&buttons[0]))
		config.nDefaultButton = defaultButtonID
	} else {
		// No custom buttons, use standard OK button
		config.dwCommonButtons = 0x0001 // TDCBF_OK_BUTTON
	}

	// Store callback data
	callbackID := uintptr(unsafe.Pointer(callbackData))
	taskDialogMutex.Lock()
	taskDialogCallbacks[callbackID] = callbackData
	taskDialogMutex.Unlock()

	// Show dialog
	var buttonPressed int32
	var radioButtonPressed int32
	var verificationChecked int32

	ret, _, _ := procTaskDialogIndirect.Call(
		uintptr(unsafe.Pointer(&config)),
		uintptr(unsafe.Pointer(&buttonPressed)),
		uintptr(unsafe.Pointer(&radioButtonPressed)),
		uintptr(unsafe.Pointer(&verificationChecked)),
	)

	// Clean up callback data
	taskDialogMutex.Lock()
	delete(taskDialogCallbacks, callbackID)
	taskDialogMutex.Unlock()

	// Handle button callback
	if ret == 0 && buttonPressed >= IDCUSTOM {
		if callback, ok := callbackData.callbacks[buttonPressed]; ok && callback != nil {
			callback()
		}
	} else if ret == 0 && buttonPressed > 0 && buttonPressed < IDCUSTOM {
		// Handle standard button presses
		// Map standard button IDs to button labels for callback matching
		standardButtonMap := map[int32]string{
			IDOK:     "OK",
			IDCANCEL: "Cancel",
			IDYES:    "Yes",
			IDNO:     "No",
			IDRETRY:  "Retry",
			IDABORT:  "Abort",
			IDIGNORE: "Ignore",
		}
		
		if label, ok := standardButtonMap[buttonPressed]; ok {
			for _, btn := range m.dialog.Buttons {
				if btn.Label == label && btn.Callback != nil {
					btn.Callback()
					break
				}
			}
		}
	}
}

// makeIntResource converts an integer to a resource pointer
func makeIntResource(id uint16) *uint16 {
	return (*uint16)(unsafe.Pointer(uintptr(id)))
}

// tryTaskDialog attempts to use TaskDialog API, returns true if successful
func tryTaskDialog(dialog *MessageDialog) bool {
	// Check if TaskDialog is available
	if err := procTaskDialogIndirect.Find(); err != nil {
		return false
	}

	// Use TaskDialog for dialogs with custom buttons
	if len(dialog.Buttons) > 0 {
		impl := newTaskDialogImpl(dialog)
		impl.show()
		return true
	}

	return false
}