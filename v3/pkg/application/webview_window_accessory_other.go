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
