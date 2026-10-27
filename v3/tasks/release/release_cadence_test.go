package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"golang.org/x/mod/semver"
)

// Release model under test (see docs/mpress/content/contributing/release-policy.md):
//   weekly  vX.Y.0                   normal release, make_latest, npm latest
//   hotfix  vX.Y.Z (Z>0, manual)     normal release, make_latest, npm latest
//   nightly vX.(Y+1).0-nightly.DATE  prerelease, not latest, npm nightly
//   legacy  vX.Y.Z-beta.N            prerelease, not latest, npm latest (pre-GA)

func day(y int, m time.Month, d int) time.Time { return time.Date(y, m, d, 15, 0, 0, 0, time.UTC) }

// goLatest simulates the Go proxy: highest tag without a pre-release suffix.
func goLatest(tags []string) string { return highestOf(tags, true) }

func TestPolicyByVersionShape(t *testing.T) {
	tests := []struct {
		version string
		want    releasePolicy
	}{
		{"v3.0.0-beta.34", releasePolicy{kindPrerelease, true, "false", "latest"}},
		{"v3.0.0", releasePolicy{kindWeekly, false, "true", "latest"}},
		{"v3.1.0", releasePolicy{kindWeekly, false, "true", "latest"}},
		{"v3.1.1", releasePolicy{kindHotfix, false, "true", "latest"}},
		{"v3.1.0-nightly.20261026", releasePolicy{kindNightly, true, "false", "nightly"}},
		{"v3.1.0-nightly.20261026.2", releasePolicy{kindNightly, true, "false", "nightly"}},
	}
	for _, tt := range tests {
		got, err := policyFor(tt.version)
		if err != nil || got != tt.want {
			t.Errorf("policyFor(%q) = %+v, %v; want %+v", tt.version, got, err, tt.want)
		}
	}
	if _, err := policyFor("garbage"); err == nil {
		t.Error("invalid version must be rejected")
	}
}

// TestVersionSequence walks the documented timeline and checks tag, flags,
// npm tag, and that simulated Go @latest never selects a nightly.
func TestVersionSequence(t *testing.T) {
	tags := []string{"v3.0.0"}
	step := func(channel string, now time.Time, want string, wantKind string) {
		t.Helper()
		got, err := nextAutoVersion("v3.0.0", channel, tags, now)
		if err != nil || got != want {
			t.Fatalf("next(%s,%s) = %q, %v; want %q", channel, now.Format("2006-01-02"), got, err, want)
		}
		if err := checkMonotonic(got, highestOf(tags, false)); err != nil {
			t.Fatal(err)
		}
		pol, _ := policyFor(got)
		if pol.Kind != wantKind {
			t.Fatalf("%s kind = %s want %s", got, pol.Kind, wantKind)
		}
		t.Logf("%-28s kind=%-8s prerelease=%-5t make_latest=%-5s npm=%-6s npmver=%s", got, pol.Kind, pol.Prerelease, pol.MakeLatest, pol.NPMTag, strings.TrimPrefix(got, "v"))
		tags = append(tags, got)
		if latest := goLatest(tags); strings.Contains(latest, "nightly") {
			t.Fatalf("go @latest picked a nightly: %s", latest)
		}
	}
	// Mon-Sat nightlies preview v3.1.0; the base never moves while v3.1.0 does not exist.
	for d := 26; d <= 31; d++ {
		step("nightly", day(2026, time.October, d), fmt.Sprintf("v3.1.0-nightly.202610%d", d), kindNightly)
		if goLatest(tags) != "v3.0.0" {
			t.Fatalf("@latest = %s, want v3.0.0", goLatest(tags))
		}
	}
	// Weekly sorts above all of its own nightlies and becomes @latest.
	step("weekly", day(2026, time.November, 1), "v3.1.0", kindWeekly)
	if goLatest(tags) != "v3.1.0" {
		t.Fatalf("@latest = %s", goLatest(tags))
	}
	for _, n := range tags {
		if strings.HasPrefix(n, "v3.1.0-nightly") && semver.Compare("v3.1.0", n) <= 0 {
			t.Fatalf("v3.1.0 must sort above %s", n)
		}
	}
	// New nightlies now preview v3.2.0.
	step("nightly", day(2026, time.November, 2), "v3.2.0-nightly.20261102", kindNightly)
	step("nightly", day(2026, time.November, 3), "v3.2.0-nightly.20261103", kindNightly)
	// Same-day re-run appends .N.
	step("nightly", day(2026, time.November, 3), "v3.2.0-nightly.20261103.1", kindNightly)
	step("nightly", day(2026, time.November, 3), "v3.2.0-nightly.20261103.2", kindNightly)
	// Hotfix is an explicit manual release (not computed); after it, base is v3.1.1 yet nightlies still preview v3.2.0.
	tags = append(tags, "v3.1.1")
	if p, _ := policyFor("v3.1.1"); p.Kind != kindHotfix || p.MakeLatest != "true" {
		t.Fatalf("hotfix policy %+v", p)
	}
	if goLatest(tags) != "v3.1.1" {
		t.Fatalf("@latest after hotfix = %s", goLatest(tags))
	}
	step("nightly", day(2026, time.November, 5), "v3.2.0-nightly.20261105", kindNightly)
	// Next weekly is v3.2.0, above the hotfix and every nightly.
	step("weekly", day(2026, time.November, 8), "v3.2.0", kindWeekly)
	if goLatest(tags) != "v3.2.0" {
		t.Fatalf("@latest = %s", goLatest(tags))
	}
	step("nightly", day(2026, time.November, 9), "v3.3.0-nightly.20261109", kindNightly)
}

func TestFromV300AndBetaSeries(t *testing.T) {
	// First nightly / weekly after the GA tag.
	if got, _ := nextAutoVersion("v3.0.0", "nightly", []string{"v3.0.0"}, day(2026, 10, 26)); got != "v3.1.0-nightly.20261026" {
		t.Errorf("first nightly = %s", got)
	}
	if got, _ := nextAutoVersion("v3.0.0", "weekly", []string{"v3.0.0"}, day(2026, 11, 1)); got != "v3.1.0" {
		t.Errorf("first weekly = %s", got)
	}
	// Beta series untouched until GA: legacy increment, tags ignored if only prereleases.
	if got, _ := nextAutoVersion("v3.0.0-beta.34", "nightly", []string{"v3.0.0-beta.34"}, day(2026, 10, 26)); got != "v3.0.0-beta.35" {
		t.Errorf("beta nightly = %s", got)
	}
	if _, err := nextAutoVersion("v3.0.0-beta.34", "weekly", nil, day(2026, 10, 26)); err == nil {
		t.Error("weekly from a beta needs explicit v3.0.0")
	}
	// Stale version.txt does not matter once release tags exist.
	if got, _ := nextAutoVersion("v3.0.0-beta.34", "nightly", []string{"v3.0.0", "v3.1.0"}, day(2026, 11, 2)); got != "v3.2.0-nightly.20261102" {
		t.Errorf("stale version.txt nightly = %s", got)
	}
	if _, err := nextAutoVersion("v3.0.0", "bogus", nil, day(2026, 1, 1)); err == nil {
		t.Error("unknown channel must error")
	}
}

func TestReleaseBodiesByKind(t *testing.T) {
	cl := "### Added\n- x"
	cases := []struct {
		version   string
		wantIn    []string
		wantNotIn []string
	}{
		{"v3.1.0", []string{"## Wails v3 Weekly Release - v3.1.0", "not with features", "@v3.1.0", "pin"}, []string{"stable", "Nightly preview of"}},
		{"v3.1.1", []string{"Hotfix Release - v3.1.1"}, []string{"stable"}},
		{"v3.2.0-nightly.20261102", []string{"## Wails v3 Nightly - v3.2.0-nightly.20261102", "Nightly preview of the upcoming v3.2.0.", "not selected by `@latest`", "pin"}, []string{"Weekly Release", "stable"}},
		{"v3.0.0-beta.35", []string{"Beta Release", "Beta Warning"}, []string{"Weekly", "Nightly preview of"}},
	}
	for _, c := range cases {
		body := buildReleaseBody(c.version, cl)
		for _, w := range c.wantIn {
			if !strings.Contains(body, w) {
				t.Errorf("%s body missing %q:\n%s", c.version, w, body)
			}
		}
		for _, w := range c.wantNotIn {
			if strings.Contains(body, w) {
				t.Errorf("%s body unexpectedly contains %q:\n%s", c.version, w, body)
			}
		}
	}
}

func gitIn(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir, "-c", "user.name=t", "-c", "user.email=t@t"}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func TestTagListingAndMonotonic(t *testing.T) {
	dir := t.TempDir()
	gitIn(t, dir, "init", "-q")
	gitIn(t, dir, "commit", "-q", "--allow-empty", "-m", "x")
	for _, tag := range []string{"v3.0.0-beta.34", "v3.0.0-tui", "v2.12.0", "v3.1.0", "v3.2.0-nightly.20261102"} {
		gitIn(t, dir, "tag", tag)
	}
	tags := listReleaseTags(dir)
	if len(tags) != 3 {
		t.Fatalf("tags = %v (stray/v2 must be ignored)", tags)
	}
	if got := highestReleaseTag(dir); got != "v3.1.0" {
		t.Errorf("highest release = %q (nightly/beta must be ignored)", got)
	}
	if got := highestOf(tags, false); got != "v3.2.0-nightly.20261102" {
		t.Errorf("highest overall = %q", got)
	}
	if checkMonotonic("v3.0.0-beta.35", "v3.1.0") == nil || checkMonotonic("v3.1.0", "v3.1.0") == nil {
		t.Error("lower/duplicate must be refused")
	}
	if checkMonotonic("v3.2.0-nightly.20261103", "v3.2.0-nightly.20261102") != nil {
		t.Error("later nightly must pass")
	}
	// A weekly sorts above its own nightlies; a nightly sorts below its weekly.
	if checkMonotonic("v3.2.0", "v3.2.0-nightly.20261109") != nil || checkMonotonic("v3.2.0-nightly.20261109", "v3.2.0") == nil {
		t.Error("pre-release < release ordering broken")
	}
}

// dryRun drives runRelease end to end in dry-run mode (no network, no git writes).
func dryRun(t *testing.T, current string, opts releaseOptions, withChangelog bool, tags ...string) (string, error) {
	t.Helper()
	cleanup, root := setupTestEnvironment(t)
	defer cleanup()
	content := getUnreleasedChangelogTemplate()
	if withChangelog {
		content = strings.Replace(content, "## Added\n", "## Added\n- Add a thing\n", 1)
	}
	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	must(os.WriteFile(filepath.Join(root, "v3", "UNRELEASED_CHANGELOG.md"), []byte(content), 0o644))
	vdir := filepath.Join(root, "v3", "internal", "version")
	must(os.MkdirAll(vdir, 0o755))
	must(os.WriteFile(filepath.Join(vdir, "version.txt"), []byte(current), 0o644))
	if len(tags) > 0 {
		gitIn(t, root, "init", "-q")
		gitIn(t, root, "commit", "-q", "--allow-empty", "-m", "x")
		for _, tag := range tags {
			gitIn(t, root, "tag", tag)
		}
	}
	out := filepath.Join(t.TempDir(), "out")
	must(os.WriteFile(out, nil, 0o644))
	t.Setenv("GITHUB_OUTPUT", out)
	old := nowUTC
	nowUTC = func() time.Time { return day(2026, time.November, 2) }
	defer func() { nowUTC = old }()
	opts.dryRun = true
	opts.branch, opts.target = defaultReleaseBranch, defaultReleaseTarget
	err := runRelease(opts)
	b, _ := os.ReadFile(out)
	return string(b), err
}

func TestDryRunScenarios(t *testing.T) {
	ga := []string{"v3.0.0", "v3.1.0"}
	cases := []struct {
		name, current string
		opts          releaseOptions
		changelog     bool
		tags          []string
		want          []string
		wantErr       string
	}{
		{"beta.34 auto (pre-GA)", "v3.0.0-beta.34", releaseOptions{}, true, nil, []string{"release_version=v3.0.0-beta.35", "release_kind=prerelease", "release_prerelease=true", "release_npm_tag=latest"}, ""},
		{"beta.34 -> explicit GA v3.0.0", "v3.0.0-beta.34", releaseOptions{version: "v3.0.0"}, true, nil, []string{"release_version=v3.0.0", "release_kind=weekly", "release_prerelease=false", "release_npm_tag=latest"}, ""},
		{"v3.0.0 first nightly", "v3.0.0", releaseOptions{}, true, []string{"v3.0.0"}, []string{"release_version=v3.1.0-nightly.20261102", "release_kind=nightly", "release_prerelease=true", "release_npm_tag=nightly"}, ""},
		{"v3.1.0 nightly previews v3.2.0", "v3.1.0", releaseOptions{}, true, ga, []string{"release_version=v3.2.0-nightly.20261102", "release_kind=nightly"}, ""},
		{"same-day re-run", "v3.1.0", releaseOptions{}, true, append([]string{"v3.2.0-nightly.20261102"}, ga...), []string{"release_version=v3.2.0-nightly.20261102.1"}, ""},
		{"first weekly", "v3.0.0", releaseOptions{channel: "weekly"}, true, []string{"v3.0.0", "v3.1.0-nightly.20261031"}, []string{"release_version=v3.1.0", "release_kind=weekly", "release_prerelease=false", "release_npm_tag=latest"}, ""},
		{"weekly, empty changelog still releases", "v3.0.0", releaseOptions{channel: "weekly"}, false, []string{"v3.0.0"}, []string{"release_version=v3.1.0"}, ""},
		{"nightly, empty changelog skips", "v3.1.0", releaseOptions{}, false, ga, []string{"release_skipped=true"}, ""},
		{"hotfix v3.1.1 explicit", "v3.1.0", releaseOptions{version: "v3.1.1"}, true, ga, []string{"release_version=v3.1.1", "release_kind=hotfix", "release_prerelease=false", "release_npm_tag=latest"}, ""},
		{"weekly after hotfix", "v3.1.1", releaseOptions{channel: "weekly"}, true, append([]string{"v3.1.1"}, ga...), []string{"release_version=v3.2.0", "release_kind=weekly"}, ""},
		{"beta weekly refused", "v3.0.0-beta.34", releaseOptions{channel: "weekly"}, true, nil, nil, "explicit --version"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := dryRun(t, c.current, c.opts, c.changelog, c.tags...)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("want error containing %q, got %v", c.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			for _, w := range c.want {
				if !strings.Contains(got, w+"\n") {
					t.Errorf("output missing %q:\n%s", w, got)
				}
			}
		})
	}
}

func TestWeeklyTagOnlyIgnoresStaleCommit(t *testing.T) {
	// Green commit's version.txt is old; tags are the source of truth.
	got, err := dryRun(t, "v3.0.0-beta.34", releaseOptions{channel: "weekly", tagOnly: true}, true, "v3.0.0", "v3.1.0", "v3.2.0-nightly.20261031")
	if err != nil || !strings.Contains(got, "release_version=v3.2.0\n") {
		t.Fatalf("err=%v\n%s", err, got)
	}
}

// The npm runtime package and the Go module release share one version: the
// package.json and package-lock.json versions must equal version.txt without
// its leading "v", so a release tag, a Go module version and an npm version
// can never drift apart.
func TestRuntimePackageVersionMatchesVersionTxt(t *testing.T) {
	read := func(rel string) []byte {
		data, err := os.ReadFile(filepath.Join("..", "..", rel))
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		return data
	}
	version := strings.TrimPrefix(strings.TrimSpace(string(read("internal/version/version.txt"))), "v")
	if version == "" {
		t.Fatal("version.txt is empty")
	}
	for _, rel := range []string{
		"internal/runtime/desktop/@wailsio/runtime/package.json",
		"internal/runtime/desktop/@wailsio/runtime/package-lock.json",
	} {
		var doc struct {
			Version string `json:"version"`
		}
		if err := json.Unmarshal(read(rel), &doc); err != nil {
			t.Fatalf("parse %s: %v", rel, err)
		}
		if doc.Version != version {
			t.Errorf("%s version = %q, want %q (version.txt without the v)", rel, doc.Version, version)
		}
	}
}

func TestSyncRuntimeVersionCommandPreservesNewerDependencies(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "release")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	build := exec.Command("go", "build", "-o", binary, "release.go")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build release command: %v\n%s", err, out)
	}
	root := t.TempDir()
	releaseDir := filepath.Join(root, "v3", "tasks", "release")
	pkgDir := filepath.Join(root, "v3", "internal", "runtime", "desktop", "@wailsio", "runtime")
	for _, dir := range []string{releaseDir, pkgDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
	packageJSON := `{
  "version": "3.0.0",
  "dependencies": {"example": "1.0.1"}
}
`
	lockJSON := `{
  "version": "3.0.0",
  "lockfileVersion": 3,
  "packages": {
    "": {
      "version": "3.0.0"
    },
    "node_modules/example": {
      "version": "1.0.1",
      "integrity": "newer-integrity"
    }
  }
}
`
	fixtures := map[string]string{"package.json": packageJSON, "package-lock.json": lockJSON}
	for name, data := range fixtures {
		if err := os.WriteFile(filepath.Join(pkgDir, name), []byte(data), 0644); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command(binary, "--sync-runtime-version", "v3.1.0")
	cmd.Dir = releaseDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("sync version: %v\n%s", err, out)
	}
	for name, original := range fixtures {
		data, err := os.ReadFile(filepath.Join(pkgDir, name))
		if err != nil {
			t.Fatal(err)
		}
		want := strings.ReplaceAll(original, "3.0.0", "3.1.0")
		if string(data) != want {
			t.Errorf("%s changed more than package versions:\n%s", name, data)
		}
	}
}
