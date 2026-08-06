package application

import (
	"errors"
	"testing"
	"unsafe"
)

func TestMacScrollEdgeEffectStyleValues(t *testing.T) {
	if MacScrollEdgeEffectStyleAutomatic != 0 || MacScrollEdgeEffectStyleSoft != 1 ||
		MacScrollEdgeEffectStyleHard != 2 {
		t.Fatalf("unexpected scroll-edge style values: %d, %d, %d",
			MacScrollEdgeEffectStyleAutomatic, MacScrollEdgeEffectStyleSoft,
			MacScrollEdgeEffectStyleHard)
	}
	for _, style := range []MacScrollEdgeEffectStyle{
		MacScrollEdgeEffectStyleAutomatic,
		MacScrollEdgeEffectStyleSoft,
		MacScrollEdgeEffectStyleHard,
	} {
		if !validMacScrollEdgeEffectStyle(style) {
			t.Fatalf("validMacScrollEdgeEffectStyle(%d) = false", style)
		}
	}
	if validMacScrollEdgeEffectStyle(MacScrollEdgeEffectStyle(99)) {
		t.Fatal("unknown scroll-edge style was accepted")
	}
}

func TestWrapMacAccessoryViewControllerRejectsNil(t *testing.T) {
	controller, err := WrapMacAccessoryViewController(nil)
	if controller != nil {
		t.Fatal("nil native controller produced a wrapper")
	}
	if !errors.Is(err, ErrMacAccessoryControllerRequired) {
		t.Fatalf("error = %v, want ErrMacAccessoryControllerRequired", err)
	}
}

func TestMacAccessoryViewControllerRejectsUnknownStyleBeforeNativeCall(t *testing.T) {
	marker := byte(0)
	controller := &MacAccessoryViewController{
		native: unsafe.Pointer(&marker),
		kind:   MacAccessoryViewControllerKindTitlebar,
	}
	if err := controller.SetPreferredScrollEdgeEffectStyle(MacScrollEdgeEffectStyle(99)); err == nil {
		t.Fatal("unknown scroll-edge style was accepted")
	}
}

func TestNilMacAccessoryViewControllerIsSafe(t *testing.T) {
	var controller *MacAccessoryViewController
	if controller.NativeController() != nil {
		t.Fatal("nil wrapper returned a native controller")
	}
	if controller.Kind() != MacAccessoryViewControllerKindUnknown {
		t.Fatal("nil wrapper returned a controller kind")
	}
	if controller.SupportsPreferredScrollEdgeEffectStyle() {
		t.Fatal("nil wrapper reported style support")
	}
	if !errors.Is(controller.SetPreferredScrollEdgeEffectStyle(MacScrollEdgeEffectStyleAutomatic),
		ErrMacAccessoryControllerRequired) {
		t.Fatal("nil wrapper did not report the required-controller error")
	}
	if _, err := controller.PreferredScrollEdgeEffectStyle(); !errors.Is(err, ErrMacAccessoryControllerRequired) {
		t.Fatalf("getter error = %v, want ErrMacAccessoryControllerRequired", err)
	}
}
