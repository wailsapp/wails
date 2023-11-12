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
	WebviewIsTransparent bool
	WindowIsTranslucent  bool
	Preferences          *Preferences
	// ActivationPolicy     ActivationPolicy
	About      *AboutInfo
	OnFileOpen func(filePath string) `json:"-"`
	OnUrlOpen  func(filePath string) `json:"-"`
	// URLHandlers          map[string]func(string)
}
