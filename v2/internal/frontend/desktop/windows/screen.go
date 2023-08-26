//go:build windows
// +build windows

package windows

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
)

func MonitorsEqual(first w32.MONITORINFO, second w32.MONITORINFO) bool {
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

func GetMonitorInfo(hMonitor w32.HMONITOR) (*w32.MONITORINFO, error) {
	// Adapted from winc.utils.getMonitorInfo TODO: add this to win32
	// See docs for
	//https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getmonitorinfoa

	var info w32.MONITORINFO
	info.CbSize = uint32(unsafe.Sizeof(info))
	succeeded := w32.GetMonitorInfo(hMonitor, &info)
	if !succeeded {
		return &info, errors.New("Windows call to getMonitorInfo failed")
	}
	return &info, nil
}

func EnumProc(hMonitor w32.HMONITOR, hdcMonitor w32.HDC, lprcMonitor *w32.RECT, screenContainer *ScreenContainer) uintptr {
	// adapted from https://stackoverflow.com/a/23492886/4188138

	// see docs for the following pages to better understand this function
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-enumdisplaymonitors
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nc-winuser-monitorenumproc
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-monitorinfo
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-monitorfromwindow

	ourMonitorData := Screen{}
	currentMonHndl := w32.MonitorFromWindow(screenContainer.mainWinHandle, w32.MONITOR_DEFAULTTONEAREST)
	currentMonInfo, currErr := GetMonitorInfo(currentMonHndl)

	if currErr != nil {
		screenContainer.errors = append(screenContainer.errors, currErr)
		screenContainer.monitors = append(screenContainer.monitors, Screen{})
		// not sure what the consequences of returning false are, so let's just return true and handle it ourselves
		return w32.TRUE
	}

	monInfo, err := GetMonitorInfo(hMonitor)
	if err != nil {
		screenContainer.errors = append(screenContainer.errors, err)
		screenContainer.monitors = append(screenContainer.monitors, Screen{})
		return w32.TRUE
	}

	width := lprcMonitor.Right - lprcMonitor.Left
	height := lprcMonitor.Bottom - lprcMonitor.Top
	ourMonitorData.IsPrimary = monInfo.DwFlags&w32.MONITORINFOF_PRIMARY == 1
	ourMonitorData.Height = int(height)
	ourMonitorData.Width = int(width)
	ourMonitorData.IsCurrent = MonitorsEqual(*currentMonInfo, *monInfo)

	ourMonitorData.PhysicalSize.Width = int(width)
	ourMonitorData.PhysicalSize.Height = int(height)

	var dpiX, dpiY uint
	w32.GetDPIForMonitor(hMonitor, w32.MDT_EFFECTIVE_DPI, &dpiX, &dpiY)
	if dpiX == 0 || dpiY == 0 {
		screenContainer.errors = append(screenContainer.errors, fmt.Errorf("unable to get DPI for screen"))
		screenContainer.monitors = append(screenContainer.monitors, Screen{})
		return w32.TRUE
	}
	ourMonitorData.Size.Width = winc.ScaleToDefaultDPI(ourMonitorData.PhysicalSize.Width, dpiX)
	ourMonitorData.Size.Height = winc.ScaleToDefaultDPI(ourMonitorData.PhysicalSize.Height, dpiY)

	// the reason we need a container is that we have don't know how many times this function will be called
	// this "append" call could potentially do an allocation and rewrite the pointer to monitors. So we save the pointer in screenContainer.monitors
	// and retrieve the values after all EnumProc calls
	// If EnumProc is multi-threaded, this could be problematic. Although, I don't think it is.
	screenContainer.monitors = append(screenContainer.monitors, ourMonitorData)
	// let's keep screenContainer.errors the same size as screenContainer.monitors in case we want to match them up later if necessary
	screenContainer.errors = append(screenContainer.errors, nil)
	return w32.TRUE
}

type ScreenContainer struct {
	monitors      []Screen
	errors        []error
	mainWinHandle w32.HWND
}

func GetAllScreens(mainWinHandle w32.HWND) ([]Screen, error) {
	// TODO fix hack of container sharing by having a proper data sharing mechanism between windows and the runtime
	monitorContainer := ScreenContainer{mainWinHandle: mainWinHandle}
	returnErr := error(nil)
	errorStrings := []string{}

	dc := w32.GetDC(0)
	defer w32.ReleaseDC(0, dc)
	succeeded := w32.EnumDisplayMonitors(dc, nil, syscall.NewCallback(EnumProc), unsafe.Pointer(&monitorContainer))
	if !succeeded {
		return monitorContainer.monitors, errors.New("Windows call to EnumDisplayMonitors failed")
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
