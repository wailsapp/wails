package edge

import "github.com/wailsapp/wails/webview2/webviewloader"

type Capability string

var UnsupportedCapabilityError = &unsupportedCapabilityError{}

type unsupportedCapabilityError struct{}

func (u *unsupportedCapabilityError) Error() string {
	return "unsupported capability"
}

// Capabilities is a list of capabilities with their corresponding minimum runtime version
// Internal Capabilities are not exposed to the user
// Larger capabilities such as DragAndDrop should be exported with a capital letter

// WebView2 Runtime Version 131.0.2903.40 (Released: September 2023)
const (
	ScreenCapture       = Capability("131.0.2903.40") // Screen capture support
	NonClientRegion     = Capability("131.0.2903.40") // Non-client region customization
	DownloadDialog      = Capability("131.0.2903.40") // Download dialog handling
	BrowserExtension    = Capability("131.0.2903.40") // Browser extension support
	BasicAuthentication = Capability("131.0.2903.40") // Basic authentication handling
	SaveFileDialog      = Capability("131.0.2903.40") // Save file dialog support
)

// WebView2 Runtime Version 113.0.1774.30 (Released: Unknown)
const (
	GetAdditionalObjects = Capability("113.0.1774.30") // Additional objects support
)

// WebView2 Runtime Version 100.0.1185.39 (Released: April 2022)
const (
	AllowExternalDrop       = Capability("100.0.1185.39") // External drop support
	GeneralAutofillEnabled  = Capability("100.0.1185.39") // General autofill features
	PasswordAutosaveEnabled = Capability("100.0.1185.39") // Password autosave support
)

// WebView2 Runtime Version 98.0.1108.43 (Released: February 2022)
const (
	CustomScheme      = Capability("98.0.1108.43") // Custom scheme support
	PrintToPdf        = Capability("98.0.1108.43") // Print to PDF functionality
	SharedBuffer      = Capability("98.0.1108.43") // Shared buffer for performance
	ServerCertificate = Capability("98.0.1108.43") // Server certificate handling
	FrameNavigation   = Capability("98.0.1108.43") // Frame navigation events
)

// WebView2 Runtime Version 97.0.1072.69 (Released: January 2022)
const (
	ClientCertificate = Capability("97.0.1072.69") // Client certificate selection
	ContextMenus      = Capability("97.0.1072.69") // Custom context menus
	BackgroundColor   = Capability("97.0.1072.69") // Background color customization
	ScriptEnabled     = Capability("97.0.1072.69") // JavaScript execution control
	StatusBar         = Capability("97.0.1072.69") // Status bar customization
)

// WebView2 Runtime Version 95.0.1020.44 (Released: October 2021)
const (
	WebMessageReceived   = Capability("95.0.1020.44") // Web message handling
	NewWindowRequested   = Capability("95.0.1020.44") // New window request handling
	DocumentTitleChanged = Capability("95.0.1020.44") // Document title change events
	ContainsFullScreen   = Capability("95.0.1020.44") // Fullscreen mode detection
	WebResourceRequested = Capability("95.0.1020.44") // Web resource request handling
)

// WebView2 Runtime Version 94.0.992.31 (Released: September 2021)
const (
	NavigationStarting      = Capability("94.0.992.31") // Navigation start events
	NavigationCompleted     = Capability("94.0.992.31") // Navigation completion events
	FrameNavigationStarting = Capability("94.0.992.31") // Frame navigation start
	SourceChanged           = Capability("94.0.992.31") // Source change detection
	HistoryChanged          = Capability("94.0.992.31") // Browser history changes
	SwipeNavigation         = Capability("94.0.992.31") // Swipe navigation support
)

// WebView2 Runtime Version 93.0.961.52 (Released: August 2021)
const (
	DOMContentLoaded    = Capability("93.0.961.52") // DOM content loaded events
	WebResourceLoaded   = Capability("93.0.961.52") // Resource load events
	ScriptDialogOpening = Capability("93.0.961.52") // Script dialog handling
	PermissionRequested = Capability("93.0.961.52") // Permission request handling
	ProcessFailed       = Capability("93.0.961.52") // Process failure detection
)

// WebView2 Runtime Version 92.0.902.78 (Released: July 2021)
const (
	AcceleratorKeyPressed = Capability("92.0.902.78") // Accelerator key handling
	ZoomFactorChanged     = Capability("92.0.902.78") // Zoom factor change events
	MoveFocusRequested    = Capability("92.0.902.78") // Focus movement handling
	DevToolsProtocol      = Capability("92.0.902.78") // DevTools protocol support
	BrowserProcessExited  = Capability("92.0.902.78") // Browser process exit handling
)

// WebView2 Runtime Version 91.0.864.41 (Released: June 2021)
const (
	DefaultDownloadDialog    = Capability("91.0.864.41") // Default download dialog
	DefaultContextMenus      = Capability("91.0.864.41") // Default context menus
	FaviconChanged           = Capability("91.0.864.41") // Favicon change events
	WindowCloseRequested     = Capability("91.0.864.41") // Window close request handling
	RasterizationScale       = Capability("91.0.864.41") // Display scaling support
	SecurityUpdated          = Capability("91.0.864.41") // Security state updates
	ProcessInfoReceived      = Capability("91.0.864.41") // Process info events
	FramePermissionRequested = Capability("91.0.864.41") // Frame permission requests
	ClearBrowsingData        = Capability("91.0.864.41") // Clear browsing data support
	IsMutedChanged           = Capability("91.0.864.41") // Audio mute state changes
)

// WebView2 Runtime Version 90.0.818.41 (Released: May 2021)
const (
	WebResourceResponseReceived = Capability("90.0.818.41") // Web resource response handling
	DOMContentLoaded90          = Capability("90.0.818.41") // DOM content loaded events (v90)
	WebResourceRequested90      = Capability("90.0.818.41") // Web resource requests (v90)
	NewWindowWithOptions        = Capability("90.0.818.41") // New window with options
	CookieManagement            = Capability("90.0.818.41") // Cookie management support
)

// WebView2 Runtime Version 89.0.774.75 (Released: April 2021)
const (
	IsBuiltInErrorPageEnabled        = Capability("89.0.774.75") // Built-in error page support
	WebResourceResponse              = Capability("89.0.774.75") // Web resource response handling
	ScriptToExecuteOnDocumentCreated = Capability("89.0.774.75") // Document creation scripts
	EnvironmentOptions               = Capability("89.0.774.75") // Environment options support
	FrameNavigation89                = Capability("89.0.774.75") // Frame navigation (v89)
)

// WebView2 Runtime Version 88.0.705.74 (Released: March 2021)
const (
	WebResourceRequested88 = Capability("88.0.705.74") // Web resource requests (v88)
	PermissionRequested88  = Capability("88.0.705.74") // Permission requests (v88)
	ProcessFailed88        = Capability("88.0.705.74") // Process failure handling (v88)
	AddHostObjectToScript  = Capability("88.0.705.74") // Host object scripting
	IsMuted                = Capability("88.0.705.74") // Audio mute state
)

// WebView2 Runtime Version 87.0.664.75 (Released: February 2021)
const (
	WebMessageReceived87       = Capability("87.0.664.75") // Web message handling (v87)
	CallDevToolsProtocolMethod = Capability("87.0.664.75") // DevTools protocol method calls
	NewWindow87                = Capability("87.0.664.75") // New window creation (v87)
	DocumentTitleChanged87     = Capability("87.0.664.75") // Document title changes (v87)
	IsSuspended                = Capability("87.0.664.75") // Suspension state
)

// WebView2 Runtime Version 86.0.622.58 (Released: January 2021)
const (
	NavigationStarting86      = Capability("86.0.622.58") // Navigation start events (v86)
	NavigationCompleted86     = Capability("86.0.622.58") // Navigation completion (v86)
	FrameNavigationStarting86 = Capability("86.0.622.58") // Frame navigation start (v86)
	BasicWebView              = Capability("86.0.622.58") // Basic WebView2 functionality
	WindowBounds              = Capability("86.0.622.58") // Window bounds control
)

// WebView2 Runtime Version 85.0.564.70 (Released: December 2020)
const (
	WebView2Environment  = Capability("85.0.564.70") // WebView2 environment
	WebView2Controller   = Capability("85.0.564.70") // WebView2 controller
	BasicSettings        = Capability("85.0.564.70") // Basic settings support
	UserDataFolder       = Capability("85.0.564.70") // User data folder
	BrowserVersionString = Capability("85.0.564.70") // Browser version info
)

func HasCapability(webview2RuntimeVersion string, capability Capability) bool {
	result, err := webviewloader.CompareBrowserVersions(webview2RuntimeVersion, string(capability))
	if err != nil {
		return false
	}
	return result >= 0
}
