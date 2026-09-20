# Release record: v3.0.0-beta.0

Working files for the beta. **Nothing here is published.** The comms drafts are
drafts: they carry structure, constraints and the facts that are already
settled, and they leave every claim that needs a human decision blank.

## GitHub release train

Issue [#5844](https://github.com/wailsapp/wails/issues/5844) is the authoritative
coordination issue. Repository milestones mirror this record:

- [Beta 0 staging](https://github.com/wailsapp/wails/milestone/23) — 11 August 2026
- [Beta 1 public beta](https://github.com/wailsapp/wails/milestone/17) — 25 August 2026
- [Beta 2 code freeze](https://github.com/wailsapp/wails/milestone/20) — 1 September 2026
- [Beta 3 final stabilisation](https://github.com/wailsapp/wails/milestone/21) — 8 September 2026
- [RC1 dry run](https://github.com/wailsapp/wails/milestone/18) — 12 September 2026
- [v3.0.0 GA](https://github.com/wailsapp/wails/milestone/10) — 15 September 2026

| File | What it is |
|---|---|
| `plan.json` | The release plan: schedule, gates, and the audit configuration for this repository. Read and written by `relman`. |
| `comms/launch-sequence.md` | The ordered launch-day runbook, with the Wails-specific steps and the abort criteria. |
| `comms/packagers.md` | Heads-up for downstream packagers, plugin and template authors and translators. Facts are filled in, dates are not. |
| `comms/faq.md` | Objection and question list to answer before the launch thread, not during it. |
| `comms/announcement.md`, `comms/show-hn.md`, `comms/social-*.md`, `comms/email.md`, `comms/discord-forum.md`, `comms/reddit.md` | Per-channel scaffolds with the length limits and audience notes for each. Copy is deliberately absent. |

## What is deliberately empty

`plan.json` has empty `positioning` and `headlines`, and the channel drafts have
TODOs where those would be interpolated. That is not an oversight: the one-liner,
the problem statement, the proof point and the call to action are the release's
argument for itself, and inventing them would produce copy that reads plausibly
and says something nobody decided to say.

Fill those five fields and re-run `relman kit -plan .release/v3.0.0-beta.0/plan.json`
to regenerate every channel draft with real copy.

## Audit state at the time of writing

Against this branch: 17 checks passing, 0 warnings, 2 blockers.

- `changelog-entry` clears at release time, when `release.go` folds
  `UNRELEASED_CHANGELOG.md` into the changelog.
- `ci-green` requires the branch name in `plan.json` to match the branch CI runs
  against.

## Known preconditions before tagging

1. The six `APPLE_*` signing secrets do not exist on the repository. Without them
   the macOS job stops in preflight with the names listed, or continues unsigned
   if `allow_unsigned_macos` is set deliberately.
2. `release-v3.yml` has never run. Rehearse it with `draft: true` on a throwaway
   tag in the `alpha2.` series before using it for real.
3. This branch drifts behind master quickly. Re-sync and re-audit immediately
   before tagging; `relman audit` now blocks on it.

## Findings worth acting on after the beta

Noticed while preparing this release. None are release-blocking; all of them
cost something later if forgotten.

**v3 CI compiles the GTK4 backend against exactly one GTK version.** The v3 test
matrix is `[windows-latest, ubuntu-latest, macos-latest]`, so `linux_cgo.go` is
only ever built against whatever GTK ubuntu-latest carries. The single
`ubuntu-22.04` job in the repository lives in `build-and-test.yml`, which is the
v2 workflow and never compiles that file. A contributor can therefore use a GTK
symbol newer than the project's real floor and CI will pass.

In practice the floor is set by `#cgo linux pkg-config: gtk4 webkitgtk-6.0`:
requiring WebKitGTK 6.0 already excludes the older distributions, which is why
this has not bitten yet. The cheap mitigation is to state a minimum GTK4 version
in the installation docs so contributors know what they may use, rather than
adding a matrix entry for distributions that cannot satisfy the WebKitGTK
requirement anyway.

**Fork pull requests cannot converge while the changelog bot is active.** Master
requires branches to be up to date, the changelog bot pushes a commit after a
merge, and a fork PR brought up to date needs its workflow runs approved by hand
before any required check reports. Each fork PR therefore needs a roughly
twenty-minute cycle that anything else merging invalidates. There were 500 runs
sitting in `action_required` while this release was being prepared. Merging in a
quiet window works and is what was done here, but it does not scale; the durable
fixes are settings changes.

**`go build ./...` from `v3/` fails on several `examples/` packages**, equally on
master and on this branch, so it is not a regression. CI stays green because it
builds examples through `task test:examples`. It does mean the most obvious
command a new contributor runs does not work.

**Six diagram placeholders render in the documentation.**
`contributing/index.mdx`, `contributing/architecture.mdx` and
`concepts/build-system.mdx` contain literal `**[... Diagram Placeholder]**` text
where D2 diagrams were removed after MDX parse failures. Contributor-facing
rather than user-facing, so the diagrams were left alone rather than
reconstructed speculatively.
