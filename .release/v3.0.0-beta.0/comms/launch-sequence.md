# Launch sequence: Wails v3 beta

Publish in this order. The rule behind it: nothing is announced before it is
installable, and published-channel installation checks do not begin before the
docs describing the release are live. Every step needs one named owner, because
a step owned by everyone is owned by nobody at 3am.

Owners are blank on purpose. Fill them in before the go/no-go, not on the day.

## Before the day

| # | Step | Owner |
|---|---|---|
| 1 | Confirm that the maintainer-approved policy is still tag-only and that no desktop binaries are expected. If policy changes, update this ledger before proceeding | |
| 2 | Re-sync `releases/v3-beta` with current `master` and re-run the audit. The ledger branch drifts quickly; `relman audit` blocks on it | |
| 3 | Run the tag-only release task in dry-run mode against the audited `master` candidate; verify the proposed version, tag, changelog and release notes without publishing | |
| 4 | Record the dry-run commands and results in `plan.json`, and prepare the post-publication clean-machine check for the exact tag. The install gate remains open until that check runs | |
| 5 | Build the documentation from the audited candidate and verify a staged preview, including the version switcher, changelog, migration content and translated locales. Record the preview evidence; do not publish if it fails | |
| 6 | Warn downstream: packagers, template and plugin authors, translators (see `packagers.md`) | |

## Release day

| When | Step | Owner |
|---|---|---|
| T-60m | Final `relman gono`; blockers green or waived in writing | |
| T-45m | Run the approved tag-only release path. `release.go` bumps `version.txt`, folds `UNRELEASED_CHANGELOG.md` into the changelog, and publishes the tag/GitHub prerelease | |
| T-40m | Verify the GitHub prerelease points at the intended commit and carries the approved notes. Zero attached desktop assets are expected under the current policy | |
| T-30m | Docs deploy (Cloudflare Pages, "wails-v3-site"); confirm the version switcher and every locale | |
| T-25m | `version.txt` moving triggers `publish-npm.yml`, which publishes `@wailsio/runtime` at the matching version | |
| T-20m | Verify `go install github.com/wailsapp/wails/v3/cmd/wails3@<tag>` on a clean machine and confirm `wails3 version` reports the beta | |
| T-15m | Repeat the published-channel install on each supported CLI platform; do not use a local build directory | |
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

- A dry run or other candidate verification fails before publication: hold
  publication and the announcement, then fix and reverify the candidate
- A published-tag install failure, data-loss defect, or security regression
  found before T+0: hold the announcement and publish an advisory when appropriate
- The same problem found after T+0: correct the public communications immediately
- For either published-version path, identify the affected Go tag and mitigation
  in the advisory, deprecate the matching npm version when unpublishing is not
  permitted, and recover with the next unused, monotonically higher prerelease
  tag in the approved release line (for example, beta.10 after beta.9)
- Docs site down or the upgrade guide missing: before T+0 hold the announcement;
  after T+0 correct the public communications, then ship the fix
- Anything else: continue, log it, decide at the T+3h sweep

## Rollback

Deleting the GitHub release or tag does not retract a version already cached by
the Go module proxy. An npm version may be unpublished only when the current
[registry criteria](https://docs.npmjs.com/policies/unpublish/) permit it;
otherwise deprecate it, and an unpublished name/version still cannot be reused.
In practice, recovery after publishing uses the next unused, monotonically
higher tag in the approved release line, which is why the hotfix path has to be
rehearsed rather than assumed.

Write the exact commands here once step 3 has actually been executed, so this
section describes something that has been done rather than something believed.
