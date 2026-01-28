package application

import (
	"fmt"
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
	destroyed     bool
	destroyedLock sync.RWMutex

	// Original window size when panel was created (for anchor calculations)
	originalWindowWidth  int
	originalWindowHeight int
	// Original panel bounds (for anchor calculations)
	originalBounds Rect
}

// NewPanel creates a new WebviewPanel with the given options.
// The panel must be associated with a parent window via window.AddPanel().
func NewPanel(options WebviewPanelOptions) *WebviewPanel {
	id := getNextPanelID()

	// Apply defaults
	if options.Width == 0 {
		options.Width = 400
	}
	if options.Height == 0 {
		options.Height = 300
	}
	if options.ZIndex == 0 {
		options.ZIndex = 1
	}
	if options.Zoom == 0 {
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

	// Normalize URL via asset server for local paths
	if options.URL != "" {
		normalizedURL, _ := assetserver.GetStartURL(options.URL)
		options.URL = normalizedURL
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
	p.options.X = bounds.X
	p.options.Y = bounds.Y
	p.options.Width = bounds.Width
	p.options.Height = bounds.Height

	// Update anchor baseline so future resizes calculate from the new position
	p.updateAnchorBaseline(bounds)

	if p.impl != nil && !p.isDestroyed() {
		InvokeSync(func() {
			p.impl.setBounds(bounds)
		})
	}
	return p
}

// updateAnchorBaseline updates the original bounds and window size used for anchor calculations.
// Called when the user manually changes panel bounds.
func (p *WebviewPanel) updateAnchorBaseline(bounds Rect) {
	p.originalBounds = bounds
	if p.parent != nil {
		p.originalWindowWidth, p.originalWindowHeight = p.parent.Size()
	}
}

// Bounds returns the current bounds of the panel
func (p *WebviewPanel) Bounds() Rect {
	if p.impl != nil && !p.isDestroyed() {
		return InvokeSyncWithResult(p.impl.bounds)
	}
	return Rect{
		X:      p.options.X,
		Y:      p.options.Y,
		Width:  p.options.Width,
		Height: p.options.Height,
	}
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
	p.options.ZIndex = zIndex
	if p.impl != nil && !p.isDestroyed() {
		InvokeSync(func() {
			p.impl.setZIndex(zIndex)
		})
	}
	return p
}

// ZIndex returns the current z-index of the panel
func (p *WebviewPanel) ZIndex() int {
	return p.options.ZIndex
}

// SetURL navigates the panel to the specified URL
// Local paths (e.g., "/panel.html") are normalized via the asset server.
func (p *WebviewPanel) SetURL(url string) *WebviewPanel {
	// Normalize URL via asset server for local paths
	normalizedURL, _ := assetserver.GetStartURL(url)
	p.options.URL = normalizedURL
	if p.impl != nil && !p.isDestroyed() {
		InvokeSync(func() {
			p.impl.setURL(normalizedURL)
		})
	}
	return p
}

// URL returns the current URL of the panel
func (p *WebviewPanel) URL() string {
	return p.options.URL
}

// Reload reloads the current page
func (p *WebviewPanel) Reload() {
	if p.impl != nil && !p.isDestroyed() {
		InvokeSync(p.impl.reload)
	}
}

// ForceReload reloads the current page, bypassing the cache
func (p *WebviewPanel) ForceReload() {
	if p.impl != nil && !p.isDestroyed() {
		InvokeSync(p.impl.forceReload)
	}
}

// Show makes the panel visible
func (p *WebviewPanel) Show() *WebviewPanel {
	visible := true
	p.options.Visible = &visible
	if p.impl != nil && !p.isDestroyed() {
		InvokeSync(p.impl.show)
	}
	return p
}

// Hide hides the panel
func (p *WebviewPanel) Hide() *WebviewPanel {
	visible := false
	p.options.Visible = &visible
	if p.impl != nil && !p.isDestroyed() {
		InvokeSync(p.impl.hide)
	}
	return p
}

// IsVisible returns whether the panel is currently visible
func (p *WebviewPanel) IsVisible() bool {
	if p.impl != nil && !p.isDestroyed() {
		return InvokeSyncWithResult(p.impl.isVisible)
	}
	return p.options.Visible != nil && *p.options.Visible
}

// SetZoom sets the zoom level of the panel
func (p *WebviewPanel) SetZoom(zoom float64) *WebviewPanel {
	p.options.Zoom = zoom
	if p.impl != nil && !p.isDestroyed() {
		InvokeSync(func() {
			p.impl.setZoom(zoom)
		})
	}
	return p
}

// GetZoom returns the current zoom level of the panel
func (p *WebviewPanel) GetZoom() float64 {
	if p.impl != nil && !p.isDestroyed() {
		return InvokeSyncWithResult(p.impl.getZoom)
	}
	return p.options.Zoom
}

// OpenDevTools opens the developer tools for this panel
func (p *WebviewPanel) OpenDevTools() {
	if p.impl != nil && !p.isDestroyed() {
		InvokeSync(p.impl.openDevTools)
	}
}

// Focus gives focus to this panel
func (p *WebviewPanel) Focus() {
	if p.impl != nil && !p.isDestroyed() {
		InvokeSync(p.impl.focus)
	}
}

// IsFocused returns whether this panel currently has focus
func (p *WebviewPanel) IsFocused() bool {
	if p.impl != nil && !p.isDestroyed() {
		return InvokeSyncWithResult(p.impl.isFocused)
	}
	return false
}

// Destroy removes the panel from its parent window and releases resources
func (p *WebviewPanel) Destroy() {
	if p.isDestroyed() {
		return
	}

	p.destroyedLock.Lock()
	p.destroyed = true
	p.destroyedLock.Unlock()

	if p.impl != nil {
		InvokeSync(p.impl.destroy)
	}

	// Remove from parent
	if p.parent != nil {
		p.parent.removePanel(p.id)
	}
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
	globalApplication.debug("[Panel] run() called", "panelID", p.id, "panelName", p.name)

	p.destroyedLock.Lock()
	if p.impl != nil || p.destroyed {
		globalApplication.debug("[Panel] run() skipped - impl already exists or destroyed",
			"panelID", p.id, "hasImpl", p.impl != nil, "destroyed", p.destroyed)
		p.destroyedLock.Unlock()
		return
	}

	// Check parent window state before creating impl
	if p.parent == nil {
		globalApplication.error("[Panel] run() failed - parent window is nil", "panelID", p.id)
		p.destroyedLock.Unlock()
		return
	}
	if p.parent.impl == nil {
		globalApplication.error("[Panel] run() failed - parent window impl is nil", "panelID", p.id, "windowID", p.parent.id)
		p.destroyedLock.Unlock()
		return
	}

	globalApplication.debug("[Panel] Creating platform impl", "panelID", p.id, "parentWindowID", p.parent.id)
	p.impl = newPanelImpl(p)
	p.destroyedLock.Unlock()

	if p.impl == nil {
		globalApplication.error("[Panel] newPanelImpl returned nil", "panelID", p.id)
		return
	}

	globalApplication.debug("[Panel] Calling impl.create()", "panelID", p.id)
	InvokeSync(p.impl.create)
	globalApplication.debug("[Panel] impl.create() completed", "panelID", p.id)
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
	width, height := p.parent.Size()
	return p.SetBounds(Rect{X: 0, Y: 0, Width: width, Height: height})
}

// DockLeft positions the panel on the left side of the window with the specified width.
// Height fills the window. Useful for sidebars and navigation panels.
func (p *WebviewPanel) DockLeft(width int) *WebviewPanel {
	if p.parent == nil {
		return p
	}
	_, height := p.parent.Size()
	return p.SetBounds(Rect{X: 0, Y: 0, Width: width, Height: height})
}

// DockRight positions the panel on the right side of the window with the specified width.
// Height fills the window. Useful for property panels and inspectors.
func (p *WebviewPanel) DockRight(width int) *WebviewPanel {
	if p.parent == nil {
		return p
	}
	windowWidth, height := p.parent.Size()
	return p.SetBounds(Rect{X: windowWidth - width, Y: 0, Width: width, Height: height})
}

// DockTop positions the panel at the top of the window with the specified height.
// Width fills the window. Useful for toolbars and header areas.
func (p *WebviewPanel) DockTop(height int) *WebviewPanel {
	if p.parent == nil {
		return p
	}
	width, _ := p.parent.Size()
	return p.SetBounds(Rect{X: 0, Y: 0, Width: width, Height: height})
}

// DockBottom positions the panel at the bottom of the window with the specified height.
// Width fills the window. Useful for status bars and terminal panels.
func (p *WebviewPanel) DockBottom(height int) *WebviewPanel {
	if p.parent == nil {
		return p
	}
	width, windowHeight := p.parent.Size()
	return p.SetBounds(Rect{X: 0, Y: windowHeight - height, Width: width, Height: height})
}

// FillBeside fills the remaining space beside another panel.
// The direction specifies whether to fill to the right, left, above, or below the reference panel.
func (p *WebviewPanel) FillBeside(refPanel *WebviewPanel, direction string) *WebviewPanel {
	if p.parent == nil || refPanel == nil {
		return p
	}

	windowWidth, windowHeight := p.parent.Size()
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
	p.originalWindowWidth, p.originalWindowHeight = p.parent.Size()
}

// handleWindowResize recalculates the panel's bounds based on its anchor settings.
// This is called automatically when the parent window is resized.
func (p *WebviewPanel) handleWindowResize(newWindowWidth, newWindowHeight int) {
	if p.isDestroyed() || p.options.Anchor == AnchorNone {
		return
	}

	newBounds := p.calculateAnchoredBounds(newWindowWidth, newWindowHeight)
	// Use internal setBounds to avoid updating anchor baseline during resize
	p.setBoundsInternal(newBounds)
}

// setBoundsInternal sets bounds without updating anchor baseline.
// Used internally during window resize handling.
func (p *WebviewPanel) setBoundsInternal(bounds Rect) {
	p.options.X = bounds.X
	p.options.Y = bounds.Y
	p.options.Width = bounds.Width
	p.options.Height = bounds.Height

	if p.impl != nil && !p.isDestroyed() {
		InvokeSync(func() {
			p.impl.setBounds(bounds)
		})
	}
}

// calculateAnchoredBounds computes the new bounds based on anchor settings.
func (p *WebviewPanel) calculateAnchoredBounds(newWindowWidth, newWindowHeight int) Rect {
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
