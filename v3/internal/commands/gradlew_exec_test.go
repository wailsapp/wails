package commands

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/wailsapp/wails/v3/internal/gosod"
)

// TestEmbeddedGradlewExtractsExecutable guards against regressing #5606. The
// Android gradlew wrapper ships inside the embedded build_assets FS, where
// go:embed reports it as 0444 (no exec bit). Extracting it must still produce an
// executable file so `wails3 task android:package` can run ./gradlew. A test that
// only exercises a fstest.MapFS with a 0755 source would not catch this, because
// such a source never occurs in the real embedded path.
func TestEmbeddedGradlewExtractsExecutable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("executable bits are not represented in file modes on Windows")
	}

	// Sanity-check the premise: the embedded source really is non-executable.
	if info, err := fs.Stat(buildAssets, "build_assets/android/gradlew"); err != nil {
		t.Fatalf("stat embedded gradlew: %v", err)
	} else if info.Mode().Perm()&0111 != 0 {
		t.Fatalf("expected embedded gradlew to report a non-exec mode, got %v", info.Mode().Perm())
	}

	// The android subtree has no .tmpl files, so it extracts with nil data.
	androidFS, err := fs.Sub(buildAssets, "build_assets/android")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := gosod.New(androidFS).Extract(dir, nil); err != nil {
		t.Fatalf("extract android build assets: %v", err)
	}

	info, err := os.Stat(filepath.Join(dir, "gradlew"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0111 == 0 {
		t.Fatalf("expected extracted gradlew to be executable, got mode %v", info.Mode().Perm())
	}
}
