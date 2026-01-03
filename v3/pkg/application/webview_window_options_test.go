package application

import (
	"testing"
)

func TestNewRGBA(t *testing.T) {
	rgba := NewRGBA(100, 150, 200, 255)

	if rgba.Red != 100 {
		t.Errorf("Red = %d, want 100", rgba.Red)
	}
	if rgba.Green != 150 {
		t.Errorf("Green = %d, want 150", rgba.Green)
	}
	if rgba.Blue != 200 {
		t.Errorf("Blue = %d, want 200", rgba.Blue)
	}
	if rgba.Alpha != 255 {
		t.Errorf("Alpha = %d, want 255", rgba.Alpha)
	}
}

func TestNewRGB(t *testing.T) {
	rgba := NewRGB(100, 150, 200)

	if rgba.Red != 100 {
		t.Errorf("Red = %d, want 100", rgba.Red)
	}
	if rgba.Green != 150 {
		t.Errorf("Green = %d, want 150", rgba.Green)
	}
	if rgba.Blue != 200 {
		t.Errorf("Blue = %d, want 200", rgba.Blue)
	}
	if rgba.Alpha != 255 {
		t.Errorf("Alpha = %d, want 255 (default)", rgba.Alpha)
	}
}

func TestNewRGBPtr(t *testing.T) {
	ptr := NewRGBPtr(0x12, 0x34, 0x56)

	if ptr == nil {
		t.Fatal("NewRGBPtr returned nil")
	}

	// RGB is packed as 0x00BBGGRR
	expected := uint32(0x12) | (uint32(0x34) << 8) | (uint32(0x56) << 16)
	if *ptr != expected {
		t.Errorf("*ptr = 0x%X, want 0x%X", *ptr, expected)
	}
}

func TestBackgroundType_Constants(t *testing.T) {
	if BackgroundTypeSolid != 0 {
		t.Error("BackgroundTypeSolid should be 0")
	}
	if BackgroundTypeTransparent != 1 {
		t.Error("BackgroundTypeTransparent should be 1")
	}
	if BackgroundTypeTranslucent != 2 {
		t.Error("BackgroundTypeTranslucent should be 2")
	}
}

func TestBackdropType_Constants(t *testing.T) {
	if Auto != 0 {
		t.Error("Auto should be 0")
	}
	if None != 1 {
		t.Error("None should be 1")
	}
	if Mica != 2 {
		t.Error("Mica should be 2")
	}
	if Acrylic != 3 {
		t.Error("Acrylic should be 3")
	}
	if Tabbed != 4 {
		t.Error("Tabbed should be 4")
	}
}

func TestTheme_Constants(t *testing.T) {
	if SystemDefault != 0 {
		t.Error("SystemDefault should be 0")
	}
	if Dark != 1 {
		t.Error("Dark should be 1")
	}
	if Light != 2 {
		t.Error("Light should be 2")
	}
}

func TestMacBackdrop_Constants(t *testing.T) {
	if MacBackdropNormal != 0 {
		t.Error("MacBackdropNormal should be 0")
	}
	if MacBackdropTransparent != 1 {
		t.Error("MacBackdropTransparent should be 1")
	}
	if MacBackdropTranslucent != 2 {
		t.Error("MacBackdropTranslucent should be 2")
	}
	if MacBackdropLiquidGlass != 3 {
		t.Error("MacBackdropLiquidGlass should be 3")
	}
}

func TestMacToolbarStyle_Constants(t *testing.T) {
	if MacToolbarStyleAutomatic != 0 {
		t.Error("MacToolbarStyleAutomatic should be 0")
	}
	if MacToolbarStyleExpanded != 1 {
		t.Error("MacToolbarStyleExpanded should be 1")
	}
	if MacToolbarStylePreference != 2 {
		t.Error("MacToolbarStylePreference should be 2")
	}
	if MacToolbarStyleUnified != 3 {
		t.Error("MacToolbarStyleUnified should be 3")
	}
	if MacToolbarStyleUnifiedCompact != 4 {
		t.Error("MacToolbarStyleUnifiedCompact should be 4")
	}
}

func TestWebviewGpuPolicy_Constants(t *testing.T) {
	if WebviewGpuPolicyAlways != 0 {
		t.Error("WebviewGpuPolicyAlways should be 0")
	}
	if WebviewGpuPolicyOnDemand != 1 {
		t.Error("WebviewGpuPolicyOnDemand should be 1")
	}
	if WebviewGpuPolicyNever != 2 {
		t.Error("WebviewGpuPolicyNever should be 2")
	}
}

func TestMacTitleBarDefault(t *testing.T) {
	titleBar := MacTitleBarDefault

	if titleBar.AppearsTransparent != false {
		t.Error("MacTitleBarDefault.AppearsTransparent should be false")
	}
	if titleBar.Hide != false {
		t.Error("MacTitleBarDefault.Hide should be false")
	}
	if titleBar.HideTitle != false {
		t.Error("MacTitleBarDefault.HideTitle should be false")
	}
	if titleBar.FullSizeContent != false {
		t.Error("MacTitleBarDefault.FullSizeContent should be false")
	}
	if titleBar.UseToolbar != false {
		t.Error("MacTitleBarDefault.UseToolbar should be false")
	}
	if titleBar.HideToolbarSeparator != false {
		t.Error("MacTitleBarDefault.HideToolbarSeparator should be false")
	}
}

func TestMacTitleBarHidden(t *testing.T) {
	titleBar := MacTitleBarHidden

	if titleBar.AppearsTransparent != true {
		t.Error("MacTitleBarHidden.AppearsTransparent should be true")
	}
	if titleBar.Hide != false {
		t.Error("MacTitleBarHidden.Hide should be false")
	}
	if titleBar.HideTitle != true {
		t.Error("MacTitleBarHidden.HideTitle should be true")
	}
	if titleBar.FullSizeContent != true {
		t.Error("MacTitleBarHidden.FullSizeContent should be true")
	}
	if titleBar.UseToolbar != false {
		t.Error("MacTitleBarHidden.UseToolbar should be false")
	}
	if titleBar.HideToolbarSeparator != false {
		t.Error("MacTitleBarHidden.HideToolbarSeparator should be false")
	}
}

func TestMacTitleBarHiddenInset(t *testing.T) {
	titleBar := MacTitleBarHiddenInset

	if titleBar.AppearsTransparent != true {
		t.Error("MacTitleBarHiddenInset.AppearsTransparent should be true")
	}
	if titleBar.Hide != false {
		t.Error("MacTitleBarHiddenInset.Hide should be false")
	}
	if titleBar.HideTitle != true {
		t.Error("MacTitleBarHiddenInset.HideTitle should be true")
	}
	if titleBar.FullSizeContent != true {
		t.Error("MacTitleBarHiddenInset.FullSizeContent should be true")
	}
	if titleBar.UseToolbar != true {
		t.Error("MacTitleBarHiddenInset.UseToolbar should be true")
	}
	if titleBar.HideToolbarSeparator != true {
		t.Error("MacTitleBarHiddenInset.HideToolbarSeparator should be true")
	}
}

func TestMacTitleBarHiddenInsetUnified(t *testing.T) {
	titleBar := MacTitleBarHiddenInsetUnified

	if titleBar.AppearsTransparent != true {
		t.Error("MacTitleBarHiddenInsetUnified.AppearsTransparent should be true")
	}
	if titleBar.ToolbarStyle != MacToolbarStyleUnified {
		t.Error("MacTitleBarHiddenInsetUnified.ToolbarStyle should be MacToolbarStyleUnified")
	}
}

func TestMacAppearanceType_Constants(t *testing.T) {
	tests := []struct {
		name     string
		value    MacAppearanceType
		expected string
	}{
		{"DefaultAppearance", DefaultAppearance, ""},
		{"NSAppearanceNameAqua", NSAppearanceNameAqua, "NSAppearanceNameAqua"},
		{"NSAppearanceNameDarkAqua", NSAppearanceNameDarkAqua, "NSAppearanceNameDarkAqua"},
		{"NSAppearanceNameVibrantLight", NSAppearanceNameVibrantLight, "NSAppearanceNameVibrantLight"},
	}

	for _, tt := range tests {
		if string(tt.value) != tt.expected {
			t.Errorf("%s = %q, want %q", tt.name, string(tt.value), tt.expected)
		}
	}
}

func TestMacWindowLevel_Constants(t *testing.T) {
	tests := []struct {
		name     string
		value    MacWindowLevel
		expected string
	}{
		{"MacWindowLevelNormal", MacWindowLevelNormal, "normal"},
		{"MacWindowLevelFloating", MacWindowLevelFloating, "floating"},
		{"MacWindowLevelTornOffMenu", MacWindowLevelTornOffMenu, "tornOffMenu"},
		{"MacWindowLevelModalPanel", MacWindowLevelModalPanel, "modalPanel"},
		{"MacWindowLevelMainMenu", MacWindowLevelMainMenu, "mainMenu"},
		{"MacWindowLevelStatus", MacWindowLevelStatus, "status"},
		{"MacWindowLevelPopUpMenu", MacWindowLevelPopUpMenu, "popUpMenu"},
		{"MacWindowLevelScreenSaver", MacWindowLevelScreenSaver, "screenSaver"},
	}

	for _, tt := range tests {
		if string(tt.value) != tt.expected {
			t.Errorf("%s = %q, want %q", tt.name, string(tt.value), tt.expected)
		}
	}
}

func TestWebviewWindowOptions_Defaults(t *testing.T) {
	opts := WebviewWindowOptions{}

	// Verify zero values
	if opts.Name != "" {
		t.Error("Name should default to empty string")
	}
	if opts.Title != "" {
		t.Error("Title should default to empty string")
	}
	if opts.Width != 0 {
		t.Error("Width should default to 0")
	}
	if opts.Height != 0 {
		t.Error("Height should default to 0")
	}
	if opts.AlwaysOnTop != false {
		t.Error("AlwaysOnTop should default to false")
	}
	if opts.Frameless != false {
		t.Error("Frameless should default to false")
	}
}

func TestWindowsWindow_Defaults(t *testing.T) {
	opts := WindowsWindow{}

	if opts.BackdropType != Auto {
		t.Error("BackdropType should default to Auto")
	}
	if opts.DisableIcon != false {
		t.Error("DisableIcon should default to false")
	}
	if opts.Theme != SystemDefault {
		t.Error("Theme should default to SystemDefault")
	}
}

func TestMacWindow_Defaults(t *testing.T) {
	opts := MacWindow{}

	if opts.Backdrop != MacBackdropNormal {
		t.Error("Backdrop should default to MacBackdropNormal")
	}
	if opts.DisableShadow != false {
		t.Error("DisableShadow should default to false")
	}
}

func TestLinuxWindow_Defaults(t *testing.T) {
	opts := LinuxWindow{}

	if opts.WindowIsTranslucent != false {
		t.Error("WindowIsTranslucent should default to false")
	}
	if opts.WebviewGpuPolicy != WebviewGpuPolicyAlways {
		t.Error("WebviewGpuPolicy should default to WebviewGpuPolicyAlways")
	}
}

func TestCoreWebView2PermissionKind_Constants(t *testing.T) {
	if CoreWebView2PermissionKindUnknownPermission != 0 {
		t.Error("CoreWebView2PermissionKindUnknownPermission should be 0")
	}
	if CoreWebView2PermissionKindMicrophone != 1 {
		t.Error("CoreWebView2PermissionKindMicrophone should be 1")
	}
	if CoreWebView2PermissionKindCamera != 2 {
		t.Error("CoreWebView2PermissionKindCamera should be 2")
	}
}

func TestCoreWebView2PermissionState_Constants(t *testing.T) {
	if CoreWebView2PermissionStateDefault != 0 {
		t.Error("CoreWebView2PermissionStateDefault should be 0")
	}
	if CoreWebView2PermissionStateAllow != 1 {
		t.Error("CoreWebView2PermissionStateAllow should be 1")
	}
	if CoreWebView2PermissionStateDeny != 2 {
		t.Error("CoreWebView2PermissionStateDeny should be 2")
	}
}

func TestMacLiquidGlassStyle_Constants(t *testing.T) {
	if LiquidGlassStyleAutomatic != 0 {
		t.Error("LiquidGlassStyleAutomatic should be 0")
	}
	if LiquidGlassStyleLight != 1 {
		t.Error("LiquidGlassStyleLight should be 1")
	}
	if LiquidGlassStyleDark != 2 {
		t.Error("LiquidGlassStyleDark should be 2")
	}
	if LiquidGlassStyleVibrant != 3 {
		t.Error("LiquidGlassStyleVibrant should be 3")
	}
}

func TestNSVisualEffectMaterial_Constants(t *testing.T) {
	if NSVisualEffectMaterialAppearanceBased != 0 {
		t.Error("NSVisualEffectMaterialAppearanceBased should be 0")
	}
	if NSVisualEffectMaterialLight != 1 {
		t.Error("NSVisualEffectMaterialLight should be 1")
	}
	if NSVisualEffectMaterialAuto != -1 {
		t.Error("NSVisualEffectMaterialAuto should be -1")
	}
}
