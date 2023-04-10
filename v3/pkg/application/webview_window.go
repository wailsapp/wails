package application

import (
	"fmt"
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
		restore()
		setBackgroundColour(color *RGBA)
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
		disableSizeConstraints()
		setFullscreenButtonEnabled(enabled bool)
		show()
		hide()
		getScreen() (*Screen, error)
		setFrameless(bool)
		openContextMenu(menu *Menu, data *ContextMenuData)
	}
)

type WebviewWindow struct {
	options  *WebviewWindowOptions
	impl     webviewWindowImpl
	implLock sync.RWMutex
	id       uint

	eventListeners     map[uint][]func(ctx *WindowEventContext)
	eventListenersLock sync.RWMutex

	contextMenus     map[string]*Menu
	contextMenusLock sync.RWMutex
}

var windowID uint
var windowIDLock sync.RWMutex

func getWindowID() uint {
	windowIDLock.Lock()
	defer windowIDLock.Unlock()
	windowID++
	return windowID
}

func NewWindow(options *WebviewWindowOptions) *WebviewWindow {
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
		eventListeners: make(map[uint][]func(ctx *WindowEventContext)),
		contextMenus:   make(map[string]*Menu),
	}

	return result
}

func (w *WebviewWindow) SetTitle(title string) *WebviewWindow {
	w.implLock.RLock()
	defer w.implLock.RUnlock()
	w.options.Title = title
	if w.impl != nil {
		w.impl.setTitle(title)
	}
	return w
}

func (w *WebviewWindow) Name() string {
	return w.options.Name
}

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
		w.impl.setSize(width, height)
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
	w.impl.run()
}

func (w *WebviewWindow) SetAlwaysOnTop(b bool) *WebviewWindow {
	w.options.AlwaysOnTop = b
	if w.impl == nil {
		w.impl.setAlwaysOnTop(b)
	}
	return w
}

func (w *WebviewWindow) Show() *WebviewWindow {
	if globalApplication.impl == nil {
		return w
	}
	if w.impl == nil {
		w.run()
		return w
	}
	w.impl.show()
	return w
}
func (w *WebviewWindow) Hide() *WebviewWindow {
	w.options.Hidden = true
	if w.impl != nil {
		w.impl.hide()
	}
	return w
}

func (w *WebviewWindow) SetURL(s string) *WebviewWindow {
	w.options.URL = s
	if w.impl != nil {
		w.impl.setURL(s)
	}
	return w
}

func (w *WebviewWindow) SetZoom(magnification float64) *WebviewWindow {
	w.options.Zoom = magnification
	if w.impl != nil {
		w.impl.setZoom(magnification)
	}
	return w
}

func (w *WebviewWindow) GetZoom() float64 {
	if w.impl != nil {
		return w.impl.getZoom()
	}
	return 1
}

func (w *WebviewWindow) SetResizable(b bool) *WebviewWindow {
	w.options.DisableResize = !b
	if w.impl != nil {
		w.impl.setResizable(b)
	}
	return w
}

func (w *WebviewWindow) Resizable() bool {
	return !w.options.DisableResize
}

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
			w.impl.setSize(newWidth, newHeight)
		}
		w.impl.setMinSize(minWidth, minHeight)
	}
	return w
}

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
			w.impl.setSize(newWidth, newHeight)
		}
		w.impl.setMaxSize(maxWidth, maxHeight)
	}
	return w
}

func (w *WebviewWindow) ExecJS(js string) {
	if w.impl == nil {
		return
	}
	w.impl.execJS(js)
}

func (w *WebviewWindow) Fullscreen() *WebviewWindow {
	if w.impl == nil {
		w.options.StartState = WindowStateFullscreen
		return w
	}
	if !w.IsFullscreen() {
		w.disableSizeConstraints()
		w.impl.fullscreen()
	}
	return w
}

func (w *WebviewWindow) SetFullscreenButtonEnabled(enabled bool) *WebviewWindow {
	w.options.FullscreenButtonEnabled = enabled
	if w.impl != nil {
		w.impl.setFullscreenButtonEnabled(enabled)
	}
	return w
}

// IsMinimised returns true if the window is minimised
func (w *WebviewWindow) IsMinimised() bool {
	if w.impl == nil {
		return false
	}
	return w.impl.isMinimised()
}

// IsMaximised returns true if the window is maximised
func (w *WebviewWindow) IsMaximised() bool {
	if w.impl == nil {
		return false
	}
	return w.impl.isMaximised()
}

// Size returns the size of the window
func (w *WebviewWindow) Size() (width int, height int) {
	if w.impl == nil {
		return 0, 0
	}
	return w.impl.size()
}

// IsFullscreen returns true if the window is fullscreen
func (w *WebviewWindow) IsFullscreen() bool {
	w.implLock.RLock()
	defer w.implLock.RUnlock()
	if w.impl == nil {
		return false
	}
	return w.impl.isFullscreen()
}

func (w *WebviewWindow) SetBackgroundColour(colour *RGBA) *WebviewWindow {
	w.options.BackgroundColour = colour
	if w.impl != nil {
		w.impl.setBackgroundColour(colour)
	}
	return w
}

func (w *WebviewWindow) handleMessage(message string) {
	w.info(message)
	// Check for special messages
	if message == "test" {
		w.SetTitle("Hello World")
	}
	w.info("ProcessMessage from front end:", message)

}

func (w *WebviewWindow) Center() {
	if w.impl == nil {
		return
	}
	w.impl.center()
}

func (w *WebviewWindow) On(eventType events.WindowEventType, callback func(ctx *WindowEventContext)) {
	eventID := uint(eventType)
	w.eventListenersLock.Lock()
	defer w.eventListenersLock.Unlock()
	w.eventListeners[eventID] = append(w.eventListeners[eventID], callback)
	if w.impl != nil {
		w.impl.on(eventID)
	}
}

func (w *WebviewWindow) handleWindowEvent(id uint) {
	w.eventListenersLock.RLock()
	for _, callback := range w.eventListeners[id] {
		go callback(blankWindowEventContext)
	}
	w.eventListenersLock.RUnlock()
}

func (w *WebviewWindow) Width() int {
	if w.impl == nil {
		return 0
	}
	return w.impl.width()
}

func (w *WebviewWindow) Height() int {
	if w.impl == nil {
		return 0
	}
	return w.impl.height()
}

func (w *WebviewWindow) Position() (int, int) {
	w.implLock.RLock()
	defer w.implLock.RUnlock()
	if w.impl == nil {
		return 0, 0
	}
	return w.impl.position()
}

func (w *WebviewWindow) Destroy() {
	if w.impl == nil {
		return
	}
	w.impl.destroy()
}

func (w *WebviewWindow) Reload() {
	if w.impl == nil {
		return
	}
	w.impl.reload()
}

func (w *WebviewWindow) ForceReload() {
	if w.impl == nil {
		return
	}
	w.impl.forceReload()
}

func (w *WebviewWindow) ToggleFullscreen() {
	if w.impl == nil {
		return
	}
	if w.IsFullscreen() {
		w.UnFullscreen()
	} else {
		w.Fullscreen()
	}
}

func (w *WebviewWindow) ToggleDevTools() {
	if w.impl == nil {
		return
	}
	w.impl.toggleDevTools()
}

func (w *WebviewWindow) ZoomReset() *WebviewWindow {
	if w.impl != nil {
		w.impl.zoomReset()
	}
	return w

}

func (w *WebviewWindow) ZoomIn() {
	if w.impl == nil {
		return
	}
	w.impl.zoomIn()
}

func (w *WebviewWindow) ZoomOut() {
	if w.impl == nil {
		return
	}
	w.impl.zoomOut()
}

func (w *WebviewWindow) Close() {
	if w.impl == nil {
		return
	}
	w.impl.close()
}

func (w *WebviewWindow) Minimize() {
	if w.impl == nil {
		return
	}
	w.impl.minimise()
}

func (w *WebviewWindow) Zoom() {
	if w.impl == nil {
		return
	}
	w.impl.zoom()
}

func (w *WebviewWindow) SetHTML(html string) *WebviewWindow {
	w.options.HTML = html
	if w.impl != nil {
		w.impl.setHTML(html)
	}
	return w
}

func (w *WebviewWindow) SetPosition(x, y int) *WebviewWindow {
	w.options.X = x
	w.options.Y = y
	if w.impl != nil {
		w.impl.setPosition(x, y)
	}
	return w
}

func (w *WebviewWindow) Minimise() *WebviewWindow {
	if w.impl == nil {
		w.options.StartState = WindowStateMinimised
		return w
	}
	if !w.IsMinimised() {
		w.impl.minimise()
	}
	return w
}

func (w *WebviewWindow) Maximise() *WebviewWindow {
	if w.impl == nil {
		w.options.StartState = WindowStateMaximised
		return w
	}
	if !w.IsMaximised() {
		w.disableSizeConstraints()
		w.impl.maximise()
	}
	return w
}

func (w *WebviewWindow) UnMinimise() {
	if w.impl == nil {
		return
	}
	w.impl.unminimise()
}

func (w *WebviewWindow) UnMaximise() {
	if w.impl == nil {
		return
	}
	w.enableSizeConstraints()
	w.impl.unmaximise()
}

func (w *WebviewWindow) UnFullscreen() {
	if w.impl == nil {
		return
	}
	w.enableSizeConstraints()
	w.impl.unfullscreen()
}

func (w *WebviewWindow) Restore() {
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

func (w *WebviewWindow) disableSizeConstraints() {
	if w.impl == nil {
		return
	}
	w.impl.setMinSize(0, 0)
	w.impl.setMaxSize(0, 0)
}

func (w *WebviewWindow) enableSizeConstraints() {
	if w.impl == nil {
		return
	}
	w.SetMinSize(w.options.MinWidth, w.options.MinHeight)
	w.SetMaxSize(w.options.MaxWidth, w.options.MaxHeight)
}

func (w *WebviewWindow) GetScreen() (*Screen, error) {
	if w.impl == nil {
		return nil, nil
	}
	return w.impl.getScreen()
}

func (w *WebviewWindow) SetFrameless(frameless bool) *WebviewWindow {
	w.options.Frameless = frameless
	if w.impl != nil {
		w.impl.setFrameless(frameless)
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
		listener(ctx)
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

func (w *WebviewWindow) RegisterContextMenu(name string, menu *Menu) {
	w.contextMenusLock.Lock()
	defer w.contextMenusLock.Unlock()
	w.contextMenus[name] = menu
}
