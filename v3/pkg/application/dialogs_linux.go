package application

func (a *linuxApp) showAboutDialog(title string, message string, icon []byte) {
	window := globalApplication.getWindowForID(a.getCurrentWindowID())
	var parent uintptr
	if window != nil {
		parent, _ = window.(*WebviewWindow).NativeWindowHandle()
	}
	about := newMessageDialog(InfoDialogType)
	about.SetTitle(title).
		SetMessage(message).
		SetIcon(icon)
	InvokeAsync(func() {
		runQuestionDialog(
			pointer(parent),
			about,
		)
	})
}

type linuxDialog struct {
	dialog *MessageDialog
}

func (m *linuxDialog) show() {
	windowId := getNativeApplication().getCurrentWindowID()
	window := globalApplication.getWindowForID(windowId)
	var parent uintptr
	if window != nil {
		parent, _ = window.(*WebviewWindow).NativeWindowHandle()
	}

	InvokeAsync(func() {
		response := runQuestionDialog(pointer(parent), m.dialog)
		if response >= 0 && response < len(m.dialog.Buttons) {
			button := m.dialog.Buttons[response]
			if button.Callback != nil {
				go func() {
					defer handlePanic()
					button.Callback()
				}()
			}
		}
	})
}

func newDialogImpl(d *MessageDialog) *linuxDialog {
	return &linuxDialog{
		dialog: d,
	}
}

type linuxOpenFileDialog struct {
	dialog *OpenFileDialogStruct
}

func newOpenFileDialogImpl(d *OpenFileDialogStruct) *linuxOpenFileDialog {
	return &linuxOpenFileDialog{
		dialog: d,
	}
}

func (m *linuxOpenFileDialog) show() (chan string, error) {
	return runOpenFileDialog(m.dialog)
}

type linuxSaveFileDialog struct {
	dialog *SaveFileDialogStruct
}

func newSaveFileDialogImpl(d *SaveFileDialogStruct) *linuxSaveFileDialog {
	return &linuxSaveFileDialog{
		dialog: d,
	}
}

func (m *linuxSaveFileDialog) show() (chan string, error) {
	return runSaveFileDialog(m.dialog)
}
