# Contribute to the Wails v3.0.0 documentation

Use the [M-Press contributor walkthrough](content/contributing/documentation.md)
for downloads, installation on macOS/Linux/Windows, live preview, authoring,
validation, and submitting a pull request. The same guide is published at
[Write documentation with M-Press](https://wails.io/contributing/documentation/).

Documentation-only corrections do not need a prior issue or a failing code test.
Fork `wailsapp/wails`, then create your contribution branch from
`master` and target that branch in your PR.

## Quick start

Install **M-Press v1.0.25** from the
[release archives](https://github.com/leaanthony/mpress/releases/tag/v1.0.25), matching the pinned repository build script.
The walkthrough explains archive checksums, extraction and PATH setup. No Go or
Node.js is needed when using the executable. Alternatively, with Go 1.26.10 or
later (or automatic toolchain selection):

```sh
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.25
mpress version
```

From the Wails repository root:

```sh
mpress dev
```

Open `http://127.0.0.1:3000`. Edit Markdown in `docs/mpress/content/` while the
preview runs. Preserve page metadata and matching `@...` / `@end` components.
Stop the preview with Ctrl+C before validation:

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Use `python` on Windows if that is the installed Python 3 command. Translation
coverage failures block the check; stale translation findings are warnings unless
`--strict` is requested. Existing English pages can be improved while translation
sync catches up; new pages need complete translations in every published locale.
Do not insert English placeholders. Translated changelog segment findings remain
excluded by the checker's existing historical exception; file coverage still applies.

Review your source diff and changed examples, then submit the PR against
`master` with the checks you ran. Layout changes need screenshots.
Generated `docs/mpress/site/` files and local caches stay out of commits.
Deployment credentials and translation providers are not needed for this workflow.

Wails code contributions and feature proposals retain their separate project rules.

## Documentation releases

The published root serves the highest stable version in `docs/mpress/releases.json`.
`/versions/nightly/` tracks the checked-out HEAD; `/versions/v3.0.0/` and other
registered stable URLs serve static snapshots of their Git tags. Local `mpress dev`
continues to preview the working source.

The release tool validates documentation before committing release metadata and
creating a tag. Stable releases also register the documentation version in that
commit; prereleases do not create stored nightly versions. Branch and tag are pushed
atomically. Weekly tag-only releases sync the registry back to master afterwards.
The pipeline then captures and verifies the tagged documentation and stores a
GitHub Actions artifact containing the snapshot, source commit, renderer version,
and SHA256 checksum. Rendered output is never committed.

To generate the retrospective v3.0.0 archive locally:

```sh
python3 docs/mpress/scripts/release_docs.py snapshot --version v3.0.0 --output /tmp/wails-docs-release
```

Set `MPRESS_BIN` to the pinned executable if it is outside PATH. The documentation
workflow also accepts `snapshot_version` when run manually, so an expired artifact
can be regenerated from its existing tag. This does not create or move a release tag.
The v3.0.0 backfill corrects legacy v2 and runtime API link destinations in a
temporary source copy, and records that routing correction in its metadata.
Versioned snapshots use their own canonical URL prefix. Tagged source is unchanged.
The renderer is pinned for repeatable tooling; capture manifests and archive
checksums verify each generated artifact rather than promising byte-identical builds.

For the complete production layout, run `bash docs/mpress/scripts/build.sh` from a
Git checkout with Python 3.12 or newer and Git available. It fetches release tags,
checks translations, generates stable and HEAD views, verifies links and snapshots,
and replaces the output only when both views succeed. Snapshots are cached under
`.mpress/release-docs/`; deleting this cache forces reconstruction. The build updates
cache headers so stable version URLs are immutable and Nightly is revalidated.
This build produces files for the existing hosting integration; it does not deploy.
