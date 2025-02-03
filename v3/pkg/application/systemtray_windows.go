//go:build windows

package application

import (
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/icons"
	"syscall"
	"time"
	"unsafe"

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

	// Platform specific implementation
	uid           uint32
	hwnd          w32.HWND
	lightModeIcon w32.HICON
	darkModeIcon  w32.HICON
	currentIcon   w32.HICON
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

	// Show the menu at the tray bounds
	s.menu.ShowAt(trayBounds.X, trayBounds.Y)
}

func (s *windowsSystemTray) positionWindow(window *WebviewWindow, offset int) error {
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
	iconIsInTrayBounds, err := s.iconIsInTrayBounds()
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
	bounds, err := w32.GetSystrayBounds(s.hwnd, s.uid)
	if err != nil {
		return nil, err
	}

	monitor := w32.MonitorFromWindow(s.hwnd, w32.MONITOR_DEFAULTTONEAREST)
	if monitor == 0 {
		return nil, fmt.Errorf("failed to get monitor")
	}

	return &Rect{
		X:      int(bounds.Left),
		Y:      int(bounds.Top),
		Width:  int(bounds.Right - bounds.Left),
		Height: int(bounds.Bottom - bounds.Top),
	}, nil
}

func (s *windowsSystemTray) iconIsInTrayBounds() (bool, error) {
	bounds, err := w32.GetSystrayBounds(s.hwnd, s.uid)
	if err != nil {
		return false, err
	}

	taskbarRect := w32.GetTaskbarPosition()

	inTasksBar := w32.RectInRect(bounds, &taskbarRect.Rc)
	if inTasksBar {
		return true, nil
	}

	return false, nil
}

func (s *windowsSystemTray) getScreen() (*Screen, error) {
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

	nid := w32.NOTIFYICONDATA{
		HWnd:             s.hwnd,
		UID:              uint32(s.parent.id),
		UFlags:           w32.NIF_ICON | w32.NIF_MESSAGE,
		HIcon:            s.currentIcon,
		UCallbackMessage: WM_USER_SYSTRAY,
	}
	nid.CbSize = uint32(unsafe.Sizeof(nid))

	for retries := range 6 {
		if !w32.ShellNotifyIcon(w32.NIM_ADD, &nid) {
			if retries == 5 {
				globalApplication.fatal("Failed to register system tray icon: %v", syscall.GetLastError())
			}

			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}

	nid.UVersion = w32.NOTIFYICON_VERSION

	if !w32.ShellNotifyIcon(w32.NIM_SETVERSION, &nid) {
		panic(syscall.GetLastError())
	}

	// Get the application icon if available
	defaultIcon := w32.LoadIconWithResourceID(w32.GetModuleHandle(""), w32.RT_ICON)
	if defaultIcon != 0 {
		s.lightModeIcon = defaultIcon
		s.darkModeIcon = defaultIcon
	} else {
		s.lightModeIcon = lo.Must(w32.CreateSmallHIconFromImage(icons.SystrayLight))
		s.darkModeIcon = lo.Must(w32.CreateSmallHIconFromImage(icons.SystrayDark))
	}

	if s.parent.icon != nil {
		s.lightModeIcon = lo.Must(w32.CreateSmallHIconFromImage(s.parent.icon))
	}
	if s.parent.darkModeIcon != nil {
		s.darkModeIcon = lo.Must(w32.CreateSmallHIconFromImage(s.parent.darkModeIcon))
	}
	s.uid = nid.UID

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
	globalApplication.OnApplicationEvent(events.Windows.SystemThemeChanged, func(event *ApplicationEvent) {
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

	s.currentIcon = newIcon
	nid := s.newNotifyIconData()
	nid.UFlags = w32.NIF_ICON
	if s.currentIcon != 0 {
		nid.HIcon = s.currentIcon
	}

	if !w32.ShellNotifyIcon(w32.NIM_MODIFY, &nid) {
		panic(syscall.GetLastError())
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
	var err error
	s.lightModeIcon, err = w32.CreateSmallHIconFromImage(icon)
	if err != nil {
		panic(syscall.GetLastError())
	}
	if s.darkModeIcon == 0 {
		s.darkModeIcon = s.lightModeIcon
	}
	// Update the icon
	s.updateIcon()
}

func (s *windowsSystemTray) setDarkModeIcon(icon []byte) {
	var err error
	s.darkModeIcon, err = w32.CreateSmallHIconFromImage(icon)
	if err != nil {
		panic(syscall.GetLastError())
	}
	if s.lightModeIcon == 0 {
		s.lightModeIcon = s.darkModeIcon
	}
	// Update the icon
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

// ---- Unsupported ----

func (s *windowsSystemTray) setLabel(_ string) {
	// Unsupported - do nothing
}

func (s *windowsSystemTray) setTemplateIcon(_ []byte) {
	// Unsupported - do nothing
}

func (s *windowsSystemTray) setIconPosition(position IconPosition) {
	// Unsupported - do nothing
}

func (s *windowsSystemTray) destroy() {
	// Remove and delete the system tray
	getNativeApplication().unregisterSystemTray(s)
	if s.menu != nil {
		s.menu.Destroy()
	}
	w32.DestroyWindow(s.hwnd)
	// destroy the notification icon
	nid := s.newNotifyIconData()
	if !w32.ShellNotifyIcon(w32.NIM_DELETE, &nid) {
		globalApplication.debug(syscall.GetLastError().Error())
	}
}

func (s *windowsSystemTray) Show() {
	// No-op
}

func (s *windowsSystemTray) Hide() {
	// No-op
}
