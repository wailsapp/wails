package application

import (
	"github.com/wailsapp/wails/v3/pkg/events"
)

type Callback interface {
	CallError(callID string, result string)
	CallResponse(callID string, result string)
	DialogError(dialogID string, result string)
	DialogResponse(dialogID string, result string, isJSON bool)
}

type Window interface {
	Callback
	AbsolutePosition() (int, int)
	Center()
	Close()
	Destroy()
	DisableSizeConstraints()
	DispatchWailsEvent(event *WailsEvent)
	EnableSizeConstraints()
	Error(message string, args ...any)
	ExecJS(callID, js string)
	Focus()
	ForceReload()
	Fullscreen() Window
	GetScreen() (*Screen, error)
	GetZoom() float64
	HandleDragAndDropMessage(filenames []string)
	HandleMessage(message string)
	HandleWindowEvent(id uint)
	Height() int
	Hide() Window
	ID() uint
	Info(message string, args ...any)
	IsFullscreen() bool
	IsMaximised() bool
	IsMinimised() bool
	HandleKeyEvent(acceleratorString string)
	Maximise() Window
	Minimise() Window
	Name() string
	On(eventType events.WindowEventType, callback func(event *WindowEvent)) func()
	OpenContextMenu(data *ContextMenuData)
	RegisterContextMenu(name string, menu *Menu)
	RelativePosition() (int, int)
	Reload()
	Resizable() bool
	Restore()
	Run()
	SetAbsolutePosition(x, y int)
	SetAlwaysOnTop(b bool) Window
	SetBackgroundColour(colour RGBA) Window
	SetFrameless(frameless bool) Window
	SetFullscreenButtonEnabled(enabled bool) Window
	SetHTML(html string) Window
	SetMaxSize(maxWidth, maxHeight int) Window
	SetMinSize(minWidth, minHeight int) Window
	SetRelativePosition(x, y int) Window
	SetResizable(b bool) Window
	SetSize(width, height int) Window
	SetTitle(title string) Window
	SetURL(s string) Window
	SetZoom(magnification float64) Window
	Show() Window
	Size() (width int, height int)
	ToggleDevTools()
	ToggleFullscreen()
	UnFullscreen()
	UnMaximise()
	UnMinimise()
	Width() int
	Zoom()
	ZoomIn()
	ZoomOut()
	ZoomReset() Window
}
