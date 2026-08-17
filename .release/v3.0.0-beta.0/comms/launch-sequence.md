# Launch sequence: Wails v3 beta

Publish in this order. The rule behind it: nothing is announced before it is
installable, and nothing is installable before the docs describing it are live.
Every step needs one named owner, because a step owned by everyone is owned by
nobody at 3am.

Owners are blank on purpose. Fill them in before the go/no-go, not on the day.

## Before the day

| # | Step | Owner |
|---|---|---|
| 1 | Confirm that the maintainer-approved policy is still tag-only and that no desktop binaries are expected. If policy changes, update this ledger before proceeding | |
| 2 | Re-sync `releases/v3-beta` with current `master` and re-run the audit. The ledger branch drifts quickly; `relman audit` blocks on it | |
| 3 | Run the tag-only release task in dry-run mode against the audited `master` candidate; verify the proposed version, tag, changelog and release notes without publishing | |
| 4 | Record the dry-run commands and results in `plan.json`, and prepare the post-publication clean-machine check for the exact tag. The install gate remains open until that check runs | |
| 5 | Warn downstream: packagers, template and plugin authors, translators (see `packagers.md`) | |

## Release day

| When | Step | Owner |
|---|---|---|
| T-60m | Final `relman gono`; blockers green or waived in writing | |
| T-45m | Run the approved tag-only release path. `release.go` bumps `version.txt`, folds `UNRELEASED_CHANGELOG.md` into the changelog, and publishes the tag/GitHub prerelease | |
| T-40m | Verify the GitHub prerelease points at the intended commit and carries the approved notes. Zero attached desktop assets are expected under the current policy | |
| T-25m | `version.txt` moving triggers `publish-npm.yml`, which publishes `@wailsio/runtime` at the matching version | |
| T-20m | Verify `go install github.com/wailsapp/wails/v3/cmd/wails3@<tag>` on a clean machine and confirm `wails3 version` reports the beta | |
| T-15m | Repeat the published-channel install on each supported CLI platform; do not use a local build directory | |
| T-10m | Docs deploy (Cloudflare Pages, "wails-v3-site"); confirm the version switcher and every locale | |
| T-5m | Support rota on duty; hotfix path confirmed | |
| T+0 | Announcement post live | |
| T+5m | Community chat and forum: existing users hear it first | |
| T+10m | Show HN plus the first comment | |
| T+15m | Social threads, newsletter | |
| T+30m | Tell packagers the source tag and Go module version are live; no prebuilt CLI assets are published under the current policy | |
| T+1h | First sweep: install failures, docs 404s, thread questions | |
| T+3h | Second sweep; decide on a same-day patch | |
| T+24h | Summarise, file issues, hand over the rota | |

## Abort criteria

Agreed in advance, while nobody is under pressure:

- A pre-publication install fails on any supported platform: hold publication
  and the announcement, then fix and reverify the candidate
- A post-publication install failure, data-loss defect, or security regression:
  hold the announcement, publish an advisory when appropriate, and recover with
  a new semver patch tag rather than replacing the cached version
- Docs site down or the upgrade guide missing: hold the announcement, ship the fix
- Anything else: continue, log it, decide at the T+3h sweep

## Rollback

Deleting the GitHub release or tag does not retract a version already cached by
the Go module proxy. An npm version may be unpublished only when the current
[registry criteria](https://docs.npmjs.com/policies/unpublish/) permit it;
otherwise deprecate it, and an unpublished name/version still cannot be reused.
In practice, the recovery for anything found after publishing is a fast patch
release, which is why the hotfix path has to be rehearsed rather than assumed.

Write the exact commands here once step 2 has actually been executed, so this
section describes something that has been done rather than something believed.
