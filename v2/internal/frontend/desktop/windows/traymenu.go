//go:build windows
// +build windows

package windows

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/win32"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"log"
	"sync"
	"unsafe"
)

var uids uint32
var lock sync.RWMutex

func newUID() uint32 {
	lock.Lock()
	result := uids
	uids++
	lock.Unlock()
	return result

}

type Win32TrayMenu struct {
	hwnd uintptr
	uid  uint32
	icon uintptr
}

func (w *Win32TrayMenu) SetLabel(label string) {}

func (w *Win32TrayMenu) SetMenu(menu *menu.Menu) {}

func (w *Win32TrayMenu) SetImage(image *menu.TrayImage) {
	data := w.newNotifyIconData()
	bitmap := image.GetBestBitmap(1, false)
	icon, err := win32.CreateIconFromResourceEx(uintptr(unsafe.Pointer(&bitmap[0])), uint32(len(bitmap)), true, 0x30000, 0, 0, 0)
	if err != nil {
		log.Fatal(err.Error())
	}
	data.UFlags |= win32.NIF_ICON
	data.HIcon = icon
	if _, err := win32.NotifyIcon(win32.NIM_MODIFY, data); err != nil {
		log.Fatal(err.Error())
	}
}

func (f *Frontend) NewWin32TrayMenu(trayMenu *menu.TrayMenu) *Win32TrayMenu {

	result := &Win32TrayMenu{
		hwnd: f.mainWindow.Handle(),
		uid:  newUID(),
	}

	data := result.newNotifyIconData()
	data.UFlags |= win32.NIF_MESSAGE | win32.NIF_ICON
	data.UCallbackMessage = win32.WM_APP + result.uid
	if _, err := win32.NotifyIcon(win32.NIM_ADD, data); err != nil {
		log.Fatal(err.Error())
	}

	return result
}

func (w *Win32TrayMenu) newNotifyIconData() *win32.NOTIFYICONDATA {
	var data win32.NOTIFYICONDATA
	data.CbSize = uint32(unsafe.Sizeof(data))
	data.UFlags = win32.NIF_GUID
	data.HWnd = w.hwnd
	data.UID = w.uid
	return &data
}

func (f *Frontend) TrayMenuAdd(trayMenu *menu.TrayMenu) menu.TrayMenuImpl {
	win32TrayMenu := f.NewWin32TrayMenu(trayMenu)
	return win32TrayMenu
}
