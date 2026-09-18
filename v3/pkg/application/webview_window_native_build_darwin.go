//go:build darwin && !ios && !server && wails_native

package application

import (
	"errors"
	"unsafe"
)

// macosWebviewWindow is a compile-time-disabled placeholder in a native build.
// Keeping the private implementation shape lets the additive v3 public API
// compile, while all actual WKWebView code and Objective-C sources are absent.
type macosWebviewWindow struct {
	parent          *WebviewWindow
	nsWindow        unsafe.Pointer
	activeToolbar   *MacToolbar
	activeSplitView *MacSplitView
}

var errWebviewUnavailableInNativeBuild = errors.New("WebviewWindow is unavailable in a wails_native build")

func newWindowImpl(parent *WebviewWindow) *macosWebviewWindow {
	return &macosWebviewWindow{parent: parent}
}

func platformTitlebarDoubleClickPreference() string { return "" }

func (w *macosWebviewWindow) reportUnavailable() {
	if w != nil && w.parent != nil {
		w.parent.Error("%s", errWebviewUnavailableInNativeBuild)
	}
}

func (w *macosWebviewWindow) run()                   { w.reportUnavailable() }
func (*macosWebviewWindow) setTitle(string)          {}
func (*macosWebviewWindow) setSize(int, int)         {}
func (*macosWebviewWindow) setAlwaysOnTop(bool)      {}
func (*macosWebviewWindow) setURL(string)            {}
func (*macosWebviewWindow) setResizable(bool)        {}
func (*macosWebviewWindow) setMinSize(int, int)      {}
func (*macosWebviewWindow) setMaxSize(int, int)      {}
func (*macosWebviewWindow) execJS(string)            {}
func (*macosWebviewWindow) execJSDragOver([]byte)    {}
func (*macosWebviewWindow) setBackgroundColour(RGBA) {}
func (*macosWebviewWindow) center()                  {}
func (*macosWebviewWindow) size() (int, int)         { return 0, 0 }
func (*macosWebviewWindow) width() int               { return 0 }
func (*macosWebviewWindow) height() int              { return 0 }
func (*macosWebviewWindow) destroy()                 {}
func (*macosWebviewWindow) reload()                  {}
func (*macosWebviewWindow) forceReload()             {}
func (*macosWebviewWindow) openDevTools()            {}
func (*macosWebviewWindow) zoomReset()               {}
func (*macosWebviewWindow) zoomIn()                  {}
func (*macosWebviewWindow) zoomOut()                 {}
func (*macosWebviewWindow) getZoom() float64         { return 1 }
func (*macosWebviewWindow) setZoom(float64)          {}
func (*macosWebviewWindow) close()                   {}
func (*macosWebviewWindow) zoom()                    {}
func (*macosWebviewWindow) setHTML(string)           {}
func (*macosWebviewWindow) on(uint)                  {}
func (*macosWebviewWindow) minimise()                {}
func (*macosWebviewWindow) unminimise()              {}
func (*macosWebviewWindow) maximise()                {}
func (*macosWebviewWindow) unmaximise()              {}
func (*macosWebviewWindow) fullscreen()              {}
func (*macosWebviewWindow) unfullscreen()            {}
func (*macosWebviewWindow) isMinimised() bool        { return false }
func (*macosWebviewWindow) isMaximised() bool        { return false }
func (*macosWebviewWindow) isFullscreen() bool       { return false }
func (*macosWebviewWindow) isNormal() bool           { return true }
func (*macosWebviewWindow) isVisible() bool          { return false }
func (*macosWebviewWindow) isFocused() bool          { return false }
func (*macosWebviewWindow) focus()                   {}
func (*macosWebviewWindow) show()                    {}
func (*macosWebviewWindow) hide()                    {}
func (*macosWebviewWindow) getScreen() (*Screen, error) {
	return nil, errWebviewUnavailableInNativeBuild
}
func (*macosWebviewWindow) setFrameless(bool)                                   {}
func (*macosWebviewWindow) openContextMenu(*Menu, *ContextMenuData)             {}
func (*macosWebviewWindow) nativeWindow() unsafe.Pointer                        { return nil }
func (*macosWebviewWindow) startDrag() error                                    { return errWebviewUnavailableInNativeBuild }
func (*macosWebviewWindow) startResize(string) error                            { return errWebviewUnavailableInNativeBuild }
func (*macosWebviewWindow) print() error                                        { return errWebviewUnavailableInNativeBuild }
func (*macosWebviewWindow) setEnabled(bool)                                     {}
func (*macosWebviewWindow) physicalBounds() Rect                                { return Rect{} }
func (*macosWebviewWindow) setPhysicalBounds(Rect)                              {}
func (*macosWebviewWindow) bounds() Rect                                        { return Rect{} }
func (*macosWebviewWindow) setBounds(Rect)                                      {}
func (*macosWebviewWindow) position() (int, int)                                { return 0, 0 }
func (*macosWebviewWindow) setPosition(int, int)                                {}
func (*macosWebviewWindow) centerOnScreen(*Screen)                              {}
func (*macosWebviewWindow) relativePosition() (int, int)                        { return 0, 0 }
func (*macosWebviewWindow) setRelativePosition(int, int)                        {}
func (*macosWebviewWindow) flash(bool)                                          {}
func (*macosWebviewWindow) handleKeyEvent(string)                               {}
func (*macosWebviewWindow) getBorderSizes() *LRTB                               { return &LRTB{} }
func (*macosWebviewWindow) setMinimiseButtonState(ButtonState)                  {}
func (*macosWebviewWindow) setMaximiseButtonState(ButtonState)                  {}
func (*macosWebviewWindow) setCloseButtonState(ButtonState)                     {}
func (*macosWebviewWindow) setFullscreenButtonState(ButtonState)                {}
func (*macosWebviewWindow) isIgnoreMouseEvents() bool                           { return false }
func (*macosWebviewWindow) setIgnoreMouseEvents(bool)                           {}
func (*macosWebviewWindow) cut()                                                {}
func (*macosWebviewWindow) copy()                                               {}
func (*macosWebviewWindow) paste()                                              {}
func (*macosWebviewWindow) undo()                                               {}
func (*macosWebviewWindow) delete()                                             {}
func (*macosWebviewWindow) selectAll()                                          {}
func (*macosWebviewWindow) redo()                                               {}
func (*macosWebviewWindow) showMenuBar()                                        {}
func (*macosWebviewWindow) hideMenuBar()                                        {}
func (*macosWebviewWindow) toggleMenuBar()                                      {}
func (*macosWebviewWindow) setMenu(*Menu)                                       {}
func (*macosWebviewWindow) snapAssist()                                         {}
func (*macosWebviewWindow) setContentProtection(bool)                           {}
func (*macosWebviewWindow) attachModal(*WebviewWindow)                          {}
func (*macosWebviewWindow) setNonClientHitTestRegions([]nonClientHitTestRegion) {}
