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
	InvalidFixedWebview2 string
	WebView2ProcessCrash string
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
	Auto    BackdropType = 0
	None    BackdropType = 1
	Mica    BackdropType = 2
	Acrylic BackdropType = 3
	Tabbed  BackdropType = 4
)

func RGB(r, g, b uint8) int32 {
	col := int32(b)
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

	IsZoomControlEnabled bool
	ZoomFactor           float64

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

	// Select the type of translucent backdrop. Requires Windows 11 22621 or later.
	BackdropType BackdropType

	// User messages that can be customised
	Messages *Messages

	// ResizeDebounceMS is the amount of time to debounce redraws of webview2
	// when resizing the window
	ResizeDebounceMS uint16

	// OnSuspend is called when Windows enters low power mode
	OnSuspend func()

	// OnResume is called when Windows resumes from low power mode
	OnResume func()

	// WebviewGpuIsDisabled is used to enable / disable GPU acceleration for the webview
	WebviewGpuIsDisabled bool

	// WebviewDisableRendererCodeIntegrity disables the `RendererCodeIntegrity` of WebView2. Some Security Endpoint
	// Protection Software inject themself into the WebView2 with unsigned or wrongly signed dlls, which is not allowed
	// and will stop the WebView2 processes. Those security software need an update to fix this issue or one can disable
	// the integrity check with this flag.
	//
	// The event viewer log contains `Code Integrity Errors` like mentioned here: https://github.com/MicrosoftEdge/WebView2Feedback/issues/2051
	//
	// !! Please keep in mind when disabling this feature, this also allows malicious software to inject into the WebView2 !!
	WebviewDisableRendererCodeIntegrity bool

	// Configure whether swipe gestures should be enabled
	EnableSwipeGestures bool
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
		InvalidFixedWebview2: "The WebView2 runtime is manually specified, but It is not valid. Check minimum required version and webview2 path.",
		WebView2ProcessCrash: "The WebView2 process crashed and the application needs to be restarted.",
	}
}
