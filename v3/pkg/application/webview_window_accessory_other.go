//go:build !darwin || ios || server

package application

import "unsafe"

func macAccessoryControllerKind(unsafe.Pointer) (MacAccessoryViewControllerKind, error) {
	return MacAccessoryViewControllerKindUnknown, ErrMacAccessoryControllerUnsupported
}

func macAccessoryControllerSupportsScrollEdgeEffectStyle(unsafe.Pointer) bool {
	return false
}

func macAccessoryControllerSetScrollEdgeEffectStyle(unsafe.Pointer, MacScrollEdgeEffectStyle) error {
	return ErrMacAccessoryControllerUnsupported
}

func macAccessoryControllerScrollEdgeEffectStyle(unsafe.Pointer) (MacScrollEdgeEffectStyle, error) {
	return MacScrollEdgeEffectStyleAutomatic, ErrMacAccessoryControllerUnsupported
}

// Native accessories exist only with AppKit. The model still builds
// everywhere so shared code can construct accessories unconditionally; the
// attach methods return ErrMacAccessoryUnsupported and setters are no-ops.
const macAccessoriesSupported = false

func macAccessoryAttachTitlebar(macAccessoryWindow, *MacAccessory) error {
	return ErrMacAccessoryUnsupported
}

func macAccessoryAttachPane(*MacSplitPane, *MacAccessory, bool) error {
	return ErrMacAccessoryUnsupported
}

func macAccessoryDetachNative(*MacAccessory)      {}
func macAccessoryFlushNativeWindow(*NativeWindow) {}

func macAccessorySetHidden(unsafe.Pointer, bool)                   {}
func macAccessorySetHeight(unsafe.Pointer, float64)                {}
func macAccessorySetFullScreenMinHeight(unsafe.Pointer, float64)   {}
func macAccessorySetAutomaticallyAdjustsSize(unsafe.Pointer, bool) {}
func macAccessorySetAppliesContentInsets(unsafe.Pointer, bool)     {}
func macAccessoryControlSetEnabled(unsafe.Pointer, uint64, bool)   {}
func macAccessoryControlSetHidden(unsafe.Pointer, uint64, bool)    {}
func macAccessoryControlSetTooltip(unsafe.Pointer, uint64, string) {}
func macAccessoryControlSetWidth(unsafe.Pointer, uint64, float64)  {}
func macAccessoryControlSetText(unsafe.Pointer, uint64, string)    {}
func macAccessoryControlSetSymbol(unsafe.Pointer, uint64, string)  {}
func macAccessoryControlSetPlaceholder(unsafe.Pointer, uint64, string) {
}
func macAccessoryControlSetIncremental(unsafe.Pointer, uint64, bool) {}
func macAccessoryControlSetSegments(unsafe.Pointer, uint64, []string, []string, int) {
}
func macAccessoryControlSetSelectedSegment(unsafe.Pointer, uint64, int) {}
func macAccessoryControlSetMenu(unsafe.Pointer, uint64, *Menu)          {}
