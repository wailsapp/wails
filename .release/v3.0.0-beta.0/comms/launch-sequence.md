# Launch sequence: Wails v3 beta

Publish in this order. The rule behind it: nothing is announced before it is
installable, and nothing is installable before the docs describing it are live.
Every step needs one named owner, because a step owned by everyone is owned by
nobody at 3am.

Owners are blank on purpose. Fill them in before the go/no-go, not on the day.

## Before the day

| # | Step | Owner |
|---|---|---|
| 1 | Configure the six Apple secrets (`APPLE_SIGNING_CERT`, `APPLE_CERT_PASSWORD`, `APPLE_SIGNING_IDENTITY`, `APPLE_NOTARIZE_USER`, `APPLE_NOTARIZE_PASSWORD`, `APPLE_TEAM_ID`), or decide to ship unsigned macOS binaries with `allow_unsigned_macos` | |
| 2 | Rehearse the release: `release-v3.yml` via `workflow_dispatch` with `draft: true` on a throwaway tag. Use a tag in the `alpha2.` series (`v3.0.0-alpha2.999`); a tag like `v3.0.0-test.1` would outrank the real beta and hijack `@latest` | |
| 3 | Delete the draft release and the throwaway tag once the rehearsal passes | |
| 4 | Re-sync `releases/v3-beta` with master and re-run the audit. The branch drifted 31 commits, then 6 more within hours; `relman audit` blocks on it now | |
| 5 | Warn downstream: packagers, template and plugin authors, translators (see `packagers.md`) | |

## Release day

| When | Step | Owner |
|---|---|---|
| T-60m | Final `relman gono`; blockers green or waived in writing | |
| T-45m | Tag from `releases/v3-beta`. `release.go` bumps `version.txt` and folds `UNRELEASED_CHANGELOG.md` into the changelog | |
| T-40m | The tag push triggers `release-v3.yml`: preflight, three platform builds, sign and notarize, `SHA256SUMS`, provenance attestation | |
| T-25m | `version.txt` moving triggers `publish-npm.yml`, which publishes `@wailsio/runtime` at the matching version | |
| T-20m | Verify: download a binary, check it against `SHA256SUMS`, verify the attestation, run `wails3 version` and confirm it reports the beta | |
| T-15m | Install on a clean machine per platform, from the published channel, not the build directory | |
| T-10m | Docs deploy (Cloudflare Pages, "wails-v3-site"); confirm the version switcher and every locale | |
| T-5m | Support rota on duty; hotfix path confirmed | |
| T+0 | Announcement post live | |
| T+5m | Community chat and forum: existing users hear it first | |
| T+10m | Show HN plus the first comment | |
| T+15m | Social threads, newsletter | |
| T+30m | Tell packagers the artifacts are live | |
| T+1h | First sweep: install failures, docs 404s, thread questions | |
| T+3h | Second sweep; decide on a same-day patch | |
| T+24h | Summarise, file issues, hand over the rota | |

## Abort criteria

Agreed in advance, while nobody is under pressure:

- Install fails on any supported platform: stop, unpublish, fix
- Data loss or a security regression: stop, unpublish, advisory
- Docs site down or the upgrade guide missing: hold the announcement, ship the fix
- Anything else: continue, log it, decide at the T+3h sweep

## Rollback

The GitHub release can be deleted and the tag removed, but **npm cannot be
unpublished after 72 hours** and the Go module proxy caches versions
permanently. In practice, the recovery for anything found after publishing is a
fast patch release, which is why the hotfix path has to be rehearsed rather than
assumed.

Write the exact commands here once step 2 has actually been executed, so this
section describes something that has been done rather than something believed.
