//go:build windows

package windows

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"github.com/leaanthony/winc"
	"github.com/leaanthony/winc/w32"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
)

type Window struct {
	winc.Form
	frontendOptions                   *options.App
	applicationMenu                   *menu.Menu
	notifyParentWindowPositionChanged func() error
}

func NewWindow(parent winc.Controller, appoptions *options.App) *Window {
	result := new(Window)
	result.frontendOptions = appoptions
	result.SetIsForm(true)

	var exStyle int
	if appoptions.Windows != nil {
		exStyle = w32.WS_EX_CONTROLPARENT | w32.WS_EX_APPWINDOW
		if appoptions.Windows.WindowIsTranslucent {
			exStyle |= w32.WS_EX_NOREDIRECTIONBITMAP
		}
	}
	if appoptions.AlwaysOnTop {
		exStyle |= w32.WS_EX_TOPMOST
	}

	var dwStyle = w32.WS_OVERLAPPEDWINDOW

	winc.RegClassOnlyOnce("wailsWindow")
	result.SetHandle(winc.CreateWindow("wailsWindow", parent, uint(exStyle), uint(dwStyle)))
	winc.RegMsgHandler(result)
	result.SetParent(parent)

	loadIcon := true
	if appoptions.Windows != nil && appoptions.Windows.DisableWindowIcon == true {
		loadIcon = false
	}
	if loadIcon {
		if ico, err := winc.NewIconFromResource(winc.GetAppInstance(), uint16(winc.AppIconID)); err == nil {
			result.SetIcon(0, ico)
		}
	}

	result.SetSize(appoptions.Width, appoptions.Height)
	result.SetText(appoptions.Title)
	result.EnableSizable(!appoptions.DisableResize)
	if !appoptions.Fullscreen {
		result.EnableMaxButton(!appoptions.DisableResize)
		result.SetMinSize(appoptions.MinWidth, appoptions.MinHeight)
		result.SetMaxSize(appoptions.MaxWidth, appoptions.MaxHeight)
	}

	if appoptions.Windows != nil {
		if appoptions.Windows.WindowIsTranslucent {
			result.SetTranslucentBackground()
		}

		if appoptions.Windows.DisableWindowIcon {
			result.DisableIcon()
		}
	}

	// Dlg forces display of focus rectangles, as soon as the user starts to type.
	w32.SendMessage(result.Handle(), w32.WM_CHANGEUISTATE, w32.UIS_INITIALIZE, 0)

	result.SetFont(winc.DefaultFont)

	if appoptions.Menu != nil {
		result.SetApplicationMenu(appoptions.Menu)
	}

	return result
}

func (w *Window) Run() int {
	return winc.RunMainLoop()
}

func (w *Window) WndProc(msg uint32, wparam, lparam uintptr) uintptr {

	switch msg {
	case w32.WM_NCLBUTTONDOWN:
		w32.SetFocus(w.Handle())
	case w32.WM_MOVE, w32.WM_MOVING:
		if w.notifyParentWindowPositionChanged != nil {
			w.notifyParentWindowPositionChanged()
		}

	// TODO move WM_DPICHANGED handling into winc
	case 0x02E0: //w32.WM_DPICHANGED
		newWindowSize := (*w32.RECT)(unsafe.Pointer(lparam))
		w32.SetWindowPos(w.Handle(),
			uintptr(0),
			int(newWindowSize.Left),
			int(newWindowSize.Top),
			int(newWindowSize.Right-newWindowSize.Left),
			int(newWindowSize.Bottom-newWindowSize.Top),
			w32.SWP_NOZORDER|w32.SWP_NOACTIVATE)
	}

	if w.frontendOptions.Frameless {
		switch msg {
		case w32.WM_ACTIVATE:
			// If we want to have a frameless window but with border, extend the client area outside of the window. This
			// Option is not affected by returning 0 in WM_NCCALCSIZE.
			// As a result we have hidden the titlebar but still have the default border drawn.
			// See: https://docs.microsoft.com/en-us/windows/win32/api/dwmapi/nf-dwmapi-dwmextendframeintoclientarea#remarks
			if winoptions := w.frontendOptions.Windows; winoptions != nil && winoptions.EnableFramelessBorder {
				if err := dwmExtendFrameIntoClientArea(w.Handle(), w32.MARGINS{-1, -1, -1, -1}); err != nil {
					log.Fatal(fmt.Errorf("DwmExtendFrameIntoClientArea failed: %s", err))
				}
			}
		case w32.WM_NCCALCSIZE:
			// Disable the standard frame by allowing the client area to take the full
			// window size.
			// See: https://docs.microsoft.com/en-us/windows/win32/winmsg/wm-nccalcsize#remarks
			// This hides the titlebar and also disables the resizing from user interaction because the standard frame is not
			// shown. We still need the WS_THICKFRAME style to enable resizing from the frontend.
			if wparam != 0 {
				style := uint32(w32.GetWindowLong(w.Handle(), w32.GWL_STYLE))
				if style&w32.WS_MAXIMIZE != 0 {
					// If the window is maximized we must adjust the client area to the work area of the monitor. Otherwise
					// some content goes beyond the visible part of the monitor.
					monitor := w32.MonitorFromWindow(w.Handle(), w32.MONITOR_DEFAULTTONEAREST)

					var monitorInfo w32.MONITORINFO
					monitorInfo.CbSize = uint32(unsafe.Sizeof(monitorInfo))
					if w32.GetMonitorInfo(monitor, &monitorInfo) {
						rgrc := (*w32.RECT)(unsafe.Pointer(lparam))
						*rgrc = monitorInfo.RcWork
					}
				}

				return 0
			}
		}
	}
	return w.Form.WndProc(msg, wparam, lparam)
}

// TODO this should be put into the winc if we are happy with this solution.
var (
	modkernel32 = syscall.NewLazyDLL("dwmapi.dll")

	procDwmExtendFrameIntoClientArea = modkernel32.NewProc("DwmExtendFrameIntoClientArea")
)

func dwmExtendFrameIntoClientArea(hwnd w32.HWND, margins w32.MARGINS) error {
	ret, _, _ := procDwmExtendFrameIntoClientArea.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&margins)))

	if ret != 0 {
		return syscall.GetLastError()
	}

	return nil
}
