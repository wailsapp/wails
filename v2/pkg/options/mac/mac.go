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

	// DisableFullscreenEscapeIntercept controls the Escape key behavior in fullscreen mode.
	// When false (default), pressing Escape in fullscreen will dispatch the event to the WebView
	// instead of exiting fullscreen, allowing JavaScript to handle it (e.g., to close modals).
	// When true, Escape will exit fullscreen using native macOS behavior.
	DisableFullscreenEscapeIntercept bool
}
