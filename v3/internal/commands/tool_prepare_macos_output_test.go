package commands

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRelocateMacOSOutputPreservesFilesAndRejectsCollision(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions vary on Windows")
	}
	project := t.TempDir()
	outputRoot := t.TempDir()
	t.Setenv("WAILS_MACOS_OUTPUT_ROOT", outputRoot)
	if err := markMacOSOutputLocal(outputRoot); err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(project, "bin")
	if err := os.Mkdir(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	artifact := filepath.Join(bin, "artifact")
	if err := os.WriteFile(artifact, []byte("original"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := relocateMacOSOutput(project, bin); err != nil {
		t.Fatal(err)
	}
	if linked, err := isOutputLink(bin); err != nil || !linked {
		t.Fatalf("output link: linked=%v, error=%v", linked, err)
	}
	if !macOSOutputIsLocal(bin) {
		t.Fatal("relocated output was not recognized through its symlink")
	}
	if data, err := os.ReadFile(artifact); err != nil || string(data) != "original" {
		t.Fatalf("artifact was not preserved: %q, %v", data, err)
	}
	target, err := os.Readlink(bin)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(bin); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(artifact, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := relocateMacOSOutput(project, bin); err == nil {
		t.Fatal("relocation overwrote an existing local output directory")
	}
	for path, want := range map[string]string{artifact: "new", filepath.Join(target, "artifact"): "original"} {
		if data, err := os.ReadFile(path); err != nil || string(data) != want {
			t.Fatalf("collision damaged %s: %q, %v", path, data, err)
		}
	}
}

func TestMacOSOutputMarkerDoesNotTrustMovedDirectory(t *testing.T) {
	root := t.TempDir()
	original := filepath.Join(root, "original")
	if err := os.Mkdir(original, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := markMacOSOutputLocal(original); err != nil {
		t.Fatal(err)
	}
	moved := filepath.Join(root, "moved")
	if err := os.Rename(original, moved); err != nil {
		t.Fatal(err)
	}
	if macOSOutputIsLocal(moved) {
		t.Fatal("moved output retained trust in its previous filesystem location")
	}
}

func TestMacOSOutputRejectsLegacyMarker(t *testing.T) {
	path := t.TempDir()
	resolved, err := resolveMacOSOutput(path)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(path, ".wails-provider-checked"), []byte(resolved), 0o644))
	require.False(t, macOSOutputIsLocal(path), "legacy probe-only result was trusted")
}

func TestOutputLinkRejectsFileTarget(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions vary on Windows")
	}
	root := t.TempDir()
	target := filepath.Join(root, "file")
	if err := os.WriteFile(target, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "bin")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if _, err := isOutputLink(link); err == nil {
		t.Fatal("a file was accepted as the build output directory")
	}
}

func TestMacOSOutputLockSerializesBuilders(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("macOS build lock")
	}
	path := filepath.Join(t.TempDir(), "bin")
	unlock, err := lockMacOSOutput(path)
	if err != nil {
		t.Fatal(err)
	}
	acquired := make(chan error, 1)
	go func() {
		release, err := lockMacOSOutput(path)
		if err == nil {
			release()
		}
		acquired <- err
	}()
	select {
	case err := <-acquired:
		unlock()
		t.Fatalf("second builder acquired the held lock: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
	unlock()
	select {
	case err := <-acquired:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("second builder did not acquire the released lock")
	}
}

func TestMacOSOutputVolumeCheck(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("macOS output volume check")
	}
	path := t.TempDir()
	if same, err := sameMacOSFilesystem(path, path); err != nil || !same {
		t.Fatalf("same-volume output: same=%v, error=%v", same, err)
	}
	if same, err := sameMacOSFilesystem(path, "/dev"); err != nil || same {
		t.Fatalf("different-volume output: same=%v, error=%v", same, err)
	}
}
