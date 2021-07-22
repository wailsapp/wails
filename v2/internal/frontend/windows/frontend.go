package windows

import (
	"github.com/tadvi/winc"
	"github.com/tadvi/winc/w32"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"runtime"
	"time"
)

type Frontend struct {
	options *options.App
	logger  *logger.Logger

	// main window handle
	mainWindow *winc.Form
}

func (f *Frontend) Run() error {

	exStyle := w32.WS_EX_CONTROLPARENT | w32.WS_EX_APPWINDOW
	if f.options.Windows.WindowBackgroundIsTranslucent {
		exStyle |= w32.WS_EX_NOREDIRECTIONBITMAP
	}

	var dwStyle uint
	if f.options.Frameless {
		dwStyle = w32.WS_POPUP
	}

	mainWindow := winc.NewCustomForm(nil, exStyle, dwStyle)
	f.mainWindow = mainWindow
	mainWindow.SetSize(f.options.Width, f.options.Height)
	mainWindow.SetText(f.options.Title)
	mainWindow.EnableSizable(!f.options.DisableResize)
	mainWindow.EnableMaxButton(!f.options.DisableResize)

	if f.options.Windows.WindowBackgroundIsTranslucent {
		mainWindow.SetTranslucentBackground()
	}

	if f.options.Fullscreen {
		mainWindow.Fullscreen()
	}

	mainWindow.Center()
	mainWindow.Show()
	mainWindow.OnClose().Bind(func(arg *winc.Event) {
		if f.options.HideWindowOnClose {

			go func() {
				time.Sleep(1 * time.Second)
				f.WindowShow()
			}()
		} else {
			f.Quit()
		}
	})

	winc.RunMainLoop()
	return nil
}

func (f *Frontend) WindowFullscreen() {
	runtime.LockOSThread()
	f.mainWindow.Fullscreen()
}

func (f *Frontend) WindowUnFullscreen() {
	runtime.LockOSThread()
	f.mainWindow.UnFullscreen()
}

func (f *Frontend) WindowShow() {
	runtime.LockOSThread()
	f.mainWindow.Show()
}

func (f *Frontend) WindowHide() {
	runtime.LockOSThread()
	f.mainWindow.Hide()
}

func (f *Frontend) Quit() {
	winc.Exit()
}

func NewFrontend(appoptions *options.App, myLogger *logger.Logger) *Frontend {

	return &Frontend{
		options: appoptions,
		logger:  myLogger,
	}
}
