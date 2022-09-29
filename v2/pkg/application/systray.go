package application

import (
	"github.com/wailsapp/wails/v2/internal/platform"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// SystemTray defines a system tray!
type SystemTray struct {
	title         string
	hidden        bool
	lightModeIcon []byte
	darkModeIcon  []byte
	tooltip       string
	startHidden   bool

	// The platform specific implementation
	impl platform.SysTray
}

func newSystemTray(options *options.SystemTray) *SystemTray {
	return &SystemTray{
		impl:          platform.NewSysTray(),
		title:         options.Title,
		lightModeIcon: options.LightModeIcon,
		darkModeIcon:  options.DarkModeIcon,
		tooltip:       options.Tooltip,
		startHidden:   options.StartHidden,
	}
}

func (t *SystemTray) run() {
	t.impl.SetTitle(t.title)
	t.impl.SetIcons(t.lightModeIcon, t.darkModeIcon)
	t.impl.SetTooltip(t.tooltip)
	if !t.startHidden {
		t.impl.Show()
	}
	t.impl.Run()
}

func (t *SystemTray) SetTitle(title string) {
	t.title = title
	t.impl.SetTitle(title)
}
func (t *SystemTray) AppendMenu(label string, callback menu.Callback) {
	t.impl.AppendMenu(label, callback)
}

func (t *SystemTray) Run() error {
	t.run()
	return nil
}

func (t *SystemTray) Close() {
	t.impl.Close()
}

func (t *SystemTray) AppendSeperator() {
	t.impl.AppendSeparator()
}

func (t *SystemTray) AppendMenuItem(item *menu.MenuItem) {
	t.impl.AppendMenuItem(item)
}
