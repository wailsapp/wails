// Command releasegate is the weekly release gate. It lists everything that
// must block a weekly release.
//
// Run it from the v3 module: `go run ./tasks/releasegate`.
//
// Exit codes: 0 = clean (or overridden), 1 = blockers found, 2 = could not
// evaluate. The gate FAILS CLOSED: an API it cannot read is never treated as
// "clean".
//
// What blocks
//  1. open issues/PRs labelled P0 or release-blocker
//  2. open issues labelled regression created since the last weekly
//  3. open Dependabot / code-scanning alerts of critical|high severity and any
//     open secret-scanning alert - only if the alert touches SHIPPED code
//  4. repository security advisories still in "triage"
//
// TRUST RULE: a label counts only if its most recent `labeled` event (with no
// later `unlabeled`) was made by an allowlisted actor, so random users cannot
// block or unblock a release. The `high-priority` label is deliberately NOT
// used: .github/issue-labeler.yml applies it automatically from issue text.
//
// Token scopes (use WAILS_REPO_TOKEN; the default GITHUB_TOKEN cannot read these):
//
//	fine-grained PAT / GitHub App on wailsapp/wails with read access to:
//	  Issues, Pull requests, Metadata, Administration (collaborator permission),
//	  Dependabot alerts, Code scanning alerts, Secret scanning alerts,
//	  Repository security advisories (Security advisories), and
//	  organisation Members (read) for the team-membership check;
//	  plus Issues: write ONLY if STATUS_ISSUE is set (to post the status comment).
//	classic PAT equivalent: repo, read:org, security_events.
//
// Environment (an empty value means the default)
//
//	GH_TOKEN              API token (falls back to GITHUB_TOKEN)
//	GITHUB_API_URL        API base URL (default https://api.github.com)
//	REPO                  default wailsapp/wails
//	ALLOW_TEAM            org/team slug allowlist source        (default wailsapp/maintainers)
//	ALLOW_COLLAB_PERMS    collaborator permissions that count   (default "maintain admin")
//	SHIPPED_ALERT_IGNORE  OWNER REVIEW: regex of paths that are NOT shipped code;
//	                      alerts located only there are ignored
//	LAST_WEEKLY_DATE      ISO date of the last v3.Y.0 tag (default: derived from git)
//	STATUS_ISSUE          number of the single "Weekly release status" issue (optional)
//	OVERRIDE_REASON       non-empty => report blockers but exit 0
//	REPORT_FILE           markdown report path (default release-gate-report.md)
//	BLOCKER_LABELS        space-separated labels that block     (default "P0 release-blocker")
//	REGRESSION_LABEL      regression label                      (default regression)
//	GITHUB_ACTOR          recorded with an override             (default unknown)
//	GITHUB_STEP_SUMMARY   if set, the report is appended to it
package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	marker = "<!-- weekly-release-status -->"
	// OWNER REVIEW: everything that is not shipped code. Anything else counts.
	defaultShippedAlertIgnore = `^(v3/examples/|v2/examples/|website/|docs/|v2/website/|\.github/)|(^|/)examples?/`
)

type config struct {
	repo            string
	allowTeam       string
	allowPerms      []string
	shippedIgnore   string
	since           string
	statusIssue     string
	overrideReason  string
	reportFile      string
	blockerLabels   []string
	regressionLabel string
	actor           string
	stepSummary     string
	apiURL          string
	token           string
}

// configFromEnv reads the configuration. Like the shell's ${X:-default},
// an empty variable means the default.
func configFromEnv(getenv func(string) string) config {
	envOr := func(key, def string) string {
		if v := getenv(key); v != "" {
			return v
		}
		return def
	}
	token := getenv("GH_TOKEN")
	if token == "" {
		token = getenv("GITHUB_TOKEN")
	}
	since := getenv("LAST_WEEKLY_DATE")
	if since == "" {
		since = lastWeeklyFromGit()
	}
	return config{
		repo:            envOr("REPO", "wailsapp/wails"),
		allowTeam:       envOr("ALLOW_TEAM", "wailsapp/maintainers"),
		allowPerms:      strings.Fields(envOr("ALLOW_COLLAB_PERMS", "maintain admin")),
		shippedIgnore:   envOr("SHIPPED_ALERT_IGNORE", defaultShippedAlertIgnore),
		since:           normaliseSince(since),
		statusIssue:     getenv("STATUS_ISSUE"),
		overrideReason:  getenv("OVERRIDE_REASON"),
		reportFile:      envOr("REPORT_FILE", "release-gate-report.md"),
		blockerLabels:   strings.Fields(envOr("BLOCKER_LABELS", "P0 release-blocker")),
		regressionLabel: envOr("REGRESSION_LABEL", "regression"),
		actor:           envOr("GITHUB_ACTOR", "unknown"),
		stepSummary:     getenv("GITHUB_STEP_SUMMARY"),
		apiURL:          getenv("GITHUB_API_URL"),
		token:           token,
	}
}

var weeklyTag = regexp.MustCompile(`^v3\.[0-9]+\.0$`)

// lastWeeklyFromGit returns the committer date of the newest v3.Y.0 tag, or
// "" if there is none (then regressions are not filtered by date).
func lastWeeklyFromGit() string {
	out, err := exec.Command("git", "tag", "--list", "v3.*.0", "--sort=-version:refname").Output()
	if err != nil {
		return ""
	}
	for _, tag := range strings.Split(string(out), "\n") {
		tag = strings.TrimSpace(tag)
		if !weeklyTag.MatchString(tag) {
			continue
		}
		date, err := exec.Command("git", "log", "-1", "--format=%cI", tag).Output()
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(date))
	}
	return ""
}

// normaliseSince converts the date to UTC `YYYY-MM-DDTHH:MM:SSZ` so it
// compares as a string against GitHub's created_at. Unparseable values are
// used as given.
func normaliseSince(since string) string {
	if since == "" {
		return ""
	}
	if t, err := time.Parse(time.RFC3339, since); err == nil {
		return t.UTC().Format("2006-01-02T15:04:05Z")
	}
	for _, layout := range []string{"2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02"} {
		if t, err := time.ParseInLocation(layout, since, time.Local); err == nil {
			return t.UTC().Format("2006-01-02T15:04:05Z")
		}
	}
	return since
}

type gate struct {
	cfg    config
	gh     *client
	stderr io.Writer
	// ignore is nil when SHIPPED_ALERT_IGNORE does not compile; every path
	// then counts as shipped.
	ignore   *regexp.Regexp
	blockers []string
	errors   []string
	actorOK  map[string]bool
}

func main() {
	os.Exit(run(configFromEnv(os.Getenv), os.Stdout, os.Stderr))
}

// run evaluates the gate, writes the report and returns the exit code.
func run(cfg config, stdout, stderr io.Writer) int {
	g := &gate{
		cfg:     cfg,
		gh:      newClient(cfg.apiURL, cfg.token),
		stderr:  stderr,
		actorOK: map[string]bool{},
	}
	re, err := regexp.Compile(cfg.shippedIgnore)
	if err != nil {
		// Fail towards blocking: nothing is ignored.
		fmt.Fprintf(stderr, "warning: invalid SHIPPED_ALERT_IGNORE (%v); treating every alert as shipped code\n", err)
	} else {
		g.ignore = re
	}

	g.preflight()
	for _, label := range cfg.blockerLabels {
		g.checkLabel(label, "", false)
	}
	g.checkLabel(cfg.regressionLabel, cfg.since, true)
	g.checkDependabot()
	g.checkCodeScanning()
	g.checkSecretScanning()
	g.checkAdvisories()

	report := g.report()
	if err := os.WriteFile(cfg.reportFile, []byte(report), 0o644); err != nil {
		fmt.Fprintf(stderr, "warning: could not write %s: %v\n", cfg.reportFile, err)
	}
	fmt.Fprint(stdout, report)
	if cfg.stepSummary != "" {
		appendFile(cfg.stepSummary, report, stderr)
	}

	// Only post when there is something to say (a failure or an override).
	if len(g.blockers) > 0 || len(g.errors) > 0 {
		if err := g.postStatus(report); err != nil {
			fmt.Fprintln(stderr, "warning: could not post status comment")
		}
	}

	switch {
	case cfg.overrideReason != "":
		return 0
	case len(g.errors) > 0:
		return 2
	case len(g.blockers) > 0:
		return 1
	}
	return 0
}

func appendFile(path, content string, stderr io.Writer) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(stderr, "warning: could not write %s: %v\n", path, err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		fmt.Fprintf(stderr, "warning: could not write %s: %v\n", path, err)
	}
}

func (g *gate) failClosed(msg string) { g.errors = append(g.errors, msg) }

// teamParts splits ALLOW_TEAM like the shell's ${T%%/*} and ${T##*/}: the org
// is everything before the first "/", the team slug everything after the last.
func (g *gate) teamParts() (org, slug string) {
	t := g.cfg.allowTeam
	org, slug = t, t
	if i := strings.Index(t, "/"); i >= 0 {
		org = t[:i]
	}
	if i := strings.LastIndex(t, "/"); i >= 0 {
		slug = t[i+1:]
	}
	return org, slug
}

// preflight proves the token can see the allowlist; otherwise every label
// would look "applied by a stranger" and the gate would silently pass.
func (g *gate) preflight() {
	org, slug := g.teamParts()
	var team struct{}
	if err := g.gh.get(fmt.Sprintf("orgs/%s/teams/%s", org, slug), &team); err != nil {
		g.failClosed(fmt.Sprintf("cannot read allowlist team %s (token needs org Members read)", g.cfg.allowTeam))
	}
}

// actorAllowed reports whether login is an active member of ALLOW_TEAM or a
// collaborator whose permission is in ALLOW_COLLAB_PERMS. A failed lookup
// means "not allowed" (preflight already proved the allowlist is readable).
// Results are cached per login.
func (g *gate) actorAllowed(login string) (bool, error) {
	if ok, seen := g.actorOK[login]; seen {
		return ok, nil
	}
	org, slug := g.teamParts()
	var membership struct {
		State string `json:"state"`
	}
	err := g.gh.get(fmt.Sprintf("orgs/%s/teams/%s/memberships/%s", org, slug, login), &membership)
	switch {
	case err == nil && membership.State == "active":
		g.actorOK[login] = true
		return true, nil
	case err != nil && !isNotFound(err):
		// Only a 404 means "not a member". Any other failure (5xx, rate
		// limit, network) is unknown, and an unknown answer must not let a
		// maintainer's blocker label be ignored: the gate fails closed.
		return false, err
	}
	var collab struct {
		Permission string `json:"permission"`
	}
	if err := g.gh.get(fmt.Sprintf("repos/%s/collaborators/%s/permission", g.cfg.repo, login), &collab); err != nil {
		if !isNotFound(err) {
			return false, err
		}
		g.actorOK[login] = false
		return false, nil
	}
	for _, p := range g.cfg.allowPerms {
		if collab.Permission == p {
			g.actorOK[login] = true
			return true, nil
		}
	}
	g.actorOK[login] = false
	return false, nil
}

// isNotFound reports a GitHub 404 (a definite "no", not an unknown).
func isNotFound(err error) bool {
	var ae *apiError
	return errors.As(err, &ae) && ae.status == 404
}

type issueEvent struct {
	Event string `json:"event"`
	Label struct {
		Name string `json:"name"`
	} `json:"label"`
	Actor struct {
		Login string `json:"login"`
	} `json:"actor"`
	CreatedAt string `json:"created_at"`
}

// labelApplier implements the TRUST RULE's first half: it returns the login
// of the actor whose `labeled` event is the current state of the label on
// issue n, or "" if the latest event for that label is an `unlabeled` (so a
// label re-added by someone else after a maintainer removed it is attributed
// to that someone else, not to the maintainer).
func (g *gate) labelApplier(n int, label string) (string, error) {
	events, err := list[issueEvent](g.gh, fmt.Sprintf("repos/%s/issues/%d/events?per_page=100", g.cfg.repo, n))
	if err != nil {
		return "", err
	}
	var relevant []issueEvent
	for _, e := range events {
		if (e.Event == "labeled" || e.Event == "unlabeled") && e.Label.Name == label {
			relevant = append(relevant, e)
		}
	}
	if len(relevant) == 0 {
		return "", nil
	}
	// Stable, so events with equal timestamps keep API order (as jq sort_by).
	sort.SliceStable(relevant, func(i, j int) bool { return relevant[i].CreatedAt < relevant[j].CreatedAt })
	last := relevant[len(relevant)-1]
	if last.Event != "labeled" {
		return "", nil
	}
	return last.Actor.Login, nil
}

type issue struct {
	Number      int    `json:"number"`
	Title       string `json:"title"`
	HTMLURL     string `json:"html_url"`
	CreatedAt   string `json:"created_at"`
	PullRequest *any   `json:"pull_request"`
}

// checkLabel blocks on open issues (and, unless issuesOnly, PRs) carrying
// label, created at or after since (if set), whose label was applied by an
// allowlisted actor.
func (g *gate) checkLabel(label, since string, issuesOnly bool) {
	issues, err := list[issue](g.gh, fmt.Sprintf("repos/%s/issues?state=open&labels=%s&per_page=100", g.cfg.repo, url.QueryEscape(label)))
	if err != nil {
		g.failClosed(fmt.Sprintf("cannot list issues labelled '%s'", label))
		return
	}
	for _, is := range issues {
		if issuesOnly && is.PullRequest != nil {
			continue
		}
		// Both are UTC ISO-8601 timestamps, so string order is time order.
		if since != "" && is.CreatedAt < since {
			continue
		}
		applier, err := g.labelApplier(is.Number, label)
		if err != nil {
			g.failClosed(fmt.Sprintf("cannot read label events of #%d", is.Number))
			continue
		}
		if applier == "" {
			continue
		}
		allowed, aerr := g.actorAllowed(applier)
		if aerr != nil {
			g.failClosed(fmt.Sprintf("cannot check whether @%s may set '%s' on #%d", applier, label, is.Number))
			continue
		}
		if allowed {
			g.blockers = append(g.blockers, fmt.Sprintf("[%s] #%d %s - %s (label set by @%s)", label, is.Number, is.Title, is.HTMLURL, applier))
		} else {
			fmt.Fprintf(g.stderr, "ignoring '%s' on #%d: applied by non-allowlisted @%s\n", label, is.Number, applier)
		}
	}
}

// ignoredPath is the SHIPPED-CODE FILTER: true when path matches
// SHIPPED_ALERT_IGNORE (examples, docs, website, CI config). An empty path
// (location unknown) is never ignored, so it counts as shipped.
func (g *gate) ignoredPath(path string) bool {
	return path != "" && g.ignore != nil && g.ignore.MatchString(path)
}

func blockingSeverity(sev string) bool { return sev == "critical" || sev == "high" }

func (g *gate) checkDependabot() {
	type alert struct {
		Number     int    `json:"number"`
		HTMLURL    string `json:"html_url"`
		Dependency struct {
			ManifestPath string `json:"manifest_path"`
		} `json:"dependency"`
		SecurityAdvisory struct {
			Severity string `json:"severity"`
			Summary  string `json:"summary"`
		} `json:"security_advisory"`
	}
	alerts, err := list[alert](g.gh, fmt.Sprintf("repos/%s/dependabot/alerts?state=open&per_page=100", g.cfg.repo))
	if err != nil {
		g.failClosed("cannot read Dependabot alerts (token scope?)")
		return
	}
	for _, a := range alerts {
		sev, path := a.SecurityAdvisory.Severity, a.Dependency.ManifestPath
		if !blockingSeverity(sev) {
			continue
		}
		if g.ignoredPath(path) {
			fmt.Fprintf(g.stderr, "ignoring Dependabot #%d (%s): not shipped code\n", a.Number, path)
			continue
		}
		g.blockers = append(g.blockers, fmt.Sprintf("[dependabot/%s] #%d %s (%s) - %s", sev, a.Number, a.SecurityAdvisory.Summary, path, a.HTMLURL))
	}
}

func (g *gate) checkCodeScanning() {
	type alert struct {
		Number             int    `json:"number"`
		HTMLURL            string `json:"html_url"`
		MostRecentInstance struct {
			Location struct {
				Path string `json:"path"`
			} `json:"location"`
		} `json:"most_recent_instance"`
		Rule struct {
			SecuritySeverityLevel string  `json:"security_severity_level"`
			Description           *string `json:"description"`
			ID                    *string `json:"id"`
		} `json:"rule"`
	}
	alerts, err := list[alert](g.gh, fmt.Sprintf("repos/%s/code-scanning/alerts?state=open&per_page=100", g.cfg.repo))
	if err != nil {
		g.failClosed("cannot read code-scanning alerts (token scope?)")
		return
	}
	for _, a := range alerts {
		sev, path := a.Rule.SecuritySeverityLevel, a.MostRecentInstance.Location.Path
		if !blockingSeverity(sev) {
			continue
		}
		if g.ignoredPath(path) {
			fmt.Fprintf(g.stderr, "ignoring code-scanning #%d (%s): not shipped code\n", a.Number, path)
			continue
		}
		rule := firstOf(a.Rule.Description, a.Rule.ID, "")
		g.blockers = append(g.blockers, fmt.Sprintf("[code-scanning/%s] #%d %s (%s) - %s", sev, a.Number, rule, path, a.HTMLURL))
	}
}

// checkSecretScanning blocks on every open secret-scanning alert unless all
// of its file locations are outside shipped code. An alert with no file
// location at all (found in issue/PR/commit text) always counts.
func (g *gate) checkSecretScanning() {
	type alert struct {
		Number                int     `json:"number"`
		HTMLURL               string  `json:"html_url"`
		SecretTypeDisplayName *string `json:"secret_type_display_name"`
		SecretType            *string `json:"secret_type"`
	}
	type location struct {
		Details struct {
			Path string `json:"path"`
		} `json:"details"`
	}
	alerts, err := list[alert](g.gh, fmt.Sprintf("repos/%s/secret-scanning/alerts?state=open&per_page=100", g.cfg.repo))
	if err != nil {
		g.failClosed("cannot read secret-scanning alerts (token scope?)")
		return
	}
	for _, a := range alerts {
		locs, err := list[location](g.gh, fmt.Sprintf("repos/%s/secret-scanning/alerts/%d/locations", g.cfg.repo, a.Number))
		if err != nil {
			g.failClosed(fmt.Sprintf("cannot read locations of secret-scanning alert #%d", a.Number))
			continue
		}
		hasPath, shipped := false, false
		for _, l := range locs {
			if l.Details.Path == "" {
				continue
			}
			hasPath = true
			if !g.ignoredPath(l.Details.Path) {
				shipped = true
			}
		}
		if !hasPath || shipped {
			g.blockers = append(g.blockers, fmt.Sprintf("[secret-scanning] #%d %s - %s", a.Number, firstOf(a.SecretTypeDisplayName, a.SecretType, "secret"), a.HTMLURL))
		}
	}
}

func (g *gate) checkAdvisories() {
	type advisory struct {
		GHSAID  string `json:"ghsa_id"`
		HTMLURL string `json:"html_url"`
		Summary string `json:"summary"`
	}
	advisories, err := list[advisory](g.gh, fmt.Sprintf("repos/%s/security-advisories?state=triage&per_page=100", g.cfg.repo))
	if err != nil {
		g.failClosed("cannot read repository security advisories (token scope?)")
		return
	}
	for _, a := range advisories {
		if a.GHSAID == "" {
			continue
		}
		g.blockers = append(g.blockers, fmt.Sprintf("[advisory/triage] %s %s - %s", a.GHSAID, a.Summary, a.HTMLURL))
	}
}

// firstOf mirrors jq's `a // b // def`: the first non-null value.
func firstOf(a, b *string, def string) string {
	if a != nil {
		return *a
	}
	if b != nil {
		return *b
	}
	return def
}

func (g *gate) report() string {
	var b strings.Builder
	line := func(format string, args ...any) { fmt.Fprintf(&b, format+"\n", args...) }
	line("%s", marker)
	line("## Weekly release gate")
	line("")
	if len(g.errors) > 0 {
		line("**Could not evaluate the gate (failing closed):**")
		for _, e := range g.errors {
			line("- %s", e)
		}
		line("")
	}
	if len(g.blockers) > 0 {
		line("**Blockers (%d):**", len(g.blockers))
		for _, bl := range g.blockers {
			line("- %s", bl)
		}
	} else if len(g.errors) == 0 {
		line("No blockers.")
	}
	if g.cfg.overrideReason != "" {
		line("")
		line("**OVERRIDDEN** by `%s`: %s", g.cfg.actor, g.cfg.overrideReason)
		line("The blockers and errors above were skipped.")
	}
	return b.String()
}

// postStatus keeps ONE status comment on STATUS_ISSUE: it updates the last
// comment carrying the marker in place, or creates it. If the comments cannot
// be listed it creates a new one.
func (g *gate) postStatus(report string) error {
	if g.cfg.statusIssue == "" {
		return nil
	}
	type comment struct {
		ID   int64  `json:"id"`
		Body string `json:"body"`
	}
	var id int64
	comments, err := list[comment](g.gh, fmt.Sprintf("repos/%s/issues/%s/comments?per_page=100", g.cfg.repo, g.cfg.statusIssue))
	if err == nil {
		for _, c := range comments {
			if strings.Contains(c.Body, marker) {
				id = c.ID
			}
		}
	}
	body := map[string]string{"body": report}
	if id != 0 {
		return g.gh.send(http.MethodPatch, fmt.Sprintf("repos/%s/issues/comments/%d", g.cfg.repo, id), body)
	}
	return g.gh.send(http.MethodPost, fmt.Sprintf("repos/%s/issues/%s/comments", g.cfg.repo, g.cfg.statusIssue), body)
}
