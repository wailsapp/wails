package changelog

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type fixture struct {
	t           *testing.T
	root, start string
	responses   map[string][]json.RawMessage
}

const repo = "fixture/private"
const pulls = "repos/" + repo + "/pulls?state=closed&per_page=100"

func newFixture(t *testing.T) *fixture {
	t.Helper()
	f := &fixture{t: t, root: t.TempDir(), responses: map[string][]json.RawMessage{}}
	f.write("v3/.keep", "")
	f.git("init", "-q", "-b", "master")
	f.git("config", "user.name", "Fixture")
	f.git("config", "user.email", "fixture@example.invalid")
	f.start = f.commit("baseline")
	f.git("tag", "v3.0.0")
	return f
}

func (f *fixture) git(args ...string) string {
	f.t.Helper()
	result, err := git(f.root, args...)
	if err != nil {
		f.t.Fatal(err)
	}
	return result
}

func (f *fixture) commit(message string) string {
	f.git("commit", "-q", "--allow-empty", "-m", message)
	return f.git("rev-parse", "HEAD")
}

func (f *fixture) write(name, content string) {
	f.t.Helper()
	path := filepath.Join(f.root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		f.t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		f.t.Fatal(err)
	}
}

func (f *fixture) notes() string {
	f.t.Helper()
	data, err := os.ReadFile(filepath.Join(f.root, "v3/release_changelog.md"))
	if err != nil {
		f.t.Fatal(err)
	}
	return string(data)
}

func raw(v any) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func filesPath(n int) string { return fmt.Sprintf("repos/%s/pulls/%d/files?per_page=100", repo, n) }

func (f *fixture) pr(n int, sha, title string, extra map[string]any) {
	pr := map[string]any{"number": n, "merge_commit_sha": sha, "merged_at": "2026-10-31", "title": title, "base": map[string]any{"repo": map[string]string{"full_name": repo}}, "user": map[string]string{"login": "fixture"}}
	for k, v := range extra {
		pr[k] = v
	}
	f.responses[pulls] = append(f.responses[pulls], raw(pr))
	f.responses[filesPath(n)] = []json.RawMessage{raw(map[string]string{"filename": "v3/window.go"})}
}

func (f *fixture) api(endpoint string) ([]json.RawMessage, error) { return f.responses[endpoint], nil }

func (f *fixture) generate(options Options) *Receipt {
	f.t.Helper()
	options.Repository = repo
	r, err := Generate(f.root, options, f.api)
	if err != nil {
		f.t.Fatal(err)
	}
	return r
}

func numbers(r *Receipt) []int {
	result := []int{}
	for _, e := range r.Entries {
		result = append(result, e.Number)
	}
	return result
}

func equal(t *testing.T, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v; want %#v", got, want)
	}
}

func fails(t *testing.T, err error, contains string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), contains) {
		t.Fatalf("expected error containing %q, got %v", contains, err)
	}
}

func TestBoundaries(t *testing.T) {
	f := newFixture(t)
	f.pr(1, f.commit("fix"), "fix: Correct size", nil)
	f.git("tag", "v3.1.0-nightly.20261030")
	f.pr(2, f.commit("feature"), "feat(window): Add feature", nil)
	equal(t, numbers(f.generate(Options{})), []int{2})
	weekly := f.generate(Options{Channel: "weekly"})
	equal(t, numbers(weekly), []int{1, 2})
	equal(t, weekly.Baseline, f.start)
	equal(t, f.generate(Options{Version: "v3.1.1"}).Entries, weekly.Entries)
}

func TestMergeMembershipAndDeduplication(t *testing.T) {
	f := newFixture(t)
	f.git("switch", "-q", "-c", "feature")
	f.commit("feature")
	f.git("switch", "-q", "master")
	f.git("merge", "-q", "--no-ff", "feature", "-m", "merge")
	f.pr(2, f.git("rev-parse", "HEAD"), "feat: Feature", nil)
	f.responses[pulls] = append(f.responses[pulls], f.responses[pulls][0])
	first := f.commit("rebased first")
	last := f.commit("rebased last")
	f.pr(3, last, "fix: Rebased", nil)
	f.responses[pulls] = append(f.responses[pulls], f.responses[pulls][2])
	f.pr(4, first, "fix: Open", map[string]any{"merged_at": nil})
	f.pr(5, strings.Repeat("f", 40), "fix: Later", nil)
	f.pr(6, first, "fix: Foreign", map[string]any{"base": map[string]any{"repo": map[string]string{"full_name": "other/repo"}}})
	equal(t, numbers(f.generate(Options{})), []int{2, 3})
	equal(t, strings.Count(f.notes(), "/pull/2)"), 1)
}

func TestCategoriesAndScope(t *testing.T) {
	f := newFixture(t)
	sha := f.commit("changes")
	f.pr(1, sha, "chore: Bookkeeping", nil)
	f.pr(2, sha, "fix: V2", nil)
	f.responses[filesPath(2)] = []json.RawMessage{raw(map[string]string{"filename": "v2/file.go"})}
	f.pr(3, sha, "fix: Hidden", map[string]any{"labels": []map[string]string{{"name": "skip-changelog"}}})
	r := f.generate(Options{})
	equal(t, numbers(r), []int{})
	equal(t, r.SkippedPRs, []int{1, 2, 3})
	equal(t, f.notes(), "# Unreleased Changes\n\n")
	f.pr(4, sha, "chore: Vulnerability", map[string]any{"labels": []map[string]string{{"name": "changelog:security"}}})
	equal(t, f.generate(Options{}).Entries[0].Section, "Security")
}

func TestInvalidBoundaries(t *testing.T) {
	f := newFixture(t)
	sha := f.commit("fix")
	_, err := Generate(f.root, Options{Repository: repo, Target: f.start}, f.api)
	fails(t, err, "tested HEAD")
	f.git("tag", "-d", "v3.0.0")
	_, err = Generate(f.root, Options{Repository: repo}, f.api)
	fails(t, err, "bootstrap")
	equal(t, f.generate(Options{Since: f.start}).Target, sha)
	f.git("switch", "-q", "--orphan", "unrelated")
	f.commit("unrelated")
	_, err = Generate(f.root, Options{Repository: repo, Since: sha}, f.api)
	fails(t, err, "merge-base")
}

func TestFailuresInvalidateOutputs(t *testing.T) {
	f := newFixture(t)
	sha := f.commit("fix")
	f.pr(1, sha, "fix: Size", nil)
	f.generate(Options{})
	_, err := Generate(f.root, Options{Repository: repo}, func(string) ([]json.RawMessage, error) { return nil, fmt.Errorf("GitHub unavailable") })
	fails(t, err, "unavailable")
	for _, name := range []string{"md", "json"} {
		if _, err := os.Stat(filepath.Join(f.root, "v3/release_changelog."+name)); !os.IsNotExist(err) {
			t.Fatalf("stale %s output", name)
		}
	}
	f.pr(2, sha, "fix: Conflicting", map[string]any{"labels": []map[string]string{{"name": "changelog:fixed"}, {"name": "changelog:added"}}})
	_, err = Generate(f.root, Options{Repository: repo}, f.api)
	fails(t, err, "conflicting")
}

func TestFileLimit(t *testing.T) {
	f := newFixture(t)
	f.pr(1, f.commit("large"), "fix: Large", nil)
	files := make([]json.RawMessage, 3000)
	for i := range files {
		files[i] = raw(map[string]string{"filename": "v3/file.go"})
	}
	f.responses[filesPath(1)] = files
	_, err := Generate(f.root, Options{Repository: repo}, f.api)
	fails(t, err, "limit")
}

func TestRenderingAndReceipt(t *testing.T) {
	f := newFixture(t)
	f.pr(1, f.commit("fix"), "fix: Correct <script>\n## injected heading", nil)
	first := f.generate(Options{})
	notes := f.notes()
	if !strings.Contains(notes, "&lt;script&gt; ## injected heading") {
		t.Fatal(notes)
	}
	equal(t, f.generate(Options{}), first)
	equal(t, f.notes(), notes)
	data, err := os.ReadFile(filepath.Join(f.root, "v3/release_changelog.json"))
	if err != nil {
		t.Fatal(err)
	}
	var receipt Receipt
	if err := json.Unmarshal(data, &receipt); err != nil {
		t.Fatal(err)
	}
	equal(t, &receipt, first)
}

func TestGitHubPagination(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX gh fixture exercises the Linux release command")
	}
	dir := t.TempDir()
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	path := filepath.Join(dir, "gh")
	for _, tc := range []struct {
		output string
		valid  bool
	}{{`[[{"number":1}],[{"number":2}]]`, true}, {`[{"message":"failure"}]`, false}, {`null`, false}, {`[null]`, false}} {
		script := "#!/bin/sh\ncase \" $* \" in *' --paginate --slurp '*) ;; *) exit 1;; esac\nprintf '%s\\n' '" + tc.output + "'\n"
		if err := os.WriteFile(path, []byte(script), 0755); err != nil {
			t.Fatal(err)
		}
		items, err := GitHubPages("repos/fixture/private/example")
		if tc.valid {
			if err != nil {
				t.Fatal(err)
			}
			equal(t, len(items), 2)
		} else {
			fails(t, err, "paginated lists")
		}
	}
}

func TestDescriptionAndAllCategories(t *testing.T) {
	f := newFixture(t)
	sha := f.commit("categories")
	for i, kind := range []string{"feat", "perf", "fix", "deprecate", "remove", "security"} {
		f.pr(i+1, sha, kind+": Summary", nil)
	}
	f.pr(7, sha, "fix: Generic title", map[string]any{"body": "## Description\nIgnore this\n## Release notes\n<!-- ignored -->\n- User-visible fix.\n## Tests\nIgnore tests"})
	r := f.generate(Options{})
	for i, section := range sections {
		equal(t, r.Entries[i].Section, section)
	}
	equal(t, r.Entries[6].Title, "User-visible fix.")
	f.pr(8, sha, "fix: CRLF generic title", map[string]any{
		"body": "## Release notes\r\n<!-- ignored -->\r\n- Windows-edited summary.\r\n## Tests\r\nIgnore tests",
	})
	equal(t, f.generate(Options{}).Entries[7].Title, "Windows-edited summary.")
}

func TestRenameOutOfV3(t *testing.T) {
	f := newFixture(t)
	f.pr(1, f.commit("move"), "remove: Legacy tool", nil)
	f.responses[filesPath(1)] = []json.RawMessage{raw(map[string]string{"filename": "tools/legacy.go", "previous_filename": "v3/legacy.go"})}
	equal(t, f.generate(Options{}).Entries[0].Section, "Removed")
}

func TestUnreachableAndUnrelatedTags(t *testing.T) {
	f := newFixture(t)
	f.pr(1, f.commit("fix"), "fix: Size", nil)
	f.git("tag", "v3.999.0-tui")
	f.git("tag", "v3.not-a-release")
	equal(t, f.generate(Options{}).BaselineRef, "v3.0.0")
	f.git("switch", "-q", "--orphan", "other")
	f.commit("other")
	f.git("tag", "v3.999.0")
	f.git("switch", "-q", "master")
	equal(t, f.generate(Options{Channel: "weekly"}).BaselineRef, "v3.0.0")
}

func TestMigrationScopeAndExpiry(t *testing.T) {
	f := newFixture(t)
	sha := f.commit("pending")
	f.pr(1, sha, "fix: Different wording", nil)
	entry := Entry{1, "Fixed", "Original approved pending entry", "fixture", sha}
	f.write("v3/tasks/release/changelog-migration.json", string(raw(map[string]any{"repository": repo, "entries": []Entry{entry}})))
	equal(t, f.generate(Options{}).Entries, []Entry{entry})
	equal(t, strings.Count(f.notes(), entry.Title), 1)
	other, err := migratePending(f.root, "other/repo", []string{sha}, []Entry{})
	if err != nil {
		t.Fatal(err)
	}
	equal(t, other, []Entry{})
	f.git("tag", "v3.1.0-nightly.20261030")
	f.commit("bookkeeping")
	equal(t, numbers(f.generate(Options{})), []int{})
	equal(t, f.generate(Options{Channel: "weekly"}).Entries, []Entry{entry})
}

func TestArchiveSync(t *testing.T) {
	f := newFixture(t)
	header := "Header\n\n## [Unreleased]\n\n"
	old := "## v3.0.0 - 2026-10-01\n\n## Fixed\n- Old fix\n"
	newSection := "## v3.1.0 - 2026-10-31\n\n## Added\n- Weekly feature\n\n"
	f.write(Archive, header+newSection+old)
	f.git("add", Archive)
	f.git("commit", "-q", "-m", "weekly")
	f.git("tag", "v3.1.0")
	concurrent := "## v3.2.0-nightly.20261101 - 2026-11-01\n\n## Fixed\n- Concurrent fix\n\n"
	f.write(Archive, header+concurrent+old)
	if err := SyncArchive(f.root, "v3.1.0"); err != nil {
		t.Fatal(err)
	}
	read := func() string {
		data, err := os.ReadFile(filepath.Join(f.root, Archive))
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}
	expected := read()
	for _, text := range []string{strings.TrimSpace(newSection), strings.TrimSpace(concurrent), old} {
		if !strings.Contains(expected, text) {
			t.Fatal(expected)
		}
	}
	if err := SyncArchive(f.root, "v3.1.0"); err != nil {
		t.Fatal(err)
	}
	equal(t, read(), expected)
	f.write(Archive, strings.ReplaceAll(expected, "Weekly feature", "Conflict"))
	fails(t, SyncArchive(f.root, "v3.1.0"), "differs")
	f.git("tag", "v3.2.0")
	fails(t, SyncArchive(f.root, "v3.2.0"), "no release changelog section")
}

func TestCommandLine(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX gh fixture exercises the Linux release command")
	}
	binary := filepath.Join(t.TempDir(), "changelog")
	build := exec.Command("go", "build", "-o", binary, "../../scripts/generate-release-changelog.go")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	f := newFixture(t)
	sha := f.commit("fix")
	f.pr(42, sha, "fix: Command-line fixture", nil)
	bin := t.TempDir()
	script := "#!/bin/sh\ncase \"$*\" in *'/files?'*) printf '%s\\n' '[ [{\"filename\":\"v3/window.go\"}] ]';; *) printf '%s\\n' '" + string(raw([][]json.RawMessage{f.responses[pulls]})) + "';; esac\n"
	if err := os.WriteFile(filepath.Join(bin, "gh"), []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("GITHUB_REPOSITORY", repo)
	cmd := exec.Command(binary, "generate", "--target", sha, "--channel", "nightly", "--version", "")
	cmd.Dir = f.root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if !strings.Contains(f.notes(), "Command-line fixture in [PR](https://github.com/fixture/private/pull/42)") {
		t.Fatal(f.notes())
	}
	f.write(Archive, "# Changelog\n\n## [Unreleased]\n\n## v3.1.0 - 2026-10-31\n\n## Fixed\n- Fixture\n")
	f.git("add", Archive)
	f.git("commit", "-q", "-m", "release")
	f.git("tag", "v3.1.0")
	f.write(Archive, "# Changelog\n\n## [Unreleased]\n")
	cmd = exec.Command(binary, "sync", "--tag", "v3.1.0")
	cmd.Dir = f.root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("sync: %v\n%s", err, out)
	}
	data, err := os.ReadFile(filepath.Join(f.root, Archive))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "- Fixture") {
		t.Fatal(string(data))
	}
	for _, args := range [][]string{{}, {"invalid"}, {"generate", "--channel", "invalid"}, {"sync"}, {"generate", "extra"}} {
		cmd = exec.Command(binary, args...)
		cmd.Dir = f.root
		if out, err := cmd.CombinedOutput(); err == nil {
			t.Fatalf("accepted %v: %s", args, out)
		}
	}
}
