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

// stubTrayImpl embeds the interface so only the methods a test exercises need bodies.
type stubTrayImpl struct{ systemTrayImpl }

func (stubTrayImpl) setLabel(string)   {}
func (stubTrayImpl) setTooltip(string) {}

// runningTray returns a tray with a stub impl attached, as if Run had happened.
func runningTray(t *testing.T) *SystemTray {
	t.Helper()
	probe := &threadProbeApp{}
	probe.onMain.Store(true)
	prev := globalApplication
	globalApplication = &App{impl: probe}
	t.Cleanup(func() { globalApplication = prev })

	tray := newSystemTray(1)
	tray.impl = stubTrayImpl{}
	return tray
}

// Label() must report the value set after the tray is running.
func TestSystemTrayLabelUpdatedAfterRun(t *testing.T) {
	tray := runningTray(t)
	tray.SetLabel("after")
	if got := tray.Label(); got != "after" {
		t.Fatalf("Label() = %q, want %q", got, "after")
	}
}

// An explicit tooltip wins over the label until it is cleared, whatever the call order.
func TestSystemTrayTooltipWinsOverLabel(t *testing.T) {
	tray := runningTray(t)
	tray.SetTooltip("explicit")
	tray.SetLabel("label")
	if got := tray.tooltipOrLabel(); got != "explicit" {
		t.Fatalf("tooltipOrLabel() = %q, want %q", got, "explicit")
	}
	tray.SetTooltip("")
	if got := tray.tooltipOrLabel(); got != "label" {
		t.Fatalf("tooltipOrLabel() = %q, want %q", got, "label")
	}
}
