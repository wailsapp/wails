package application

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


type RGBA struct {
	Red, Green, Blue, Alpha uint8
}
