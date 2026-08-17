package application

import (
	"testing"
)

func TestActivationPolicy_Constants(t *testing.T) {
	if ActivationPolicyRegular != 0 {
		t.Error("ActivationPolicyRegular should be 0")
	}
	if ActivationPolicyAccessory != 1 {
		t.Error("ActivationPolicyAccessory should be 1")
	}
	if ActivationPolicyProhibited != 2 {
		t.Error("ActivationPolicyProhibited should be 2")
	}
}

func TestNativeTabIcon_Constants(t *testing.T) {
	tests := []struct {
		name     string
		icon     NativeTabIcon
		expected string
	}{
		{"NativeTabIconNone", NativeTabIconNone, ""},
		{"NativeTabIconHouse", NativeTabIconHouse, "house"},
		{"NativeTabIconGear", NativeTabIconGear, "gear"},
		{"NativeTabIconStar", NativeTabIconStar, "star"},
		{"NativeTabIconPerson", NativeTabIconPerson, "person"},
		{"NativeTabIconBell", NativeTabIconBell, "bell"},
		{"NativeTabIconMagnify", NativeTabIconMagnify, "magnifyingglass"},
		{"NativeTabIconList", NativeTabIconList, "list.bullet"},
		{"NativeTabIconFolder", NativeTabIconFolder, "folder"},
	}

	for _, tt := range tests {
		if string(tt.icon) != tt.expected {
			t.Errorf("%s = %q, want %q", tt.name, string(tt.icon), tt.expected)
		}
	}
}

func TestOptions_Defaults(t *testing.T) {
	opts := Options{}

	if opts.Name != "" {
		t.Error("Name should default to empty string")
	}
	if opts.Description != "" {
		t.Error("Description should default to empty string")
	}
	if opts.Icon != nil {
		t.Error("Icon should default to nil")
	}
	if opts.Logger != nil {
		t.Error("Logger should default to nil")
	}
	if opts.DisableDefaultSignalHandler != false {
		t.Error("DisableDefaultSignalHandler should default to false")
	}
}

func TestMacOptions_Defaults(t *testing.T) {
	opts := MacOptions{}

	if opts.ActivationPolicy != ActivationPolicyRegular {
		t.Error("ActivationPolicy should default to ActivationPolicyRegular")
	}
	if opts.ApplicationShouldTerminateAfterLastWindowClosed != false {
		t.Error("ApplicationShouldTerminateAfterLastWindowClosed should default to false")
	}
}

func TestWindowsOptions_Defaults(t *testing.T) {
	opts := WindowsOptions{}

	if opts.WndClass != "" {
		t.Error("WndClass should default to empty string")
	}
	if opts.DisableQuitOnLastWindowClosed != false {
		t.Error("DisableQuitOnLastWindowClosed should default to false")
	}
	if opts.WebviewUserDataPath != "" {
		t.Error("WebviewUserDataPath should default to empty string")
	}
	if opts.WebviewBrowserPath != "" {
		t.Error("WebviewBrowserPath should default to empty string")
	}
}

func TestLinuxOptions_Defaults(t *testing.T) {
	opts := LinuxOptions{}

	if opts.DisableQuitOnLastWindowClosed != false {
		t.Error("DisableQuitOnLastWindowClosed should default to false")
	}
	if opts.ProgramName != "" {
		t.Error("ProgramName should default to empty string")
	}
}

func TestIOSOptions_Defaults(t *testing.T) {
	opts := IOSOptions{}

	if opts.DisableInputAccessoryView != false {
		t.Error("DisableInputAccessoryView should default to false")
	}
	if opts.DisableScroll != false {
		t.Error("DisableScroll should default to false")
	}
	if opts.DisableBounce != false {
		t.Error("DisableBounce should default to false")
	}
	if opts.EnableBackForwardNavigationGestures != false {
		t.Error("EnableBackForwardNavigationGestures should default to false")
	}
	if opts.EnableNativeTabs != false {
		t.Error("EnableNativeTabs should default to false")
	}
}

func TestAndroidOptions_Defaults(t *testing.T) {
	opts := AndroidOptions{}

	if opts.DisableScroll != false {
		t.Error("DisableScroll should default to false")
	}
	if opts.DisableOverscroll != false {
		t.Error("DisableOverscroll should default to false")
	}
	if opts.EnableZoom != false {
		t.Error("EnableZoom should default to false")
	}
	if opts.DisableHardwareAcceleration != false {
		t.Error("DisableHardwareAcceleration should default to false")
	}
}

func TestNativeTabItem_Fields(t *testing.T) {
	item := NativeTabItem{
		Title:       "Home",
		SystemImage: NativeTabIconHouse,
	}

	if item.Title != "Home" {
		t.Errorf("Title = %q, want %q", item.Title, "Home")
	}
	if item.SystemImage != NativeTabIconHouse {
		t.Errorf("SystemImage = %q, want %q", item.SystemImage, NativeTabIconHouse)
	}
}
