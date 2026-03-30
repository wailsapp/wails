//go:build server

package application

import (
	"fmt"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// BrowserWindow represents a browser client connection in server mode.
// It implements the Window interface so browser clients can be treated
// uniformly with native windows throughout the codebase.
type BrowserWindow struct {
	id       uint
	name     string
	clientID string // The runtime's nanoid for this client
}

// NewBrowserWindow creates a new browser window with the given ID.
func NewBrowserWindow(id uint, clientID string) *BrowserWindow {
	return &BrowserWindow{
		id:       id,
		name:     fmt.Sprintf("browser-%d", id),
		clientID: clientID,
	}
}

// Core identification methods

func (b *BrowserWindow) ID() uint      { return b.id }
func (b *BrowserWindow) Name() string  { return b.name }
func (b *BrowserWindow) ClientID() string { return b.clientID }

// Event methods - these are meaningful for browser windows

func (b *BrowserWindow) DispatchWailsEvent(event *CustomEvent) {
	// Events are dispatched via WebSocket broadcast, not per-window
}

func (b *BrowserWindow) EmitEvent(name string, data ...any) bool {
	return globalApplication.Event.Emit(name, data...)
}

// Logging methods

func (b *BrowserWindow) Error(message string, args ...any) {
	globalApplication.error(message, args...)
}

func (b *BrowserWindow) Info(message string, args ...any) {
	globalApplication.info(message, args...)
}

// No-op methods - these don't apply to browser windows


func (b *BrowserWindow) Center()                                      {}
func (b *BrowserWindow) Close()                                       {}
func (b *BrowserWindow) DisableSizeConstraints()                      {}
func (b *BrowserWindow) EnableSizeConstraints()                       {}
func (b *BrowserWindow) ExecJS(js string)                             {}
func (b *BrowserWindow) Focus()                                       {}
func (b *BrowserWindow) ForceReload()                                 {}
func (b *BrowserWindow) Fullscreen() Window                           { return b }
func (b *BrowserWindow) GetBorderSizes() *LRTB                        { return nil }
func (b *BrowserWindow) GetScreen() (*Screen, error)                  { return nil, nil }
func (b *BrowserWindow) GetZoom() float64                             { return 1.0 }
func (b *BrowserWindow) handleDragAndDropMessage(filenames []string, dropTarget *DropTargetDetails) {}
func (b *BrowserWindow) InitiateFrontendDropProcessing(filenames []string, x int, y int) {}
func (b *BrowserWindow) HandleMessage(message string)                 {}
func (b *BrowserWindow) HandleWindowEvent(id uint)                    {}
func (b *BrowserWindow) Height() int                                  { return 0 }
func (b *BrowserWindow) Hide() Window                                 { return b }
func (b *BrowserWindow) HideMenuBar()                                 {}
func (b *BrowserWindow) IsFocused() bool                              { return false }
func (b *BrowserWindow) IsFullscreen() bool                           { return false }
func (b *BrowserWindow) IsIgnoreMouseEvents() bool                    { return false }
func (b *BrowserWindow) IsMaximised() bool                            { return false }
func (b *BrowserWindow) IsMinimised() bool                            { return false }
func (b *BrowserWindow) HandleKeyEvent(acceleratorString string)      {}
func (b *BrowserWindow) Maximise() Window                             { return b }
func (b *BrowserWindow) Minimise() Window                             { return b }
func (b *BrowserWindow) OnWindowEvent(eventType events.WindowEventType, callback func(event *WindowEvent)) func() {
	return func() {}
}
func (b *BrowserWindow) OpenContextMenu(data *ContextMenuData)        {}
func (b *BrowserWindow) Position() (int, int)                         { return 0, 0 }
func (b *BrowserWindow) RelativePosition() (int, int)                 { return 0, 0 }
func (b *BrowserWindow) Reload()                                      {}
func (b *BrowserWindow) Resizable() bool                              { return false }
func (b *BrowserWindow) Restore()                                     {}
func (b *BrowserWindow) Run()                                         {}
func (b *BrowserWindow) SetPosition(x, y int)                         {}
func (b *BrowserWindow) SetAlwaysOnTop(b2 bool) Window                { return b }
func (b *BrowserWindow) SetBackgroundColour(colour RGBA) Window       { return b }
func (b *BrowserWindow) SetFrameless(frameless bool) Window           { return b }
func (b *BrowserWindow) SetHTML(html string) Window                   { return b }
func (b *BrowserWindow) SetMinimiseButtonState(state ButtonState) Window { return b }
func (b *BrowserWindow) SetMaximiseButtonState(state ButtonState) Window { return b }
func (b *BrowserWindow) SetCloseButtonState(state ButtonState) Window { return b }
func (b *BrowserWindow) SetMaxSize(maxWidth, maxHeight int) Window    { return b }
func (b *BrowserWindow) SetMinSize(minWidth, minHeight int) Window    { return b }
func (b *BrowserWindow) SetRelativePosition(x, y int) Window          { return b }
func (b *BrowserWindow) SetResizable(b2 bool) Window                  { return b }
func (b *BrowserWindow) SetIgnoreMouseEvents(ignore bool) Window      { return b }
func (b *BrowserWindow) SetSize(width, height int) Window             { return b }
func (b *BrowserWindow) SetTitle(title string) Window                 { return b }
func (b *BrowserWindow) SetURL(s string) Window                       { return b }
func (b *BrowserWindow) SetZoom(magnification float64) Window         { return b }
func (b *BrowserWindow) Show() Window                                 { return b }
func (b *BrowserWindow) ShowMenuBar()                                 {}
func (b *BrowserWindow) Size() (width int, height int)                { return 0, 0 }
func (b *BrowserWindow) OpenDevTools()                                {}
func (b *BrowserWindow) ToggleFullscreen()                            {}
func (b *BrowserWindow) ToggleMaximise()                              {}
func (b *BrowserWindow) ToggleMenuBar()                               {}
func (b *BrowserWindow) ToggleFrameless()                             {}
func (b *BrowserWindow) UnFullscreen()                                {}
func (b *BrowserWindow) UnMaximise()                                  {}
func (b *BrowserWindow) UnMinimise()                                  {}
func (b *BrowserWindow) Width() int                                   { return 0 }
func (b *BrowserWindow) IsVisible() bool                              { return true }
func (b *BrowserWindow) Bounds() Rect                                 { return Rect{} }
func (b *BrowserWindow) SetBounds(bounds Rect)                        {}
func (b *BrowserWindow) Zoom()                                        {}
func (b *BrowserWindow) ZoomIn()                                      {}
func (b *BrowserWindow) ZoomOut()                                     {}
func (b *BrowserWindow) ZoomReset() Window                            { return b }
func (b *BrowserWindow) SetMenu(menu *Menu)                           {}
func (b *BrowserWindow) SnapAssist()                                  {}
func (b *BrowserWindow) SetContentProtection(protection bool) Window  { return b }
func (b *BrowserWindow) NativeWindow() unsafe.Pointer                 { return nil }
func (b *BrowserWindow) SetEnabled(enabled bool)                      {}
func (b *BrowserWindow) Flash(enabled bool)                           {}
func (b *BrowserWindow) Print() error                                 { return nil }
func (b *BrowserWindow) RegisterHook(eventType events.WindowEventType, callback func(event *WindowEvent)) func() {
	return func() {}
}
func (b *BrowserWindow) shouldUnconditionallyClose() bool             { return true }

// Editing methods
func (b *BrowserWindow) cut()       {}
func (b *BrowserWindow) copy()      {}
func (b *BrowserWindow) paste()     {}
func (b *BrowserWindow) undo()      {}
func (b *BrowserWindow) redo()      {}
func (b *BrowserWindow) delete()    {}
func (b *BrowserWindow) selectAll() {}
