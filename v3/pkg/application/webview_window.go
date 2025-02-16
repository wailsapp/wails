package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"slices"
	"strings"
	"sync"

	"github.com/leaanthony/u"

	"github.com/samber/lo"
	"github.com/wailsapp/wails/v3/internal/assetserver"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Enabled means the feature should be enabled
var Enabled = u.True

// Disabled means the feature should be disabled
var Disabled = u.False

// LRTB is a struct that holds Left, Right, Top, Bottom values
type LRTB struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}

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
		destroy()
		reload()
		forceReload()
		openDevTools()
		zoomReset()
		zoomIn()
		zoomOut()
		getZoom() float64
		setZoom(zoom float64)
		close()
		zoom()
		setHTML(html string)
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
		isFocused() bool
		focus()
		show()
		hide()
		getScreen() (*Screen, error)
		setFrameless(bool)
		openContextMenu(menu *Menu, data *ContextMenuData)
		nativeWindowHandle() uintptr
		startDrag() error
		startResize(border string) error
		print() error
		setEnabled(enabled bool)
		physicalBounds() Rect
		setPhysicalBounds(physicalBounds Rect)
		bounds() Rect
		setBounds(bounds Rect)
		position() (int, int)
		setPosition(x int, y int)
		relativePosition() (int, int)
		setRelativePosition(x int, y int)
		flash(enabled bool)
		handleKeyEvent(acceleratorString string)
		getBorderSizes() *LRTB
		setMinimiseButtonState(state ButtonState)
		setMaximiseButtonState(state ButtonState)
		setCloseButtonState(state ButtonState)
		isIgnoreMouseEvents() bool
		setIgnoreMouseEvents(ignore bool)
		cut()
		copy()
		paste()
		undo()
		delete()
		selectAll()
		redo()
		showMenuBar()
		hideMenuBar()
		toggleMenuBar()
	}
)

type WindowEvent struct {
	ctx       *WindowEventContext
	Cancelled bool
}

func (w *WindowEvent) Context() *WindowEventContext {
	return w.ctx
}

func NewWindowEvent() *WindowEvent {
	return &WindowEvent{}
}

func (w *WindowEvent) Cancel() {
	w.Cancelled = true
}

type WindowEventListener struct {
	callback func(event *WindowEvent)
}

type WebviewWindow struct {
	options WebviewWindowOptions
	impl    webviewWindowImpl
	id      uint

	eventListeners     map[uint][]*WindowEventListener
	eventListenersLock sync.RWMutex
	eventHooks         map[uint][]*WindowEventListener
	eventHooksLock     sync.RWMutex

	// A map of listener cancellation functions
	cancellersLock sync.RWMutex
	cancellers     []func()

	// keyBindings holds the keybindings for the window
	keyBindings     map[string]func(*WebviewWindow)
	keyBindingsLock sync.RWMutex

	// menuBindings holds the menu bindings for the window
	menuBindings     map[string]*MenuItem
	menuBindingsLock sync.RWMutex

	// Indicates that the window is destroyed
	destroyed     bool
	destroyedLock sync.RWMutex

	// Flags for managing the runtime
	// runtimeLoaded indicates that the runtime has been loaded
	runtimeLoaded bool
	// pendingJS holds JS that was sent to the window before the runtime was loaded
	pendingJS []string

	// unconditionallyClose marks the window to be unconditionally closed
	unconditionallyClose bool
}

// EmitEvent emits an event from the window
func (w *WebviewWindow) EmitEvent(name string, data ...any) {
	globalApplication.emitEvent(&CustomEvent{
		Name:   name,
		Data:   data,
		Sender: w.Name(),
	})
}

var windowID uint
var windowIDLock sync.RWMutex

func getWindowID() uint {
	windowIDLock.Lock()
	defer windowIDLock.Unlock()
	windowID++
	return windowID
}

// FIXME: This should like be an interface method (TDM)
// Use onApplicationEvent to register a callback for an application event from a window.
// This will handle tidying up the callback when the window is destroyed
func (w *WebviewWindow) onApplicationEvent(eventType events.ApplicationEventType, callback func(*ApplicationEvent)) {
	cancelFn := globalApplication.OnApplicationEvent(eventType, callback)
	w.addCancellationFunction(cancelFn)
}

func (w *WebviewWindow) markAsDestroyed() {
	w.destroyedLock.Lock()
	defer w.destroyedLock.Unlock()
	w.destroyed = true
}

func (w *WebviewWindow) setupEventMapping() {

	var mapping map[events.WindowEventType]events.WindowEventType
	switch runtime.GOOS {
	case "darwin":
		mapping = w.options.Mac.EventMapping
	case "windows":
		mapping = w.options.Windows.EventMapping
	case "linux":
		// TBD
	}
	if mapping == nil {
		mapping = events.DefaultWindowEventMapping()
	}

	for source, target := range mapping {
		source := source
		target := target
		w.OnWindowEvent(source, func(event *WindowEvent) {
			w.emit(target)
		})
	}
}

// NewWindow creates a new window with the given options
func NewWindow(options WebviewWindowOptions) *WebviewWindow {

	thisWindowID := getWindowID()

	if options.Width == 0 {
		options.Width = 800
	}
	if options.Height == 0 {
		options.Height = 600
	}
	if options.URL == "" {
		options.URL = "/"
	}

	if options.Name == "" {
		options.Name = fmt.Sprintf("window-%d", thisWindowID)
	}

	result := &WebviewWindow{
		id:             thisWindowID,
		options:        options,
		eventListeners: make(map[uint][]*WindowEventListener),
		eventHooks:     make(map[uint][]*WindowEventListener),
		menuBindings:   make(map[string]*MenuItem),
	}

	result.setupEventMapping()

	// Listen for window closing events and de
	result.OnWindowEvent(events.Common.WindowClosing, func(event *WindowEvent) {
		result.unconditionallyClose = true
		InvokeSync(result.markAsDestroyed)
		InvokeSync(result.impl.close)
		globalApplication.deleteWindowByID(result.id)
	})

	// Process keybindings
	if result.options.KeyBindings != nil {
		result.keyBindings = processKeyBindingOptions(result.options.KeyBindings)
	}

	return result
}

func processKeyBindingOptions(keyBindings map[string]func(window *WebviewWindow)) map[string]func(window *WebviewWindow) {
	result := make(map[string]func(window *WebviewWindow))
	for key, callback := range keyBindings {
		// Parse the key to an accelerator
		acc, err := parseAccelerator(key)
		if err != nil {
			globalApplication.error("Invalid keybinding: %s", err.Error())
			continue
		}
		result[acc.String()] = callback
		globalApplication.debug("Added Keybinding", "accelerator", acc.String())
	}
	return result
}

func (w *WebviewWindow) addCancellationFunction(canceller func()) {
	w.cancellersLock.Lock()
	defer w.cancellersLock.Unlock()
	w.cancellers = append(w.cancellers, canceller)
}

// formatJS ensures the 'data' provided marshals to valid json or panics
func (w *WebviewWindow) formatJS(f string, callID string, data string) string {
	j, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(f, callID, j)
}

func (w *WebviewWindow) CallError(callID string, result string) {
	if w.impl != nil {
		w.impl.execJS(w.formatJS("_wails.callErrorHandler('%s', %s);", callID, result))
	}
}

func (w *WebviewWindow) CallResponse(callID string, result string) {
	if w.impl != nil {
		w.impl.execJS(w.formatJS("_wails.callResultHandler('%s', %s, true);", callID, result))
	}
}

func (w *WebviewWindow) DialogError(dialogID string, result string) {
	if w.impl != nil {
		w.impl.execJS(w.formatJS("_wails.dialogErrorCallback('%s', %s);", dialogID, result))
	}
}

func (w *WebviewWindow) DialogResponse(dialogID string, result string, isJSON bool) {
	if w.impl != nil {
		if isJSON {
			w.impl.execJS(w.formatJS("_wails.dialogResultCallback('%s', %s, true);", dialogID, result))
		} else {
			w.impl.execJS(fmt.Sprintf("_wails.dialogResultCallback('%s', '%s', false);", dialogID, result))
		}
	}
}

func (w *WebviewWindow) ID() uint {
	return w.id
}

// SetTitle sets the title of the window
func (w *WebviewWindow) SetTitle(title string) Window {
	w.options.Title = title
	if w.impl != nil {
		InvokeSync(func() {
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
func (w *WebviewWindow) SetSize(width, height int) Window {
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
		InvokeSync(func() {
			w.impl.setSize(width, height)
		})
	}
	return w
}

func (w *WebviewWindow) Run() {
	if w.impl != nil {
		return
	}
	w.impl = newWindowImpl(w)

	InvokeSync(w.impl.run)
}

// SetAlwaysOnTop sets the window to be always on top.
func (w *WebviewWindow) SetAlwaysOnTop(b bool) Window {
	w.options.AlwaysOnTop = b
	if w.impl != nil {
		InvokeSync(func() {
			w.impl.setAlwaysOnTop(b)
		})
	}
	return w
}

// Show shows the window.
func (w *WebviewWindow) Show() Window {
	if globalApplication.impl == nil {
		return w
	}
	if w.impl == nil || w.isDestroyed() {
		InvokeSync(w.Run)
		return w
	}
	w.options.Hidden = false
	InvokeSync(w.impl.show)
	return w
}

// Hide hides the window.
func (w *WebviewWindow) Hide() Window {
	w.options.Hidden = true
	if w.impl != nil {
		InvokeSync(w.impl.hide)
	}
	return w
}

func (w *WebviewWindow) SetURL(s string) Window {
	url, _ := assetserver.GetStartURL(s)
	w.options.URL = url
	if w.impl != nil {
		InvokeSync(func() {
			w.impl.setURL(url)
		})
	}
	return w
}

func (w *WebviewWindow) GetBorderSizes() *LRTB {
	if w.impl != nil {
		return InvokeSyncWithResult(w.impl.getBorderSizes)
	}
	return &LRTB{}
}

// SetZoom sets the zoom level of the window.
func (w *WebviewWindow) SetZoom(magnification float64) Window {
	w.options.Zoom = magnification
	if w.impl != nil {
		InvokeSync(func() {
			w.impl.setZoom(magnification)
		})
	}
	return w
}

// GetZoom returns the current zoom level of the window.
func (w *WebviewWindow) GetZoom() float64 {
	if w.impl != nil {
		return InvokeSyncWithResult(w.impl.getZoom)
	}
	return 1
}

// SetResizable sets whether the window is resizable.
func (w *WebviewWindow) SetResizable(b bool) Window {
	w.options.DisableResize = !b
	if w.impl != nil {
		InvokeSync(func() {
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
func (w *WebviewWindow) SetMinSize(minWidth, minHeight int) Window {
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
			InvokeSync(func() {
				w.impl.setSize(newWidth, newHeight)
			})
		}
		InvokeSync(func() {
			w.impl.setMinSize(minWidth, minHeight)
		})
	}
	return w
}

// SetMaxSize sets the maximum size of the window.
func (w *WebviewWindow) SetMaxSize(maxWidth, maxHeight int) Window {
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
			InvokeSync(func() {
				w.impl.setSize(newWidth, newHeight)
			})
		}
		InvokeSync(func() {
			w.impl.setMaxSize(maxWidth, maxHeight)
		})
	}
	return w
}

// ExecJS executes the given javascript in the context of the window.
func (w *WebviewWindow) ExecJS(js string) {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	if w.runtimeLoaded {
		InvokeSync(func() {
			w.impl.execJS(js)
		})
	} else {
		w.pendingJS = append(w.pendingJS, js)
	}
}

// Fullscreen sets the window to fullscreen mode. Min/Max size constraints are disabled.
func (w *WebviewWindow) Fullscreen() Window {
	if w.impl == nil || w.isDestroyed() {
		w.options.StartState = WindowStateFullscreen
		return w
	}
	if !w.IsFullscreen() {
		w.DisableSizeConstraints()
		InvokeSync(w.impl.fullscreen)
	}
	return w
}

func (w *WebviewWindow) SetMinimiseButtonState(state ButtonState) Window {
	w.options.MinimiseButtonState = state
	if w.impl != nil {
		InvokeSync(func() {
			w.impl.setMinimiseButtonState(state)
		})
	}
	return w
}

func (w *WebviewWindow) SetMaximiseButtonState(state ButtonState) Window {
	w.options.MaximiseButtonState = state
	if w.impl != nil {
		InvokeSync(func() {
			w.impl.setMaximiseButtonState(state)
		})
	}
	return w
}

func (w *WebviewWindow) SetCloseButtonState(state ButtonState) Window {
	w.options.CloseButtonState = state
	if w.impl != nil {
		InvokeSync(func() {
			w.impl.setCloseButtonState(state)
		})
	}
	return w
}

// Flash flashes the window's taskbar button/icon.
// Useful to indicate that attention is required. Windows only.
func (w *WebviewWindow) Flash(enabled bool) {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(func() {
		w.impl.flash(enabled)
	})
}

// IsMinimised returns true if the window is minimised
func (w *WebviewWindow) IsMinimised() bool {
	if w.impl == nil || w.isDestroyed() {
		return false
	}
	return InvokeSyncWithResult(w.impl.isMinimised)
}

// IsVisible returns true if the window is visible
func (w *WebviewWindow) IsVisible() bool {
	if w.impl == nil || w.isDestroyed() {
		return false
	}
	return InvokeSyncWithResult(w.impl.isVisible)
}

// IsMaximised returns true if the window is maximised
func (w *WebviewWindow) IsMaximised() bool {
	if w.impl == nil || w.isDestroyed() {
		return false
	}
	return InvokeSyncWithResult(w.impl.isMaximised)
}

// Size returns the size of the window
func (w *WebviewWindow) Size() (int, int) {
	if w.impl == nil || w.isDestroyed() {
		return 0, 0
	}
	var width, height int
	InvokeSync(func() {
		width, height = w.impl.size()
	})
	return width, height
}

// IsFocused returns true if the window is currently focused
func (w *WebviewWindow) IsFocused() bool {
	if w.impl == nil || w.isDestroyed() {
		return false
	}
	return InvokeSyncWithResult(w.impl.isFocused)
}

// IsFullscreen returns true if the window is fullscreen
func (w *WebviewWindow) IsFullscreen() bool {
	if w.impl == nil || w.isDestroyed() {
		return false
	}
	return InvokeSyncWithResult(w.impl.isFullscreen)
}

// SetBackgroundColour sets the background colour of the window
func (w *WebviewWindow) SetBackgroundColour(colour RGBA) Window {
	w.options.BackgroundColour = colour
	if w.impl != nil {
		InvokeSync(func() {
			w.impl.setBackgroundColour(colour)
		})
	}
	return w
}

func (w *WebviewWindow) HandleMessage(message string) {
	// Check for special messages
	switch true {
	case message == "wails:drag":
		if !w.IsFullscreen() {
			InvokeSync(func() {
				err := w.startDrag()
				if err != nil {
					w.Error("Failed to start drag: %s", err)
				}
			})
		}
	case strings.HasPrefix(message, "wails:resize:"):
		if !w.IsFullscreen() {
			sl := strings.Split(message, ":")
			if len(sl) != 3 {
				w.Error("Unknown message returned from dispatcher", "message", message)
				return
			}
			err := w.startResize(sl[2])
			if err != nil {
				w.Error(err.Error())
			}
		}
	case message == "wails:runtime:ready":
		w.emit(events.Common.WindowRuntimeReady)
		w.runtimeLoaded = true
		w.SetResizable(!w.options.DisableResize)
		for _, js := range w.pendingJS {
			w.ExecJS(js)
		}
	default:
		w.Error("Unknown message sent via 'invoke' on frontend: %v", message)
	}
}

func (w *WebviewWindow) startResize(border string) error {
	if w.impl == nil || w.isDestroyed() {
		return nil
	}
	return InvokeSyncWithResult(func() error {
		return w.impl.startResize(border)
	})
}

// Center centers the window on the screen
func (w *WebviewWindow) Center() {
	if w.impl == nil || w.isDestroyed() {
		w.options.InitialPosition = WindowCentered
		return
	}
	InvokeSync(w.impl.center)
}

// OnWindowEvent registers a callback for the given window event
func (w *WebviewWindow) OnWindowEvent(eventType events.WindowEventType, callback func(event *WindowEvent)) func() {
	eventID := uint(eventType)
	windowEventListener := &WindowEventListener{
		callback: callback,
	}
	w.eventListenersLock.Lock()
	w.eventListeners[eventID] = append(w.eventListeners[eventID], windowEventListener)
	w.eventListenersLock.Unlock()
	if w.impl != nil {
		w.impl.on(eventID)
	}

	return func() {
		// Check if eventListener is already locked
		w.eventListenersLock.Lock()
		w.eventListeners[eventID] = lo.Without(w.eventListeners[eventID], windowEventListener)
		w.eventListenersLock.Unlock()
	}
}

// RegisterHook registers a hook for the given window event
func (w *WebviewWindow) RegisterHook(eventType events.WindowEventType, callback func(event *WindowEvent)) func() {
	eventID := uint(eventType)
	w.eventHooksLock.Lock()
	defer w.eventHooksLock.Unlock()
	windowEventHook := &WindowEventListener{
		callback: callback,
	}
	w.eventHooks[eventID] = append(w.eventHooks[eventID], windowEventHook)

	return func() {
		w.eventHooksLock.Lock()
		defer w.eventHooksLock.Unlock()
		w.eventHooks[eventID] = lo.Without(w.eventHooks[eventID], windowEventHook)
	}
}

func (w *WebviewWindow) HandleWindowEvent(id uint) {
	// Get hooks
	w.eventHooksLock.RLock()
	hooks := w.eventHooks[id]
	w.eventHooksLock.RUnlock()

	// Create new WindowEvent
	thisEvent := NewWindowEvent()

	for _, thisHook := range hooks {
		thisHook.callback(thisEvent)
		if thisEvent.Cancelled {
			return
		}
	}

	// Copy the w.eventListeners
	w.eventListenersLock.RLock()
	var tempListeners = slices.Clone(w.eventListeners[id])
	w.eventListenersLock.RUnlock()

	for _, listener := range tempListeners {
		go func() {
			defer handlePanic()
			listener.callback(thisEvent)
		}()
	}
	w.dispatchWindowEvent(id)
}

// Width returns the width of the window
func (w *WebviewWindow) Width() int {
	if w.impl == nil || w.isDestroyed() {
		return 0
	}
	return InvokeSyncWithResult(w.impl.width)
}

// Height returns the height of the window
func (w *WebviewWindow) Height() int {
	if w.impl == nil || w.isDestroyed() {
		return 0
	}
	return InvokeSyncWithResult(w.impl.height)
}

// PhysicalBounds returns the physical bounds of the window
func (w *WebviewWindow) PhysicalBounds() Rect {
	if w.impl == nil || w.isDestroyed() {
		return Rect{}
	}
	var rect Rect
	InvokeSync(func() {
		rect = w.impl.physicalBounds()
	})
	return rect
}

// SetPhysicalBounds sets the physical bounds of the window
func (w *WebviewWindow) SetPhysicalBounds(physicalBounds Rect) {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(func() {
		w.impl.setPhysicalBounds(physicalBounds)
	})
}

// Bounds returns the DIP bounds of the window
func (w *WebviewWindow) Bounds() Rect {
	if w.impl == nil || w.isDestroyed() {
		return Rect{}
	}
	var rect Rect
	InvokeSync(func() {
		rect = w.impl.bounds()
	})
	return rect
}

// SetBounds sets the DIP bounds of the window
func (w *WebviewWindow) SetBounds(bounds Rect) {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(func() {
		w.impl.setBounds(bounds)
	})
}

// Position returns the absolute position of the window
func (w *WebviewWindow) Position() (int, int) {
	if w.impl == nil || w.isDestroyed() {
		return 0, 0
	}
	var x, y int
	InvokeSync(func() {
		x, y = w.impl.position()
	})
	return x, y
}

// SetPosition sets the absolute position of the window.
func (w *WebviewWindow) SetPosition(x int, y int) {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(func() {
		w.impl.setPosition(x, y)
	})
}

// RelativePosition returns the position of the window relative to the screen WorkArea on which it is
func (w *WebviewWindow) RelativePosition() (int, int) {
	if w.impl == nil || w.isDestroyed() {
		return 0, 0
	}
	var x, y int
	InvokeSync(func() {
		x, y = w.impl.relativePosition()
	})
	return x, y
}

// SetRelativePosition sets the position of the window relative to the screen WorkArea on which it is.
func (w *WebviewWindow) SetRelativePosition(x, y int) Window {
	w.options.X = x
	w.options.Y = y
	if w.impl != nil {
		InvokeSync(func() {
			w.impl.setRelativePosition(x, y)
		})
	}
	return w
}

func (w *WebviewWindow) destroy() {
	if w.impl == nil || w.isDestroyed() {
		return
	}

	// Cancel the callbacks
	for _, cancelFunc := range w.cancellers {
		cancelFunc()
	}

	InvokeSync(w.impl.destroy)
}

// Reload reloads the page assets
func (w *WebviewWindow) Reload() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(w.impl.reload)
}

// ForceReload forces the window to reload the page assets
func (w *WebviewWindow) ForceReload() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(w.impl.forceReload)
}

// ToggleFullscreen toggles the window between fullscreen and normal
func (w *WebviewWindow) ToggleFullscreen() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(func() {
		if w.IsFullscreen() {
			w.UnFullscreen()
		} else {
			w.Fullscreen()
		}
	})
}

// ToggleMaximise toggles the window between maximised and normal
func (w *WebviewWindow) ToggleMaximise() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(func() {
		if w.IsMaximised() {
			w.UnMaximise()
		} else {
			w.Maximise()
		}
	})
}

func (w *WebviewWindow) OpenDevTools() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(w.impl.openDevTools)
}

// ZoomReset resets the zoom level of the webview content to 100%
func (w *WebviewWindow) ZoomReset() Window {
	if w.impl != nil {
		InvokeSync(w.impl.zoomReset)
	}
	return w
}

// ZoomIn increases the zoom level of the webview content
func (w *WebviewWindow) ZoomIn() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(w.impl.zoomIn)
}

// ZoomOut decreases the zoom level of the webview content
func (w *WebviewWindow) ZoomOut() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(w.impl.zoomOut)
}

// Close closes the window
func (w *WebviewWindow) Close() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(func() {
		w.emit(events.Common.WindowClosing)
	})
}

func (w *WebviewWindow) Zoom() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(w.impl.zoom)
}

// SetHTML sets the HTML of the window to the given html string.
func (w *WebviewWindow) SetHTML(html string) Window {
	w.options.HTML = html
	if w.impl != nil {
		InvokeSync(func() {
			w.impl.setHTML(html)
		})
	}
	return w
}

// Minimise minimises the window.
func (w *WebviewWindow) Minimise() Window {
	if w.impl == nil || w.isDestroyed() {
		w.options.StartState = WindowStateMinimised
		return w
	}
	if !w.IsMinimised() {
		InvokeSync(w.impl.minimise)
	}
	return w
}

// Maximise maximises the window. Min/Max size constraints are disabled.
func (w *WebviewWindow) Maximise() Window {
	if w.impl == nil || w.isDestroyed() {
		w.options.StartState = WindowStateMaximised
		return w
	}
	if !w.IsMaximised() {
		w.DisableSizeConstraints()
		InvokeSync(w.impl.maximise)
	}
	return w
}

// UnMinimise un-minimises the window. Min/Max size constraints are re-enabled.
func (w *WebviewWindow) UnMinimise() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	if w.IsMinimised() {
		InvokeSync(w.impl.unminimise)
	}
}

// UnMaximise un-maximises the window. Min/Max size constraints are re-enabled.
func (w *WebviewWindow) UnMaximise() {
	if w.IsMaximised() {
		w.EnableSizeConstraints()
		InvokeSync(w.impl.unmaximise)
	}
}

// UnFullscreen un-fullscreens the window. Min/Max size constraints are re-enabled.
func (w *WebviewWindow) UnFullscreen() {
	if w.IsFullscreen() {
		w.EnableSizeConstraints()
		InvokeSync(w.impl.unfullscreen)
	}
}

// Restore restores the window to its previous state if it was previously minimised, maximised or fullscreen.
func (w *WebviewWindow) Restore() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(func() {
		if w.IsMinimised() {
			w.UnMinimise()
		} else if w.IsFullscreen() {
			w.UnFullscreen()
		} else if w.IsMaximised() {
			w.UnMaximise()
		}
	})
}

func (w *WebviewWindow) DisableSizeConstraints() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(func() {
		if w.options.MinWidth > 0 && w.options.MinHeight > 0 {
			w.impl.setMinSize(0, 0)
		}
		if w.options.MaxWidth > 0 && w.options.MaxHeight > 0 {
			w.impl.setMaxSize(0, 0)
		}
	})
}

func (w *WebviewWindow) EnableSizeConstraints() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(func() {
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
	if w.impl == nil || w.isDestroyed() {
		return nil, nil
	}
	return InvokeSyncWithResultAndError(w.impl.getScreen)
}

// SetFrameless removes the window frame and title bar
func (w *WebviewWindow) SetFrameless(frameless bool) Window {
	w.options.Frameless = frameless
	if w.impl != nil {
		InvokeSync(func() {
			w.impl.setFrameless(frameless)
		})
	}
	return w
}

func (w *WebviewWindow) DispatchWailsEvent(event *CustomEvent) {
	msg := fmt.Sprintf("_wails.dispatchWailsEvent(%s);", event.ToJSON())
	w.ExecJS(msg)
}

func (w *WebviewWindow) dispatchWindowEvent(id uint) {
	// TODO: Make this more efficient by keeping a list of which events have been registered
	// and only dispatching those.
	jsEvent := &CustomEvent{
		Name: events.JSEvent(id),
	}
	w.DispatchWailsEvent(jsEvent)
}

func (w *WebviewWindow) Info(message string, args ...any) {
	var messageArgs []interface{}
	messageArgs = append(messageArgs, args...)
	messageArgs = append(messageArgs, "sender", w.Name())
	globalApplication.info(message, messageArgs...)
}

func (w *WebviewWindow) Error(message string, args ...any) {
	var messageArgs []interface{}
	messageArgs = append(messageArgs, args...)
	messageArgs = append(messageArgs, "sender", w.Name())
	globalApplication.error(message, messageArgs...)
}

func (w *WebviewWindow) HandleDragAndDropMessage(filenames []string) {
	thisEvent := NewWindowEvent()
	ctx := newWindowEventContext()
	ctx.setDroppedFiles(filenames)
	thisEvent.ctx = ctx
	for _, listener := range w.eventListeners[uint(events.Common.WindowFilesDropped)] {
		listener.callback(thisEvent)
	}
}

func (w *WebviewWindow) OpenContextMenu(data *ContextMenuData) {
	// try application level context menu
	menu, ok := globalApplication.getContextMenu(data.Id)
	if !ok {
		w.Error("No context menu found for id: %s", data.Id)
		return
	}
	menu.setContextData(data)
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(func() {
		w.impl.openContextMenu(menu.Menu, data)
	})
}

// NativeWindowHandle returns the platform native window handle for the window.
func (w *WebviewWindow) NativeWindowHandle() (uintptr, error) {
	if w.impl == nil || w.isDestroyed() {
		return 0, errors.New("native handle unavailable as window is not running")
	}
	return w.impl.nativeWindowHandle(), nil
}

func (w *WebviewWindow) Focus() {
	InvokeSync(w.impl.focus)
}

func (w *WebviewWindow) emit(eventType events.WindowEventType) {
	windowEvents <- &windowEvent{
		WindowID: w.id,
		EventID:  uint(eventType),
	}
}

func (w *WebviewWindow) startDrag() error {
	if w.impl == nil || w.isDestroyed() {
		return nil
	}
	return InvokeSyncWithError(w.impl.startDrag)
}

func (w *WebviewWindow) Print() error {
	if w.impl == nil || w.isDestroyed() {
		return nil
	}
	return InvokeSyncWithError(w.impl.print)
}

func (w *WebviewWindow) SetEnabled(enabled bool) {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(func() {
		w.impl.setEnabled(enabled)
	})
}

func (w *WebviewWindow) processKeyBinding(acceleratorString string) bool {
	// Check menu bindings
	if w.menuBindings != nil {
		w.menuBindingsLock.RLock()
		defer w.menuBindingsLock.RUnlock()
		if menuItem := w.menuBindings[acceleratorString]; menuItem != nil {
			menuItem.handleClick()
			return true
		}
	}

	// Check key bindings
	if w.keyBindings != nil {
		w.keyBindingsLock.RLock()
		defer w.keyBindingsLock.RUnlock()
		if callback := w.keyBindings[acceleratorString]; callback != nil {
			// Execute callback
			go func() {
				defer handlePanic()
				callback(w)
			}()
			return true
		}
	}

	return globalApplication.processKeyBinding(acceleratorString, w)
}

func (w *WebviewWindow) HandleKeyEvent(acceleratorString string) {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(func() {
		w.impl.handleKeyEvent(acceleratorString)
	})
}

func (w *WebviewWindow) isDestroyed() bool {
	w.destroyedLock.RLock()
	defer w.destroyedLock.RUnlock()
	return w.destroyed
}

func (w *WebviewWindow) removeMenuBinding(a *accelerator) {
	w.menuBindingsLock.Lock()
	defer w.menuBindingsLock.Unlock()
	w.menuBindings[a.String()] = nil
}

func (w *WebviewWindow) addMenuBinding(a *accelerator, menuItem *MenuItem) {
	w.menuBindingsLock.Lock()
	defer w.menuBindingsLock.Unlock()
	w.menuBindings[a.String()] = menuItem
}

func (w *WebviewWindow) IsIgnoreMouseEvents() bool {
	if w.impl == nil || w.isDestroyed() {
		return false
	}
	return InvokeSyncWithResult(w.impl.isIgnoreMouseEvents)
}

func (w *WebviewWindow) SetIgnoreMouseEvents(ignore bool) Window {
	w.options.IgnoreMouseEvents = ignore
	if w.impl == nil || w.isDestroyed() {
		return w
	}
	InvokeSync(func() {
		w.impl.setIgnoreMouseEvents(ignore)
	})
	return w
}

func (w *WebviewWindow) cut() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	w.impl.cut()
}

func (w *WebviewWindow) copy() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	w.impl.copy()
}

func (w *WebviewWindow) paste() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	w.impl.paste()
}

func (w *WebviewWindow) selectAll() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	w.impl.selectAll()
}

func (w *WebviewWindow) undo() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	w.impl.undo()
}

func (w *WebviewWindow) delete() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	w.impl.delete()
}

func (w *WebviewWindow) redo() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	w.impl.redo()
}

// ShowMenuBar shows the menu bar for the window.
func (w *WebviewWindow) ShowMenuBar() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(w.impl.showMenuBar)
}

// HideMenuBar hides the menu bar for the window.
func (w *WebviewWindow) HideMenuBar() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(w.impl.hideMenuBar)
}

// ToggleMenuBar toggles the menu bar for the window.
func (w *WebviewWindow) ToggleMenuBar() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(w.impl.toggleMenuBar)
}
