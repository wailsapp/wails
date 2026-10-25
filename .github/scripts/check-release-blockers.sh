#!/usr/bin/env bash
# Weekly release gate. Lists everything that must block a weekly release.
#
# Exit codes: 0 = clean (or overridden), 1 = blockers found, 2 = could not
# evaluate (FAILS CLOSED: an unreadable alert API is never treated as "clean").
#
# What blocks
#   1. open issues/PRs labelled P0 or release-blocker
#   2. open issues labelled regression created since the last weekly
#   3. open Dependabot / code-scanning alerts of critical|high severity and any
#      open secret-scanning alert - only if the alert touches SHIPPED code
#   4. repository security advisories still in "triage"
#
# TRUST RULE: a label counts only if its most recent `labeled` event (with no
# later `unlabeled`) was made by an allowlisted actor, so random users cannot
# block or unblock a release. The `high-priority` label is deliberately NOT
# used: .github/issue-labeler.yml applies it automatically from issue text.
#
# Token scopes (use WAILS_REPO_TOKEN; the default GITHUB_TOKEN cannot read these):
#   fine-grained PAT / GitHub App on wailsapp/wails with read access to:
#     Issues, Pull requests, Metadata, Administration (collaborator permission),
#     Dependabot alerts, Code scanning alerts, Secret scanning alerts,
#     Repository security advisories (Security advisories), and
#     organisation Members (read) for the team-membership check;
#     plus Issues: write ONLY if STATUS_ISSUE is set (to post the status comment).
#   classic PAT equivalent: repo, read:org, security_events.
#
# Environment
#   REPO                  default wailsapp/wails
#   GH                    gh binary (default gh)
#   ALLOW_TEAM            org/team slug allowlist source        (default wailsapp/maintainers)
#   ALLOW_COLLAB_PERMS    collaborator permissions that count   (default "maintain admin")
#   SHIPPED_ALERT_IGNORE  OWNER REVIEW: regex of paths that are NOT shipped code;
#                         alerts located only there are ignored
#   LAST_WEEKLY_DATE      ISO date of the last vX.Y.0 tag (default: derived from git)
#   STATUS_ISSUE          number of the single "Weekly release status" issue (optional)
#   OVERRIDE_REASON       non-empty => report blockers but exit 0
#   REPORT_FILE           markdown report path (default release-gate-report.md)
set -uo pipefail

REPO="${REPO:-wailsapp/wails}"
GH="${GH:-gh}"
ALLOW_TEAM="${ALLOW_TEAM:-wailsapp/maintainers}"
ALLOW_COLLAB_PERMS="${ALLOW_COLLAB_PERMS:-maintain admin}"
# OWNER REVIEW: everything that is not shipped code. Anything else counts.
SHIPPED_ALERT_IGNORE="${SHIPPED_ALERT_IGNORE:-^(v3/examples/|v2/examples/|website/|docs/|v2/website/|\.github/)|(^|/)examples?/}"
OVERRIDE_REASON="${OVERRIDE_REASON:-}"
STATUS_ISSUE="${STATUS_ISSUE:-}"
REPORT_FILE="${REPORT_FILE:-release-gate-report.md}"
MARKER='<!-- weekly-release-status -->'
BLOCKER_LABELS="${BLOCKER_LABELS:-P0 release-blocker}"
REGRESSION_LABEL="${REGRESSION_LABEL:-regression}"

blockers=()
errors=()
declare -A ACTOR_OK

fail_closed() { errors+=("$1"); }

# gh_json PATH... -> stdout JSON; returns non-zero on API failure.
api() { "$GH" api --paginate "$@" 2>/dev/null | jq -s 'add // []'; }

# actor_allowed LOGIN: team member, or collaborator with an allowed permission.
actor_allowed() {
  local login="$1" state perm p
  [ -n "${ACTOR_OK[$login]:-}" ] && { [ "${ACTOR_OK[$login]}" = 1 ]; return; }
  ACTOR_OK[$login]=0
  state="$("$GH" api "orgs/${ALLOW_TEAM%%/*}/teams/${ALLOW_TEAM##*/}/memberships/$login" --jq .state 2>/dev/null || true)"
  if [ "$state" = "active" ]; then ACTOR_OK[$login]=1; return 0; fi
  perm="$("$GH" api "repos/$REPO/collaborators/$login/permission" --jq .permission 2>/dev/null || true)"
  for p in $ALLOW_COLLAB_PERMS; do
    if [ "$perm" = "$p" ]; then ACTOR_OK[$login]=1; return 0; fi
  done
  return 1
}

# label_applier ISSUE LABEL -> login of the actor whose `labeled` event is the
# current state of that label, or empty if the label was later removed.
label_applier() {
  local n="$1" label="$2" events
  events="$(api "repos/$REPO/issues/$n/events?per_page=100")" || return 2
  jq -r --arg l "$label" '
    [ .[] | select((.event=="labeled" or .event=="unlabeled") and .label.name==$l) ]
    | sort_by(.created_at) | last
    | if .event=="labeled" then .actor.login else empty end' <<<"$events"
}

# check_label LABEL SINCE(optional) ISSUES_ONLY(0/1)
check_label() {
  local label="$1" since="${2:-}" issues_only="${3:-0}" list line n title url created applier
  list="$(api "repos/$REPO/issues?state=open&labels=$label&per_page=100")" || { fail_closed "cannot list issues labelled '$label'"; return; }
  while IFS=$'\t' read -r n title url created is_pr; do
    [ -z "$n" ] && continue
    [ "$issues_only" = 1 ] && [ "$is_pr" = true ] && continue
    if [ -n "$since" ] && [[ "$created" < "$since" ]]; then continue; fi
    applier="$(label_applier "$n" "$label")" || { fail_closed "cannot read label events of #$n"; continue; }
    [ -z "$applier" ] && continue
    if actor_allowed "$applier"; then
      blockers+=("[$label] #$n $title - $url (label set by @$applier)")
    else
      echo "ignoring '$label' on #$n: applied by non-allowlisted @$applier" >&2
    fi
  done < <(jq -r '.[] | [.number, .title, .html_url, .created_at, (has("pull_request"))] | @tsv' <<<"$list")
}

ignored_path() { [ -n "$1" ] && [[ "$1" =~ $SHIPPED_ALERT_IGNORE ]]; }

check_dependabot() {
  local list n url path summary sev
  list="$(api "repos/$REPO/dependabot/alerts?state=open&per_page=100")" || { fail_closed "cannot read Dependabot alerts (token scope?)"; return; }
  while IFS=$'\t' read -r n url path sev summary; do
    [ -z "$n" ] && continue
    [[ "$sev" == critical || "$sev" == high ]] || continue
    ignored_path "$path" && { echo "ignoring Dependabot #$n ($path): not shipped code" >&2; continue; }
    blockers+=("[dependabot/$sev] #$n $summary ($path) - $url")
  done < <(jq -r '.[] | [.number, .html_url, (.dependency.manifest_path // ""), (.security_advisory.severity // ""), (.security_advisory.summary // "")] | @tsv' <<<"$list")
}

check_code_scanning() {
  local list n url path sev rule
  list="$(api "repos/$REPO/code-scanning/alerts?state=open&per_page=100")" || { fail_closed "cannot read code-scanning alerts (token scope?)"; return; }
  while IFS=$'\t' read -r n url path sev rule; do
    [ -z "$n" ] && continue
    [[ "$sev" == critical || "$sev" == high ]] || continue
    ignored_path "$path" && { echo "ignoring code-scanning #$n ($path): not shipped code" >&2; continue; }
    blockers+=("[code-scanning/$sev] #$n $rule ($path) - $url")
  done < <(jq -r '.[] | [.number, .html_url, (.most_recent_instance.location.path // ""), (.rule.security_severity_level // ""), (.rule.description // .rule.id // "")] | @tsv' <<<"$list")
}

check_secret_scanning() {
  local list n url type locs shipped
  list="$(api "repos/$REPO/secret-scanning/alerts?state=open&per_page=100")" || { fail_closed "cannot read secret-scanning alerts (token scope?)"; return; }
  while IFS=$'\t' read -r n url type; do
    [ -z "$n" ] && continue
    locs="$(api "repos/$REPO/secret-scanning/alerts/$n/locations")" || { fail_closed "cannot read locations of secret-scanning alert #$n"; continue; }
    shipped=0
    local paths; paths="$(jq -r '.[] | (.details.path // "")' <<<"$locs")"
    if [ -z "$paths" ]; then shipped=1; fi   # no file location (issue/PR/commit text): always counts
    while IFS= read -r path; do
      [ -z "$path" ] && continue
      ignored_path "$path" || shipped=1
    done <<<"$paths"
    [ "$shipped" = 1 ] && blockers+=("[secret-scanning] #$n $type - $url")
  done < <(jq -r '.[] | [.number, .html_url, (.secret_type_display_name // .secret_type // "secret")] | @tsv' <<<"$list")
}

check_advisories() {
  local list ghsa url summary
  list="$(api "repos/$REPO/security-advisories?state=triage&per_page=100")" || { fail_closed "cannot read repository security advisories (token scope?)"; return; }
  while IFS=$'\t' read -r ghsa url summary; do
    [ -z "$ghsa" ] && continue
    blockers+=("[advisory/triage] $ghsa $summary - $url")
  done < <(jq -r '.[] | [.ghsa_id, .html_url, (.summary // "")] | @tsv' <<<"$list")
}

preflight() {
  # Prove the token can see the allowlist; otherwise every label would look
  # "applied by a stranger" and the gate would silently pass.
  "$GH" api "orgs/${ALLOW_TEAM%%/*}/teams/${ALLOW_TEAM##*/}" --jq .slug >/dev/null 2>&1 \
    || fail_closed "cannot read allowlist team $ALLOW_TEAM (token needs org Members read)"
}

since="${LAST_WEEKLY_DATE:-}"
if [ -z "$since" ]; then
  last="$(git tag --list 'v3.*.0' --sort=-version:refname 2>/dev/null | grep -E '^v3\.[0-9]+\.0$' | head -1 || true)"
  [ -n "$last" ] && since="$(git log -1 --format=%cI "$last" 2>/dev/null || true)"
fi

[ -n "$since" ] && since="$(date -u -d "$since" +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo "$since")"

preflight
for l in $BLOCKER_LABELS; do check_label "$l"; done
check_label "$REGRESSION_LABEL" "$since" 1
check_dependabot
check_code_scanning
check_secret_scanning
check_advisories

{
  echo "$MARKER"
  echo "## Weekly release gate"
  echo
  if [ "${#errors[@]}" -gt 0 ]; then
    echo "**Could not evaluate the gate (failing closed):**"
    for e in "${errors[@]}"; do echo "- $e"; done
    echo
  fi
  if [ "${#blockers[@]}" -gt 0 ]; then
    echo "**Blockers (${#blockers[@]}):**"
    for b in "${blockers[@]}"; do echo "- $b"; done
  elif [ "${#errors[@]}" -eq 0 ]; then
    echo "No blockers."
  fi
  if [ -n "$OVERRIDE_REASON" ]; then
    echo
    echo "**OVERRIDDEN** by \`${GITHUB_ACTOR:-unknown}\`: ${OVERRIDE_REASON}"
    echo "The blockers and errors above were skipped."
  fi
} >"$REPORT_FILE"
cat "$REPORT_FILE"
[ -n "${GITHUB_STEP_SUMMARY:-}" ] && cat "$REPORT_FILE" >>"$GITHUB_STEP_SUMMARY"

post_status() {
  [ -z "$STATUS_ISSUE" ] && return 0
  local cid
  cid="$("$GH" api --paginate "repos/$REPO/issues/$STATUS_ISSUE/comments?per_page=100" 2>/dev/null \
    | jq -r --arg m "$MARKER" '[.[] | select(.body | contains($m))] | last | .id // empty')"
  if [ -n "$cid" ]; then
    "$GH" api -X PATCH "repos/$REPO/issues/comments/$cid" -F "body=@$REPORT_FILE" >/dev/null
  else
    "$GH" api -X POST "repos/$REPO/issues/$STATUS_ISSUE/comments" -F "body=@$REPORT_FILE" >/dev/null
  fi
}

# Only post when there is something to say (a failure or an override).
if [ "${#blockers[@]}" -gt 0 ] || [ "${#errors[@]}" -gt 0 ]; then post_status || echo "warning: could not post status comment" >&2; fi

if [ -n "$OVERRIDE_REASON" ]; then exit 0; fi
[ "${#errors[@]}" -gt 0 ] && exit 2
[ "${#blockers[@]}" -gt 0 ] && exit 1
exit 0
