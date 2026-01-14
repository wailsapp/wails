package application

import (
	"fmt"
	"runtime"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"encoding/json"

	"github.com/leaanthony/u"
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
		nativeWindow() unsafe.Pointer
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
		setMenu(menu *Menu)
		snapAssist()
		setContentProtection(enabled bool)
	}
)

type WindowEvent struct {
	ctx       *WindowEventContext
	cancelled atomic.Bool
}

func (w *WindowEvent) Context() *WindowEventContext {
	return w.ctx
}

func NewWindowEvent() *WindowEvent {
	return &WindowEvent{}
}

func (w *WindowEvent) IsCancelled() bool {
	return w.cancelled.Load()
}

func (w *WindowEvent) Cancel() {
	w.cancelled.Store(true)
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
	keyBindings     map[string]func(Window)
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

	// unconditionallyClose marks the window to be unconditionally closed (atomic)
	unconditionallyClose uint32

	// Embedded panels management
	panels     map[uint]*WebviewPanel
	panelsLock sync.RWMutex
}

func (w *WebviewWindow) SetMenu(menu *Menu) {
	switch runtime.GOOS {
	case "darwin":
		return
	case "windows":
		w.options.Windows.Menu = menu
	case "linux":
		w.options.Linux.Menu = menu
	}
	if w.impl != nil {
		InvokeSync(func() {
			w.impl.setMenu(menu)
		})
	}
}

// EmitEvent emits a custom event with the specified name and associated data.
// It returns a boolean indicating whether the event was cancelled by a hook.
// The [CustomEvent.Sender] field will be set to the window name.
//
// If the given event name is registered, EmitEvent validates the data parameter
// against the expected data type. In case of a mismatch, EmitEvent reports an error
// to the registered error handler for the application and cancels the event.
func (w *WebviewWindow) EmitEvent(name string, data ...any) bool {
	event := &CustomEvent{
		Name:   name,
		Sender: w.Name(),
	}

	if len(data) == 1 {
		event.Data = data[0]
	} else if len(data) > 1 {
		event.Data = data
	}

	return globalApplication.Event.EmitEvent(event)
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
func (w *WebviewWindow) onApplicationEvent(
	eventType events.ApplicationEventType,
	callback func(*ApplicationEvent),
) {
	cancelFn := globalApplication.Event.OnApplicationEvent(eventType, callback)
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
		panels:         make(map[uint]*WebviewPanel),
	}

	result.setupEventMapping()

	// Listen for window closing events and cleanup panels
	result.OnWindowEvent(events.Common.WindowClosing, func(event *WindowEvent) {
		atomic.StoreUint32(&result.unconditionallyClose, 1)
		InvokeSync(result.destroyAllPanels)
		InvokeSync(result.markAsDestroyed)
		InvokeSync(result.impl.close)
		globalApplication.Window.Remove(result.id)
	})

	// Process keybindings
	if result.options.KeyBindings != nil {
		result.keyBindings = processKeyBindingOptions(result.options.KeyBindings)
	}

	return result
}

func processKeyBindingOptions(
	keyBindings map[string]func(window Window),
) map[string]func(window Window) {
	result := make(map[string]func(window Window))
	for key, callback := range keyBindings {
		// Parse the key to an accelerator
		acc, err := parseAccelerator(key)
		if err != nil {
			globalApplication.error("invalid keybinding: %w", err)
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

	// Start any panels that were added before the window was run
	w.runPanels()
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

func (w *WebviewWindow) SetContentProtection(b bool) Window {
	if w.impl == nil {
		w.options.ContentProtectionEnabled = b
	} else {
		InvokeSync(func() {
			w.impl.setContentProtection(b)
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
					w.Error("failed to start drag: %w", err)
				}
			})
		}
	case strings.HasPrefix(message, "wails:resize:"):
		if !w.IsFullscreen() {
			sl := strings.Split(message, ":")
			if len(sl) != 3 {
				w.Error("unknown message returned from dispatcher: %s", message)
				return
			}
			err := w.startResize(sl[2])
			if err != nil {
				w.Error("%w", err)
			}
		}
	case message == "wails:runtime:ready":
		w.emit(events.Common.WindowRuntimeReady)
		w.runtimeLoaded = true
		w.SetResizable(!w.options.DisableResize)
		for _, js := range w.pendingJS {
			w.ExecJS(js)
		}
		w.pendingJS = nil
	default:
		w.Error("unknown message sent via 'invoke' on frontend: %v", message)
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
func (w *WebviewWindow) OnWindowEvent(
	eventType events.WindowEventType,
	callback func(event *WindowEvent),
) func() {
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
		w.eventListeners[eventID] = slices.DeleteFunc(w.eventListeners[eventID], func(l *WindowEventListener) bool {
			return l == windowEventListener
		})
		w.eventListenersLock.Unlock()
	}
}

// RegisterHook registers a hook for the given window event
func (w *WebviewWindow) RegisterHook(
	eventType events.WindowEventType,
	callback func(event *WindowEvent),
) func() {
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
		w.eventHooks[eventID] = slices.DeleteFunc(w.eventHooks[eventID], func(l *WindowEventListener) bool {
			return l == windowEventHook
		})
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
		if thisEvent.IsCancelled() {
			return
		}
	}

	// Copy the w.eventListeners
	w.eventListenersLock.RLock()
	var tempListeners = slices.Clone(w.eventListeners[id])
	w.eventListenersLock.RUnlock()

	for _, listener := range tempListeners {
		go func() {
			if thisEvent.IsCancelled() {
				return
			}
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

// ToggleFrameless toggles the window between frameless and normal
func (w *WebviewWindow) ToggleFrameless() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(func() {
		w.SetFrameless(!w.options.Frameless)
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
	msg := fmt.Sprintf("window._wails.dispatchWailsEvent(%s);", event.ToJSON())
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
	args = append([]any{w.Name()}, args...)
	globalApplication.error("in window '%s': "+message, args...)
}

func (w *WebviewWindow) handleDragAndDropMessage(filenames []string, dropTarget *DropTargetDetails) {
	thisEvent := NewWindowEvent()
	ctx := newWindowEventContext()
	ctx.setDroppedFiles(filenames)
	if dropTarget != nil {
		ctx.setDropTargetDetails(dropTarget)
	}
	thisEvent.ctx = ctx

	listeners := w.eventListeners[uint(events.Common.WindowFilesDropped)]
	for _, listener := range listeners {
		if listener == nil {
			continue
		}
		listener.callback(thisEvent)
	}
}

func (w *WebviewWindow) OpenContextMenu(data *ContextMenuData) {
	// try application level context menu
	menu, ok := globalApplication.ContextMenu.Get(data.Id)
	if !ok {
		w.Error("no context menu found for id: %s", data.Id)
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

// NativeWindow returns the platform-specific native window handle
func (w *WebviewWindow) NativeWindow() unsafe.Pointer {
	if w.impl == nil || w.isDestroyed() {
		return nil
	}
	return w.impl.nativeWindow()
}

// shouldUnconditionallyClose returns whether the window should close unconditionally
func (w *WebviewWindow) shouldUnconditionallyClose() bool {
	return atomic.LoadUint32(&w.unconditionallyClose) != 0
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

	return globalApplication.KeyBinding.Process(acceleratorString, w)
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

func (w *WebviewWindow) InitiateFrontendDropProcessing(filenames []string, x int, y int) {
	if w.impl == nil || w.isDestroyed() {
		return
	}

	filenamesJSON, err := json.Marshal(filenames)
	if err != nil {
		w.Error("Error marshalling filenames for drop processing: %s", err)
		return
	}

	jsCall := fmt.Sprintf(
		"window.wails.Window.HandlePlatformFileDrop(%s, %d, %d);",
		string(filenamesJSON),
		x,
		y,
	)

	// Ensure JS is executed after runtime is loaded
	if !w.runtimeLoaded {
		w.pendingJS = append(w.pendingJS, jsCall)
		return
	}

	InvokeSync(func() {
		w.impl.execJS(jsCall)
	})
}

// HandleDragEnter is called when drag enters the window (Linux only, since GTK intercepts drag events)
func (w *WebviewWindow) HandleDragEnter() {
	if w.impl == nil || w.isDestroyed() || !w.runtimeLoaded {
		return
	}

	// Reset drag hover state for new drag session
	dragHover.lastSentX = 0
	dragHover.lastSentY = 0

	w.impl.execJS("window._wails.handleDragEnter();")
}

// Drag hover throttle state
var dragHover struct {
	lastSentX int
	lastSentY int
}

// HandleDragOver is called during drag-motion to update hover state in JS
// This is called from the GTK main thread, so we can call execJS directly
func (w *WebviewWindow) HandleDragOver(x int, y int) {
	if w.impl == nil || w.isDestroyed() || !w.runtimeLoaded {
		return
	}

	// Throttle: only send if moved at least 5 pixels
	dx := x - dragHover.lastSentX
	dy := y - dragHover.lastSentY
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	if dx < 5 && dy < 5 {
		return
	}
	dragHover.lastSentX = x
	dragHover.lastSentY = y

	// Use platform-specific zero-alloc implementation if available
	if impl, ok := w.impl.(interface{ execJSDragOver(x, y int) }); ok {
		impl.execJSDragOver(x, y)
	} else {
		w.impl.execJS(fmt.Sprintf("window._wails.handleDragOver(%d,%d)", x, y))
	}
}

// HandleDragLeave is called when drag leaves the window
func (w *WebviewWindow) HandleDragLeave() {
	if w.impl == nil || w.isDestroyed() || !w.runtimeLoaded {
		return
	}

	// Don't use InvokeSync - execJS already handles main thread dispatch internally
	w.impl.execJS("window._wails.handleDragLeave();")
}

// SnapAssist triggers the Windows Snap Assist feature by simulating Win+Z key combination.
// On Windows, this opens the snap layout options. On Linux and macOS, this is a no-op.
func (w *WebviewWindow) SnapAssist() {
	if w.impl == nil || w.isDestroyed() {
		return
	}
	InvokeSync(w.impl.snapAssist)
}

// ============================================================================
// Panel Management Methods
// ============================================================================

// NewPanel creates a new WebviewPanel with the given options and adds it to this window.
// The panel is a secondary webview that can be positioned anywhere within the window.
// This is similar to Electron's BrowserView or the deprecated webview tag.
//
// Example:
//
//	panel := window.NewPanel(application.WebviewPanelOptions{
//		X:      0,
//		Y:      0,
//		Width:  300,
//		Height: 400,
//		URL:    "https://example.com",
//	})
func (w *WebviewWindow) NewPanel(options WebviewPanelOptions) *WebviewPanel {
	panel := NewPanel(options)
	panel.parent = w

	w.panelsLock.Lock()
	w.panels[panel.id] = panel
	w.panelsLock.Unlock()

	// If window is already running, start the panel immediately
	if w.impl != nil && !w.isDestroyed() {
		InvokeSync(panel.run)
	}

	return panel
}

// GetPanel returns a panel by its name, or nil if not found.
func (w *WebviewWindow) GetPanel(name string) *WebviewPanel {
	w.panelsLock.RLock()
	defer w.panelsLock.RUnlock()

	for _, panel := range w.panels {
		if panel.name == name {
			return panel
		}
	}
	return nil
}

// GetPanelByID returns a panel by its ID, or nil if not found.
func (w *WebviewWindow) GetPanelByID(id uint) *WebviewPanel {
	w.panelsLock.RLock()
	defer w.panelsLock.RUnlock()
	return w.panels[id]
}

// GetPanels returns all panels attached to this window.
func (w *WebviewWindow) GetPanels() []*WebviewPanel {
	w.panelsLock.RLock()
	defer w.panelsLock.RUnlock()

	panels := make([]*WebviewPanel, 0, len(w.panels))
	for _, panel := range w.panels {
		panels = append(panels, panel)
	}
	return panels
}

// RemovePanel removes a panel from this window by its name.
// Returns true if the panel was found and removed.
func (w *WebviewWindow) RemovePanel(name string) bool {
	panel := w.GetPanel(name)
	if panel == nil {
		return false
	}
	panel.Destroy()
	return true
}

// RemovePanelByID removes a panel from this window by its ID.
// Returns true if the panel was found and removed.
func (w *WebviewWindow) RemovePanelByID(id uint) bool {
	panel := w.GetPanelByID(id)
	if panel == nil {
		return false
	}
	panel.Destroy()
	return true
}

// removePanel is called by WebviewPanel.Destroy() to remove itself from the parent
func (w *WebviewWindow) removePanel(id uint) {
	w.panelsLock.Lock()
	defer w.panelsLock.Unlock()
	delete(w.panels, id)
}

// runPanels starts all panels that haven't been started yet.
// This is called after the window's impl is created.
func (w *WebviewWindow) runPanels() {
	// Collect panels under lock, then run them outside the lock
	w.panelsLock.RLock()
	panels := make([]*WebviewPanel, 0, len(w.panels))
	for _, panel := range w.panels {
		if panel.impl == nil {
			panels = append(panels, panel)
		}
	}
	w.panelsLock.RUnlock()

	for _, panel := range panels {
		panel.run()
	}
}

// destroyAllPanels destroys all panels in this window.
// This is called when the window is closing.
func (w *WebviewWindow) destroyAllPanels() {
	w.panelsLock.Lock()
	panels := make([]*WebviewPanel, 0, len(w.panels))
	for _, panel := range w.panels {
		panels = append(panels, panel)
	}
	w.panelsLock.Unlock()

	for _, panel := range panels {
		panel.Destroy()
	}
}
