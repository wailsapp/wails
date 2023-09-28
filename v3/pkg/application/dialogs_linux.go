package application

func (m *linuxApp) showAboutDialog(title string, message string, icon []byte) {
	window := globalApplication.getWindowForID(m.getCurrentWindowID())
	var parent pointer
	if window != nil {
		parent = (window.(*WebviewWindow).impl).(*linuxWebviewWindow).window
	}
	about := newMessageDialog(InfoDialogType)
	about.SetTitle(title).
		SetMessage(message).
		SetIcon(icon)
	runQuestionDialog(
		parent,
		about,
	)
}

type linuxDialog struct {
	dialog *MessageDialog
}

func (m *linuxDialog) show() {
	windowId := getNativeApplication().getCurrentWindowID()
	window := globalApplication.getWindowForID(windowId)
	var parent pointer
	if window != nil {
		parent = (window.(*WebviewWindow).impl).(*linuxWebviewWindow).window
	}

	response := runQuestionDialog(parent, m.dialog)
	if response >= 0 && response < len(m.dialog.Buttons) {
		button := m.dialog.Buttons[response]
		if button.Callback != nil {
			go button.Callback()
		}
	}
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

func (m *linuxOpenFileDialog) show() ([]string, error) {
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

func (m *linuxSaveFileDialog) show() (string, error) {
	return runSaveFileDialog(m.dialog)
}
