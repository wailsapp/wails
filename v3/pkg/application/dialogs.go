package application

import (
	"strings"
	"sync"
)

type DialogType int

var dialogMapID = make(map[uint]struct{})
var dialogIDLock sync.RWMutex

func getDialogID() uint {
	dialogIDLock.Lock()
	defer dialogIDLock.Unlock()
	var dialogID uint
	for {
		if _, ok := dialogMapID[dialogID]; !ok {
			dialogMapID[dialogID] = struct{}{}
			break
		}
		dialogID++
		if dialogID == 0 {
			panic("no more dialog IDs")
		}
	}
	return dialogID
}

func freeDialogID(id uint) {
	dialogIDLock.Lock()
	defer dialogIDLock.Unlock()
	delete(dialogMapID, id)
}

var openFileResponses = make(map[uint]chan string)
var saveFileResponses = make(map[uint]chan string)

const (
	InfoDialogType DialogType = iota
	QuestionDialogType
	WarningDialogType
	ErrorDialogType
)

type Button struct {
	Label     string
	IsCancel  bool
	IsDefault bool
	Callback  func()
}

func (b *Button) OnClick(callback func()) *Button {
	b.Callback = callback
	return b
}

func (b *Button) SetAsDefault() *Button {
	b.IsDefault = true
	return b
}

func (b *Button) SetAsCancel() *Button {
	b.IsCancel = true
	return b
}

type messageDialogImpl interface {
	show()
}

type MessageDialogOptions struct {
	DialogType DialogType
	Title      string
	Message    string
	Buttons    []*Button
	Icon       []byte
	window     *WebviewWindow
}

type MessageDialog struct {
	MessageDialogOptions

	// platform independent
	impl messageDialogImpl
}

var defaultTitles = map[DialogType]string{
	InfoDialogType:     "Information",
	QuestionDialogType: "Question",
	WarningDialogType:  "Warning",
	ErrorDialogType:    "Error",
}

func newMessageDialog(dialogType DialogType) *MessageDialog {
	return &MessageDialog{
		MessageDialogOptions: MessageDialogOptions{
			DialogType: dialogType,
		},
		impl: nil,
	}
}

func (d *MessageDialog) SetTitle(title string) *MessageDialog {
	d.Title = title
	return d
}

func (d *MessageDialog) Show() {
	if d.impl == nil {
		d.impl = newDialogImpl(d)
	}
	InvokeSync(d.impl.show)
}

func (d *MessageDialog) SetIcon(icon []byte) *MessageDialog {
	d.Icon = icon
	return d
}

func (d *MessageDialog) AddButton(s string) *Button {
	result := &Button{
		Label: s,
	}
	d.Buttons = append(d.Buttons, result)
	return result
}

func (d *MessageDialog) AddButtons(buttons []*Button) *MessageDialog {
	d.Buttons = buttons
	return d
}

func (d *MessageDialog) AttachToWindow(window Window) *MessageDialog {
	d.window = window.(*WebviewWindow)
	return d
}

func (d *MessageDialog) SetDefaultButton(button *Button) *MessageDialog {
	for _, b := range d.Buttons {
		b.IsDefault = false
	}
	button.IsDefault = true
	return d
}

func (d *MessageDialog) SetCancelButton(button *Button) *MessageDialog {
	for _, b := range d.Buttons {
		b.IsCancel = false
	}
	button.IsCancel = true
	return d
}

func (d *MessageDialog) SetMessage(message string) *MessageDialog {
	d.Message = message
	return d
}

type openFileDialogImpl interface {
	show() (chan string, error)
}

type FileFilter struct {
	DisplayName string // Filter information EG: "Image Files (*.jpg, *.png)"
	Pattern     string // semicolon separated list of extensions, EG: "*.jpg;*.png"
}

type OpenFileDialogOptions struct {
	CanChooseDirectories            bool
	CanChooseFiles                  bool
	CanCreateDirectories            bool
	ShowHiddenFiles                 bool
	ResolvesAliases                 bool
	AllowsMultipleSelection         bool
	HideExtension                   bool
	CanSelectHiddenExtension        bool
	TreatsFilePackagesAsDirectories bool
	AllowsOtherFileTypes            bool
	Filters                         []FileFilter
	Window                          *WebviewWindow

	Title      string
	Message    string
	ButtonText string
	Directory  string
}

type OpenFileDialogStruct struct {
	id                              uint
	canChooseDirectories            bool
	canChooseFiles                  bool
	canCreateDirectories            bool
	showHiddenFiles                 bool
	resolvesAliases                 bool
	allowsMultipleSelection         bool
	hideExtension                   bool
	canSelectHiddenExtension        bool
	treatsFilePackagesAsDirectories bool
	allowsOtherFileTypes            bool
	filters                         []FileFilter

	title      string
	message    string
	buttonText string
	directory  string
	window     *WebviewWindow

	impl openFileDialogImpl
}

func (d *OpenFileDialogStruct) CanChooseFiles(canChooseFiles bool) *OpenFileDialogStruct {
	d.canChooseFiles = canChooseFiles
	return d
}

func (d *OpenFileDialogStruct) CanChooseDirectories(canChooseDirectories bool) *OpenFileDialogStruct {
	d.canChooseDirectories = canChooseDirectories
	return d
}

func (d *OpenFileDialogStruct) CanCreateDirectories(canCreateDirectories bool) *OpenFileDialogStruct {
	d.canCreateDirectories = canCreateDirectories
	return d
}

func (d *OpenFileDialogStruct) AllowsOtherFileTypes(allowsOtherFileTypes bool) *OpenFileDialogStruct {
	d.allowsOtherFileTypes = allowsOtherFileTypes
	return d
}

func (d *OpenFileDialogStruct) ShowHiddenFiles(showHiddenFiles bool) *OpenFileDialogStruct {
	d.showHiddenFiles = showHiddenFiles
	return d
}

func (d *OpenFileDialogStruct) HideExtension(hideExtension bool) *OpenFileDialogStruct {
	d.hideExtension = hideExtension
	return d
}

func (d *OpenFileDialogStruct) TreatsFilePackagesAsDirectories(treatsFilePackagesAsDirectories bool) *OpenFileDialogStruct {
	d.treatsFilePackagesAsDirectories = treatsFilePackagesAsDirectories
	return d
}

func (d *OpenFileDialogStruct) AttachToWindow(window Window) *OpenFileDialogStruct {
	d.window = window.(*WebviewWindow)
	return d
}

func (d *OpenFileDialogStruct) ResolvesAliases(resolvesAliases bool) *OpenFileDialogStruct {
	d.resolvesAliases = resolvesAliases
	return d
}

func (d *OpenFileDialogStruct) SetTitle(title string) *OpenFileDialogStruct {
	d.title = title
	return d
}

func (d *OpenFileDialogStruct) PromptForSingleSelection() (string, error) {
	d.allowsMultipleSelection = false
	if d.impl == nil {
		d.impl = newOpenFileDialogImpl(d)
	}

	var result string
	selections, err := InvokeSyncWithResultAndError(d.impl.show)
	if err == nil {
		result = <-selections
	}

	return result, err
}

// AddFilter adds a filter to the dialog. The filter is a display name and a semicolon separated list of extensions.
// EG: AddFilter("Image Files", "*.jpg;*.png")
func (d *OpenFileDialogStruct) AddFilter(displayName, pattern string) *OpenFileDialogStruct {
	d.filters = append(d.filters, FileFilter{
		DisplayName: strings.TrimSpace(displayName),
		Pattern:     strings.TrimSpace(pattern),
	})
	return d
}

func (d *OpenFileDialogStruct) PromptForMultipleSelection() ([]string, error) {
	d.allowsMultipleSelection = true
	if d.impl == nil {
		d.impl = newOpenFileDialogImpl(d)
	}

	selections, err := InvokeSyncWithResultAndError(d.impl.show)

	var result []string
	for filename := range selections {
		result = append(result, filename)
	}

	return result, err
}

func (d *OpenFileDialogStruct) SetMessage(message string) *OpenFileDialogStruct {
	d.message = message
	return d
}

func (d *OpenFileDialogStruct) SetButtonText(text string) *OpenFileDialogStruct {
	d.buttonText = text
	return d
}

func (d *OpenFileDialogStruct) SetDirectory(directory string) *OpenFileDialogStruct {
	d.directory = directory
	return d
}

func (d *OpenFileDialogStruct) CanSelectHiddenExtension(canSelectHiddenExtension bool) *OpenFileDialogStruct {
	d.canSelectHiddenExtension = canSelectHiddenExtension
	return d
}

func (d *OpenFileDialogStruct) SetOptions(options *OpenFileDialogOptions) {
	d.title = options.Title
	d.message = options.Message
	d.buttonText = options.ButtonText
	d.directory = options.Directory
	d.canChooseDirectories = options.CanChooseDirectories
	d.canChooseFiles = options.CanChooseFiles
	d.canCreateDirectories = options.CanCreateDirectories
	d.showHiddenFiles = options.ShowHiddenFiles
	d.resolvesAliases = options.ResolvesAliases
	d.allowsMultipleSelection = options.AllowsMultipleSelection
	d.hideExtension = options.HideExtension
	d.canSelectHiddenExtension = options.CanSelectHiddenExtension
	d.treatsFilePackagesAsDirectories = options.TreatsFilePackagesAsDirectories
	d.allowsOtherFileTypes = options.AllowsOtherFileTypes
	d.filters = options.Filters
	d.window = options.Window
}

func newOpenFileDialog() *OpenFileDialogStruct {
	return &OpenFileDialogStruct{
		id:                   getDialogID(),
		canChooseDirectories: false,
		canChooseFiles:       true,
		canCreateDirectories: true,
		resolvesAliases:      false,
	}
}

func newSaveFileDialog() *SaveFileDialogStruct {
	return &SaveFileDialogStruct{
		id:                   getDialogID(),
		canCreateDirectories: true,
	}
}

type SaveFileDialogOptions struct {
	CanCreateDirectories            bool
	ShowHiddenFiles                 bool
	CanSelectHiddenExtension        bool
	AllowOtherFileTypes             bool
	HideExtension                   bool
	TreatsFilePackagesAsDirectories bool
	Title                           string
	Message                         string
	Directory                       string
	Filename                        string
	ButtonText                      string
	Filters                         []FileFilter
	Window                          *WebviewWindow
}

type SaveFileDialogStruct struct {
	id                              uint
	canCreateDirectories            bool
	showHiddenFiles                 bool
	canSelectHiddenExtension        bool
	allowOtherFileTypes             bool
	hideExtension                   bool
	treatsFilePackagesAsDirectories bool
	message                         string
	directory                       string
	filename                        string
	buttonText                      string
	filters                         []FileFilter

	window *WebviewWindow

	impl  saveFileDialogImpl
	title string
}

type saveFileDialogImpl interface {
	show() (chan string, error)
}

func (d *SaveFileDialogStruct) SetOptions(options *SaveFileDialogOptions) {
	d.title = options.Title
	d.canCreateDirectories = options.CanCreateDirectories
	d.showHiddenFiles = options.ShowHiddenFiles
	d.canSelectHiddenExtension = options.CanSelectHiddenExtension
	d.allowOtherFileTypes = options.AllowOtherFileTypes
	d.hideExtension = options.HideExtension
	d.treatsFilePackagesAsDirectories = options.TreatsFilePackagesAsDirectories
	d.message = options.Message
	d.directory = options.Directory
	d.filename = options.Filename
	d.buttonText = options.ButtonText
	d.filters = options.Filters
	d.window = options.Window
}

// AddFilter adds a filter to the dialog. The filter is a display name and a semicolon separated list of extensions.
// EG: AddFilter("Image Files", "*.jpg;*.png")
func (d *SaveFileDialogStruct) AddFilter(displayName, pattern string) *SaveFileDialogStruct {
	d.filters = append(d.filters, FileFilter{
		DisplayName: strings.TrimSpace(displayName),
		Pattern:     strings.TrimSpace(pattern),
	})
	return d
}

func (d *SaveFileDialogStruct) CanCreateDirectories(canCreateDirectories bool) *SaveFileDialogStruct {
	d.canCreateDirectories = canCreateDirectories
	return d
}

func (d *SaveFileDialogStruct) CanSelectHiddenExtension(canSelectHiddenExtension bool) *SaveFileDialogStruct {
	d.canSelectHiddenExtension = canSelectHiddenExtension
	return d
}

func (d *SaveFileDialogStruct) ShowHiddenFiles(showHiddenFiles bool) *SaveFileDialogStruct {
	d.showHiddenFiles = showHiddenFiles
	return d
}

func (d *SaveFileDialogStruct) SetMessage(message string) *SaveFileDialogStruct {
	d.message = message
	return d
}

func (d *SaveFileDialogStruct) SetDirectory(directory string) *SaveFileDialogStruct {
	d.directory = directory
	return d
}

func (d *SaveFileDialogStruct) AttachToWindow(window Window) *SaveFileDialogStruct {
	d.window = window.(*WebviewWindow)
	return d
}

func (d *SaveFileDialogStruct) PromptForSingleSelection() (string, error) {
	if d.impl == nil {
		d.impl = newSaveFileDialogImpl(d)
	}

	var result string
	selections, err := InvokeSyncWithResultAndError(d.impl.show)
	if err == nil {
		result = <-selections
	}
	return result, err
}

func (d *SaveFileDialogStruct) SetButtonText(text string) *SaveFileDialogStruct {
	d.buttonText = text
	return d
}

func (d *SaveFileDialogStruct) SetFilename(filename string) *SaveFileDialogStruct {
	d.filename = filename
	return d
}

func (d *SaveFileDialogStruct) AllowsOtherFileTypes(allowOtherFileTypes bool) *SaveFileDialogStruct {
	d.allowOtherFileTypes = allowOtherFileTypes
	return d
}

func (d *SaveFileDialogStruct) HideExtension(hideExtension bool) *SaveFileDialogStruct {
	d.hideExtension = hideExtension
	return d
}

func (d *SaveFileDialogStruct) TreatsFilePackagesAsDirectories(treatsFilePackagesAsDirectories bool) *SaveFileDialogStruct {
	d.treatsFilePackagesAsDirectories = treatsFilePackagesAsDirectories
	return d
}
