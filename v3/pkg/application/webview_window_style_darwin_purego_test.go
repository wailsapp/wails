//go:build darwin && purego && !ios && !server

package application

import (
	"testing"
)

// AppKit SDK values used by the expected masks below:
// Titled 1<<0, Closable 1<<1, Miniaturizable 1<<2, Resizable 1<<3,
// Borderless 0, UtilityWindow 1<<4, NonactivatingPanel 1<<7, FullSizeContentView 1<<15.
const (
	wantTitledMask    uint = 1<<0 | 1<<1 | 1<<2 | 1<<3 // 15
	wantBorderlessRsz uint = 1<<3 | 1<<2               // 12 (Borderless | Resizable | Miniaturizable)
	wantFullSizeBit   uint = 1 << 15
)

func TestMacWindowStyleMask(t *testing.T) {
	tests := []struct {
		name      string
		frameless bool
		square    bool
		radius    float64
		want      uint
	}{
		{name: "standard window", want: wantTitledMask},
		{name: "square corners without frameless are ignored", square: true, want: wantTitledMask},
		{name: "radius without frameless is ignored", radius: 24, want: wantTitledMask},
		{name: "frameless keeps the titled frame under the title bar", frameless: true, want: wantTitledMask | wantFullSizeBit},
		{name: "frameless with a negative radius takes the native-frame branch", frameless: true, radius: -1, want: wantTitledMask | wantFullSizeBit},
		{name: "frameless square corners are borderless", frameless: true, square: true, want: wantBorderlessRsz},
		{name: "frameless square ignores the radius", frameless: true, square: true, radius: 24, want: wantBorderlessRsz},
		{name: "frameless custom radius is borderless", frameless: true, radius: 24, want: wantBorderlessRsz},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := macWindowStyleMask(tt.frameless, tt.square, tt.radius); got != tt.want {
				t.Errorf("macWindowStyleMask(%v, %v, %v) = %#x, want %#x", tt.frameless, tt.square, tt.radius, got, tt.want)
			}
		})
	}
}

func TestMacFramelessStyleMaskKeepsOnlyPanelBits(t *testing.T) {
	const (
		nonActivating = nsWindowStyleMaskNonactivatingPanel
		utility       = nsWindowStyleMaskUtilityWindow
	)
	tests := []struct {
		name      string
		current   uint
		frameless bool
		square    bool
		radius    float64
		want      uint
	}{
		{name: "unframing drops the full-size content bit", current: wantTitledMask | wantFullSizeBit,
			frameless: false, want: wantTitledMask},
		{name: "frameless square keeps the panel bits", current: wantTitledMask | nonActivating | utility,
			frameless: true, square: true, want: wantBorderlessRsz | nonActivating | utility},
		{name: "frameless default drops the old borderless bits but keeps the panel bits",
			current: wantBorderlessRsz | nonActivating, frameless: true, want: wantTitledMask | wantFullSizeBit | nonActivating},
		{name: "a stale full-size bit is not carried over", current: wantTitledMask | wantFullSizeBit,
			frameless: true, radius: 12, want: wantBorderlessRsz},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := macFramelessStyleMask(tt.current, tt.frameless, tt.square, tt.radius); got != tt.want {
				t.Errorf("macFramelessStyleMask(%#x, %v, %v, %v) = %#x, want %#x", tt.current, tt.frameless, tt.square, tt.radius, got, tt.want)
			}
		})
	}
}

func TestMacCreationContentLayer(t *testing.T) {
	tests := []struct {
		name      string
		frameless bool
		square    bool
		radius    float64
		want      macContentLayerAction
	}{
		{name: "standard window keeps its layer", want: macContentLayerKeep},
		{name: "frameless default keeps AppKit's corners", frameless: true, want: macContentLayerKeep},
		{name: "frameless square has no radius to apply", frameless: true, square: true, radius: 24, want: macContentLayerKeep},
		{name: "frameless custom radius rounds the content", frameless: true, radius: 24, want: macContentLayerRound},
		{name: "custom radius without frameless does nothing", radius: 24, want: macContentLayerKeep},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := macCreationContentLayer(tt.frameless, tt.square, tt.radius); got != tt.want {
				t.Errorf("macCreationContentLayer(%v, %v, %v) = %v, want %v", tt.frameless, tt.square, tt.radius, got, tt.want)
			}
		})
	}
}

func TestMacFramelessContentLayer(t *testing.T) {
	tests := []struct {
		name      string
		frameless bool
		square    bool
		radius    float64
		want      macContentLayerAction
	}{
		{name: "unframing resets the layer", want: macContentLayerReset},
		{name: "unframing a square window resets the layer", square: true, radius: 24, want: macContentLayerReset},
		{name: "frameless default leaves the layer for AppKit's corners", frameless: true, want: macContentLayerKeep},
		{name: "frameless square resets the layer", frameless: true, square: true, radius: 24, want: macContentLayerReset},
		{name: "frameless custom radius rounds the content", frameless: true, radius: 24, want: macContentLayerRound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := macFramelessContentLayer(tt.frameless, tt.square, tt.radius); got != tt.want {
				t.Errorf("macFramelessContentLayer(%v, %v, %v) = %v, want %v", tt.frameless, tt.square, tt.radius, got, tt.want)
			}
		})
	}
}

func TestMacCreationHidesNativeTitlebar(t *testing.T) {
	tests := []struct {
		name      string
		frameless bool
		square    bool
		radius    float64
		want      bool
	}{
		{name: "standard window", want: false},
		{name: "frameless default hides the title bar", frameless: true, want: true},
		{name: "frameless square does not", frameless: true, square: true, want: false},
		{name: "frameless custom radius does not", frameless: true, radius: 24, want: false},
		{name: "frameless negative radius does not (cgo tests radius == 0)", frameless: true, radius: -1, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := macCreationHidesNativeTitlebar(tt.frameless, tt.square, tt.radius); got != tt.want {
				t.Errorf("macCreationHidesNativeTitlebar(%v, %v, %v) = %v, want %v", tt.frameless, tt.square, tt.radius, got, tt.want)
			}
		})
	}
}

// TestMacWindowStyleMaskAcceptedByAppKit checks the bit constants against
// AppKit: each mask must survive a round trip through a real NSWindow. A wrong
// constant would be stored as a different mask.
func TestMacWindowStyleMaskAcceptedByAppKit(t *testing.T) {
	masks := []uint{
		macWindowStyleMask(false, false, 0),
		macWindowStyleMask(true, false, 0),
		macWindowStyleMask(true, true, 0),
		macWindowStyleMask(true, false, 24),
	}
	// AppKit aborts if windows are made off the main thread: create them
	// there, assert here.
	got := make([]uint, len(masks))
	runOnMain(func() {
		for i, mask := range masks {
			win := class("NSWindow").send("alloc").send("initWithContentRect:styleMask:backing:defer:",
				rect(0, 0, 100, 100), mask, uint(2), false)
			got[i] = get[uint](win, "styleMask")
			win.send("release")
		}
	})
	for i, mask := range masks {
		if got[i] != mask {
			t.Errorf("NSWindow styleMask = %#x after init with %#x", got[i], mask)
		}
	}
}

// TestContentLayerRoundAndResetOnAppKitView checks applyContentLayer against a
// real NSView: round makes the view layer-backed with the radius and
// masksToBounds, reset clears them again.
func TestContentLayerRoundAndResetOnAppKitView(t *testing.T) {
	var (
		hasLayer                 bool
		roundRadius, resetRadius CGFloat
		roundMasks, resetMasks   bool
	)
	runOnMain(func() {
		view := class("NSView").send("alloc").send("initWithFrame:", rect(0, 0, 100, 100))
		defer view.send("release")
		applyContentLayer(view, macContentLayerRound, 24)
		layer := view.send("layer")
		if hasLayer = !layer.isNil(); !hasLayer {
			return
		}
		roundRadius = get[CGFloat](layer, "cornerRadius")
		roundMasks = get[bool](layer, "masksToBounds")
		applyContentLayer(view, macContentLayerReset, 0)
		layer = view.send("layer")
		resetRadius = get[CGFloat](layer, "cornerRadius")
		resetMasks = get[bool](layer, "masksToBounds")
	})
	if !hasLayer {
		t.Fatal("round left the view without a layer")
	}
	if roundRadius != 24 {
		t.Errorf("cornerRadius after round = %v, want 24", roundRadius)
	}
	if !roundMasks {
		t.Error("masksToBounds after round = false, want true")
	}
	if resetRadius != 0 {
		t.Errorf("cornerRadius after reset = %v, want 0", resetRadius)
	}
	if resetMasks {
		t.Error("masksToBounds after reset = true, want false")
	}
}

// TestWindowSetFramelessStyleOnAppKitWindow checks the frameless style that
// windowSetFrameless leaves on a real NSWindow, for the shapes that matter.
func TestWindowSetFramelessStyleOnAppKitWindow(t *testing.T) {
	tests := []struct {
		name      string
		frameless bool
		square    bool
		radius    float64
	}{
		{name: "frameless default", frameless: true},
		{name: "frameless square", frameless: true, square: true},
		{name: "frameless custom radius", frameless: true, radius: 24},
		{name: "back to a standard window", frameless: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got uint
			runOnMain(func() {
				win := class("NSWindow").send("alloc").send("initWithContentRect:styleMask:backing:defer:",
					rect(0, 0, 100, 100), uint(wantTitledMask), uint(2), false)
				defer win.send("release")
				win.send("setContentView:", class("NSView").send("alloc").send("initWithFrame:", rect(0, 0, 100, 100)))
				windowSetFrameless(win, tt.frameless, tt.square, tt.radius)
				got = get[uint](win, "styleMask")
			})
			if want := macWindowStyleMask(tt.frameless, tt.square, tt.radius); got != want {
				t.Errorf("styleMask = %#x, want %#x", got, want)
			}
		})
	}
}
