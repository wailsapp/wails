package cmds

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGoInstallPreservesAllPackages(t *testing.T) {
	root, output, err := runOfflineGo(t, "go install -buildvcs=false ./first ./second", map[string]string{
		"go.mod":         "module example.com/wake-install-check\n\ngo 1.25.0\n",
		"first/main.go":  "package main\n\nfunc main() {}\n",
		"second/main.go": "package main\n\nfunc main() {}\n",
	})
	if err != nil {
		t.Fatalf("offline install failed: %v; output: %s", err, output)
	}
	extension := ""
	if runtime.GOOS == "windows" {
		extension = ".exe"
	}
	for _, name := range []string{"first", "second"} {
		info, err := os.Stat(filepath.Join(root, "bin", name+extension))
		if err != nil || !info.Mode().IsRegular() {
			t.Errorf("requested executable %s was not installed: %v", name, err)
		}
	}
}
