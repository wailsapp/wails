//go:build windows

package application

func (m *windowsApp) showAboutDialog(title string, message string, icon []byte) {
	panic("implement me")
}

type windowsDialog struct {
	dialog *MessageDialog

	//dialogImpl unsafe.Pointer
}

func (m *windowsDialog) show() {
	globalApplication.dispatchOnMainThread(func() {
		//
		//// Mac can only have 4 Buttons on a dialog
		//if len(m.dialog.Buttons) > 4 {
		//	m.dialog.Buttons = m.dialog.Buttons[:4]
		//}
		//
		//if m.nsDialog != nil {
		//	C.releaseDialog(m.nsDialog)
		//}
		//var title *C.char
		//if m.dialog.Title != "" {
		//	title = C.CString(m.dialog.Title)
		//}
		//var message *C.char
		//if m.dialog.Message != "" {
		//	message = C.CString(m.dialog.Message)
		//}
		//var iconData unsafe.Pointer
		//var iconLength C.int
		//if m.dialog.Icon != nil {
		//	iconData = unsafe.Pointer(&m.dialog.Icon[0])
		//	iconLength = C.int(len(m.dialog.Icon))
		//} else {
		//	// if it's an error, use the application Icon
		//	if m.dialog.DialogType == ErrorDialog {
		//		iconData = unsafe.Pointer(&globalApplication.options.Icon[0])
		//		iconLength = C.int(len(globalApplication.options.Icon))
		//	}
		//}
		//
		//alertType, ok := alertTypeMap[m.dialog.DialogType]
		//if !ok {
		//	alertType = C.NSAlertStyleInformational
		//}
		//
		//m.nsDialog = C.createAlert(alertType, title, message, iconData, iconLength)
		//
		//// Reverse the Buttons so that the default is on the right
		//reversedButtons := make([]*Button, len(m.dialog.Buttons))
		//var count = 0
		//for i := len(m.dialog.Buttons) - 1; i >= 0; i-- {
		//	button := m.dialog.Buttons[i]
		//	C.alertAddButton(m.nsDialog, C.CString(button.Label), C.bool(button.IsDefault), C.bool(button.IsCancel))
		//	reversedButtons[count] = m.dialog.Buttons[i]
		//	count++
		//}
		//
		//buttonPressed := int(C.dialogRunModal(m.nsDialog))
		//if len(m.dialog.Buttons) > buttonPressed {
		//	button := reversedButtons[buttonPressed]
		//	if button.callback != nil {
		//		button.callback()
		//	}
		//}
		panic("implement me")
	})

}

func newDialogImpl(d *MessageDialog) *windowsDialog {
	return &windowsDialog{
		dialog: d,
	}
}

type windowOpenFileDialog struct {
	dialog *OpenFileDialog
}

func newOpenFileDialogImpl(d *OpenFileDialog) *windowOpenFileDialog {
	return &windowOpenFileDialog{
		dialog: d,
	}
}

func (m *windowOpenFileDialog) show() ([]string, error) {
	//openFileResponses[m.dialog.id] = make(chan string)
	//nsWindow := unsafe.Pointer(nil)
	//if m.dialog.window != nil {
	//	// get NSWindow from window
	//	nsWindow = m.dialog.window.impl.(*windowsWebviewWindow).nsWindow
	//}
	//
	//// Massage filter patterns into macOS format
	//// We iterate all filter patterns, tidy them up and then join them with a semicolon
	//// This should produce a single string of extensions like "png;jpg;gif"
	//var filterPatterns string
	//if len(m.dialog.filters) > 0 {
	//	var allPatterns []string
	//	for _, filter := range m.dialog.filters {
	//		patternComponents := strings.Split(filter.Pattern, ";")
	//		for i, component := range patternComponents {
	//			filterPattern := strings.TrimSpace(component)
	//			filterPattern = strings.TrimPrefix(filterPattern, "*.")
	//			patternComponents[i] = filterPattern
	//		}
	//		allPatterns = append(allPatterns, strings.Join(patternComponents, ";"))
	//	}
	//	filterPatterns = strings.Join(allPatterns, ";")
	//}
	//
	//C.showOpenFileDialog(C.uint(m.dialog.id),
	//	C.bool(m.dialog.canChooseFiles),
	//	C.bool(m.dialog.canChooseDirectories),
	//	C.bool(m.dialog.canCreateDirectories),
	//	C.bool(m.dialog.showHiddenFiles),
	//	C.bool(m.dialog.allowsMultipleSelection),
	//	C.bool(m.dialog.resolvesAliases),
	//	C.bool(m.dialog.hideExtension),
	//	C.bool(m.dialog.treatsFilePackagesAsDirectories),
	//	C.bool(m.dialog.allowsOtherFileTypes),
	//	toCString(filterPatterns),
	//	C.uint(len(filterPatterns)),
	//	toCString(m.dialog.message),
	//	toCString(m.dialog.directory),
	//	toCString(m.dialog.buttonText),
	//	nsWindow)
	//var result []string
	//for filename := range openFileResponses[m.dialog.id] {
	//	result = append(result, filename)
	//}
	//return result, nil
	panic("implement me")
}

type windowSaveFileDialog struct {
	dialog *SaveFileDialog
}

func newSaveFileDialogImpl(d *SaveFileDialog) *windowSaveFileDialog {
	return &windowSaveFileDialog{
		dialog: d,
	}
}

func (m *windowSaveFileDialog) show() (string, error) {
	//saveFileResponses[m.dialog.id] = make(chan string)
	//nsWindow := unsafe.Pointer(nil)
	//if m.dialog.window != nil {
	//	// get NSWindow from window
	//	nsWindow = m.dialog.window.impl.(*macosWebviewWindow).nsWindow
	//}
	//C.showSaveFileDialog(C.uint(m.dialog.id),
	//	C.bool(m.dialog.canCreateDirectories),
	//	C.bool(m.dialog.showHiddenFiles),
	//	C.bool(m.dialog.canSelectHiddenExtension),
	//	C.bool(m.dialog.hideExtension),
	//	C.bool(m.dialog.treatsFilePackagesAsDirectories),
	//	C.bool(m.dialog.allowOtherFileTypes),
	//	toCString(m.dialog.message),
	//	toCString(m.dialog.directory),
	//	toCString(m.dialog.buttonText),
	//	toCString(m.dialog.filename),
	//	nsWindow)
	//return <-saveFileResponses[m.dialog.id], nil
	panic("implement me")
}
