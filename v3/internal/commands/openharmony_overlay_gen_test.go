package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenHarmonyOverlayGenWritesTemplate(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "overlay.json")
	if err := OpenHarmonyOverlayGen(&OpenHarmonyOverlayGenOptions{Out: out}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("overlay is empty")
	}
	gen := filepath.Join(dir, "gen", "main_openharmony.gen.go")
	if _, err := os.Stat(gen); err != nil {
		t.Fatalf("generated file missing: %v", err)
	}
}
