//go:build !windows

// These tests execute the Bash scripts used by the Linux release workflows.
package releaseci

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type step struct {
	Name            string            `yaml:"name"`
	Uses            string            `yaml:"uses"`
	Run             string            `yaml:"run"`
	With            map[string]any    `yaml:"with"`
	Env             map[string]string `yaml:"env"`
	ContinueOnError bool              `yaml:"continue-on-error"`
}

type job struct {
	Steps []step         `yaml:"steps"`
	Needs dependencies   `yaml:"needs"`
	With  map[string]any `yaml:"with"`
}

type dependencies []string

func (d *dependencies) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		*d = []string{node.Value}
		return nil
	}
	var values []string
	if err := node.Decode(&values); err != nil {
		return err
	}
	*d = values
	return nil
}

func root(t *testing.T) string {
	t.Helper()
	path, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func workflow(t *testing.T, name string) map[string]job {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root(t), ".github/workflows", name+".yml"))
	if err != nil {
		t.Fatal(err)
	}
	var w struct {
		Jobs map[string]job `yaml:"jobs"`
	}
	if err := yaml.Unmarshal(data, &w); err != nil {
		t.Fatal(err)
	}
	return w.Jobs
}

func named(t *testing.T, j job, name string) step {
	t.Helper()
	for _, s := range j.Steps {
		if s.Name == name {
			return s
		}
	}
	t.Fatalf("missing step %q", name)
	return step{}
}

func equal(t *testing.T, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v; want %#v", got, want)
	}
}

func contains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("missing %q in %s", want, got)
	}
}

func run(dir string, env map[string]string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func success(t *testing.T, dir string, env map[string]string, name string, args ...string) string {
	t.Helper()
	out, err := run(dir, env, name, args...)
	if err != nil {
		t.Fatalf("%s: %v\n%s", name, err, out)
	}
	return strings.TrimSpace(out)
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func read(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(data))
}

func TestPRScopeUsesAuthorPermissions(t *testing.T) {
	script := named(t, workflow(t, "build-and-test-v3")["changes"], "Choose CI scope from PR author permissions").With["script"].(string)
	harness := `
const assert = require('node:assert/strict');
const context = {eventName: 'pull_request', actor: 'maintainer-rerunning-ci',
  repo: {owner: 'wailsapp', repo: 'wails'},
  payload: {pull_request: {user: {login: 'pr-author'}}}};
const core = {setOutput: (k, v) => console.log(JSON.stringify(v)), warning: () => {}, info: () => {}};
const github = {rest: {repos: {getCollaboratorPermissionLevel: async (args) => {
  assert.equal(args.username, 'pr-author');
  if (process.env.TEST_PERMISSION === 'error') throw new Error('API unavailable');
  return {data: {permission: process.env.TEST_PERMISSION}};
}}}};
`
	for _, permission := range []string{"admin", "maintain", "write", "read", "triage", "none", "error"} {
		t.Run(permission, func(t *testing.T) {
			out := success(t, "", map[string]string{"TEST_PERMISSION": permission}, "node", "-e", harness+"(async () => {\n"+script+"\n})()")
			var result bool
			if err := json.Unmarshal([]byte(out), &result); err != nil {
				t.Fatal(err)
			}
			equal(t, result, !slices.Contains([]string{"admin", "maintain", "write"}, permission))
		})
	}
}

func TestReleaseRunsIgnorePRPathFilters(t *testing.T) {
	for _, tc := range []struct{ script, expected string }{{named(t, workflow(t, "build-and-test-v3")["changes"], "Check for changes under v3/").Run, "v3=true"}, {named(t, workflow(t, "verify-runtime-assets")["verify"], "Detect runtime asset or source changes").Run, "relevant=true"}} {
		for _, event := range []string{"schedule", "workflow_dispatch"} {
			dir := t.TempDir()
			output := filepath.Join(dir, "output")
			success(t, dir, map[string]string{"EVENT_NAME": event, "GITHUB_OUTPUT": output}, "bash", "-e", "-c", tc.script)
			equal(t, read(t, output), tc.expected)
		}
	}
}

func TestPRBranchesOnlyAllowDryRuns(t *testing.T) {
	script := named(t, workflow(t, "nightly-release-v3")["cadence"], "Resolve cadence and branch").Run
	for _, tc := range []struct {
		dry, ref, version string
		allowed           bool
	}{{"true", "codex/nightly-ci-gate", "", true}, {"false", "codex/nightly-ci-gate", "", false}, {"false", "master", "", true}, {"false", "releases/3.1", "", false}, {"false", "releases/3.1", "v3.1.1", true}} {
		_, err := run("", map[string]string{"EVENT_SCHEDULE": "", "INPUT_CHANNEL": "nightly", "INPUT_REF": tc.ref, "INPUT_VERSION": tc.version, "INPUT_DRY_RUN": tc.dry, "GITHUB_OUTPUT": filepath.Join(t.TempDir(), "output")}, "bash", "-e", "-c", script)
		equal(t, err == nil, tc.allowed)
	}
}

func TestResultsRejectIncompleteChecks(t *testing.T) {
	script := workflow(t, "build-and-test-v3")["build_results"].Steps[0].Run
	checks := []string{"test_go", "test_js", "test_templates", "test_mac_private_apis", "cross-compile", "gtk3-menus", "msix"}
	values := map[string]string{"needs.changes.result": "success", "needs.changes.outputs.v3": "true", "needs.changes.outputs.full": "true"}
	for _, check := range checks {
		values["needs."+check+".result"] = "success"
	}
	pattern := regexp.MustCompile(`\$\{\{\s*(.*?)\s*\}\}`)
	passed := func(overrides map[string]string) bool {
		rendered := pattern.ReplaceAllStringFunc(script, func(match string) string {
			key := strings.TrimSpace(pattern.FindStringSubmatch(match)[1])
			if value, ok := overrides[key]; ok {
				return value
			}
			value, ok := values[key]
			if !ok {
				t.Fatalf("unhandled expression: %s", key)
			}
			return value
		})
		_, err := run("", nil, "bash", "-e", "-c", rendered)
		return err == nil
	}
	equal(t, passed(nil), true)
	for _, check := range append([]string{"changes"}, checks...) {
		for _, status := range []string{"failure", "cancelled", "skipped"} {
			equal(t, passed(map[string]string{"needs." + check + ".result": status}), false)
		}
	}
	quick := map[string]string{"needs.changes.outputs.full": "false"}
	for _, check := range checks[2:] {
		quick["needs."+check+".result"] = "skipped"
	}
	equal(t, passed(quick), true)
	quick["needs.test_go.result"] = "skipped"
	equal(t, passed(quick), false)
}

func TestChecksShareSnapshot(t *testing.T) {
	release := workflow(t, "nightly-release-v3")
	build := workflow(t, "build-and-test-v3")
	snapshot := "${{ needs.cadence.outputs.sha }}"
	needs := slices.Clone([]string(release["nightly-release"].Needs))
	slices.Sort(needs)
	equal(t, needs, []string{"cadence", "changelog-tests", "full-ci", "runtime-assets", "security-scan"})
	equal(t, release["full-ci"].With["full"], true)
	for _, name := range []string{"full-ci", "runtime-assets", "changelog-tests", "security-scan"} {
		equal(t, release[name].With["ref"], snapshot)
	}
	equal(t, named(t, release["nightly-release"], "Checkout code").With["ref"], snapshot)
	for _, name := range []string{"build-and-test-v3", "cross-compile-test-v3", "gtk3-menu-test-v3", "msix-smoke-v3", "verify-runtime-assets", "changelog-publication-tests", "semgrep"} {
		for _, j := range workflow(t, name) {
			for _, s := range j.Steps {
				if strings.HasPrefix(s.Uses, "actions/checkout@") {
					equal(t, s.With["ref"], "${{ inputs.ref || github.sha }}")
				}
			}
		}
	}
	for _, name := range []string{"cross-compile", "gtk3-menus", "msix"} {
		equal(t, slices.Contains(build["build_results"].Needs, name), true)
		equal(t, build[name].With["ref"], "${{ inputs.ref || github.sha }}")
	}
}

func TestNotesGeneratedBeforePublication(t *testing.T) {
	j := workflow(t, "nightly-release-v3")["nightly-release"]
	generate := named(t, j, "Generate changelog from merged PRs in the tested snapshot")
	equal(t, generate.Env["TESTED_SHA"], "${{ needs.cadence.outputs.sha }}")
	contains(t, generate.Run, `go -C v3 run ./scripts/generate-release-changelog.go generate`)
	contains(t, generate.Run, `--target "$TESTED_SHA"`)
	contains(t, generate.Run, `--version "$RELEASE_VERSION"`)
	index := func(name string) int { return slices.IndexFunc(j.Steps, func(s step) bool { return s.Name == name }) }
	for _, name := range []string{"Check for unreleased changelog content", "Run release script"} {
		equal(t, index(generate.Name) < index(name), true)
	}
	equal(t, generate.ContinueOnError, false)
	contains(t, named(t, j, "Sync release commit back to master").Run, `go -C v3 run ./scripts/generate-release-changelog.go sync --tag "$TAG"`)
	for _, path := range []string{".github/workflows/auto-changelog-v3.yml", ".github/workflows/drop-unreleased-changelog.yml", "v3/UNRELEASED_CHANGELOG.md"} {
		if _, err := os.Stat(filepath.Join(root(t), path)); !os.IsNotExist(err) {
			t.Fatalf("retired file exists: %s", path)
		}
	}
}

func TestRetiredFileRejected(t *testing.T) {
	script := named(t, workflow(t, "changelog-publication-tests")["publication"], "Reject the retired shared changelog file").Run
	dir := t.TempDir()
	success(t, dir, nil, "git", "init", "-q")
	success(t, dir, nil, "bash", "-e", "-c", script)
	write(t, filepath.Join(dir, "v3/UNRELEASED_CHANGELOG.md"), "recreated")
	success(t, dir, nil, "git", "add", ".")
	out, err := run(dir, nil, "bash", "-e", "-c", script)
	if err == nil {
		t.Fatal("recreated file accepted")
	}
	contains(t, out, "DO NOT recreate")
}

func TestBranchMovementRejectsStaleSnapshot(t *testing.T) {
	script := named(t, workflow(t, "nightly-release-v3")["nightly-release"], "Refuse a stale release snapshot").Run
	dir := t.TempDir()
	remote, source, release := filepath.Join(dir, "remote.git"), filepath.Join(dir, "source"), filepath.Join(dir, "release")
	git := func(dir string, args ...string) string { return success(t, dir, nil, "git", args...) }
	git(dir, "init", "--bare", remote)
	git(dir, "init", "-b", "master", source)
	git(source, "config", "user.name", "CI test")
	git(source, "config", "user.email", "ci@example.invalid")
	write(t, filepath.Join(source, "code"), "tested\n")
	git(source, "add", ".")
	git(source, "commit", "-m", "tested snapshot")
	git(source, "remote", "add", "origin", remote)
	git(source, "push", "origin", "master")
	git(dir, "clone", "--branch", "master", remote, release)
	sha := git(release, "rev-parse", "HEAD")
	git(release, "checkout", "--detach", sha)
	env := map[string]string{"TESTED_SHA": sha, "RELEASE_REF": "master"}
	success(t, release, env, "bash", "-e", "-c", script)
	write(t, filepath.Join(source, "code"), "untested\n")
	git(source, "commit", "-am", "concurrent merge")
	git(source, "push", "origin", "master")
	if _, err := run(release, env, "bash", "-e", "-c", script); err == nil {
		t.Fatal("stale snapshot accepted")
	}
	equal(t, git(release, "rev-parse", "HEAD"), sha)
	git(release, "config", "user.name", "CI test")
	git(release, "config", "user.email", "ci@example.invalid")
	write(t, filepath.Join(release, "version"), "metadata\n")
	git(release, "add", ".")
	git(release, "commit", "-m", "metadata")
	if _, err := run(release, nil, "git", "push", "origin", "HEAD:master"); err == nil {
		t.Fatal("competing push accepted")
	}
}
