# Generate v3 release notes

Release notes come from merged PRs included between a previous release tag and
the exact tested source commit. There is no shared unreleased changelog file to
edit or revert. `v3/UNRELEASED_CHANGELOG.md` has been removed.

From the repository root, with Git tags fetched and read access through `gh`:

```sh
go -C v3 run ./scripts/generate-release-changelog.go generate --repo wailsapp/wails --target HEAD --channel nightly
```

Use the repository being released; the generator has no public-repository default.
This command only reads GitHub. It writes ignored `v3/release_changelog.md` and
`v3/release_changelog.json`, containing notes and a provenance receipt respectively.
Failures invalidate any earlier generated input and fail the command. Keep the
receipt with the notes for auditing the selected PRs and exact source/baseline.

Nightlies use the nearest reachable v3 release tag. Weekly and explicitly stable
releases use the nearest reachable **stable** v3 tag so changes already shipped
in nightlies appear in the stable release. `--version v3.1.1` selects stable
boundaries even when the channel defaults to nightly. With no eligible reachable
tag, an operator must explicitly choose an ancestor with `--since`; the tool
fails rather than guessing a bootstrap boundary.

Closed PRs are listed once, following all API pages. A PR is included once only
if it is merged into this repository and its
merge commit is actually in the range. This covers normal, squash, and rebase
merges, and excludes open PRs, later merges and other repositories' PRs. PRs
must affect `v3/` or `docs/mpress/`; v2-only PRs are excluded. Direct commits
without an associated merged PR do not produce release-note entries.

Use a clear, user-visible PR title. An optional `## Release notes` section in the
PR description overrides its prose. Conventional types `feat` and `fix` select
Added and Fixed; other user-visible types, including `perf` and `refactor`, select
Changed. `security`, `deprecate` and `remove` select their corresponding sections.
Internal types `build`, `chore`, `ci`, `docs`, `style` and `test` are omitted unless
an explicit `changelog:added|changed|fixed|deprecated|removed|security` label is
present. `skip-changelog` always excludes a PR. Multiple category labels fail the
run. The generator uses no external model or reviewer service.

After generation, the existing release commands consume the same notes:

```sh
cd v3/tasks/release
go run release.go --check-only
go run release.go --extract-changelog
go run release.go --create-release-notes /tmp/release-notes.md
```

The nightly/weekly workflow generates and preserves the input and receipt before
checking whether a release is needed. Missing input, failed GitHub reads or an
incomplete changed-files response fail the pipeline; they do not count as an
empty release. An internal-only range skips an automatic nightly. Explicit
versions and weekly cadence retain their existing empty-range release behavior.

Live releases use the generated notes for both the GitHub release body and a new
versioned section of `docs/mpress/content/changelog.md`. The generated inputs are
never staged or reset. Weekly tag-only releases include the changelog section in
the tag; the subsequent sync copies only that section onto current master,
preserving newer sections and detecting conflicting edits to the same release.
The sync is idempotent and does not recreate a staging file.

To run the focused checks without publishing:

```sh
go -C v3 test ./tasks/changelog
go -C v3 test ./tasks/releaseci
cd v3
go test -race ./tasks/release ./tasks/releasegate
```

The release tests use isolated Git repositories, local bare remotes and a local
HTTP server. They do not create GitHub releases or touch repository branches.
The migration preserves the pending PR #6519 entry in a repository-scoped, immutable migration seed. Its original wording is used only while its merge commit is inside the selected release range; it is not an ongoing staging file.
