package application

import "C"
import (
	"sync"
)

type DialogType int

var dialogID uint
var dialogIDLock sync.RWMutex

func getDialogID() uint {
	dialogIDLock.Lock()
	defer dialogIDLock.Unlock()
	dialogID++
	return dialogID
}

var openFileResponses = make(map[uint]chan string)

const (
	InfoDialog DialogType = iota
	QuestionDialog
	WarningDialog
	ErrorDialog
	OpenDirectoryDialog
)

type Button struct {
	label     string
	isCancel  bool
	isDefault bool
	callback  func()
}

func (b *Button) OnClick(callback func()) {
	b.callback = callback
}

type messageDialogImpl interface {
	show()
}

type MessageDialog struct {
	dialogType DialogType
	title      string
	message    string
	buttons    []*Button
	icon       []byte

	// platform independent
	impl messageDialogImpl
}

var defaultTitles = map[DialogType]string{
	InfoDialog:     "Information",
	QuestionDialog: "Question",
	WarningDialog:  "Warning",
	ErrorDialog:    "Error",
}

func newMessageDialog(dialogType DialogType) *MessageDialog {
	return &MessageDialog{
		dialogType: dialogType,
		title:      defaultTitles[dialogType],
	}
}

func (d *MessageDialog) SetTitle(title string) *MessageDialog {
	d.title = title
	return d
}

func (d *MessageDialog) SetMessage(message string) *MessageDialog {
	d.message = message
	return d
}

func (d *MessageDialog) Show() {
	if d.impl == nil {
		d.impl = newDialogImpl(d)
	}
	d.impl.show()
}

func (d *MessageDialog) SetIcon(icon []byte) *MessageDialog {
	d.icon = icon
	return d
}

func (d *MessageDialog) AddButton(s string) *Button {
	result := &Button{
		label: s,
	}
	d.buttons = append(d.buttons, result)
	return result
}

func (d *MessageDialog) SetDefaultButton(button *Button) *MessageDialog {
	for _, b := range d.buttons {
		b.isDefault = false
	}
	button.isDefault = true
	return d
}

func (d *MessageDialog) SetCancelButton(button *Button) *MessageDialog {
	for _, b := range d.buttons {
		b.isCancel = false
	}
	button.isCancel = true
	return d
}

type openFileDialogImpl interface {
	show() ([]string, error)
}

type OpenFileDialog struct {
	id                      uint
	canChooseDirectories    bool
	canChooseFiles          bool
	canCreateDirectories    bool
	showHiddenFiles         bool
	allowsMultipleSelection bool
	window                  *Window

	impl openFileDialogImpl
}

func (d *OpenFileDialog) CanChooseFiles(canChooseFiles bool) *OpenFileDialog {
	d.canChooseFiles = canChooseFiles
	return d
}

func (d *OpenFileDialog) CanChooseDirectories(canChooseDirectories bool) *OpenFileDialog {
	d.canChooseDirectories = canChooseDirectories
	return d
}

func (d *OpenFileDialog) CanCreateDirectories(canCreateDirectories bool) *OpenFileDialog {
	d.canCreateDirectories = canCreateDirectories
	return d
}

func (d *OpenFileDialog) ShowHiddenFiles(showHiddenFiles bool) *OpenFileDialog {
	d.showHiddenFiles = showHiddenFiles
	return d
}

func (d *OpenFileDialog) AttachToWindow(window *Window) *OpenFileDialog {
	d.window = window
	return d
}

func (d *OpenFileDialog) PromptForSingleFile() (string, error) {
	d.allowsMultipleSelection = false
	if d.impl == nil {
		d.impl = newOpenFileDialogImpl(d)
	}
	selection, err := d.impl.show()
	var result string
	if len(selection) > 0 {
		result = selection[0]
	}

	return result, err
}

func (d *OpenFileDialog) PromptForMultipleFiles() ([]string, error) {
	d.allowsMultipleSelection = true
	if d.impl == nil {
		d.impl = newOpenFileDialogImpl(d)
	}
	return d.impl.show()
}

func newOpenFileDialog() *OpenFileDialog {
	return &OpenFileDialog{
		id:                   getDialogID(),
		canChooseDirectories: false,
		canChooseFiles:       true,
		canCreateDirectories: false,
	}
}
