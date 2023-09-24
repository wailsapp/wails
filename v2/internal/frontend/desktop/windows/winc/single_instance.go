//go:build windows
// +build windows

package winc

import (
	"encoding/json"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
	"github.com/wailsapp/wails/v2/pkg/options"
	"golang.org/x/sys/windows"
	"os"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

type COPYDATASTRUCT struct {
	dwData uintptr
	cbData uint32
	lpData uintptr
}

var mainHWND w32.HWND

const WMCOPYDATA_SINGLE_INSTANCE_DATA = 1542
const SC_RESTORE = 0xF120

func SendMessage(hWnd w32.HWND, data string) {
	arrUtf16, _ := syscall.UTF16FromString(data)

	pCopyData := new(COPYDATASTRUCT)
	pCopyData.dwData = WMCOPYDATA_SINGLE_INSTANCE_DATA
	pCopyData.cbData = uint32(len(arrUtf16)*2 + 1)
	pCopyData.lpData = uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(data)))

	w32.SendMessage(hWnd, w32.WM_COPYDATA, 0, uintptr(unsafe.Pointer(pCopyData)))
}

// single instance windows app
func SetupSingleInstance(title string, lock *options.SingleInstanceLock) {
	// TODO: create unique id based on title. Maybe better to have something more unique?
	id := "wails-app-" + title

	className := id + "-sic"
	windowName := id + "-siw"
	mutexName := id + "sim"

	_, err := windows.CreateMutex(nil, false, windows.StringToUTF16Ptr(mutexName))

	if err != nil {
		if err == windows.ERROR_ALREADY_EXISTS {
			// app is already running
			hwnd := w32.FindWindowW(windows.StringToUTF16Ptr(className), windows.StringToUTF16Ptr(windowName))

			if hwnd != 0 {
				data := options.SecondInstanceData{
					Args: os.Args[1:],
				}
				serialized, _ := json.Marshal(data)

				SendMessage(hwnd, string(serialized))
				// exit second instance of app after sending message
				os.Exit(0)
			}

			panic("unknown error on mutex")
		}
	}

	createEventTargetWindow(className, windowName, lock)
}

func SingleInstanceHWND(hwnd w32.HWND) {
	mainHWND = hwnd
}

func createEventTargetWindow(className string, windowName string, lock *options.SingleInstanceLock) w32.HWND {
	// callback handler in the event target window
	wndProc := func(
		hwnd w32.HWND, msg uint32, wparam w32.WPARAM, lparam w32.LPARAM,
	) w32.LRESULT {
		if msg == w32.WM_COPYDATA {
			ldata := (*COPYDATASTRUCT)(unsafe.Pointer(lparam))

			if ldata.dwData == WMCOPYDATA_SINGLE_INSTANCE_DATA {
				serialized := uintptrToString(ldata.lpData)

				var secondInstanceData options.SecondInstanceData

				json.Unmarshal([]byte(serialized), &secondInstanceData)

				go lock.OnSecondInstanceLaunch(secondInstanceData)

				if lock.ActivateAppOnSubsequentLaunch && mainHWND != 0 {
					w32.SendMessage(mainHWND, w32.WM_SYSCOMMAND, SC_RESTORE, 0)                                                  // restore the minimize window
					w32.SetWindowPos(mainHWND, w32.HWND_TOPMOST, 0, 0, 0, 0, w32.SWP_SHOWWINDOW|w32.SWP_NOSIZE|w32.SWP_NOMOVE)   // force set our main window on top
					w32.SetWindowPos(mainHWND, w32.HWND_NOTOPMOST, 0, 0, 0, 0, w32.SWP_SHOWWINDOW|w32.SWP_NOSIZE|w32.SWP_NOMOVE) // remove topmost to allow normal windows manipulations
					w32.SetForegroundWindow(mainHWND)                                                                            // put window on tops
				}
			}

			return w32.LRESULT(0)
		}

		return w32.DefWindowProc(hwnd, msg, wparam, lparam)
	}

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
	class.ClassName = windows.StringToUTF16Ptr(className)
	class.IconSm = 0

	w32.RegisterClassEx(&class)

	// create event window that will not be visible for user
	hwnd := w32.CreateWindowEx(
		w32.WS_EX_NOACTIVATE|
			w32.WS_EX_TRANSPARENT|
			w32.WS_EX_LAYERED|
			// WS_EX_TOOLWINDOW prevents this window from ever showing up in the taskbar, which
			// we want to avoid. If you remove this style, this window won't show up in the
			// taskbar *initially*, but it can show up at some later point. This can sometimes
			// happen on its own after several hours have passed, although this has proven
			// difficult to reproduce. Alternatively, it can be manually triggered by killing
			// `explorer.exe` and then starting the process back up.
			// It is unclear why the bug is triggered by waiting for several hours.
			w32.WS_EX_TOOLWINDOW,
		windows.StringToUTF16Ptr(className),
		windows.StringToUTF16Ptr(windowName),
		w32.WS_OVERLAPPED,
		0,
		0,
		0,
		0,
		0,
		0,
		w32.GetModuleHandle(""),
		nil,
	)

	w32.SetWindowLongPtr(
		hwnd,
		w32.GWL_STYLE,
		// The window technically has to be visible to receive WM_PAINT messages (which are used
		// for delivering events during resizes), but it isn't displayed to the user because of
		// the LAYERED style.
		w32.WS_VISIBLE|w32.WS_POPUP,
	)

	return hwnd
}

func uintptrToString(cstr uintptr) string {
	if cstr != 0 {
		us := make([]uint16, 0, 256)
		for p := cstr; ; p += 2 {
			u := *(*uint16)(unsafe.Pointer(p))
			if u == 0 {
				return string(utf16.Decode(us))
			}
			us = append(us, u)
		}
	}
	return ""
}
