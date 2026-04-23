package application

// WebviewPanelOptions contains options for creating a WebviewPanel.
// Panels are absolutely positioned webview containers within a window,
// similar to Electron's BrowserView or iframes in web development.
//
// Example - Simple panel:
//
//	panel := window.NewPanel(application.WebviewPanelOptions{
//		Name:   "browser",
//		URL:    "https://example.com",
//		X:      0, Y: 50, Width: 800, Height: 600,
//	})
//
// Example - Panel with custom headers and anchoring:
//
//	sidebar := window.NewPanel(application.WebviewPanelOptions{
//		Name:    "api-panel",
//		URL:     "https://api.example.com/dashboard",
//		Headers: map[string]string{"Authorization": "Bearer token123"},
//		X:       0, Y: 0, Width: 200, Height: 600,
//		Anchor:  application.AnchorTop | application.AnchorBottom | application.AnchorLeft,
//	})
type WebviewPanelOptions struct {
	// Name is a unique identifier for the panel within its parent window.
	// Used for retrieving panels via window.GetPanel(name).
	// If empty, a name will be auto-generated (e.g., "panel-1").
	Name string

	// ==================== Content ====================

	// URL is the URL to load in the panel.
	// Can be:
	//   - An external URL (e.g., "https://example.com")
	//   - A local path served by the asset server (e.g., "/panel.html")
	URL string

	// Headers are custom HTTP headers to send with the initial request.
	// These headers are only applied to the initial navigation.
	// Example: {"Authorization": "Bearer token", "X-Custom-Header": "value"}
	Headers map[string]string

	// UserAgent overrides the default user agent string for this panel.
	// If empty, uses the default WebView2/WebKit user agent.
	UserAgent string

	// ==================== Position & Size ====================

	// X is the horizontal position of the panel relative to the parent window's content area.
	// Uses CSS pixels (device-independent).
	X int

	// Y is the vertical position of the panel relative to the parent window's content area.
	// Uses CSS pixels (device-independent).
	Y int

	// Width is the width of the panel in CSS pixels.
	// If 0, defaults to 400.
	Width int

	// Height is the height of the panel in CSS pixels.
	// If 0, defaults to 300.
	Height int

	// ZIndex controls the stacking order of panels within the window.
	// Higher values appear on top of lower values.
	// The main webview has an effective ZIndex of 0.
	// Default: 1
	ZIndex int

	// Anchor specifies how the panel should respond to window resizing.
	// When anchored to an edge, the panel maintains its distance from that edge.
	//
	// Examples:
	//   - AnchorLeft | AnchorTop: Panel stays in top-left corner
	//   - AnchorLeft | AnchorTop | AnchorBottom: Left sidebar that stretches vertically
	//   - AnchorFill: Panel fills the entire window
	//
	// See also: DockLeft(), DockRight(), DockTop(), DockBottom(), FillWindow()
	Anchor AnchorType

	// ==================== Appearance ====================

	// Visible controls whether the panel is initially visible.
	// Default: true
	Visible *bool

	// BackgroundColour is the background color of the panel.
	// Only used when Transparent is false.
	BackgroundColour RGBA

	// Transparent makes the panel background transparent.
	// Useful for overlays or panels with rounded corners.
	// Default: false
	Transparent bool

	// Frameless removes the default styling/border around the panel.
	// Default: false
	Frameless bool

	// Zoom is the initial zoom level of the panel.
	// 1.0 = 100%, 1.5 = 150%, etc.
	// Default: 1.0
	Zoom float64

	// ==================== Developer Options ====================

	// DevToolsEnabled enables the developer tools for this panel.
	// Default: follows the application's debug mode setting
	DevToolsEnabled *bool

	// OpenInspectorOnStartup will open the inspector when the panel is first shown.
	// Only works when DevToolsEnabled is true or app is in debug mode.
	OpenInspectorOnStartup bool
}

// AnchorType defines how a panel is anchored within its parent window.
// Multiple anchors can be combined using bitwise OR.
//
// When a window is resized:
//   - Anchored edges maintain their distance from the window edge
//   - Non-anchored edges allow the panel to stretch/shrink
//
// Example combinations:
//   - AnchorLeft: Panel stays on left, doesn't resize
//   - AnchorLeft | AnchorRight: Panel stretches horizontally with window
//   - AnchorTop | AnchorLeft | AnchorBottom: Left sidebar that stretches vertically
type AnchorType uint8

const (
	// AnchorNone - panel uses absolute positioning only (default)
	// Panel position and size remain fixed regardless of window size changes.
	AnchorNone AnchorType = 0

	// AnchorTop - panel maintains distance from top edge
	AnchorTop AnchorType = 1 << iota

	// AnchorBottom - panel maintains distance from bottom edge
	AnchorBottom

	// AnchorLeft - panel maintains distance from left edge
	AnchorLeft

	// AnchorRight - panel maintains distance from right edge
	AnchorRight

	// AnchorFill - panel fills the entire window (anchored to all edges)
	// Equivalent to: AnchorTop | AnchorBottom | AnchorLeft | AnchorRight
	AnchorFill AnchorType = AnchorTop | AnchorBottom | AnchorLeft | AnchorRight
)

// HasAnchor checks if the anchor type includes a specific anchor.
func (a AnchorType) HasAnchor(anchor AnchorType) bool {
	return a&anchor == anchor
}

// Note: Rect is defined in screenmanager.go
