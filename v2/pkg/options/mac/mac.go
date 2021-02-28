package mac

import "github.com/wailsapp/wails/v2/pkg/menu"

type ActivationPolicy int

const (
	NSApplicationActivationPolicyRegular    ActivationPolicy = 0
	NSApplicationActivationPolicyAccessory  ActivationPolicy = 1
	NSApplicationActivationPolicyProhibited ActivationPolicy = 2
)

// Options are options specific to Mac
type Options struct {
	TitleBar                      *TitleBar
	Appearance                    AppearanceType
	WebviewIsTransparent          bool
	WindowBackgroundIsTranslucent bool
	Menu                          *menu.Menu
	TrayMenus                     []*menu.TrayMenu
	ContextMenus                  []*menu.ContextMenu
	ActivationPolicy              ActivationPolicy
	URLHandlers                   map[string]func(string)
}
