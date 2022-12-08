package options

type WindowState int

const (
	WindowStateNormal WindowState = iota
	WindowStateMinimised
	WindowStateMaximised
	WindowStateFullscreen
)

type Window struct {
	Title            string
	Width, Height    int
	AlwaysOnTop      bool
	URL              string
	DisableResize    bool
	MinWidth         int
	MinHeight        int
	MaxWidth         int
	MaxHeight        int
	EnableDevTools   bool
	StartState       WindowState
	Mac              *MacWindow
	BackgroundColour *RGBA
}

type RGBA struct {
	Red, Green, Blue, Alpha uint8
}
