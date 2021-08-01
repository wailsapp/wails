package windows

import (
	"embed"
	"github.com/wailsapp/wails/v2/pkg/menu"
)

// Options are options specific to Windows
type Options struct {
	WebviewIsTransparent          bool
	WindowBackgroundIsTranslucent bool
	DisableWindowIcon             bool
	Menu                          *menu.Menu
	Assets                        *embed.FS
}
