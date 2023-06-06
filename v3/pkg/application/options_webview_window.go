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
	Width                           int
	Height                          int
	AlwaysOnTop                     bool
	URL                             string
	DisableResize                   bool
	Frameless                       bool
	MinWidth                        int
	MinHeight                       int
	MaxWidth                        int
	MaxHeight                       int
	StartState                      WindowState
	Centered                        bool
	BackgroundType                  BackgroundType
	BackgroundColour                RGBA
	HTML                            string
	JS                              string
	CSS                             string
	X                               int
	Y                               int
	HideOnClose                     bool
	FullscreenButtonEnabled         bool
	Hidden                          bool
	EnableFraudulentWebsiteWarnings bool
	Zoom                            float64
	ZoomControlEnabled              bool
	EnableDragAndDrop               bool
	OpenInspectorOnStartup          bool
	Mac                             MacWindow
	Windows                         WindowsWindow
	Focused                         bool
	Menu                            *Menu
}

var WebviewWindowDefaults = &WebviewWindowOptions{
	Title:  "",
	Width:  800,
	Height: 600,
	URL:    "",
	BackgroundColour: RGBA{
		Red:   255,
		Green: 255,
		Blue:  255,
		Alpha: 255,
	},
}

type RGBA struct {
	Red, Green, Blue, Alpha uint8
}

type BackgroundType int

const (
	BackgroundTypeSolid BackgroundType = iota
	BackgroundTypeTransparent
	BackgroundTypeTranslucent
)
