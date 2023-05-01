//go:build windows

package w32

import (
	"fmt"
	"syscall"
	"unsafe"
)

func MonitorsEqual(first MONITORINFO, second MONITORINFO) bool {
	// Checks to make sure all the fields are the same.
	// A cleaner way would be to check identity of devices. but I couldn't find a way of doing that using the win32 API
	return first.DwFlags == second.DwFlags &&
		first.RcMonitor.Top == second.RcMonitor.Top &&
		first.RcMonitor.Bottom == second.RcMonitor.Bottom &&
		first.RcMonitor.Right == second.RcMonitor.Right &&
		first.RcMonitor.Left == second.RcMonitor.Left &&
		first.RcWork.Top == second.RcWork.Top &&
		first.RcWork.Bottom == second.RcWork.Bottom &&
		first.RcWork.Right == second.RcWork.Right &&
		first.RcWork.Left == second.RcWork.Left
}

func GetMonitorInformation(hMonitor HMONITOR) (*MONITORINFO, error) {
	// Adapted from winc.utils.getMonitorInfo
	// See docs for
	//https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getmonitorinfoa

	var info MONITORINFO
	info.CbSize = uint32(unsafe.Sizeof(info))
	succeeded := GetMonitorInfo(hMonitor, &info)
	if !succeeded {
		return &info, fmt.Errorf("Windows call to getMonitorInfo failed")
	}
	return &info, nil
}

type Screen struct {
	IsCurrent bool
	IsPrimary bool
	Width     int
	Height    int
}

func EnumProc(hMonitor HMONITOR, hdcMonitor HDC, lprcMonitor *RECT, screenContainer *ScreenContainer) uintptr {
	// adapted from https://stackoverflow.com/a/23492886/4188138

	// see docs for the following pages to better understand this function
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-enumdisplaymonitors
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nc-winuser-monitorenumproc
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-monitorinfo
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-monitorfromwindow

	ourMonitorData := Screen{}
	currentMonHndl := MonitorFromWindow(screenContainer.mainWinHandle, MONITOR_DEFAULTTONEAREST)
	currentMonInfo, currErr := GetMonitorInformation(currentMonHndl)

	if currErr != nil {
		screenContainer.errors = append(screenContainer.errors, currErr)
		screenContainer.monitors = append(screenContainer.monitors, Screen{})
		// not sure what the consequences of returning false are, so let's just return true and handle it ourselves
		return TRUE
	}

	monInfo, err := GetMonitorInformation(hMonitor)
	if err != nil {
		screenContainer.errors = append(screenContainer.errors, err)
		screenContainer.monitors = append(screenContainer.monitors, Screen{})
		return TRUE
	}

	height := lprcMonitor.Right - lprcMonitor.Left
	width := lprcMonitor.Bottom - lprcMonitor.Top
	ourMonitorData.IsPrimary = monInfo.DwFlags&MONITORINFOF_PRIMARY == 1
	ourMonitorData.Height = int(width)
	ourMonitorData.Width = int(height)
	ourMonitorData.IsCurrent = MonitorsEqual(*currentMonInfo, *monInfo)

	// the reason we need a container is that we have don't know how many times this function will be called
	// this "append" call could potentially do an allocation and rewrite the pointer to monitors. So we save the pointer in screenContainer.monitors
	// and retrieve the values after all EnumProc calls
	// If EnumProc is multi-threaded, this could be problematic. Although, I don't think it is.
	screenContainer.monitors = append(screenContainer.monitors, ourMonitorData)
	// let's keep screenContainer.errors the same size as screenContainer.monitors in case we want to match them up later if necessary
	screenContainer.errors = append(screenContainer.errors, nil)
	return TRUE
}

type ScreenContainer struct {
	monitors      []Screen
	errors        []error
	mainWinHandle HWND
}

func GetAllScreens(mainWinHandle HWND) ([]Screen, error) {
	// TODO fix hack of container sharing by having a proper data sharing mechanism between windows and the runtime
	monitorContainer := ScreenContainer{mainWinHandle: mainWinHandle}
	returnErr := error(nil)
	var errorStrings []string

	dc := GetDC(0)
	defer ReleaseDC(0, dc)
	succeeded := EnumDisplayMonitors(dc, nil, syscall.NewCallback(EnumProc), unsafe.Pointer(&monitorContainer))
	if !succeeded {
		return monitorContainer.monitors, fmt.Errorf("Windows call to EnumDisplayMonitors failed")
	}
	for idx, err := range monitorContainer.errors {
		if err != nil {
			errorStrings = append(errorStrings, fmt.Sprintf("Error from monitor #%v, %v", idx+1, err))
		}
	}

	if len(errorStrings) > 0 {
		returnErr = fmt.Errorf("%v errors encountered: %v", len(errorStrings), errorStrings)
	}
	return monitorContainer.monitors, returnErr
}
