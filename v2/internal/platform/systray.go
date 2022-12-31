package platform

import (
	"github.com/wailsapp/wails/v2/internal/platform/systray"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
)
import "github.com/samber/lo"

type SysTray interface {
	// SetTitle sets the title of the tray menu
	SetTitle(title string)
	SetTooltip(tooltip string) error
	Show() error
	Hide() error
	Run() error
	Close()
	SetMenu(menu *menu.Menu) error
	SetIcons(lightModeIcon, darkModeIcon *options.SystemTrayIcon) error
	Update() error
	OnLeftClick(func())
	OnRightClick(func())
	OnLeftDoubleClick(func())
	OnRightDoubleClick(func())
	OnMenuClose(func())
	OnMenuOpen(func())
}

func NewSysTray() SysTray {
	return lo.Must(systray.New())
}
