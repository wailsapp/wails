//go:build darwin && !ios && !server

package application

import "testing"

// These tests cover the MacWindow option resolution a NativeWindow performs
// before touching AppKit. They run without a GUI session.

func TestNativeMacWindowLevelPrecedence(t *testing.T) {
	tests := []struct {
		name    string
		options NativeWindowOptions
		want    MacWindowLevel
	}{
		{"default", NativeWindowOptions{}, MacWindowLevelNormal},
		{"always on top", NativeWindowOptions{AlwaysOnTop: true}, MacWindowLevelFloating},
		{"floating panel", NativeWindowOptions{Mac: MacWindow{
			WindowClass:      MacWindowClassPanel,
			PanelPreferences: MacPanelPreferences{FloatingPanel: true},
		}}, MacWindowLevelFloating},
		{"floating preference without panel class", NativeWindowOptions{Mac: MacWindow{
			PanelPreferences: MacPanelPreferences{FloatingPanel: true},
		}}, MacWindowLevelNormal},
		{"explicit level wins over always on top", NativeWindowOptions{
			AlwaysOnTop: true,
			Mac:         MacWindow{WindowLevel: MacWindowLevelStatus},
		}, MacWindowLevelStatus},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := effectiveNativeMacWindowLevel(tt.options); got != tt.want {
				t.Fatalf("effectiveNativeMacWindowLevel = %q, want %q", got, tt.want)
			}
		})
	}
	// The WebviewWindow resolver is the source of truth; both must agree.
	webview := WebviewWindowOptions{AlwaysOnTop: true, Mac: MacWindow{WindowLevel: MacWindowLevelPopUpMenu}}
	native := NativeWindowOptions{AlwaysOnTop: true, Mac: MacWindow{WindowLevel: MacWindowLevelPopUpMenu}}
	if effectiveMacWindowLevel(webview) != effectiveNativeMacWindowLevel(native) {
		t.Fatal("native and webview window level resolution diverged")
	}
}

func TestNativeMacWindowLevelCodes(t *testing.T) {
	codes := map[MacWindowLevel]int{
		MacWindowLevelNormal:      nativeMacWindowLevelCodeNormal,
		MacWindowLevelFloating:    nativeMacWindowLevelCodeFloating,
		MacWindowLevelTornOffMenu: nativeMacWindowLevelCodeTornOffMenu,
		MacWindowLevelModalPanel:  nativeMacWindowLevelCodeModalPanel,
		MacWindowLevelMainMenu:    nativeMacWindowLevelCodeMainMenu,
		MacWindowLevelStatus:      nativeMacWindowLevelCodeStatus,
		MacWindowLevelPopUpMenu:   nativeMacWindowLevelCodePopUpMenu,
		MacWindowLevelScreenSaver: nativeMacWindowLevelCodeScreenSaver,
		MacWindowLevel("bogus"):   nativeMacWindowLevelCodeNormal,
	}
	seen := make(map[int]MacWindowLevel)
	for level, want := range codes {
		got := nativeMacWindowLevelCode(level)
		if got != want {
			t.Errorf("nativeMacWindowLevelCode(%q) = %d, want %d", level, got, want)
		}
		if previous, duplicate := seen[got]; duplicate && level != "bogus" && previous != "bogus" {
			t.Errorf("levels %q and %q share code %d", previous, level, got)
		}
		seen[got] = level
	}
}

func TestNativeMacTabbingModeResolution(t *testing.T) {
	if got := effectiveNativeMacTabbingMode(MacWindowTabbingModeDefault); got != MacWindowTabbingModeDisallowed {
		t.Fatalf("default tabbing mode resolved to %d, want disallowed", got)
	}
	// NSWindowTabbingMode: automatic 0, preferred 1, disallowed 2.
	tests := map[MacWindowTabbingMode]int{
		MacWindowTabbingModeDefault:    2,
		MacWindowTabbingModeAutomatic:  0,
		MacWindowTabbingModePreferred:  1,
		MacWindowTabbingModeDisallowed: 2,
	}
	for mode, want := range tests {
		if got := nativeMacTabbingModeValue(mode); got != want {
			t.Errorf("nativeMacTabbingModeValue(%d) = %d, want %d", mode, got, want)
		}
	}
}

func TestResolveNativeMacFrame(t *testing.T) {
	tests := []struct {
		name string
		mac  MacWindow
		want nativeMacFrame
	}{
		{"titled window ignores corner options", MacWindow{CornerType: MacWindowCornerTypeSquare, CornerRadius: 9}, nativeMacFrame{}},
		{"hidden titlebar keeps the AppKit frame", MacWindow{TitleBar: MacTitleBar{Hide: true}}, nativeMacFrame{frameless: true}},
		{"square corners are borderless", MacWindow{TitleBar: MacTitleBar{Hide: true}, CornerType: MacWindowCornerTypeSquare, CornerRadius: 9},
			nativeMacFrame{frameless: true, borderless: true}},
		{"custom radius is borderless and masked", MacWindow{TitleBar: MacTitleBar{Hide: true}, CornerRadius: 14},
			nativeMacFrame{frameless: true, borderless: true, cornerRadius: 14}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveNativeMacFrame(tt.mac); got != tt.want {
				t.Fatalf("resolveNativeMacFrame = %+v, want %+v", got, tt.want)
			}
		})
	}
	// The AppKit-frame decision must match the WebviewWindow helper.
	for _, mac := range []MacWindow{{}, {CornerRadius: 3}, {CornerType: MacWindowCornerTypeSquare}} {
		mac.TitleBar.Hide = true
		if usesNativeMacFramelessFrame(mac) == resolveNativeMacFrame(mac).borderless {
			t.Fatalf("frameless frame decision diverged from usesNativeMacFramelessFrame for %+v", mac)
		}
	}
}

func TestResolveNativeMacContentLayout(t *testing.T) {
	tests := []struct {
		name string
		mac  MacWindow
		want MacContentLayout
	}{
		{"automatic without full size content", MacWindow{}, MacContentLayoutBelowToolbar},
		{"automatic with full size content", MacWindow{TitleBar: MacTitleBar{FullSizeContent: true}}, MacContentLayoutEdgeToEdge},
		{"explicit below toolbar", MacWindow{ContentLayout: MacContentLayoutBelowToolbar, TitleBar: MacTitleBar{FullSizeContent: true}}, MacContentLayoutBelowToolbar},
		{"explicit edge to edge", MacWindow{ContentLayout: MacContentLayoutEdgeToEdge}, MacContentLayoutEdgeToEdge},
		{"invalid falls back to automatic", MacWindow{ContentLayout: MacContentLayout(99)}, MacContentLayoutBelowToolbar},
		{"frameless automatic is edge to edge", MacWindow{TitleBar: MacTitleBar{Hide: true}}, MacContentLayoutEdgeToEdge},
		{"frameless explicit below toolbar", MacWindow{TitleBar: MacTitleBar{Hide: true}, ContentLayout: MacContentLayoutBelowToolbar}, MacContentLayoutBelowToolbar},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveNativeMacContentLayout(tt.mac); got != tt.want {
				t.Fatalf("resolveNativeMacContentLayout = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestNativeMacSurfaceIsTransparent(t *testing.T) {
	if nativeMacSurfaceIsTransparent(MacWindow{}) {
		t.Fatal("normal backdrop on a titled window must be opaque")
	}
	for _, backdrop := range []MacBackdrop{MacBackdropTransparent, MacBackdropTranslucent, MacBackdropLiquidGlass} {
		if !nativeMacSurfaceIsTransparent(MacWindow{Backdrop: backdrop}) {
			t.Fatalf("backdrop %d must use a clear window surface", backdrop)
		}
	}
	rounded := MacWindow{TitleBar: MacTitleBar{Hide: true}, CornerRadius: 10}
	if !nativeMacSurfaceIsTransparent(rounded) {
		t.Fatal("a custom corner radius needs a clear window surface")
	}
	if nativeMacSurfaceIsTransparent(MacWindow{CornerRadius: 10}) {
		t.Fatal("a corner radius without TitleBar.Hide has no effect and must stay opaque")
	}
}

func TestNativeMacLiquidGlassCornerRadius(t *testing.T) {
	if got := nativeMacLiquidGlassCornerRadius(MacLiquidGlass{CornerRadius: -4}); got != 0 {
		t.Fatalf("negative radius resolved to %v, want 0", got)
	}
	if got := nativeMacLiquidGlassCornerRadius(MacLiquidGlass{CornerRadius: 6}); got != 6 {
		t.Fatalf("radius resolved to %v, want 6", got)
	}
}

func TestNativeWindowConfigMapsPanelAndFrameOptions(t *testing.T) {
	options := NativeWindowOptions{
		Width:       320,
		Height:      200,
		HideOnClose: true,
		Mac: MacWindow{
			WindowClass: MacWindowClassPanel,
			PanelPreferences: MacPanelPreferences{
				FloatingPanel:          true,
				BecomesKeyOnlyIfNeeded: true,
				NonActivating:          true,
				UtilityWindow:          true,
			},
			DisableEscapeExitsFullscreen: true,
			TitleBar:                     MacTitleBar{Hide: true},
			CornerRadius:                 12,
		},
	}
	window := &macosNativeWindow{parent: newNativeWindow(options)}
	config := window.windowConfig(resolveNativeMacFrame(options.Mac))
	if int(config.width) != 320 || int(config.height) != 200 || !bool(config.hideOnClose) {
		t.Fatalf("size or hideOnClose not mapped: %+v", config)
	}
	if !bool(config.isPanel) || !bool(config.floatingPanel) || !bool(config.becomesKeyOnlyIfNeeded) ||
		!bool(config.nonActivating) || !bool(config.utilityWindow) {
		t.Fatalf("panel preferences not mapped: %+v", config)
	}
	if !bool(config.disableEscapeExitsFullscreen) {
		t.Fatal("DisableEscapeExitsFullscreen not mapped")
	}
	if !bool(config.frameless) || !bool(config.borderless) || float64(config.cornerRadius) != 12 {
		t.Fatalf("frame not mapped: %+v", config)
	}
	if uint(config.windowID) != window.parent.id {
		t.Fatal("window ID not mapped")
	}
}
