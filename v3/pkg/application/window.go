package application

import (
	"github.com/wailsapp/wails/v3/pkg/events"
)

type Callback interface {
	CallError(callID string, result string, isJSON bool)
	CallResponse(callID string, result string)
	DialogError(dialogID string, result string)
	DialogResponse(dialogID string, result string, isJSON bool)
}

type Window interface {
	Callback
	Center()
	Close()
	DisableSizeConstraints()
	DispatchWailsEvent(event *CustomEvent)
	EmitEvent(name string, data ...any)
	EnableSizeConstraints()
	Error(message string, args ...any)
	ExecJS(js string)
	Focus()
	ForceReload()
	Fullscreen() Window
	GetBorderSizes() *LRTB
	GetScreen() (*Screen, error)
	GetZoom() float64
	HandleDragAndDropMessage(filenames []string)
	HandleMessage(message string)
	HandleWindowEvent(id uint)
	Height() int
	Hide() Window
	HideMenuBar()
	ID() uint
	Info(message string, args ...any)
	IsFocused() bool
	IsFullscreen() bool
	IsIgnoreMouseEvents() bool
	IsMaximised() bool
	IsMinimised() bool
	HandleKeyEvent(acceleratorString string)
	Maximise() Window
	Minimise() Window
	Name() string
	OnWindowEvent(eventType events.WindowEventType, callback func(event *WindowEvent)) func()
	OpenContextMenu(data *ContextMenuData)
	Position() (int, int)
	RelativePosition() (int, int)
	Reload()
	Resizable() bool
	Restore()
	Run()
	SetPosition(x, y int)
	SetAlwaysOnTop(b bool) Window
	SetBackgroundColour(colour RGBA) Window
	SetFrameless(frameless bool) Window
	SetHTML(html string) Window
	SetMinimiseButtonState(state ButtonState) Window
	SetMaximiseButtonState(state ButtonState) Window
	SetCloseButtonState(state ButtonState) Window
	SetMaxSize(maxWidth, maxHeight int) Window
	SetMinSize(minWidth, minHeight int) Window
	SetRelativePosition(x, y int) Window
	SetResizable(b bool) Window
	SetIgnoreMouseEvents(ignore bool) Window
	SetSize(width, height int) Window
	SetTitle(title string) Window
	SetURL(s string) Window
	SetZoom(magnification float64) Window
	Show() Window
	ShowMenuBar()
	Size() (width int, height int)
	OpenDevTools()
	ToggleFullscreen()
	ToggleMaximise()
	ToggleMenuBar()
	UnFullscreen()
	UnMaximise()
	UnMinimise()
	Width() int
	Zoom()
	ZoomIn()
	ZoomOut()
	ZoomReset() Window
}
