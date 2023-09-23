package application

import (
	"io/fs"
	"log/slog"
	"net/http"
)

type Options struct {
	// Name is the name of the application
	Name string

	// Description is the description of the application (used in the default about box)
	Description string

	// Icon is the icon of the application (used in the default about box)
	Icon []byte

	// Mac is the Mac specific configuration for Mac builds
	Mac MacOptions

	// Windows is the Windows specific configuration for Windows builds
	Windows WindowsApplicationOptions

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

	// PanicHandler is a way to register a custom panic handler
	PanicHandler func(any)

	// KeyBindings is a map of key bindings to functions
	KeyBindings map[string]func(window *WebviewWindow)
}

// AssetOptions defines the configuration of the AssetServer.
type AssetOptions struct {
	// FS defines the static assets to be used. A GET request is first tried to be served from this FS. If the FS returns
	// `os.ErrNotExist` for that file, the request handling will fallback to the Handler and tries to serve the GET
	// request from it.
	//
	// If set to nil, all GET requests will be forwarded to Handler.
	FS fs.FS

	// Handler will be called for every GET request that can't be served from FS, due to `os.ErrNotExist`. Furthermore all
	// non GET requests will always be served from this Handler.
	//
	// If not defined, the result is the following in cases where the Handler would have been called:
	//   GET request:   `http.StatusNotFound`
	//   Other request: `http.StatusMethodNotAllowed`
	Handler http.Handler

	// Middleware is HTTP Middleware which allows to hook into the AssetServer request chain. It allows to skip the default
	// request handler dynamically, e.g. implement specialized Routing etc.
	// The Middleware is called to build a new `http.Handler` used by the AssetSever and it also receives the default
	// handler used by the AssetServer as an argument.
	//
	// If not defined, the default AssetServer request chain is executed.
	//
	// Multiple Middlewares can be chained together with:
	//   ChainMiddleware(middleware ...Middleware) Middleware
	Middleware Middleware

	// External URL can be set to a development server URL so that all requests are forwarded to it. This is useful
	// when using a development server like `vite` or `snowpack` which serves the assets on a different port.
	ExternalURL string
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
