package application

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"unsafe"
)

// macSplitPaneRole is the semantic AppKit role of a split-view pane. The role
// selects the NSSplitViewItem factory and therefore the system's default
// material, sizing, and collapse behavior.
type macSplitPaneRole uint8

const (
	macSplitPaneSidebar macSplitPaneRole = iota
	macSplitPanePrimary
	macSplitPaneInspector
)

// MacSplitView is a native, side-by-side NSSplitViewController layout for one
// WebviewWindow (macOS only). Construct panes in leading-to-trailing order
// with the Add methods, then attach the layout with
// WebviewWindow.SetSplitView before the window is shown.
//
// The window's existing WebView becomes the pane added with
// AddPrimaryContent. A sidebar added with AddSidebar is a native AppKit source
// list and does not create another WebView. On platforms without AppKit the
// window keeps its ordinary single WebView and the native handles are inert.
type MacSplitView struct {
	// lock protects every field below. Native installation state is only
	// mutated on the application thread while this lock is held.
	lock         sync.RWMutex
	autosaveName string
	panes        []*MacSplitWebviewPane
	owner        macSplitWindow
	frozen       bool
	installed    bool
	native       unsafe.Pointer
}

// macSplitWindow is the temporary internal host boundary shared by
// WebviewWindow and NativeWindow. Keeping it private lets v4 redesign the
// public Window interface without exposing another compatibility surface.
type macSplitWindow interface {
	ID() uint
	Error(string, ...any)
	macSplitOptions() MacWindow
}

// NewMacSplitView creates an empty native split-view layout.
func NewMacSplitView() *MacSplitView {
	return &MacSplitView{}
}

// SetAutosaveName enables AppKit's divider-position persistence. Passing an
// empty name disables persistence. The name should be stable per logical
// window and split layout; it is read once when the layout is installed.
func (s *MacSplitView) SetAutosaveName(name string) *MacSplitView {
	if s == nil {
		return s
	}
	s.lock.Lock()
	s.autosaveName = name
	s.lock.Unlock()
	return s
}

// AddSidebar adds a native AppKit source list. It normally appears first so it
// occupies the leading edge. A sidebar can belong to only one split view.
func (s *MacSplitView) AddSidebar(sidebar *MacSidebar) *MacSplitPane {
	if sidebar == nil {
		return nil
	}
	pane := s.addPane(macSplitPaneSidebar, false)
	if pane == nil {
		return nil
	}
	sidebar.lock.Lock()
	if sidebar.pane != nil {
		sidebar.lock.Unlock()
		s.lock.Lock()
		s.panes = s.panes[:len(s.panes)-1]
		s.lock.Unlock()
		return nil
	}
	sidebar.pane = pane.MacSplitPane
	sidebar.lock.Unlock()
	pane.sidebar = sidebar
	return pane.MacSplitPane
}

// AddPrimaryContent places the window's existing WKWebView at this position.
// Exactly one primary content pane is required. Its initial URL and WebView
// preferences come from WebviewWindowOptions, and the outer window's WebView
// methods keep addressing it.
func (s *MacSplitView) AddPrimaryContent() *MacSplitWebviewPane {
	return s.addPane(macSplitPanePrimary, true)
}

// AddTextEditor adds a native NSTextView-backed primary content pane. A split
// layout contains either this pane for NativeWindow or AddPrimaryContent for
// WebviewWindow, never both.
func (s *MacSplitView) AddTextEditor(editor *MacTextEditor) *MacSplitPane {
	if editor == nil {
		return nil
	}
	pane := s.addPane(macSplitPanePrimary, true)
	if pane == nil {
		return nil
	}
	editor.lock.Lock()
	if editor.pane != nil {
		editor.lock.Unlock()
		s.lock.Lock()
		s.panes = s.panes[:len(s.panes)-1]
		s.lock.Unlock()
		return nil
	}
	editor.pane = pane.MacSplitPane
	editor.lock.Unlock()
	pane.editor = editor
	return pane.MacSplitPane
}

// AddInspector adds a native AppKit property inspector. It normally appears
// last so it occupies the trailing edge. An inspector can belong to only one
// split view. On macOS 11 and newer Wails uses AppKit's semantic inspector
// split-item role; older supported releases retain the same controls and
// collapse behavior in a regular trailing split item.
func (s *MacSplitView) AddInspector(inspector *MacInspector) *MacSplitPane {
	if inspector == nil {
		return nil
	}
	pane := s.addPane(macSplitPaneInspector, false)
	if pane == nil {
		return nil
	}
	inspector.lock.Lock()
	if inspector.pane != nil {
		inspector.lock.Unlock()
		s.lock.Lock()
		s.panes = s.panes[:len(s.panes)-1]
		s.lock.Unlock()
		return nil
	}
	inspector.pane = pane.MacSplitPane
	inspector.lock.Unlock()
	pane.inspector = inspector
	return pane.MacSplitPane
}

var macSplitPaneID uint64

func nextMacSplitPaneID() uint64 {
	return atomic.AddUint64(&macSplitPaneID, 1)
}

func (s *MacSplitView) addPane(role macSplitPaneRole, primary bool) *MacSplitWebviewPane {
	if s == nil {
		return nil
	}
	s.lock.Lock()
	if s.frozen {
		owner := s.owner
		s.lock.Unlock()
		reportMacSplitError(owner, "split view structure cannot be changed after attachment")
		return nil
	}
	pane := &MacSplitWebviewPane{
		MacSplitPane: &MacSplitPane{
			internalID: nextMacSplitPaneID(),
			role:       role,
			split:      s,
		},
		primary: primary,
	}
	s.panes = append(s.panes, pane)
	s.lock.Unlock()
	return pane
}

func (s *MacSplitView) paneSnapshot() []*MacSplitWebviewPane {
	if s == nil {
		return nil
	}
	s.lock.RLock()
	defer s.lock.RUnlock()
	return append([]*MacSplitWebviewPane(nil), s.panes...)
}

func (s *MacSplitView) primaryPane() *MacSplitWebviewPane {
	for _, pane := range s.paneSnapshot() {
		if pane.primary {
			return pane
		}
	}
	return nil
}

func (s *MacSplitView) hasSidebarPane() bool {
	for _, pane := range s.paneSnapshot() {
		if pane.role == macSplitPaneSidebar {
			return true
		}
	}
	return false
}

func (s *MacSplitView) hasInspectorPane() bool {
	for _, pane := range s.paneSnapshot() {
		if pane.role == macSplitPaneInspector {
			return true
		}
	}
	return false
}

func (s *MacSplitView) inspectorPane() *MacSplitPane {
	for _, pane := range s.paneSnapshot() {
		if pane.role == macSplitPaneInspector {
			return pane.MacSplitPane
		}
	}
	return nil
}

func (s *MacSplitView) ownerWindow() macSplitWindow {
	if s == nil {
		return nil
	}
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.owner
}

func (s *MacSplitView) ownerWebviewWindow() *WebviewWindow {
	owner, _ := s.ownerWindow().(*WebviewWindow)
	return owner
}

func (s *MacSplitView) ownerNativeWindow() *NativeWindow {
	owner, _ := s.ownerWindow().(*NativeWindow)
	return owner
}

func (s *MacSplitView) isInstalled() bool {
	if s == nil {
		return false
	}
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.installed
}

// reportMacSplitError logs through the owning window when one exists. It stays
// silent otherwise so building a layout does not require a running
// application.
func reportMacSplitError(owner macSplitWindow, format string, args ...any) {
	if owner != nil && globalApplication != nil {
		owner.Error(format, args...)
	}
}

// validateMacSplitView checks a layout before any native mutation. A failure
// leaves the current window unchanged.
func validateMacSplitView(s *MacSplitView) error {
	panes := s.paneSnapshot()
	if len(panes) < 2 {
		return fmt.Errorf("split view requires at least two panes")
	}
	primaries := 0
	for _, pane := range panes {
		if pane == nil {
			return fmt.Errorf("split view contains a nil pane")
		}
		if pane.primary {
			primaries++
		}
		if pane.role == macSplitPaneSidebar && pane.sidebar == nil {
			return fmt.Errorf("split view sidebar pane requires a native sidebar")
		}
		if pane.role == macSplitPaneInspector && pane.inspector == nil {
			return fmt.Errorf("split view inspector pane requires a native inspector")
		}
		if err := pane.validate(); err != nil {
			return err
		}
	}
	if primaries != 1 {
		return fmt.Errorf("split view requires exactly one primary content pane, found %d", primaries)
	}
	return nil
}

func claimMacSplitView(split *MacSplitView, window macSplitWindow) error {
	split.lock.Lock()
	defer split.lock.Unlock()
	if split.owner != nil && split.owner != window {
		return fmt.Errorf("split view is already attached to another window")
	}
	split.owner = window
	split.frozen = true
	return nil
}

func releaseMacSplitViewOwnership(split *MacSplitView, window macSplitWindow) {
	if split == nil {
		return
	}
	split.lock.Lock()
	if split.owner == window && !split.installed {
		split.owner = nil
		split.frozen = false
	}
	split.lock.Unlock()
}

// adoptPendingPrimaryState forwards a URL or JavaScript that was configured on
// the primary pane before the split had a window to delegate to.
func (s *MacSplitView) adoptPendingPrimaryState(window *WebviewWindow) {
	primary := s.primaryPane()
	if primary == nil {
		return
	}
	primary.lock.Lock()
	url := primary.url
	primary.url = ""
	pending := primary.pendingJS
	primary.pendingJS = nil
	primary.lock.Unlock()
	if url != "" {
		window.SetURL(url)
	}
	for _, js := range pending {
		window.ExecJS(js)
	}
}

// MacSplitPane configures and observes one NSSplitViewItem. Configuration
// made before the native window exists is applied during installation; later
// calls update the native item on the application thread.
type MacSplitPane struct {
	lock sync.RWMutex

	internalID uint64
	role       macSplitPaneRole
	split      *MacSplitView

	minimumThickness         float64
	maximumThickness         float64
	preferredThickness       float64
	preferredThicknessSet    bool
	holdingPriority          float64
	holdingPrioritySet       bool
	collapsible              bool
	collapsibleSet           bool
	canCollapseFromResize    bool
	canCollapseFromResizeSet bool
	collapsed                bool
	onCollapsedChange        func(*Context, bool)

	native unsafe.Pointer
	dead   bool
}

func (p *MacSplitPane) validate() error {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if !isFiniteMacSplitNumber(p.minimumThickness) || !isFiniteMacSplitNumber(p.maximumThickness) {
		return fmt.Errorf("split pane thickness must be finite")
	}
	if p.minimumThickness < 0 || p.maximumThickness < 0 {
		return fmt.Errorf("split pane thickness must not be negative")
	}
	if p.minimumThickness > 0 && p.maximumThickness > 0 && p.minimumThickness > p.maximumThickness {
		return fmt.Errorf("split pane minimum thickness %.1f exceeds maximum %.1f", p.minimumThickness, p.maximumThickness)
	}
	if p.preferredThicknessSet && (!isFiniteMacSplitNumber(p.preferredThickness) || p.preferredThickness <= 0 || p.preferredThickness > 1) {
		return fmt.Errorf("split pane preferred thickness fraction must be within (0, 1]")
	}
	if p.holdingPrioritySet && (!isFiniteMacSplitNumber(p.holdingPriority) || p.holdingPriority < 1 || p.holdingPriority > 1000) {
		return fmt.Errorf("split pane holding priority must be within 1...1000")
	}
	return nil
}

func isFiniteMacSplitNumber(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

// SetMinimumThickness sets the pane's minimum width in points. Passing zero
// restores AppKit's role-specific default. Invalid values are ignored and
// reported through the owning window's error handler.
func (p *MacSplitPane) SetMinimumThickness(points float64) *MacSplitPane {
	if p == nil {
		return p
	}
	p.lock.Lock()
	if p.dead {
		p.lock.Unlock()
		return p
	}
	if !isFiniteMacSplitNumber(points) || points < 0 ||
		(points > 0 && p.maximumThickness > 0 && points > p.maximumThickness) {
		p.lock.Unlock()
		reportMacSplitError(p.split.ownerWindow(), "SetMinimumThickness: value %.1f is invalid for the pane's current maximum", points)
		return p
	}
	p.minimumThickness = points
	p.lock.Unlock()
	macSplitPaneApplyMinimumThickness(p, points)
	return p
}

// SetMaximumThickness sets the pane's maximum width in points. Passing zero
// restores AppKit's role-specific default. Invalid values are ignored and
// reported through the owning window's error handler.
func (p *MacSplitPane) SetMaximumThickness(points float64) *MacSplitPane {
	if p == nil {
		return p
	}
	p.lock.Lock()
	if p.dead {
		p.lock.Unlock()
		return p
	}
	if !isFiniteMacSplitNumber(points) || points < 0 ||
		(points > 0 && p.minimumThickness > 0 && p.minimumThickness > points) {
		p.lock.Unlock()
		reportMacSplitError(p.split.ownerWindow(), "SetMaximumThickness: value %.1f is invalid for the pane's current minimum", points)
		return p
	}
	p.maximumThickness = points
	p.lock.Unlock()
	macSplitPaneApplyMaximumThickness(p, points)
	return p
}

// SetPreferredThicknessFraction sets the fraction of the split view's width
// this pane prefers, within (0, 1].
func (p *MacSplitPane) SetPreferredThicknessFraction(fraction float64) *MacSplitPane {
	if p == nil {
		return p
	}
	p.lock.Lock()
	if p.dead {
		p.lock.Unlock()
		return p
	}
	if !isFiniteMacSplitNumber(fraction) || fraction <= 0 || fraction > 1 {
		p.lock.Unlock()
		reportMacSplitError(p.split.ownerWindow(), "SetPreferredThicknessFraction: value must be within (0, 1]")
		return p
	}
	p.preferredThickness = fraction
	p.preferredThicknessSet = true
	p.lock.Unlock()
	macSplitPaneApplyPreferredFraction(p, fraction)
	return p
}

// SetHoldingPriority sets the Auto Layout priority with which the pane resists
// resizing, within AppKit's useful 1...1000 range. Lower values give up space
// first.
func (p *MacSplitPane) SetHoldingPriority(priority float64) *MacSplitPane {
	if p == nil {
		return p
	}
	p.lock.Lock()
	if p.dead {
		p.lock.Unlock()
		return p
	}
	if !isFiniteMacSplitNumber(priority) || priority < 1 || priority > 1000 {
		p.lock.Unlock()
		reportMacSplitError(p.split.ownerWindow(), "SetHoldingPriority: value must be within 1...1000")
		return p
	}
	p.holdingPriority = priority
	p.holdingPrioritySet = true
	p.lock.Unlock()
	macSplitPaneApplyHoldingPriority(p, priority)
	return p
}

// SetCollapsible controls whether the user can collapse the pane. Panes keep
// their AppKit role default until this is called.
func (p *MacSplitPane) SetCollapsible(collapsible bool) *MacSplitPane {
	if p == nil {
		return p
	}
	p.lock.Lock()
	if p.dead {
		p.lock.Unlock()
		return p
	}
	p.collapsible = collapsible
	p.collapsibleSet = true
	p.lock.Unlock()
	macSplitPaneApplyCollapsible(p, collapsible)
	return p
}

// SetCanCollapseFromWindowResize controls whether shrinking the window may
// collapse the pane, on macOS releases that support the behavior.
func (p *MacSplitPane) SetCanCollapseFromWindowResize(allowed bool) *MacSplitPane {
	if p == nil {
		return p
	}
	p.lock.Lock()
	if p.dead {
		p.lock.Unlock()
		return p
	}
	p.canCollapseFromResize = allowed
	p.canCollapseFromResizeSet = true
	p.lock.Unlock()
	macSplitPaneApplyCanCollapseFromResize(p, allowed)
	return p
}

// SetCollapsed collapses or expands the pane. Before installation it sets the
// pane's initial state; afterwards the change animates and is reported
// through OnCollapsedChange like any other collapse source.
func (p *MacSplitPane) SetCollapsed(collapsed bool) *MacSplitPane {
	if p == nil {
		return p
	}
	p.lock.RLock()
	dead := p.dead
	p.lock.RUnlock()
	if dead {
		return p
	}
	if p.split.isInstalled() {
		macSplitPaneApplyCollapsed(p, collapsed)
		return p
	}
	p.lock.Lock()
	p.collapsed = collapsed
	p.lock.Unlock()
	return p
}

// Toggle flips the pane's collapse state.
func (p *MacSplitPane) Toggle() *MacSplitPane {
	if p == nil {
		return p
	}
	p.lock.RLock()
	dead := p.dead
	p.lock.RUnlock()
	if dead {
		return p
	}
	if p.split.isInstalled() {
		macSplitPaneApplyToggle(p)
		return p
	}
	p.lock.Lock()
	p.collapsed = !p.collapsed
	p.lock.Unlock()
	return p
}

// IsCollapsed reports the pane's last known collapse state. After a
// SetCollapsed or Toggle call on an installed pane the value updates when
// AppKit reports the change.
func (p *MacSplitPane) IsCollapsed() bool {
	if p == nil {
		return false
	}
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.collapsed
}

// OnCollapsedChange sets the callback invoked whenever the pane's collapse
// state changes, regardless of whether the change came from the sidebar
// toggle, a menu, a divider gesture, window resizing, or SetCollapsed.
// Passing nil clears the callback.
func (p *MacSplitPane) OnCollapsedChange(callback func(*Context, bool)) *MacSplitPane {
	if p == nil {
		return p
	}
	p.lock.Lock()
	if p.dead {
		p.lock.Unlock()
		return p
	}
	p.onCollapsedChange = callback
	p.lock.Unlock()
	return p
}

func (p *MacSplitWebviewPane) markDead() {
	if p == nil {
		return
	}
	p.lock.Lock()
	p.dead = true
	p.native = nil
	p.loaded = false
	p.pendingJS = nil
	p.onLoaded = nil
	p.onCollapsedChange = nil
	sidebar := p.sidebar
	inspector := p.inspector
	p.url = ""
	p.lock.Unlock()
	if sidebar != nil {
		sidebar.markDead()
	}
	if inspector != nil {
		inspector.markDead()
	}
	if p.editor != nil {
		p.editor.markDead()
	}
}

func (p *MacSplitPane) isDead() bool {
	if p == nil {
		return true
	}
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.dead
}

// MacSplitWebviewPane is the split view's primary web-content pane. It
// addresses the window's existing WKWebView; native sidebars and inspectors
// use MacSplitPane and never expose web-navigation methods.
type MacSplitWebviewPane struct {
	*MacSplitPane

	// The embedded pane's lock protects these fields as well.
	url     string
	primary bool
	// contentLayout overrides MacWindow.ContentLayout when non-automatic.
	contentLayout MacContentLayout
	sidebar       *MacSidebar
	inspector     *MacInspector
	editor        *MacTextEditor
	loaded        bool
	// navigationGeneration increments synchronously when WebKit starts a new
	// navigation. Completion events capture it before leaving AppKit so a stale
	// completion cannot flush a newer navigation's JavaScript queue.
	navigationGeneration uint64
	pendingJS            []string
	onLoaded             func(*Context)
}

// SetContentLayout controls whether this primary WebView is constrained below
// the window toolbar or extends underneath it. Automatic inherits the owning
// window's MacWindow.ContentLayout, which itself follows FullSizeContent when
// left automatic. Edge-to-edge layout lets WKWebView participate in AppKit's
// automatic scroll-edge effect on macOS 26 and newer.
//
// Calls made after installation update the native constraints immediately.
func (p *MacSplitWebviewPane) SetContentLayout(layout MacContentLayout) *MacSplitWebviewPane {
	if p == nil {
		return p
	}
	if !validMacContentLayout(layout) {
		reportMacSplitError(p.split.ownerWindow(), "SetContentLayout: unknown layout value %d", layout)
		return p
	}
	p.lock.Lock()
	if p.dead {
		p.lock.Unlock()
		return p
	}
	p.contentLayout = layout
	p.lock.Unlock()
	owner := p.split.ownerWebviewWindow()
	if owner != nil {
		macSplitPaneApplyContentLayout(p.MacSplitPane, resolveMacContentLayout(owner.options.Mac, layout))
	}
	return p
}

// ContentLayout returns the pane's configured layout override. Automatic
// means the pane inherits the owning window's content-layout policy.
func (p *MacSplitWebviewPane) ContentLayout() MacContentLayout {
	if p == nil {
		return MacContentLayoutAutomatic
	}
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.contentLayout
}

// SetURL navigates the pane. For the primary content pane this is equivalent
// to calling SetURL on the window itself.
func (p *MacSplitWebviewPane) SetURL(url string) *MacSplitWebviewPane {
	if p == nil {
		return p
	}
	p.lock.RLock()
	dead := p.dead
	p.lock.RUnlock()
	if dead {
		return p
	}
	if !p.primary {
		return p
	}
	p.lock.Lock()
	p.navigationGeneration++
	p.loaded = false
	p.lock.Unlock()
	owner := p.split.ownerWebviewWindow()
	if owner != nil {
		owner.SetURL(url)
		return p
	}
	p.lock.Lock()
	p.url = url
	p.lock.Unlock()
	return p
}

// Reload reloads the pane's current page.
func (p *MacSplitWebviewPane) Reload() *MacSplitWebviewPane {
	if p == nil {
		return p
	}
	p.lock.RLock()
	dead := p.dead
	p.lock.RUnlock()
	if dead {
		return p
	}
	if !p.primary {
		return p
	}
	p.lock.Lock()
	p.navigationGeneration++
	p.loaded = false
	p.lock.Unlock()
	owner := p.split.ownerWebviewWindow()
	if owner != nil {
		owner.Reload()
	}
	return p
}

// ExecJS executes JavaScript in the pane. Auxiliary panes queue the script
// until the pane finishes loading, including after a later SetURL call. The
// primary content pane follows the window's own ExecJS semantics.
func (p *MacSplitWebviewPane) ExecJS(js string) *MacSplitWebviewPane {
	if p == nil {
		return p
	}
	p.lock.RLock()
	dead := p.dead
	p.lock.RUnlock()
	if dead {
		return p
	}
	if !p.primary {
		return p
	}
	owner := p.split.ownerWebviewWindow()
	if owner != nil && owner.impl != nil {
		owner.ExecJS(js)
		return p
	}
	p.lock.Lock()
	if owner == nil || owner.impl == nil {
		p.pendingJS = append(p.pendingJS, js)
	}
	p.lock.Unlock()
	return p
}

// OnLoaded sets the callback invoked when the pane finishes loading a page.
// Passing nil clears the callback.
func (p *MacSplitWebviewPane) OnLoaded(callback func(*Context)) *MacSplitWebviewPane {
	if p == nil {
		return p
	}
	p.lock.Lock()
	if p.dead {
		p.lock.Unlock()
		return p
	}
	p.onLoaded = callback
	p.lock.Unlock()
	return p
}

// macSplitPaneRegistry routes native pane callbacks to their Go handles.
// Entries live from native pane creation until teardown begins.
var macSplitPaneRegistry = make(map[uint64]*MacSplitWebviewPane)
var macSplitPaneRegistryLock sync.RWMutex

func registerMacSplitPane(pane *MacSplitWebviewPane) {
	macSplitPaneRegistryLock.Lock()
	macSplitPaneRegistry[pane.internalID] = pane
	macSplitPaneRegistryLock.Unlock()
}

func unregisterMacSplitPane(id uint64) {
	macSplitPaneRegistryLock.Lock()
	delete(macSplitPaneRegistry, id)
	macSplitPaneRegistryLock.Unlock()
}

func macSplitPaneByID(id uint64) *MacSplitWebviewPane {
	macSplitPaneRegistryLock.RLock()
	defer macSplitPaneRegistryLock.RUnlock()
	return macSplitPaneRegistry[id]
}

type splitPaneLoadEvent struct {
	paneID     uint64
	generation uint64
}

var splitPaneLoadedEvents = make(chan splitPaneLoadEvent, 128)

type splitPaneCollapseEvent struct {
	paneID    uint64
	collapsed bool
}

var splitPaneCollapseEvents = make(chan splitPaneCollapseEvent, 128)

func newMacSplitPaneLoadEvent(paneID uint64) (splitPaneLoadEvent, bool) {
	pane := macSplitPaneByID(paneID)
	if pane == nil {
		return splitPaneLoadEvent{}, false
	}
	pane.lock.RLock()
	defer pane.lock.RUnlock()
	if pane.dead {
		return splitPaneLoadEvent{}, false
	}
	return splitPaneLoadEvent{paneID: paneID, generation: pane.navigationGeneration}, true
}

// handleMacSplitPaneLoaded marks one pane loaded, flushes only that pane's
// queued JavaScript in insertion order, and dispatches OnLoaded. Callbacks
// never run while a pane lock is held.
func handleMacSplitPaneLoaded(event splitPaneLoadEvent) {
	defer handlePanic()
	pane := macSplitPaneByID(event.paneID)
	if pane == nil {
		return
	}
	pane.lock.Lock()
	if pane.dead || event.generation != pane.navigationGeneration {
		pane.lock.Unlock()
		return
	}
	wasLoaded := pane.loaded
	pane.loaded = true
	pending := pane.pendingJS
	pane.pendingJS = nil
	callback := pane.onLoaded
	pane.lock.Unlock()
	for _, js := range pending {
		if owner := pane.split.ownerWebviewWindow(); owner != nil && owner.impl != nil {
			owner.ExecJS(js)
		}
	}
	if !wasLoaded && callback != nil {
		callback(newContext())
	}
}

func handleMacSplitPaneNavigationStarted(paneID uint64) {
	pane := macSplitPaneByID(paneID)
	if pane == nil {
		return
	}
	pane.lock.Lock()
	if pane.dead {
		pane.lock.Unlock()
		return
	}
	pane.navigationGeneration++
	pane.loaded = false
	pane.lock.Unlock()
}

// handleMacSplitPaneCollapsed stores the pane's new collapse state and
// invokes OnCollapsedChange when the value actually changed.
func handleMacSplitPaneCollapsed(paneID uint64, collapsed bool) {
	defer handlePanic()
	pane := macSplitPaneByID(paneID)
	if pane == nil {
		return
	}
	pane.lock.Lock()
	if pane.dead {
		pane.lock.Unlock()
		return
	}
	changed := pane.collapsed != collapsed
	pane.collapsed = collapsed
	callback := pane.onCollapsedChange
	pane.lock.Unlock()
	if changed && callback != nil {
		callback(newContext(), collapsed)
	}
}
