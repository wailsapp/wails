package commands

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCpMultipleDestinations(t *testing.T) {
	previousArgs, previousFooter := os.Args, DisableFooter
	t.Cleanup(func() { os.Args, DisableFooter = previousArgs, previousFooter })
	for _, mode := range []string{"files", "directories"} {
		t.Run(mode, func(t *testing.T) {
			root := t.TempDir()
			source := filepath.Join(root, "source file.txt")
			const content = "copy this content to every destination"
			if err := os.WriteFile(source, []byte(content), 0600); err != nil {
				t.Fatal(err)
			}
			first, second := filepath.Join(root, "first directory"), filepath.Join(root, "second directory")
			for _, dir := range []string{first, second} {
				if err := os.Mkdir(dir, 0700); err != nil {
					t.Fatal(err)
				}
			}
			copies := []string{filepath.Join(first, filepath.Base(source)), filepath.Join(second, filepath.Base(source)), filepath.Join(root, "renamed copy.txt")}
			if err := os.WriteFile(copies[0], []byte("older content that is much longer than the source file contents"), 0600); err != nil {
				t.Fatal(err)
			}
			destinations := copies
			if mode == "directories" {
				destinations = []string{first, second, copies[2]}
			}
			os.Args = append([]string{"wails3", "tool", "cp", source, source}, destinations...)
			if err := Cp(&CpOptions{}); err != nil {
				t.Fatal(err)
			}
			for _, path := range append(copies, source) {
				data, err := os.ReadFile(path)
				if err != nil || string(data) != content {
					t.Errorf("%s = %q, error=%v", path, data, err)
				}
			}
		})
	}
}

func TestCpRequiresSourceAndDestination(t *testing.T) {
	previousArgs, previousFooter := os.Args, DisableFooter
	t.Cleanup(func() { os.Args, DisableFooter = previousArgs, previousFooter })
	for _, args := range [][]string{
		{"wails3", "tool", "cp"},
		{"wails3", "tool", "cp", "source"},
	} {
		os.Args = args
		if err := Cp(&CpOptions{}); err == nil {
			t.Errorf("accepted incomplete arguments: %v", args)
		}
	}
}

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
				contents, err := os.ReadFile(path)
				require.NoError(t, err)
				require.Equal(t, "payload", string(contents), path)
			}
		})
	}
}
