/*
 * Based on code originally from https://github.com/tadvi/systray. Copyright (C) 2019 The Systray Authors. All Rights Reserved.
 */

package systray

import (
	"errors"
	"github.com/samber/lo"
	"github.com/wailsapp/wails/v2/internal/platform/win32"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"syscall"
	"unsafe"
)

var (
	user32 = syscall.MustLoadDLL("user32.dll")

	DefWindowProc   = user32.MustFindProc("DefWindowProcW")
	RegisterClassEx = user32.MustFindProc("RegisterClassExW")
	CreateWindowEx  = user32.MustFindProc("CreateWindowExW")

	windowClasses = map[string]win32.HINSTANCE{}
)

type Systray struct {
	id     uint32
	hwnd   win32.HWND
	hinst  win32.HINSTANCE
	lclick func()
	rclick func()

	appIcon       win32.HICON
	lightModeIcon win32.HICON
	darkModeIcon  win32.HICON
	currentIcon   win32.HICON

	Menu []*menu.MenuItem

	quit chan struct{}
	icon *options.SystemTrayIcon
}

func (p *Systray) Close() {
	err := p.Stop()
	if err != nil {
		println(err.Error())
	}
}

// SetTitle is unused on Windows
func (p *Systray) SetTitle(_ string) {}

func New() (*Systray, error) {
	ni := &Systray{
		lclick: func() {},
		rclick: func() {},
	}

	MainClassName := "WailsSystray"
	ni.hinst, _ = RegisterWindow(MainClassName, ni.WinProc)

	mhwnd := win32.CreateWindowEx(
		win32.WS_EX_CONTROLPARENT,
		win32.MustStringToUTF16Ptr(MainClassName),
		win32.MustStringToUTF16Ptr(""),
		win32.WS_OVERLAPPEDWINDOW|win32.WS_CLIPSIBLINGS,
		win32.CW_USEDEFAULT,
		win32.CW_USEDEFAULT,
		win32.CW_USEDEFAULT,
		win32.CW_USEDEFAULT,
		0,
		0,
		0,
		unsafe.Pointer(nil))

	if mhwnd == 0 {
		return nil, errors.New("create main win failed")
	}

	NotifyIconClassName := "NotifyIconForm"
	_, err := RegisterWindow(NotifyIconClassName, ni.WinProc)
	if err != nil {
		return nil, err
	}

	hwnd, _, _ := CreateWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(win32.MustStringToUTF16Ptr(NotifyIconClassName))),
		0,
		0,
		0,
		0,
		0,
		0,
		uintptr(win32.HWND_MESSAGE),
		0,
		0,
		0)
	if hwnd == 0 {
		return nil, errors.New("create notify win failed")
	}

	ni.hwnd = win32.HWND(hwnd) // Important to keep this inside struct.

	nid := win32.NOTIFYICONDATA{
		HWnd:             win32.HWND(hwnd),
		UFlags:           win32.NIF_MESSAGE | win32.NIF_STATE,
		DwState:          win32.NIS_HIDDEN,
		DwStateMask:      win32.NIS_HIDDEN,
		UCallbackMessage: win32.NotifyIconMessageId,
	}
	nid.CbSize = uint32(unsafe.Sizeof(nid))

	ret := win32.ShellNotifyIcon(win32.NIM_ADD, &nid)
	if ret == 0 {
		return nil, errors.New("shell notify create failed")
	}

	nid.UVersion = win32.NOTIFYICON_VERSION

	ret = win32.ShellNotifyIcon(win32.NIM_SETVERSION, &nid)
	if ret == 0 {
		return nil, errors.New("shell notify version failed")
	}

	ni.appIcon = win32.LoadIconWithResourceID(0, uintptr(win32.IDI_APPLICATION))
	ni.lightModeIcon = ni.appIcon
	ni.darkModeIcon = ni.appIcon
	ni.id = nid.UID
	return ni, nil
}

func (p *Systray) HWND() win32.HWND {
	return p.hwnd
}

// AppendMenu add menu item.
func (p *Systray) AppendMenu(label string, onclick menu.Callback) {
	p.Menu = append(p.Menu, &menu.MenuItem{Type: menu.TextType, Label: label, Click: onclick})
}

// AppendMenuItem add menu item.
func (p *Systray) AppendMenuItem(item *menu.MenuItem) {
	p.Menu = append(p.Menu, item)
}

// AppendSeparator to the menu.
func (p *Systray) AppendSeparator() {
	p.Menu = append(p.Menu, menu.Separator())
}

func (p *Systray) Stop() error {
	nid := p.newNotifyIconData()
	ret := win32.ShellNotifyIcon(win32.NIM_DELETE, &nid)
	if ret == 0 {
		return errors.New("shell notify delete failed")
	}
	return nil
}

func (p *Systray) Click(fn func()) {
	p.lclick = fn
}

func (p *Systray) OnRightClick(fn func()) {
	p.rclick = fn
}

func (p *Systray) SetTooltip(tooltip string) error {
	nid := p.newNotifyIconData()
	nid.UFlags = win32.NIF_TIP
	copy(nid.SzTip[:], win32.MustUTF16FromString(tooltip))

	ret := win32.ShellNotifyIcon(win32.NIM_MODIFY, &nid)
	if ret == 0 {
		return errors.New("shell notify tooltip failed")
	}
	return nil
}

func (p *Systray) ShowMessage(title, msg string, bigIcon bool) error {
	nid := p.newNotifyIconData()
	if bigIcon == true {
		nid.DwInfoFlags = win32.NIIF_USER
	}

	nid.CbSize = uint32(unsafe.Sizeof(nid))

	nid.UFlags = win32.NIF_INFO
	copy(nid.SzInfoTitle[:], win32.MustUTF16FromString(title))
	copy(nid.SzInfo[:], win32.MustUTF16FromString(msg))

	ret := win32.ShellNotifyIcon(win32.NIM_MODIFY, &nid)
	if ret == 0 {
		return errors.New("shell notify tooltip failed")
	}
	return nil
}

func (p *Systray) newNotifyIconData() win32.NOTIFYICONDATA {
	nid := win32.NOTIFYICONDATA{
		UID:  p.id,
		HWnd: p.hwnd,
	}
	nid.CbSize = uint32(unsafe.Sizeof(nid))
	return nid
}

func (p *Systray) Show() error {
	return p.setVisible(true)
}

func (p *Systray) Hide() error {
	return p.setVisible(false)
}

func (p *Systray) setVisible(visible bool) error {
	nid := p.newNotifyIconData()
	nid.UFlags = win32.NIF_STATE
	nid.DwStateMask = win32.NIS_HIDDEN
	if !visible {
		nid.DwState = win32.NIS_HIDDEN
	}

	ret := win32.ShellNotifyIcon(win32.NIM_MODIFY, &nid)
	if ret == 0 {
		return errors.New("shell notify tooltip failed")
	}
	return nil
}

func (p *Systray) SetIcons(lightModeIcon, darkModeIcon *options.SystemTrayIcon) error {
	var newLightModeIcon, newDarkModeIcon win32.HICON
	if lightModeIcon != nil && lightModeIcon.Data != nil {
		newLightModeIcon = p.getIcon(lightModeIcon.Data)
	}
	if darkModeIcon != nil && darkModeIcon.Data != nil {
		newDarkModeIcon = p.getIcon(darkModeIcon.Data)
	}
	p.lightModeIcon, _ = lo.Coalesce(newLightModeIcon, newDarkModeIcon, p.appIcon)
	p.darkModeIcon, _ = lo.Coalesce(newDarkModeIcon, newLightModeIcon, p.appIcon)
	return p.updateIcon()
}

func (p *Systray) getIcon(icon []byte) win32.HICON {
	result, err := win32.CreateHIconFromPNG(icon)
	if err != nil {
		result = p.appIcon
	}
	return result
}

func (p *Systray) setIcon(hicon win32.HICON) error {
	nid := p.newNotifyIconData()
	nid.UFlags = win32.NIF_ICON
	if hicon == 0 {
		nid.HIcon = 0
	} else {
		nid.HIcon = hicon
	}

	ret := win32.ShellNotifyIcon(win32.NIM_MODIFY, &nid)
	if ret == 0 {
		return errors.New("shell notify icon failed")
	}
	return nil
}

func (p *Systray) WinProc(hwnd win32.HWND, msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case win32.NotifyIconMessageId:
		if lparam == win32.WM_LBUTTONUP {
			p.lclick()
			if len(p.Menu) > 0 {
				err := displayMenu(p.hwnd, p.Menu)
				if err != nil {
					return 0
				}
			}
		} else if lparam == win32.WM_RBUTTONUP {
			p.rclick()
			if len(p.Menu) > 0 {
				err := displayMenu(p.hwnd, p.Menu)
				if err != nil {
					return 0
				}
			}
		}
	case win32.WM_SETTINGCHANGE:
		settingChanged := win32.UTF16PtrToString(lparam)
		if settingChanged == "ImmersiveColorSet" {
			err := p.updateIcon()
			if err != nil {
				println("update icon failed", err.Error())
			}
		}
		return 0
	case win32.WM_COMMAND:
		cmdMsgID := int(wparam & 0xffff)
		switch cmdMsgID {
		default:
			if cmdMsgID >= win32.MenuItemMsgID && cmdMsgID < (win32.MenuItemMsgID+len(p.Menu)) {
				itemIndex := cmdMsgID - win32.MenuItemMsgID
				menuItem := p.Menu[itemIndex]
				menuItem.Click(nil)
			}
		}
	}

	result, _, _ := DefWindowProc.Call(uintptr(hwnd), uintptr(msg), wparam, lparam)
	return result
}

func (p *Systray) Run() error {
	var msg win32.MSG
	for {
		rt := win32.GetMessage(&msg)
		switch int(rt) {
		case 0:
			println("Quitting Run()")
			return nil
		case -1:
			return errors.New("run failed")
		}

		if win32.IsDialogMessage(p.hwnd, &msg) == 0 {
			win32.TranslateMessage(&msg)
			win32.DispatchMessage(&msg)
		}
	}
}

func (p *Systray) updateIcon() error {

	var newIcon win32.HICON
	if win32.IsCurrentlyDarkMode() {
		newIcon = p.darkModeIcon
	} else {
		newIcon = p.lightModeIcon
	}
	if p.currentIcon == newIcon {
		return nil
	}
	p.currentIcon = newIcon
	return p.setIcon(newIcon)
}

func RegisterWindow(name string, proc win32.WindowProc) (win32.HINSTANCE, error) {
	instance, exists := windowClasses[name]
	if exists {
		return instance, nil
	}
	hinst := win32.GetModuleHandle(0)
	if hinst == 0 {
		return 0, errors.New("get module handle failed")
	}
	hicon := win32.LoadIconWithResourceID(0, uintptr(win32.IDI_APPLICATION))
	if hicon == 0 {
		return 0, errors.New("load icon failed")
	}
	hcursor := win32.LoadCursorWithResourceID(0, uintptr(win32.IDC_ARROW))
	if hcursor == 0 {
		return 0, errors.New("load cursor failed")
	}

	hi := win32.HINSTANCE(hinst)

	var wc win32.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.LpfnWndProc = syscall.NewCallback(proc)
	wc.HInstance = win32.HINSTANCE(hinst)
	wc.HIcon = hicon
	wc.HCursor = hcursor
	wc.HbrBackground = win32.COLOR_BTNFACE + 1
	wc.LpszClassName = win32.MustStringToUTF16Ptr(name)

	atom, _, e := RegisterClassEx.Call(uintptr(unsafe.Pointer(&wc)))
	if atom == 0 {
		println(e.Error())
		return 0, errors.New("register class failed")
	}

	windowClasses[name] = hi
	return hi, nil
}
