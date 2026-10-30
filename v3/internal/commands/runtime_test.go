package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/wailsapp/wails/v3/internal/assetserver/bundledassets"
)

func TestGenerateRuntimeWithoutSourceTree(t *testing.T) {
	previous := DisableFooter
	t.Cleanup(func() { DisableFooter = previous })
	dir := t.TempDir()
	if err := GenerateRuntime(&RuntimeOptions{Directory: dir}); err != nil {
		t.Fatal(err)
	}
	production, debug := bundledassets.RuntimeFiles()
	for name, want := range map[string][]byte{"runtime.js": production, "runtime-debug.js": debug} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil || len(data) == 0 || !bytes.Equal(data, want) {
			t.Errorf("generated %s differs from embedded bundle: length=%d, error=%v", name, len(data), err)
		}
	}
}
