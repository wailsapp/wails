package platform

import (
	"github.com/wailsapp/wails/v2/internal/platform/systray"
	"github.com/wailsapp/wails/v2/pkg/menu"
)
import "github.com/samber/lo"

type SysTray interface {
	// SetTitle sets the title of the tray menu
	SetTitle(title string)
	SetIcons(lightModeIcon []byte, darkModeIcon []byte) error
	SetTooltip(tooltip string) error
	Show() error
	Hide() error
	Run() error
	Close()
	AppendMenu(label string, callback menu.Callback)
	AppendMenuItem(item *menu.MenuItem)
	AppendSeparator()
}

func NewSysTray() SysTray {
	return lo.Must(systray.New())
}
