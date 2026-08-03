package application

import (
	"fmt"
	"sync"
	"unsafe"
)

// MacSplitOrientation controls how a MacSplitWindow's panes are arranged.
type MacSplitOrientation int

const (
	// SplitHorizontal arranges panes side by side (left to right).
	SplitHorizontal MacSplitOrientation = iota
	// SplitVertical stacks panes top to bottom.
	SplitVertical
)

// MacSplitPaneBehavior maps to NSSplitViewItemBehavior. Sidebar and Inspector
// panes get automatic translucent-material styling and standard collapse
// behaviour from AppKit with no extra configuration.
type MacSplitPaneBehavior int

const (
	PaneBehaviorDefault MacSplitPaneBehavior = iota
	PaneBehaviorSidebar
	PaneBehaviorContentList
	PaneBehaviorInspector
)

// MacPaneContent is a sealed interface: the only implementations are the
// ones declared in this package. MacWebviewContent is the only kind
// implemented so far; native text/PDF/QuickLook panes are a follow-up.
type MacPaneContent interface {
	isMacPaneContent()
}

// MacWebviewContent hosts an independent WKWebView inside a split pane,
// with its own URL and preferences.
type MacWebviewContent struct {
	URL                string
	WebviewPreferences MacWebviewPreferences
}

func (MacWebviewContent) isMacPaneContent() {}

// MacSplitPaneOptions describes one pane of a MacSplitWindow.
type MacSplitPaneOptions struct {
	// Name identifies the pane for later lookup via MacSplitWindow.Pane. It
	// must be unique within a single MacSplitWindow.
	Name    string
	Content MacPaneContent
	// Behavior selects the NSSplitViewItem factory used to build this pane;
	// Sidebar/Inspector get automatic AppKit glass and default sizing.
	Behavior MacSplitPaneBehavior
	// MinThickness / MaxThickness / PreferredThicknessFraction / HoldingPriority
	// leave AppKit's own default when <= 0.
	MinThickness               float64
	MaxThickness               float64
	PreferredThicknessFraction float64
	HoldingPriority            float64
	Collapsible                bool
	StartCollapsed              bool
}

// MacSplitWindowOptions configures a new MacSplitWindow.
type MacSplitWindowOptions struct {
	Title  string
	Width  int
	Height int
	// Orientation controls whether panes sit side by side (SplitHorizontal,
	// the default) or stacked (SplitVertical).
	Orientation MacSplitOrientation
	// AutosaveName persists divider positions across relaunches via
	// NSSplitView's own autosave mechanism (NSUserDefaults-backed).
	AutosaveName string
	// Panes must contain at least one entry; the first pane reuses the
	// window's own webview instead of creating a new one.
	Panes []MacSplitPaneOptions

	TitleBar           MacTitleBar
	Appearance         MacAppearanceType
	Backdrop           MacBackdrop
	WindowLevel        MacWindowLevel
	CollectionBehavior MacWindowCollectionBehavior
}

// MacSplitPane is a handle to one pane of a MacSplitWindow.
type MacSplitPane struct {
	id       uint
	name     string
	behavior MacSplitPaneBehavior
	content  MacPaneContent

	handleLock sync.RWMutex
	// nativeItem is the NSSplitViewItem*; nativeWebview is the WKWebView*
	// hosted by this pane. Both are nil until the owning window's Run() has
	// actually executed (which may be deferred until app.Run() starts), and
	// stay nil forever on non-macOS builds.
	nativeItem    unsafe.Pointer
	nativeWebview unsafe.Pointer
}

func (p *MacSplitPane) ID() uint    { return p.id }
func (p *MacSplitPane) Name() string { return p.name }

// Content returns the pane's typed content handle: a *MacWebviewPane for a
// MacWebviewContent pane.
func (p *MacSplitPane) Content() any {
	switch p.content.(type) {
	case MacWebviewContent:
		return &MacWebviewPane{pane: p}
	default:
		return nil
	}
}

func (p *MacSplitPane) ready() (unsafe.Pointer, bool) {
	p.handleLock.RLock()
	defer p.handleLock.RUnlock()
	return p.nativeItem, p.nativeItem != nil
}

// IsCollapsed reports whether the pane is currently collapsed. Returns false
// if the pane isn't natively constructed yet.
func (p *MacSplitPane) IsCollapsed() bool {
	item, ok := p.ready()
	if !ok {
		return false
	}
	return macSplitPaneIsCollapsed(item)
}

// SetCollapsed collapses or reveals the pane. A no-op, logged, if the pane
// isn't natively constructed yet (the owning MacSplitWindow's Run() hasn't
// executed, e.g. called before app.Run()).
func (p *MacSplitPane) SetCollapsed(collapsed bool) {
	item, ok := p.ready()
	if !ok {
		globalApplication.warning("MacSplitPane %q: SetCollapsed called before the window was ready", p.name)
		return
	}
	InvokeSync(func() {
		macSplitPaneSetCollapsed(item, collapsed)
	})
}

// MacWebviewPane is the typed content handle for a pane whose content is
// MacWebviewContent.
type MacWebviewPane struct {
	pane *MacSplitPane
}

func (wp *MacWebviewPane) webview() (unsafe.Pointer, bool) {
	wp.pane.handleLock.RLock()
	defer wp.pane.handleLock.RUnlock()
	return wp.pane.nativeWebview, wp.pane.nativeWebview != nil
}

func (wp *MacWebviewPane) SetURL(url string) {
	webview, ok := wp.webview()
	if !ok {
		globalApplication.warning("MacWebviewPane %q: SetURL called before the window was ready", wp.pane.name)
		return
	}
	InvokeSync(func() {
		macSplitPaneWebviewSetURL(webview, url)
	})
}

func (wp *MacWebviewPane) Reload() {
	webview, ok := wp.webview()
	if !ok {
		globalApplication.warning("MacWebviewPane %q: Reload called before the window was ready", wp.pane.name)
		return
	}
	InvokeAsync(func() {
		macSplitPaneWebviewReload(webview)
	})
}

func (wp *MacWebviewPane) ExecJS(js string) {
	webview, ok := wp.webview()
	if !ok {
		globalApplication.warning("MacWebviewPane %q: ExecJS called before the window was ready", wp.pane.name)
		return
	}
	InvokeAsync(func() {
		macSplitPaneWebviewExecJS(webview, js)
	})
}

// macSplitPendingConfig is stashed on the underlying WebviewWindow between
// NewSplitWindow and the point where that window's Run() actually executes
// (which may be deferred until app.Run() starts, e.g. when NewSplitWindow is
// called from main() before app.Run() -- the same construction-before-Run
// gap SetToolbar has to handle). Darwin's run() consumes it via
// installSplitPanes; it's never populated on any other platform.
type macSplitPendingConfig struct {
	vertical     bool
	autosaveName string
	paneOptions  []MacSplitPaneOptions
	panes        []*MacSplitPane
}

// MacSplitWindow is a native macOS window whose content is an
// NSSplitViewController rather than a single webview. It is built on top of
// a real WebviewWindow, so it gets full window chrome (toolbar, standard
// titlebar behaviour, etc.) for free; only the content area is split into
// panes.
type MacSplitWindow struct {
	window *WebviewWindow
	panes  []*MacSplitPane
}

// NewSplitWindow creates a new native split-pane window. It is macOS only;
// on every other platform it returns an error.
func NewSplitWindow(options MacSplitWindowOptions) (*MacSplitWindow, error) {
	if len(options.Panes) == 0 {
		return nil, fmt.Errorf("MacSplitWindowOptions.Panes must contain at least one pane")
	}
	seenNames := make(map[string]bool, len(options.Panes))
	for i, paneOptions := range options.Panes {
		if paneOptions.Name == "" {
			return nil, fmt.Errorf("pane %d: Name is required", i)
		}
		if seenNames[paneOptions.Name] {
			return nil, fmt.Errorf("pane %d: duplicate pane name %q", i, paneOptions.Name)
		}
		seenNames[paneOptions.Name] = true
		if _, ok := paneOptions.Content.(MacWebviewContent); !ok {
			return nil, fmt.Errorf("pane %q: only MacWebviewContent is currently supported", paneOptions.Name)
		}
	}
	return newSplitWindow(options)
}

func (sw *MacSplitWindow) Show() *MacSplitWindow {
	// WebviewWindow.Show(), when the native window hasn't run yet (true here
	// until app.Run() starts, since NewSplitWindow always constructs with
	// Hidden:true to avoid flashing an unsplit window), triggers Run() but
	// -- unlike its already-running branch -- never clears options.Hidden.
	// Left alone, a Show() called before app.Run() (the normal call site,
	// right after NewSplitWindow in main()) would construct the window and
	// silently never display it. Same shape as the SetToolbar bug fixed
	// earlier in this window's Run() path; clear it here explicitly.
	sw.window.options.Hidden = false
	sw.window.Show()
	return sw
}

func (sw *MacSplitWindow) Hide() *MacSplitWindow {
	sw.window.Hide()
	return sw
}

func (sw *MacSplitWindow) Close() {
	sw.window.Close()
}

func (sw *MacSplitWindow) Focus() *MacSplitWindow {
	sw.window.Focus()
	return sw
}

func (sw *MacSplitWindow) SetPosition(x, y int) *MacSplitWindow {
	sw.window.SetPosition(x, y)
	return sw
}

func (sw *MacSplitWindow) SetSize(width, height int) *MacSplitWindow {
	sw.window.SetSize(width, height)
	return sw
}

// Window returns the underlying WebviewWindow, e.g. to attach a MacToolbar
// via SetToolbar.
func (sw *MacSplitWindow) Window() *WebviewWindow {
	return sw.window
}

func (sw *MacSplitWindow) Panes() []*MacSplitPane {
	result := make([]*MacSplitPane, len(sw.panes))
	copy(result, sw.panes)
	return result
}

func (sw *MacSplitWindow) Pane(name string) (*MacSplitPane, bool) {
	for _, pane := range sw.panes {
		if pane.name == name {
			return pane, true
		}
	}
	return nil, false
}
