//go:build darwin

package updater

import "testing"

func TestBundleTarget(t *testing.T) {
	tests := []struct{ name, exe, want string }{
		{"inside app bundle", "/Applications/MyApp.app/Contents/MacOS/MyApp", "/Applications/MyApp.app"},
		{"nested in subdir", "/Users/user/Desktop/My App.app/Contents/MacOS/myapp", "/Users/user/Desktop/My App.app"},
		{"plain binary no bundle", "/usr/local/bin/mytool", "/usr/local/bin/mytool"},
		{"relative path no bundle", "./myapp", "./myapp"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := bundleTarget(tc.exe)
			if got != tc.want {
				t.Errorf("bundleTarget(%q) = %q, want %q", tc.exe, got, tc.want)
			}
		})
	}
}
