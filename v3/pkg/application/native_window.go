package application

import (
	"fmt"
	"sync"
	"unsafe"
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
	// are ignored because a NativeWindow does not create a WKWebView.
	Mac MacWindow
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

// Run creates the native window. Applications normally let App.Run invoke it.
func (w *NativeWindow) Run() {
	w.lock.Lock()
	if w.destroyed || w.impl != nil {
		w.lock.Unlock()
		return
	}
	w.impl = newNativeWindowImpl(w)
	impl := w.impl
	w.lock.Unlock()
	if err := impl.run(); err != nil {
		w.Error("NativeWindow.Run: %s", err)
	}
}

func (w *NativeWindow) Show() *NativeWindow {
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

// SetSplitView installs a native AppKit split layout. The layout may use the
// existing WebView primary pane when attached to WebviewWindow, or a native
// text editor primary pane when attached to NativeWindow.
func (w *NativeWindow) SetSplitView(split *MacSplitView) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.impl != nil {
		return fmt.Errorf("split view must be configured before the native window is created")
	}
	if split == nil {
		previous := w.split
		w.split = nil
		releaseMacSplitViewOwnership(previous, w)
		return nil
	}
	if err := validateMacSplitView(split); err != nil {
		return err
	}
	if err := claimMacSplitView(split, w); err != nil {
		return err
	}
	previous := w.split
	w.split = split
	if previous != nil && previous != split {
		releaseMacSplitViewOwnership(previous, w)
	}
	return nil
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
