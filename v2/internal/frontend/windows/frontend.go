package windows

import (
	"github.com/tadvi/winc"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"runtime"
)

type Frontend struct {
	frontendOptions *options.App
	logger          *logger.Logger

	// main window handle
	mainWindow                               *Window
	minWidth, minHeight, maxWidth, maxHeight int
}

func (f *Frontend) Run() error {

	mainWindow := NewWindow(nil, f.frontendOptions)
	f.mainWindow = mainWindow

	f.WindowCenter()

	if !f.frontendOptions.StartHidden {
		mainWindow.Show()
	}

	mainWindow.OnClose().Bind(func(arg *winc.Event) {
		if f.frontendOptions.HideWindowOnClose {
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
func (f *Frontend) WindowGetPos() (int, int) {
	runtime.LockOSThread()
	return f.mainWindow.Pos()
}

func (f *Frontend) WindowSetSize(width, height int) {
	runtime.LockOSThread()
	f.mainWindow.SetSize(width, height)
}

func (f *Frontend) WindowGetSize() (int, int) {
	runtime.LockOSThread()
	return f.mainWindow.Size()
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

func (f *Frontend) WindowSetMinSize(width int, height int) {
	runtime.LockOSThread()
	f.mainWindow.SetMinSize(width, height)
}
func (f *Frontend) WindowSetMaxSize(width int, height int) {
	runtime.LockOSThread()
	f.mainWindow.SetMaxSize(width, height)
}

func (f *Frontend) Quit() {
	winc.Exit()
}

func NewFrontend(appoptions *options.App, myLogger *logger.Logger) *Frontend {

	return &Frontend{
		frontendOptions: appoptions,
		logger:          myLogger,
	}
}
