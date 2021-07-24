package windows

import (
	"github.com/tadvi/winc"
	"github.com/tadvi/winc/w32"
	"github.com/wailsapp/wails/v2/pkg/options"
)

type Window struct {
	winc.Form
	frontendOptions *options.App
}

func NewWindow(parent winc.Controller, options *options.App) *Window {
	result := new(Window)
	result.frontendOptions = options
	result.SetIsForm(true)

	exStyle := w32.WS_EX_CONTROLPARENT | w32.WS_EX_APPWINDOW
	if options.Windows.WindowBackgroundIsTranslucent {
		exStyle |= w32.WS_EX_NOREDIRECTIONBITMAP
	}

	var dwStyle = w32.WS_OVERLAPPEDWINDOW
	if options.Frameless {
		dwStyle = w32.WS_POPUP
	}

	winc.RegClassOnlyOnce("wailsWindow")
	result.SetHandle(winc.CreateWindow("wailsWindow", parent, uint(exStyle), uint(dwStyle)))
	result.SetParent(parent)

	// result might fail if icon resource is not embedded in the binary
	if ico, err := winc.NewIconFromResource(winc.GetAppInstance(), uint16(winc.AppIconID)); err == nil {
		result.SetIcon(0, ico)
	}

	result.SetSize(options.Width, options.Height)
	result.SetText(options.Title)
	result.EnableSizable(!options.DisableResize)
	result.EnableMaxButton(!options.DisableResize)
	result.SetMinSize(options.MinWidth, options.MinHeight)
	result.SetMaxSize(options.MaxWidth, options.MaxHeight)

	// Dlg forces display of focus rectangles, as soon as the user starts to type.
	w32.SendMessage(result.Handle(), w32.WM_CHANGEUISTATE, w32.UIS_INITIALIZE, 0)
	winc.RegMsgHandler(result)

	result.SetFont(winc.DefaultFont)

	if options.Windows.WindowBackgroundIsTranslucent {
		result.SetTranslucentBackground()
	}

	if options.Windows.DisableWindowIcon {
		result.DisableIcon()
	}

	if options.Fullscreen {
		result.Fullscreen()
	}

	return result
}
