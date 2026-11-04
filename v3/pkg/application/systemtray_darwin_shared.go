//go:build darwin && !ios && !server

package application

import (
	"errors"
	"strings"
	"unsafe"
)

type macosSystemTray struct {
	id    uint
	label string
	icon  []byte
	menu  *Menu

	nsStatusItem      unsafe.Pointer
	nsImage           unsafe.Pointer
	nsMenu            unsafe.Pointer
	iconPosition      IconPosition
	isTemplateIcon    bool
	parent            *SystemTray
	lastClickedScreen unsafe.Pointer

	// macosSystemTrayNative is the backend's own state: empty for cgo.
	macosSystemTrayNative

	// displayedMenu is the menu on screen, if any. SetMenu can replace s.menu
	// while it is open, so the close must be reported to this one.
	displayedMenu *Menu
}

func (s *macosSystemTray) setDarkModeIcon(icon []byte) {
	s.setIcon(icon)
}

func (s *macosSystemTray) setIconPosition(position IconPosition) {
	s.iconPosition = position
}

func (s *macosSystemTray) setTooltip(tooltip string) {
	// Tooltips not supported on macOS
}

func hasANSICodes(s string) bool {
	return strings.Contains(s, "\033[")
}

func newSystemTrayImpl(s *SystemTray) systemTrayImpl {
	result := &macosSystemTray{
		parent:         s,
		id:             s.id,
		label:          s.label,
		icon:           s.icon,
		menu:           s.menu,
		iconPosition:   s.iconPosition,
		isTemplateIcon: s.isTemplateIcon,
	}
	systemTrayMap[s.id] = result
	return result
}

func (s *macosSystemTray) processClick(b button) {
	switch b {
	case leftButtonDown:
		// Check if we have a callback
		if s.parent.clickHandler != nil {
			s.parent.clickHandler()
			return
		}
		if s.parent.attachedWindow.Window != nil {
			s.parent.defaultClickHandler()
			return
		}
		if s.menu != nil {
			s.showMenu()
		}
	case rightButtonDown:
		// Check if we have a callback
		if s.parent.rightClickHandler != nil {
			s.parent.rightClickHandler()
			return
		}
		if s.menu != nil {
			if s.parent.attachedWindow.Window != nil {
				s.parent.attachedWindow.Window.Hide()
			}
			s.showMenu()
			return
		}
		if s.parent.attachedWindow.Window != nil {
			s.parent.defaultClickHandler()
		}
	}
}

func (s *macosSystemTray) openMenu() {
	if s.nsMenu == nil {
		return
	}
	s.showMenu()
}

func (s *macosSystemTray) setMenu(menu *Menu) {
	s.menu = menu
	if s.nsStatusItem != nil && menu != nil {
		menu.Update()
		s.nsMenu = (menu.impl).(*macosMenu).nsMenu
		s.cacheMenu()
	}
}

func (s *macosSystemTray) run() {
	globalApplication.dispatchOnMainThread(func() {
		if s.nsStatusItem != nil {
			Fatal("System tray '%d' already running", s.id)
		}
		s.nsStatusItem = s.newStatusItem()

		if s.label != "" {
			s.setLabel(s.label)
		}
		if s.icon != nil {
			s.applyIcon(s.icon)
		}
		if s.menu != nil {
			s.menu.Update()
			// Convert impl to macosMenu object
			s.nsMenu = (s.menu.impl).(*macosMenu).nsMenu
			s.cacheMenu()
		}
	})
}

func (s *macosSystemTray) setIcon(icon []byte) {
	s.icon = icon
	globalApplication.dispatchOnMainThread(func() {
		s.applyIcon(icon)
	})
}

func (s *macosSystemTray) setTemplateIcon(icon []byte) {
	s.icon = icon
	s.isTemplateIcon = true
	globalApplication.dispatchOnMainThread(func() {
		s.applyIcon(icon)
	})
}

func (s *macosSystemTray) getScreen() (*Screen, error) {
	if s.lastClickedScreen != nil {
		result := &Screen{
			Bounds: s.lastClickedScreenBounds(),
		}
		return result, nil
	}
	return nil, errors.New("no screen available")
}

// systrayPreClick is called from the NSEvent local monitor BEFORE the
// button processes the mouse-down.  It returns 1 when the framework should
// show the menu via native tracking (proper highlight, no app activation),
// or 0 to let the action handler fire for custom click/window behaviour.
func systrayPreClick(trayID uint, buttonID int) int {
	systemTray := systemTrayMap[trayID]
	if systemTray == nil || systemTray.nsMenu == nil {
		return 0
	}
	b := button(buttonID)
	switch b {
	case leftButtonDown:
		if systemTray.parent.clickHandler == nil &&
			systemTray.parent.attachedWindow.Window == nil {
			return 1
		}
	case rightButtonDown:
		if systemTray.parent.rightClickHandler == nil {
			// Hide the attached window before the menu appears.
			if systemTray.parent.attachedWindow.Window != nil &&
				systemTray.parent.attachedWindow.Window.IsVisible() {
				systemTray.parent.attachedWindow.Window.Hide()
			}
			return 1
		}
	}
	return 0
}

// systrayMenuDisplayed is called on the main thread when the tray
// menu opens (displayed=1) or closes (displayed=0), so Menu.Clear keeps the IDs
// of the rows on screen resolvable however often the menu is updated (#4719).
func systrayMenuDisplayed(trayID uint, displayed int) {
	systemTray := systemTrayMap[trayID]
	if systemTray == nil {
		return
	}
	if systemTray.displayedMenu != nil {
		systemTray.displayedMenu.displayed.Store(false)
		systemTray.displayedMenu = nil
	}
	if displayed != 0 && systemTray.menu != nil {
		systemTray.displayedMenu = systemTray.menu
		systemTray.displayedMenu.displayed.Store(true)
	}
}

func systrayClick(trayID uint, buttonID int) {
	// Get the system tray
	systemTray := systemTrayMap[trayID]
	if systemTray == nil {
		globalApplication.error("system tray not found: %v", trayID)
		return
	}
	systemTray.processClick(button(buttonID))
}
