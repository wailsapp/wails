//go:build windows

package application

import (
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/icons"
	"syscall"
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

func (s *windowsSystemTray) positionWindow(window *WebviewWindow, offset int) error {

	// Get the trayBounds of this system tray
	trayBounds, err := s.bounds()
	if err != nil {
		return err
	}

	// Get the current screen trayBounds
	currentScreen, err := s.getScreen()
	if err != nil {
		return err
	}

	screenBounds := currentScreen.WorkArea

	taskbarBounds := w32.GetTaskbarPosition()
	switch taskbarBounds.UEdge {
	case w32.ABE_LEFT:
		if trayBounds == nil {
			// Move it to the bottom left corner of the screen
			window.SetRelativePosition(offset, screenBounds.Height-window.Height())
			return nil
		}
		newHeight := trayBounds.Y - (window.Height() / 2)
		if newHeight < 0 {
			newHeight = 0
		}
		// Move it to the top left corner of the screen
		window.SetRelativePosition(offset, newHeight)
	case w32.ABE_TOP:
		if trayBounds == nil {
			// Move it to the top right corner of the screen
			window.SetRelativePosition(screenBounds.Width-window.Width(), offset)
			return nil
		}
		newWidth := trayBounds.X - (window.Width() / 2)
		if newWidth > screenBounds.Width-window.Width() {
			newWidth = screenBounds.Width - window.Width()
		}
		// Move it to the top left corner of the screen
		window.SetRelativePosition(newWidth, offset)
	case w32.ABE_RIGHT:
		if trayBounds == nil {
			// Move it to the bottom right corner of the screen
			window.SetRelativePosition(screenBounds.Width-window.Width()-offset, screenBounds.Height-window.Height())
			return nil
		}
		newHeight := trayBounds.Y - (window.Height() / 2)
		if newHeight > screenBounds.Height-window.Height() {
			newHeight = screenBounds.Height - window.Height()
		}
		window.SetRelativePosition(screenBounds.Width-window.Width()-offset, newHeight)
	case w32.ABE_BOTTOM:
		if trayBounds == nil {
			// Move it to the bottom right corner of the screen
			window.SetRelativePosition(screenBounds.Width-window.Width(), screenBounds.Height-window.Height()-offset)
			return nil
		}
		newWidth := trayBounds.X - (window.Width() / 2)
		if newWidth > screenBounds.Width-window.Width() {
			newWidth = screenBounds.Width - window.Width()
		}
		window.SetRelativePosition(newWidth, screenBounds.Height-window.Height()-offset)
	}
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

	// Get the taskbar rect
	taskbarRect := w32.GetTaskbarPosition()

	flyoutOpen := !w32.RectInRect(bounds, &taskbarRect.Rc)
	if flyoutOpen {
		return nil, nil
	}

	return &Rect{
		X:      int(bounds.Left),
		Y:      int(bounds.Top),
		Width:  int(bounds.Right - bounds.Left),
		Height: int(bounds.Bottom - bounds.Top),
	}, nil
}

func (s *windowsSystemTray) getScreen() (*Screen, error) {
	// Get the screen for this systray
	return getScreen(s.hwnd)
}

func (s *windowsSystemTray) setMenu(menu *Menu) {
	s.updateMenu(menu)
}

func (s *windowsSystemTray) run() {
	s.hwnd = w32.CreateWindowEx(
		0,
		windowClassName,
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

	if !w32.ShellNotifyIcon(w32.NIM_ADD, &nid) {
		panic(syscall.GetLastError())
	}

	nid.UVersion = w32.NOTIFYICON_VERSION

	if !w32.ShellNotifyIcon(w32.NIM_SETVERSION, &nid) {
		panic(syscall.GetLastError())
	}

	if s.parent.icon != nil {
		s.lightModeIcon = lo.Must(w32.CreateSmallHIconFromImage(s.parent.icon))
	} else {
		s.lightModeIcon = lo.Must(w32.CreateSmallHIconFromImage(icons.SystrayLight))
	}
	if s.parent.darkModeIcon != nil {
		s.darkModeIcon = lo.Must(w32.CreateSmallHIconFromImage(s.parent.darkModeIcon))
	} else {
		s.darkModeIcon = lo.Must(w32.CreateSmallHIconFromImage(icons.SystrayDark))
	}
	s.uid = nid.UID

	if s.parent.menu != nil {
		s.updateMenu(s.parent.menu)
	}

	// Set Default Callbacks
	if s.parent.clickHandler == nil {
		s.parent.clickHandler = func() {
			println("Left Button Clicked")
		}
	}
	if s.parent.rightClickHandler == nil {
		s.parent.rightClickHandler = func() {
			if s.menu != nil {
				s.menu.ShowAtCursor()
			}
		}
	}

	// Update the icon
	s.updateIcon()

	// Listen for dark mode changes
	globalApplication.On(events.Windows.SystemThemeChanged, func() {
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
	return
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
		//println(w32.WMMessageToString(msg))

	// Menu processing
	case w32.WM_COMMAND:
		cmdMsgID := int(wParam & 0xffff)
		switch cmdMsgID {
		default:
			s.menu.ProcessCommand(cmdMsgID)
		}
	default:
		//msg := int(wParam & 0xffff)
		//println(w32.WMMessageToString(uintptr(msg)))
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

func (s *windowsSystemTray) setIconPosition(position int) {
	// Unsupported - do nothing
}

func (s *windowsSystemTray) destroy() {
	// Remove and delete the system tray
	getNativeApplication().unregisterSystemTray(s)
	s.menu.Destroy()
	w32.DestroyWindow(s.hwnd)
	// Destroy the notification icon
	nid := s.newNotifyIconData()
	if !w32.ShellNotifyIcon(w32.NIM_DELETE, &nid) {
		globalApplication.info(syscall.GetLastError().Error())
	}
}
