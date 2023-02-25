package application

import (
	"io/fs"
	"net/http"
)

type WindowState int

const (
	WindowStateNormal WindowState = iota
	WindowStateMinimised
	WindowStateMaximised
	WindowStateFullscreen
)

type WebviewWindowOptions struct {
	Name                            string
	Title                           string
	Width, Height                   int
	AlwaysOnTop                     bool
	URL                             string
	DisableResize                   bool
	Frameless                       bool
	MinWidth                        int
	MinHeight                       int
	MaxWidth                        int
	MaxHeight                       int
	StartState                      WindowState
	Mac                             MacWindow
	BackgroundColour                *RGBA
	Assets                          AssetOptions
	HTML                            string
	JS                              string
	CSS                             string
	X                               int
	Y                               int
	FullscreenButtonEnabled         bool
	Hidden                          bool
	EnableFraudulentWebsiteWarnings bool
	Zoom                            float64
	EnableDragAndDrop               bool
}

var WebviewWindowDefaults = &WebviewWindowOptions{
	Title:  "",
	Width:  800,
	Height: 600,
	URL:    "",
}

type AssetOptions struct {
	// FS to use for loading assets from
	FS fs.FS
	// Handler is a custom handler to use for serving assets. If this is set, the `URL` and `FS` fields are ignored.
	Handler http.Handler
	// Middleware is a custom middleware to use for serving assets. If this is set, the `URL` and `FS` fields are ignored.
	Middleware func(http.Handler) http.Handler
}

type RGBA struct {
	Red, Green, Blue, Alpha uint8
}
