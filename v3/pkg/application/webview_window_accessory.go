package application

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

// MacScrollEdgeEffectStyle is AppKit's preferred scroll-edge treatment for a
// titlebar or split-view-item accessory controller. AppKit controls the exact
// blur radius, fade depth, and adaptation; these values select only its
// documented automatic, soft, and hard policies.
type MacScrollEdgeEffectStyle int

const (
	MacScrollEdgeEffectStyleAutomatic MacScrollEdgeEffectStyle = iota
	MacScrollEdgeEffectStyleSoft
	MacScrollEdgeEffectStyleHard
)

func validMacScrollEdgeEffectStyle(style MacScrollEdgeEffectStyle) bool {
	return style >= MacScrollEdgeEffectStyleAutomatic && style <= MacScrollEdgeEffectStyleHard
}

var (
	// ErrMacAccessoryControllerRequired means a nil native controller was
	// supplied to WrapMacAccessoryViewController.
	ErrMacAccessoryControllerRequired = errors.New("a native macOS accessory view controller is required")
	// ErrMacAccessoryControllerType means the native object is neither an
	// NSTitlebarAccessoryViewController nor an
	// NSSplitViewItemAccessoryViewController.
	ErrMacAccessoryControllerType = errors.New("native object is not a supported macOS accessory view controller")
	// ErrMacScrollEdgeEffectStyleUnavailable means the requested explicit style
	// requires macOS 26.1 or newer. Automatic remains a valid no-op fallback.
	ErrMacScrollEdgeEffectStyleUnavailable = errors.New("preferred scroll-edge effect styles require macOS 26.1 or newer")
	// ErrMacAccessoryControllerUnsupported means native AppKit accessory
	// controllers are unavailable on the current platform.
	ErrMacAccessoryControllerUnsupported = errors.New("macOS accessory view controllers are unavailable on this platform")
)

// MacAccessoryViewController is a non-owning, type-checked wrapper around an
// NSTitlebarAccessoryViewController or
// NSSplitViewItemAccessoryViewController. It exists for native integrations
// that construct or receive an AppKit accessory controller and need to use
// the shared scroll-edge style API without private selectors.
//
// The native owner (normally NSWindow or NSSplitViewItem) must keep the
// controller alive for as long as this wrapper is used.
type MacAccessoryViewController struct {
	native unsafe.Pointer
	kind   MacAccessoryViewControllerKind
}

// MacAccessoryViewControllerKind identifies the AppKit controller class held
// by a MacAccessoryViewController.
type MacAccessoryViewControllerKind uint8

const (
	MacAccessoryViewControllerKindUnknown MacAccessoryViewControllerKind = iota
	MacAccessoryViewControllerKindTitlebar
	MacAccessoryViewControllerKindSplitItem
)

// WrapMacAccessoryViewController validates and wraps an existing native
// NSTitlebarAccessoryViewController or
// NSSplitViewItemAccessoryViewController. The wrapper does not retain or
// release the native object.
func WrapMacAccessoryViewController(native unsafe.Pointer) (*MacAccessoryViewController, error) {
	if native == nil {
		return nil, ErrMacAccessoryControllerRequired
	}
	kind, err := macAccessoryControllerKind(native)
	if err != nil {
		return nil, err
	}
	if kind != MacAccessoryViewControllerKindTitlebar && kind != MacAccessoryViewControllerKindSplitItem {
		return nil, ErrMacAccessoryControllerType
	}
	return &MacAccessoryViewController{native: native, kind: kind}, nil
}

// Kind reports whether the wrapped native object is a titlebar or split-item
// accessory view controller.
func (c *MacAccessoryViewController) Kind() MacAccessoryViewControllerKind {
	if c == nil {
		return MacAccessoryViewControllerKindUnknown
	}
	return c.kind
}

// NativeController returns the wrapped AppKit controller pointer. The pointer
// is borrowed and has the same lifetime as its native owner.
func (c *MacAccessoryViewController) NativeController() unsafe.Pointer {
	if c == nil {
		return nil
	}
	return c.native
}

// SupportsPreferredScrollEdgeEffectStyle reports whether this controller can
// currently apply Apple's preferredScrollEdgeEffectStyle property. The API is
// available for both supported accessory controller classes on macOS 26.1+.
func (c *MacAccessoryViewController) SupportsPreferredScrollEdgeEffectStyle() bool {
	return c != nil && c.native != nil && macAccessoryControllerSupportsScrollEdgeEffectStyle(c.native)
}

// SetPreferredScrollEdgeEffectStyle sets AppKit's preferred effect for content
// scrolling behind this accessory. Automatic, Soft, and Hard map directly to
// NSScrollEdgeEffectStyle. The system still owns the precise fade and blur.
//
// On macOS versions before 26.1, Automatic succeeds as a no-op because it is
// the platform default; explicit Soft or Hard returns
// ErrMacScrollEdgeEffectStyleUnavailable.
func (c *MacAccessoryViewController) SetPreferredScrollEdgeEffectStyle(style MacScrollEdgeEffectStyle) error {
	if c == nil || c.native == nil {
		return ErrMacAccessoryControllerRequired
	}
	if !validMacScrollEdgeEffectStyle(style) {
		return fmt.Errorf("unknown macOS scroll-edge effect style %d", style)
	}
	return macAccessoryControllerSetScrollEdgeEffectStyle(c.native, style)
}

// PreferredScrollEdgeEffectStyle returns the controller's current AppKit
// preference. Before macOS 26.1 it returns Automatic, the system behavior.
func (c *MacAccessoryViewController) PreferredScrollEdgeEffectStyle() (MacScrollEdgeEffectStyle, error) {
	if c == nil || c.native == nil {
		return MacScrollEdgeEffectStyleAutomatic, ErrMacAccessoryControllerRequired
	}
	return macAccessoryControllerScrollEdgeEffectStyle(c.native)
}

// MacAccessoryLayout selects where a MacAccessory is placed. Leading,
// Trailing, and Bottom are titlebar placements (NSLayoutAttributeLeading,
// NSLayoutAttributeTrailing, and NSLayoutAttributeBottom). Top and Bottom are
// the alignments accepted by split-view panes through
// MacSplitPane.AddTopAccessory and MacSplitPane.AddBottomAccessory.
type MacAccessoryLayout int

const (
	// MacAccessoryLayoutLeading places a titlebar accessory next to the
	// window buttons.
	MacAccessoryLayoutLeading MacAccessoryLayout = iota
	// MacAccessoryLayoutTrailing places a titlebar accessory at the trailing
	// edge of the titlebar or toolbar.
	MacAccessoryLayoutTrailing
	// MacAccessoryLayoutBottom is a full-width strip below the titlebar and
	// toolbar, or the bottom-aligned accessory of a split-view pane.
	MacAccessoryLayoutBottom
	// MacAccessoryLayoutTop is the top-aligned accessory of a split-view
	// pane, the place for a Finder-style filter strip above a sidebar.
	MacAccessoryLayoutTop
)

func (l MacAccessoryLayout) String() string {
	switch l {
	case MacAccessoryLayoutLeading:
		return "leading"
	case MacAccessoryLayoutTrailing:
		return "trailing"
	case MacAccessoryLayoutBottom:
		return "bottom"
	case MacAccessoryLayoutTop:
		return "top"
	default:
		return fmt.Sprintf("MacAccessoryLayout(%d)", int(l))
	}
}

func validMacAccessoryLayout(l MacAccessoryLayout) bool {
	return l >= MacAccessoryLayoutLeading && l <= MacAccessoryLayoutTop
}

// titlebar reports whether the layout is valid for a window titlebar.
func (l MacAccessoryLayout) titlebar() bool {
	return l == MacAccessoryLayoutLeading || l == MacAccessoryLayoutTrailing || l == MacAccessoryLayoutBottom
}

// splitItem reports whether the layout is valid for a split-view pane.
func (l MacAccessoryLayout) splitItem() bool {
	return l == MacAccessoryLayoutTop || l == MacAccessoryLayoutBottom
}

const (
	// MacAccessoryDefaultTitlebarHeight is the height used by a titlebar
	// accessory when SetHeight is not called.
	MacAccessoryDefaultTitlebarHeight = 28.0
	// MacAccessoryDefaultSplitItemHeight is the height used by a split-view
	// pane accessory when SetHeight is not called.
	MacAccessoryDefaultSplitItemHeight = 36.0
	// macAccessoryDefaultSearchWidth is the width of a search field whose
	// width is left at zero (fit), because NSSearchField has no useful
	// intrinsic width.
	macAccessoryDefaultSearchWidth = 180.0
)

var (
	// ErrMacAccessoryRequired means a nil accessory was passed to an attach
	// method.
	ErrMacAccessoryRequired = errors.New("a macOS accessory is required")
	// ErrMacAccessoryLayoutInvalid means the accessory's layout does not fit
	// the attach point: Top cannot be used in a titlebar, and pane accessories
	// must use Top for AddTopAccessory or Bottom for AddBottomAccessory.
	ErrMacAccessoryLayoutInvalid = errors.New("macOS accessory layout is not valid for this attach point")
	// ErrMacAccessoryAttached means the accessory is already attached (or
	// queued for attachment) elsewhere. Call Remove first.
	ErrMacAccessoryAttached = errors.New("macOS accessory is already attached")
	// ErrMacAccessoryUnsupported means native accessories are unavailable on
	// the current platform. Constructors and setters still work so shared
	// code can build accessories unconditionally.
	ErrMacAccessoryUnsupported = errors.New("native macOS accessories are unavailable on this platform")
	// ErrMacSplitItemAccessoryUnavailable means AppKit's split-view item
	// accessories (NSSplitViewItemAccessoryViewController) require macOS 26
	// or newer. Titlebar accessories remain available on every supported
	// macOS release.
	ErrMacSplitItemAccessoryUnavailable = errors.New("split-view pane accessories require macOS 26 or newer")
	// ErrMacAccessoryPaneRequired means the split pane does not belong to a
	// split view or has already been torn down.
	ErrMacAccessoryPaneRequired = errors.New("the split pane is not part of a live split view")
)

// macAccessoryWindow is the window surface an accessory needs from its
// titlebar host. Both WebviewWindow and NativeWindow satisfy it.
type macAccessoryWindow interface {
	ID() uint
	Error(message string, args ...any)
}

// macAccessoryTarget records where an accessory is attached or queued.
type macAccessoryTarget struct {
	window macAccessoryWindow
	pane   *MacSplitPane
	// top is the alignment of a pane accessory.
	top bool
}

func (t macAccessoryTarget) empty() bool {
	return t.window == nil && t.pane == nil
}

// key is the pending-queue key: the host window or pane.
func (t macAccessoryTarget) key() any {
	if t.pane != nil {
		return t.pane
	}
	if t.window != nil {
		return t.window
	}
	return nil
}

// MacAccessory is a strip of native AppKit controls attached to a window's
// titlebar (NSTitlebarAccessoryViewController) or to the top or bottom of a
// split-view pane (NSSplitViewItemAccessoryViewController, macOS 26+). Build
// it with NewMacAccessory, add controls with the Add methods, then attach it
// with WebviewWindow.AddTitlebarAccessory, NativeWindow.AddTitlebarAccessory,
// MacSplitPane.AddTopAccessory, or MacSplitPane.AddBottomAccessory.
//
// Controls are laid out in insertion order in a horizontal NSStackView with
// standard spacing and 8pt horizontal insets. An accessory can be attached
// before or after its window exists; attachment before creation is queued
// and applied when the native window is built. Every setter is safe to call
// at any time: the value is stored and, once attached, applied on the
// application thread.
//
// An accessory belongs to one attach point at a time. Remove detaches it
// natively so it can be attached elsewhere.
type MacAccessory struct {
	lock sync.RWMutex

	layout   MacAccessoryLayout
	controls []*MacAccessoryControl

	hidden              bool
	height              float64
	fullScreenMinHeight float64
	adjustsSize         bool
	adjustsSizeSet      bool
	contentInsets       bool
	scrollEdgeStyle     MacScrollEdgeEffectStyle
	scrollEdgeStyleSet  bool

	// target is set while the accessory is queued or attached. native and
	// controller are set only while the native controller exists; both are
	// mutated on the application thread.
	target     macAccessoryTarget
	native     unsafe.Pointer
	controller *MacAccessoryViewController
}

// NewMacAccessory creates an empty accessory for the given placement. The
// layout decides which attach points accept it; see MacAccessoryLayout.
func NewMacAccessory(layout MacAccessoryLayout) *MacAccessory {
	if !validMacAccessoryLayout(layout) {
		layout = MacAccessoryLayoutBottom
	}
	return &MacAccessory{layout: layout, contentInsets: true}
}

// Layout returns the placement chosen at construction.
func (a *MacAccessory) Layout() MacAccessoryLayout {
	if a == nil {
		return MacAccessoryLayoutBottom
	}
	return a.layout
}

// Controller returns the type-checked wrapper around the native accessory
// view controller, or nil while the accessory is not attached. It gives
// native integrations the AppKit controller without unsafe pointer juggling.
func (a *MacAccessory) Controller() *MacAccessoryViewController {
	if a == nil {
		return nil
	}
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.controller
}

// IsAttached reports whether the native controller currently exists.
func (a *MacAccessory) IsAttached() bool {
	if a == nil {
		return false
	}
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.native != nil
}

// Controls returns the controls in layout order.
func (a *MacAccessory) Controls() []*MacAccessoryControl {
	if a == nil {
		return nil
	}
	a.lock.RLock()
	defer a.lock.RUnlock()
	return append([]*MacAccessoryControl(nil), a.controls...)
}

// SetHidden collapses the accessory to zero height without detaching it.
// AppKit animates the change for attached accessories.
func (a *MacAccessory) SetHidden(hidden bool) *MacAccessory {
	if a == nil {
		return a
	}
	a.lock.Lock()
	a.hidden = hidden
	a.lock.Unlock()
	a.update(func(native unsafe.Pointer) { macAccessorySetHidden(native, hidden) })
	return a
}

// IsHidden reports the last requested hidden state.
func (a *MacAccessory) IsHidden() bool {
	if a == nil {
		return false
	}
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.hidden
}

// SetHeight sets the accessory's height in points. Zero restores the default:
// MacAccessoryDefaultTitlebarHeight for titlebar accessories and
// MacAccessoryDefaultSplitItemHeight for pane accessories. Negative or
// non-finite values are ignored.
func (a *MacAccessory) SetHeight(points float64) *MacAccessory {
	if a == nil || !isFiniteMacSplitNumber(points) || points < 0 {
		return a
	}
	a.lock.Lock()
	a.height = points
	a.lock.Unlock()
	a.update(func(native unsafe.Pointer) { macAccessorySetHeight(native, a.resolvedHeight()) })
	return a
}

// Height returns the configured height, or zero when the default applies.
func (a *MacAccessory) Height() float64 {
	if a == nil {
		return 0
	}
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.height
}

// resolvedHeight is the height applied natively.
func (a *MacAccessory) resolvedHeight() float64 {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.resolvedHeightLocked()
}

func (a *MacAccessory) resolvedHeightLocked() float64 {
	if a.height > 0 {
		return a.height
	}
	if a.target.pane != nil || a.layout == MacAccessoryLayoutTop {
		return MacAccessoryDefaultSplitItemHeight
	}
	return MacAccessoryDefaultTitlebarHeight
}

// SetFullScreenMinHeight controls how much of a Bottom titlebar accessory
// stays visible in full screen while the menu bar is hidden (AppKit's
// fullScreenMinHeight). Zero hides it fully; the accessory height keeps it
// fully shown. It is ignored for other layouts and for pane accessories.
func (a *MacAccessory) SetFullScreenMinHeight(points float64) *MacAccessory {
	if a == nil || !isFiniteMacSplitNumber(points) || points < 0 {
		return a
	}
	a.lock.Lock()
	a.fullScreenMinHeight = points
	a.lock.Unlock()
	a.update(func(native unsafe.Pointer) { macAccessorySetFullScreenMinHeight(native, points) })
	return a
}

// SetAutomaticallyAdjustsSize lets AppKit size a Bottom titlebar accessory
// from its view's constraints instead of its frame (macOS 11+,
// automaticallyAdjustsSize). It is ignored for other layouts and for pane
// accessories.
func (a *MacAccessory) SetAutomaticallyAdjustsSize(adjusts bool) *MacAccessory {
	if a == nil {
		return a
	}
	a.lock.Lock()
	a.adjustsSize = adjusts
	a.adjustsSizeSet = true
	a.lock.Unlock()
	a.update(func(native unsafe.Pointer) { macAccessorySetAutomaticallyAdjustsSize(native, adjusts) })
	return a
}

// SetAutomaticallyAppliesContentInsets controls whether a pane accessory
// receives AppKit's standard content insets
// (automaticallyAppliesContentInsets, default true). Titlebar accessories
// ignore it.
func (a *MacAccessory) SetAutomaticallyAppliesContentInsets(applies bool) *MacAccessory {
	if a == nil {
		return a
	}
	a.lock.Lock()
	a.contentInsets = applies
	a.lock.Unlock()
	a.update(func(native unsafe.Pointer) { macAccessorySetAppliesContentInsets(native, applies) })
	return a
}

// SetPreferredScrollEdgeEffectStyle sets AppKit's preferred scroll-edge
// treatment for content scrolling behind the accessory. It delegates to
// MacAccessoryViewController.SetPreferredScrollEdgeEffectStyle once attached;
// explicit Soft or Hard styles require macOS 26.1 and are reported through
// the owning window's error handler when unavailable.
func (a *MacAccessory) SetPreferredScrollEdgeEffectStyle(style MacScrollEdgeEffectStyle) *MacAccessory {
	if a == nil {
		return a
	}
	if !validMacScrollEdgeEffectStyle(style) {
		a.logError("SetPreferredScrollEdgeEffectStyle: unknown scroll-edge effect style %d", style)
		return a
	}
	a.lock.Lock()
	a.scrollEdgeStyle = style
	a.scrollEdgeStyleSet = true
	a.lock.Unlock()
	a.applyScrollEdgeStyle()
	return a
}

// PreferredScrollEdgeEffectStyle returns the requested style, or Automatic.
func (a *MacAccessory) PreferredScrollEdgeEffectStyle() MacScrollEdgeEffectStyle {
	if a == nil {
		return MacScrollEdgeEffectStyleAutomatic
	}
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.scrollEdgeStyle
}

// applyScrollEdgeStyle forwards the stored style to the attached controller.
func (a *MacAccessory) applyScrollEdgeStyle() {
	a.lock.RLock()
	controller := a.controller
	style := a.scrollEdgeStyle
	set := a.scrollEdgeStyleSet
	a.lock.RUnlock()
	if controller == nil || !set {
		return
	}
	if err := controller.SetPreferredScrollEdgeEffectStyle(style); err != nil {
		a.logError("SetPreferredScrollEdgeEffectStyle: %s", err)
	}
}

// update runs a native update on the application thread when the accessory
// is attached. It re-checks the native handle under the lock on that thread
// so a concurrent Remove cannot free the controller underneath the update.
func (a *MacAccessory) update(apply func(native unsafe.Pointer)) {
	if a == nil || globalApplication == nil {
		return
	}
	a.lock.RLock()
	attached := a.native != nil
	a.lock.RUnlock()
	if !attached {
		return
	}
	InvokeSync(func() {
		a.lock.RLock()
		defer a.lock.RUnlock()
		if a.native != nil {
			apply(a.native)
		}
	})
}

// logError reports through the host window when there is one, and through
// the application logger otherwise.
func (a *MacAccessory) logError(format string, args ...any) {
	if a == nil || globalApplication == nil {
		return
	}
	a.lock.RLock()
	target := a.target
	a.lock.RUnlock()
	if target.window != nil {
		target.window.Error(format, args...)
		return
	}
	if target.pane != nil && target.pane.split != nil {
		if owner := target.pane.split.ownerWindow(); owner != nil {
			owner.Error(format, args...)
			return
		}
	}
	globalApplication.error(format, args...)
}

// claim reserves the accessory for one attach point.
func (a *MacAccessory) claim(target macAccessoryTarget) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if !a.target.empty() {
		return ErrMacAccessoryAttached
	}
	a.target = target
	return nil
}

// releaseClaim undoes claim after a failed attachment.
func (a *MacAccessory) releaseClaim() {
	a.lock.Lock()
	if a.native == nil {
		a.target = macAccessoryTarget{}
	}
	a.lock.Unlock()
}

// Remove detaches the accessory from its window or pane, or drops it from the
// pending queue when the native window does not exist yet. Control callbacks
// stop firing immediately. The accessory keeps its controls and settings and
// can be attached again.
func (a *MacAccessory) Remove() {
	if a == nil {
		return
	}
	a.lock.RLock()
	target := a.target
	a.lock.RUnlock()
	if key := target.key(); key != nil {
		dequeueMacAccessory(key, a)
	}
	macAccessoryDetachNative(a)
	unregisterMacAccessoryControls(a)
	a.lock.Lock()
	a.target = macAccessoryTarget{}
	a.native = nil
	a.controller = nil
	a.lock.Unlock()
}

// add appends a control. Controls added after attachment are not installed
// natively; add every control before attaching the accessory.
func (a *MacAccessory) add(control *MacAccessoryControl) *MacAccessoryControl {
	a.lock.Lock()
	if a.native != nil {
		a.lock.Unlock()
		a.logError("accessory controls must be added before the accessory is attached")
		return control
	}
	a.controls = append(a.controls, control)
	a.lock.Unlock()
	return control
}

// AddSearch adds a native NSSearchField. OnSearch receives the query when the
// user presses Return or clears the field, or on every keystroke after
// SetIncremental(true). The default width is 180 points until SetWidth is
// called.
func (a *MacAccessory) AddSearch(placeholder string) *MacAccessoryControl {
	if a == nil {
		return nil
	}
	control := newMacAccessoryControl(a, MacAccessoryControlSearch)
	control.placeholder = placeholder
	return a.add(control)
}

// AddSegmented adds a native NSSegmentedControl with one segment per label.
// selected is the initially selected segment; a negative value selects none.
// Use SetSegmentSymbols to show SF Symbols instead of the labels.
func (a *MacAccessory) AddSegmented(labels []string, selected int) *MacAccessoryControl {
	if a == nil {
		return nil
	}
	control := newMacAccessoryControl(a, MacAccessoryControlSegmented)
	control.segments = append([]string(nil), labels...)
	if selected >= len(labels) {
		selected = -1
	}
	control.selected = selected
	return a.add(control)
}

// AddButton adds a native push button with a text label. Use SetSymbol to
// add an SF Symbol image.
func (a *MacAccessory) AddButton(label string) *MacAccessoryControl {
	if a == nil {
		return nil
	}
	control := newMacAccessoryControl(a, MacAccessoryControlButton)
	control.text = label
	return a.add(control)
}

// AddSymbolButton adds a native push button showing only an SF Symbol.
func (a *MacAccessory) AddSymbolButton(symbol string) *MacAccessoryControl {
	if a == nil {
		return nil
	}
	control := newMacAccessoryControl(a, MacAccessoryControlButton)
	control.symbol = symbol
	return a.add(control)
}

// AddMenuButton adds a native button that pops up menu when clicked. Menu
// item clicks fire the usual MenuItem.OnClick callbacks and Menu.Update
// refreshes the native menu in place. Use SetSymbol for an icon-only button.
func (a *MacAccessory) AddMenuButton(label string, menu *Menu) *MacAccessoryControl {
	if a == nil {
		return nil
	}
	control := newMacAccessoryControl(a, MacAccessoryControlMenuButton)
	control.text = label
	control.menu = menu
	return a.add(control)
}

// AddLabel adds a secondary-style status label. SetText updates it live and
// SetSymbol places an SF Symbol before the text.
func (a *MacAccessory) AddLabel(text string) *MacAccessoryControl {
	if a == nil {
		return nil
	}
	control := newMacAccessoryControl(a, MacAccessoryControlLabel)
	control.text = text
	return a.add(control)
}

// AddFlexibleSpace adds a spacer that absorbs any spare width, pushing later
// controls to the trailing edge. It only has an effect in layouts whose width
// is set by AppKit (Bottom titlebar accessories and pane accessories);
// Leading and Trailing titlebar accessories are sized to fit their controls.
func (a *MacAccessory) AddFlexibleSpace() *MacAccessoryControl {
	if a == nil {
		return nil
	}
	return a.add(newMacAccessoryControl(a, MacAccessoryControlFlexibleSpace))
}

// AddNativeView adds an existing AppKit view (an NSView pointer) to the
// strip. This is an unsafe escape hatch: the pointer is not type-checked,
// the view must outlive the accessory or be retained by its creator, it must
// not already have a superview, and it must be created and mutated only on
// the application thread. The view is sized by its own constraints or by
// SetWidth. It is ignored on platforms without AppKit.
func (a *MacAccessory) AddNativeView(view unsafe.Pointer) *MacAccessoryControl {
	if a == nil {
		return nil
	}
	control := newMacAccessoryControl(a, MacAccessoryControlNativeView)
	control.nativeView = view
	return a.add(control)
}

// MacAccessoryControlKind identifies the AppKit control behind a
// MacAccessoryControl.
type MacAccessoryControlKind int

const (
	MacAccessoryControlSearch MacAccessoryControlKind = iota
	MacAccessoryControlSegmented
	MacAccessoryControlButton
	MacAccessoryControlMenuButton
	MacAccessoryControlLabel
	MacAccessoryControlFlexibleSpace
	MacAccessoryControlNativeView
)

// MacAccessoryControl is one native control inside a MacAccessory. Setters
// may be called at any time; changes are applied live once the accessory is
// attached. Setters that do not apply to the control's kind are no-ops.
type MacAccessoryControl struct {
	lock sync.RWMutex

	accessory  *MacAccessory
	internalID uint64
	kind       MacAccessoryControlKind

	text           string
	placeholder    string
	symbol         string
	tooltip        string
	width          float64
	disabled       bool
	hidden         bool
	incremental    bool
	segments       []string
	segmentSymbols []string
	selected       int
	menu           *Menu
	nativeView     unsafe.Pointer

	onSearch          func(*Context, string)
	onSelectionChange func(*Context, int, string)
	onClick           func(*Context)
}

var macAccessoryControlID uint64

func nextMacAccessoryControlID() uint64 {
	return atomic.AddUint64(&macAccessoryControlID, 1)
}

func newMacAccessoryControl(accessory *MacAccessory, kind MacAccessoryControlKind) *MacAccessoryControl {
	return &MacAccessoryControl{
		accessory:  accessory,
		internalID: nextMacAccessoryControlID(),
		kind:       kind,
		selected:   -1,
	}
}

// Kind reports the control's native kind.
func (c *MacAccessoryControl) Kind() MacAccessoryControlKind {
	if c == nil {
		return MacAccessoryControlFlexibleSpace
	}
	return c.kind
}

// Accessory returns the owning accessory.
func (c *MacAccessoryControl) Accessory() *MacAccessory {
	if c == nil {
		return nil
	}
	return c.accessory
}

// OnSearch sets the callback for a search field. Passing nil clears it.
func (c *MacAccessoryControl) OnSearch(callback func(*Context, string)) *MacAccessoryControl {
	if c == nil {
		return c
	}
	c.lock.Lock()
	c.onSearch = callback
	c.lock.Unlock()
	return c
}

// OnSelectionChange sets the callback for a segmented control. It receives
// the selected index and its label. Passing nil clears it.
func (c *MacAccessoryControl) OnSelectionChange(callback func(*Context, int, string)) *MacAccessoryControl {
	if c == nil {
		return c
	}
	c.lock.Lock()
	c.onSelectionChange = callback
	c.lock.Unlock()
	return c
}

// OnClick sets the callback for a button. Menu buttons open their menu
// instead and ignore it. Passing nil clears it.
func (c *MacAccessoryControl) OnClick(callback func(*Context)) *MacAccessoryControl {
	if c == nil {
		return c
	}
	c.lock.Lock()
	c.onClick = callback
	c.lock.Unlock()
	return c
}

// SetEnabled enables or disables the control.
func (c *MacAccessoryControl) SetEnabled(enabled bool) *MacAccessoryControl {
	if c == nil {
		return c
	}
	c.lock.Lock()
	c.disabled = !enabled
	c.lock.Unlock()
	c.update(func(native unsafe.Pointer, id uint64) { macAccessoryControlSetEnabled(native, id, enabled) })
	return c
}

// SetHidden hides or shows the control. Hidden controls give up their space.
func (c *MacAccessoryControl) SetHidden(hidden bool) *MacAccessoryControl {
	if c == nil {
		return c
	}
	c.lock.Lock()
	c.hidden = hidden
	c.lock.Unlock()
	c.update(func(native unsafe.Pointer, id uint64) { macAccessoryControlSetHidden(native, id, hidden) })
	return c
}

// SetTooltip sets the control's tooltip. An empty string removes it.
func (c *MacAccessoryControl) SetTooltip(tooltip string) *MacAccessoryControl {
	if c == nil {
		return c
	}
	c.lock.Lock()
	c.tooltip = tooltip
	c.lock.Unlock()
	c.update(func(native unsafe.Pointer, id uint64) { macAccessoryControlSetTooltip(native, id, tooltip) })
	return c
}

// SetWidth fixes the control's width in points. Zero lets the control fit
// its content (search fields fall back to 180 points). Flexible spaces ignore
// it.
func (c *MacAccessoryControl) SetWidth(points float64) *MacAccessoryControl {
	if c == nil || !isFiniteMacSplitNumber(points) || points < 0 {
		return c
	}
	c.lock.Lock()
	c.width = points
	c.lock.Unlock()
	c.update(func(native unsafe.Pointer, id uint64) { macAccessoryControlSetWidth(native, id, points) })
	return c
}

// SetText updates a label's text, a button's title, or a search field's
// current query.
func (c *MacAccessoryControl) SetText(text string) *MacAccessoryControl {
	if c == nil {
		return c
	}
	c.lock.Lock()
	c.text = text
	c.lock.Unlock()
	c.update(func(native unsafe.Pointer, id uint64) { macAccessoryControlSetText(native, id, text) })
	return c
}

// Text returns the label text, button title, or last known search query.
func (c *MacAccessoryControl) Text() string {
	if c == nil {
		return ""
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.text
}

// SetSymbol sets the SF Symbol shown by a button, menu button, or label. An
// empty name removes the image.
func (c *MacAccessoryControl) SetSymbol(symbol string) *MacAccessoryControl {
	if c == nil {
		return c
	}
	c.lock.Lock()
	c.symbol = symbol
	c.lock.Unlock()
	c.update(func(native unsafe.Pointer, id uint64) { macAccessoryControlSetSymbol(native, id, symbol) })
	return c
}

// SetPlaceholder sets a search field's placeholder text.
func (c *MacAccessoryControl) SetPlaceholder(placeholder string) *MacAccessoryControl {
	if c == nil {
		return c
	}
	c.lock.Lock()
	c.placeholder = placeholder
	c.lock.Unlock()
	c.update(func(native unsafe.Pointer, id uint64) { macAccessoryControlSetPlaceholder(native, id, placeholder) })
	return c
}

// SetIncremental makes a search field report every keystroke through
// OnSearch instead of waiting for Return.
func (c *MacAccessoryControl) SetIncremental(incremental bool) *MacAccessoryControl {
	if c == nil {
		return c
	}
	c.lock.Lock()
	c.incremental = incremental
	c.lock.Unlock()
	c.update(func(native unsafe.Pointer, id uint64) { macAccessoryControlSetIncremental(native, id, incremental) })
	return c
}

// SetSegmentSymbols shows an SF Symbol on each segment of a segmented
// control, in segment order; the labels become tooltips. Extra names are
// ignored and missing or empty names leave the label in place.
func (c *MacAccessoryControl) SetSegmentSymbols(symbols ...string) *MacAccessoryControl {
	if c == nil {
		return c
	}
	c.lock.Lock()
	c.segmentSymbols = append([]string(nil), symbols...)
	segments := append([]string(nil), c.segments...)
	selected := c.selected
	c.lock.Unlock()
	c.update(func(native unsafe.Pointer, id uint64) {
		macAccessoryControlSetSegments(native, id, segments, symbols, selected)
	})
	return c
}

// SetSelectedSegment selects a segment; a negative index clears the
// selection. OnSelectionChange is not invoked for programmatic changes.
func (c *MacAccessoryControl) SetSelectedSegment(index int) *MacAccessoryControl {
	if c == nil {
		return c
	}
	c.lock.Lock()
	if index >= len(c.segments) {
		index = -1
	}
	c.selected = index
	c.lock.Unlock()
	c.update(func(native unsafe.Pointer, id uint64) { macAccessoryControlSetSelectedSegment(native, id, index) })
	return c
}

// SelectedSegment returns the selected segment index, or -1.
func (c *MacAccessoryControl) SelectedSegment() int {
	if c == nil {
		return -1
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.selected
}

// SetMenu replaces a menu button's menu.
func (c *MacAccessoryControl) SetMenu(menu *Menu) *MacAccessoryControl {
	if c == nil {
		return c
	}
	c.lock.Lock()
	c.menu = menu
	c.lock.Unlock()
	c.update(func(native unsafe.Pointer, id uint64) { macAccessoryControlSetMenu(native, id, menu) })
	return c
}

// update runs a native control update on the application thread while the
// owning accessory is attached.
func (c *MacAccessoryControl) update(apply func(native unsafe.Pointer, id uint64)) {
	if c == nil || c.accessory == nil {
		return
	}
	id := c.internalID
	c.accessory.update(func(native unsafe.Pointer) { apply(native, id) })
}

// macAccessoryControlSnapshot is a lock-free copy of one control's state.
type macAccessoryControlSnapshot struct {
	id             uint64
	kind           MacAccessoryControlKind
	text           string
	placeholder    string
	symbol         string
	tooltip        string
	width          float64
	disabled       bool
	hidden         bool
	incremental    bool
	segments       []string
	segmentSymbols []string
	selected       int
	menu           *Menu
	nativeView     unsafe.Pointer
}

func (c *MacAccessoryControl) snapshot() macAccessoryControlSnapshot {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return macAccessoryControlSnapshot{
		id:             c.internalID,
		kind:           c.kind,
		text:           c.text,
		placeholder:    c.placeholder,
		symbol:         c.symbol,
		tooltip:        c.tooltip,
		width:          c.width,
		disabled:       c.disabled,
		hidden:         c.hidden,
		incremental:    c.incremental,
		segments:       append([]string(nil), c.segments...),
		segmentSymbols: append([]string(nil), c.segmentSymbols...),
		selected:       c.selected,
		menu:           c.menu,
		nativeView:     c.nativeView,
	}
}

// macAccessorySnapshot is a lock-free copy of an accessory's configuration.
type macAccessorySnapshot struct {
	layout              MacAccessoryLayout
	hidden              bool
	height              float64
	fullScreenMinHeight float64
	adjustsSize         bool
	adjustsSizeSet      bool
	contentInsets       bool
	controls            []*MacAccessoryControl
}

func (a *MacAccessory) snapshot() macAccessorySnapshot {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return macAccessorySnapshot{
		layout:              a.layout,
		hidden:              a.hidden,
		height:              a.resolvedHeightLocked(),
		fullScreenMinHeight: a.fullScreenMinHeight,
		adjustsSize:         a.adjustsSize,
		adjustsSizeSet:      a.adjustsSizeSet,
		contentInsets:       a.contentInsets,
		controls:            append([]*MacAccessoryControl(nil), a.controls...),
	}
}

// Control registry: native callbacks carry a control identifier that is
// resolved here. Entries exist only while the owning accessory is attached.
var macAccessoryControlRegistry = make(map[uint64]*MacAccessoryControl)
var macAccessoryControlRegistryLock sync.RWMutex

func registerMacAccessoryControls(a *MacAccessory) {
	if a == nil {
		return
	}
	controls := a.Controls()
	macAccessoryControlRegistryLock.Lock()
	for _, control := range controls {
		macAccessoryControlRegistry[control.internalID] = control
	}
	macAccessoryControlRegistryLock.Unlock()
}

func unregisterMacAccessoryControls(a *MacAccessory) {
	if a == nil {
		return
	}
	controls := a.Controls()
	macAccessoryControlRegistryLock.Lock()
	for _, control := range controls {
		delete(macAccessoryControlRegistry, control.internalID)
	}
	macAccessoryControlRegistryLock.Unlock()
}

func macAccessoryControlByID(id uint64) *MacAccessoryControl {
	macAccessoryControlRegistryLock.RLock()
	defer macAccessoryControlRegistryLock.RUnlock()
	return macAccessoryControlRegistry[id]
}

// Pending queue: accessories attached before their native host exists,
// keyed by the host window or pane and drained on the application thread
// once the host is built.
var macAccessoryPending = make(map[any][]*MacAccessory)
var macAccessoryPendingLock sync.Mutex

func queueMacAccessory(key any, a *MacAccessory) {
	macAccessoryPendingLock.Lock()
	macAccessoryPending[key] = append(macAccessoryPending[key], a)
	macAccessoryPendingLock.Unlock()
}

func dequeueMacAccessory(key any, a *MacAccessory) {
	macAccessoryPendingLock.Lock()
	defer macAccessoryPendingLock.Unlock()
	queue := macAccessoryPending[key]
	kept := queue[:0]
	for _, candidate := range queue {
		if candidate != a {
			kept = append(kept, candidate)
		}
	}
	if len(kept) == 0 {
		delete(macAccessoryPending, key)
		return
	}
	macAccessoryPending[key] = kept
}

func takePendingMacAccessories(key any) []*MacAccessory {
	macAccessoryPendingLock.Lock()
	defer macAccessoryPendingLock.Unlock()
	queue := macAccessoryPending[key]
	delete(macAccessoryPending, key)
	return queue
}

func pendingMacAccessories(key any) []*MacAccessory {
	macAccessoryPendingLock.Lock()
	defer macAccessoryPendingLock.Unlock()
	return append([]*MacAccessory(nil), macAccessoryPending[key]...)
}

// AddTitlebarAccessory attaches a native control strip to the window's
// titlebar (macOS only). The accessory may use the Leading, Trailing, or
// Bottom layout. Calls made before the native window exists are queued and
// applied when it is created; later calls attach immediately. Remove the
// accessory with MacAccessory.Remove. On platforms without AppKit the call
// returns ErrMacAccessoryUnsupported.
func (w *WebviewWindow) AddTitlebarAccessory(accessory *MacAccessory) error {
	if w == nil {
		return errors.New("window is nil")
	}
	return addMacTitlebarAccessory(w, accessory, w.impl != nil)
}

// AddTitlebarAccessory attaches a native control strip to the window's
// titlebar. See WebviewWindow.AddTitlebarAccessory.
func (w *NativeWindow) AddTitlebarAccessory(accessory *MacAccessory) error {
	if w == nil {
		return errors.New("window is nil")
	}
	return addMacTitlebarAccessory(w, accessory, w.implementation() != nil)
}

// flushPendingMacAccessories attaches accessories queued before the native
// window existed. NativeWindow.Run calls it right after scheduling native
// creation; the platform implementation runs after creation on the
// application thread.
func (w *NativeWindow) flushPendingMacAccessories() {
	if w == nil {
		return
	}
	macAccessoryFlushNativeWindow(w)
}

func addMacTitlebarAccessory(window macAccessoryWindow, accessory *MacAccessory, created bool) error {
	if accessory == nil {
		return ErrMacAccessoryRequired
	}
	if !accessory.layout.titlebar() {
		return fmt.Errorf("%w: %s accessories cannot be added to a titlebar", ErrMacAccessoryLayoutInvalid, accessory.layout)
	}
	if !macAccessoriesSupported {
		return ErrMacAccessoryUnsupported
	}
	if err := accessory.claim(macAccessoryTarget{window: window}); err != nil {
		return err
	}
	if !created || globalApplication == nil {
		queueMacAccessory(window, accessory)
		return nil
	}
	if err := macAccessoryAttachTitlebar(window, accessory); err != nil {
		accessory.releaseClaim()
		return err
	}
	return nil
}

// AddTopAccessory attaches a native control strip to the top of the pane
// (NSSplitViewItemAccessoryViewController, macOS 26 and newer). The
// accessory must use the Top layout. Calls made before the split view is
// installed are queued; later calls attach immediately. Remove the accessory
// with MacAccessory.Remove.
func (p *MacSplitPane) AddTopAccessory(accessory *MacAccessory) error {
	return addMacSplitPaneAccessory(p, accessory, true)
}

// AddBottomAccessory attaches a native control strip to the bottom of the
// pane. The accessory must use the Bottom layout. See AddTopAccessory.
func (p *MacSplitPane) AddBottomAccessory(accessory *MacAccessory) error {
	return addMacSplitPaneAccessory(p, accessory, false)
}

func addMacSplitPaneAccessory(pane *MacSplitPane, accessory *MacAccessory, top bool) error {
	if accessory == nil {
		return ErrMacAccessoryRequired
	}
	if pane == nil || pane.split == nil || pane.isDead() {
		return ErrMacAccessoryPaneRequired
	}
	want := MacAccessoryLayoutBottom
	if top {
		want = MacAccessoryLayoutTop
	}
	if accessory.layout != want {
		return fmt.Errorf("%w: a %s accessory cannot be %s-aligned in a pane", ErrMacAccessoryLayoutInvalid, accessory.layout, want)
	}
	if !macAccessoriesSupported {
		return ErrMacAccessoryUnsupported
	}
	if err := accessory.claim(macAccessoryTarget{pane: pane, top: top}); err != nil {
		return err
	}
	// The installed check and the queue append happen under the pending lock
	// so an installation that completes concurrently drains this accessory
	// instead of racing past it.
	macAccessoryPendingLock.Lock()
	if !pane.split.isInstalled() || globalApplication == nil {
		macAccessoryPending[pane] = append(macAccessoryPending[pane], accessory)
		macAccessoryPendingLock.Unlock()
		return nil
	}
	macAccessoryPendingLock.Unlock()
	if err := macAccessoryAttachPane(pane, accessory, top); err != nil {
		accessory.releaseClaim()
		return err
	}
	return nil
}

// Native control events are routed through one channel drained by the
// chrome event loop registered below.
type macAccessoryEventKind int

const (
	macAccessoryEventClick macAccessoryEventKind = iota
	macAccessoryEventSearch
	macAccessoryEventSelection
)

type macAccessoryEvent struct {
	controlID uint64
	kind      macAccessoryEventKind
	text      string
	index     int
}

var macAccessoryEvents = make(chan macAccessoryEvent, 64)

func init() {
	registerChromeEventLoop(func(*App) {
		for {
			event := <-macAccessoryEvents
			go handleMacAccessoryEvent(event)
		}
	})
}

// handleMacAccessoryEvent dispatches one native event to its control's
// callback. Callbacks never run while a control lock is held.
func handleMacAccessoryEvent(event macAccessoryEvent) {
	defer handlePanic()
	control := macAccessoryControlByID(event.controlID)
	if control == nil {
		return
	}
	switch event.kind {
	case macAccessoryEventClick:
		control.lock.RLock()
		callback := control.onClick
		control.lock.RUnlock()
		if callback != nil {
			callback(newContext())
		}
	case macAccessoryEventSearch:
		control.lock.Lock()
		control.text = event.text
		callback := control.onSearch
		control.lock.Unlock()
		if callback != nil {
			callback(newContext(), event.text)
		}
	case macAccessoryEventSelection:
		control.lock.Lock()
		control.selected = event.index
		label := ""
		if event.index >= 0 && event.index < len(control.segments) {
			label = control.segments[event.index]
		}
		callback := control.onSelectionChange
		control.lock.Unlock()
		if callback != nil {
			callback(newContext(), event.index, label)
		}
	}
}
