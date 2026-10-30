package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Exercise generation, release payload, tagged archive and push using only
// isolated Git repositories, local bare remotes and a local HTTP server.
func TestGeneratedNotesReachReleaseAndTag(t *testing.T) {
	script, err := filepath.Abs("../../scripts/generate-release-changelog.py")
	if err != nil {
		t.Fatal(err)
	}
	for _, tagOnly := range []bool{false, true} {
		name, version := "nightly", "v3.1.0-nightly.20261031"
		if tagOnly {
			name, version = "weekly", "v3.1.0"
		}
		t.Run(name, func(t *testing.T) {
			cleanup, root := setupTestEnvironment(t)
			defer cleanup()
			write := func(path, content string) {
				t.Helper()
				path = filepath.Join(root, path)
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
			}
			write(".gitignore", "/v3/release_changelog.md\n/v3/release_changelog.json\n")
			write("v3/internal/version/version.txt", "v3.0.0")
			write(runtimePackageJSON, "{\n  \"version\": \"3.0.0\"\n}\n")
			write(runtimePackageLock, "{\n  \"version\": \"3.0.0\",\n  \"packages\": {\n    \"\": {\n      \"version\": \"3.0.0\"\n    }\n  }\n}\n")
			write("docs/mpress/content/changelog.md", "# Changelog\n\n## [Unreleased]\n\n## v3.0.0 - 2026-10-01\n\n## Fixed\n- Historical fix\n")
			gitIn(t, root, "init", "-q", "-b", "master")
			gitIn(t, root, "add", ".")
			gitIn(t, root, "commit", "-q", "-m", "baseline")
			gitIn(t, root, "tag", "v3.0.0")
			write("v3/change.go", "package fixture\n")
			gitIn(t, root, "add", "v3/change.go")
			gitIn(t, root, "commit", "-q", "-m", "fix: Correct fixture size")
			source, err := newGitRunner(root).runCapture("rev-parse", "HEAD")
			if err != nil {
				t.Fatal(err)
			}
			generator := `import importlib.util, pathlib, sys
spec = importlib.util.spec_from_file_location("generator", sys.argv[1])
m = importlib.util.module_from_spec(spec)
spec.loader.exec_module(m)
root, sha = pathlib.Path(sys.argv[2]), sys.argv[3]
def api(endpoint):
    if "/pulls?" in endpoint:
        return [{"number": 42, "title": "fix: Correct fixture size", "merged_at": "2026-10-31", "merge_commit_sha": sha, "base": {"repo": {"full_name": "fixture/private"}}, "user": {"login": "fixture"}}]
    return [{"filename": "v3/change.go"}]
m.generate(root, "fixture/private", sha, sys.argv[4], api=api)
`
			command := exec.Command("python3", "-c", generator, script, root, source, name)
			if output, err := command.CombinedOutput(); err != nil {
				t.Fatalf("generator: %v\n%s", err, output)
			}
			remote := filepath.Join(t.TempDir(), "release.git")
			gitIn(t, root, "init", "--bare", "-q", remote)
			gitIn(t, root, "config", "url."+remote+".insteadOf", buildAuthURL("fixture/private", "test"))
			payloads := make(chan map[string]any, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost && r.URL.Path == "/repos/fixture/private/releases" {
					var payload map[string]any
					if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
						t.Error(err)
					}
					payloads <- payload
					w.WriteHeader(http.StatusCreated)
					_, _ = w.Write([]byte(`{"id":1,"html_url":"https://example.invalid/release"}`))
					return
				}
				if r.Method == http.MethodGet && r.URL.Path == "/repos/fixture/private" {
					_, _ = w.Write([]byte(`{}`))
					return
				}
				t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
				w.WriteHeader(http.StatusNotFound)
			}))
			defer server.Close()
			t.Setenv("GITHUB_API_URL", server.URL)
			t.Setenv("GITHUB_REPOSITORY", "fixture/private")
			t.Setenv("WAILS_REPO_TOKEN", "test")
			t.Setenv("GITHUB_TOKEN", "")
			t.Setenv("GITHUB_OUTPUT", filepath.Join(t.TempDir(), "outputs"))
			// Tampering or selecting another snapshot must fail before mutation.
			if err := validateGeneratedChangelog(root, "v3.0.0"); err == nil {
				t.Fatal("stale target accepted")
			}
			notes, err := os.ReadFile(unreleasedChangelogFile)
			if err != nil {
				t.Fatal(err)
			}
			write("v3/release_changelog.md", "## Fixed\n- Tampered\n")
			if err := validateGeneratedChangelog(root, source); err == nil {
				t.Fatal("tampered notes accepted")
			}
			write("v3/release_changelog.md", string(notes))
			if err := runRelease(releaseOptions{version: version, branch: "master", target: source,
				channel: name, tagOnly: tagOnly}); err != nil {
				t.Fatal(err)
			}
			bullet := "- Correct fixture size in [PR](https://github.com/fixture/private/pull/42) by @fixture"
			payload := <-payloads
			body, _ := payload["body"].(string)
			if !strings.Contains(body, bullet) || payload["target_commitish"] != source {
				t.Fatalf("release payload: %v", payload)
			}
			published := gitIn(t, remote, "show", "refs/tags/"+version+":docs/mpress/content/changelog.md")
			if !strings.Contains(published, bullet) || !strings.Contains(published, "Historical fix") {
				t.Fatalf("tagged archive: %s", published)
			}
			tracked := gitIn(t, root, "ls-files")
			if strings.Contains(tracked, "release_changelog") || strings.Contains(tracked, "UNRELEASED_CHANGELOG") {
				t.Fatal("release-local inputs were committed")
			}
			if input, err := os.ReadFile(unreleasedChangelogFile); err != nil || !strings.Contains(string(input), bullet) {
				t.Fatal("generated input was lost or reset")
			}
		})
	}
}
