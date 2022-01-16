//go:build windows

package windows

import (
	"github.com/leaanthony/winc"
	"github.com/leaanthony/winc/w32"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
)

type Window struct {
	winc.Form
	frontendOptions *options.App
	applicationMenu *menu.Menu
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
	if appoptions.Frameless {
		dwStyle = w32.WS_POPUP
		if winoptions := appoptions.Windows; winoptions != nil && winoptions.EnableFramelessBorder {
			dwStyle |= w32.WS_BORDER
		}
	}

	winc.RegClassOnlyOnce("wailsWindow")
	result.SetHandle(winc.CreateWindow("wailsWindow", parent, uint(exStyle), uint(dwStyle)))
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
	if appoptions.Frameless == false && !appoptions.Fullscreen {
		result.EnableMaxButton(!appoptions.DisableResize)
		result.SetMinSize(appoptions.MinWidth, appoptions.MinHeight)
		result.SetMaxSize(appoptions.MaxWidth, appoptions.MaxHeight)
	}
	result.EnableSizable(!appoptions.DisableResize)

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
	winc.RegMsgHandler(result)

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
		w.frontendOptions.Windows.NotifyParentWindowPositionChanged()
	}
	return w.Form.WndProc(msg, wparam, lparam)
}
