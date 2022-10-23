package application

import (
	"github.com/wailsapp/wails/v2/internal/platform"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// SystemTray defines a system tray!
type SystemTray struct {
	title              string
	hidden             bool
	lightModeIcon      *options.SystemTrayIcon
	darkModeIcon       *options.SystemTrayIcon
	tooltip            string
	startHidden        bool
	menu               *menu.Menu
	onLeftClick        func()
	onRightClick       func()
	onLeftDoubleClick  func()
	onRightDoubleClick func()
	onMenuClose        func()
	onMenuOpen         func()

	// The platform specific implementation
	impl platform.SysTray
}

func newSystemTray(options *options.SystemTray) *SystemTray {
	return &SystemTray{
		title:              options.Title,
		lightModeIcon:      options.LightModeIcon,
		darkModeIcon:       options.DarkModeIcon,
		tooltip:            options.Tooltip,
		startHidden:        options.StartHidden,
		menu:               options.Menu,
		onLeftClick:        options.OnLeftClick,
		onRightClick:       options.OnRightClick,
		onLeftDoubleClick:  options.OnLeftDoubleClick,
		onRightDoubleClick: options.OnRightDoubleClick,
		onMenuOpen:         options.OnMenuOpen,
		onMenuClose:        options.OnMenuClose,
	}
}

func (t *SystemTray) run() {
	t.impl = platform.NewSysTray()
	t.impl.SetTitle(t.title)
	t.impl.SetIcons(t.lightModeIcon, t.darkModeIcon)
	t.impl.SetTooltip(t.tooltip)
	t.impl.OnLeftClick(t.onLeftClick)
	t.impl.OnRightClick(t.onRightClick)
	t.impl.OnLeftDoubleClick(t.onLeftDoubleClick)
	t.impl.OnRightDoubleClick(t.onRightDoubleClick)
	t.impl.OnMenuOpen(t.onMenuOpen)
	t.impl.OnMenuClose(t.onMenuClose)
	if !t.startHidden {
		t.impl.Show()
	}
	t.impl.SetMenu(t.menu)
	t.impl.Run()
}

func (t *SystemTray) SetTitle(title string) {
	if t.impl != nil {
		t.impl.SetTitle(title)
	} else {
		t.title = title
	}
}

func (t *SystemTray) Run() error {
	t.run()
	return nil
}

func (t *SystemTray) Close() {
	if t.impl != nil {
		t.impl.Close()
		t.impl = nil
	}
}

func (t *SystemTray) SetMenu(items *menu.Menu) {
	if t.impl != nil {
		t.impl.SetMenu(t.menu)
	} else {
		t.menu = items
	}
}

func (t *SystemTray) Update() error {
	if t.impl != nil {
		return t.impl.Update()
	}
	return nil
}

func (t *SystemTray) SetTooltip(s string) {
	if t.impl != nil {
		t.impl.SetTooltip(s)
	} else {
		t.tooltip = s
	}
}

func (t *SystemTray) SetIcons(lightModeIcon *options.SystemTrayIcon, darkModeIcon *options.SystemTrayIcon) {
	if t.impl != nil {
		t.impl.SetIcons(lightModeIcon, darkModeIcon)
	} else {
		t.lightModeIcon = lightModeIcon
		t.darkModeIcon = darkModeIcon
	}

}

func (t *SystemTray) OnLeftClick(fn func()) {
	if t.impl != nil {
		t.impl.OnLeftClick(fn)
	}
}

func (t *SystemTray) OnRightClick(fn func()) {
	if t.impl != nil {
		t.impl.OnRightClick(fn)
	}
}

func (t *SystemTray) OnLeftDoubleClick(fn func()) {
	if t.impl != nil {
		t.impl.OnLeftDoubleClick(fn)
	}
}

func (t *SystemTray) OnRightDoubleClick(fn func()) {
	if t.impl != nil {
		t.impl.OnRightDoubleClick(fn)
	}
}

func (t *SystemTray) OnMenuOpen(fn func()) {
	if t.impl != nil {
		t.impl.OnMenuOpen(fn)
	}
}

func (t *SystemTray) OnMenuClose(fn func()) {
	if t.impl != nil {
		t.impl.OnMenuClose(fn)
	}
}
