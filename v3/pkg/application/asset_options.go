//go:build !wails_native

package application

import (
	"io/fs"
	"net/http"

	"github.com/wailsapp/wails/v3/internal/assetserver"
)

// AssetOptions defines the configuration of the AssetServer.
type AssetOptions struct {
	Handler        http.Handler
	Middleware     Middleware
	DisableLogging bool
}

// Middleware defines HTTP middleware that can be applied to the AssetServer.
type Middleware func(next http.Handler) http.Handler

// ChainMiddleware allows chaining multiple middlewares into one middleware.
func ChainMiddleware(middleware ...Middleware) Middleware {
	return func(handler http.Handler) http.Handler {
		for i := len(middleware) - 1; i >= 0; i-- {
			handler = middleware[i](handler)
		}
		return handler
	}
}

func AssetFileServerFS(assets fs.FS) http.Handler {
	return assetserver.NewAssetFileServer(assets)
}

func BundledAssetFileServer(assets fs.FS) http.Handler {
	return assetserver.NewBundledAssetFileServer(assets)
}
