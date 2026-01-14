//go:build ios

package application

// iOS stub implementation for WebviewPanel
// Panels are not yet supported on iOS

type iosPanelImpl struct {
	panel *WebviewPanel
}

func newPanelImpl(panel *WebviewPanel) webviewPanelImpl {
	return &iosPanelImpl{panel: panel}
}

func (p *iosPanelImpl) create()               {}
func (p *iosPanelImpl) destroy()              {}
func (p *iosPanelImpl) setBounds(bounds Rect) {}
func (p *iosPanelImpl) bounds() Rect          { return Rect{} }
func (p *iosPanelImpl) setZIndex(zIndex int)  {}
func (p *iosPanelImpl) setURL(url string)     {}
func (p *iosPanelImpl) setHTML(html string)   {}
func (p *iosPanelImpl) execJS(js string)      {}
func (p *iosPanelImpl) reload()               {}
func (p *iosPanelImpl) forceReload()          {}
func (p *iosPanelImpl) show()                 {}
func (p *iosPanelImpl) hide()                 {}
func (p *iosPanelImpl) isVisible() bool       { return false }
func (p *iosPanelImpl) setZoom(zoom float64)  {}
func (p *iosPanelImpl) getZoom() float64      { return 1.0 }
func (p *iosPanelImpl) openDevTools()         {}
func (p *iosPanelImpl) focus()                {}
func (p *iosPanelImpl) isFocused() bool       { return false }
