package application

import (
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/wailsapp/wails/v3/internal/assetserver"
)

// Options contains the options for the application
type Options struct {
	// Name is the name of the application (used in the default about box)
	Name string

	// Description is the description of the application (used in the default about box)
	Description string

	// Icon is the icon of the application (used in the default about box)
	Icon []byte

	// Mac is the Mac specific configuration for Mac builds
	Mac MacOptions

	// Windows is the Windows specific configuration for Windows builds
	Windows WindowsOptions

	// Linux is the Linux specific configuration for Linux builds
	Linux LinuxOptions

	// IOS is the iOS specific configuration for iOS builds
	IOS IOSOptions

	// Android is the Android specific configuration for Android builds
	Android AndroidOptions

	// Services allows you to bind Go methods to the frontend.
	Services []Service

	// MarshalError will be called if non-nil
	// to marshal to JSON the error values returned by service methods.
	//
	// MarshalError is not allowed to fail,
	// but it may return a nil slice to fall back
	// to the default error handling mechanism.
	//
	// If the returned slice is not nil, it must contain valid JSON.
	MarshalError func(error) []byte

	// BindAliases allows you to specify alias IDs for your bound methods.
	// Example: `BindAliases: map[uint32]uint32{1: 1411160069}` states that alias ID 1 maps to the Go method with ID 1411160069.
	BindAliases map[uint32]uint32

	// Logger is a slog.Logger instance used for logging Wails system messages (not application messages).
	// If not defined, a default logger is used.
	Logger *slog.Logger

	// LogLevel defines the log level of the Wails system logger.
	LogLevel slog.Level

	// Assets are the application assets to be used.
	Assets AssetOptions

	// Flags are key value pairs that are available to the frontend.
	// This is also used by Wails to provide information to the frontend.
	Flags map[string]any

	// PanicHandler is called when a panic occurs
	PanicHandler func(*PanicDetails)

	// DisableDefaultSignalHandler disables the default signal handler
	DisableDefaultSignalHandler bool

	// KeyBindings is a map of key bindings to functions
	KeyBindings map[string]func(window Window)

	// OnShutdown is called when the application is about to terminate.
	// This is useful for cleanup tasks.
	// The shutdown process blocks until this function returns.
	OnShutdown func()

	// PostShutdown is called after the application
	// has finished shutting down, just before process termination.
	// This is useful for testing and logging purposes
	// on platforms where the Run() method does not return.
	// When PostShutdown is called, the application instance is not usable anymore.
	// The shutdown process blocks until this function returns.
	PostShutdown func()

	// ShouldQuit is a function that is called when the user tries to quit the application.
	// If the function returns true, the application will quit.
	// If the function returns false, the application will not quit.
	ShouldQuit func() bool

	// RawMessageHandler is called when the frontend sends a raw message.
	// This is useful for implementing custom frontend-to-backend communication.
	RawMessageHandler func(window Window, message string, originInfo *OriginInfo)

	// WarningHandler is called when a warning occurs
	WarningHandler func(string)

	// ErrorHandler is called when an error occurs
	ErrorHandler func(err error)

	// File extensions associated with the application
	// Example: [".txt", ".md"]
	// The '.' is required
	FileAssociations []string

	// SingleInstance options for single instance functionality
	SingleInstance *SingleInstanceOptions

	// Transport allows you to provide a custom IPC transport layer.
	// When set, Wails will use your transport instead of the default HTTP fetch-based transport.
	// This allows you to use WebSockets, custom protocols, or any other transport mechanism
	// while retaining all Wails generated bindings and event communication.
	//
	// The default transport uses HTTP fetch requests to /wails/runtime + events via js.Exec in webview.
	// If not specified, the default transport is used.
	//
	// Example use case: Implementing WebSocket-based or PostMessage IPC.
	Transport Transport
}

// AssetOptions defines the configuration of the AssetServer.
type AssetOptions struct {
	// Handler which serves all the content to the WebView.
	Handler http.Handler

	// Middleware is a HTTP Middleware which allows to hook into the AssetServer request chain. It allows to skip the default
	// request handler dynamically, e.g. implement specialized Routing etc.
	// The Middleware is called to build a new `http.Handler` used by the AssetSever and it also receives the default
	// handler used by the AssetServer as an argument.
	//
	// This middleware injects itself before any of Wails internal middlewares.
	//
	// If not defined, the default AssetServer request chain is executed.
	//
	// Multiple Middlewares can be chained together with:
	//   ChainMiddleware(middleware ...Middleware) Middleware
	Middleware Middleware

	// DisableLogging disables logging of the AssetServer. By default, the AssetServer logs every request.
	DisableLogging bool
}

// Middleware defines HTTP middleware that can be applied to the AssetServer.
// The handler passed as next is the next handler in the chain. One can decide to call the next handler
// or implement a specialized handling.
type Middleware func(next http.Handler) http.Handler

// ChainMiddleware allows chaining multiple middlewares to one middleware.
func ChainMiddleware(middleware ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		for i := len(middleware) - 1; i >= 0; i-- {
			h = middleware[i](h)
		}
		return h
	}
}

// AssetFileServerFS returns a http handler which serves the assets from the fs.FS.
// If an external devserver has been provided 'FRONTEND_DEVSERVER_URL' the files are being served
// from the external server, ignoring the `assets`.
func AssetFileServerFS(assets fs.FS) http.Handler {
	return assetserver.NewAssetFileServer(assets)
}

// BundledAssetFileServer returns a http handler which serves the assets from the fs.FS.
// If an external devserver has been provided 'FRONTEND_DEVSERVER_URL' the files are being served
// from the external server, ignoring the `assets`.
// It also serves the compiled runtime.js file at `/wails/runtime.js`.
// It will provide the production runtime.js file from the embedded assets if the `production` tag is used.
func BundledAssetFileServer(assets fs.FS) http.Handler {
	return assetserver.NewBundledAssetFileServer(assets)
}

/******** Mac Options ********/

// ActivationPolicy is the activation policy for the application.
type ActivationPolicy int

const (
	// ActivationPolicyRegular is used for applications that have a user interface,
	ActivationPolicyRegular ActivationPolicy = iota
	// ActivationPolicyAccessory is used for applications that do not have a main window,
	// such as system tray applications or background applications.
	ActivationPolicyAccessory
	ActivationPolicyProhibited
)

// MacOptions contains options for macOS applications.
type MacOptions struct {
	// ActivationPolicy is the activation policy for the application. Defaults to
	// applicationActivationPolicyRegular.
	ActivationPolicy ActivationPolicy
	// If set to true, the application will terminate when the last window is closed.
	ApplicationShouldTerminateAfterLastWindowClosed bool
}

/****** Windows Options *******/

// WindowsOptions contains options for Windows applications.
type WindowsOptions struct {

	// Window class name
	// Default: WailsWebviewWindow
	WndClass string

	// WndProcInterceptor is a function that will be called for every message sent in the application.
	// Use this to hook into the main message loop. This is useful for handling custom window messages.
	// If `shouldReturn` is `true` then `returnCode` will be returned by the main message loop.
	// If `shouldReturn` is `false` then returnCode will be ignored and the message will be processed by the main message loop.
	WndProcInterceptor func(hwnd uintptr, msg uint32, wParam, lParam uintptr) (returnCode uintptr, shouldReturn bool)

	// DisableQuitOnLastWindowClosed disables the auto quit of the application if the last window has been closed.
	DisableQuitOnLastWindowClosed bool

	// Path where the WebView2 stores the user data. If empty %APPDATA%\[BinaryName.exe] will be used.
	// If the path is not valid, a messagebox will be displayed with the error and the app will exit with error code.
	WebviewUserDataPath string

	// Path to the directory with WebView2 executables. If empty WebView2 installed in the system will be used.
	WebviewBrowserPath string
}

/********* Linux Options *********/

// LinuxOptions contains options for Linux applications.
type LinuxOptions struct {
	// DisableQuitOnLastWindowClosed disables the auto quit of the application if the last window has been closed.
	DisableQuitOnLastWindowClosed bool

	// ProgramName is used to set the program's name for the window manager via GTK's g_set_prgname().
	//This name should not be localized. [see the docs]
	//
	//When a .desktop file is created this value helps with window grouping and desktop icons when the .desktop file's Name
	//property differs form the executable's filename.
	//
	//[see the docs]: https://docs.gtk.org/glib/func.set_prgname.html
	ProgramName string
}

/********* iOS Options *********/

// IOSOptions contains options for iOS applications.
type IOSOptions struct {
    // DisableInputAccessoryView controls whether the iOS WKWebView shows the
    // input accessory toolbar (the bar with Next/Previous/Done) above the keyboard.
    // Default: false (accessory bar is shown).
    // true  => accessory view is disabled/hidden
    // false => accessory view is enabled/shown
    DisableInputAccessoryView bool

    // Scrolling & Bounce (defaults: scroll/bounce/indicators are enabled on iOS)
    // Use Disable* to keep default true behavior without surprising zero-values.
    DisableScroll           bool
    DisableBounce           bool
    DisableScrollIndicators bool

    // Navigation gestures (default false)
    EnableBackForwardNavigationGestures bool

    // Link previews (default true on iOS)
    // Use Disable* so default (false) means previews are enabled.
    DisableLinkPreview bool

    // Media playback
    // Inline playback (default false) -> Enable*
    EnableInlineMediaPlayback        bool
    // Autoplay without user action (default false) -> Enable*
    EnableAutoplayWithoutUserAction  bool

    // Inspector / Debug (default true in dev)
    // Use Disable* so default (false) keeps inspector enabled.
    DisableInspectable bool

    // User agent customization
    // If empty, defaults apply. ApplicationNameForUserAgent defaults to "wails.io".
    UserAgent                   string
    ApplicationNameForUserAgent string

    // App-wide background colour for the main iOS window prior to any WebView creation.
    // If AppBackgroundColourSet is true, the delegate will apply this colour to the app window
    // during didFinishLaunching. Otherwise, it defaults to white.
    AppBackgroundColourSet bool
    BackgroundColour       RGBA

    // EnableNativeTabs enables a native iOS UITabBar at the bottom of the screen.
    // When enabled, the native tab bar will dispatch a 'nativeTabSelected' CustomEvent
    // to the window with detail: { index: number }.
    // NOTE: If NativeTabsItems has one or more entries, native tabs are auto-enabled
    // regardless of this flag, and the provided items will be used.
    EnableNativeTabs bool

    // NativeTabsItems configures the labels and optional SF Symbol icons for the
    // native UITabBar. If one or more items are provided, native tabs are automatically
    // enabled. If empty and EnableNativeTabs is true, default items are used.
    NativeTabsItems []NativeTabItem
}

// NativeTabItem describes a single item in the iOS native UITabBar.
// SystemImage is the SF Symbols name to use for the icon (iOS 13+). If empty or
// unavailable on the current OS, no icon is shown.
type NativeTabItem struct {
    Title       string        `json:"Title"`
    SystemImage NativeTabIcon `json:"SystemImage"`
}

// NativeTabIcon is a string-based enum for SF Symbols.
// It allows using predefined constants for common symbols while still accepting
// any valid SF Symbols name as a plain string.
//
// Example:
//  NativeTabsItems: []NativeTabItem{
//    { Title: "Home", SystemImage: NativeTabIconHouse },
//    { Title: "Settings", SystemImage: "gearshape" }, // arbitrary string still allowed
//  }
type NativeTabIcon string

const (
    // Common icons
    NativeTabIconNone    NativeTabIcon = ""
    NativeTabIconHouse   NativeTabIcon = "house"
    NativeTabIconGear    NativeTabIcon = "gear"
    NativeTabIconStar    NativeTabIcon = "star"
    NativeTabIconPerson  NativeTabIcon = "person"
    NativeTabIconBell    NativeTabIcon = "bell"
    NativeTabIconMagnify NativeTabIcon = "magnifyingglass"
    NativeTabIconList    NativeTabIcon = "list.bullet"
    NativeTabIconFolder  NativeTabIcon = "folder"
)
<<<<<<< HEAD
=======

/********* Android Options *********/

// AndroidOptions contains options for Android applications.
type AndroidOptions struct {
	// DisableScroll disables scrolling in the WebView
	DisableScroll bool

	// DisableBounce disables the overscroll bounce effect
	DisableOverscroll bool

	// EnableZoom allows pinch-to-zoom in the WebView (default: false)
	EnableZoom bool

	// UserAgent sets a custom user agent string
	UserAgent string

	// BackgroundColour sets the background colour of the WebView
	BackgroundColour RGBA

	// DisableHardwareAcceleration disables hardware acceleration for the WebView
	DisableHardwareAcceleration bool
}
>>>>>>> origin/v3-alpha-feature/android-support
