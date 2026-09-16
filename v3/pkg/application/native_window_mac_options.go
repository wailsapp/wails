package application

// Option resolution shared by the macOS NativeWindow implementation and its
// tests. Everything in this file is plain Go so the mapping from MacWindow
// values to the native configuration can be verified without AppKit.

// nativeMacWindowLevelCode is the integer contract between Go and the
// NSWindow-level helper in mac_window_chrome_darwin.m. The values are private
// to Wails; they are not NSWindowLevel constants.
const (
	nativeMacWindowLevelCodeNormal = iota
	nativeMacWindowLevelCodeFloating
	nativeMacWindowLevelCodeTornOffMenu
	nativeMacWindowLevelCodeModalPanel
	nativeMacWindowLevelCodeMainMenu
	nativeMacWindowLevelCodeStatus
	nativeMacWindowLevelCodePopUpMenu
	nativeMacWindowLevelCodeScreenSaver
)

// nativeMacWindowLevelCode maps a MacWindowLevel to the code understood by
// windowChromeSetLevel. Unknown levels resolve to the normal level, matching
// WebviewWindow, whose level switch ignores unrecognised values.
func nativeMacWindowLevelCode(level MacWindowLevel) int {
	switch level {
	case MacWindowLevelFloating:
		return nativeMacWindowLevelCodeFloating
	case MacWindowLevelTornOffMenu:
		return nativeMacWindowLevelCodeTornOffMenu
	case MacWindowLevelModalPanel:
		return nativeMacWindowLevelCodeModalPanel
	case MacWindowLevelMainMenu:
		return nativeMacWindowLevelCodeMainMenu
	case MacWindowLevelStatus:
		return nativeMacWindowLevelCodeStatus
	case MacWindowLevelPopUpMenu:
		return nativeMacWindowLevelCodePopUpMenu
	case MacWindowLevelScreenSaver:
		return nativeMacWindowLevelCodeScreenSaver
	default:
		return nativeMacWindowLevelCodeNormal
	}
}

// effectiveNativeMacWindowLevel resolves the initial level of a NativeWindow
// with the same precedence WebviewWindow uses: an explicit Mac.WindowLevel wins,
// then AlwaysOnTop or a floating NSPanel selects the floating level.
func effectiveNativeMacWindowLevel(options NativeWindowOptions) MacWindowLevel {
	if options.Mac.WindowLevel != "" {
		return options.Mac.WindowLevel
	}
	if options.AlwaysOnTop ||
		(options.Mac.WindowClass == MacWindowClassPanel && options.Mac.PanelPreferences.FloatingPanel) {
		return MacWindowLevelFloating
	}
	return MacWindowLevelNormal
}

// effectiveNativeMacTabbingMode resolves the zero-value sentinel the same way
// WebviewWindow does: an unset tabbing mode disallows tabbing.
func effectiveNativeMacTabbingMode(mode MacWindowTabbingMode) MacWindowTabbingMode {
	if mode == MacWindowTabbingModeDefault {
		return MacWindowTabbingModeDisallowed
	}
	return mode
}

// nativeMacTabbingModeValue converts a resolved MacWindowTabbingMode to the
// NSWindowTabbingMode value AppKit expects (the Go enum is offset by one).
func nativeMacTabbingModeValue(mode MacWindowTabbingMode) int {
	return int(effectiveNativeMacTabbingMode(mode)) - 1
}

// nativeMacFrame describes how the NSWindow frame of a NativeWindow is built.
//
// NativeWindowOptions has no Frameless field, so TitleBar.Hide is the native
// frameless switch. When it is set the window follows exactly the rules
// WebviewWindowOptions.Frameless uses on macOS: CornerType and CornerRadius
// select between AppKit's own rounded frame (window buttons hidden) and a true
// borderless window with square or custom rounded corners.
type nativeMacFrame struct {
	// frameless reports that no titlebar is shown (TitleBar.Hide).
	frameless bool
	// borderless reports that NSWindowStyleMaskBorderless is used, which
	// happens for square corners or a custom corner radius.
	borderless bool
	// cornerRadius is the custom rounded-corner radius applied to the content,
	// or 0 when AppKit's own corners (or square corners) are used.
	cornerRadius float64
}

func resolveNativeMacFrame(mac MacWindow) nativeMacFrame {
	frame := nativeMacFrame{frameless: mac.TitleBar.Hide}
	if !frame.frameless {
		return frame
	}
	if mac.CornerType == MacWindowCornerTypeSquare {
		frame.borderless = true
		return frame
	}
	if mac.CornerRadius > 0 {
		frame.borderless = true
		frame.cornerRadius = mac.CornerRadius
	}
	return frame
}

// resolveNativeMacContentLayout resolves the window-level content layout for
// the native primary pane. Automatic follows TitleBar.FullSizeContent exactly
// as it does for a WebviewWindow. A frameless window (TitleBar.Hide) is built
// with a full-size content view, so Automatic is edge-to-edge there too, which
// is what windowApplyContentLayout derives from a frameless WebviewWindow's
// style mask.
func resolveNativeMacContentLayout(mac MacWindow) MacContentLayout {
	if mac.TitleBar.Hide && !validExplicitMacContentLayout(mac.ContentLayout) {
		return MacContentLayoutEdgeToEdge
	}
	return resolveMacContentLayout(mac, MacContentLayoutAutomatic)
}

func validExplicitMacContentLayout(layout MacContentLayout) bool {
	return validMacContentLayout(layout) && layout != MacContentLayoutAutomatic
}

// nativeMacSurfaceIsTransparent reports whether the NSWindow surface must be
// clear. Every backdrop except Normal is drawn through a clear window, and a
// borderless window with a custom corner radius needs a clear surface so the
// rounded content mask is visible.
func nativeMacSurfaceIsTransparent(mac MacWindow) bool {
	if mac.Backdrop != MacBackdropNormal {
		return true
	}
	return resolveNativeMacFrame(mac).cornerRadius > 0
}

// nativeMacLiquidGlassCornerRadius validates the glass corner radius the same
// way WebviewWindow.applyLiquidGlass does.
func nativeMacLiquidGlassCornerRadius(glass MacLiquidGlass) float64 {
	if glass.CornerRadius < 0 {
		return 0
	}
	return glass.CornerRadius
}
