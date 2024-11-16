package application

import (
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/wailsapp/wails/v3/internal/assetserver"
)

// Service wraps a bound type instance.
// The zero value of Service is invalid.
// Valid values may only be obtained by calling [NewService].
type Service struct {
	instance any
	options  ServiceOptions
}

type ServiceOptions struct {
	// Name can be set to override the name of the service
	// This is useful for logging and debugging purposes
	Name string
	// Route is the path to the assets
	Route string
}

var DefaultServiceOptions = ServiceOptions{
	Route: "",
}

// NewService returns a Service value wrapping the given pointer.
// If T is not a named type, the returned value is invalid.
// The prefix is used if Service implements a http.Handler only one allowed
func NewService[T any](instance *T, options ...ServiceOptions) Service {
	if len(options) == 1 {
		return Service{instance, options[0]}
	}
	return Service{instance, DefaultServiceOptions}
}

func (s Service) Instance() any {
	return s.instance
}

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

	// Services allows you to bind Go methods to the frontend.
	Services []Service

	// BindAliases allows you to specify alias IDs for your bound methods.
	// Example: `BindAliases: map[uint32]uint32{1: 1411160069}` states that alias ID 1 maps to the Go method with ID 1411160069.
	BindAliases map[uint32]uint32

	// Logger i a slog.Logger instance used for logging Wails system messages (not application messages).
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
	PanicHandler func(any)

	// DisableDefaultSignalHandler disables the default signal handler
	DisableDefaultSignalHandler bool

	// KeyBindings is a map of key bindings to functions
	KeyBindings map[string]func(window *WebviewWindow)

	// OnShutdown is called when the application is about to terminate.
	// This is useful for cleanup tasks.
	// The shutdown process blocks until this function returns
	OnShutdown func()

	// ShouldQuit is a function that is called when the user tries to quit the application.
	// If the function returns true, the application will quit.
	// If the function returns false, the application will not quit.
	ShouldQuit func() bool

	// RawMessageHandler is called when the frontend sends a raw message.
	// This is useful for implementing custom frontend-to-backend communication.
	RawMessageHandler func(window Window, message string)

	// ErrorHandler is called when an error occurs
	ErrorHandler func(err error)

	// File extensions associated with the application
	// Example: [".txt", ".md"]
	// The '.' is required
	FileAssociations []string

	// This blank field ensures types from other packages
	// are never convertible to Options.
	// This property, in turn, improves the accuracy of the binding generator.
	_ struct{}
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
