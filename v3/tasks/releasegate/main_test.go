package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
)

const testToken = "test-token"

// route is one fake endpoint. Each entry of pages is served for ?page=N
// (default 1) with a Link rel="next" header while more pages remain.
type route struct {
	status int
	pages  []string
}

// fakeGitHub is an httptest fake of the GitHub REST endpoints the gate uses.
// Routes are keyed by path without the query string, so every label query
// shares the issues endpoint (tests use one label at a time). Unknown paths
// answer 404. POST/PATCH calls are recorded with their bodies.
type fakeGitHub struct {
	t      *testing.T
	dir    string
	srv    *httptest.Server
	mu     sync.Mutex
	routes map[string]route
	calls  []string // "METHOD path" for writes
	bodies []string // JSON "body" of each write
	gets   []string // request URI of every GET
}

func newGate(t *testing.T) *fakeGitHub {
	t.Helper()
	g := &fakeGitHub{t: t, dir: t.TempDir(), routes: map[string]route{}}
	g.srv = httptest.NewServer(http.HandlerFunc(g.serve))
	t.Cleanup(g.srv.Close)
	// Healthy baseline: allowlist readable, nothing open.
	g.set("orgs/wailsapp/teams/maintainers", `{"slug":"maintainers"}`)
	g.set("orgs/wailsapp/teams/maintainers/memberships/lea", `{"state":"active"}`)
	g.set("repos/wailsapp/wails/collaborators/taliesin-ai/permission", `{"permission":"maintain"}`)
	g.set("repos/wailsapp/wails/collaborators/stranger/permission", `{"permission":"read"}`)
	for _, ep := range []string{"issues", "dependabot/alerts", "code-scanning/alerts", "secret-scanning/alerts", "security-advisories"} {
		g.set("repos/wailsapp/wails/"+ep, "[]")
	}
	return g
}

func (g *fakeGitHub) serve(w http.ResponseWriter, r *http.Request) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if r.Header.Get("Authorization") != "Bearer "+testToken {
		http.Error(w, `{"message":"Bad credentials"}`, http.StatusUnauthorized)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/")
	if r.Method != http.MethodGet {
		var body struct {
			Body string `json:"body"`
		}
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		g.calls = append(g.calls, r.Method+" "+path)
		g.bodies = append(g.bodies, body.Body)
		fmt.Fprint(w, `{}`)
		return
	}
	g.gets = append(g.gets, r.URL.RequestURI())
	rt, ok := g.routes[path]
	if !ok {
		http.Error(w, `{"message":"Not Found"}`, http.StatusNotFound)
		return
	}
	if rt.status != http.StatusOK {
		http.Error(w, `{"message":"error"}`, rt.status)
		return
	}
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		page, _ = strconv.Atoi(p)
	}
	if page < 1 || page > len(rt.pages) {
		fmt.Fprint(w, "[]")
		return
	}
	if page < len(rt.pages) {
		q := r.URL.Query()
		q.Set("page", strconv.Itoa(page+1))
		next := g.srv.URL + r.URL.Path + "?" + q.Encode()
		w.Header().Set("Link", fmt.Sprintf(`<%s>; rel="next", <%s>; rel="last"`, next, next))
	}
	fmt.Fprint(w, rt.pages[page-1])
}

func (g *fakeGitHub) set(ep, content string) { g.setPages(ep, content) }

func (g *fakeGitHub) setPages(ep string, pages ...string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.routes[strings.SplitN(ep, "?", 2)[0]] = route{status: http.StatusOK, pages: pages}
}

func (g *fakeGitHub) setStatus(ep string, status int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.routes[strings.SplitN(ep, "?", 2)[0]] = route{status: status}
}

func (g *fakeGitHub) issues(list string) { g.set("repos/wailsapp/wails/issues", list) }
func (g *fakeGitHub) events(n int, ev string) {
	g.set(fmt.Sprintf("repos/wailsapp/wails/issues/%d/events", n), ev)
}

func (g *fakeGitHub) writes() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	return strings.Join(g.calls, "\n")
}

// run evaluates the gate in-process against the fake. env entries are
// KEY=VALUE and override the defaults below; nothing is read from the real
// environment. It returns stdout and stderr combined, and the exit code.
func (g *fakeGitHub) run(env ...string) (string, int) {
	g.t.Helper()
	vars := map[string]string{
		"GH_TOKEN":         testToken,
		"GITHUB_API_URL":   g.srv.URL,
		"LAST_WEEKLY_DATE": "2026-11-01T00:00:00Z",
		"REPORT_FILE":      filepath.Join(g.dir, "report.md"),
	}
	for _, kv := range env {
		k, v, _ := strings.Cut(kv, "=")
		vars[k] = v
	}
	var out bytes.Buffer
	code := run(configFromEnv(func(k string) string { return vars[k] }), &out, &out)
	return out.String(), code
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
	if w := g.writes(); w != "" {
		t.Fatalf("a clean gate must not post: %q", w)
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
	// The fake returns this for every label query, so make events say a DIFFERENT label.
	g.events(9, `[{"event":"labeled","label":{"name":"high-priority"},"actor":{"login":"lea"},"created_at":"2026-11-03T01:00:00Z"}]`)
	if out, code := g.run(); code != 0 {
		t.Fatalf("high-priority must not block: code=%d\n%s", code, out)
	}
	for _, uri := range g.gets {
		if strings.Contains(uri, "high-priority") {
			t.Fatalf("high-priority must never be queried: %s", uri)
		}
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

// An alert without a manifest path cannot be proven to be outside shipped
// code, so it counts. (The shell version's `read` collapsed the empty field
// and shifted the severity out of place, silently dropping such alerts.)
func TestGate_DependabotWithoutPathBlocks(t *testing.T) {
	g := newGate(t)
	g.set("repos/wailsapp/wails/dependabot/alerts", dep("", "high"))
	if out, code := g.run(); code != 1 || !strings.Contains(out, "[dependabot/high] #3 bad dep () - https://x/3") {
		t.Fatalf("code=%d\n%s", code, out)
	}
}

func TestGate_CodeScanning(t *testing.T) {
	alert := func(path, sev string) string {
		return `[{"number":4,"html_url":"https://x/4","most_recent_instance":{"location":{"path":"` + path + `"}},"rule":{"id":"go/sql-injection","description":"SQL injection","security_severity_level":"` + sev + `"}}]`
	}
	g := newGate(t)
	g.set("repos/wailsapp/wails/code-scanning/alerts", alert(".github/workflows/x.yml", "high"))
	if out, code := g.run(); code != 0 || !strings.Contains(out, "ignoring code-scanning #4") {
		t.Fatalf("CI-only alert must be ignored: code=%d\n%s", code, out)
	}
	g.set("repos/wailsapp/wails/code-scanning/alerts", alert("v3/pkg/application/app.go", "high"))
	if out, code := g.run(); code != 1 || !strings.Contains(out, "[code-scanning/high] #4 SQL injection (v3/pkg/application/app.go) - https://x/4") {
		t.Fatalf("code=%d\n%s", code, out)
	}
	g.set("repos/wailsapp/wails/code-scanning/alerts", alert("v3/pkg/application/app.go", "low"))
	if out, code := g.run(); code != 0 {
		t.Fatalf("low must not block: code=%d\n%s", code, out)
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
	// No file location (found in issue/PR/commit text): always counts.
	g.set("repos/wailsapp/wails/secret-scanning/alerts/1/locations", `[{"type":"issue_body","details":{}}]`)
	if out, code := g.run(); code != 1 || !strings.Contains(out, "[secret-scanning] #1 GitHub token - u") {
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
		for _, status := range []int{http.StatusForbidden, http.StatusInternalServerError, http.StatusNotFound} {
			g := newGate(t)
			g.setStatus("repos/wailsapp/wails/"+ep, status)
			if out, code := g.run(); code != 2 || !strings.Contains(out, "failing closed") {
				t.Fatalf("%s %d: code=%d\n%s", ep, status, code, out)
			}
		}
	}
	// Unreadable allowlist must also fail closed, never pass.
	for _, status := range []int{http.StatusForbidden, http.StatusInternalServerError} {
		g := newGate(t)
		g.setStatus("orgs/wailsapp/teams/maintainers", status)
		if out, code := g.run(); code != 2 || !strings.Contains(out, "allowlist team") {
			t.Fatalf("%d: code=%d\n%s", status, code, out)
		}
	}
}

func TestGate_PerItemFailuresFailClosed(t *testing.T) {
	// Label events unreadable: the issue may be a blocker, so fail closed.
	g := newGate(t)
	g.issues(issueP0)
	g.setStatus("repos/wailsapp/wails/issues/7/events", http.StatusForbidden)
	if out, code := g.run(); code != 2 || !strings.Contains(out, "cannot read label events of #7") {
		t.Fatalf("code=%d\n%s", code, out)
	}
	// Secret locations unreadable.
	g = newGate(t)
	g.set("repos/wailsapp/wails/secret-scanning/alerts", `[{"number":1,"html_url":"u","secret_type":"token"}]`)
	g.setStatus("repos/wailsapp/wails/secret-scanning/alerts/1/locations", http.StatusInternalServerError)
	if out, code := g.run(); code != 2 || !strings.Contains(out, "cannot read locations of secret-scanning alert #1") {
		t.Fatalf("code=%d\n%s", code, out)
	}
	// A response that is not the expected JSON array is unreadable too.
	g = newGate(t)
	g.set("repos/wailsapp/wails/dependabot/alerts", `{"message":"Must have admin rights"}`)
	if out, code := g.run(); code != 2 || !strings.Contains(out, "cannot read Dependabot alerts") {
		t.Fatalf("code=%d\n%s", code, out)
	}
}

func TestGate_TokenSources(t *testing.T) {
	g := newGate(t)
	if out, code := g.run("GH_TOKEN=", "GITHUB_TOKEN="+testToken); code != 0 {
		t.Fatalf("GITHUB_TOKEN fallback: code=%d\n%s", code, out)
	}
	// No token: every authenticated read fails, so the gate fails closed.
	if out, code := g.run("GH_TOKEN="); code != 2 || !strings.Contains(out, "allowlist team") {
		t.Fatalf("no token: code=%d\n%s", code, out)
	}
}

func TestGate_Pagination(t *testing.T) {
	g := newGate(t)
	// The blocker is on the second page of issues.
	g.setPages("repos/wailsapp/wails/issues",
		`[{"number":1,"title":"other","html_url":"u1","created_at":"2026-11-03T00:00:00Z"}]`,
		issueP0)
	g.events(1, `[]`)
	// The deciding `labeled` event is on the second page of events.
	g.setPages("repos/wailsapp/wails/issues/7/events",
		`[{"event":"labeled","label":{"name":"P0"},"actor":{"login":"stranger"},"created_at":"2026-11-02T01:00:00Z"}]`,
		`[{"event":"unlabeled","label":{"name":"P0"},"actor":{"login":"lea"},"created_at":"2026-11-02T05:00:00Z"},
		  {"event":"labeled","label":{"name":"P0"},"actor":{"login":"lea"},"created_at":"2026-11-03T01:00:00Z"}]`)
	// The marker comment is on the second page of comments.
	g.setPages("repos/wailsapp/wails/issues/12/comments",
		`[{"id":50,"body":"unrelated"}]`,
		`[{"id":99,"body":"<!-- weekly-release-status -->\nold"},{"id":100,"body":"later"}]`)
	out, code := g.run("STATUS_ISSUE=12")
	if code != 1 || !strings.Contains(out, "[P0] #7 Crash on start") || !strings.Contains(out, "@lea") {
		t.Fatalf("code=%d\n%s", code, out)
	}
	if w := g.writes(); w != "PATCH repos/wailsapp/wails/issues/comments/99" {
		t.Fatalf("writes=%q", w)
	}

	// An `unlabeled` on the second page of events means the label no longer counts.
	g = newGate(t)
	g.issues(issueP0)
	g.setPages("repos/wailsapp/wails/issues/7/events",
		labeled("lea"),
		`[{"event":"unlabeled","label":{"name":"P0"},"actor":{"login":"lea"},"created_at":"2026-11-04T01:00:00Z"}]`)
	if out, code := g.run(); code != 0 {
		t.Fatalf("code=%d\n%s", code, out)
	}

	// A failure on a later page fails the whole list closed.
	g = newGate(t)
	g.setPages("repos/wailsapp/wails/dependabot/alerts", "[]", "not json")
	if out, code := g.run(); code != 2 {
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
	if !strings.Contains(g.writes(), "POST repos/wailsapp/wails/issues/12/comments") {
		t.Fatalf("status comment not posted: %q", g.writes())
	}
	// Override also works when the gate could not be evaluated.
	g.setStatus("repos/wailsapp/wails/dependabot/alerts", http.StatusInternalServerError)
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
	if code != 1 || !strings.Contains(g.writes(), "PATCH repos/wailsapp/wails/issues/comments/99") {
		t.Fatalf("code=%d calls=%q", code, g.writes())
	}
}

func TestGate_ReportOutputs(t *testing.T) {
	g := newGate(t)
	g.issues(issueP0)
	g.events(7, labeled("lea"))
	g.set("repos/wailsapp/wails/issues/12/comments", "[]")
	summary := filepath.Join(g.dir, "summary.md")
	if err := os.WriteFile(summary, []byte("earlier step\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := g.run("STATUS_ISSUE=12", "GITHUB_STEP_SUMMARY="+summary)
	want := "<!-- weekly-release-status -->\n" +
		"## Weekly release gate\n" +
		"\n" +
		"**Blockers (1):**\n" +
		"- [P0] #7 Crash on start - https://github.com/wailsapp/wails/issues/7 (label set by @lea)\n"
	if code != 1 || out != want {
		t.Fatalf("code=%d\ngot:\n%s\nwant:\n%s", code, out, want)
	}
	if b, _ := os.ReadFile(filepath.Join(g.dir, "report.md")); string(b) != want {
		t.Fatalf("report file:\n%s", b)
	}
	if b, _ := os.ReadFile(summary); string(b) != "earlier step\n"+want {
		t.Fatalf("step summary:\n%s", b)
	}
	if len(g.bodies) != 1 || g.bodies[0] != want {
		t.Fatalf("status comment body: %q", g.bodies)
	}
}

func TestNextLink(t *testing.T) {
	for header, want := range map[string]string{
		"": "",
		`<https://api.github.com/x?page=2>; rel="next", <https://api.github.com/x?page=5>; rel="last"`: "https://api.github.com/x?page=2",
		`<https://api.github.com/x?page=1>; rel="prev", <https://api.github.com/x?page=3>; rel="next"`: "https://api.github.com/x?page=3",
		`<https://api.github.com/x?page=1>; rel="first"`:                                               "",
	} {
		if got := nextLink(header); got != want {
			t.Errorf("nextLink(%q) = %q, want %q", header, got, want)
		}
	}
}

func TestNormaliseSince(t *testing.T) {
	for in, want := range map[string]string{
		"":                          "",
		"2026-11-01T00:00:00Z":      "2026-11-01T00:00:00Z",
		"2026-11-01T11:00:00+11:00": "2026-11-01T00:00:00Z",
		"not a date":                "not a date",
	} {
		if got := normaliseSince(in); got != want {
			t.Errorf("normaliseSince(%q) = %q, want %q", in, got, want)
		}
	}
}

// An unknown answer about who set a blocker label fails closed: a 5xx (or a
// rate limit) while checking the applier must not let the label be ignored
// and the release pass. A 404 stays a definite "not allowed".
func TestGate_ActorCheckErrorFailsClosed(t *testing.T) {
	g := newGate(t)
	g.issues(issueP0)
	g.events(7, labeled("flaky"))
	g.setStatus("orgs/wailsapp/teams/maintainers/memberships/flaky", http.StatusInternalServerError)
	out, code := g.run()
	if code != 2 || !strings.Contains(out, "@flaky") {
		t.Fatalf("code=%d (want 2, could not evaluate)\n%s", code, out)
	}
}
