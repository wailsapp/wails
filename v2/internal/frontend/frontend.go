package frontend

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// FileFilter defines a filter for dialog boxes
type FileFilter struct {
	DisplayName string // Filter information EG: "Image Files (*.jpg, *.png)"
	Pattern     string // semicolon separated list of extensions, EG: "*.jpg;*.png"
}

// OpenDialogOptions contains the options for the OpenDialogOptions runtime method
type OpenDialogOptions struct {
	DefaultDirectory           string
	DefaultFilename            string
	Title                      string
	Filters                    []FileFilter
	ShowHiddenFiles            bool
	CanCreateDirectories       bool
	ResolvesAliases            bool
	TreatPackagesAsDirectories bool
}

// SaveDialogOptions contains the options for the SaveDialog runtime method
type SaveDialogOptions struct {
	DefaultDirectory           string
	DefaultFilename            string
	Title                      string
	Filters                    []FileFilter
	ShowHiddenFiles            bool
	CanCreateDirectories       bool
	TreatPackagesAsDirectories bool
}

type DialogType string

const (
	InfoDialog     DialogType = "info"
	WarningDialog  DialogType = "warning"
	ErrorDialog    DialogType = "error"
	QuestionDialog DialogType = "question"
)

type Screen struct {
	IsCurrent bool `json:"isCurrent"`
	IsPrimary bool `json:"isPrimary"`

	// Deprecated: Please use Size and PhysicalSize
	Width int `json:"width"`
	// Deprecated: Please use Size and PhysicalSize
	Height int `json:"height"`

	// Size is the size of the screen in logical pixel space, used when setting sizes in Wails
	Size ScreenSize `json:"size"`
	// PhysicalSize is the physical size of the screen in pixels
	PhysicalSize ScreenSize `json:"physicalSize"`
}

type ScreenSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// MessageDialogOptions contains the options for the Message dialogs, EG Info, Warning, etc runtime methods
type MessageDialogOptions struct {
	Type          DialogType
	Title         string
	Message       string
	Buttons       []string
	DefaultButton string
	CancelButton  string
	Icon          []byte
}

type Frontend interface {
	Run(ctx context.Context) error
	RunMainLoop()
	ExecJS(js string)
	Hide()
	Show()
	Quit()

	// Dialog
	OpenFileDialog(dialogOptions OpenDialogOptions) (string, error)
	OpenMultipleFilesDialog(dialogOptions OpenDialogOptions) ([]string, error)
	OpenDirectoryDialog(dialogOptions OpenDialogOptions) (string, error)
	SaveFileDialog(dialogOptions SaveDialogOptions) (string, error)
	MessageDialog(dialogOptions MessageDialogOptions) (string, error)

	// Window
	WindowSetTitle(title string)
	WindowShow()
	WindowHide()
	WindowCenter()
	WindowToggleMaximise()
	WindowMaximise()
	WindowUnmaximise()
	WindowMinimise()
	WindowUnminimise()
	WindowSetAlwaysOnTop(b bool)
	WindowSetPosition(x int, y int)
	WindowGetPosition() (int, int)
	WindowSetSize(width int, height int)
	WindowGetSize() (int, int)
	WindowSetMinSize(width int, height int)
	WindowSetMaxSize(width int, height int)
	WindowFullscreen()
	WindowUnfullscreen()
	WindowSetBackgroundColour(col *options.RGBA)
	WindowReload()
	WindowReloadApp()
	WindowSetSystemDefaultTheme()
	WindowSetLightTheme()
	WindowSetDarkTheme()
	WindowIsMaximised() bool
	WindowIsMinimised() bool
	WindowIsNormal() bool
	WindowIsFullscreen() bool
	WindowClose()
	WindowPrint()

	// Screen
	ScreenGetAll() ([]Screen, error)

	// Menus
	MenuSetApplicationMenu(menu *menu.Menu)
	MenuUpdateApplicationMenu()

	// Events
	Notify(name string, data ...interface{})

	// Browser
	BrowserOpenURL(url string)

	// Clipboard
	ClipboardGetText() (string, error)
	ClipboardSetText(text string) error
}
