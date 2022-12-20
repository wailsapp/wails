//go:build linux

/*
 * Based on code originally from https://github.com/tadvi/systray. Copyright (C) 2019 The Systray Authors. All Rights Reserved.
 */

package systray

import (
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
)

type Systray struct {
}

func (p *Systray) Close() {
	err := p.Stop()
	if err != nil {
		println(err.Error())
	}
}

func (p *Systray) Update() error {
	return nil
}

// SetTitle is unused on Windows
func (p *Systray) SetTitle(_ string) {}

func New() (*Systray, error) {
	return nil, nil
}

func (p *Systray) SetMenu(popupMenu *menu.Menu) (err error) {
	return
}

func (p *Systray) Stop() error {
	return nil
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
	return nil
}

func (p *Systray) Show() error {
	return p.setVisible(true)
}

func (p *Systray) Hide() error {
	return p.setVisible(false)
}

func (p *Systray) setVisible(visible bool) error {
	return nil
}

func (p *Systray) SetIcons(lightModeIcon, darkModeIcon *options.SystemTrayIcon) error {
	return nil
}

func (p *Systray) Run() error {
	return nil
}
