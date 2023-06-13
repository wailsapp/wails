package application

import (
	"errors"
	"fmt"
	"github.com/samber/lo"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/logger"

	"github.com/wailsapp/wails/v3/pkg/events"
)

type (
	webviewWindowImpl interface {
		setTitle(title string)
		setSize(width, height int)
		setAlwaysOnTop(alwaysOnTop bool)
		setURL(url string)
		setResizable(resizable bool)
		setMinSize(width, height int)
		setMaxSize(width, height int)
		execJS(js string)
		setBackgroundColour(color RGBA)
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
		zoomReset()
		zoomIn()
		zoomOut()
		getZoom() float64
		setZoom(zoom float64)
		close()
		zoom()
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
		isNormal() bool
		isVisible() bool
		setFullscreenButtonEnabled(enabled bool)
		focus()
		show()
		hide()
		getScreen() (*Screen, error)
		setFrameless(bool)
		openContextMenu(menu *Menu, data *ContextMenuData)
		nativeWindowHandle() uintptr
		startDrag() error
	}
)

type WindowEventListener struct {
	callback func(ctx *WindowEventContext)
}

type WebviewWindow struct {
	options  WebviewWindowOptions
	impl     webviewWindowImpl
	implLock sync.RWMutex
	id       uint

	eventListeners     map[uint][]*WindowEventListener
	eventListenersLock sync.RWMutex

	contextMenus     map[string]*Menu
	contextMenusLock sync.RWMutex

	// A map of listener cancellation functions
	cancellersLock sync.RWMutex
	cancellers     []func()
}

var windowID uint
var windowIDLock sync.RWMutex

func getWindowID() uint {
	windowIDLock.Lock()
	defer windowIDLock.Unlock()
	windowID++
	return windowID
}

// Use onApplicationEvent to register a callback for an application event from a window.
// This will handle tidying up the callback when the window is destroyed
func (w *WebviewWindow) onApplicationEvent(eventType events.ApplicationEventType, callback func()) {
	cancelFn := globalApplication.On(eventType, callback)
	w.addCancellationFunction(cancelFn)
}

// NewWindow creates a new window with the given options
func NewWindow(options WebviewWindowOptions) *WebviewWindow {
	if options.Width == 0 {
		options.Width = 800
	}
	if options.Height == 0 {
		options.Height = 600
	}
	if options.URL == "" {
		options.URL = "/"
	}

	result := &WebviewWindow{
		id:             getWindowID(),
		options:        options,
		eventListeners: make(map[uint][]*WindowEventListener),
		contextMenus:   make(map[string]*Menu),
	}

	return result
}

func (w *WebviewWindow) addCancellationFunction(canceller func()) {
	w.cancellersLock.Lock()
	defer w.cancellersLock.Unlock()
	w.cancellers = append(w.cancellers, canceller)
}

// SetTitle sets the title of the window
func (w *WebviewWindow) SetTitle(title string) *WebviewWindow {
	w.implLock.RLock()
	defer w.implLock.RUnlock()
	w.options.Title = title
	if w.impl != nil {
		invokeSync(func() {
			w.impl.setTitle(title)
		})
	}
	return w
}

// Name returns the name of the window
func (w *WebviewWindow) Name() string {
	return w.options.Name
}

// SetSize sets the size of the window
func (w *WebviewWindow) SetSize(width, height int) *WebviewWindow {
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
		invokeSync(func() {
			w.impl.setSize(width, height)
		})
	}
	return w
}

func (w *WebviewWindow) run() {
	if w.impl != nil {
		return
	}
	w.implLock.Lock()
	w.impl = newWindowImpl(w)
	w.implLock.Unlock()
	invokeSync(w.impl.run)
}

// SetAlwaysOnTop sets the window to be always on top.
func (w *WebviewWindow) SetAlwaysOnTop(b bool) *WebviewWindow {
	w.options.AlwaysOnTop = b
	if w.impl != nil {
		invokeSync(func() {
			w.impl.setAlwaysOnTop(b)
		})
	}
	return w
}

// Show shows the window.
func (w *WebviewWindow) Show() *WebviewWindow {
	if globalApplication.impl == nil {
		return w
	}
	if w.impl == nil {
		w.run()
		return w
	}
	invokeSync(w.impl.show)
	w.emit(events.Common.WindowShow)
	return w
}

// Hide hides the window.
func (w *WebviewWindow) Hide() *WebviewWindow {
	w.options.Hidden = true
	if w.impl != nil {
		invokeSync(w.impl.hide)
		w.emit(events.Common.WindowHide)
	}
	return w
}

func (w *WebviewWindow) SetURL(s string) *WebviewWindow {
	w.options.URL = s
	if w.impl != nil {
		invokeSync(func() {
			w.impl.setURL(s)
		})
	}
	return w
}

// SetZoom sets the zoom level of the window.
func (w *WebviewWindow) SetZoom(magnification float64) *WebviewWindow {
	w.options.Zoom = magnification
	if w.impl != nil {
		invokeSync(func() {
			w.impl.setZoom(magnification)
		})
	}
	return w
}

// GetZoom returns the current zoom level of the window.
func (w *WebviewWindow) GetZoom() float64 {
	if w.impl != nil {
		return invokeSyncWithResult(w.impl.getZoom)
	}
	return 1
}

// SetResizable sets whether the window is resizable.
func (w *WebviewWindow) SetResizable(b bool) *WebviewWindow {
	w.options.DisableResize = !b
	if w.impl != nil {
		invokeSync(func() {
			w.impl.setResizable(b)
		})
	}
	return w
}

// Resizable returns true if the window is resizable.
func (w *WebviewWindow) Resizable() bool {
	return !w.options.DisableResize
}

// SetMinSize sets the minimum size of the window.
func (w *WebviewWindow) SetMinSize(minWidth, minHeight int) *WebviewWindow {
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
			invokeSync(func() {
				w.impl.setSize(newWidth, newHeight)
			})
		}
		invokeSync(func() {
			w.impl.setMinSize(minWidth, minHeight)
		})
	}
	return w
}

// SetMaxSize sets the maximum size of the window.
func (w *WebviewWindow) SetMaxSize(maxWidth, maxHeight int) *WebviewWindow {
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
			invokeSync(func() {
				w.impl.setSize(newWidth, newHeight)
			})
		}
		invokeSync(func() {
			w.impl.setMaxSize(maxWidth, maxHeight)
		})
	}
	return w
}

// ExecJS executes the given javascript in the context of the window.
func (w *WebviewWindow) ExecJS(js string) {
	if w.impl == nil {
		return
	}
	w.impl.execJS(js)
}

// Fullscreen sets the window to fullscreen mode. Min/Max size constraints are disabled.
func (w *WebviewWindow) Fullscreen() *WebviewWindow {
	if w.impl == nil {
		w.options.StartState = WindowStateFullscreen
		return w
	}
	if !w.IsFullscreen() {
		w.disableSizeConstraints()
		invokeSync(w.impl.fullscreen)
	}
	return w
}

func (w *WebviewWindow) SetFullscreenButtonEnabled(enabled bool) *WebviewWindow {
	w.options.FullscreenButtonEnabled = enabled
	if w.impl != nil {
		invokeSync(func() {
			w.impl.setFullscreenButtonEnabled(enabled)
		})
	}
	return w
}

// IsMinimised returns true if the window is minimised
func (w *WebviewWindow) IsMinimised() bool {
	if w.impl == nil {
		return false
	}
	return invokeSyncWithResult(w.impl.isMinimised)
}

// IsVisible returns true if the window is visible
func (w *WebviewWindow) IsVisible() bool {
	if w.impl == nil {
		return false
	}
	return invokeSyncWithResult(w.impl.isVisible)
}

// IsMaximised returns true if the window is maximised
func (w *WebviewWindow) IsMaximised() bool {
	if w.impl == nil {
		return false
	}
	return invokeSyncWithResult(w.impl.isMaximised)
}

// Size returns the size of the window
func (w *WebviewWindow) Size() (int, int) {
	if w.impl == nil {
		return 0, 0
	}
	var width, height int
	invokeSync(func() {
		width, height = w.impl.size()
	})
	return width, height
}

// IsFullscreen returns true if the window is fullscreen
func (w *WebviewWindow) IsFullscreen() bool {
	w.implLock.RLock()
	defer w.implLock.RUnlock()
	if w.impl == nil {
		return false
	}
	return invokeSyncWithResult(w.impl.isFullscreen)
}

// SetBackgroundColour sets the background colour of the window
func (w *WebviewWindow) SetBackgroundColour(colour RGBA) *WebviewWindow {
	w.options.BackgroundColour = colour
	if w.impl != nil {
		invokeSync(func() {
			w.impl.setBackgroundColour(colour)
		})
	}
	return w
}

func (w *WebviewWindow) handleMessage(message string) {
	w.info(message)
	// Check for special messages
	if message == "drag" {
		if !w.IsFullscreen() {
			invokeSync(func() {
				err := w.startDrag()
				if err != nil {
					w.error("Failed to start drag: %s", err)
				}
			})
		}
	}
	w.info("ProcessMessage from front end: %s", message)

}

// Center centers the window on the screen
func (w *WebviewWindow) Center() {
	if w.impl == nil {
		w.options.Centered = true
		return
	}
	invokeSync(w.impl.center)
}

// On registers a callback for the given window event
func (w *WebviewWindow) On(eventType events.WindowEventType, callback func(ctx *WindowEventContext)) func() {
	eventID := uint(eventType)
	w.eventListenersLock.Lock()
	defer w.eventListenersLock.Unlock()
	windowEventListener := &WindowEventListener{
		callback: callback,
	}
	w.eventListeners[eventID] = append(w.eventListeners[eventID], windowEventListener)
	if w.impl != nil {
		w.impl.on(eventID)
	}

	return func() {
		w.eventListenersLock.Lock()
		defer w.eventListenersLock.Unlock()
		w.eventListeners[eventID] = lo.Without(w.eventListeners[eventID], windowEventListener)
	}

}

func (w *WebviewWindow) handleWindowEvent(id uint) {
	w.eventListenersLock.RLock()
	for _, listener := range w.eventListeners[id] {
		go listener.callback(blankWindowEventContext)
	}
	w.eventListenersLock.RUnlock()
}

// Width returns the width of the window
func (w *WebviewWindow) Width() int {
	if w.impl == nil {
		return 0
	}
	return invokeSyncWithResult(w.impl.width)
}

// Height returns the height of the window
func (w *WebviewWindow) Height() int {
	if w.impl == nil {
		return 0
	}
	return invokeSyncWithResult(w.impl.height)
}

// Position returns the position of the window
func (w *WebviewWindow) Position() (int, int) {
	w.implLock.RLock()
	defer w.implLock.RUnlock()
	if w.impl == nil {
		return 0, 0
	}
	var x, y int
	invokeSync(func() {
		x, y = w.impl.position()
	})
	return x, y
}

func (w *WebviewWindow) Destroy() {
	if w.impl == nil {
		return
	}
	// Cancel the callbacks
	for _, cancelFunc := range w.cancellers {
		cancelFunc()
	}
	invokeSync(w.impl.destroy)
}

// Reload reloads the page assets
func (w *WebviewWindow) Reload() {
	if w.impl == nil {
		return
	}
	invokeSync(w.impl.reload)
}

// ForceReload forces the window to reload the page assets
func (w *WebviewWindow) ForceReload() {
	if w.impl == nil {
		return
	}
	invokeSync(w.impl.forceReload)
}

// ToggleFullscreen toggles the window between fullscreen and normal
func (w *WebviewWindow) ToggleFullscreen() {
	if w.impl == nil {
		return
	}
	invokeSync(func() {
		if w.IsFullscreen() {
			w.UnFullscreen()
		} else {
			w.Fullscreen()
		}
	})
}

func (w *WebviewWindow) ToggleDevTools() {
	if w.impl == nil {
		return
	}
	invokeSync(w.impl.toggleDevTools)
}

// ZoomReset resets the zoom level of the webview content to 100%
func (w *WebviewWindow) ZoomReset() *WebviewWindow {
	if w.impl != nil {
		invokeSync(w.impl.zoomReset)
		w.emit(events.Common.WindowZoomReset)
	}
	return w

}

// ZoomIn increases the zoom level of the webview content
func (w *WebviewWindow) ZoomIn() {
	if w.impl == nil {
		return
	}
	invokeSync(w.impl.zoomIn)
	w.emit(events.Common.WindowZoomIn)

}

// ZoomOut decreases the zoom level of the webview content
func (w *WebviewWindow) ZoomOut() {
	if w.impl == nil {
		return
	}
	invokeSync(w.impl.zoomOut)
	w.emit(events.Common.WindowZoomOut)
}

// Close closes the window
func (w *WebviewWindow) Close() {
	if w.impl == nil {
		return
	}
	if w.options.HideOnClose {
		invokeSync(func() { w.Hide() })
		return
	}
	invokeSync(w.impl.close)
	w.emit(events.Common.WindowClose)
}

func (w *WebviewWindow) Zoom() {
	if w.impl == nil {
		return
	}
	invokeSync(w.impl.zoom)
	w.emit(events.Common.WindowZoom)
}

// SetHTML sets the HTML of the window to the given html string.
func (w *WebviewWindow) SetHTML(html string) *WebviewWindow {
	w.options.HTML = html
	if w.impl != nil {
		invokeSync(func() {
			w.impl.setHTML(html)
		})
	}
	return w
}

// SetPosition sets the position of the window.
func (w *WebviewWindow) SetPosition(x, y int) *WebviewWindow {
	w.options.X = x
	w.options.Y = y
	if w.impl != nil {
		invokeSync(func() {
			w.impl.setPosition(x, y)
		})
	}
	return w
}

// Minimise minimises the window.
func (w *WebviewWindow) Minimise() *WebviewWindow {
	if w.impl == nil {
		w.options.StartState = WindowStateMinimised
		return w
	}
	if !w.IsMinimised() {
		invokeSync(w.impl.minimise)
		w.emit(events.Common.WindowMinimise)
	}
	return w
}

// Maximise maximises the window. Min/Max size constraints are disabled.
func (w *WebviewWindow) Maximise() *WebviewWindow {
	if w.impl == nil {
		w.options.StartState = WindowStateMaximised
		return w
	}
	if !w.IsMaximised() {
		w.disableSizeConstraints()
		invokeSync(w.impl.maximise)
		w.emit(events.Common.WindowMaximise)
	}
	return w
}

// UnMinimise un-minimises the window. Min/Max size constraints are re-enabled.
func (w *WebviewWindow) UnMinimise() {
	if w.impl == nil {
		return
	}
	if w.IsMinimised() {
		invokeSync(w.impl.unminimise)
		w.emit(events.Common.WindowUnMinimise)
	}
}

// UnMaximise un-maximises the window.
func (w *WebviewWindow) UnMaximise() {
	if w.impl == nil {
		return
	}
	if w.IsMaximised() {
		w.enableSizeConstraints()
		invokeSync(w.impl.unmaximise)
		w.emit(events.Common.WindowUnMaximise)
	}
}

// UnFullscreen un-fullscreens the window.
func (w *WebviewWindow) UnFullscreen() {
	if w.impl == nil {
		return
	}
	if w.IsFullscreen() {
		w.enableSizeConstraints()
		invokeSync(w.impl.unfullscreen)
		w.emit(events.Common.WindowUnFullscreen)
	}
}

// Restore restores the window to its previous state if it was previously minimised, maximised or fullscreen.
func (w *WebviewWindow) Restore() {
	if w.impl == nil {
		return
	}
	invokeSync(func() {
		if w.IsMinimised() {
			w.UnMinimise()
		} else if w.IsMaximised() {
			w.UnMaximise()
		} else if w.IsFullscreen() {
			w.UnFullscreen()
		}
		w.emit(events.Common.WindowRestore)
	})
}

func (w *WebviewWindow) disableSizeConstraints() {
	if w.impl == nil {
		return
	}
	invokeSync(func() {
		if w.options.MinWidth > 0 && w.options.MinHeight > 0 {
			w.impl.setMinSize(0, 0)
		}
		if w.options.MaxWidth > 0 && w.options.MaxHeight > 0 {
			w.impl.setMaxSize(0, 0)
		}
	})
}

func (w *WebviewWindow) enableSizeConstraints() {
	if w.impl == nil {
		return
	}
	invokeSync(func() {
		if w.options.MinWidth > 0 && w.options.MinHeight > 0 {
			w.SetMinSize(w.options.MinWidth, w.options.MinHeight)
		}
		if w.options.MaxWidth > 0 && w.options.MaxHeight > 0 {
			w.SetMaxSize(w.options.MaxWidth, w.options.MaxHeight)
		}
	})
}

// GetScreen returns the screen that the window is on
func (w *WebviewWindow) GetScreen() (*Screen, error) {
	if w.impl == nil {
		return nil, nil
	}
	return invokeSyncWithResultAndError(w.impl.getScreen)
}

// SetFrameless removes the window frame and title bar
func (w *WebviewWindow) SetFrameless(frameless bool) *WebviewWindow {
	w.options.Frameless = frameless
	if w.impl != nil {
		invokeSync(func() {
			w.impl.setFrameless(frameless)
		})
	}
	return w
}

func (w *WebviewWindow) dispatchWailsEvent(event *WailsEvent) {
	msg := fmt.Sprintf("_wails.dispatchWailsEvent(%s);", event.ToJSON())
	w.ExecJS(msg)
}

func (w *WebviewWindow) info(message string, args ...any) {

	globalApplication.Log(&logger.Message{
		Level:   "INFO",
		Message: message,
		Data:    args,
		Sender:  w.Name(),
		Time:    time.Now(),
	})
}
func (w *WebviewWindow) error(message string, args ...any) {

	globalApplication.Log(&logger.Message{
		Level:   "ERROR",
		Message: message,
		Data:    args,
		Sender:  w.Name(),
		Time:    time.Now(),
	})
}

func (w *WebviewWindow) handleDragAndDropMessage(event *dragAndDropMessage) {
	ctx := newWindowEventContext()
	ctx.setDroppedFiles(event.filenames)
	for _, listener := range w.eventListeners[uint(events.FilesDropped)] {
		listener.callback(ctx)
	}
}

func (w *WebviewWindow) openContextMenu(data *ContextMenuData) {
	menu, ok := w.contextMenus[data.Id]
	if !ok {
		// try application level context menu
		menu, ok = globalApplication.getContextMenu(data.Id)
		if !ok {
			w.error("No context menu found for id: %s", data.Id)
			return
		}
	}
	menu.setContextData(data)
	if w.impl == nil {
		return
	}
	w.impl.openContextMenu(menu, data)
}

// RegisterContextMenu registers a context menu and assigns it the given name.
func (w *WebviewWindow) RegisterContextMenu(name string, menu *Menu) {
	w.contextMenusLock.Lock()
	defer w.contextMenusLock.Unlock()
	w.contextMenus[name] = menu
}

// NativeWindowHandle returns the platform native window handle for the window.
func (w *WebviewWindow) NativeWindowHandle() (uintptr, error) {
	if w.impl == nil {
		return 0, errors.New("native handle unavailable as window is not running")
	}
	return w.impl.nativeWindowHandle(), nil
}

func (w *WebviewWindow) Focus() {
	if w.impl == nil {
		w.options.Focused = true
		return
	}
	invokeSync(w.impl.focus)
	w.emit(events.Common.WindowFocus)
}

func (w *WebviewWindow) emit(eventType events.WindowEventType) {
	windowEvents <- &WindowEvent{
		WindowID: w.id,
		EventID:  uint(eventType),
	}
}

func (w *WebviewWindow) startDrag() error {
	if w.impl == nil {
		return nil
	}
	return invokeSyncWithError(w.impl.startDrag)
}
