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
	TitleBar *TitleBar
}

// TitleBar contains options for the Mac titlebar
type TitleBar struct {
	AppearsTransparent   bool
	Hide                 bool
	HideTitle            bool
	FullSizeContent      bool
	UseToolbar           bool
	HideToolbarSeparator bool
}

// TitleBarDefault results in the default Mac Titlebar
func TitleBarDefault() *TitleBar {
	return &TitleBar{
		AppearsTransparent:   false,
		Hide:                 false,
		HideTitle:            false,
		FullSizeContent:      false,
		UseToolbar:           false,
		HideToolbarSeparator: false,
	}
}

// Credit: Comments from Electron site

// TitleBarHidden results in a hidden title bar and a full size content window,
// yet the title bar still has the standard window controls (“traffic lights”)
// in the top left.
func TitleBarHidden() *TitleBar {
	return &TitleBar{
		AppearsTransparent:   true,
		Hide:                 false,
		HideTitle:            true,
		FullSizeContent:      true,
		UseToolbar:           false,
		HideToolbarSeparator: false,
	}
}

// TitleBarHiddenInset results in a hidden title bar with an alternative look where
// the traffic light buttons are slightly more inset from the window edge.
func TitleBarHiddenInset() *TitleBar {

	return &TitleBar{
		AppearsTransparent:   true,
		Hide:                 false,
		HideTitle:            true,
		FullSizeContent:      true,
		UseToolbar:           true,
		HideToolbarSeparator: true,
	}

}
