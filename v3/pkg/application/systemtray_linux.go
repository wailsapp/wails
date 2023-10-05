//go:build linux

package application

import "log"

type linuxSystemTray struct {
	id    uint
	label string
	icon  []byte
	menu  *Menu

	iconPosition   int
	isTemplateIcon bool
}

func (s *linuxSystemTray) setIconPosition(position int) {
	s.iconPosition = position
	log.Println("linuxSystemTray.setIconPosition() not implemented")
}

func (s *linuxSystemTray) setMenu(menu *Menu) {
	s.menu = menu
	log.Println("linuxSystemTray.setMenu() not implemented")
}

func (s *linuxSystemTray) positionWindow(window *WebviewWindow, offset int) error {
	log.Println("linuxSystemTray.positionWindow() not implemented")
}

func (s *linuxSystemTray) getScreen() (*Screen, error) {
	log.Println("linuxSystemTray.getScreen() not implemented")
}

func (s *linuxSystemTray) bounds() (*Rect, error) {
	log.Println("linuxSystemTray.bounds() not implemented")
}

func (s *linuxSystemTray) run() {
	log.Println("linuxSystemTray.run() - not implemented")
}

func (s *linuxSystemTray) setIcon(icon []byte) {
	s.icon = icon
	log.Println("linuxSystemTray.setIcon() not implemented")
}

func (s *linuxSystemTray) setDarkModeIcon(icon []byte) {
	s.icon = icon
	log.Println("linuxSystemTray.setDarkModeIcon() not implemented")
}

func (s *linuxSystemTray) setTemplateIcon(icon []byte) {
	s.icon = icon
	s.isTemplateIcon = true
	log.Println("linuxSystemTray.setTemplateIcon() not implemented")
}

func newSystemTrayImpl(s *SystemTray) systemTrayImpl {
	return &linuxSystemTray{
		id:             s.id,
		label:          s.label,
		icon:           s.icon,
		menu:           s.menu,
		iconPosition:   s.iconPosition,
		isTemplateIcon: s.isTemplateIcon,
	}
}

func (s *linuxSystemTray) openMenu() {
	log.Println("linuxSystemTray.openMenu() not implemented")
}

func (s *linuxSystemTray) setLabel(label string) {
	s.label = label
	log.Println("linuxSystemTray.setLabel() not implemented")
}

func (s *linuxSystemTray) destroy() {
	// Remove the status item from the status bar and its associated menu
	log.Println("linuxSystemTray.destroy() not implemented")
}
