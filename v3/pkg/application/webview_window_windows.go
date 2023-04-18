//go:build windows

package application

import (
	"unsafe"
)

var showDevTools = func(window unsafe.Pointer) {}

type windowsWebviewWindow struct {
	windowImpl unsafe.Pointer
	parent     *WebviewWindow
}

func (w *windowsWebviewWindow) setTitle(title string) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setSize(width, height int) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setAlwaysOnTop(alwaysOnTop bool) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setURL(url string) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setResizable(resizable bool) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setMinSize(width, height int) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setMaxSize(width, height int) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) execJS(js string) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) restore() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setBackgroundColour(color *RGBA) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) run() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) center() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) size() (int, int) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) width() int {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) height() int {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) position() (int, int) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) destroy() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) reload() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) forceReload() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) toggleDevTools() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) zoomReset() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) zoomIn() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) zoomOut() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) getZoom() float64 {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setZoom(zoom float64) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) close() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) zoom() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setHTML(html string) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setPosition(x int, y int) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) on(eventID uint) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) minimise() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) unminimise() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) maximise() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) unmaximise() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) fullscreen() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) unfullscreen() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) isMinimised() bool {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) isMaximised() bool {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) isFullscreen() bool {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) disableSizeConstraints() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setFullscreenButtonEnabled(enabled bool) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) show() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) hide() {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) getScreen() (*Screen, error) {
	//TODO implement me
	panic("implement me")
}

func (w *windowsWebviewWindow) setFrameless(b bool) {
	//TODO implement me
	panic("implement me")
}

func newWindowImpl(parent *WebviewWindow) *windowsWebviewWindow {
	result := &windowsWebviewWindow{
		parent: parent,
	}
	return result
}

func (w *windowsWebviewWindow) openContextMenu(menu *Menu, data *ContextMenuData) {
	// Create the menu
	thisMenu := newMenuImpl(menu)
	thisMenu.update()
	//C.windowShowMenu(w.nsWindow, thisMenu.nsMenu, C.int(data.X), C.int(data.Y))
}
