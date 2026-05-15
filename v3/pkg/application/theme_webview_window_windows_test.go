//go:build windows

package application

import "testing"

func TestResolveWindowsEffectiveTheme(t *testing.T) {
	tests := []struct {
		name               string
		winTheme           WinTheme
		appTheme           AppTheme
		wantTheme          theme
		wantFollowApp      bool
	}{
		// Explicit window overrides — never follow the application.
		{"WinDark overrides app", WinDark, AppLight, dark, false},
		{"WinLight overrides app", WinLight, AppDark, light, false},
		{"WinSystemDefault overrides app", WinSystemDefault, AppDark, systemDefault, false},

		// WinAppDefault (and empty / unset) — always follow the application.
		{"WinAppDefault + AppDark → dark", WinAppDefault, AppDark, dark, true},
		{"WinAppDefault + AppLight → light", WinAppDefault, AppLight, light, true},
		{"WinAppDefault + AppSystemDefault → systemDefault", WinAppDefault, AppSystemDefault, systemDefault, true},

		// Empty WinTheme is treated as WinAppDefault.
		{"empty WinTheme + AppDark → dark", "", AppDark, dark, true},
		{"empty WinTheme + AppLight → light", "", AppLight, light, true},
		{"empty WinTheme + AppSystemDefault → systemDefault", "", AppSystemDefault, systemDefault, true},
		{"empty WinTheme + empty AppTheme → systemDefault", "", "", systemDefault, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotTheme, gotFollow := resolveWindowsEffectiveTheme(tc.winTheme, tc.appTheme)
			if gotTheme != tc.wantTheme {
				t.Errorf("theme: got %v, want %v", gotTheme, tc.wantTheme)
			}
			if gotFollow != tc.wantFollowApp {
				t.Errorf("followApp: got %v, want %v", gotFollow, tc.wantFollowApp)
			}
		})
	}
}
