package windows

type Theme int

type Messages struct {
	InstallationRequired string
	UpdateRequired       string
	MissingRequirements  string
	Webview2NotInstalled string
	Error                string
	FailedToInstall      string
	DownloadPage         string
	PressOKToInstall     string
	ContactAdmin         string
}

const (
	// SystemDefault will use whatever the system theme is. The application will follow system theme changes.
	SystemDefault Theme = 0
	// Dark Mode
	Dark Theme = 1
	// Light Mode
	Light Theme = 2
)

type BackdropType int32

const (
	Auto            BackdropType = 0
	Disable         BackdropType = 1 // None
	MainWindow      BackdropType = 2 // Mica
	TransientWindow BackdropType = 3 // Acrylic
	TabbedWindow    BackdropType = 4 // Tabbed
)

func RGB(r, g, b uint8) int32 {
	var col = int32(b)
	col = col<<8 | int32(g)
	col = col<<8 | int32(r)
	return col
}

// ThemeSettings contains optional colours to use.
// They may be set using the hex values: 0x00BBGGRR
type ThemeSettings struct {
	DarkModeTitleBar           int32
	DarkModeTitleBarInactive   int32
	DarkModeTitleText          int32
	DarkModeTitleTextInactive  int32
	DarkModeBorder             int32
	DarkModeBorderInactive     int32
	LightModeTitleBar          int32
	LightModeTitleBarInactive  int32
	LightModeTitleText         int32
	LightModeTitleTextInactive int32
	LightModeBorder            int32
	LightModeBorderInactive    int32
}

// Options are options specific to Windows
type Options struct {
	WebviewIsTransparent bool
	WindowIsTranslucent  bool
	DisableWindowIcon    bool

	// Disable all window decorations in Frameless mode, which means no "Aero Shadow" and no "Rounded Corner" will be shown.
	// "Rounded Corners" are only available on Windows 11.
	DisableFramelessWindowDecorations bool

	// Path where the WebView2 stores the user data. If empty %APPDATA%\[BinaryName.exe] will be used.
	// If the path is not valid, a messagebox will be displayed with the error and the app will exit with error code.
	WebviewUserDataPath string

	// Path to the directory with WebView2 executables. If empty WebView2 installed in the system will be used.
	WebviewBrowserPath string

	// Dark/Light or System Default Theme
	Theme Theme

	// Custom settings for dark/light mode
	CustomTheme *ThemeSettings

	// Windows 11 22579 minimum
	TranslucencyType BackdropType

	// User messages that can be customised
	Messages *Messages
}

func DefaultMessages() *Messages {
	return &Messages{
		InstallationRequired: "The WebView2 runtime is required. Press Ok to download and install. Note: The installer will download silently so please wait.",
		UpdateRequired:       "The WebView2 runtime needs updating. Press Ok to download and install. Note: The installer will download silently so please wait.",
		MissingRequirements:  "Missing Requirements",
		Webview2NotInstalled: "WebView2 runtime not installed",
		Error:                "Error",
		FailedToInstall:      "The runtime failed to install correctly. Please try again.",
		DownloadPage:         "This application requires the WebView2 runtime. Press OK to open the download page. Minimum version required: ",
		PressOKToInstall:     "Press Ok to install.",
		ContactAdmin:         "The WebView2 runtime is required to run this application. Please contact your system administrator.",
	}
}
