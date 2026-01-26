//go:build ios

package application

// iOS doesn't have system tray support
// These are placeholder implementations

func (t *SystemTray) update() {}

func (t *SystemTray) setMenu(menu *Menu) {
	// iOS doesn't have system tray
}

func (t *SystemTray) close() {
	// iOS doesn't have system tray
}


func (t *SystemTray) attachWindow(window *WebviewWindow) {
	// iOS doesn't have system tray
}

func (t *SystemTray) detachWindow(windowID uint) {
	// iOS doesn't have system tray
}

type iosSystemTray struct {
	parent *SystemTray
}

func newSystemTrayImpl(s *SystemTray) systemTrayImpl {
	return &iosSystemTray{
		parent: s,
	}
}

func (s *iosSystemTray) run() {
	// iOS doesn't have system tray
}

func (s *iosSystemTray) setLabel(_ string) {
	// iOS doesn't have system tray
}

func (s *iosSystemTray) setMenu(_ *Menu) {
	// iOS doesn't have system tray
}

func (s *iosSystemTray) setIcon(_ []byte) {
	// iOS doesn't have system tray
}

func (s *iosSystemTray) setDarkModeIcon(_ []byte) {
	// iOS doesn't have system tray
}

func (s *iosSystemTray) destroy() {
	// iOS doesn't have system tray
}

func (s *iosSystemTray) setIconPosition(_ IconPosition) {
	// iOS doesn't have system tray
}

func (s *iosSystemTray) positionWindow(_ Window, _ int) error {
	return nil
}

func (s *iosSystemTray) detachWindowPositioning(_ uint) {
	// iOS doesn't have system tray
}

func (s *iosSystemTray) setTemplateIcon(_ []byte) {
	// iOS doesn't have system tray
}

func (s *iosSystemTray) openMenu() {
	// iOS doesn't have system tray
}

func (s *iosSystemTray) setTooltip(_ string) {
	// iOS doesn't have system tray
}

func (s *iosSystemTray) bounds() (*Rect, error) {
	return nil, nil
}

func (s *iosSystemTray) getScreen() (*Screen, error) {
	screens, err := getScreens()
	if err != nil || len(screens) == 0 {
		return nil, err
	}
	return screens[0], nil
}

func (s *iosSystemTray) Show() {
	// iOS doesn't have system tray
}

func (s *iosSystemTray) Hide() {
	// iOS doesn't have system tray
}

func (s *iosSystemTray) onAttachedWindowHidden() {
	// iOS doesn't have system tray
}

func (s *iosSystemTray) onAttachedWindowShown() {
	// iOS doesn't have system tray
}