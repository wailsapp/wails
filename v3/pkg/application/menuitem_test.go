package application_test

import (
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestMenuItem_GetAccelerator(t *testing.T) {
	tests := []struct {
		name        string
		menuItem    *application.MenuItem
		expectedAcc string
	}{
		{
			name:        "Get existing accelerator",
			menuItem:    application.NewMenuItem("Item 1").SetAccelerator("ctrl+a"),
			expectedAcc: "Ctrl+A",
		},
		{
			name:        "Get non-existing accelerator",
			menuItem:    application.NewMenuItem("Item 2"),
			expectedAcc: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			acc := test.menuItem.GetAccelerator()
			if acc != test.expectedAcc {
				t.Errorf("Expected accelerator to be %v, but got %v", test.expectedAcc, acc)
			}
		})
	}
}

func TestMenuItem_RemoveAccelerator(t *testing.T) {
	tests := []struct {
		name     string
		menuItem *application.MenuItem
	}{
		{
			name:     "Remove existing accelerator",
			menuItem: application.NewMenuItem("Item 1").SetAccelerator("Ctrl+A"),
		},
		{
			name:     "Remove non-existing accelerator",
			menuItem: application.NewMenuItem("Item 2"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.menuItem.RemoveAccelerator()
			acc := test.menuItem.GetAccelerator()
			if acc != "" {
				t.Errorf("Expected accelerator to be removed, but got %v", acc)
			}
		})
	}
}
