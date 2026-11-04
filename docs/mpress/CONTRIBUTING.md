# Contribute to the Wails v3.0.0 documentation

Use the [M-Press contributor walkthrough](content/contributing/documentation.md)
for downloads, installation on macOS/Linux/Windows, live preview, authoring,
validation, and submitting a pull request. The same guide is published at
[Write documentation with M-Press](https://wails.io/contributing/documentation/).

Documentation-only corrections do not need a prior issue or a failing code test.
Fork `wailsapp/wails`, then create your contribution branch from
`master` and target that branch in your PR.

## Quick start

Install the **pinned M-Press v1.0.27 release** from the
[release archives](https://github.com/leaanthony/mpress/releases/tag/v1.0.27), matching the repository build script.
The walkthrough explains archive checksums, extraction and PATH setup. No Go or
Node.js is needed when using the executable. Alternatively, with Go 1.26.10 or
later (or automatic toolchain selection):

```sh
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.27
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

## Audit external source and example links

Normal builds and `mpress check` validate local references. They do not establish
whether external links still exist. After building, use an M-Press version that
includes `check external` (v1.0.27 or newer):

```sh
mpress check external --host github.com \
  --url-prefix https://github.com/wailsapp/wails/tree/ \
  --url-prefix https://github.com/wailsapp/wails/blob/ \
  --json > /tmp/wails-external-links.json
```

This explicitly contacts GitHub and redirect destinations. Every unique URL is
requested once per probe sequence, even when repeated across locales, and the
report retains every generated page path. HTTP 404/410 is confirmed broken;
timeouts, 403, rate limits and server failures are inconclusive. Exit 0 means
all selected URLs returned success, 1 means confirmed broken links or an
operational error, and 2 means inconclusive results without confirmed breakage.
Inspect stderr if no report was produced. Retry inconclusive results later;
do not rewrite links or refresh translation exceptions to silence them.
Defaults: 10 seconds per attempt, one additional transient retry, one second
between attempts, at most ten redirects. `Retry-After` defers to a later run.
See `mpress check external --help` for timeout/retry/redirect controls.

The independent **Documentation external links** Actions workflow audits current
`master` daily, on branch deletion and on manual dispatch. This catches URLs
that start failing after a release branch is removed. It does not gate normal
builds, release publishing or Cloudflare deployment. A failed job is visible in
Actions and through maintainers' configured Actions notifications; its summary
lists affected pages, and the `docs-external-links` artifact contains the full
JSON. Both confirmed and inconclusive results fail this independent job with
their distinct exit statuses. GitHub's own scheduling, branch-deletion event
limits and notification preferences apply; the daily run is the backstop.
The Wails audit intentionally selects repository `tree`/`blob` links, covering
source files and examples without thousands of changelog/profile requests.
Remove the URL-prefix and host filters to audit all external HTML `href`/`src`
destinations.

The build and audit use `docs/mpress/mpress-version.txt` (v1.0.27), with
archive checksum and executable version verification. Update that pin and
the setup instructions together when upgrading the toolchain.
For local candidate validation, without downloading a release:

```sh
MPRESS_BIN=/absolute/path/to/candidate/mpress \
  bash docs/mpress/scripts/build.sh --external-audit /tmp/wails-external-links.json
```

The build script also accepts `MPRESS_BIN` for ordinary validation without the
audit flag. The candidate must be an absolute executable path. Routine builds
do not gain an external-link network step.

An HTTP success is not proof that a destination is the correct example. Verify
published Git tags and file/directory types before changing GitHub `tree`/`blob`
links. Remote anchors, dynamically generated links, CSS references and clone
commands in code blocks are not covered by this availability audit.

### Apply a URL-only correction to existing translations

Save the English page from before the correction, update only its inline HTTP(S)
link destinations, then preview one locale without a translation provider:

```sh
git show BASE_REVISION:docs/mpress/content/features/windows/options.md \
  > /tmp/options-before.md
mpress translate sync-links --lang fr --file features/windows/options.md \
  --previous-source /tmp/options-before.md
# Apply after reviewing the source revision and intended destinations:
mpress translate sync-links --lang fr --file features/windows/options.md \
  --previous-source /tmp/options-before.md --write
```

Repeat for each affected page and locale. This narrowly scoped command preserves
translated prose, labels and formatting and refuses mixed prose/layout changes,
damaged translations or URLs a human has already changed independently. It does
not replace missing examples by guesswork. Unchanged translation state and audit
exceptions retain their historical provenance; review any stale/manual findings
and exception hashes explicitly using the existing translation checks. It does
not automatically grant fresh approval or overwrite manual translations through
`--force`. Reference links, raw HTML, autolinks, MPD and clone commands require
manual handling.
