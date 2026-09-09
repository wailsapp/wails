package updater

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

// Live archive-extraction tests against the wailsapp/updater-demo v2.0.0
// release artifacts. The release ships .tar.gz for darwin+linux and .zip for
// windows, all produced by a real release pipeline (goreleaser) — exercising
// the extractor against archives the in-tree tar/zip writers may not produce
// the same way.
//
// Gated behind WAILS_UPDATER_LIVE=1 so default `go test` doesn't hit the
// network. Run with:
//   WAILS_UPDATER_LIVE=1 go test -count=1 -run LiveArchive ./pkg/updater/...

const liveTestEnv = "WAILS_UPDATER_LIVE"

func skipUnlessLive(t *testing.T) {
	t.Helper()
	if os.Getenv(liveTestEnv) != "1" {
		t.Skipf("set %s=1 to run live network-dependent tests", liveTestEnv)
	}
}

type liveArchive struct {
	name        string
	url         string
	want256     string
	entryPrefix string // expected basename of extracted top-level entry
}

func TestLiveArchive_DarwinTarGz(t *testing.T) {
	skipUnlessLive(t)
	runLiveArchive(t, liveArchive{
		name:        "updater-demo_darwin_arm64.tar.gz",
		url:         "https://github.com/wailsapp/updater-demo/releases/download/v2.0.0/updater-demo_darwin_arm64.tar.gz",
		want256:     "5e57a9c150da871562ed3ff0b924013c35770ff8d06613c1de75ce2fdd65a2c7",
		entryPrefix: "updater-demo",
	})
}

func TestLiveArchive_LinuxTarGz(t *testing.T) {
	skipUnlessLive(t)
	runLiveArchive(t, liveArchive{
		name:        "updater-demo_linux_amd64.tar.gz",
		url:         "https://github.com/wailsapp/updater-demo/releases/download/v2.0.0/updater-demo_linux_amd64.tar.gz",
		want256:     "154b7c1766070318fc3253e025a5f08d05c71a0a05b50a3c0e9bc80147f41fc2",
		entryPrefix: "updater-demo",
	})
}

func TestLiveArchive_WindowsZip(t *testing.T) {
	skipUnlessLive(t)
	runLiveArchive(t, liveArchive{
		name:        "updater-demo_windows_amd64.zip",
		url:         "https://github.com/wailsapp/updater-demo/releases/download/v2.0.0/updater-demo_windows_amd64.zip",
		want256:     "ac8aed23683c9b376c3effc1a5d0994de18fc44ca6e7b23b46a58a22bbd9a6ff",
		entryPrefix: "updater-demo",
	})
}

func runLiveArchive(t *testing.T, a liveArchive) {
	dir := t.TempDir()
	staged := filepath.Join(dir, a.name)

	resp, err := http.Get(a.url)
	if err != nil {
		t.Fatalf("download: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d for %s", resp.StatusCode, a.url)
	}
	f, err := os.Create(staged)
	if err != nil {
		t.Fatal(err)
	}
	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(f, h), resp.Body); err != nil {
		_ = f.Close()
		t.Fatal(err)
	}
	_ = f.Close()
	gotHash := hex.EncodeToString(h.Sum(nil))
	if gotHash != a.want256 {
		t.Fatalf("download SHA256: got %s want %s — release may have been re-cut", gotHash, a.want256)
	}
	t.Logf("downloaded %s (%d bytes, SHA256 matches published digest)", a.name, h.Size())

	out, did, err := maybeExtractInto(staged)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if !did {
		t.Fatalf("extractor did not recognise %s as an archive", a.name)
	}
	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("stat extracted: %v", err)
	}
	if _, err := os.Stat(staged); err == nil {
		t.Errorf("staged archive %s should have been removed after extract", staged)
	}
	if got := filepath.Base(out); len(got) < len(a.entryPrefix) || got[:len(a.entryPrefix)] != a.entryPrefix {
		t.Errorf("extracted entry name %q doesn't start with %q", got, a.entryPrefix)
	}
	t.Logf("extracted %s (%v, %d bytes)", filepath.Base(out), info.Mode(), info.Size())
}
