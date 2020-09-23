package mac

// Options ae options speific to Mac
type Options struct {
	TitleBar *TitleBar
}

// TitleBar contains options for the Mac titlebar
type TitleBar struct {
	TitlebarAppearsTransparent bool
	HideTitle                  bool
	HideTitleBar               bool
	FullSizeContent            bool
	UseToolbar                 bool
	HideToolbarSeparator       bool
}
