package application

import (
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

var (
	// ErrNativeWindowContentRequired is returned by NativeWindow.Run when the
	// window has no split layout to display. Native creation is deferred until
	// SetSplitView, or the NativeWindowOptions.SplitView option, supplies one.
	ErrNativeWindowContentRequired = errors.New("NativeWindow requires a MacSplitView content layout; call SetSplitView or set NativeWindowOptions.SplitView")
	// ErrNativeWindowEditorRequired is returned when a split layout offered to
	// a NativeWindow does not use AddTextEditor for its primary content pane.
	ErrNativeWindowEditorRequired = errors.New("NativeWindow requires exactly one AddTextEditor primary content pane")
	// ErrNativeWindowClosed is returned by operations on a closed window.
	ErrNativeWindowClosed = errors.New("NativeWindow has been closed")
)

// NativeWindowOptions configures an experimental native-content window.
// The deliberately small v3 surface can change when the common Window API is
// redesigned for v4.
type NativeWindowOptions struct {
	Name   string
	Title  string
	Width  int
	Height int
	Hidden bool

	MinWidth  int
	MinHeight int
	MaxWidth  int
	MaxHeight int

	DisableResize bool
	AlwaysOnTop   bool
	HideOnClose   bool

	InitialPosition WindowStartPosition
	X               int
	Y               int

	// Mac reuses the existing macOS window chrome options. WebviewPreferences
	// are ignored because a NativeWindow does not create a WKWebView, and
	// Mac.SplitView and Mac.Toolbar are ignored in favour of the fields below.
	Mac MacWindow

	// SplitView is the window's content layout, applied exactly as
	// SetSplitView would before the window is created. It must contain one
	// AddTextEditor primary pane. With it, a window created after App.Run is
	// constructed in one call; without it the window stays deferred until
	// SetSplitView supplies a layout.
	SplitView *MacSplitView

	// Toolbar is attached when the window is created, exactly as SetToolbar
	// would before creation. It is claimed before SplitView so it is in place
	// when the window is created.
	Toolbar *MacToolbar
}

// nativeWindowImplFactory creates the platform implementation. It is a
// variable so tests can substitute a fake without a native window.
var nativeWindowImplFactory = newNativeWindowImpl

// appIsRunning reports whether App.Run has started. It is read under the
// application's run lock, the same lock runOrDeferToAppRun uses.
func appIsRunning() bool {
	app := globalApplication
	if app == nil {
		return false
	}
	app.runLock.Lock()
	defer app.runLock.Unlock()
	return app.running
}

// validateNativeWindowSplit checks a layout for use as NativeWindow content:
// the shared split rules plus a native text editor as the primary pane.
func validateNativeWindowSplit(split *MacSplitView) error {
	if split == nil {
		return ErrNativeWindowContentRequired
	}
	if err := validateMacSplitView(split); err != nil {
		return err
	}
	primary := split.primaryPane()
	if primary == nil || primary.editor == nil {
		return ErrNativeWindowEditorRequired
	}
	return nil
}

type nativeWindowImpl interface {
	run() error
	show()
	hide()
	focus()
	close()
	isVisible() bool
	setTitle(string)
	nativeWindow() unsafe.Pointer
	setToolbar(*MacToolbar) error
	installSplitView() error
}

// NativeWindow is an experimental v3 window whose root content is native.
// It intentionally exposes no URL, JavaScript, zoom, or developer-tool API.
type NativeWindow struct {
	id      uint
	options NativeWindowOptions
	impl    nativeWindowImpl

	lock      sync.RWMutex
	destroyed bool

	// runLock serialises Run so the application's start-up queue and a
	// SetSplitView call in a running application cannot both create the
	// native window.
	runLock sync.Mutex

	toolbar *MacToolbar
	split   *MacSplitView
}

func newNativeWindow(options NativeWindowOptions) *NativeWindow {
	id := getWindowID()
	if options.Width == 0 {
		options.Width = 900
	}
	if options.Height == 0 {
		options.Height = 640
	}
	if options.Name == "" {
		options.Name = fmt.Sprintf("native-window-%d", id)
	}
	return &NativeWindow{id: id, options: options}
}

func (w *NativeWindow) ID() uint     { return w.id }
func (w *NativeWindow) Name() string { return w.options.Name }

// Run creates the native window. NativeWindowManager schedules it: a window
// created before App.Run is run when the application starts, and a window
// created afterwards is run as soon as it has content. A window without a
// split layout stays deferred and Run returns ErrNativeWindowContentRequired;
// SetSplitView then creates it.
//
// Run returns the first validation or installation error, which is also
// reported through Error. A window whose creation failed keeps no
// half-built implementation, so the caller can supply a valid layout and run
// it again. Calling Run on a created window is a no-op that returns nil.
func (w *NativeWindow) Run() error {
	w.runLock.Lock()
	defer w.runLock.Unlock()

	w.lock.RLock()
	destroyed, created, split := w.destroyed, w.impl != nil, w.split
	w.lock.RUnlock()
	if destroyed {
		return ErrNativeWindowClosed
	}
	if created {
		return nil
	}
	if err := validateNativeWindowSplit(split); err != nil {
		if !errors.Is(err, ErrNativeWindowContentRequired) {
			w.Error("NativeWindow.Run: %s", err)
		}
		return err
	}

	impl := nativeWindowImplFactory(w)
	w.lock.Lock()
	if w.destroyed {
		w.lock.Unlock()
		return ErrNativeWindowClosed
	}
	w.impl = impl
	w.lock.Unlock()
	if err := impl.run(); err != nil {
		w.lock.Lock()
		w.impl = nil
		w.lock.Unlock()
		w.Error("NativeWindow.Run: %s", err)
		return err
	}
	// Accessories added before the native window existed are attached once
	// the scheduled creation has run.
	w.flushPendingMacAccessories()
	return nil
}

// runIfDeferred creates a window whose native creation was deferred until it
// had content, once the application is running. Before App.Run the start-up
// queue creates it. The error is already reported by Run.
func (w *NativeWindow) runIfDeferred() error {
	if !appIsRunning() {
		return nil
	}
	return w.Run()
}

// Show shows the window, creating it first when creation was deferred and
// the application is running.
func (w *NativeWindow) Show() *NativeWindow {
	if w.implementation() == nil {
		if err := w.runIfDeferred(); err != nil {
			return w
		}
	}
	if impl := w.implementation(); impl != nil {
		InvokeSync(impl.show)
	}
	return w
}

func (w *NativeWindow) Hide() *NativeWindow {
	if impl := w.implementation(); impl != nil {
		InvokeSync(impl.hide)
	}
	return w
}

func (w *NativeWindow) Focus() {
	if impl := w.implementation(); impl != nil {
		InvokeSync(impl.focus)
	}
}

func (w *NativeWindow) Close() {
	w.lock.Lock()
	if w.destroyed {
		w.lock.Unlock()
		return
	}
	w.destroyed = true
	impl := w.impl
	w.lock.Unlock()
	if impl != nil {
		InvokeSync(impl.close)
	}
	if globalApplication != nil && globalApplication.NativeWindow != nil {
		globalApplication.NativeWindow.remove(w.id)
	}
}

func (w *NativeWindow) IsVisible() bool {
	if impl := w.implementation(); impl != nil {
		return InvokeSyncWithResult(impl.isVisible)
	}
	return false
}

func (w *NativeWindow) SetTitle(title string) *NativeWindow {
	w.lock.Lock()
	w.options.Title = title
	impl := w.impl
	w.lock.Unlock()
	if impl != nil {
		InvokeSync(func() { impl.setTitle(title) })
	}
	return w
}

func (w *NativeWindow) NativeWindow() unsafe.Pointer {
	if impl := w.implementation(); impl != nil {
		return InvokeSyncWithResult(impl.nativeWindow)
	}
	return nil
}

// SetSplitView supplies the window's native AppKit content layout, which
// must contain one AddTextEditor primary pane next to native sidebar,
// content-list, or inspector panes. A NativeWindow is created only once it
// has content: in a running application SetSplitView creates the window
// immediately (and shows it unless the Hidden option is set), returning any
// creation error; before App.Run the layout is queued for the start-up
// queue. The layout of a created window cannot be replaced; such a call
// returns ErrMacSplitViewAlreadyInstalled. Passing nil before creation
// clears a queued layout and releases it for use elsewhere.
func (w *NativeWindow) SetSplitView(split *MacSplitView) error {
	w.lock.Lock()
	if w.destroyed {
		w.lock.Unlock()
		return ErrNativeWindowClosed
	}
	if w.impl != nil {
		w.lock.Unlock()
		return ErrMacSplitViewAlreadyInstalled
	}
	if split == nil {
		previous := w.split
		w.split = nil
		w.lock.Unlock()
		releaseMacSplitViewOwnership(previous, w)
		return nil
	}
	if err := validateNativeWindowSplit(split); err != nil {
		w.lock.Unlock()
		return err
	}
	if err := claimMacSplitView(split, w); err != nil {
		w.lock.Unlock()
		return err
	}
	previous := w.split
	w.split = split
	w.lock.Unlock()
	if previous != nil && previous != split {
		releaseMacSplitViewOwnership(previous, w)
	}
	return w.runIfDeferred()
}

// SetToolbar attaches or replaces an NSToolbar. It uses the same MacToolbar
// type supported by WebviewWindow.
func (w *NativeWindow) SetToolbar(toolbar *MacToolbar) error {
	if toolbar != nil {
		if err := validateToolbarItems(toolbar.itemSnapshot()); err != nil {
			return err
		}
		if _, err := claimMacToolbar(toolbar, w); err != nil {
			return err
		}
	}
	w.lock.Lock()
	previous := w.toolbar
	w.toolbar = toolbar
	impl := w.impl
	w.lock.Unlock()
	if impl != nil {
		var err error
		InvokeSync(func() { err = impl.setToolbar(toolbar) })
		if err != nil {
			w.lock.Lock()
			w.toolbar = previous
			w.lock.Unlock()
			if toolbar != nil {
				releaseMacToolbarOwnership(toolbar, w)
			}
			return err
		}
	}
	if previous != nil && previous != toolbar {
		releaseMacToolbarOwnership(previous, w)
	}
	return nil
}

func (w *NativeWindow) implementation() nativeWindowImpl {
	w.lock.RLock()
	defer w.lock.RUnlock()
	if w.destroyed {
		return nil
	}
	return w.impl
}

func (w *NativeWindow) Error(message string, args ...any) {
	if globalApplication != nil {
		globalApplication.error("NativeWindow \"%s\": %s", w.Name(), fmt.Sprintf(message, args...))
	}
}

func (w *NativeWindow) macInspectorPane() *MacSplitPane {
	w.lock.RLock()
	split := w.split
	w.lock.RUnlock()
	if split == nil {
		return nil
	}
	return split.inspectorPane()
}

func (w *NativeWindow) macSplitOptions() MacWindow { return w.options.Mac }
