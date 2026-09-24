# Correct the Wails v3 documentation

Correction PRs are welcome: fix a typo, repair a link, clarify an explanation,
update an example, or correct a translation. You do not need to open an issue
or add a failing code test for a documentation-only correction.

1. Fork `wailsapp/wails` and create a documentation branch from `master`.
2. Install M-Press v1.0.17 from its release archives, or run:

   ```sh
   go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
   ```

3. From the repository root, run `mpress dev`. Edit the `.md` files in
   `docs/mpress/content/`; English lives directly in that directory and
   translations live under their language codes (for example `fr/`).
4. Run `python3 docs/mpress/scripts/check_translations.py`,
   `mpress build --strict --no-purge-css`, `mpress check`, and
   `python3 docs/mpress/scripts/check_site.py docs/mpress/site`.
5. Open a PR against `master`, explain the correction and list the checks run.
   Include before/after screenshots for layout changes. Verify any code examples
   you change, and identify the platform and Wails version used.

The production site is static. You do not need Node.js, Cloudflare credentials,
a translation provider, or private M-Press services to preview or correct it.
Every published language must contain a complete translation of every English
page. Do not copy English prose into a language folder or rely on fallbacks.
A translation correction can change just the affected language. If you change
English meaning, update the corresponding translations; unrelated pages do not
need to be regenerated.

The translated `changelog.md` files are retained for publication and coverage,
but their historical segment audit is excluded until the generated release notes
are re-synchronised. Translation audits remain strict for every other page.

Markdown keeps headings, paragraphs, lists and fenced code readable. Preserve its
`---` YAML metadata block and matching `@...` / `@end` component boundaries. See
https://github.com/leaanthony/mpress for the authoring reference.

The M-Press site is built from `docs/mpress/content/`; do not edit generated
`docs/mpress/site/` files.

Wails code contributions and feature proposals retain the project's existing
rules; this guide applies to documentation corrections.
