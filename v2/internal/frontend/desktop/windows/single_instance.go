//go:build windows

package windows

import (
	"encoding/json"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
	"github.com/wailsapp/wails/v2/pkg/options"
	"golang.org/x/sys/windows"
	"log"
	"os"
	"syscall"
	"unsafe"
)

type COPYDATASTRUCT struct {
	dwData uintptr
	cbData uint32
	lpData uintptr
}

// WMCOPYDATA_SINGLE_INSTANCE_DATA we define our own type for WM_COPYDATA message
const WMCOPYDATA_SINGLE_INSTANCE_DATA = 1542

func SendMessage(hwnd w32.HWND, data string) {
	arrUtf16, _ := syscall.UTF16FromString(data)

	pCopyData := new(COPYDATASTRUCT)
	pCopyData.dwData = WMCOPYDATA_SINGLE_INSTANCE_DATA
	pCopyData.cbData = uint32(len(arrUtf16)*2 + 1)
	pCopyData.lpData = uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(data)))

	w32.SendMessage(hwnd, w32.WM_COPYDATA, 0, uintptr(unsafe.Pointer(pCopyData)))
}

// SetupSingleInstance single instance Windows app
func SetupSingleInstance(uniqueId string) {
	id := "wails-app-" + uniqueId

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
				data.WorkingDirectory, err = os.Getwd()
				if err != nil {
					log.Printf("Failed to get working directory: %v", err)
					return
				}
				serialized, err := json.Marshal(data)
				if err != nil {
					log.Printf("Failed to marshal data: %v", err)
					return
				}

				SendMessage(hwnd, string(serialized))
				// exit second instance of app after sending message
				os.Exit(0)
			}
			// if we got any other unknown error we will just start new application instance
		}
	} else {
		createEventTargetWindow(className, windowName)
	}
}

func createEventTargetWindow(className string, windowName string) w32.HWND {
	// callback handler in the event target window
	wndProc := func(
		hwnd w32.HWND, msg uint32, wparam w32.WPARAM, lparam w32.LPARAM,
	) w32.LRESULT {
		if msg == w32.WM_COPYDATA {
			ldata := (*COPYDATASTRUCT)(unsafe.Pointer(lparam))

			if ldata.dwData == WMCOPYDATA_SINGLE_INSTANCE_DATA {
				serialized := windows.UTF16PtrToString((*uint16)(unsafe.Pointer(ldata.lpData)))

				var secondInstanceData options.SecondInstanceData

				err := json.Unmarshal([]byte(serialized), &secondInstanceData)

				if err == nil {
					secondInstanceBuffer <- secondInstanceData
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
		0,
		windows.StringToUTF16Ptr(className),
		windows.StringToUTF16Ptr(windowName),
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
