//go:build windows && !server

package application

import (
	"errors"
	"strconv"
	"time"

	"github.com/wailsapp/wails/v3/pkg/w32"
	"golang.org/x/sys/windows"
)

const (
	// screenCacheStartupAttempts is how many times the screens are read at
	// startup before a failure is treated as fatal.
	screenCacheStartupAttempts = 3
	// screenCacheRetryDelay is the pause before reading the screens again
	// after a failure. Display failures are usually transient: a monitor
	// that is being reconfigured or removed, or a session switch.
	screenCacheRetryDelay = 250 * time.Millisecond
)

// retryScreenCache calls cache up to attempts times, pausing for delay between
// attempts, and returns nil as soon as one attempt succeeds. If every attempt
// fails it returns the error of the last one.
func retryScreenCache(attempts int, delay time.Duration, cache func() error) error {
	var err error
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 {
			time.Sleep(delay)
		}
		if err = cache(); err == nil {
			return nil
		}
	}
	return err
}

// retryCacheScreensLater schedules one more attempt to read the screens after
// a failed re-read on a display change. Without it the previous layout would
// stay cached until the next display message, which might never come if the
// failed read was the last message of the change. Failures that arrive while an
// attempt is already scheduled share it.
func (m *windowsApp) retryCacheScreensLater() {
	if !m.screenRetryPending.CompareAndSwap(false, true) {
		return
	}
	time.AfterFunc(screenCacheRetryDelay, func() {
		m.screenRetryPending.Store(false)
		if err := m.processAndCacheScreens(); err != nil {
			m.parent.handleError(err)
		}
	})
}

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

	err = m.parent.Screen.LayoutScreens(screens)
	if err != nil {
		return err
	}

	return nil
}

// NOTE: should be moved to *App after DPI is implemented in all platforms
func (m *windowsApp) getScreens() ([]*Screen, error) {
	return m.parent.Screen.GetAll(), nil
}

// NOTE: should be moved to *App after DPI is implemented in all platforms
func (m *windowsApp) getPrimaryScreen() (*Screen, error) {
	return m.parent.Screen.GetPrimary(), nil
}

func getScreenForWindow(window *windowsWebviewWindow) (*Screen, error) {
	return ScreenNearestPhysicalRect(window.physicalBounds()), nil
}

func getScreenForWindowHwnd(hwnd w32.HWND) (*Screen, error) {
	hMonitor := w32.MonitorFromWindow(hwnd, w32.MONITOR_DEFAULTTONEAREST)
	screenID := hMonitorToScreenID(hMonitor)
	for _, screen := range globalApplication.Screen.GetAll() {
		if screen.ID == screenID {
			return screen, nil
		}
	}
	return nil, errors.New("screen not found for window")
}

func hMonitorToScreenID(hMonitor uintptr) string {
	return strconv.Itoa(int(hMonitor))
}
