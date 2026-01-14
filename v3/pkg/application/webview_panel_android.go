//go:build android

package application

// Android stub implementation for WebviewPanel
// Panels are not yet supported on Android

type androidPanelImpl struct {
	panel *WebviewPanel
}

func newPanelImpl(panel *WebviewPanel) webviewPanelImpl {
	return &androidPanelImpl{panel: panel}
}

func (p *androidPanelImpl) create()               {}
func (p *androidPanelImpl) destroy()              {}
func (p *androidPanelImpl) setBounds(bounds Rect) {}
func (p *androidPanelImpl) bounds() Rect          { return Rect{} }
func (p *androidPanelImpl) setZIndex(zIndex int)  {}
func (p *androidPanelImpl) setURL(url string)     {}
func (p *androidPanelImpl) setHTML(html string)   {}
func (p *androidPanelImpl) execJS(js string)      {}
func (p *androidPanelImpl) reload()               {}
func (p *androidPanelImpl) forceReload()          {}
func (p *androidPanelImpl) show()                 {}
func (p *androidPanelImpl) hide()                 {}
func (p *androidPanelImpl) isVisible() bool       { return false }
func (p *androidPanelImpl) setZoom(zoom float64)  {}
func (p *androidPanelImpl) getZoom() float64      { return 1.0 }
func (p *androidPanelImpl) openDevTools()         {}
func (p *androidPanelImpl) focus()                {}
func (p *androidPanelImpl) isFocused() bool       { return false }
