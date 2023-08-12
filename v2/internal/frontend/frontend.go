package frontend

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// FileFilter defines a filter for dialog boxes
type FileFilter struct {
	DisplayName string `json:"displayName" ` // Filter information EG: "Image Files (*.jpg, *.png)"
	Pattern     string `json:"pattern"`      // semicolon separated list of extensions, EG: "*.jpg;*.png"
}

// OpenDialogOptions contains the options for the OpenDialogOptions runtime method
type OpenDialogOptions struct {
	DefaultDirectory           string       `json:"defaultDirectory"`
	DefaultFilename            string       `json:"defaultFilename"`
	Title                      string       `json:"title"`
	Filters                    []FileFilter `json:"filters"`
	ShowHiddenFiles            bool         `json:"showHiddenFiles"`
	CanCreateDirectories       bool         `json:"canCreateDirectories"`
	ResolvesAliases            bool         `json:"resolvesAliases"`
	TreatPackagesAsDirectories bool         `json:"treatPackagesAsDirectories"`
}

// SaveDialogOptions contains the options for the SaveDialog runtime method
type SaveDialogOptions struct {
	DefaultDirectory           string       `json:"defaultDirectory"`
	DefaultFilename            string       `json:"defaultFilename"`
	Title                      string       `json:"title"`
	Filters                    []FileFilter `json:"filters"`
	ShowHiddenFiles            bool         `json:"showHiddenFiles"`
	CanCreateDirectories       bool         `json:"canCreateDirectories"`
	TreatPackagesAsDirectories bool         `json:"treatPackagesAsDirectories"`
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
	Width     int  `json:"width"`
	Height    int  `json:"height"`
}

// MessageDialogOptions contains the options for the Message dialogs, EG Info, Warning, etc runtime methods
type MessageDialogOptions struct {
	Type          DialogType `json:"type"`
	Title         string     `json:"title"`
	Message       string     `json:"message"`
	Buttons       []string   `json:"buttons"`
	DefaultButton string     `json:"defaultButton"`
	CancelButton  string     `json:"cancelButton"`
	Icon          []byte     `json:"icon"`
}

type Frontend interface {
	Run(context.Context) error
	RunMainLoop()
	ExecJS(js string)
	Hide()
	Show()
	Quit()

	// Dialog
	OpenFileDialog(dialogOptions OpenDialogOptions) (string, error)
	OpenMultipleFilesDialog(dialogOptions OpenDialogOptions) ([]string, error)
	OpenDirectoryDialog(dialogOptions OpenDialogOptions) (string, error)
	OpenMultipleDirectoriesDialog(dialogOptions OpenDialogOptions) ([]string, error)
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

	//Screen
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
