//go:build darwin && !ios

package mac

import (
	"path/filepath"
	"testing"
)

func TestResourcePathForExecutable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		executable string
		want       string
		ok         bool
	}{
		{
			name:       "application bundle",
			executable: "/Applications/My App.app/Contents/MacOS/My App",
			want:       "/Applications/My App.app/Contents/Resources",
			ok:         true,
		},
		{
			name:       "nested application bundle",
			executable: "/Users/user/Desktop/My App.app/Contents/MacOS/My App",
			want:       "/Users/user/Desktop/My App.app/Contents/Resources",
			ok:         true,
		},
		{
			name:       "normalises executable path",
			executable: "/Applications/My App.app/Contents/MacOS/../MacOS/My App",
			want:       "/Applications/My App.app/Contents/Resources",
			ok:         true,
		},
		{
			name:       "executable outside application bundle",
			executable: "/usr/local/bin/my-app",
			ok:         false,
		},
		{
			name:       "missing Contents directory",
			executable: "/Applications/My App.app/MacOS/My App",
			ok:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := resourcePathForExecutable(tt.executable)
			if ok != tt.ok {
				t.Fatalf("resourcePathForExecutable(%q) returned ok=%t, want %t", tt.executable, ok, tt.ok)
			}
			if got != filepath.FromSlash(tt.want) {
				t.Errorf("resourcePathForExecutable(%q) = %q, want %q", tt.executable, got, filepath.FromSlash(tt.want))
			}
		})
	}
}
