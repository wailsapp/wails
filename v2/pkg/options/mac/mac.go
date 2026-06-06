package mac

//type ActivationPolicy int
//
//const (
//	NSApplicationActivationPolicyRegular    ActivationPolicy = 0
//	NSApplicationActivationPolicyAccessory  ActivationPolicy = 1
//	NSApplicationActivationPolicyProhibited ActivationPolicy = 2
//)

type AboutInfo struct {
	Title   string
	Message string
	Icon    []byte
}

// Options are options specific to Mac
type Options struct {
	TitleBar             *TitleBar
	Appearance           AppearanceType
	ContentProtection    bool
	WebviewIsTransparent bool
	WindowIsTranslucent  bool
	Preferences          *Preferences
	DisableZoom          bool
	// ActivationPolicy     ActivationPolicy
	About      *AboutInfo
	OnFileOpen func(filePath string) `json:"-"`
	OnUrlOpen  func(filePath string) `json:"-"`
	// URLHandlers          map[string]func(string)

	// DisableEscapeExitsFullscreen prevents the Escape key from exiting fullscreen mode.
	// When true, web content can handle the Escape key (e.g., to close modals) without
	// triggering the macOS system behaviour that exits the fullscreen window.
	DisableEscapeExitsFullscreen bool
}
