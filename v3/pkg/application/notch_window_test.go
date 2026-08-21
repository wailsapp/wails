package application

import (
	"testing"
	"time"
)

func TestUnsupportedNotchWindowIsInert(t *testing.T) {
	if notchWindowsSupported {
		t.Skip("notch windows are supported on this platform")
	}

	window := (&WindowManager{}).NewNotchWindow(NotchWindowOptions{})
	if window == nil {
		t.Fatal("NewNotchWindow returned nil")
	}
	if window.Show() != window || window.Hide() != window || window.Visibility() {
		t.Fatal("unsupported notch window handle was not inert")
	}
	window.Close()
}

func TestNormaliseNotchWindowOptions(t *testing.T) {
	tests := []struct {
		name  string
		input NotchWindowOptions
		want  NotchWindowOptions
	}{
		{
			name:  "zero values use defaults",
			input: NotchWindowOptions{},
			want: NotchWindowOptions{
				Width:          defaultNotchWindowWidth,
				Height:         defaultNotchWindowHeight,
				AnimationSpeed: defaultNotchWindowAnimationSpeed,
			},
		},
		{
			name: "custom values are preserved",
			input: NotchWindowOptions{
				Width:          800,
				Height:         120,
				Animated:       true,
				AnimationSpeed: 600 * time.Millisecond,
			},
			want: NotchWindowOptions{
				Width:          800,
				Height:         120,
				Animated:       true,
				AnimationSpeed: 600 * time.Millisecond,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := normaliseNotchWindowOptions(test.input)
			if got.Width != test.want.Width || got.Height != test.want.Height ||
				got.Animated != test.want.Animated || got.AnimationSpeed != test.want.AnimationSpeed {
				t.Fatalf("normaliseNotchWindowOptions() = %+v, want %+v", got, test.want)
			}
		})
	}
}

func TestWebviewOptionsForNotchWindow(t *testing.T) {
	targetScreen := &Screen{ID: "42"}
	input := NotchWindowOptions{
		Width:  700,
		Height: 100,
		Screen: targetScreen,
		WindowOptions: WebviewWindowOptions{
			Name:       "second-alert",
			Hidden:     true,
			URL:        "/second",
			MinWidth:   500,
			MaxHeight:  900,
			StartState: WindowStateFullscreen,
			Mac: MacWindow{
				Appearance:       NSAppearanceNameDarkAqua,
				PanelPreferences: MacPanelPreferences{UtilityWindow: true},
			},
		},
	}

	got := webviewOptionsForNotchWindow(input)
	if got.Width != 760 || got.Height != 146 {
		t.Fatalf("outer size = %dx%d, want 760x146", got.Width, got.Height)
	}
	if got.Name != "second-alert" || !got.Hidden || got.URL != "/second" {
		t.Fatalf("content options were not preserved: %+v", got)
	}
	if got.Screen != targetScreen {
		t.Fatalf("target screen = %p, want %p", got.Screen, targetScreen)
	}
	if got.Mac.Appearance != NSAppearanceNameDarkAqua {
		t.Fatalf("unowned macOS option was not preserved: %q", got.Mac.Appearance)
	}
	if !got.Frameless || !got.DisableResize || got.BackgroundType != BackgroundTypeTransparent {
		t.Fatalf("native notch geometry options not applied: %+v", got)
	}
	if got.MinWidth != 0 || got.MaxHeight != 0 || got.StartState != WindowStateNormal {
		t.Fatalf("incompatible geometry options were not cleared: %+v", got)
	}
	if got.Mac.WindowClass != MacWindowClassPanel || !got.Mac.PanelPreferences.NonActivating ||
		!got.Mac.PanelPreferences.FloatingPanel || got.Mac.PanelPreferences.UtilityWindow {
		t.Fatalf("panel options not applied: %+v", got.Mac)
	}
	if got.Mac.WindowLevel != MacWindowLevelPopUpMenu {
		t.Fatalf("window level = %q, want %q", got.Mac.WindowLevel, MacWindowLevelPopUpMenu)
	}
}
