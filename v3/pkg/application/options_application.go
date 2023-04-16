package application

import (
	"io/fs"
	"net/http"

	"github.com/wailsapp/wails/v3/pkg/logger"
)

type Options struct {
	Name        string
	Description string
	Icon        []byte
	Mac         MacOptions
	Bind        []any
	Logger      struct {
		Silent        bool
		CustomLoggers []logger.Output
	}
	Assets  AssetOptions
	Plugins map[string]Plugin
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

	// Middleware is a HTTP Middleware which allows to hook into the AssetServer request chain. It allows to skip the default
	// request handler dynamically, e.g. implement specialized Routing etc.
	// The Middleware is called to build a new `http.Handler` used by the AssetSever and it also receives the default
	// handler used by the AssetServer as an argument.
	//
	// If not defined, the default AssetServer request chain is executed.
	//
	// Multiple Middlewares can be chained together with:
	//   ChainMiddleware(middleware ...Middleware) Middleware
	Middleware Middleware
}

// Middleware defines a HTTP middleware that can be applied to the AssetServer.
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
