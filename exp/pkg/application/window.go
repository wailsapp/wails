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
		unminimise()
		maximise()
		unmaximise()
		fullscreen()
		unfullscreen()
		isMinimised() bool
		isMaximised() bool
		isFullscreen() bool
		disableSizeConstraints()
		setFullscreenButtonEnabled(enabled bool)
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
	// Don't set size if fullscreen
	if w.IsFullscreen() {
		return w
	}
	w.options.Width = width
	w.options.Height = height

	var newMaxWidth = w.options.MaxWidth
	var newMaxHeight = w.options.MaxHeight
	if width > w.options.MaxWidth && w.options.MaxWidth != 0 {
		newMaxWidth = width
	}
	if height > w.options.MaxHeight && w.options.MaxHeight != 0 {
		newMaxHeight = height
	}

	if newMaxWidth != 0 || newMaxHeight != 0 {
		w.SetMaxSize(newMaxWidth, newMaxHeight)
	}

	var newMinWidth = w.options.MinWidth
	var newMinHeight = w.options.MinHeight
	if width < w.options.MinWidth && w.options.MinWidth != 0 {
		newMinWidth = width
	}
	if height < w.options.MinHeight && w.options.MinHeight != 0 {
		newMinHeight = height
	}

	if newMinWidth != 0 || newMinHeight != 0 {
		w.SetMinSize(newMinWidth, newMinHeight)
	}

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
	w.options.MinHeight = minHeight

	currentWidth, currentHeight := w.Size()
	newWidth, newHeight := currentWidth, currentHeight

	var newSize bool
	if minHeight != 0 && currentHeight < minHeight {
		newHeight = minHeight
		w.options.Height = newHeight
		newSize = true
	}
	if minWidth != 0 && currentWidth < minWidth {
		newWidth = minWidth
		w.options.Width = newWidth
		newSize = true
	}
	if w.impl != nil {
		if newSize {
			w.impl.setSize(newWidth, newHeight)
		}
		w.impl.setMinSize(minWidth, minHeight)
	}
	return w
}

func (w *Window) SetMaxSize(maxWidth, maxHeight int) *Window {
	w.options.MaxWidth = maxWidth
	w.options.MaxHeight = maxHeight

	currentWidth, currentHeight := w.Size()
	newWidth, newHeight := currentWidth, currentHeight

	var newSize bool
	if maxHeight != 0 && currentHeight > maxHeight {
		newHeight = maxHeight
		w.options.Height = maxHeight
		newSize = true
	}
	if maxWidth != 0 && currentWidth > maxWidth {
		newWidth = maxWidth
		w.options.Width = maxWidth
		newSize = true
	}
	if w.impl != nil {
		if newSize {
			w.impl.setSize(newWidth, newHeight)
		}
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

func (w *Window) Fullscreen() *Window {
	if w.impl == nil {
		w.options.StartState = options.WindowStateFullscreen
		return w
	}
	if !w.IsFullscreen() {
		w.disableSizeConstraints()
		w.impl.fullscreen()
	}
	return w
}

func (w *Window) SetFullscreenButtonEnabled(enabled bool) *Window {
	w.options.FullscreenButtonEnabled = enabled
	if w.impl != nil {
		w.impl.setFullscreenButtonEnabled(enabled)
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

// Size returns the size of the window
func (w *Window) Size() (width int, height int) {
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
	if w.IsFullscreen() {
		w.UnFullscreen()
	} else {
		w.Fullscreen()
	}
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
		w.disableSizeConstraints()
		w.impl.maximise()
	}
	return w
}

func (w *Window) UnMinimise() {
	if w.impl == nil {
		return
	}
	w.impl.unminimise()
}

func (w *Window) UnMaximise() {
	if w.impl == nil {
		return
	}
	w.enableSizeConstraints()
	w.impl.unmaximise()
}

func (w *Window) UnFullscreen() {
	if w.impl == nil {
		return
	}
	w.enableSizeConstraints()
	w.impl.unfullscreen()
}

func (w *Window) Restore() {
	if w.impl == nil {
		return
	}
	if w.IsMinimised() {
		w.UnMinimise()
	} else if w.IsMaximised() {
		w.UnMaximise()
	} else if w.IsFullscreen() {
		w.UnFullscreen()
	}
}

func (w *Window) disableSizeConstraints() {
	if w.impl == nil {
		return
	}
	w.impl.setMinSize(0, 0)
	w.impl.setMaxSize(0, 0)
}

func (w *Window) enableSizeConstraints() {
	if w.impl == nil {
		return
	}
	w.SetMinSize(w.options.MinWidth, w.options.MinHeight)
	w.SetMaxSize(w.options.MaxWidth, w.options.MaxHeight)
}
