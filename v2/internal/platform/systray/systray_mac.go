//go:build darwin

package systray

import (
	"errors"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
)

var NotImplementedSysTray = errors.New("not implemented")

type Systray struct {
}

func (p *Systray) Close() {
	err := p.Stop()
	if err != nil {
		println(err.Error())
	}
}

func (p *Systray) Update() error {
	return NotImplementedSysTray
}

func (p *Systray) SetTitle(_ string) {}

func New() (*Systray, error) {
	return nil, NotImplementedSysTray
}

func (p *Systray) SetMenu(popupMenu *menu.Menu) (err error) {
	return NotImplementedSysTray
}

func (p *Systray) Stop() error {
	return NotImplementedSysTray
}

func (p *Systray) OnLeftClick(fn func()) {

}

func (p *Systray) OnRightClick(fn func()) {

}

func (p *Systray) OnLeftDoubleClick(fn func()) {

}

func (p *Systray) OnRightDoubleClick(fn func()) {

}

func (p *Systray) OnMenuClose(fn func()) {

}

func (p *Systray) OnMenuOpen(fn func()) {

}

func (p *Systray) SetTooltip(tooltip string) error {
	return NotImplementedSysTray
}

func (p *Systray) ShowMessage(title, msg string, bigIcon bool) error {
	return NotImplementedSysTray
}

func (p *Systray) Show() error {
	return p.setVisible(true)
}

func (p *Systray) Hide() error {
	return p.setVisible(false)
}

func (p *Systray) setVisible(visible bool) error {
	return NotImplementedSysTray
}

func (p *Systray) SetIcons(lightModeIcon, darkModeIcon *options.SystemTrayIcon) error {
	return NotImplementedSysTray
}

func (p *Systray) Run() error {
	return NotImplementedSysTray
}
