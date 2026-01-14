//go:build android

package application

// Android stub implementation for WebviewPanel
// Panels are not yet supported on Android.
// All methods are no-ops until Android platform support is implemented.

type androidPanelImpl struct {
	panel *WebviewPanel
}

func newPanelImpl(panel *WebviewPanel) webviewPanelImpl {
	return &androidPanelImpl{panel: panel}
}

func (p *androidPanelImpl) create() {
	// Not implemented on Android
}

func (p *androidPanelImpl) destroy() {
	// Not implemented on Android
}

func (p *androidPanelImpl) setBounds(_ Rect) {
	// Not implemented on Android
}

func (p *androidPanelImpl) bounds() Rect {
	return Rect{}
}

func (p *androidPanelImpl) setZIndex(_ int) {
	// Not implemented on Android
}

func (p *androidPanelImpl) setURL(_ string) {
	// Not implemented on Android
}

func (p *androidPanelImpl) setHTML(_ string) {
	// Not implemented on Android
}

func (p *androidPanelImpl) execJS(_ string) {
	// Not implemented on Android
}

func (p *androidPanelImpl) reload() {
	// Not implemented on Android
}

func (p *androidPanelImpl) forceReload() {
	// Not implemented on Android
}

func (p *androidPanelImpl) show() {
	// Not implemented on Android
}

func (p *androidPanelImpl) hide() {
	// Not implemented on Android
}

func (p *androidPanelImpl) isVisible() bool {
	return false
}

func (p *androidPanelImpl) setZoom(_ float64) {
	// Not implemented on Android
}

func (p *androidPanelImpl) getZoom() float64 {
	return 1.0
}

func (p *androidPanelImpl) openDevTools() {
	// Not implemented on Android
}

func (p *androidPanelImpl) focus() {
	// Not implemented on Android
}

func (p *androidPanelImpl) isFocused() bool {
	return false
}
