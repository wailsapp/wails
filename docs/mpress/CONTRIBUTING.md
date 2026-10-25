# Contribute to the Wails v3.0.0 documentation

Use the [M-Press contributor walkthrough](content/contributing/documentation.md)
for downloads, installation on macOS/Linux/Windows, live preview, authoring,
validation, and submitting a pull request. The same guide is published at
[Write documentation with M-Press](https://v3.wails.io/contributing/documentation/).

Documentation-only corrections do not need a prior issue or a failing code test.
Fork `wailsapp/wails`, then create your contribution branch from
`releases/v3.0.0` and target that branch in your PR.

## Quick start

Install **M-Press v1.0.19** from the
[release archives](https://github.com/leaanthony/mpress/releases/tag/v1.0.19).
The walkthrough explains archive checksums, extraction and PATH setup. No Go or
Node.js is needed when using the executable. Alternatively, with Go 1.26.6 or
later (or automatic toolchain selection):

```sh
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.19
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
`releases/v3.0.0` with the checks you ran. Layout changes need screenshots.
Generated `docs/mpress/site/` files and local caches stay out of commits.
Deployment credentials and translation providers are not needed for this workflow.

Wails code contributions and feature proposals retain their separate project rules.
