package application

import (
	"fmt"
	"maps"
	"math"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/internal/assetserver"
)

// webviewPanelImpl is the platform-specific interface for WebviewPanel
type webviewPanelImpl interface {
	// Lifecycle
	create()
	destroy()

	// Position and size
	setBounds(bounds Rect)
	bounds() Rect
	setZIndex(zIndex int)

	// Content
	setURL(url string)
	reload()
	forceReload()

	// Visibility
	show()
	hide()
	isVisible() bool

	// Zoom
	setZoom(zoom float64)
	getZoom() float64

	// DevTools
	openDevTools()

	// Focus
	focus()
	isFocused() bool
}

var panelID uint32

func getNextPanelID() uint {
	return uint(atomic.AddUint32(&panelID, 1))
}

// WebviewPanel represents an embedded webview within a window.
// Unlike WebviewWindow, a WebviewPanel is a child view that exists within
// a parent window and can be positioned anywhere within that window.
// This is similar to Electron's BrowserView or the deprecated webview tag.
type WebviewPanel struct {
	id      uint
	name    string
	options WebviewPanelOptions
	impl    webviewPanelImpl
	parent  *WebviewWindow

	// Track if the panel has been destroyed
	creating        bool // Accessed only on the UI thread.
	nativeDestroyed bool // Accessed only on the UI thread.
	destroyed       bool
	destroyedLock   sync.RWMutex

	// Original window size when panel was created (for anchor calculations)
	originalWindowWidth  int
	originalWindowHeight int
	// Original panel bounds (for anchor calculations)
	originalBounds Rect
}

// NewPanel creates a new WebviewPanel with the given options.
// Typically called via window.NewPanel() to associate the panel with a parent window.
func NewPanel(options WebviewPanelOptions) *WebviewPanel {
	id := getNextPanelID()

	// Apply defaults
	if options.Width <= 0 {
		options.Width = 400
	}
	if options.Height <= 0 {
		options.Height = 300
	}
	if options.ZIndex == 0 {
		options.ZIndex = 1
	}
	if options.Zoom <= 0 || math.IsNaN(options.Zoom) || math.IsInf(options.Zoom, 0) {
		options.Zoom = 1.0
	}
	if options.Name == "" {
		options.Name = fmt.Sprintf("panel-%d", id)
	}
	// Default to visible
	if options.Visible == nil {
		visible := true
		options.Visible = &visible
	}

	options.Headers = maps.Clone(options.Headers)
	if options.Visible != nil {
		visible := *options.Visible
		options.Visible = &visible
	}
	if options.DevToolsEnabled != nil {
		enabled := *options.DevToolsEnabled
		options.DevToolsEnabled = &enabled
	}

	// Normalize URL via asset server for local paths
	if options.URL != "" {
		if normalizedURL, err := assetserver.GetStartURL(options.URL); err == nil && normalizedURL != "" {
			options.URL = normalizedURL
		} else {
			options.URL = ""
		}
	}

	// Store original bounds for anchor calculations
	originalBounds := Rect{
		X:      options.X,
		Y:      options.Y,
		Width:  options.Width,
		Height: options.Height,
	}

	return &WebviewPanel{
		id:             id,
		name:           options.Name,
		options:        options,
		originalBounds: originalBounds,
	}
}

// ID returns the unique identifier for this panel
func (p *WebviewPanel) ID() uint {
	return p.id
}

// Name returns the name of this panel
func (p *WebviewPanel) Name() string {
	return p.name
}

// Parent returns the parent window of this panel
func (p *WebviewPanel) Parent() *WebviewWindow {
	return p.parent
}

// SetBounds sets the position and size of the panel within its parent window.
// This also updates the anchor baseline so future window resizes calculate from the new position.
func (p *WebviewPanel) SetBounds(bounds Rect) *WebviewPanel {
	bounds.Width = max(1, bounds.Width)
	bounds.Height = max(1, bounds.Height)
	p.updateAnchorBaseline(bounds)
	p.setBoundsInternal(bounds)
	return p
}

// updateAnchorBaseline updates the original bounds and window size used for anchor calculations.
// Called when the user manually changes panel bounds.
func (p *WebviewPanel) updateAnchorBaseline(bounds Rect) {
	var width, height int
	if p.parent != nil {
		width, height = p.parentSize()
	}
	p.destroyedLock.Lock()
	defer p.destroyedLock.Unlock()
	p.originalBounds = bounds
	p.originalWindowWidth, p.originalWindowHeight = width, height
}

// Bounds returns the current bounds of the panel
func (p *WebviewPanel) Bounds() Rect {
	options := p.snapshotOptions()
	bounds := Rect{X: options.X, Y: options.Y, Width: options.Width, Height: options.Height}
	p.withImpl(func(impl webviewPanelImpl) { bounds = impl.bounds() })
	return bounds
}

// SetPosition sets the position of the panel within its parent window
func (p *WebviewPanel) SetPosition(x, y int) *WebviewPanel {
	bounds := p.Bounds()
	bounds.X = x
	bounds.Y = y
	return p.SetBounds(bounds)
}

// Position returns the current position of the panel
func (p *WebviewPanel) Position() (int, int) {
	bounds := p.Bounds()
	return bounds.X, bounds.Y
}

// SetSize sets the size of the panel
func (p *WebviewPanel) SetSize(width, height int) *WebviewPanel {
	bounds := p.Bounds()
	bounds.Width = width
	bounds.Height = height
	return p.SetBounds(bounds)
}

// Size returns the current size of the panel
func (p *WebviewPanel) Size() (int, int) {
	bounds := p.Bounds()
	return bounds.Width, bounds.Height
}

// SetZIndex sets the stacking order of the panel
func (p *WebviewPanel) SetZIndex(zIndex int) *WebviewPanel {
	p.destroyedLock.Lock()
	p.options.ZIndex = zIndex
	p.destroyedLock.Unlock()
	p.withImpl(func(impl webviewPanelImpl) { impl.setZIndex(zIndex) })
	return p
}

// ZIndex returns the current z-index of the panel
func (p *WebviewPanel) ZIndex() int {
	return p.snapshotOptions().ZIndex
}

// SetURL navigates the panel to the specified URL
// Local paths (e.g., "/panel.html") are normalized via the asset server.
func (p *WebviewPanel) SetURL(url string) *WebviewPanel {
	// Empty means no navigation. Invalid URLs never reach the native webview.
	if url == "" {
		return p
	}
	normalized, err := assetserver.GetStartURL(url)
	if err != nil || normalized == "" {
		return p
	}
	p.destroyedLock.Lock()
	p.options.URL = normalized
	p.destroyedLock.Unlock()
	p.withImpl(func(impl webviewPanelImpl) { impl.setURL(normalized) })
	return p
}

// URL returns the current URL of the panel
func (p *WebviewPanel) URL() string {
	return p.snapshotOptions().URL
}

// Reload reloads the current page
func (p *WebviewPanel) Reload() {
	p.withImpl(func(impl webviewPanelImpl) { impl.reload() })
}

// ForceReload reloads the current page, bypassing the cache
func (p *WebviewPanel) ForceReload() {
	p.withImpl(func(impl webviewPanelImpl) { impl.forceReload() })
}

// Show makes the panel visible
func (p *WebviewPanel) Show() *WebviewPanel {
	visible := true
	p.destroyedLock.Lock()
	p.options.Visible = &visible
	p.destroyedLock.Unlock()
	p.withImpl(func(impl webviewPanelImpl) { impl.show() })
	return p
}

// Hide hides the panel
func (p *WebviewPanel) Hide() *WebviewPanel {
	visible := false
	p.destroyedLock.Lock()
	p.options.Visible = &visible
	p.destroyedLock.Unlock()
	p.withImpl(func(impl webviewPanelImpl) { impl.hide() })
	return p
}

// IsVisible returns whether the panel is currently visible
func (p *WebviewPanel) IsVisible() bool {
	if p.isDestroyed() {
		return false
	}
	options := p.snapshotOptions()
	visible := options.Visible != nil && *options.Visible
	p.withImpl(func(impl webviewPanelImpl) { visible = impl.isVisible() })
	return visible
}

// SetZoom sets the zoom level of the panel
func (p *WebviewPanel) SetZoom(zoom float64) *WebviewPanel {
	if zoom <= 0 || math.IsNaN(zoom) || math.IsInf(zoom, 0) {
		return p
	}
	p.destroyedLock.Lock()
	p.options.Zoom = zoom
	p.destroyedLock.Unlock()
	p.withImpl(func(impl webviewPanelImpl) { impl.setZoom(zoom) })
	return p
}

// GetZoom returns the current zoom level of the panel
func (p *WebviewPanel) GetZoom() float64 {
	zoom := p.snapshotOptions().Zoom
	p.withImpl(func(impl webviewPanelImpl) { zoom = impl.getZoom() })
	return zoom
}

// OpenDevTools opens the developer tools for this panel
func (p *WebviewPanel) OpenDevTools() {
	options := p.snapshotOptions()
	enabled := globalApplication != nil && globalApplication.isDebugMode
	if options.DevToolsEnabled != nil {
		enabled = *options.DevToolsEnabled
	}
	if enabled {
		p.withImpl(func(impl webviewPanelImpl) { impl.openDevTools() })
	}
}

// Focus gives focus to this panel
func (p *WebviewPanel) Focus() {
	p.withImpl(func(impl webviewPanelImpl) { impl.focus() })
}

// IsFocused returns whether this panel currently has focus
func (p *WebviewPanel) IsFocused() bool {
	focused := false
	p.withImpl(func(impl webviewPanelImpl) { focused = impl.isFocused() })
	return focused
}

// Destroy removes the panel from its parent window and releases resources
func (p *WebviewPanel) Destroy() {
	p.destroyedLock.Lock()
	if p.destroyed {
		p.destroyedLock.Unlock()
		return
	}
	p.destroyed = true
	impl := p.impl
	p.destroyedLock.Unlock()
	if impl != nil {
		InvokeSync(func() { p.destroyNative(impl) })
	}
	if p.parent != nil {
		p.parent.removePanel(p.id)
	}
}

// destroyNative runs only on the UI thread. Native creation may pump events,
// so cleanup waits for create to return. A queued Destroy callback may arrive
// after run has already cleaned up; both paths must release the view only once.
func (p *WebviewPanel) destroyNative(impl webviewPanelImpl) {
	if p.creating || p.nativeDestroyed {
		return
	}
	p.nativeDestroyed = true
	impl.destroy()
}

// isDestroyed returns whether the panel has been destroyed
func (p *WebviewPanel) isDestroyed() bool {
	p.destroyedLock.RLock()
	defer p.destroyedLock.RUnlock()
	return p.destroyed
}

// run initializes the platform-specific implementation
// This is called by the parent window when the panel is added
func (p *WebviewPanel) run() {
	if p.parent == nil || p.isDestroyed() {
		return
	}
	InvokeSync(func() {
		p.destroyedLock.Lock()
		if p.impl != nil || p.destroyed || p.parent.isDestroyed() || p.parent.impl == nil {
			p.destroyedLock.Unlock()
			return
		}
		impl := newPanelImpl(p)
		p.impl = impl
		p.destroyedLock.Unlock()
		if impl == nil {
			return
		}
		p.initializeAnchor()
		initialURL := p.snapshotOptions().URL
		p.creating = true
		impl.create()
		p.creating = false
		if p.isDestroyed() {
			p.destroyNative(impl)
			return
		}
		// Changes can arrive while a native create call pumps its event loop.
		// Reapply the latest configuration before exposing the completed view.
		options := p.snapshotOptions()
		impl.setBounds(Rect{X: options.X, Y: options.Y, Width: options.Width, Height: options.Height})
		impl.setZoom(options.Zoom)
		if *options.Visible {
			impl.show()
		} else {
			impl.hide()
		}
		if options.URL != initialURL && options.URL != "" {
			impl.setURL(options.URL)
		}
		impl.setZIndex(options.ZIndex)
	})
}

// =========================================================================
// Layout Helper Methods
// =========================================================================

// FillWindow makes the panel fill the entire parent window.
// This is a convenience method equivalent to setting position to (0,0)
// and size to the window's content size.
func (p *WebviewPanel) FillWindow() *WebviewPanel {
	if p.parent == nil {
		return p
	}
	width, height := p.parentSize()
	return p.SetBounds(Rect{X: 0, Y: 0, Width: width, Height: height})
}

// DockLeft positions the panel on the left side of the window with the specified width.
// Height fills the window. Useful for sidebars and navigation panels.
func (p *WebviewPanel) DockLeft(width int) *WebviewPanel {
	if p.parent == nil {
		return p
	}
	_, height := p.parentSize()
	return p.SetBounds(Rect{X: 0, Y: 0, Width: width, Height: height})
}

// DockRight positions the panel on the right side of the window with the specified width.
// Height fills the window. Useful for property panels and inspectors.
func (p *WebviewPanel) DockRight(width int) *WebviewPanel {
	if p.parent == nil {
		return p
	}
	windowWidth, height := p.parentSize()
	return p.SetBounds(Rect{X: windowWidth - width, Y: 0, Width: width, Height: height})
}

// DockTop positions the panel at the top of the window with the specified height.
// Width fills the window. Useful for toolbars and header areas.
func (p *WebviewPanel) DockTop(height int) *WebviewPanel {
	if p.parent == nil {
		return p
	}
	width, _ := p.parentSize()
	return p.SetBounds(Rect{X: 0, Y: 0, Width: width, Height: height})
}

// DockBottom positions the panel at the bottom of the window with the specified height.
// Width fills the window. Useful for status bars and terminal panels.
func (p *WebviewPanel) DockBottom(height int) *WebviewPanel {
	if p.parent == nil {
		return p
	}
	width, windowHeight := p.parentSize()
	return p.SetBounds(Rect{X: 0, Y: windowHeight - height, Width: width, Height: height})
}

// FillBeside fills the remaining space beside another panel.
// The direction specifies whether to fill to the right, left, above, or below the reference panel.
func (p *WebviewPanel) FillBeside(refPanel *WebviewPanel, direction string) *WebviewPanel {
	if p.parent == nil || refPanel == nil || refPanel == p || refPanel.parent != p.parent {
		return p
	}

	windowWidth, windowHeight := p.parentSize()
	refBounds := refPanel.Bounds()

	var bounds Rect
	switch direction {
	case "right":
		bounds = Rect{
			X:      refBounds.X + refBounds.Width,
			Y:      refBounds.Y,
			Width:  windowWidth - (refBounds.X + refBounds.Width),
			Height: refBounds.Height,
		}
	case "left":
		bounds = Rect{
			X:      0,
			Y:      refBounds.Y,
			Width:  refBounds.X,
			Height: refBounds.Height,
		}
	case "below":
		bounds = Rect{
			X:      refBounds.X,
			Y:      refBounds.Y + refBounds.Height,
			Width:  refBounds.Width,
			Height: windowHeight - (refBounds.Y + refBounds.Height),
		}
	case "above":
		bounds = Rect{
			X:      refBounds.X,
			Y:      0,
			Width:  refBounds.Width,
			Height: refBounds.Y,
		}
	default:
		return p
	}

	return p.SetBounds(bounds)
}

// =========================================================================
// Anchor/Responsive Layout Methods
// =========================================================================

// initializeAnchor stores the original window size for anchor calculations.
// This is called when the panel is first attached to a window.
func (p *WebviewPanel) initializeAnchor() {
	if p.parent == nil {
		return
	}
	width, height := p.parentSize()
	p.destroyedLock.Lock()
	p.originalWindowWidth, p.originalWindowHeight = width, height
	p.destroyedLock.Unlock()
}

// handleWindowResize recalculates the panel's bounds based on its anchor settings.
// This is called automatically when the parent window is resized.
func (p *WebviewPanel) handleWindowResize(newWindowWidth, newWindowHeight int) {
	if p.isDestroyed() {
		return
	}
	options := p.snapshotOptions()
	if options.Anchor == AnchorNone {
		// A DPI change still requires rescaling an absolutely positioned child.
		bounds := Rect{X: options.X, Y: options.Y, Width: options.Width, Height: options.Height}
		p.withImpl(func(impl webviewPanelImpl) { impl.setBounds(bounds) })
		return
	}

	newBounds := p.calculateAnchoredBounds(newWindowWidth, newWindowHeight)
	// Use internal setBounds to avoid updating anchor baseline during resize
	p.setBoundsInternal(newBounds)
}

// setBoundsInternal sets bounds without updating anchor baseline.
// Used internally during window resize handling.
func (p *WebviewPanel) setBoundsInternal(bounds Rect) {
	p.destroyedLock.Lock()
	p.options.X, p.options.Y = bounds.X, bounds.Y
	p.options.Width, p.options.Height = bounds.Width, bounds.Height
	p.destroyedLock.Unlock()
	p.withImpl(func(impl webviewPanelImpl) { impl.setBounds(bounds) })
}

// calculateAnchoredBounds computes the new bounds based on anchor settings.
func (p *WebviewPanel) calculateAnchoredBounds(newWindowWidth, newWindowHeight int) Rect {
	p.destroyedLock.RLock()
	defer p.destroyedLock.RUnlock()
	anchor := p.options.Anchor
	orig := p.originalBounds
	origWinW := p.originalWindowWidth
	origWinH := p.originalWindowHeight

	// If original window size was not recorded, use current bounds
	if origWinW == 0 || origWinH == 0 {
		return Rect{
			X:      p.options.X,
			Y:      p.options.Y,
			Width:  p.options.Width,
			Height: p.options.Height,
		}
	}

	// Calculate distances from edges
	distanceFromRight := origWinW - (orig.X + orig.Width)
	distanceFromBottom := origWinH - (orig.Y + orig.Height)

	newX := orig.X
	newY := orig.Y
	newWidth := orig.Width
	newHeight := orig.Height

	// Handle horizontal anchoring
	hasLeft := anchor.HasAnchor(AnchorLeft)
	hasRight := anchor.HasAnchor(AnchorRight)

	if hasLeft && hasRight {
		// Anchored to both sides - stretch horizontally
		newX = orig.X
		newWidth = newWindowWidth - orig.X - distanceFromRight
	} else if hasRight {
		// Anchored to right only - maintain distance from right
		newX = newWindowWidth - distanceFromRight - orig.Width
	}
	// If hasLeft only or no horizontal anchor, X stays the same

	// Handle vertical anchoring
	hasTop := anchor.HasAnchor(AnchorTop)
	hasBottom := anchor.HasAnchor(AnchorBottom)

	if hasTop && hasBottom {
		// Anchored to both sides - stretch vertically
		newY = orig.Y
		newHeight = newWindowHeight - orig.Y - distanceFromBottom
	} else if hasBottom {
		// Anchored to bottom only - maintain distance from bottom
		newY = newWindowHeight - distanceFromBottom - orig.Height
	}
	// If hasTop only or no vertical anchor, Y stays the same

	// Ensure minimum dimensions
	if newWidth < 1 {
		newWidth = 1
	}
	if newHeight < 1 {
		newHeight = 1
	}

	return Rect{
		X:      newX,
		Y:      newY,
		Width:  newWidth,
		Height: newHeight,
	}
}

// snapshotOptions returns a consistent view of the panel configuration.
// Reference-valued options are copied at construction and never mutated in place.
func (p *WebviewPanel) snapshotOptions() WebviewPanelOptions {
	p.destroyedLock.RLock()
	defer p.destroyedLock.RUnlock()
	return p.options
}

// withImpl checks lifecycle state again on the UI thread, so queued operations
// cannot access a native view that has already been destroyed.
func (p *WebviewPanel) withImpl(fn func(webviewPanelImpl)) {
	p.destroyedLock.RLock()
	ready := p.impl != nil && !p.destroyed
	p.destroyedLock.RUnlock()
	if !ready {
		return
	}
	InvokeSync(func() {
		p.destroyedLock.RLock()
		impl, destroyed := p.impl, p.destroyed
		p.destroyedLock.RUnlock()
		if impl != nil && !destroyed && !p.creating {
			fn(impl)
		}
	})
}

// parentSize uses the configured content size before the native window exists.
func (p *WebviewPanel) parentSize() (int, int) {
	if p.parent.impl == nil {
		return p.parent.options.Width, p.parent.options.Height
	}
	return p.parent.Size()
}

// sortedSiblings snapshots stacking keys before sorting so concurrent setters
// cannot change the comparator midway through a native reordering pass.
func (p *WebviewPanel) sortedSiblings() []*WebviewPanel {
	panels := p.parent.GetPanels()
	keys := make(map[uint]int, len(panels))
	for _, panel := range panels {
		keys[panel.id] = panel.ZIndex()
	}
	sort.Slice(panels, func(i, j int) bool {
		if keys[panels[i].id] == keys[panels[j].id] {
			return panels[i].id < panels[j].id
		}
		return keys[panels[i].id] < keys[panels[j].id]
	})
	return panels
}
