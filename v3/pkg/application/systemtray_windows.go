//go:build windows

package application

import (
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/w32"
	"syscall"
	"unsafe"
)

type windowsSystemTray struct {
	parent *SystemTray

	// Platform specific implementation
	uid           uint32
	hwnd          w32.HWND
	appIcon       w32.HICON
	lightModeIcon w32.HICON
	darkModeIcon  w32.HICON
	currentIcon   w32.HICON
	//menu          *w32.PopupMenu
}

func (s *windowsSystemTray) setIconPosition(position int) {
	// Unsupported - do nothing
}

func (s *windowsSystemTray) setMenu(menu *Menu) {
	//s.menu = menu
	panic("implement me")
}

func (s *windowsSystemTray) run() {

	NotifyIconClassName := "WailsSystray"
	_, err := w32.RegisterWindow(NotifyIconClassName, getNativeApplication().wndProc)
	if err != nil {
		panic(err)
	}

	s.hwnd = w32.CreateWindowEx(
		0,
		w32.MustStringToUTF16Ptr(NotifyIconClassName),
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
		UCallbackMessage: w32.WM_USER + uint32(s.parent.id),
	}
	nid.CbSize = uint32(unsafe.Sizeof(nid))

	if !w32.ShellNotifyIcon(w32.NIM_ADD, &nid) {
		panic(syscall.GetLastError())
	}

	nid.UVersion = w32.NOTIFYICON_VERSION

	if !w32.ShellNotifyIcon(w32.NIM_SETVERSION, &nid) {
		panic(syscall.GetLastError())
	}

	defaultIcon, err := w32.CreateHIconFromPNG(s.parent.icon)
	if err != nil {
		panic(err)
	}
	s.lightModeIcon = defaultIcon
	s.darkModeIcon = defaultIcon
	s.uid = nid.UID

	// TODO: Set Menu

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
	nid.UFlags = w32.NIF_ICON | w32.NIF_MESSAGE
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
	// TODO:
	var err error
	if w32.IsCurrentlyDarkMode() {
		s.darkModeIcon, err = w32.CreateHIconFromPNG(icon)
		if err != nil {
			panic(syscall.GetLastError())
		}
	} else {
		s.lightModeIcon, err = w32.CreateHIconFromPNG(icon)
		if err != nil {
			panic(syscall.GetLastError())
		}
	}
	// Update the icon
	s.updateIcon()
}

func newSystemTrayImpl(parent *SystemTray) systemTrayImpl {
	return &windowsSystemTray{
		parent: parent,
	}
}

func (s *windowsSystemTray) destroy() {
	panic("implement me")
}

func (s *windowsSystemTray) wndProc(msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case w32.WM_USER + uint32(s.parent.id):
		switch lParam {
		case w32.WM_LBUTTONUP:
			if s.parent.leftButtonClickHandler != nil {
				s.parent.leftButtonClickHandler()
			}
		case w32.WM_RBUTTONUP:
			if s.parent.rightButtonClickHandler != nil {
				s.parent.rightButtonClickHandler()
			}
		case w32.WM_LBUTTONDBLCLK:
			if s.parent.leftButtonDoubleClickHandler != nil {
				s.parent.leftButtonDoubleClickHandler()
			}
		case w32.WM_RBUTTONDBLCLK:
			if s.parent.rightButtonDoubleClickHandler != nil {
				s.parent.rightButtonDoubleClickHandler()
			}
		default:
			println(w32.WMMessageToString(lParam))
		}
		// TODO: Menu processing
	//case w32.WM_COMMAND:
	//	cmdMsgID := int(wparam & 0xffff)
	//	switch cmdMsgID {
	//	default:
	//		p.menu.ProcessCommand(cmdMsgID)
	//	}
	default:
		msg := int(wParam & 0xffff)
		println(w32.WMMessageToString(uintptr(msg)))
	}

	return w32.DefWindowProc(s.hwnd, msg, wParam, lParam)
}

// ---- Unsupported ----

func (s *windowsSystemTray) setLabel(_ string) {
	// Unsupported - do nothing
}

func (s *windowsSystemTray) setTemplateIcon(_ []byte) {
	// Unsupported - do nothing
}
