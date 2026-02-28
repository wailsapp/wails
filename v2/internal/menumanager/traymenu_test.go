package menumanager

import (
	"testing"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

func TestNewTrayMenu(t *testing.T) {
	tm := &menu.TrayMenu{
		Label:    "Test Label",
		Tooltip:  "Test Tooltip",
		Disabled: false,
	}

	tray := NewTrayMenu(tm)

	if tray.Label != tm.Label {
		t.Errorf("Expected label %s, got %s", tm.Label, tray.Label)
	}

	if tray.Tooltip != tm.Tooltip {
		t.Errorf("Expected tooltip %s, got %s", tm.Tooltip, tray.Tooltip)
	}

	if tray.Disabled != tm.Disabled {
		t.Errorf("Expected disabled %v, got %v", tm.Disabled, tray.Disabled)
	}
}

func TestNewTrayMenuWithANSI(t *testing.T) {
	tm := &menu.TrayMenu{
		Label: "\033[31mRed Label\033[0m",
	}

	tray := NewTrayMenu(tm)

	if tray.Label != tm.Label {
		t.Errorf("Expected label %s, got %s", tm.Label, tray.Label)
	}

	if len(tray.StyledLabel) == 0 {
		t.Error("Expected StyledLabel to be populated for ANSI text")
	}
}
