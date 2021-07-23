package windows

import (
	"github.com/tadvi/winc"
	"github.com/tadvi/winc/w32"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"runtime"
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

	if f.options.Windows.DisableWindowIcon {
		mainWindow.DisableIcon()
	}

	if f.options.StartHidden {
		mainWindow.Hide()
	}

	if f.options.Fullscreen {
		mainWindow.Fullscreen()
	}

	f.WindowCenter()

	if !f.options.StartHidden {
		mainWindow.Show()
	}

	mainWindow.OnClose().Bind(func(arg *winc.Event) {
		if f.options.HideWindowOnClose {
			f.WindowHide()
		} else {
			f.Quit()
		}
	})

	winc.RunMainLoop()
	return nil
}

func (f *Frontend) WindowCenter() {
	runtime.LockOSThread()
	f.mainWindow.Center()
}

func (f *Frontend) WindowSetPos(x, y int) {
	runtime.LockOSThread()
	f.mainWindow.SetPos(x, y)
}

func (f *Frontend) WindowSetSize(width, height int) {
	runtime.LockOSThread()
	f.mainWindow.SetSize(width, height)
}

func (f *Frontend) WindowSetTitle(title string) {
	runtime.LockOSThread()
	f.mainWindow.SetText(title)
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
func (f *Frontend) WindowMaximise() {
	runtime.LockOSThread()
	f.mainWindow.Maximise()
}
func (f *Frontend) WindowUnmaximise() {
	runtime.LockOSThread()
	f.mainWindow.Restore()
}
func (f *Frontend) WindowMinimise() {
	runtime.LockOSThread()
	f.mainWindow.Minimise()
}
func (f *Frontend) WindowUnminimise() {
	runtime.LockOSThread()
	f.mainWindow.Restore()
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
