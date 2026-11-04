package updater_test

import (
	"context"
	"crypto/sha256"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/pkg/updater"
)

// workflowReproHost is a minimal in-memory Host. It never opens a window.
type workflowReproHost struct{}

func (workflowReproHost) Emit(string, ...any) bool { return false }
func (workflowReproHost) OnEvent(string, func(any)) func() {
	return func() {}
}
func (workflowReproHost) OpenWindow(updater.WindowOptions) updater.WindowHandle { return nil }
func (workflowReproHost) Quit()                                                 {}

// workflowReproCountLeftovers counts wails-update-* staging directories
// directly under root (the directory holding the target executable), and how
// many of them contain an .artifact file.
func workflowReproCountLeftovers(root string) (dirs, artifacts int, err error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return 0, 0, err
	}
	for _, e := range entries {
		if !e.IsDir() || !strings.HasPrefix(e.Name(), "wails-update-") {
			continue
		}
		dirs++
		if _, statErr := os.Stat(filepath.Join(root, e.Name(), ".artifact")); statErr == nil {
			artifacts++
		}
	}
	return dirs, artifacts, nil
}

// workflowReproProvider writes partial bytes to dst and then fails with
// io.ErrUnexpectedEOF, simulating a cut-off download stream (no networking).
type workflowReproProvider struct {
	root string
	full []byte
	part []byte

	downloadCalls   int
	written         int
	duringDirs      int
	duringArtifacts int
	duringErr       error
}

func (p *workflowReproProvider) Name() string { return "workflow-repro" }

func (p *workflowReproProvider) Check(context.Context, updater.CheckRequest) (*updater.Release, error) {
	sum := sha256.Sum256(p.full)
	return &updater.Release{
		Version:      "2.0.0",
		Artifact:     updater.Artifact{Filename: "app.bin", Size: int64(len(p.full))},
		Verification: &updater.Verification{DigestAlgo: "sha256", Digest: sum[:]},
	}, nil
}

func (p *workflowReproProvider) Download(_ context.Context, _ *updater.Release, dst io.Writer, onProgress func(int64, int64)) error {
	p.downloadCalls++
	n, err := dst.Write(p.part)
	p.written = n
	if err != nil {
		return err
	}
	if onProgress != nil {
		onProgress(int64(n), int64(len(p.full)))
	}
	// Positive control: observe the staging directory while the download is
	// still in flight, before the cut-off error is returned.
	p.duringDirs, p.duringArtifacts, p.duringErr = workflowReproCountLeftovers(p.root)
	return io.ErrUnexpectedEOF
}

func TestWorkflowReproDownloadFailureLeavesTempArtifacts(t *testing.T) {
	// The updater stages downloads beside the target executable (not under
	// os.TempDir) so the final swap is a same-filesystem rename. Point the
	// executable resolver into an isolated directory so every wails-update-*
	// staging dir this test produces lands under root.
	root := t.TempDir()
	t.Cleanup(updater.SetSelfExecutableForTest(func() (string, error) {
		return filepath.Join(root, "app"), nil
	}))

	full := []byte("0123456789abcdefghij0123456789abcdefghij")
	prov := &workflowReproProvider{root: root, full: full, part: full[:len(full)/2]}

	u := updater.New(workflowReproHost{})
	if err := u.Init(updater.Config{
		CurrentVersion: "1.0.0",
		Providers:      []updater.Provider{prov},
		Window:         updater.WindowNone,
	}); err != nil {
		t.Fatalf("setup: Init failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rel, err := u.Check(ctx)
	if err != nil {
		t.Fatalf("setup: Check failed: %v", err)
	}
	if rel == nil {
		t.Fatal("setup: Check returned no release")
	}

	dlErr := u.DownloadAndInstall(ctx)

	// Positive control: the trigger was actually reached.
	if prov.downloadCalls != 1 {
		t.Fatalf("setup: provider Download called %d times, want 1", prov.downloadCalls)
	}
	if prov.written == 0 {
		t.Fatal("setup: provider wrote no partial bytes")
	}
	if prov.duringErr != nil {
		t.Fatalf("setup: could not inspect isolated temp root during download: %v", prov.duringErr)
	}
	if prov.duringDirs == 0 || prov.duringArtifacts == 0 {
		t.Fatalf("setup: staging dir/.artifact not observed in isolated temp root during download (dirs=%d artifacts=%d); isolation or trigger not reached", prov.duringDirs, prov.duringArtifacts)
	}
	if dlErr == nil || !errors.Is(dlErr, io.ErrUnexpectedEOF) {
		t.Fatalf("setup: DownloadAndInstall error = %v, want one wrapping io.ErrUnexpectedEOF", dlErr)
	}

	dirs, artifacts, err := workflowReproCountLeftovers(root)
	if err != nil {
		t.Fatalf("setup: could not inspect isolated temp root after failure: %v", err)
	}

	t.Logf("REPRO_OBSERVATION: partial_bytes_written=%d of %d; during_download dirs=%d artifacts=%d; state_after=%s; after_failure leftover_dirs=%d leftover_artifacts=%d",
		prov.written, len(full), prov.duringDirs, prov.duringArtifacts, u.State(), dirs, artifacts)

	if dirs != 0 || artifacts != 0 {
		t.Errorf("REPRO_ASSERTION: failed download left temporary files behind: wails-update-* dirs=%d, .artifact files=%d (expected 0 and 0)", dirs, artifacts)
	}
}
