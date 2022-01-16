package windows

// Options are options specific to Windows
type Options struct {
	WebviewIsTransparent bool
	WindowIsTranslucent  bool
	DisableWindowIcon    bool

	// Draw a border around the window, even if the window is frameless
	EnableFramelessBorder             bool
	NotifyParentWindowPositionChanged func()
}
