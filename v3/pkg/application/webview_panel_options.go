package application

// WebviewPanelOptions contains options for creating a WebviewPanel.
// Panels are absolutely positioned webview containers within a window.
type WebviewPanelOptions struct {
	// Name is a unique identifier for the panel within its parent window.
	// If empty, a name will be auto-generated.
	Name string

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

	// URL is the initial URL to load in the panel.
	// Can be a local path (e.g., "/panel.html") or external URL (e.g., "https://example.com").
	URL string

	// HTML is the initial HTML content to display in the panel.
	// If both URL and HTML are set, URL takes precedence.
	HTML string

	// JS is JavaScript to execute after the page loads.
	JS string

	// CSS is CSS to inject into the panel.
	CSS string

	// Visible controls whether the panel is initially visible.
	// Default: true
	Visible *bool

	// DevToolsEnabled enables the developer tools for this panel.
	// Default: follows the parent window's setting
	DevToolsEnabled *bool

	// Zoom is the initial zoom level of the panel.
	// Default: 1.0
	Zoom float64

	// BackgroundColour is the background color of the panel.
	BackgroundColour RGBA

	// Frameless removes the default styling/border around the panel.
	// Default: false
	Frameless bool

	// Transparent makes the panel background transparent.
	// Default: false
	Transparent bool

	// Anchor specifies how the panel should be anchored to the window edges.
	// When anchored, the panel maintains its distance from the specified edges
	// when the window is resized.
	Anchor AnchorType

	// OpenInspectorOnStartup will open the inspector when the panel is first shown.
	OpenInspectorOnStartup bool
}

// AnchorType defines how a panel is anchored within its parent window.
// Multiple anchors can be combined using bitwise OR.
type AnchorType uint8

const (
	// AnchorNone - panel uses absolute positioning only (default)
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
	AnchorFill AnchorType = AnchorTop | AnchorBottom | AnchorLeft | AnchorRight
)

// Note: Rect is defined in screenmanager.go
