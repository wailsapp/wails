package application

import "fmt"

func (m *linuxApp) showAboutDialog(title string, message string, icon []byte) {
	window := globalApplication.getWindowForID(m.getCurrentWindowID())
	var parent pointer
	if window != nil {
		parent = window.impl.(*linuxWebviewWindow).window
	}
	about := newMessageDialog(InfoDialog)
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
		parent = window.impl.(*linuxWebviewWindow).window
	}

	response := runQuestionDialog(parent, m.dialog)
	if response >= 0 {
		fmt.Println("Response: ", response)
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
	dialog *OpenFileDialog
}

func newOpenFileDialogImpl(d *OpenFileDialog) *linuxOpenFileDialog {
	return &linuxOpenFileDialog{
		dialog: d,
	}
}

func (m *linuxOpenFileDialog) show() ([]string, error) {
	openFileResponses[m.dialog.id] = make(chan string)
	//	nsWindow := unsafe.Pointer(nil)
	if m.dialog.window != nil {
		// get NSWindow from window
		//nsWindow = m.dialog.window.impl.(*macosWebviewWindow).nsWindow
	}

	// Massage filter patterns into macOS format
	// We iterate all filter patterns, tidy them up and then join them with a semicolon
	// This should produce a single string of extensions like "png;jpg;gif"
	// 	var filterPatterns string
	// if len(m.dialog.filters) > 0 {
	// 	var allPatterns []string
	// 	for _, filter := range m.dialog.filters {
	// 		patternComponents := strings.Split(filter.Pattern, ";")
	// 		for i, component := range patternComponents {
	// 			filterPattern := strings.TrimSpace(component)
	// 			filterPattern = strings.TrimPrefix(filterPattern, "*.")
	// 			patternComponents[i] = filterPattern
	// 		}
	// 		allPatterns = append(allPatterns, strings.Join(patternComponents, ";"))
	// 	}
	// 	filterPatterns = strings.Join(allPatterns, ";")
	// }

	// C.showOpenFileDialog(C.uint(m.dialog.id),
	// 	C.bool(m.dialog.canChooseFiles),
	// 	C.bool(m.dialog.canChooseDirectories),
	// 	C.bool(m.dialog.canCreateDirectories),
	// 	C.bool(m.dialog.showHiddenFiles),
	// 	C.bool(m.dialog.allowsMultipleSelection),
	// 	C.bool(m.dialog.resolvesAliases),
	// 	C.bool(m.dialog.hideExtension),
	// 	C.bool(m.dialog.treatsFilePackagesAsDirectories),
	// 	C.bool(m.dialog.allowsOtherFileTypes),
	// 	toCString(filterPatterns),
	// 	C.uint(len(filterPatterns)),
	// 	toCString(m.dialog.message),
	// 	toCString(m.dialog.directory),
	// 	toCString(m.dialog.buttonText),
	// 	nsWindow)
	var result []string
	for filename := range openFileResponses[m.dialog.id] {
		result = append(result, filename)
	}
	return result, nil
}

type linuxSaveFileDialog struct {
	dialog *SaveFileDialog
}

func newSaveFileDialogImpl(d *SaveFileDialog) *linuxSaveFileDialog {
	return &linuxSaveFileDialog{
		dialog: d,
	}
}

func (m *linuxSaveFileDialog) show() (string, error) {
	saveFileResponses[m.dialog.id] = make(chan string)
	//	nsWindow := unsafe.Pointer(nil)
	if m.dialog.window != nil {
		// get NSWindow from window
		//		nsWindow = m.dialog.window.impl.(*linuxWebviewWindow).nsWindow
	}

	// C.showSaveFileDialog(C.uint(m.dialog.id),
	// 	C.bool(m.dialog.canCreateDirectories),
	// 	C.bool(m.dialog.showHiddenFiles),
	// 	C.bool(m.dialog.canSelectHiddenExtension),
	// 	C.bool(m.dialog.hideExtension),
	// 	C.bool(m.dialog.treatsFilePackagesAsDirectories),
	// 	C.bool(m.dialog.allowOtherFileTypes),
	// 	toCString(m.dialog.message),
	// 	toCString(m.dialog.directory),
	// 	toCString(m.dialog.buttonText),
	// 	toCString(m.dialog.filename),
	// 	nsWindow)
	return <-saveFileResponses[m.dialog.id], nil
}
