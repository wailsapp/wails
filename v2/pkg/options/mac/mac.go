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
	// EnableRetinaDevicePixelRatio enables the use of a private WebKit API to
	// correctly report window.devicePixelRatio in JavaScript on Retina displays
	// when using the wails:// custom URL scheme. Without this, canvas-based
	// rendering appears blurry on HiDPI screens.
	//
	// This is opt-in because it uses a private Apple API (_setOverrideDeviceScaleFactor:)
	// which may cause App Store rejection. Only enable if you are not distributing
	// through the Mac App Store.
	EnableRetinaDevicePixelRatio bool

	About      *AboutInfo
	OnFileOpen func(filePath string) `json:"-"`
	OnUrlOpen  func(filePath string) `json:"-"`
	// URLHandlers          map[string]func(string)
}
