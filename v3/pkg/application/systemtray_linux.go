//go:build linux

package application

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
}

func (s *linuxSystemTray) setMenu(menu *Menu) {
	s.menu = menu
}

func (s *linuxSystemTray) positionWindow(window *WebviewWindow, offset int) error {
	return nil
}

func (s *linuxSystemTray) getScreen() (*Screen, error) {
	return &Screen{}, nil
}

func (s *linuxSystemTray) bounds() (*Rect, error) {
	return &Rect{}, nil
}

func (s *linuxSystemTray) run() {
	globalApplication.error("linuxSystemTray.run() - not implemented")
}

func (s *linuxSystemTray) setIcon(icon []byte) {
	s.icon = icon
}

func (s *linuxSystemTray) setDarkModeIcon(icon []byte) {
	s.icon = icon
}

func (s *linuxSystemTray) setTemplateIcon(icon []byte) {
	s.icon = icon
	s.isTemplateIcon = true
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
}

func (s *linuxSystemTray) setLabel(label string) {
	s.label = label
}

func (s *linuxSystemTray) destroy() {
	// Remove the status item from the status bar and its associated menu
}
