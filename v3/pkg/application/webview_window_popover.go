package application

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

// This file is the cross-platform surface of macOS popovers (NSPopover):
// anchored, transient panels shown from a rectangle in a window, from a
// toolbar item or from the system tray's status item. The platform
// functions it calls live in webview_window_popover_darwin.go, with
// documented stubs in webview_window_popover_other.go.
//
// A popover hosts native AppKit content: a MacAccessory control strip
// (search field, segmented control, buttons, labels, menu buttons, or any
// NSView through AddNativeView). Hosting a WKWebView page in a popover is a
// follow-up: the window's WebView cannot be reparented safely and a second
// WebView needs its own runtime bridge.

var (
	// ErrMacPopoverRequired means a nil popover was passed.
	ErrMacPopoverRequired = errors.New("a popover is required")
	// ErrMacPopoverDestroyed means the popover was destroyed and cannot be
	// shown again.
	ErrMacPopoverDestroyed = errors.New("the popover has been destroyed")
	// ErrMacPopoverAnchorRequired means the window, toolbar item or tray to
	// anchor the popover to is nil or has no native counterpart yet.
	ErrMacPopoverAnchorRequired = errors.New("a created window, installed toolbar item or running tray is required to anchor a popover")
	// ErrMacPopoverAnchorUnavailable means the toolbar item has no view to
	// anchor to on this macOS release (anchoring to items without a custom
	// view needs macOS 14).
	ErrMacPopoverAnchorUnavailable = errors.New("the toolbar item cannot anchor a popover on this macOS release")
	// ErrMacPopoverUnsupported means popovers are unavailable on the
	// current platform.
	ErrMacPopoverUnsupported = errors.New("popovers are unavailable on this platform")
	// ErrMacPopoverContentAttached means the accessory passed as content is
	// already attached elsewhere.
	ErrMacPopoverContentAttached = errors.New("the popover content accessory is already attached")
)

// MacPopoverBehavior mirrors NSPopoverBehavior: how the popover closes.
type MacPopoverBehavior int

const (
	// MacPopoverBehaviorApplicationDefined keeps the popover open until
	// Close is called. This is NSPopover's default and the zero value.
	MacPopoverBehaviorApplicationDefined MacPopoverBehavior = 0
	// MacPopoverBehaviorTransient closes the popover when the user
	// interacts with anything outside it.
	MacPopoverBehaviorTransient MacPopoverBehavior = 1
	// MacPopoverBehaviorSemitransient closes the popover when the user
	// interacts with the window that shows it, but not with other windows.
	MacPopoverBehaviorSemitransient MacPopoverBehavior = 2
)

func validMacPopoverBehavior(behavior MacPopoverBehavior) bool {
	switch behavior {
	case MacPopoverBehaviorApplicationDefined, MacPopoverBehaviorTransient, MacPopoverBehaviorSemitransient:
		return true
	}
	return false
}

// MacRectEdge names the edge of the anchor rectangle a popover attaches
// to, and so the side on which it appears. The values map to NSRectEdge
// when passed to AppKit; the zero value means below the anchor.
type MacRectEdge int

const (
	// MacRectEdgeDefault attaches to the bottom edge: the popover appears
	// below the anchor, as it does for a toolbar item or status item.
	MacRectEdgeDefault MacRectEdge = 0
	// MacRectEdgeMinX attaches to the left edge: the popover appears to the
	// left of the anchor.
	MacRectEdgeMinX MacRectEdge = 1
	// MacRectEdgeMinY attaches to the bottom edge: the popover appears
	// below the anchor.
	MacRectEdgeMinY MacRectEdge = 2
	// MacRectEdgeMaxX attaches to the right edge: the popover appears to
	// the right of the anchor.
	MacRectEdgeMaxX MacRectEdge = 3
	// MacRectEdgeMaxY attaches to the top edge: the popover appears above
	// the anchor.
	MacRectEdgeMaxY MacRectEdge = 4
)

func validMacRectEdge(edge MacRectEdge) bool {
	return edge >= MacRectEdgeDefault && edge <= MacRectEdgeMaxY
}

// nsRectEdge converts to the NSRectEdge value (NSMinXEdge 0, NSMinYEdge 1,
// NSMaxXEdge 2, NSMaxYEdge 3).
func (e MacRectEdge) nsRectEdge() int {
	switch e {
	case MacRectEdgeMinX:
		return 0
	case MacRectEdgeMaxX:
		return 2
	case MacRectEdgeMaxY:
		return 3
	default:
		return 1
	}
}

// Default popover content width in points when the options carry none.
const macPopoverDefaultWidth = 320.0

// Vertical padding added around a content strip when the options carry no
// height.
const macPopoverContentPadding = 8.0

// MacPopoverOptions configures NewMacPopover.
type MacPopoverOptions struct {
	// Width and Height are the content size in points. Zero width uses 320;
	// zero height fits the content strip.
	Width  float64
	Height float64
	// Behavior decides when the popover closes. The zero value keeps it
	// open until Close is called; most popovers want Transient.
	Behavior MacPopoverBehavior
	// DisableAnimation shows and hides the popover without the standard
	// animation.
	DisableAnimation bool
	// PreferredEdge is the edge of the anchor the popover attaches to when
	// a show method is not given one. The zero value is below the anchor.
	PreferredEdge MacRectEdge
	// Content is the native control strip shown inside the popover. It is
	// claimed when the popover is first shown and released by Destroy; add
	// every control before the first show. Control setters apply live;
	// SetHidden and SetHeight do not affect popover content (close the
	// popover or use SetContentSize instead). Nil shows an empty popover.
	Content *MacAccessory
}

// normalise fills defaults and reports the first invalid value.
func (o MacPopoverOptions) normalise() (MacPopoverOptions, error) {
	if !validMacPopoverBehavior(o.Behavior) {
		return o, fmt.Errorf("unknown macOS popover behavior %d", o.Behavior)
	}
	if !validMacRectEdge(o.PreferredEdge) {
		return o, fmt.Errorf("unknown macOS rect edge %d", o.PreferredEdge)
	}
	if o.Width < 0 || o.Height < 0 {
		return o, errors.New("popover size cannot be negative")
	}
	if o.Width == 0 {
		o.Width = macPopoverDefaultWidth
	}
	if o.Height == 0 {
		if o.Content != nil {
			o.Height = o.Content.resolvedHeight() + 2*macPopoverContentPadding
		} else {
			o.Height = MacAccessoryDefaultTitlebarHeight + 2*macPopoverContentPadding
		}
	}
	return o, nil
}

// macPopoverCallback is one OnClose registration.
type macPopoverCallback struct {
	callback func()
}

// MacPopover is a native NSPopover. Create it with NewMacPopover, show it
// with ShowRelativeTo, MacToolbarItem.ShowPopover or SystemTray.ShowPopover,
// and release it with Destroy. The native popover is created on the
// application thread the first time it is shown; every method is safe to
// call from any goroutine.
type MacPopover struct {
	id      uint64
	lock    sync.RWMutex
	options MacPopoverOptions
	// native is the NSPopover; it is created and released on the
	// application thread under lock.
	native    unsafe.Pointer
	destroyed bool

	closeCallbacks []*macPopoverCallback
}

var (
	macPopoverNextID   atomic.Uint64
	macPopoverRegistry = struct {
		sync.RWMutex
		popovers map[uint64]*MacPopover
	}{popovers: map[uint64]*MacPopover{}}
	macPopoverClosed = make(chan uint64, 16)
)

func init() {
	registerChromeEventLoop(func(*App) {
		for {
			id := <-macPopoverClosed
			go handleMacPopoverClosed(id)
		}
	})
}

func handleMacPopoverClosed(id uint64) {
	defer handlePanic()
	macPopoverRegistry.RLock()
	popover := macPopoverRegistry.popovers[id]
	macPopoverRegistry.RUnlock()
	if popover == nil {
		return
	}
	popover.lock.RLock()
	callbacks := append([]*macPopoverCallback(nil), popover.closeCallbacks...)
	popover.lock.RUnlock()
	for _, entry := range callbacks {
		entry.callback()
	}
}

// NewMacPopover creates a popover. Invalid options are reported through the
// application error handler and replaced by defaults, so the returned
// popover is always usable.
func NewMacPopover(options MacPopoverOptions) *MacPopover {
	normalised, err := options.normalise()
	if err != nil {
		if globalApplication != nil {
			globalApplication.error("NewMacPopover: %v", err)
		}
		options.Behavior = MacPopoverBehaviorApplicationDefined
		options.PreferredEdge = MacRectEdgeDefault
		if options.Width < 0 {
			options.Width = 0
		}
		if options.Height < 0 {
			options.Height = 0
		}
		normalised, _ = options.normalise()
	}
	popover := &MacPopover{id: macPopoverNextID.Add(1), options: normalised}
	macPopoverRegistry.Lock()
	macPopoverRegistry.popovers[popover.id] = popover
	macPopoverRegistry.Unlock()
	return popover
}

// ID identifies the popover; it also names it in error messages.
func (p *MacPopover) ID() uint {
	if p == nil {
		return 0
	}
	return uint(p.id)
}

// Error reports a popover problem through the application error handler.
// It satisfies the accessory host interface so the content accessory can
// log through its popover.
func (p *MacPopover) Error(message string, args ...any) {
	if globalApplication != nil {
		globalApplication.error("MacPopover %d: %s", p.ID(), fmt.Sprintf(message, args...))
	}
}

// Options returns the normalised options the popover was created with.
func (p *MacPopover) Options() MacPopoverOptions {
	if p == nil {
		return MacPopoverOptions{}
	}
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.options
}

// Content returns the accessory shown inside the popover, or nil.
func (p *MacPopover) Content() *MacAccessory {
	if p == nil {
		return nil
	}
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.options.Content
}

// OnClose registers a callback that runs after the popover closes, whether
// through Close, a click outside a transient popover, or Escape. It returns
// a function that removes the callback.
func (p *MacPopover) OnClose(callback func()) func() {
	if p == nil || callback == nil {
		return func() {}
	}
	entry := &macPopoverCallback{callback: callback}
	p.lock.Lock()
	p.closeCallbacks = append(p.closeCallbacks, entry)
	p.lock.Unlock()
	return func() {
		p.lock.Lock()
		defer p.lock.Unlock()
		remaining := p.closeCallbacks[:0]
		for _, existing := range p.closeCallbacks {
			if existing != entry {
				remaining = append(remaining, existing)
			}
		}
		p.closeCallbacks = remaining
	}
}

// IsDestroyed reports whether Destroy has been called.
func (p *MacPopover) IsDestroyed() bool {
	if p == nil {
		return true
	}
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.destroyed
}

// check reports why the popover cannot be shown, before any native work.
func (p *MacPopover) check() error {
	if p == nil {
		return ErrMacPopoverRequired
	}
	if !macPopoversSupported() {
		return ErrMacPopoverUnsupported
	}
	if p.IsDestroyed() {
		return ErrMacPopoverDestroyed
	}
	return nil
}

// ShowRelativeTo shows the popover anchored to rect, given in window
// content coordinates with the origin at the top left, attached to the
// given edge of that rectangle. The window's native window must exist. Pass
// MacRectEdgeDefault to use the options' PreferredEdge.
func (p *MacPopover) ShowRelativeTo(rect Rect, in Window, edge MacRectEdge) error {
	if err := p.check(); err != nil {
		return err
	}
	if in == nil {
		return ErrMacPopoverAnchorRequired
	}
	nsWindow := in.NativeWindow()
	if nsWindow == nil {
		return ErrMacPopoverAnchorRequired
	}
	native, err := macPopoverEnsureNative(p)
	if err != nil {
		return err
	}
	return macPopoverShowInWindow(native, nsWindow, rect, p.resolveEdge(edge))
}

// ShowRelativeToNative shows the popover anchored to rect in a NativeWindow.
// See ShowRelativeTo.
func (p *MacPopover) ShowRelativeToNative(rect Rect, in *NativeWindow, edge MacRectEdge) error {
	if err := p.check(); err != nil {
		return err
	}
	if in == nil {
		return ErrMacPopoverAnchorRequired
	}
	nsWindow := in.NativeWindow()
	if nsWindow == nil {
		return ErrMacPopoverAnchorRequired
	}
	native, err := macPopoverEnsureNative(p)
	if err != nil {
		return err
	}
	return macPopoverShowInWindow(native, nsWindow, rect, p.resolveEdge(edge))
}

func (p *MacPopover) resolveEdge(edge MacRectEdge) MacRectEdge {
	if !validMacRectEdge(edge) || edge == MacRectEdgeDefault {
		p.lock.RLock()
		edge = p.options.PreferredEdge
		p.lock.RUnlock()
	}
	if edge == MacRectEdgeDefault {
		edge = MacRectEdgeMinY
	}
	return edge
}

// Close closes the popover if it is shown. OnClose callbacks run
// afterwards.
func (p *MacPopover) Close() {
	if p == nil {
		return
	}
	p.lock.RLock()
	native := p.native
	p.lock.RUnlock()
	if native != nil {
		macPopoverClose(native)
	}
}

// IsShown reports whether the popover is on screen.
func (p *MacPopover) IsShown() bool {
	if p == nil {
		return false
	}
	p.lock.RLock()
	native := p.native
	p.lock.RUnlock()
	return native != nil && macPopoverIsShown(native)
}

// SetContentSize changes the content size in points; a shown popover
// resizes in place. Non-positive values are ignored.
func (p *MacPopover) SetContentSize(width, height float64) *MacPopover {
	if p == nil || width <= 0 || height <= 0 {
		return p
	}
	p.lock.Lock()
	p.options.Width, p.options.Height = width, height
	native := p.native
	p.lock.Unlock()
	if native != nil {
		macPopoverSetContentSize(native, width, height)
	}
	return p
}

// SetBehavior changes how the popover closes. Unknown values are rejected.
func (p *MacPopover) SetBehavior(behavior MacPopoverBehavior) error {
	if p == nil {
		return ErrMacPopoverRequired
	}
	if !validMacPopoverBehavior(behavior) {
		return fmt.Errorf("unknown macOS popover behavior %d", behavior)
	}
	p.lock.Lock()
	p.options.Behavior = behavior
	native := p.native
	p.lock.Unlock()
	if native != nil {
		macPopoverSetBehavior(native, behavior)
	}
	return nil
}

// Destroy closes the popover, detaches its content accessory and releases
// the native popover. The popover cannot be shown again; the accessory can
// be attached elsewhere.
func (p *MacPopover) Destroy() {
	if p == nil {
		return
	}
	p.lock.Lock()
	if p.destroyed {
		p.lock.Unlock()
		return
	}
	p.destroyed = true
	content := p.options.Content
	p.lock.Unlock()
	// The popover drops its reference to the content controller first, then
	// the accessory releases its own; see WailsPopoverContentController.
	macPopoverRelease(p)
	if content != nil {
		content.Remove()
	}
	macPopoverRegistry.Lock()
	delete(macPopoverRegistry.popovers, p.id)
	macPopoverRegistry.Unlock()
}

// macToolbarItemShowPopover anchors popover to a toolbar item; see
// MacToolbarItem.ShowPopover.
func macToolbarItemShowPopover(item *MacToolbarItem, popover *MacPopover) error {
	if err := popover.check(); err != nil {
		return err
	}
	if item == nil || item.toolbar == nil {
		return ErrMacPopoverAnchorRequired
	}
	native, err := macPopoverEnsureNative(popover)
	if err != nil {
		return err
	}
	result := ErrMacPopoverAnchorRequired
	item.update(func(unsafe.Pointer) {
		var nsWindow unsafe.Pointer
		if host, ok := item.toolbar.state.window.(interface{ NativeWindow() unsafe.Pointer }); ok && host != nil {
			nsWindow = host.NativeWindow()
		}
		if nsWindow == nil {
			return
		}
		result = macPopoverShowFromToolbarItem(native, nsWindow, item.identifier, popover.resolveEdge(MacRectEdgeDefault))
	})
	return result
}

// macSystemTrayShowPopover anchors popover to the tray's status item; see
// SystemTray.ShowPopover.
func macSystemTrayShowPopover(tray *SystemTray, popover *MacPopover) error {
	if err := popover.check(); err != nil {
		return err
	}
	if tray == nil || tray.impl == nil {
		return ErrMacPopoverAnchorRequired
	}
	statusItem := macSystemTrayStatusItem(tray)
	if statusItem == nil {
		return ErrMacPopoverAnchorRequired
	}
	native, err := macPopoverEnsureNative(popover)
	if err != nil {
		return err
	}
	return macPopoverShowFromStatusItem(native, statusItem, popover.resolveEdge(MacRectEdgeDefault))
}
