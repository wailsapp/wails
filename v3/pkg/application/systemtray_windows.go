//go:build windows

package application

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/icons"

	"github.com/samber/lo"

	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

const (
	WM_USER_SYSTRAY = w32.WM_USER + 1
)

type windowsSystemTray struct {
	parent *SystemTray

	menu *Win32Menu

	cancelTheme func()
	uid         uint32
	hwnd        w32.HWND

	lightModeIcon      w32.HICON
	lightModeIconOwned bool
	darkModeIcon       w32.HICON
	darkModeIconOwned  bool
	currentIcon        w32.HICON
	currentIconOwned   bool
}

// releaseIcon destroys an icon handle only when we own it and no new handle reuses it.
// Shared handles (e.g. from LoadIcon/LoadIconWithResourceID) must not be passed to DestroyIcon per https://learn.microsoft.com/windows/win32/api/winuser/nf-winuser-destroyicon.
func (s *windowsSystemTray) releaseIcon(handle w32.HICON, owned bool, keep ...w32.HICON) {
	if !owned || handle == 0 {
		return
	}
	for _, k := range keep {
		if handle == k {
			return
		}
	}
	w32.DestroyIcon(handle)
}

func (s *windowsSystemTray) openMenu() {
	if s.menu == nil {
		return
	}
	// Get the system tray bounds
	trayBounds, err := s.bounds()
	if err != nil {
		return
	}
	if trayBounds == nil {
		return
	}

	// Show the menu at the tray bounds
	s.menu.ShowAt(trayBounds.X, trayBounds.Y)
}

func (s *windowsSystemTray) positionWindow(window Window, offset int) error {
	// Get the current screen trayBounds
	currentScreen, err := s.getScreen()
	if err != nil {
		return err
	}

	screenBounds := currentScreen.WorkArea
	windowBounds := window.Bounds()

	newX := screenBounds.Width - windowBounds.Width - offset
	newY := screenBounds.Height - windowBounds.Height - offset

	// systray icons in windows can either be in the taskbar
	// or in a flyout menu.
	var iconIsInTrayBounds bool
	iconIsInTrayBounds, err = s.iconIsInTrayBounds()
	if err != nil {
		return err
	}

	var trayBounds *Rect
	var centerAlignX, centerAlignY int

	// we only need the traybounds if the icon is in the tray
	if iconIsInTrayBounds {
		trayBounds, err = s.bounds()
		if err != nil {
			return err
		}
		if trayBounds == nil {
			return errors.New("failed to get system tray bounds")
		}
		*trayBounds = PhysicalToDipRect(*trayBounds)
		centerAlignX = trayBounds.X + (trayBounds.Width / 2) - (windowBounds.Width / 2)
		centerAlignY = trayBounds.Y + (trayBounds.Height / 2) - (windowBounds.Height / 2)
	}

	taskbarBounds := w32.GetTaskbarPosition()

	// Set the window position based on the icon location
	// if the icon is in the taskbar (traybounds) then we need
	// to adjust the position so the window is centered on the icon
	switch taskbarBounds.UEdge {
	case w32.ABE_LEFT:
		if iconIsInTrayBounds && centerAlignY <= newY {
			newY = centerAlignY
		}
		newX = screenBounds.X + offset
	case w32.ABE_TOP:
		if iconIsInTrayBounds && centerAlignX <= newX {
			newX = centerAlignX
		}
		newY = screenBounds.Y + offset
	case w32.ABE_RIGHT:
		if iconIsInTrayBounds && centerAlignY <= newY {
			newY = centerAlignY
		}
	case w32.ABE_BOTTOM:
		if iconIsInTrayBounds && centerAlignX <= newX {
			newX = centerAlignX
		}
	}
	newPos := currentScreen.relativeToAbsoluteDipPoint(Point{X: newX, Y: newY})
	windowBounds.X = newPos.X
	windowBounds.Y = newPos.Y
	window.SetBounds(windowBounds)
	return nil
}

func (s *windowsSystemTray) bounds() (*Rect, error) {
	if s.hwnd == 0 {
		return nil, errors.New("system tray window handle not initialized")
	}

	bounds, err := w32.GetSystrayBounds(s.hwnd, s.uid)
	if err != nil {
		return nil, err
	}
	if bounds == nil {
		return nil, errors.New("GetSystrayBounds returned nil")
	}

	monitor := w32.MonitorFromWindow(s.hwnd, w32.MONITOR_DEFAULTTONEAREST)
	if monitor == 0 {
		return nil, errors.New("failed to get monitor")
	}

	return &Rect{
		X:      int(bounds.Left),
		Y:      int(bounds.Top),
		Width:  int(bounds.Right - bounds.Left),
		Height: int(bounds.Bottom - bounds.Top),
	}, nil
}

func (s *windowsSystemTray) iconIsInTrayBounds() (bool, error) {
	if s.hwnd == 0 {
		return false, errors.New("system tray window handle not initialized")
	}

	bounds, err := w32.GetSystrayBounds(s.hwnd, s.uid)
	if err != nil {
		return false, err
	}
	if bounds == nil {
		return false, errors.New("GetSystrayBounds returned nil")
	}

	taskbarRect := w32.GetTaskbarPosition()
	if taskbarRect == nil {
		return false, errors.New("failed to get taskbar position")
	}

	inTasksBar := w32.RectInRect(bounds, &taskbarRect.Rc)
	if inTasksBar {
		return true, nil
	}

	return false, nil
}

func (s *windowsSystemTray) getScreen() (*Screen, error) {
	if s.hwnd == 0 {
		return nil, errors.New("system tray window handle not initialized")
	}
	// Get the screen for this systray
	return getScreenForWindowHwnd(s.hwnd)
}

func (s *windowsSystemTray) setMenu(menu *Menu) {
	s.updateMenu(menu)
}

func (s *windowsSystemTray) run() {
	s.hwnd = w32.CreateWindowEx(
		0,
		w32.MustStringToUTF16Ptr(globalApplication.options.Windows.WndClass),
		nil,
		0,
		0,
		0,
		0,
		0,
		w32.HWND_MESSAGE,
		0,
		0,
		nil)
	if s.hwnd == 0 {
		panic(syscall.GetLastError())
	}

	s.uid = uint32(s.parent.id)

	if _, err := s.show(); err != nil {
		// Initial systray add can fail when the shell is not available. This is handled downstream via TaskbarCreated message.
		globalApplication.warning("initial systray add failed: %v", err)
	}

	// Resolve the base icons once so we can reuse them for light/dark modes
	defaultIcon := getNativeApplication().windowClass.Icon

	if s.parent.icon != nil {
		s.lightModeIcon = lo.Must(w32.CreateSmallHIconFromImage(s.parent.icon))
		s.lightModeIconOwned = true
	} else if defaultIcon != 0 {
		s.lightModeIcon = defaultIcon
		s.lightModeIconOwned = false
	} else {
		s.lightModeIcon = lo.Must(w32.CreateSmallHIconFromImage(icons.SystrayLight))
		s.lightModeIconOwned = true
	}

	if s.parent.darkModeIcon != nil {
		s.darkModeIcon = lo.Must(w32.CreateSmallHIconFromImage(s.parent.darkModeIcon))
		s.darkModeIconOwned = true
	} else if s.parent.icon != nil {
		s.darkModeIcon = s.lightModeIcon
		s.darkModeIconOwned = false
	} else if defaultIcon != 0 {
		s.darkModeIcon = defaultIcon
		s.darkModeIconOwned = false
	} else {
		s.darkModeIcon = lo.Must(w32.CreateSmallHIconFromImage(icons.SystrayDark))
		s.darkModeIconOwned = true
	}

	if s.parent.menu != nil {
		s.updateMenu(s.parent.menu)
	}

	// Set Default Callbacks
	if s.parent.clickHandler == nil {
		s.parent.clickHandler = func() {
			globalApplication.debug("Left Button Clicked")
		}
	}

	if s.parent.rightClickHandler == nil {
		s.parent.rightClickHandler = func() {
			if s.menu != nil {
				s.openMenu()
			}
		}
	}

	// Update the icon
	s.updateIcon()

	// Listen for dark mode changes
	if s.cancelTheme != nil {
		s.cancelTheme()
	}
	s.cancelTheme = globalApplication.Event.OnApplicationEvent(events.Windows.SystemThemeChanged, func(event *ApplicationEvent) {
		s.updateIcon()
	})

	// Register the system tray
	getNativeApplication().registerSystemTray(s)
}

func (s *windowsSystemTray) updateIcon() {
	var newIcon w32.HICON
	if w32.IsCurrentlyDarkMode() {
		newIcon = s.darkModeIcon
	} else {
		newIcon = s.lightModeIcon
	}
	if s.currentIcon == newIcon {
		return
	}

	// Store the old icon to destroy it after updating
	oldIcon := s.currentIcon
	oldIconOwned := s.currentIconOwned

	s.currentIcon = newIcon
	nid := s.newNotifyIconData()
	nid.UFlags = w32.NIF_ICON
	if s.currentIcon != 0 {
		nid.HIcon = s.currentIcon
	}

	if !w32.ShellNotifyIcon(w32.NIM_MODIFY, &nid) {
		panic(syscall.GetLastError())
	}

	// Track ownership of the current icon so we know if we can destroy it later
	currentOwned := false
	if newIcon != 0 {
		if newIcon == s.lightModeIcon && s.lightModeIconOwned {
			currentOwned = true
		} else if newIcon == s.darkModeIcon && s.darkModeIconOwned {
			currentOwned = true
		}
	}
	s.currentIconOwned = currentOwned

	// Destroy the old icon handle if it exists, we owned it, and nothing else references it
	if oldIconOwned && oldIcon != 0 && oldIcon != s.lightModeIcon && oldIcon != s.darkModeIcon {
		w32.DestroyIcon(oldIcon)
	}
}

func (s *windowsSystemTray) newNotifyIconData() w32.NOTIFYICONDATA {
	nid := w32.NOTIFYICONDATA{
		UID:  s.uid,
		HWnd: s.hwnd,
	}
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	return nid
}

func (s *windowsSystemTray) setIcon(icon []byte) {
	newIcon, err := w32.CreateSmallHIconFromImage(icon)
	if err != nil {
		panic(err.Error())
	}

	oldLight := s.lightModeIcon
	oldLightOwned := s.lightModeIconOwned
	oldDark := s.darkModeIcon
	oldDarkOwned := s.darkModeIconOwned

	s.lightModeIcon = newIcon
	s.lightModeIconOwned = true

	// Keep dark mode in sync when both modes shared the same handle (or dark was unset).
	if s.darkModeIcon == 0 || s.darkModeIcon == oldLight {
		s.darkModeIcon = newIcon
		s.darkModeIconOwned = false
	}

	// Only free previous handles we own that are no longer referenced.
	s.releaseIcon(oldLight, oldLightOwned, s.lightModeIcon, s.darkModeIcon)
	if oldDark != s.darkModeIcon {
		s.releaseIcon(oldDark, oldDarkOwned, s.lightModeIcon, s.darkModeIcon)
	}

	s.updateIcon()
}

func (s *windowsSystemTray) setDarkModeIcon(icon []byte) {
	newIcon, err := w32.CreateSmallHIconFromImage(icon)
	if err != nil {
		panic(err.Error())
	}

	oldDark := s.darkModeIcon
	oldDarkOwned := s.darkModeIconOwned
	oldLight := s.lightModeIcon
	oldLightOwned := s.lightModeIconOwned

	s.darkModeIcon = newIcon
	s.darkModeIconOwned = true

	lightReplaced := false

	// Keep light mode in sync when both modes shared the same handle (or light was unset).
	if s.lightModeIcon == 0 || s.lightModeIcon == oldDark {
		s.lightModeIcon = newIcon
		s.lightModeIconOwned = false
		lightReplaced = true
	}

	// Only free the previous handle if nothing else keeps a reference to it.
	s.releaseIcon(oldDark, oldDarkOwned, s.lightModeIcon, s.darkModeIcon)
	if lightReplaced {
		s.releaseIcon(oldLight, oldLightOwned, s.lightModeIcon, s.darkModeIcon)
	}

	s.updateIcon()
}

func newSystemTrayImpl(parent *SystemTray) systemTrayImpl {
	return &windowsSystemTray{
		parent: parent,
	}
}

func (s *windowsSystemTray) wndProc(msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_USER_SYSTRAY:
		msg := lParam & 0xffff
		switch msg {
		case w32.WM_LBUTTONUP:
			if s.parent.clickHandler != nil {
				s.parent.clickHandler()
			}
		case w32.WM_RBUTTONUP:
			if s.parent.rightClickHandler != nil {
				s.parent.rightClickHandler()
			}
		case w32.WM_LBUTTONDBLCLK:
			if s.parent.doubleClickHandler != nil {
				s.parent.doubleClickHandler()
			}
		case w32.WM_RBUTTONDBLCLK:
			if s.parent.rightDoubleClickHandler != nil {
				s.parent.rightDoubleClickHandler()
			}
		case 0x0406:
			if s.parent.mouseEnterHandler != nil {
				s.parent.mouseEnterHandler()
			}
		case 0x0407:
			if s.parent.mouseLeaveHandler != nil {
				s.parent.mouseLeaveHandler()
			}
		}
		// println(w32.WMMessageToString(msg))

	// Menu processing
	case w32.WM_COMMAND:
		cmdMsgID := int(wParam & 0xffff)
		switch cmdMsgID {
		default:
			s.menu.ProcessCommand(cmdMsgID)
		}
	default:
		// msg := int(wParam & 0xffff)
		// println(w32.WMMessageToString(uintptr(msg)))
	}

	return w32.DefWindowProc(s.hwnd, msg, wParam, lParam)
}

func (s *windowsSystemTray) updateMenu(menu *Menu) {
	s.menu = NewPopupMenu(s.hwnd, menu)
	s.menu.onMenuOpen = s.parent.onMenuOpen
	s.menu.onMenuClose = s.parent.onMenuClose
	s.menu.Update()
}

func (s *windowsSystemTray) setTooltip(tooltip string) {
	// Create a new NOTIFYICONDATA structure
	nid := s.newNotifyIconData()
	nid.UFlags = w32.NIF_TIP | w32.NIF_SHOWTIP

	// Ensure the tooltip length is within the limit (128 characters including null terminate characters for szTip for Windows 2000 and later)
	// https://learn.microsoft.com/en-us/windows/win32/api/shellapi/ns-shellapi-notifyicondataw
	tooltipUTF16, err := w32.StringToUTF16(truncateUTF16(tooltip, 127))
	if err != nil {
		return
	}

	copy(nid.SzTip[:], tooltipUTF16)

	// Modify the tray icon with the new tooltip
	if !w32.ShellNotifyIcon(w32.NIM_MODIFY, &nid) {
		return
	}
}

// ---- Unsupported ----
func (s *windowsSystemTray) setLabel(label string) {}

func (s *windowsSystemTray) setTemplateIcon(_ []byte) {
	// Unsupported - do nothing
}

func (s *windowsSystemTray) setIconPosition(position IconPosition) {
	// Unsupported - do nothing
}

func (s *windowsSystemTray) destroy() {
	if s.cancelTheme != nil {
		s.cancelTheme()
		s.cancelTheme = nil
	}
	// Remove and delete the system tray
	getNativeApplication().unregisterSystemTray(s)
	if s.menu != nil {
		s.menu.Destroy()
	}

	// destroy the notification icon
	nid := s.newNotifyIconData()
	if !w32.ShellNotifyIcon(w32.NIM_DELETE, &nid) {
		globalApplication.debug(syscall.GetLastError().Error())
	}

	// Clean up icon handles
	lightIcon := s.lightModeIcon
	darkIcon := s.darkModeIcon
	currentIcon := s.currentIcon

	s.releaseIcon(lightIcon, s.lightModeIconOwned)
	s.releaseIcon(darkIcon, s.darkModeIconOwned, lightIcon)
	s.releaseIcon(currentIcon, s.currentIconOwned, lightIcon, darkIcon)

	s.lightModeIcon = 0
	s.lightModeIconOwned = false
	s.darkModeIcon = 0
	s.darkModeIconOwned = false
	s.currentIcon = 0
	s.currentIconOwned = false

	w32.DestroyWindow(s.hwnd)
	s.hwnd = 0
}

func (s *windowsSystemTray) Show() {
	if s.hwnd == 0 {
		return
	}

	nid := s.newNotifyIconData()
	nid.UFlags = w32.NIF_STATE
	nid.DwStateMask = w32.NIS_HIDDEN
	nid.DwState = 0
	if !w32.ShellNotifyIcon(w32.NIM_MODIFY, &nid) {
		globalApplication.debug("ShellNotifyIcon NIM_MODIFY show failed: %v", syscall.GetLastError())
	}
}

func (s *windowsSystemTray) Hide() {
	if s.hwnd == 0 {
		return
	}

	nid := s.newNotifyIconData()
	nid.UFlags = w32.NIF_STATE
	nid.DwStateMask = w32.NIS_HIDDEN
	nid.DwState = w32.NIS_HIDDEN
	if !w32.ShellNotifyIcon(w32.NIM_MODIFY, &nid) {
		globalApplication.debug("ShellNotifyIcon NIM_MODIFY hide failed: %v", syscall.GetLastError())
	}
}

func (s *windowsSystemTray) show() (w32.NOTIFYICONDATA, error) {
	nid := s.newNotifyIconData()
	nid.UFlags = w32.NIF_ICON | w32.NIF_MESSAGE
	nid.HIcon = s.currentIcon
	nid.UCallbackMessage = WM_USER_SYSTRAY

	if !w32.ShellNotifyIcon(w32.NIM_ADD, &nid) {
		err := syscall.GetLastError()
		return nid, fmt.Errorf("ShellNotifyIcon NIM_ADD failed: %w", err)
	}

	nid.UVersion = w32.NOTIFYICON_VERSION
	if !w32.ShellNotifyIcon(w32.NIM_SETVERSION, &nid) {
		err := syscall.GetLastError()
		return nid, fmt.Errorf("ShellNotifyIcon NIM_SETVERSION failed: %w", err)
	}

	if s.parent.tooltip != "" {
		s.setTooltip(s.parent.tooltip)
	}

	return nid, nil
}

func truncateUTF16(s string, maxUnits int) string {
	var units int
	for i, r := range s {
		var u int

		// check if rune will take 2 UTF-16 units
		if r > 0xFFFF {
			u = 2
		} else {
			u = 1
		}

		if units+u > maxUnits {
			return s[:i]
		}
		units += u
	}

	return s
}
