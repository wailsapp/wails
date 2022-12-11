package options

type WindowState int

const (
	WindowStateNormal WindowState = iota
	WindowStateMinimised
	WindowStateMaximised
	WindowStateFullscreen
)

type Window struct {
	// Alias is a human-readable name for the window. This can be used to reference the window in the frontend.
	Alias            string
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
