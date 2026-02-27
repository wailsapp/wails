//go:build ios

package application

// iOS stub implementation for WebviewPanel
// Panels are not yet supported on iOS.
// All methods are no-ops until iOS platform support is implemented.

type iosPanelImpl struct {
	panel *WebviewPanel
}

func newPanelImpl(panel *WebviewPanel) webviewPanelImpl {
	return &iosPanelImpl{panel: panel}
}

func (p *iosPanelImpl) create() {
	// Not implemented on iOS
}

func (p *iosPanelImpl) destroy() {
	// Not implemented on iOS
}

func (p *iosPanelImpl) setBounds(_ Rect) {
	// Not implemented on iOS
}

func (p *iosPanelImpl) bounds() Rect {
	return Rect{}
}

func (p *iosPanelImpl) setZIndex(_ int) {
	// Not implemented on iOS
}

func (p *iosPanelImpl) setURL(_ string) {
	// Not implemented on iOS
}

func (p *iosPanelImpl) reload() {
	// Not implemented on iOS
}

func (p *iosPanelImpl) forceReload() {
	// Not implemented on iOS
}

func (p *iosPanelImpl) show() {
	// Not implemented on iOS
}

func (p *iosPanelImpl) hide() {
	// Not implemented on iOS
}

func (p *iosPanelImpl) isVisible() bool {
	return false
}

func (p *iosPanelImpl) setZoom(_ float64) {
	// Not implemented on iOS
}

func (p *iosPanelImpl) getZoom() float64 {
	return 1.0
}

func (p *iosPanelImpl) openDevTools() {
	// Not implemented on iOS
}

func (p *iosPanelImpl) focus() {
	// Not implemented on iOS
}

func (p *iosPanelImpl) isFocused() bool {
	return false
}
