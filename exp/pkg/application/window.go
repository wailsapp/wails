package application

import (
	"fmt"
	"sync"

	"github.com/wailsapp/wails/exp/pkg/events"
	"github.com/wailsapp/wails/exp/pkg/options"
)

type (
	windowImpl interface {
		setTitle(title string)
		setSize(width, height int)
		setAlwaysOnTop(alwaysOnTop bool)
		setURL(url string)
		setResizable(resizable bool)
		setMinSize(width, height int)
		setMaxSize(width, height int)
		execJS(js string)
		setFullscreen()
		isMinimised() bool
		isMaximised() bool
		isFullscreen() bool
		restore()
		setBackgroundColour(color *options.RGBA)
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
		setHTML(html string)
		setPosition(x int, y int)
		on(eventID uint)
		minimise()
		maximise()
		unminimise()
		unmaximise()
	}
)

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
	if options.Width == 0 {
		options.Width = 800
	}
	if options.Height == 0 {
		options.Height = 600
	}
	return &Window{
		id:             getWindowID(),
		options:        options,
		eventListeners: make(map[uint][]func()),
	}
}

func (w *Window) SetTitle(title string) *Window {
	w.implLock.RLock()
	defer w.implLock.RUnlock()
	w.options.Title = title
	if w.impl != nil {
		w.impl.setTitle(title)
	}
	return w
}

func (w *Window) SetSize(width, height int) *Window {
	w.options.Width = width
	w.options.Height = height
	if w.impl != nil {
		w.impl.setSize(width, height)
	}
	return w
}

func (w *Window) Run() {
	if w.impl != nil {
		return
	}
	w.implLock.Lock()
	w.impl = newWindowImpl(w)
	w.implLock.Unlock()
	w.impl.run()
}

func (w *Window) SetAlwaysOnTop(b bool) *Window {
	w.options.AlwaysOnTop = b
	if w.impl == nil {
		w.impl.setAlwaysOnTop(b)
	}
	return w
}

func (w *Window) SetURL(s string) *Window {
	w.options.URL = s
	if w.impl != nil {
		w.impl.setURL(s)
	}
	return w
}

func (w *Window) SetResizable(b bool) *Window {
	w.options.DisableResize = !b
	if w.impl != nil {
		w.impl.setResizable(b)
	}
	return w
}

func (w *Window) Resizable() bool {
	return !w.options.DisableResize
}

func (w *Window) SetMinSize(minWidth, minHeight int) *Window {
	w.options.MinWidth = minWidth
	if w.options.Width < minWidth {
		w.options.Width = minWidth
	}
	w.options.MinHeight = minHeight
	if w.options.Height < minHeight {
		w.options.Height = minHeight
	}
	if w.impl != nil {
		w.impl.setSize(w.options.Width, w.options.Height)
		w.impl.setMinSize(minWidth, minHeight)
	}
	return w
}

func (w *Window) SetMaxSize(maxWidth, maxHeight int) *Window {
	w.options.MinWidth = maxWidth
	if w.options.Width > maxWidth {
		w.options.Width = maxWidth
	}
	w.options.MinHeight = maxHeight
	if w.options.Height > maxHeight {
		w.options.Height = maxHeight
	}
	if w.impl != nil {
		w.impl.setSize(w.options.Width, w.options.Height)
		w.impl.setMaxSize(maxWidth, maxHeight)
	}
	return w
}

func (w *Window) ExecJS(js string) {
	if w.impl == nil {
		return
	}
	w.impl.execJS(js)
}

func (w *Window) SetFullscreen() *Window {
	if w.impl == nil {
		w.options.StartState = options.WindowStateFullscreen
		return w
	}
	if !w.IsFullscreen() {
		w.impl.setFullscreen()
	}
	return w
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

func (w *Window) SetBackgroundColour(colour *options.RGBA) *Window {
	w.options.BackgroundColour = colour
	if w.impl != nil {
		w.impl.setBackgroundColour(colour)
	}
	return w
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
	defer w.eventListenersLock.Unlock()
	w.eventListeners[eventID] = append(w.eventListeners[eventID], callback)
	if w.impl != nil {
		w.impl.on(eventID)
	}
}

func (w *Window) handleWindowEvent(id uint) {
	w.eventListenersLock.RLock()
	for _, callback := range w.eventListeners[id] {
		go callback()
	}
	w.eventListenersLock.RUnlock()
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

func (w *Window) ResetZoom() *Window {
	if w.impl != nil {
		w.impl.resetZoom()
	}
	return w

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

func (w *Window) SetHTML(html string) *Window {
	w.options.HTML = html
	if w.impl != nil {
		w.impl.setHTML(html)
	}
	return w
}

func (w *Window) SetPosition(x, y int) *Window {
	w.options.X = x
	w.options.Y = y
	if w.impl != nil {
		w.impl.setPosition(x, y)
	}
	return w
}

func (w *Window) Minimise() *Window {
	if w.impl == nil {
		w.options.StartState = options.WindowStateMinimised
		return w
	}
	if !w.IsMinimised() {
		w.impl.minimise()
	}
	return w
}

func (w *Window) Maximise() *Window {
	if w.impl == nil {
		w.options.StartState = options.WindowStateMaximised
		return w
	}
	if !w.IsMaximised() {
		w.impl.maximise()
	}
	return w
}

func (w *Window) Restore() {
	if w.impl == nil {
		return
	}
	if w.IsMinimised() {
		w.impl.unminimise()
	} else if w.IsMaximised() {
		w.impl.unmaximise()
	} else if w.IsFullscreen() {
		w.impl.toggleFullscreen()
	}
}
