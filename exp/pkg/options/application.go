package options

type Application struct {
	Mac *Mac
}

// Mac contains macOS specific options

type ActivationPolicy int

const (
	ActivationPolicyRegular ActivationPolicy = iota
	ActivationPolicyAccessory
	ActivationPolicyProhibited
)

type Mac struct {
	// ActivationPolicy is the activation policy for the application. Defaults to
	// applicationActivationPolicyRegular.
	ActivationPolicy ActivationPolicy
}

type Window struct {
	Title            string
	Width, Height    int
	AlwaysOnTop      bool
	URL              string
	DisableResize    bool
	Resizable        bool
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

type MacBackdrop int

const (
	MacBackdropNormal MacBackdrop = iota
	MacBackdropTransparent
	MacBackdropTranslucent
)

type WindowState int

const (
	WindowStateNormal WindowState = iota
	WindowStateMinimised
	WindowStateMaximised
	WindowStateFullscreen
)

// MacWindow contains macOS specific options
type MacWindow struct {
	Backdrop MacBackdrop
}
