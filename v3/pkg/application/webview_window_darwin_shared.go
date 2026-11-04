//go:build darwin && !ios && !server

package application

import (
	"time"
	"unsafe"
)

type macosWebviewWindow struct {
	nsWindow unsafe.Pointer // NSWindow*
	parent   *WebviewWindow

	// macosWindowNative is the backend's own state: empty for cgo, which keeps
	// it on the Objective-C side; the purego backend holds it here.
	macosWindowNative
}

// applyWindowButtonStates applies the title-bar button states, hiding all three
// buttons on a native-frame frameless window.
func (w *macosWebviewWindow) applyWindowButtonStates() {
	states := effectiveMacWindowButtonStates(w.parent.options)
	w.setMinimiseButtonState(states.minimise)
	w.setCloseButtonState(states.close)
	w.setMaximiseButtonState(states.zoom)
}

func (w *macosWebviewWindow) handleKeyEvent(acceleratorString string) {
	// Parse acceleratorString
	accelerator, err := parseAccelerator(acceleratorString)
	if err != nil {
		globalApplication.error("unable to parse accelerator: %w", err)
		return
	}
	w.parent.processKeyBinding(accelerator.String())
}

func (w *macosWebviewWindow) hideMenuBar() {}

func (w *macosWebviewWindow) isNormal() bool {
	return !w.isMinimised() && !w.isMaximised() && !w.isFullscreen()
}

// notchAnimation reports whether the window animates, and the show duration in
// seconds.
func (w *macosWebviewWindow) notchAnimation() (bool, float64) {
	if w.parent.notchWindow == nil {
		return false, 0
	}
	return w.parent.notchWindow.animated,
		float64(w.parent.notchWindow.animationSpeed) / float64(time.Second)
}

func (w *macosWebviewWindow) setMenu(_ *Menu) {}

func (w *macosWebviewWindow) showMenuBar() {}

func (w *macosWebviewWindow) snapAssist() {}

func (w *macosWebviewWindow) toggleMenuBar() {}

// checkUnconditionalClose reports whether the window may close without asking.
func checkUnconditionalClose(windowId uint) bool {
	window, _ := globalApplication.Window.GetByID(windowId)
	if window == nil {
		globalApplication.debug("windowShouldUnconditionallyClose: window not found", "windowId", windowId)
		return false
	}
	unconditionallyClose := window.shouldUnconditionallyClose()
	globalApplication.debug("windowShouldUnconditionallyClose check", "windowId", windowId, "unconditionallyClose", unconditionallyClose)
	return unconditionallyClose
}

func (w *macosWebviewWindow) setZoom(zoom float64) {
	zoom, ok := zoomFactor(zoom, w.parent.options.Frameless)
	if !ok {
		return
	}
	w.setNativeZoom(zoom)
}
