package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// gateEnv builds a stubbed `gh`: each API path maps to a fixture file
// (path with "/" -> "_", query string dropped). "FAIL" content simulates an
// API error. POST/PATCH calls are appended to calls.log.
type gateEnv struct {
	t   *testing.T
	dir string
}

func newGate(t *testing.T) *gateEnv {
	t.Helper()
	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq not installed")
	}
	g := &gateEnv{t: t, dir: t.TempDir()}
	stub := `#!/usr/bin/env bash
FIX="` + g.dir + `"
method=GET; jqexpr=""; ep=""
while [ $# -gt 0 ]; do
  case "$1" in
    api) ;;
    --paginate) ;;
    -X) method="$2"; shift ;;
    -F) shift ;;
    --jq) jqexpr="$2"; shift ;;
    -*) ;;
    *) [ -z "$ep" ] && ep="$1" ;;
  esac
  shift
done
if [ "$method" != GET ]; then echo "$method $ep" >> "$FIX/calls.log"; echo '{}'; exit 0; fi
f="$FIX/$(echo "${ep%%\?*}" | tr '/' '_')"
[ -f "$f" ] || exit 1
[ "$(cat "$f")" = FAIL ] && exit 1
if [ -n "$jqexpr" ]; then jq -r "$jqexpr" "$f"; else cat "$f"; fi
`
	if err := os.WriteFile(filepath.Join(g.dir, "gh"), []byte(stub), 0o755); err != nil {
		t.Fatal(err)
	}
	// Healthy baseline: allowlist readable, nothing open.
	g.set("orgs/wailsapp/teams/maintainers", `{"slug":"maintainers"}`)
	g.set("orgs/wailsapp/teams/maintainers/memberships/lea", `{"state":"active"}`)
	g.set("repos/wailsapp/wails/collaborators/taliesin-ai/permission", `{"permission":"maintain"}`)
	g.set("repos/wailsapp/wails/collaborators/stranger/permission", `{"permission":"read"}`)
	for _, ep := range []string{"issues?labels=P0", "issues?labels=release-blocker", "issues?labels=regression", "dependabot/alerts", "code-scanning/alerts", "secret-scanning/alerts", "security-advisories"} {
		g.set("repos/wailsapp/wails/"+strings.SplitN(ep, "?", 2)[0], "[]")
	}
	return g
}

func (g *gateEnv) set(ep, content string) {
	g.t.Helper()
	ep = strings.SplitN(ep, "?", 2)[0]
	if err := os.WriteFile(filepath.Join(g.dir, strings.ReplaceAll(ep, "/", "_")), []byte(content), 0o644); err != nil {
		g.t.Fatal(err)
	}
}

// labelled sets the open issue list for a label (all labels share the issues
// endpoint in the stub, so tests use one label at a time).
func (g *gateEnv) issues(list string) { g.set("repos/wailsapp/wails/issues", list) }
func (g *gateEnv) events(n int, ev string) {
	g.set(fmt.Sprintf("repos/wailsapp/wails/issues/%d/events", n), ev)
}

func (g *gateEnv) run(env ...string) (string, int) {
	g.t.Helper()
	script, _ := filepath.Abs("../../../.github/scripts/check-release-blockers.sh")
	cmd := exec.Command("bash", script)
	cmd.Dir = g.dir
	cmd.Env = append(os.Environ(), "GH="+filepath.Join(g.dir, "gh"), "LAST_WEEKLY_DATE=2026-11-01T00:00:00Z", "REPORT_FILE="+filepath.Join(g.dir, "report.md"))
	cmd.Env = append(cmd.Env, env...)
	out, err := cmd.CombinedOutput()
	code := 0
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
	} else if err != nil {
		g.t.Fatal(err)
	}
	return string(out), code
}

const issueP0 = `[{"number":7,"title":"Crash on start","html_url":"https://github.com/wailsapp/wails/issues/7","created_at":"2026-11-03T00:00:00Z"}]`

func labeled(actor string) string {
	return `[{"event":"labeled","label":{"name":"P0"},"actor":{"login":"` + actor + `"},"created_at":"2026-11-03T01:00:00Z"}]`
}

func TestGate_CleanPasses(t *testing.T) {
	g := newGate(t)
	if out, code := g.run(); code != 0 || !strings.Contains(out, "No blockers") {
		t.Fatalf("code=%d\n%s", code, out)
	}
}

func TestGate_P0FromAllowedActorBlocks(t *testing.T) {
	for _, actor := range []string{"lea", "taliesin-ai"} { // team member; bot collaborator with maintain
		g := newGate(t)
		g.issues(issueP0)
		g.events(7, labeled(actor))
		out, code := g.run()
		if code != 1 || !strings.Contains(out, "#7") || !strings.Contains(out, "@"+actor) {
			t.Fatalf("%s: code=%d\n%s", actor, code, out)
		}
	}
}

func TestGate_NonAllowlistedActorIgnored(t *testing.T) {
	g := newGate(t)
	g.issues(issueP0)
	g.events(7, labeled("stranger"))
	out, code := g.run()
	if code != 0 || !strings.Contains(out, "non-allowlisted @stranger") {
		t.Fatalf("code=%d\n%s", code, out)
	}
}

func TestGate_UnlabeledLaterNoLongerBlocks(t *testing.T) {
	g := newGate(t)
	g.issues(issueP0)
	g.events(7, `[{"event":"labeled","label":{"name":"P0"},"actor":{"login":"lea"},"created_at":"2026-11-03T01:00:00Z"},
	{"event":"unlabeled","label":{"name":"P0"},"actor":{"login":"lea"},"created_at":"2026-11-04T01:00:00Z"}]`)
	if out, code := g.run(); code != 0 {
		t.Fatalf("code=%d\n%s", code, out)
	}
	// A stranger re-adding the label after a maintainer removed it must not block.
	g.events(7, `[{"event":"labeled","label":{"name":"P0"},"actor":{"login":"lea"},"created_at":"2026-11-03T01:00:00Z"},
	{"event":"unlabeled","label":{"name":"P0"},"actor":{"login":"lea"},"created_at":"2026-11-04T01:00:00Z"},
	{"event":"labeled","label":{"name":"P0"},"actor":{"login":"stranger"},"created_at":"2026-11-05T01:00:00Z"}]`)
	if out, code := g.run(); code != 0 {
		t.Fatalf("code=%d\n%s", code, out)
	}
}

func TestGate_HighPriorityLabelIgnored(t *testing.T) {
	g := newGate(t)
	// Even if a maintainer-looking issue carries high-priority, it is never queried.
	g.set("repos/wailsapp/wails/issues", `[{"number":9,"title":"urgent critical production blocker","html_url":"u","created_at":"2026-11-03T00:00:00Z","labels":[{"name":"high-priority"}]}]`)
	// The stub returns this for every label query, so make events say a DIFFERENT label.
	g.events(9, `[{"event":"labeled","label":{"name":"high-priority"},"actor":{"login":"lea"},"created_at":"2026-11-03T01:00:00Z"}]`)
	if out, code := g.run(); code != 0 {
		t.Fatalf("high-priority must not block: code=%d\n%s", code, out)
	}
}

func TestGate_RegressionOlderThanLastWeeklyIgnored(t *testing.T) {
	g := newGate(t)
	old := `[{"number":5,"title":"old regression","html_url":"u","created_at":"2026-10-01T00:00:00Z"}]`
	ev := `[{"event":"labeled","label":{"name":"regression"},"actor":{"login":"lea"},"created_at":"2026-10-01T01:00:00Z"}]`
	g.issues(old)
	g.events(5, ev)
	if out, code := g.run("BLOCKER_LABELS="); code != 0 {
		t.Fatalf("old regression must be ignored: code=%d\n%s", code, out)
	}
	g.issues(`[{"number":5,"title":"new regression","html_url":"u","created_at":"2026-11-02T00:00:00Z"}]`)
	if out, code := g.run("BLOCKER_LABELS="); code != 1 || !strings.Contains(out, "[regression] #5") {
		t.Fatalf("new regression must block: code=%d\n%s", code, out)
	}
	// Pull requests labelled regression are not counted.
	g.issues(`[{"number":5,"title":"pr","html_url":"u","created_at":"2026-11-02T00:00:00Z","pull_request":{}}]`)
	if out, code := g.run("BLOCKER_LABELS="); code != 0 {
		t.Fatalf("PR must be ignored for regression: code=%d\n%s", code, out)
	}
}

func dep(path, sev string) string {
	return `[{"number":3,"html_url":"https://x/3","dependency":{"manifest_path":"` + path + `"},"security_advisory":{"severity":"` + sev + `","summary":"bad dep"}}]`
}

func TestGate_Dependabot(t *testing.T) {
	g := newGate(t)
	g.set("repos/wailsapp/wails/dependabot/alerts", dep("v3/examples/foo/go.mod", "critical"))
	if out, code := g.run(); code != 0 || !strings.Contains(out, "not shipped code") {
		t.Fatalf("example-only alert must be ignored: code=%d\n%s", code, out)
	}
	g.set("repos/wailsapp/wails/dependabot/alerts", dep("v3/go.mod", "critical"))
	if out, code := g.run(); code != 1 || !strings.Contains(out, "[dependabot/critical]") {
		t.Fatalf("critical go.mod alert must block: code=%d\n%s", code, out)
	}
	g.set("repos/wailsapp/wails/dependabot/alerts", dep("v3/go.mod", "medium"))
	if out, code := g.run(); code != 0 {
		t.Fatalf("medium must not block: code=%d\n%s", code, out)
	}
}

func TestGate_SecretAndAdvisory(t *testing.T) {
	g := newGate(t)
	g.set("repos/wailsapp/wails/secret-scanning/alerts", `[{"number":1,"html_url":"u","secret_type_display_name":"GitHub token"}]`)
	g.set("repos/wailsapp/wails/secret-scanning/alerts/1/locations", `[{"details":{"path":"v3/examples/x/.env"}}]`)
	if _, code := g.run(); code != 0 {
		t.Fatal("secret only in examples must be ignored")
	}
	g.set("repos/wailsapp/wails/secret-scanning/alerts/1/locations", `[{"details":{"path":"v3/pkg/app/config.go"}}]`)
	if out, code := g.run(); code != 1 || !strings.Contains(out, "[secret-scanning]") {
		t.Fatalf("code=%d\n%s", code, out)
	}
	g = newGate(t)
	g.set("repos/wailsapp/wails/security-advisories", `[{"ghsa_id":"GHSA-xxxx","html_url":"u","summary":"RCE"}]`)
	if out, code := g.run(); code != 1 || !strings.Contains(out, "[advisory/triage] GHSA-xxxx") {
		t.Fatalf("code=%d\n%s", code, out)
	}
}

func TestGate_APIFailureFailsClosed(t *testing.T) {
	for _, ep := range []string{"dependabot/alerts", "code-scanning/alerts", "secret-scanning/alerts", "security-advisories", "issues"} {
		g := newGate(t)
		g.set("repos/wailsapp/wails/"+ep, "FAIL")
		if out, code := g.run(); code != 2 || !strings.Contains(out, "failing closed") {
			t.Fatalf("%s: code=%d\n%s", ep, code, out)
		}
	}
	// Unreadable allowlist must also fail closed, never pass.
	g := newGate(t)
	g.set("orgs/wailsapp/teams/maintainers", "FAIL")
	if out, code := g.run(); code != 2 || !strings.Contains(out, "allowlist team") {
		t.Fatalf("code=%d\n%s", code, out)
	}
}

func TestGate_OverrideRecordsSkippedBlockers(t *testing.T) {
	g := newGate(t)
	g.issues(issueP0)
	g.events(7, labeled("lea"))
	g.set("repos/wailsapp/wails/issues/12/comments", "[]")
	out, code := g.run("OVERRIDE_REASON=ship the security fix", "GITHUB_ACTOR=lea", "STATUS_ISSUE=12")
	if code != 0 || !strings.Contains(out, "#7") || !strings.Contains(out, "OVERRIDDEN") || !strings.Contains(out, "ship the security fix") || !strings.Contains(out, "`lea`") {
		t.Fatalf("code=%d\n%s", code, out)
	}
	calls, _ := os.ReadFile(filepath.Join(g.dir, "calls.log"))
	if !strings.Contains(string(calls), "POST repos/wailsapp/wails/issues/12/comments") {
		t.Fatalf("status comment not posted: %q", calls)
	}
	// Override also works when the gate could not be evaluated.
	g.set("repos/wailsapp/wails/dependabot/alerts", "FAIL")
	if _, code := g.run("OVERRIDE_REASON=x"); code != 0 {
		t.Fatal("override must skip evaluation errors too")
	}
}

func TestGate_FailureUpdatesExistingStatusComment(t *testing.T) {
	g := newGate(t)
	g.issues(issueP0)
	g.events(7, labeled("lea"))
	g.set("repos/wailsapp/wails/issues/12/comments", `[{"id":99,"body":"<!-- weekly-release-status -->\nold"}]`)
	_, code := g.run("STATUS_ISSUE=12")
	calls, _ := os.ReadFile(filepath.Join(g.dir, "calls.log"))
	if code != 1 || !strings.Contains(string(calls), "PATCH repos/wailsapp/wails/issues/comments/99") {
		t.Fatalf("code=%d calls=%q", code, calls)
	}
}
