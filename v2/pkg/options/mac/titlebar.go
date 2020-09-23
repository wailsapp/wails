package mac

// TitleBar contains options for the Mac titlebar
type TitleBar struct {
	TitlebarAppearsTransparent bool
	HideTitle                  bool
	HideTitleBar               bool
	FullSizeContent            bool
	UseToolbar                 bool
	HideToolbarSeparator       bool
}

// Credit: Comments from Electron site

// TitleBarHidden results in a hidden title bar and a full size content window,
// yet the title bar still has the standard window controls (“traffic lights”)
// in the top left.
func TitleBarHidden() *TitleBar {
	return &TitleBar{
		TitlebarAppearsTransparent: true,
		HideTitle:                  true,
		HideTitleBar:               false,
		FullSizeContent:            true,
		UseToolbar:                 false,
		HideToolbarSeparator:       false,
	}
}
