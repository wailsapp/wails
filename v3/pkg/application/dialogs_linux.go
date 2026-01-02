//go:build linux && !android

package application

func (a *linuxApp) showAboutDialog(title string, message string, icon []byte) {
	window, _ := globalApplication.Window.GetByID(a.getCurrentWindowID())
	var parent uintptr
	if window != nil {
		nativeWindow := window.NativeWindow()
		if nativeWindow != nil {
			parent = uintptr(nativeWindow)
		}
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

func (m *linuxDialog) show() (string, error) {
	windowId := getNativeApplication().getCurrentWindowID()
	window, _ := globalApplication.Window.GetByID(windowId)
	var parent uintptr
	if window != nil {
		nativeWindow := window.NativeWindow()
		if nativeWindow != nil {
			parent = uintptr(nativeWindow)
		}
	}

	// Channel to receive the result
	resultChan := make(chan string, 1)

	InvokeAsync(func() {
		response := runQuestionDialog(pointer(parent), m.dialog)
		var buttonLabel string
		if response >= 0 && response < len(m.dialog.Buttons) {
			button := m.dialog.Buttons[response]
			buttonLabel = button.Label
			if button.Callback != nil {
				go func() {
					defer handlePanic()
					button.Callback()
				}()
			}
		}
		resultChan <- buttonLabel
	})

	// Wait for and return the result
	return <-resultChan, nil
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
