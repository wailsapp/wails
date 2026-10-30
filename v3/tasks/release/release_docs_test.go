package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestDocumentationFailureStopsReleaseBeforeOutputs(t *testing.T) {
	called := false
	output, err := dryRun(t, "v3.0.0", releaseOptions{
		version: "v3.0.1",
		prepareDocs: func(version string, dry bool) error {
			called = true
			if version != "v3.0.1" || !dry {
				t.Fatalf("unexpected documentation preparation: %q, %v", version, dry)
			}
			return fmt.Errorf("broken documentation link")
		},
	}, true)
	if !called || err == nil || !strings.Contains(err.Error(), "before release commit/tag") {
		t.Fatalf("documentation failure did not veto release: called=%v err=%v", called, err)
	}
	if strings.Contains(output, "release_tag=") || strings.Contains(output, "release_outcome=success") {
		t.Fatalf("failed documentation validation emitted release outputs: %s", output)
	}
}

func TestPrepareReleaseDocumentationPassesDryRunToPython(t *testing.T) {
	root := t.TempDir()
	directory := filepath.Join(root, "docs", "mpress", "scripts")
	if err := os.MkdirAll(directory, 0755); err != nil {
		t.Fatal(err)
	}
	script := `import pathlib, sys
pathlib.Path("arguments").write_text("\n".join(sys.argv[1:]))
`
	if err := os.WriteFile(filepath.Join(directory, "release_docs.py"), []byte(script), 0644); err != nil {
		t.Fatal(err)
	}
	if err := prepareReleaseDocumentation(root, "v3.0.1", true); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "arguments"))
	if err != nil {
		t.Fatal(err)
	}
	want := "prepare\n--root\n" + root + "\n--version\nv3.0.1\n--dry-run"
	if string(data) != want {
		t.Fatalf("arguments = %q; want %q", data, want)
	}
}

func TestReleasePushPublishesBranchAndTagAtomically(t *testing.T) {
	root := t.TempDir()
	revision := func(directory, ref string) string {
		t.Helper()
		output, err := exec.Command("git", "-C", directory, "rev-parse", ref).Output()
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimSpace(string(output))
	}
	gitIn(t, root, "init", "-q")
	gitIn(t, root, "commit", "-q", "--allow-empty", "-m", "release")
	gitIn(t, root, "tag", "v3.0.1")
	remote := filepath.Join(t.TempDir(), "remote.git")
	gitIn(t, root, "init", "--bare", "-q", remote)
	// Redirect only the fixture URL to a local bare repository; no network access.
	gitIn(t, root, "config", "url."+remote+".insteadOf", buildAuthURL("fixture/release", "test"))
	if err := newGitRunner(root).pushRelease("master", "v3.0.1", "fixture/release", "test", false); err != nil {
		t.Fatal(err)
	}
	if got := revision(remote, "master"); got != revision(remote, "v3.0.1") {
		t.Fatal("release refs disagree")
	}
	// A rejected tag must also reject the branch update.
	gitIn(t, root, "commit", "-q", "--allow-empty", "-m", "next")
	gitIn(t, root, "tag", "-f", "v3.0.1")
	before := revision(remote, "master")
	if err := newGitRunner(root).pushRelease("master", "v3.0.1", "fixture/release", "test", false); err == nil {
		t.Fatal("conflicting release tag accepted")
	}
	if got := revision(remote, "master"); got != before {
		t.Fatal("branch advanced despite rejected tag")
	}
}
