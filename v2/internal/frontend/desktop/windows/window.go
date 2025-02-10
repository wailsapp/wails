//go:build windows

package windows

import (
	"github.com/wailsapp/go-webview2/pkg/edge"
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/win32"
	"github.com/wailsapp/wails/v2/internal/system/operatingsystem"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	winoptions "github.com/wailsapp/wails/v2/pkg/options/windows"
)

type Window struct {
	winc.Form
	frontendOptions                          *options.App
	applicationMenu                          *menu.Menu
	minWidth, minHeight, maxWidth, maxHeight int
	versionInfo                              *operatingsystem.WindowsVersionInfo
	isDarkMode                               bool
	isActive                                 bool
	hasBeenShown                             bool

	// Theme
	theme        winoptions.Theme
	themeChanged bool

	framelessWithDecorations bool

	OnSuspend func()
	OnResume  func()

	chromium *edge.Chromium
}

func NewWindow(parent winc.Controller, appoptions *options.App, versionInfo *operatingsystem.WindowsVersionInfo, chromium *edge.Chromium) *Window {
	windowsOptions := appoptions.Windows

	result := &Window{
		frontendOptions: appoptions,
		minHeight:       appoptions.MinHeight,
		minWidth:        appoptions.MinWidth,
		maxHeight:       appoptions.MaxHeight,
		maxWidth:        appoptions.MaxWidth,
		versionInfo:     versionInfo,
		isActive:        true,
		themeChanged:    true,
		chromium:        chromium,

		framelessWithDecorations: appoptions.Frameless && (windowsOptions == nil || !windowsOptions.DisableFramelessWindowDecorations),
	}
	result.SetIsForm(true)

	var exStyle int
	if windowsOptions != nil {
		exStyle = w32.WS_EX_CONTROLPARENT | w32.WS_EX_APPWINDOW
		if windowsOptions.WindowIsTranslucent {
			exStyle |= w32.WS_EX_NOREDIRECTIONBITMAP
		}
	}
	if appoptions.AlwaysOnTop {
		exStyle |= w32.WS_EX_TOPMOST
	}

	var dwStyle = w32.WS_OVERLAPPEDWINDOW

	winc.RegClassOnlyOnce("wailsWindow")
	handle := winc.CreateWindow("wailsWindow", parent, uint(exStyle), uint(dwStyle))
	result.SetHandle(handle)
	winc.RegMsgHandler(result)
	result.SetParent(parent)

	loadIcon := true
	if windowsOptions != nil && windowsOptions.DisableWindowIcon == true {
		loadIcon = false
	}
	if loadIcon {
		if ico, err := winc.NewIconFromResource(winc.GetAppInstance(), uint16(winc.AppIconID)); err == nil {
			result.SetIcon(0, ico)
		}
	}

	if appoptions.BackgroundColour != nil {
		win32.SetBackgroundColour(result.Handle(), appoptions.BackgroundColour.R, appoptions.BackgroundColour.G, appoptions.BackgroundColour.B)
	}

	if windowsOptions != nil {
		result.theme = windowsOptions.Theme
	} else {
		result.theme = winoptions.SystemDefault
	}

	result.SetSize(appoptions.Width, appoptions.Height)
	result.SetText(appoptions.Title)
	result.EnableSizable(!appoptions.DisableResize)
	if !appoptions.Fullscreen {
		result.EnableMaxButton(!appoptions.DisableResize)
		result.SetMinSize(appoptions.MinWidth, appoptions.MinHeight)
		result.SetMaxSize(appoptions.MaxWidth, appoptions.MaxHeight)
	}

	result.UpdateTheme()

	if windowsOptions != nil {
		result.OnSuspend = windowsOptions.OnSuspend
		result.OnResume = windowsOptions.OnResume
		if windowsOptions.WindowIsTranslucent {
			if !win32.SupportsBackdropTypes() {
				result.SetTranslucentBackground()
			} else {
				win32.EnableTranslucency(result.Handle(), win32.BackdropType(windowsOptions.BackdropType))
			}
		}

		if windowsOptions.DisableWindowIcon {
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

func (w *Window) Fullscreen() {
	if w.Form.IsFullScreen() {
		return
	}
	if w.framelessWithDecorations {
		win32.ExtendFrameIntoClientArea(w.Handle(), false)
	}
	w.Form.SetMaxSize(0, 0)
	w.Form.SetMinSize(0, 0)
	w.Form.Fullscreen()
}

func (w *Window) UnFullscreen() {
	if !w.Form.IsFullScreen() {
		return
	}
	if w.framelessWithDecorations {
		win32.ExtendFrameIntoClientArea(w.Handle(), true)
	}
	w.Form.UnFullscreen()
	w.SetMinSize(w.minWidth, w.minHeight)
	w.SetMaxSize(w.maxWidth, w.maxHeight)
}

func (w *Window) Restore() {
	if w.Form.IsFullScreen() {
		w.UnFullscreen()
	} else {
		w.Form.Restore()
	}
}

func (w *Window) SetMinSize(minWidth int, minHeight int) {
	w.minWidth = minWidth
	w.minHeight = minHeight
	w.Form.SetMinSize(minWidth, minHeight)
}

func (w *Window) SetMaxSize(maxWidth int, maxHeight int) {
	w.maxWidth = maxWidth
	w.maxHeight = maxHeight
	w.Form.SetMaxSize(maxWidth, maxHeight)
}

func (w *Window) IsVisible() bool {
	return win32.IsVisible(w.Handle())
}

func (w *Window) WndProc(msg uint32, wparam, lparam uintptr) uintptr {

	switch msg {
	case win32.WM_POWERBROADCAST:
		switch wparam {
		case win32.PBT_APMSUSPEND:
			if w.OnSuspend != nil {
				w.OnSuspend()
			}
		case win32.PBT_APMRESUMEAUTOMATIC:
			if w.OnResume != nil {
				w.OnResume()
			}
		}
	case w32.WM_SETTINGCHANGE:
		settingChanged := w32.UTF16PtrToString((*uint16)(unsafe.Pointer(lparam)))
		if settingChanged == "ImmersiveColorSet" {
			w.themeChanged = true
			w.UpdateTheme()
		}
		return 0
	case w32.WM_NCLBUTTONDOWN:
		w32.SetFocus(w.Handle())
	case w32.WM_MOVE, w32.WM_MOVING:
		w.chromium.NotifyParentWindowPositionChanged()
	case w32.WM_ACTIVATE:
		//if !w.frontendOptions.Frameless {
		w.themeChanged = true
		if int(wparam) == w32.WA_INACTIVE {
			w.isActive = false
			w.UpdateTheme()
		} else {
			w.isActive = true
			w.UpdateTheme()
			//}
		}

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
			// If we want to have a frameless window but with the default frame decorations, extend the DWM client area.
			// This Option is not affected by returning 0 in WM_NCCALCSIZE.
			// As a result we have hidden the titlebar but still have the default window frame styling.
			// See: https://docs.microsoft.com/en-us/windows/win32/api/dwmapi/nf-dwmapi-dwmextendframeintoclientarea#remarks
			if w.framelessWithDecorations {
				win32.ExtendFrameIntoClientArea(w.Handle(), true)
			}
		case w32.WM_NCCALCSIZE:
			// Disable the standard frame by allowing the client area to take the full
			// window size.
			// See: https://docs.microsoft.com/en-us/windows/win32/winmsg/wm-nccalcsize#remarks
			// This hides the titlebar and also disables the resizing from user interaction because the standard frame is not
			// shown. We still need the WS_THICKFRAME style to enable resizing from the frontend.
			if wparam != 0 {
				rgrc := (*w32.RECT)(unsafe.Pointer(lparam))
				if w.Form.IsFullScreen() {
					// In Full-Screen mode we don't need to adjust anything
					w.chromium.SetPadding(edge.Rect{})
				} else if w.IsMaximised() {
					// If the window is maximized we must adjust the client area to the work area of the monitor. Otherwise
					// some content goes beyond the visible part of the monitor.
					// Make sure to use the provided RECT to get the monitor, because during maximizig there might be
					// a wrong monitor returned in multi screen mode when using MonitorFromWindow.
					// See: https://github.com/MicrosoftEdge/WebView2Feedback/issues/2549
					monitor := w32.MonitorFromRect(rgrc, w32.MONITOR_DEFAULTTONULL)

					var monitorInfo w32.MONITORINFO
					monitorInfo.CbSize = uint32(unsafe.Sizeof(monitorInfo))
					if monitor != 0 && w32.GetMonitorInfo(monitor, &monitorInfo) {
						*rgrc = monitorInfo.RcWork

						maxWidth := w.frontendOptions.MaxWidth
						maxHeight := w.frontendOptions.MaxHeight
						if maxWidth > 0 || maxHeight > 0 {
							var dpiX, dpiY uint
							w32.GetDPIForMonitor(monitor, w32.MDT_EFFECTIVE_DPI, &dpiX, &dpiY)

							maxWidth := int32(winc.ScaleWithDPI(maxWidth, dpiX))
							if maxWidth > 0 && rgrc.Right-rgrc.Left > maxWidth {
								rgrc.Right = rgrc.Left + maxWidth
							}

							maxHeight := int32(winc.ScaleWithDPI(maxHeight, dpiY))
							if maxHeight > 0 && rgrc.Bottom-rgrc.Top > maxHeight {
								rgrc.Bottom = rgrc.Top + maxHeight
							}
						}
					}
					w.chromium.SetPadding(edge.Rect{})
				} else {
					// This is needed to workaround the resize flickering in frameless mode with WindowDecorations
					// See: https://stackoverflow.com/a/6558508
					// The workaround originally suggests to decrese the bottom 1px, but that seems to bring up a thin
					// white line on some Windows-Versions, due to DrawBackground using also this reduces ClientSize.
					// Increasing the bottom also worksaround the flickering but we would loose 1px of the WebView content
					// therefore let's pad the content with 1px at the bottom.
					rgrc.Bottom += 1
					w.chromium.SetPadding(edge.Rect{Bottom: 1})
				}
				return 0
			}
		}
	}
	return w.Form.WndProc(msg, wparam, lparam)
}

func (w *Window) IsMaximised() bool {
	return win32.IsWindowMaximised(w.Handle())
}

func (w *Window) IsMinimised() bool {
	return win32.IsWindowMinimised(w.Handle())
}

func (w *Window) IsNormal() bool {
	return win32.IsWindowNormal(w.Handle())
}

func (w *Window) IsFullScreen() bool {
	return win32.IsWindowFullScreen(w.Handle())
}

func (w *Window) SetTheme(theme winoptions.Theme) {
	w.theme = theme
	w.themeChanged = true
	w.Invoke(func() {
		w.UpdateTheme()
	})
}

func invokeSync[T any](cba *Window, fn func() (T, error)) (res T, err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	cba.Invoke(func() {
		res, err = fn()
		wg.Done()
	})
	wg.Wait()
	return res, err
}
