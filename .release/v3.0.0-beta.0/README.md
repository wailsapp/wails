# Release record: v3.0.0-beta.0

Working files for the beta. **Nothing here is published.** The comms drafts are
drafts: they carry structure, constraints and the facts that are already
settled, and they leave every claim that needs a human decision blank.

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
