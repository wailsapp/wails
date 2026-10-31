//go:build linux && gtk3 && !android && !server

package application

func (w *linuxWebviewWindow) zoomIn() {
	w.setZoom(w.getZoom() * zoomInFactor)
}

func (w *linuxWebviewWindow) zoomOut() {
	w.setZoom(w.getZoom() * (1.0 / zoomInFactor))
}

// unregisterMenuItemAccelerator does nothing on GTK3: each item's accelerator
// is a closure on the item widget, which GTK disconnects when the menu bar
// holding the item is destroyed (windowDetachMenuBar).
func unregisterMenuItemAccelerator(_ uint) {}

func setWindowMenuActions(_ pointer, _ bool) {}

// NewCalloc creates a new allocator. Returns a pointer so the allocation
// pool is shared across calls without copying.
func NewCalloc() *Calloc {
	return &Calloc{}
}
