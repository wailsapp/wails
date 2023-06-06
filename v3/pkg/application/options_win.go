package application

type WindowsApplicationOptions struct {
	// WndProcInterceptor is a function that will be called for every message sent in the application.
	// Use this to hook into the main message loop. This is useful for handling custom window messages.
	// If `shouldReturn` is `true` then `returnCode` will be returned by the main message loop.
	// If `shouldReturn` is `false` then returnCode will be ignored and the message will be processed by the main message loop.
	WndProcInterceptor func(hwnd uintptr, msg uint32, wParam, lParam uintptr) (returnCode uintptr, shouldReturn bool)

	// DisableQuitOnLastWindowClosed disables the auto quit of the application if the last window has been closed.
	DisableQuitOnLastWindowClosed bool
}

type BackdropType int32

const (
	Auto    BackdropType = 0
	None    BackdropType = 1
	Mica    BackdropType = 2
	Acrylic BackdropType = 3
	Tabbed  BackdropType = 4
)

type WindowsWindow struct {
	// Select the type of translucent backdrop. Requires Windows 11 22621 or later.
	BackdropType BackdropType
	// Disable the icon in the titlebar
	DisableIcon bool
	// Theme. Defaults to SystemDefault which will use whatever the system theme is. The application will follow system theme changes.
	Theme Theme
	// Custom colours for dark/light mode
	CustomTheme *ThemeSettings

	// Disable all window decorations in Frameless mode, which means no "Aero Shadow" and no "Rounded Corner" will be shown.
	// "Rounded Corners" are only available on Windows 11.
	DisableFramelessWindowDecorations bool

	// WindowMask is used to set the window shape. Use a PNG with an alpha channel to create a custom shape.
	WindowMask          []byte
	WindowMaskDraggable bool

	// Path where the WebView2 stores the user data. If empty %APPDATA%\[BinaryName.exe] will be used.
	// If the path is not valid, a messagebox will be displayed with the error and the app will exit with error code.
	WebviewUserDataPath string

	// Path to the directory with WebView2 executables. If empty WebView2 installed in the system will be used.
	WebviewBrowserPath string

	// WebviewGpuIsDisabled is used to enable / disable GPU acceleration for the webview
	WebviewGpuIsDisabled bool

	// ResizeDebounceMS is the amount of time to debounce redraws of webview2
	// when resizing the window
	ResizeDebounceMS uint16

	// Disable the menu bar for this window
	DisableMenu bool
}

type Theme int

const (
	// SystemDefault will use whatever the system theme is. The application will follow system theme changes.
	SystemDefault Theme = 0
	// Dark Mode
	Dark Theme = 1
	// Light Mode
	Light Theme = 2
)

// ThemeSettings defines custom colours to use in dark or light mode.
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
