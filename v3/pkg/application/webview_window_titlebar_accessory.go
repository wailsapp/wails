package application

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

// MacTitlebarAccessoryPosition controls which side of the titlebar an
// accessory is placed on.
type MacTitlebarAccessoryPosition int

const (
	AccessoryLeading MacTitlebarAccessoryPosition = iota
	AccessoryTrailing
)

// MacTitlebarAccessoryOptions configures a titlebar accessory added via
// WebviewWindow.AddTitlebarAccessory.
type MacTitlebarAccessoryOptions struct {
	Name     string
	URL      string
	Position MacTitlebarAccessoryPosition
	// Width and Height default to 180x28 if left at zero.
	Width  float64
	Height float64
}

var titlebarAccessoryIDCounter uint64

func nextTitlebarAccessoryID() uint {
	return uint(atomic.AddUint64(&titlebarAccessoryIDCounter, 1))
}

// MacTitlebarAccessory is a handle to a titlebar accessory's own WKWebView,
// added via WebviewWindow.AddTitlebarAccessory.
type MacTitlebarAccessory struct {
	id   uint
	name string

	handleLock sync.RWMutex
	// nativeController is the NSTitlebarAccessoryViewController*; nativeWebview
	// is the WKWebView* it hosts. Both are nil until the owning window's
	// Run() has actually executed (may be deferred until app.Run() starts),
	// and stay nil forever on non-macOS builds.
	nativeController unsafe.Pointer
	nativeWebview     unsafe.Pointer
}

func (a *MacTitlebarAccessory) ID() uint    { return a.id }
func (a *MacTitlebarAccessory) Name() string { return a.name }

func (a *MacTitlebarAccessory) webview() (unsafe.Pointer, bool) {
	a.handleLock.RLock()
	defer a.handleLock.RUnlock()
	return a.nativeWebview, a.nativeWebview != nil
}

// SetURL navigates the accessory's webview. A no-op, logged, if the
// accessory isn't natively constructed yet (the owning window's Run()
// hasn't executed, e.g. AddTitlebarAccessory was called before app.Run()
// and the window hasn't been shown/run since).
func (a *MacTitlebarAccessory) SetURL(url string) {
	webview, ok := a.webview()
	if !ok {
		globalApplication.warning("MacTitlebarAccessory %q: SetURL called before the window was ready", a.name)
		return
	}
	InvokeSync(func() {
		macTitlebarAccessoryWebviewSetURL(webview, url)
	})
}

func (a *MacTitlebarAccessory) ExecJS(js string) {
	webview, ok := a.webview()
	if !ok {
		globalApplication.warning("MacTitlebarAccessory %q: ExecJS called before the window was ready", a.name)
		return
	}
	InvokeAsync(func() {
		macTitlebarAccessoryWebviewExecJS(webview, js)
	})
}

// macPendingTitlebarAccessory is stashed on WebviewWindow between AddTitlebarAccessory
// and the point where the window's Run() actually executes; darwin's run()
// consumes the queue via installPendingTitlebarAccessories. Never populated
// on any other platform.
type macPendingTitlebarAccessory struct {
	options   MacTitlebarAccessoryOptions
	accessory *MacTitlebarAccessory
}

// AddTitlebarAccessory adds a titlebar accessory hosting its own independent
// webview (macOS only; returns an error on every other platform). Safe to
// call before app.Run() -- construction is deferred until the window's
// native run() executes, same as SetToolbar.
func (w *WebviewWindow) AddTitlebarAccessory(options MacTitlebarAccessoryOptions) (*MacTitlebarAccessory, error) {
	if options.Name == "" {
		return nil, fmt.Errorf("MacTitlebarAccessoryOptions.Name is required")
	}
	if options.Width <= 0 {
		options.Width = 180
	}
	if options.Height <= 0 {
		options.Height = 28
	}

	accessory := &MacTitlebarAccessory{
		id:   nextTitlebarAccessoryID(),
		name: options.Name,
	}

	if w.impl == nil {
		w.titlebarAccessoriesLock.Lock()
		w.titlebarAccessoriesPending = append(w.titlebarAccessoriesPending, &macPendingTitlebarAccessory{
			options:   options,
			accessory: accessory,
		})
		w.titlebarAccessoriesLock.Unlock()
		return accessory, nil
	}

	var err error
	InvokeSync(func() {
		err = addTitlebarAccessoryNow(w, options, accessory)
	})
	if err != nil {
		return nil, err
	}
	return accessory, nil
}

// RemoveTitlebarAccessory removes a previously added titlebar accessory.
func (w *WebviewWindow) RemoveTitlebarAccessory(accessory *MacTitlebarAccessory) error {
	if accessory == nil {
		return nil
	}
	controller, ok := func() (unsafe.Pointer, bool) {
		accessory.handleLock.RLock()
		defer accessory.handleLock.RUnlock()
		return accessory.nativeController, accessory.nativeController != nil
	}()
	if !ok {
		return fmt.Errorf("MacTitlebarAccessory %q is not natively constructed yet", accessory.name)
	}
	if w.impl == nil {
		return fmt.Errorf("window is not natively constructed yet")
	}
	InvokeSync(func() {
		removeTitlebarAccessoryNow(w, controller)
	})
	return nil
}
