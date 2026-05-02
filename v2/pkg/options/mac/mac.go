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
	// DisableEscapeExitsFullscreen, when true, prevents the Esc key from
	// exiting fullscreen mode. The keypress is swallowed by the window so
	// web content can handle it (e.g. closing modals). Default false
	// preserves standard macOS fullscreen behaviour.
	DisableEscapeExitsFullscreen bool
	// ActivationPolicy     ActivationPolicy
	About      *AboutInfo
	OnFileOpen func(filePath string) `json:"-"`
	OnUrlOpen  func(filePath string) `json:"-"`
	// URLHandlers          map[string]func(string)
}
