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
	Assets AssetOptions
}

type AssetOptions struct {
	// FS to use for loading assets from
	FS fs.FS
	// Handler is a custom handler to use for serving assets. If this is set, the `URL` and `FS` fields are ignored.
	Handler http.Handler
	// Middleware is a custom middleware to use for serving assets. If this is set, the `URL` and `FS` fields are ignored.
	Middleware func(http.Handler) http.Handler
}
