//go:build windows

package windows

import (
	"github.com/leaanthony/winc"
	"github.com/leaanthony/winc/w32"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"sync"
)

type Window struct {
	winc.Form
	frontendOptions *options.App
	applicationMenu *menu.Menu
	m               sync.Mutex
	dispatchq       []func()
}

func NewWindow(parent winc.Controller, options *options.App) *Window {
	result := new(Window)
	result.frontendOptions = options
	result.SetIsForm(true)

	var exStyle int
	if options.Windows != nil {
		exStyle = w32.WS_EX_CONTROLPARENT | w32.WS_EX_APPWINDOW
		if options.Windows.WindowIsTranslucent {
			exStyle |= w32.WS_EX_NOREDIRECTIONBITMAP
		}
	}
	if options.AlwaysOnTop {
		exStyle |= w32.WS_EX_TOPMOST
	}

	var dwStyle = w32.WS_OVERLAPPEDWINDOW
	if options.Frameless {
		dwStyle = w32.WS_POPUP
	}

	winc.RegClassOnlyOnce("wailsWindow")
	result.SetHandle(winc.CreateWindow("wailsWindow", parent, uint(exStyle), uint(dwStyle)))
	result.SetParent(parent)

	loadIcon := true
	if options.Windows != nil && options.Windows.DisableWindowIcon == true {
		loadIcon = false
	}
	if loadIcon {
		if ico, err := winc.NewIconFromResource(winc.GetAppInstance(), uint16(winc.AppIconID)); err == nil {
			result.SetIcon(0, ico)
		}
	}

	result.SetSize(options.Width, options.Height)
	result.SetText(options.Title)
	if options.Frameless == false && !options.Fullscreen {
		result.EnableMaxButton(!options.DisableResize)
		result.EnableSizable(!options.DisableResize)
		result.SetMinSize(options.MinWidth, options.MinHeight)
		result.SetMaxSize(options.MaxWidth, options.MaxHeight)
	}

	if options.Windows != nil {
		if options.Windows.WindowIsTranslucent {
			result.SetTranslucentBackground()
		}

		if options.Windows.DisableWindowIcon {
			result.DisableIcon()
		}
	}

	// Dlg forces display of focus rectangles, as soon as the user starts to type.
	w32.SendMessage(result.Handle(), w32.WM_CHANGEUISTATE, w32.UIS_INITIALIZE, 0)
	winc.RegMsgHandler(result)

	result.SetFont(winc.DefaultFont)

	if options.Menu != nil {
		result.SetApplicationMenu(options.Menu)
	}

	return result
}

func (w *Window) Run() int {
	var m w32.MSG

	for w32.GetMessage(&m, 0, 0, 0) != 0 {
		if m.Message == w32.WM_APP {
			// Credit: https://github.com/jchv/go-webview2
			w.m.Lock()
			q := append([]func(){}, w.dispatchq...)
			w.dispatchq = []func(){}
			w.m.Unlock()
			for _, v := range q {
				v()
			}
		}
		if !w.PreTranslateMessage(&m) {
			w32.TranslateMessage(&m)
			w32.DispatchMessage(&m)
		}
	}

	w32.GdiplusShutdown()
	return int(m.WParam)
}

func (w *Window) Dispatch(f func()) {
	w.m.Lock()
	w.dispatchq = append(w.dispatchq, f)
	w.m.Unlock()
	w32.PostMainThreadMessage(w32.WM_APP, 0, 0)
}
