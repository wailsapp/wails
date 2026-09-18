package application

import "testing"

func TestSystemTraySmartDefaults(t *testing.T) {
	tests := []struct {
		name             string
		hasWindow        bool
		hasMenu          bool
		wantLeftHandler  bool
		wantRightHandler bool
	}{
		{name: "window only", hasWindow: true, wantLeftHandler: true},
		{name: "menu only", hasMenu: true, wantRightHandler: true},
		{name: "window and menu", hasWindow: true, hasMenu: true, wantLeftHandler: true, wantRightHandler: true},
		{name: "neither"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tray := newSystemTray(1)
			if test.hasWindow {
				tray.attachedWindow.Window = &WebviewWindow{}
			}
			if test.hasMenu {
				tray.menu = NewMenu()
			}
			tray.applySmartDefaults()

			if got := tray.clickHandler != nil; got != test.wantLeftHandler {
				t.Errorf("left handler present = %v, want %v", got, test.wantLeftHandler)
			}
			if got := tray.rightClickHandler != nil; got != test.wantRightHandler {
				t.Errorf("right handler present = %v, want %v", got, test.wantRightHandler)
			}
		})
	}
}

func TestSystemTrayExplicitHandlersOverrideSmartDefaults(t *testing.T) {
	tray := newSystemTray(1)
	tray.attachedWindow.Window = &WebviewWindow{}
	tray.menu = NewMenu()

	leftCalls := 0
	rightCalls := 0
	tray.OnClick(func() { leftCalls++ })
	tray.OnRightClick(func() { rightCalls++ })
	tray.applySmartDefaults()

	tray.clickHandler()
	tray.rightClickHandler()
	if leftCalls != 1 || rightCalls != 1 {
		t.Fatalf("custom handlers replaced: left calls = %d, right calls = %d", leftCalls, rightCalls)
	}
}

func TestSystemTrayRemovableModel(t *testing.T) {
	tray := newSystemTray(1)
	if tray.IsRemovable() || tray.AutosaveName() != "" {
		t.Fatalf("new tray should not be removable: removable=%v autosave=%q", tray.IsRemovable(), tray.AutosaveName())
	}
	if got := tray.SetRemovable(true, "com.example.tray"); got != tray {
		t.Fatal("SetRemovable should return the tray for chaining")
	}
	if !tray.IsRemovable() {
		t.Error("IsRemovable = false after SetRemovable(true)")
	}
	if tray.AutosaveName() != "com.example.tray" {
		t.Errorf("AutosaveName = %q, want com.example.tray", tray.AutosaveName())
	}
	tray.SetRemovable(false, "")
	if tray.IsRemovable() || tray.AutosaveName() != "" {
		t.Errorf("SetRemovable(false, \"\") left removable=%v autosave=%q", tray.IsRemovable(), tray.AutosaveName())
	}
}

func TestSystemTrayVisibilityModelWithoutImpl(t *testing.T) {
	tray := newSystemTray(1)
	if !tray.IsVisible() {
		t.Fatal("new tray should report visible")
	}
	tray.Hide()
	if tray.IsVisible() {
		t.Error("IsVisible = true after Hide")
	}
	tray.SetVisible(true)
	if !tray.IsVisible() {
		t.Error("IsVisible = false after SetVisible(true)")
	}
	tray.SetVisible(false)
	if tray.IsVisible() {
		t.Error("IsVisible = true after SetVisible(false)")
	}
	tray.Show()
	if !tray.IsVisible() {
		t.Error("IsVisible = false after Show")
	}
}

func TestSystemTraySymbolAndTooltipModel(t *testing.T) {
	tray := newSystemTray(1)
	if tray.Symbol() != "" || tray.Tooltip() != "" {
		t.Fatal("new tray should have no symbol or tooltip")
	}
	tray.SetSymbol("star.fill").SetSymbolConfiguration(-4, MacSymbolWeightBold)
	if tray.Symbol() != "star.fill" {
		t.Errorf("Symbol = %q, want star.fill", tray.Symbol())
	}
	if tray.symbolPointSize != 0 {
		t.Errorf("negative point size should clamp to 0, got %v", tray.symbolPointSize)
	}
	if tray.symbolWeight != MacSymbolWeightBold {
		t.Errorf("symbolWeight = %v, want bold", tray.symbolWeight)
	}
	tray.SetSymbolConfiguration(18, MacSymbolWeightUnspecified)
	if tray.symbolPointSize != 18 {
		t.Errorf("symbolPointSize = %v, want 18", tray.symbolPointSize)
	}
	tray.SetTooltip("Feedback demo")
	if tray.Tooltip() != "Feedback demo" {
		t.Errorf("Tooltip = %q", tray.Tooltip())
	}
	handlerSet := tray.OnVisibilityChange(func(bool) {})
	if handlerSet != tray || tray.visibilityHandler == nil {
		t.Error("OnVisibilityChange should store the handler and return the tray")
	}
}

func TestMacSymbolWeightString(t *testing.T) {
	tests := map[MacSymbolWeight]string{
		MacSymbolWeightUnspecified: "unspecified",
		MacSymbolWeightUltraLight:  "ultraLight",
		MacSymbolWeightThin:        "thin",
		MacSymbolWeightLight:       "light",
		MacSymbolWeightRegular:     "regular",
		MacSymbolWeightMedium:      "medium",
		MacSymbolWeightSemibold:    "semibold",
		MacSymbolWeightBold:        "bold",
		MacSymbolWeightHeavy:       "heavy",
		MacSymbolWeightBlack:       "black",
		MacSymbolWeight(42):        "unspecified",
	}
	for weight, want := range tests {
		if got := weight.String(); got != want {
			t.Errorf("MacSymbolWeight(%d).String() = %q, want %q", int(weight), got, want)
		}
	}
}
