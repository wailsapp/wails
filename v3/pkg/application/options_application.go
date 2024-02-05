package application

import (
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/wailsapp/wails/v3/internal/assetserver"
)

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

	// Bind allows you to bind Go methods to the frontend.
	Bind []any

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

	// Plugins is a map of plugins used by the application
	Plugins map[string]Plugin

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
