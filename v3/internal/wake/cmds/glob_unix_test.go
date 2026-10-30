//go:build !windows

package cmds

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRouteExpandsUnquotedGlobs(t *testing.T) {
	for _, tc := range []struct {
		pattern string
		want    []string
	}{
		{"*.go", []string{"run", "alpha.go", "beta file.go"}},
		{"alpha.?o", []string{"run", "alpha.go"}},
		{"[ab]*.go", []string{"run", "alpha.go", "beta file.go"}},
		{`"*.go"`, []string{"run", "*.go"}},
		{`\*.go`, []string{"run", "*.go"}},
		{`"[ab]*.go"`, []string{"run", "[ab]*.go"}},
	} {
		t.Run(tc.pattern, func(t *testing.T) {
			base := installFakeGo(t, `printf '%s\n' "$@"`)
			for _, name := range []string{"alpha.go", "beta file.go", "alpha.txt"} {
				if err := os.WriteFile(filepath.Join(base, name), nil, 0644); err != nil {
					t.Fatal(err)
				}
			}
			output := captureRoutedOutput(t, "go run "+tc.pattern, base)
			if want := strings.Join(tc.want, "\n") + "\n"; output != want {
				t.Errorf("arguments = %q, want %q", output, want)
			}
		})
	}
}

func TestGoRunExpandsSourceGlob(t *testing.T) {
	_, output, err := runOfflineGo(t, "go run -buildvcs=false *.go", map[string]string{
		"go.mod":     "module example.com/wake-glob-check\n\ngo 1.25.0\n",
		"main.go":    "package main\n\nimport \"fmt\"\n\nfunc main() { fmt.Print(message()) }\n",
		"message.go": "package main\n\nfunc message() string { return \"all-files-ran\" }\n",
	})
	if err != nil {
		t.Fatalf("source glob failed: %v; output: %s", err, output)
	}
	if output != "all-files-ran" {
		t.Errorf("program output = %q, want all-files-ran", output)
	}
}
