package application

import (
	"fmt"
	"sync"

	"github.com/wailsapp/wails/exp/pkg/events"
	"github.com/wailsapp/wails/exp/pkg/options"
)

type windowImpl interface {
	setTitle(title string)
	setSize(width, height int)
	setAlwaysOnTop(alwaysOnTop bool)
	navigateToURL(url string)
	setResizable(resizable bool)
	setMinSize(width, height int)
	setMaxSize(width, height int)
	enableDevTools()
	execJS(js string)
	setMaximised()
	setMinimised()
	setFullscreen()
	isMinimised() bool
	isMaximised() bool
	isFullscreen() bool
	restore()
	setBackgroundColor(color *options.RGBA)
	run()
	center()
	size() (int, int)
	width() int
	height() int
	position() (int, int)
	destroy()
	reload()
	forceReload()
	toggleFullscreen()
	toggleDevTools()
	resetZoom()
	zoomIn()
	zoomOut()
	close()
	zoom()
	minimize()
}

type Window struct {
	options  *options.Window
	impl     windowImpl
	implLock sync.RWMutex
	id       uint

	eventListeners     map[uint][]func()
	eventListenersLock sync.RWMutex
}

var windowID uint
var windowIDLock sync.RWMutex

func getWindowID() uint {
	windowIDLock.Lock()
	defer windowIDLock.Unlock()
	windowID++
	return windowID
}

func NewWindow(options *options.Window) *Window {
	return &Window{
		id:             getWindowID(),
		options:        options,
		eventListeners: make(map[uint][]func()),
	}
}

func (w *Window) SetTitle(title string) {
	w.implLock.RLock()
	defer w.implLock.RUnlock()
	if w.impl == nil {
		w.options.Title = title
		return
	}
	w.impl.setTitle(title)
}

func (w *Window) SetSize(width, height int) {
	if w.impl == nil {
		w.options.Width = width
		w.options.Height = height
		return
	}
	w.impl.setSize(width, height)
}

func (w *Window) Run() {
	w.implLock.Lock()
	w.impl = newWindowImpl(w.id, w.options)
	w.implLock.Unlock()
	w.impl.run()
}

func (w *Window) SetAlwaysOnTop(b bool) {
	if w.impl == nil {
		w.options.AlwaysOnTop = b
		return
	}
	w.impl.setAlwaysOnTop(b)
}

func (w *Window) NavigateToURL(s string) {
	if w.impl == nil {
		w.options.URL = s
		return
	}
	w.impl.navigateToURL(s)
}

func (w *Window) SetResizable(b bool) {
	w.options.DisableResize = !b
	if w.impl == nil {
		return
	}
	w.impl.setResizable(b)
}

func (w *Window) Resizable() bool {
	return !w.options.DisableResize
}

func (w *Window) SetMinSize(minWidth, minHeight int) {
	if w.impl == nil {
		w.options.MinWidth = minWidth
		if w.options.Width < minWidth {
			w.options.Width = minWidth
		}
		w.options.MinHeight = minHeight
		if w.options.Height < minHeight {
			w.options.Height = minHeight
		}
		return
	}
	w.impl.setSize(w.options.Width, w.options.Height)
	w.impl.setMinSize(minWidth, minHeight)
}
func (w *Window) SetMaxSize(maxWidth, maxHeight int) {
	if w.impl == nil {
		w.options.MinWidth = maxWidth
		if w.options.Width > maxWidth {
			w.options.Width = maxWidth
		}
		w.options.MinHeight = maxHeight
		if w.options.Height > maxHeight {
			w.options.Height = maxHeight
		}
		return
	}
	w.impl.setSize(w.options.Width, w.options.Height)
	w.impl.setMaxSize(maxWidth, maxHeight)
}

func (w *Window) EnableDevTools() {
	if w.impl == nil {
		w.options.EnableDevTools = true
		return
	}
	w.impl.enableDevTools()
}

func (w *Window) ExecJS(js string) {
	if w.impl == nil {
		return
	}
	w.impl.execJS(js)
}

// Set Maximized
func (w *Window) SetMaximized() {
	w.options.StartState = options.WindowStateMaximised
	if w.impl == nil {
		return
	}
	w.impl.setMaximised()
}

// Set Minimized
func (w *Window) SetMinimized() {
	w.options.StartState = options.WindowStateMinimised
	if w.impl == nil {
		return
	}
	w.impl.setMinimised()
}

// Set Fullscreen
func (w *Window) SetFullscreen() {
	w.options.StartState = options.WindowStateFullscreen
	if w.impl == nil {
		return
	}
	w.impl.setFullscreen()
}

// IsMinimised returns true if the window is minimised
func (w *Window) IsMinimised() bool {
	if w.impl == nil {
		return false
	}
	return w.impl.isMinimised()
}

// IsMaximised returns true if the window is maximised
func (w *Window) IsMaximised() bool {
	if w.impl == nil {
		return false
	}
	return w.impl.isMaximised()
}

// Size returns the current size of the window
func (w *Window) Size() (int, int) {
	if w.impl == nil {
		return 0, 0
	}
	return w.impl.size()
}

// IsFullscreen returns true if the window is fullscreen
func (w *Window) IsFullscreen() bool {
	w.implLock.RLock()
	defer w.implLock.RUnlock()
	if w.impl == nil {
		return false
	}
	return w.impl.isFullscreen()
}

func (w *Window) SetBackgroundColor(color *options.RGBA) {
	w.options.BackgroundColour = color
	if w.impl == nil {
		return
	}
	w.impl.setBackgroundColor(color)
}

func (w *Window) handleMessage(message string) {
	fmt.Printf("[window %d] %s", w.id, message)
	// Check for special messages
	if message == "test" {
		w.SetTitle("Hello World")
	}
}

func (w *Window) Center() {
	if w.impl == nil {
		return
	}
	w.impl.center()
}

func (w *Window) On(eventType events.WindowEventType, callback func()) {
	eventID := uint(eventType)
	w.eventListenersLock.Lock()
	w.eventListeners[eventID] = append(w.eventListeners[eventID], callback)
	w.eventListenersLock.Unlock()
}

func (w *Window) handleWindowEvent(id uint) {
	w.eventListenersLock.RLock()
	for _, callback := range w.eventListeners[id] {
		go callback()
	}
	w.eventListenersLock.RUnlock()
}

func (w *Window) ID() uint {
	return w.id
}

func (w *Window) Width() int {
	if w.impl == nil {
		return 0
	}
	return w.impl.width()
}

func (w *Window) Height() int {
	if w.impl == nil {
		return 0
	}
	return w.impl.height()
}

func (w *Window) Position() (int, int) {
	w.implLock.RLock()
	defer w.implLock.RUnlock()
	if w.impl == nil {
		return 0, 0
	}
	return w.impl.position()
}

func (w *Window) Destroy() {
	if w.impl == nil {
		return
	}
	w.impl.destroy()
}

func (w *Window) Reload() {
	if w.impl == nil {
		return
	}
	w.impl.reload()
}

func (w *Window) ForceReload() {
	if w.impl == nil {
		return
	}
	w.impl.forceReload()
}

func (w *Window) ToggleFullscreen() {
	if w.impl == nil {
		return
	}
	w.impl.toggleFullscreen()
}

func (w *Window) ToggleDevTools() {
	if w.impl == nil {
		return
	}
	w.impl.toggleDevTools()
}

func (w *Window) ResetZoom() {
	if w.impl == nil {
		return
	}
	w.impl.resetZoom()
}

func (w *Window) ZoomIn() {
	if w.impl == nil {
		return
	}
	w.impl.zoomIn()
}

func (w *Window) ZoomOut() {
	if w.impl == nil {
		return
	}
	w.impl.zoomOut()
}

func (w *Window) Close() {
	if w.impl == nil {
		return
	}
	w.impl.close()
}

func (w *Window) Minimize() {
	if w.impl == nil {
		return
	}
	w.impl.minimize()
}

func (w *Window) Zoom() {
	if w.impl == nil {
		return
	}
	w.impl.zoom()
}
