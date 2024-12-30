//go:build windows

package application

import (
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/w32"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")
)

type windowsLock struct {
	handle     syscall.Handle
	uniqueID   string
	msgString  string
	hwnd       w32.HWND
	manager    *singleInstanceManager
	className  string
	windowName string
}

func newPlatformLock(manager *singleInstanceManager) (platformLock, error) {
	return &windowsLock{
		manager: manager,
	}, nil
}

func (l *windowsLock) acquire(uniqueID string) error {
	if uniqueID == "" {
		return fmt.Errorf("UniqueID is required for single instance lock")
	}

	l.uniqueID = uniqueID
	id := "wails-app-" + uniqueID
	l.className = id + "-sic"
	l.windowName = id + "-siw"
	mutexName := id + "-sim"

	_, err := windows.CreateMutex(nil, false, windows.StringToUTF16Ptr(mutexName))
	if err != nil {
		// Find the window
		return alreadyRunningError
	} else {
		l.hwnd = createEventTargetWindow(l.className, l.windowName)
	}

	return nil
}

func (l *windowsLock) release() {
	if l.handle != 0 {
		syscall.CloseHandle(l.handle)
		l.handle = 0
	}
	if l.hwnd != 0 {
		w32.DestroyWindow(l.hwnd)
		l.hwnd = 0
	}
}

func (l *windowsLock) notify(data string) error {

	// app is already running
	hwnd := w32.FindWindowW(windows.StringToUTF16Ptr(l.className), windows.StringToUTF16Ptr(l.windowName))

	if hwnd == 0 {
		return errors.New("unable to notify other instance")
	}

	w32.SendMessageToWindow(hwnd, data)

	return nil
}

func createEventTargetWindow(className string, windowName string) w32.HWND {
	var class w32.WNDCLASSEX
	class.Size = uint32(unsafe.Sizeof(class))
	class.Style = 0
	class.WndProc = syscall.NewCallback(wndProc)
	class.ClsExtra = 0
	class.WndExtra = 0
	class.Instance = w32.GetModuleHandle("")
	class.Icon = 0
	class.Cursor = 0
	class.Background = 0
	class.MenuName = nil
	class.ClassName = w32.MustStringToUTF16Ptr(className)
	class.IconSm = 0

	w32.RegisterClassEx(&class)

	// Create hidden message-only window
	hwnd := w32.CreateWindowEx(
		0,
		w32.MustStringToUTF16Ptr(className),
		w32.MustStringToUTF16Ptr(windowName),
		0,
		0,
		0,
		0,
		0,
		w32.HWND_MESSAGE,
		0,
		w32.GetModuleHandle(""),
		nil,
	)

	return hwnd
}

func wndProc(hwnd w32.HWND, msg uint32, wparam w32.WPARAM, lparam w32.LPARAM) w32.LRESULT {
	if msg == w32.WM_COPYDATA {
		ldata := (*w32.COPYDATASTRUCT)(unsafe.Pointer(lparam))

		if ldata.DwData == w32.WMCOPYDATA_SINGLE_INSTANCE_DATA {
			serialized := windows.UTF16PtrToString((*uint16)(unsafe.Pointer(ldata.LpData)))
			secondInstanceBuffer <- serialized
		}
		return w32.LRESULT(0)
	}

	return w32.DefWindowProc(hwnd, msg, wparam, lparam)
}
