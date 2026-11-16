//go:build darwin && purego && !ios && !server

package application

// CGO-free port of the frameless style and corner decisions of the cgo window
// (webview_window_darwin.go: windowStyleMask, windowNew's content view setup,
// windowSetFrameless, and the Go-side setFrameless tail). cgo does not use any
// private AppKit API for corners: the rounded shape is the content view layer's
// cornerRadius with masksToBounds, so nothing here goes through the private API
// seam.

// AppKit style mask bits (NSWindow.h of the macOS 26 SDK).
const (
	nsWindowStyleMaskBorderless          = 0       // NSWindowStyleMaskBorderless
	nsWindowStyleMaskTitled              = 1 << 0  // NSWindowStyleMaskTitled
	nsWindowStyleMaskClosable            = 1 << 1  // NSWindowStyleMaskClosable
	nsWindowStyleMaskMiniaturizable      = 1 << 2  // NSWindowStyleMaskMiniaturizable
	nsWindowStyleMaskResizable           = 1 << 3  // NSWindowStyleMaskResizable
	nsWindowStyleMaskFullSizeContentView = 1 << 15 // NSWindowStyleMaskFullSizeContentView
)

// macWindowStyleMask is cgo windowStyleMask. A frameless window with a square
// or custom-radius corner is borderless, so AppKit does not draw its own frame
// and corners; any other frameless window keeps the titled frame but extends
// its content under the (transparent) title bar. Without frameless, the corner
// options have no effect.
func macWindowStyleMask(frameless, squareCorners bool, cornerRadius float64) uint {
	if frameless && (squareCorners || cornerRadius > 0) {
		return nsWindowStyleMaskBorderless | nsWindowStyleMaskResizable | nsWindowStyleMaskMiniaturizable
	}
	mask := uint(nsWindowStyleMaskTitled | nsWindowStyleMaskClosable | nsWindowStyleMaskMiniaturizable | nsWindowStyleMaskResizable)
	if frameless {
		mask |= nsWindowStyleMaskFullSizeContentView
	}
	return mask
}

// macFramelessStyleMask is cgo windowSetFrameless's style: the frameless mask
// for the requested shape, plus the panel bits (non-activating, utility) that
// the window already had. Any other bit of the current mask is dropped.
func macFramelessStyleMask(current uint, frameless, squareCorners bool, cornerRadius float64) uint {
	const panelBits = nsWindowStyleMaskNonactivatingPanel | nsWindowStyleMaskUtilityWindow
	return macWindowStyleMask(frameless, squareCorners, cornerRadius) | (current & panelBits)
}

// macContentLayerAction says what happens to the content view's layer.
type macContentLayerAction int

const (
	// macContentLayerKeep leaves the layer untouched.
	macContentLayerKeep macContentLayerAction = iota
	// macContentLayerRound makes the view layer-backed, with cornerRadius set
	// and masksToBounds on.
	macContentLayerRound
	// macContentLayerReset clears cornerRadius and masksToBounds on an existing
	// layer. It does not make the view layer-backed.
	macContentLayerReset
)

// macCreationContentLayer is the layer decision of cgo windowNew: only a
// frameless window with a custom radius (and not square) rounds the content.
func macCreationContentLayer(frameless, squareCorners bool, cornerRadius float64) macContentLayerAction {
	if frameless && !squareCorners && cornerRadius > 0 {
		return macContentLayerRound
	}
	return macContentLayerKeep
}

// macFramelessContentLayer is the layer decision of cgo windowSetFrameless.
// The rounded and reset cases apply to the content view; the native-frame
// frameless case (rounded corners with no radius) leaves the layer as it is,
// so AppKit's own corners stay.
func macFramelessContentLayer(frameless, squareCorners bool, cornerRadius float64) macContentLayerAction {
	if !frameless {
		return macContentLayerReset
	}
	if squareCorners || cornerRadius > 0 {
		if squareCorners {
			return macContentLayerReset
		}
		return macContentLayerRound
	}
	return macContentLayerKeep
}

// macCreationHidesNativeTitlebar reports whether cgo windowNew hides the native
// title bar of a frameless window that keeps AppKit's own frame (rounded corners
// with no custom radius). The title bar is made transparent and its title hidden.
func macCreationHidesNativeTitlebar(frameless, squareCorners bool, cornerRadius float64) bool {
	return frameless && !squareCorners && cornerRadius == 0
}

// applyContentLayer applies a content layer decision to view. Only the round
// action makes the view layer-backed; reset on a view without a layer is a
// no-op, as in cgo where view.layer is nil.
func applyContentLayer(view id, action macContentLayerAction, cornerRadius float64) {
	switch action {
	case macContentLayerRound:
		view.send("setWantsLayer:", true)
		if layer := view.send("layer"); !layer.isNil() {
			layer.send("setCornerRadius:", CGFloat(cornerRadius))
			layer.send("setMasksToBounds:", true)
		}
	case macContentLayerReset:
		if layer := view.send("layer"); !layer.isNil() {
			layer.send("setCornerRadius:", CGFloat(0))
			layer.send("setMasksToBounds:", false)
		}
	}
}

// windowSetFrameless is cgo windowSetFrameless for the NSWindow win: it sets
// the style mask and the content view layer for the requested frame. The
// title-bar and button tail is applied by setFrameless.
func windowSetFrameless(win id, frameless, squareCorners bool, cornerRadius float64) {
	win.send("setStyleMask:", macFramelessStyleMask(get[uint](win, "styleMask"), frameless, squareCorners, cornerRadius))
	applyContentLayer(win.send("contentView"), macFramelessContentLayer(frameless, squareCorners, cornerRadius), cornerRadius)
}

// setFrameless is the purego setFrameless: cgo windowSetFrameless on the main
// thread, then the cgo Go-side tail. A frameless window hides its title bar;
// otherwise the title bar presets come back. The title-bar buttons are then
// reapplied, since a native-frame frameless window hides them.
func (w *macosWebviewWindow) setFrameless(frameless bool) {
	runOnMain(func() {
		macOptions := w.parent.options.Mac
		windowSetFrameless(w.win(), frameless, macOptions.CornerType == MacWindowCornerTypeSquare, macOptions.CornerRadius)
		if frameless {
			w.setTitleBarAppearsTransparent(true)
			w.setTitleVisibility(true)
		} else {
			w.setTitleBarAppearsTransparent(macOptions.TitleBar.AppearsTransparent)
			w.setTitleVisibility(macOptions.TitleBar.HideTitle)
		}
		if usesNativeMacFramelessFrame(macOptions) {
			w.applyWindowButtonStates()
		}
	})
}
