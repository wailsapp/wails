package commands

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCpSourcePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix executable permissions")
	}
	previousArgs, previousFooter := os.Args, DisableFooter
	t.Cleanup(func() { os.Args, DisableFooter = previousArgs, previousFooter })
	for _, mode := range []os.FileMode{0755, 0600} {
		t.Run(mode.String(), func(t *testing.T) {
			root := t.TempDir()
			source := filepath.Join(root, "source")
			if err := os.WriteFile(source, []byte("payload"), mode); err != nil {
				t.Fatal(err)
			}
			sourceInfo, err := os.Stat(source)
			if err != nil {
				t.Fatal(err)
			}
			created := filepath.Join(root, "created")
			existing := filepath.Join(root, "existing")
			if err := os.WriteFile(existing, []byte("old"), 0600); err != nil {
				t.Fatal(err)
			}
			for path, want := range map[string]os.FileMode{created: sourceInfo.Mode().Perm(), existing: 0600} {
				os.Args = []string{"wails3", "tool", "cp", source, path}
				if err := Cp(&CpOptions{}); err != nil {
					t.Fatal(err)
				}
				info, err := os.Stat(path)
				if err != nil {
					t.Fatal(err)
				}
				if got := info.Mode().Perm(); got != want {
					t.Errorf("%s permissions=%#o, want %#o", path, got, want)
				}
			}
		})
	}
}
