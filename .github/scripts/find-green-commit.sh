#!/usr/bin/env bash
# Select the newest commit on origin/<branch> whose required CI checks all passed.
#
#   find-green-commit.sh            print the newest green SHA since the last weekly tag
#   find-green-commit.sh --check S  exit 0 if commit S is green, 1 otherwise (prints why)
#
# master does not run the PR workflows on push, so a merged commit usually has no
# check-runs of its own. For each candidate we therefore look at the commit's own
# check-runs first and fall back to the check-runs of the pull request head that
# was merged into it (the same checks the branch ruleset required before merge).
#
# Env: REPO (default wailsapp/wails), BRANCH (default master),
#      REQUIRED_CHECKS (newline separated check-run names),
#      GH (gh binary, default "gh"), MAX_SCAN (default 50).
set -euo pipefail

REPO="${REPO:-wailsapp/wails}"
BRANCH="${BRANCH:-master}"
GH="${GH:-gh}"
MAX_SCAN="${MAX_SCAN:-50}"
REQUIRED_CHECKS="${REQUIRED_CHECKS:-v3 Build Results
Run Go Tests
Bundled runtime matches source}"

# 0 = every required check present and passing; 1 = failing or missing.
checks_green() {
  local sha="$1" runs name
  runs="$($GH api --paginate "repos/$REPO/commits/$sha/check-runs?per_page=100" --jq '.check_runs[] | [.name, (.conclusion // "pending")] | @tsv')"
  while IFS= read -r name; do
    [ -z "$name" ] && continue
    local seen bad
    seen="$(awk -F'\t' -v n="$name" '$1==n' <<<"$runs" | wc -l)"
    bad="$(awk -F'\t' -v n="$name" '$1==n && $2!="success" && $2!="skipped" && $2!="neutral"' <<<"$runs" | wc -l)"
    if [ "$seen" -eq 0 ] || [ "$bad" -ne 0 ]; then
      echo "  $sha: required check '$name' not green (runs=$seen, not-passing=$bad)" >&2
      return 1
    fi
  done <<<"$REQUIRED_CHECKS"
  return 0
}

commit_green() {
  local sha="$1" pr_head
  if checks_green "$sha" 2>/dev/null; then return 0; fi
  pr_head="$($GH api "repos/$REPO/commits/$sha/pulls" --jq '.[0].head.sha // empty' 2>/dev/null || true)"
  if [ -n "$pr_head" ] && [ "$pr_head" != "$sha" ]; then
    checks_green "$pr_head"
    return $?
  fi
  echo "  $sha: no check-runs found" >&2
  return 1
}

if [ "${1:-}" = "--check" ]; then
  commit_green "${2:?sha required}"
  exit $?
fi

# Lower bound: the last weekly tag (vX.Y.0, X.Y > 3.0 or GA). Nothing older qualifies.
last_weekly="$(git tag --list 'v3.*.0' --sort=-version:refname | grep -E '^v3\.[0-9]+\.0$' | head -1 || true)"
range="origin/$BRANCH"
[ -n "$last_weekly" ] && range="$last_weekly..origin/$BRANCH"
echo "Scanning first-parent commits in $range (max $MAX_SCAN)" >&2

n=0
for sha in $(git rev-list --first-parent "$range"); do
  n=$((n + 1)); [ "$n" -gt "$MAX_SCAN" ] && break
  subject="$(git log -1 --format=%s "$sha")"
  case "$subject" in
    "chore(v3): bump to"*|*"[skip ci]"*|"Publish @wailsio/runtime"*|*"chore(v3): sync version"*)
      echo "  skip release-bookkeeping commit $sha ($subject)" >&2; continue ;;
  esac
  if commit_green "$sha"; then
    echo "$sha"
    exit 0
  fi
done

echo "::error::No commit with green required checks on $BRANCH since ${last_weekly:-the beginning}; refusing to cut a weekly release." >&2
exit 1
