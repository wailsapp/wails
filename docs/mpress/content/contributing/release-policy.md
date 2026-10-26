---
title: "Release Policy"
description: "How Wails v3 versions its nightly and weekly releases, how maintainers cut them, and why users should pin versions."
slug: "contributing/release-policy"
sourcePath: "contributing/release-policy.md"
---

## The model in one table

Wails v3 releases on a **clock, not on features**. After the `v3.0.0` baseline:

| Release | Tag | Cadence | GitHub release | npm `@wailsio/runtime` |
| --- | --- | --- | --- | --- |
| Nightly | `v3.<next>.0-nightly.<YYYYMMDD>` (UTC date; `.N` appended if re-run the same day) | automated, daily, only when `UNRELEASED_CHANGELOG.md` has content | pre-release, **not** Latest | same version without `v`, dist-tag `nightly` |
| Weekly | `v3.<next>.0` | automated, Sunday 20:00 Australia/Sydney **(day unconfirmed)** | normal release, Latest | same version, dist-tag `latest` |
| Hotfix | `v3.<current>.<N>`, `N > 0` | rare, manual | normal release, Latest | same version, dist-tag `latest` |
| Major | `v4.0.0` | only for breaking changes | normal release | `latest` |

`<next>` is the highest **released** v3 minor plus one. A nightly is a preview of the **next weekly**, which is why its base is the highest release tag (pre-releases are ignored): after a hotfix `v3.1.1` the next nightly is still `v3.2.0-nightly.*`. It has to name a higher version because a pre-release sorts below its own release: `v3.0.0-nightly.<date>` would sort before `v3.0.0`.

**A minor bump is not a feature milestone.** `v3.4.0` is "the weekly release after `v3.3.0`". Read the release notes, not the version number, to see what changed.

The policy is chosen by version **shape**: a `-nightly.` suffix is a nightly; any other suffix (`-alpha`, `-beta`, `-rc`) is a legacy pre-GA build; no suffix is a release.

### Lockstep versions

The CLI, the Go module and the npm package `@wailsio/runtime` carry the **same version**, taken from `v3/internal/version/version.txt` (`3.2.0-nightly.20261102` on npm for tag `v3.2.0-nightly.20261102`). Projects created with `wails3 init` pin `@wailsio/runtime` to the CLI's own version; development builds of the CLI fall back to `latest`.

## Worked example

| Date | Event | Tag | `go install ...@latest` resolves to |
| --- | --- | --- | --- |
| Sun 25 Oct | GA, manual release | `v3.0.0` | `v3.0.0` |
| Mon 26 - Sat 31 Oct | nightlies | `v3.1.0-nightly.20261026` ... `v3.1.0-nightly.20261031` | `v3.0.0` |
| Sun 1 Nov 20:00 Sydney | weekly | `v3.1.0` | `v3.1.0` |
| Mon 2 Nov | nightly (previews the next weekly) | `v3.2.0-nightly.20261102` | `v3.1.0` |
| Wed 4 Nov | hotfix, manual, from a release branch | `v3.1.1` | `v3.1.1` |
| Sun 8 Nov | weekly | `v3.2.0` | `v3.2.0` |

A weekly always sorts above every nightly of its own number, and a nightly never becomes `@latest`.

## What users run

| Goal | Command |
| --- | --- |
| Newest release (weekly or hotfix) | `go install github.com/wailsapp/wails/v3/cmd/wails3@latest` or `wails3 update cli` |
| **Reproducible (recommended)** | pin an exact tag: `v3.1.0` in `go.mod`, CI and docs |
| A nightly | its exact tag (`@v3.2.0-nightly.20261102`), or `wails3 update cli -nightly` |
| Newest build of any kind | `wails3 update cli -pre` |
| Runtime npm package | `@wailsio/runtime` (`latest` = weekly/hotfix, `nightly` = nightly) |

Why this works: Go resolves `@latest` to the highest tag **without a pre-release suffix**, and only considers pre-releases when no release exists. The GitHub "pre-release" checkbox is invisible to the Go proxy, so nightlies carry a semver pre-release suffix instead of relying on it. Always pin in production.

## Branches

| Branch | Purpose |
| --- | --- |
| `master` | Development and nightlies. |
| `releases/**` | Short-lived branches for one-off releases: the `v3.0.0` GA and hotfixes. Same required CI as `master`. |

Weeklies are cut from `master`, but **not from its tip**: the workflow picks the newest commit since the previous weekly whose required checks (`v3 Build Results`, `Run Go Tests`, `Bundled runtime matches source`) all passed, and tags exactly that commit. If none is green, nothing is released. Nightlies are cut from `master` and only warn when HEAD is not green.

## How to cut a release

Release rights are limited to the project Lead (see the governance document). Contributors never tag or publish; they add entries to `v3/UNRELEASED_CHANGELOG.md`.

Everything runs from **Nightly + Weekly Release v3** (`.github/workflows/nightly-release-v3.yml`):

- **Nightly (automatic)** at 15:00 UTC. Skipped when the unreleased changelog is empty.
- **Weekly (automatic)** at 20:00 Sydney time. GitHub cron is UTC-only, so two crons (09:00 and 10:00 UTC) are scheduled and a guard job proceeds only when it is 20:00 in Sydney (DST-aware). The weekly passes the release gate, selects the newest green commit, bumps the minor, tags that commit without touching the changelog, then updates only the version fields on `master` (`version.txt` and the runtime package versions), so dependency and lockfile updates merged after the selected commit are kept. npm publishes from the tag, not from `master`.
- **Manual.** Run the workflow with `channel`, `dry_run` (default on) and optionally `release_version`. A dry run prints version, kind, GitHub flags, npm dist-tag and notes without touching GitHub or npm.
- **GA or hotfix.** Set `release_version` (`v3.0.0`, `v3.1.1`). A hotfix for an older line can also set `ref: releases/...`; its release commit is cherry-picked back to `master`.

Local dry run: `cd v3/tasks/release && go run release.go --dry-run --channel weekly`; `go run release.go --print-policy v3.2.0-nightly.20261102` shows the derived flags.

### Safety nets

- An automatic version must be strictly greater (semver order, pre-release < release) than every existing tag, or the tool refuses.
- One release job at a time (shared concurrency group).
- npm publishing is triggered by the release **tag** and builds the tagged commit.

## Release gate

The **weekly** (never the nightly) refuses to tag while any of these exist:

1. open issues or pull requests labelled `P0` or `release-blocker`;
2. open issues labelled `regression` created since the last weekly;
3. open Dependabot alerts of severity critical or high, open code-scanning alerts of critical or high severity, and any open secret-scanning alert, **only for shipped code**. Alerts located solely under `v3/examples/`, `v2/examples/`, `website/`, `docs/`, `.github/` or any `examples/` directory are ignored (`SHIPPED_ALERT_IGNORE` in `v3/tasks/releasegate/main.go`, **owner to review**);
4. repository security advisories still in `triage`.

**Who can set a blocker.** A label counts only if its most recent `labeled` event (with no later `unlabeled`) was made by a member of the `wailsapp/maintainers` team or a collaborator with `maintain`/`admin` permission. The same label applied by anyone else is ignored and logged, so outsiders cannot block or unblock a release. The auto-applied `high-priority` label is deliberately **not** used. Bot accounts such as `taliesin-ai` pass the allowlist if they hold `maintain`; this is intended, and the triage agent may apply blocker labels with evidence.

**Failure behaviour.** The gate fails **closed**: if an alert API or the allowlist cannot be read, the weekly stops. On failure it tags nothing, writes the list (with links) to the job summary, and updates a comment on the single "Weekly release status" issue when the repository variable `WEEKLY_STATUS_ISSUE` is set.

**Setup the owner must do on GitHub** (nothing is created automatically): create the labels `P0`, `release-blocker` and `regression`; create the "Weekly release status" issue and set `WEEKLY_STATUS_ISSUE` to its number; give `WAILS_REPO_TOKEN` read access to Issues, Pull requests, Administration (collaborator permission), Dependabot alerts, Code scanning alerts, Secret scanning alerts, Repository security advisories and organisation Members, plus Issues write for the status comment. The default `GITHUB_TOKEN` cannot read these.

**Override.** `override_reason` (workflow-dispatch input) skips the gate. The reason, the actor and every skipped blocker are written to the summary and the status issue. Only the release Lead should use it. GitHub cannot restrict an input to one person, so run the release job in an environment with required reviewers (recommended) so an override needs explicit approval.

## Related

- [Testing & Continuous Integration](/contributing/testing-ci/)
- [Maintainer Triage](/contributing/maintainer-triage/)
