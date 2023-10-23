package application

import (
	"github.com/wailsapp/go-webview2/pkg/edge"
	"github.com/wailsapp/wails/v3/pkg/events"
)

type BackdropType int32
type DragEffect int32

const (
	// DragEffectNone is used to indicate that the drop target cannot accept the data.
	DragEffectNone DragEffect = 1
	// DragEffectCopy is used to indicate that the data is copied to the drop target.
	DragEffectCopy DragEffect = 2
	// DragEffectMove is used to indicate that the data is removed from the drag source.
	DragEffectMove DragEffect = 3
	// DragEffectLink is used to indicate that a link to the original data is established.
	DragEffectLink DragEffect = 4
	// DragEffectScroll is used to indicate that the target can be scrolled while dragging to locate a drop position that is not currently visible in the target.

)

const (
	Auto    BackdropType = 0
	None    BackdropType = 1
	Mica    BackdropType = 2
	Acrylic BackdropType = 3
	Tabbed  BackdropType = 4
)

type WindowsWindow struct {
	// Select the type of translucent backdrop. Requires Windows 11 22621 or later.
	// Only used when window's `BackgroundType` is set to `BackgroundTypeTranslucent`.
	// Default: Auto
	BackdropType BackdropType

	// Disable the icon in the titlebar
	// Default: false
	DisableIcon bool

	// Theme (Dark / Light / SystemDefault)
	// Default: SystemDefault - The application will follow system theme changes.
	Theme Theme

	// Specify custom colours to use for dark/light mode
	// Default: nil
	CustomTheme *ThemeSettings

	// Disable all window decorations in Frameless mode, which means no "Aero Shadow" and no "Rounded Corner" will be shown.
	// "Rounded Corners" are only available on Windows 11.
	// Default: false
	DisableFramelessWindowDecorations bool

	// WindowMask is used to set the window shape. Use a PNG with an alpha channel to create a custom shape.
	// Default: nil
	WindowMask []byte

	// WindowMaskDraggable is used to make the window draggable by clicking on the window mask.
	// Default: false
	WindowMaskDraggable bool

	// WebviewGpuIsDisabled is used to enable / disable GPU acceleration for the webview
	// Default: false
	WebviewGpuIsDisabled bool

	// ResizeDebounceMS is the amount of time to debounce redraws of webview2
	// when resizing the window
	// Default: 0
	ResizeDebounceMS uint16

	// Disable the menu bar for this window
	// Default: false
	DisableMenu bool

	// Event mapping for the window. This allows you to define a translation from one event to another.
	// Default: nil
	EventMapping map[events.WindowEventType]events.WindowEventType

	// HiddenOnTaskbar hides the window from the taskbar
	// Default: false
	HiddenOnTaskbar bool

	// EnableSwipeGestures enables swipe gestures for the window
	// Default: false
	EnableSwipeGestures bool

	// EnableFraudulentWebsiteWarnings will enable warnings for fraudulent websites.
	// Default: false
	EnableFraudulentWebsiteWarnings bool

	// Menu is the menu to use for the window.
	Menu *Menu

	// Drag Cursor Effects
	OnEnterEffect DragEffect
	OnOverEffect  DragEffect

	// Permissions map for WebView2. If empty, default permissions will be granted.
	Permissions map[edge.CoreWebView2PermissionKind]edge.CoreWebView2PermissionState
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
