package frontend

import (
	"context"
	"time"

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
	Width     int  `json:"width"`
	Height    int  `json:"height"`
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

// LinuxNotificationAction represents a clickable action on a notification
type LinuxNotificationAction struct {

	// Key is the actions's identifier
	Key string

	// Label is shown on the notification
	Label string

	// OnAction handles the activated action's signal.
	OnAction func(notificationID uint32)
}

// LinuxNotificationSound represents a sound to be played when a notification shows up
type LinuxNotificationSound struct {

	// Name is a themeable named sound from the freedesktop.org sound naming specification to play
	// when the notification pops up. An example would be "message-new-instant".
	// Sound names can be found here: http://0pointer.de/public/sound-naming-spec.html
	Name string

	// File to play when the notification pops up. Most notification servers only handle *.wav files.
	File []byte

	// Causes the server to suppress playing any sounds, if it has that ability.
	// This is usually set when the client itself is going to play its own sound.
	Suppress bool
}

// LinuxNotificationUrgency represents the notification's urgency.
type LinuxNotificationUrgency int

// String returns urgency as a string
func (u LinuxNotificationUrgency) String() string {
	switch u {
	case 0:
		return "low"
	case 1:
		return "normal"
	case 2:
		return "critical"
	default:
		return "normal"
	}
}

// String returns urgency as an unsigned integer acceptable for org.freedesktop.Notifications
func (u LinuxNotificationUrgency) Uint() uint {
	if u < 0 {
		return 0
	}
	if u > 2 {
		return 2
	}
	return uint(u)
}

// LinuxNotificationOptions contains options that are specific to linux
type LinuxNotificationOptions struct {

	// Urgency represents the notifications' urgency.
	//
	// For low and normal urgencies, server implementations may display the notifications how they choose.
	// They should, however, have a sane expiration timeout dependent on the urgency level.
	//
	// Critical notifications should not automatically expire, as they are things that the user will most
	// likely want to know about. They should only be closed when the user dismisses them, for example,
	// by clicking on the notification.
	//   - 0 = low
	//   - 1 = normal
	//   - 2 = critical
	Urgency int

	// ReplacesID is used to replace an existing notification.
	ReplacesID uint32

	// Actions to be shown as buttons on the notification.
	//
	// If an action's key is set to "default" the whole notification becomes clickable instead of creating a button.
	// Additional actions will be shown as buttons.
	Actions []LinuxNotificationAction

	// Sound represents a sound to be played when a notification shows up
	Sound *LinuxNotificationSound

	// OnClose handles the closed notification's signal.
	//   - expired
	//   - dismissed-by-user
	//   - activated-by-user
	//   - closed-by-call
	//   - unknown
	//   - other
	OnClose func(notificationID uint32, reason string)

	// OnShow is called when the notification pops up, returns the shown notification's ID
	OnShow func(notificationID uint32)
}

// WindowsNotificationAction represents a Notification action for a notification
type WindowsNotificationAction struct {

	// Type is the action's type
	Type string

	// Label action's button label
	Label string

	// Arguments to be interpreted by the notification server
	Arguments string
}

// WindowsNotificationOptions contains options that are specific to Windows.
type WindowsNotificationOptions struct {
	Actions []WindowsNotificationAction
	Sound   string
}

// MacNotificationAction represents a Notification action for a notification
type MacNotificationAction struct {

	// Label action's button label
	Label string

	// OnAction handles the activated action's signal.
	OnAction func(ActivationType string, ActivationValue string)
}

// MacOptions contains options that are specific to macOS.
type MacNotificationOptions struct {

	// SubTitle The subtitle of the notification.
	SubTitle string

	// Actions to be shown.
	Actions []MacNotificationAction

	// CloseText The notification "Close" button label.
	CloseText string

	// ContentImage is an image to be displayed attached inside the notification.
	ContentImage []byte
}

// NotificationOptions contains the options for desktop notification options
type NotificationOptions struct {

	// AppID identifies the application for the notifications, defaults to the AppID set in the main or "wails" if not set on start
	AppID string

	// AppIcon used to show an icon on the notification. Absolute path to the icon.
	AppIcon []byte

	// Title is a summary for the notification.
	Title string

	// Message is the notification's message.
	Message string

	// Timeout is the duration for how long a notification is shown.
	Timeout time.Duration

	// Close is a channel that clears an active notification
	Close <-chan bool

	// LinuxOptions holds the linux specific options for a notification.
	LinuxOptions *LinuxNotificationOptions

	// WindowsOptions holds the windows specific options for a notification.
	WindowsOptions *WindowsNotificationOptions

	// MacOptions holds the MacOs specific options for a notification.
	MacOptions *MacNotificationOptions
}

type Frontend interface {
	Run(context.Context) error
	RunMainLoop()
	ExecJS(js string)
	Hide()
	Show()
	Quit()
	AppID() string

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

	//Screen
	ScreenGetAll() ([]Screen, error)

	// Menus
	MenuSetApplicationMenu(menu *menu.Menu)
	MenuUpdateApplicationMenu()

	// Events
	Notify(name string, data ...interface{})

	// Browser
	BrowserOpenURL(url string)

	// Notification
	SendNotification(notificationOptions NotificationOptions) error
}
