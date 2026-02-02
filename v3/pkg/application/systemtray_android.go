//go:build android

package application

// Android doesn't have system tray support
// These are placeholder implementations

func (t *SystemTray) update() {}

func (t *SystemTray) setMenu(menu *Menu) {
	// Android doesn't have system tray
}

func (t *SystemTray) close() {
	// Android doesn't have system tray
}

func (t *SystemTray) attachWindow(window *WebviewWindow) {
	// Android doesn't have system tray
}

func (t *SystemTray) detachWindow(windowID uint) {
	// Android doesn't have system tray
}

type androidSystemTray struct {
	parent *SystemTray
}

func newSystemTrayImpl(s *SystemTray) systemTrayImpl {
	return &androidSystemTray{
		parent: s,
	}
}

func (s *androidSystemTray) run() {
	// Android doesn't have system tray
}

func (s *androidSystemTray) setLabel(_ string) {
	// Android doesn't have system tray
}

func (s *androidSystemTray) setMenu(_ *Menu) {
	// Android doesn't have system tray
}

func (s *androidSystemTray) setIcon(_ []byte) {
	// Android doesn't have system tray
}

func (s *androidSystemTray) setDarkModeIcon(_ []byte) {
	// Android doesn't have system tray
}

func (s *androidSystemTray) destroy() {
	// Android doesn't have system tray
}

func (s *androidSystemTray) setIconPosition(_ IconPosition) {
	// Android doesn't have system tray
}

func (s *androidSystemTray) positionWindow(_ Window, _ int) error {
	return nil
}

func (s *androidSystemTray) detachWindowPositioning(_ uint) {
	// Android doesn't have system tray
}

func (s *androidSystemTray) setTemplateIcon(_ []byte) {
	// Android doesn't have system tray
}

func (s *androidSystemTray) openMenu() {
	// Android doesn't have system tray
}

func (s *androidSystemTray) setTooltip(_ string) {
	// Android doesn't have system tray
}

func (s *androidSystemTray) bounds() (*Rect, error) {
	return nil, nil
}

func (s *androidSystemTray) getScreen() (*Screen, error) {
	screens, err := getScreens()
	if err != nil || len(screens) == 0 {
		return nil, err
	}
	return screens[0], nil
}

func (s *androidSystemTray) Show() {
	// Android doesn't have system tray
}

func (s *androidSystemTray) Hide() {
	// Android doesn't have system tray
}
