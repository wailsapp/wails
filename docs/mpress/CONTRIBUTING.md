# Correct the Wails v3 documentation

Correction PRs are welcome: fix a typo, repair a link, clarify an explanation,
update an example, or correct a translation. You do not need to open an issue
or add a failing code test for a documentation-only correction.

1. Fork `wailsapp/wails` and create a documentation branch from `master`.
2. Install M-Press v1.0.3 from its release archives, or run:

   ```sh
   go install github.com/leaanthony/mpress/cmd/mpress@v1.0.3
   ```

3. From the repository root, run `mpress dev`. Edit the `.mpd` files in
   `docs/mpress/content/`; English lives directly in that directory and
   translations live under their language codes (for example `fr/`).
4. Run `mpress build --strict --no-purge-css`, `mpress check`, and
   `python3 docs/mpress/scripts/check_site.py docs/mpress/site`.
5. Open a PR against `master`, explain the correction and list the checks run.
   Include before/after screenshots for layout changes. Verify any code examples
   you change, and identify the platform and Wails version used.

The production site is static. You do not need Node.js, Cloudflare credentials,
a translation provider, or private M-Press services to preview or correct it.
Keep existing translations; do not replace them with machine translations for
an unrelated correction. Missing translations link to the English page.

MPD keeps headings, paragraphs, lists and fenced code readable. Preserve its
`---` metadata block and matching `@...` / `@end` component boundaries. See
https://github.com/leaanthony/mpress for the authoring reference.

`docs/src/content/docs/` remains the Starlight rollback source.
The M-Press site is built from `docs/mpress/content/`; do not edit generated
`docs/mpress/site/` files.

Wails code contributions and feature proposals retain the project's existing
rules; this guide applies to documentation corrections.
