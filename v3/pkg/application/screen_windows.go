//go:build windows

package application

import (
	"errors"
	"strconv"

	"github.com/wailsapp/wails/v3/pkg/w32"
	"golang.org/x/sys/windows"
)

func (m *windowsApp) processAndCacheScreens() error {
	allScreens, err := w32.GetAllScreens()
	if err != nil {
		return err
	}

	// Convert result to []*Screen
	var screens []*Screen

	for _, screen := range allScreens {
		x := int(screen.MONITORINFOEX.RcMonitor.Left)
		y := int(screen.MONITORINFOEX.RcMonitor.Top)
		right := int(screen.MONITORINFOEX.RcMonitor.Right)
		bottom := int(screen.MONITORINFOEX.RcMonitor.Bottom)
		width := right - x
		height := bottom - y

		workArea := Rect{
			X:      int(screen.MONITORINFOEX.RcWork.Left),
			Y:      int(screen.MONITORINFOEX.RcWork.Top),
			Width:  int(screen.MONITORINFOEX.RcWork.Right - screen.MONITORINFOEX.RcWork.Left),
			Height: int(screen.MONITORINFOEX.RcWork.Bottom - screen.MONITORINFOEX.RcWork.Top),
		}

		screens = append(screens, &Screen{
			ID:               hMonitorToScreenID(screen.HMonitor),
			Name:             windows.UTF16ToString(screen.MONITORINFOEX.SzDevice[:]),
			X:                x,
			Y:                y,
			Size:             Size{Width: width, Height: height},
			Bounds:           Rect{X: x, Y: y, Width: width, Height: height},
			PhysicalBounds:   Rect{X: x, Y: y, Width: width, Height: height},
			WorkArea:         workArea,
			PhysicalWorkArea: workArea,
			IsPrimary:        screen.IsPrimary,
			ScaleFactor:      screen.ScaleFactor,
			Rotation:         0,
		})
	}

	err = m.parent.screenManager.LayoutScreens(screens)
	if err != nil {
		return err
	}

	return nil
}

// NOTE: should be moved to *App after DPI is implemented in all platforms
func (m *windowsApp) getScreens() ([]*Screen, error) {
	return m.parent.screenManager.screens, nil
}

// NOTE: should be moved to *App after DPI is implemented in all platforms
func (m *windowsApp) getPrimaryScreen() (*Screen, error) {
	return m.parent.screenManager.primaryScreen, nil
}

func getScreenForWindow(window *windowsWebviewWindow) (*Screen, error) {
	return ScreenNearestPhysicalRect(window.physicalBounds()), nil
}

func getScreenForWindowHwnd(hwnd w32.HWND) (*Screen, error) {
	hMonitor := w32.MonitorFromWindow(hwnd, w32.MONITOR_DEFAULTTONEAREST)
	screenID := hMonitorToScreenID(hMonitor)
	for _, screen := range globalApplication.screenManager.screens {
		if screen.ID == screenID {
			return screen, nil
		}
	}
	return nil, errors.New("screen not found for window")
}

func hMonitorToScreenID(hMonitor uintptr) string {
	return strconv.Itoa(int(hMonitor))
}
