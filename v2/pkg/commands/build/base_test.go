package build

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func Test_commandPrettifier(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{
			name:  "empty",
			input: []string{},
			want:  "",
		},
		{
			name:  "one arg",
			input: []string{"one"},
			want:  "one",
		},
		{
			name:  "args where one has spaces",
			input: []string{"one", "two three"},
			want:  `one "two three"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commandPrettifier(tt.input); got != tt.want {
				t.Errorf("commandPrettifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateRuntimeWrapper_permissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file modes are not applicable on Windows")
	}
	dir := t.TempDir()
	if err := generateRuntimeWrapper(&Options{WailsJSDir: dir}); err != nil {
		t.Fatal(err)
	}
	runtimeDir := filepath.Join(dir, "wailsjs", "runtime")
	count := 0
	err := filepath.Walk(runtimeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		count++
		// Must match what `wails generate module` leaves behind (#6192)
		if got := info.Mode().Perm(); got != 0o755 {
			t.Errorf("%s: mode = %o, want 755", path, got)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count < 2 {
		t.Errorf("expected extracted runtime files, found %d entries", count)
	}
}
