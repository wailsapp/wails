# Starlight Migration Report

## Reviewed migration

Imported with the published M-Press v1.0.1 binary from Wails commit `4146200f1`.
The earlier preview in closed PR #5905 contained generated HTML; this migration
imports the current editable Starlight source instead.

- Preserved 526 imported pages and added the documentation correction guide.
- English is the default at `/`; all nine existing translation trees remain.
- Both contributing route collisions are resolved: the general guides live at
  `/contributing/` and `/id/contributing/`; technical overviews live beneath
  `contributing/overview/`. Legacy `contributing-migrated` URLs redirect.
- The two underscore-prefixed showcase templates remain private in the rollback
  source. The 22 stale fragment warnings below retain valid destination pages.
- Repaired invalid link labels in eight translated home pages, two obsolete
  dock-to-systray links, and unnamed generated sidebar groups.
- Restored the English and Indonesian custom-protocol INI examples: the importer
  had interpreted `[Desktop Entry]` inside a code fence as a new tab. Restored
  the original shell comment markers in their Linux tabs and all home terminals.
- Compared 5,228 source code examples across languages and 1,985 rendered
  examples from the previous site. All match after the repairs, allowing only
  whitespace and filename comments that Starlight displayed as frame captions.
  Intentionally rewritten contribution instructions and blog archive views are
  excluded from this comparison; individual blog posts remain preserved.
- All 194 former Starlight page routes have a page or redirect. Six obsolete
  blog pagination, author and tag views redirect to the blog index.
- Strict builds produce 527 content pages plus a 404 page. The static validator
  checks internal links and anchors, canonical URLs, contribution source paths,
  generated indexes, redirects, and Cloudflare upload limits.
- Browser checks cover desktop/mobile, light/dark themes, language selection,
  search, code tabs and contribution instructions. CSS is retained at build time
  for the home-page animation; mobile branding and search hit areas are repaired.
- Diagram rendering is checked separately from source preservation: all 110 D2
  blocks across 32 pages must become working SVG images. The site validator
  rejects D2 displayed as a code block. Builds use M-Press v1.0.3, which embeds
  D2 and renders these blocks to static SVG without an external installation.

The Starlight source remains available for rollback. Release and changelog
automation now uses MPD paths. Cloudflare builds production from `master`
through the existing GitHub integration. See `../README.md` for configuration.

- Source: `./docs`
- Content source: `docs/src/content/docs`
- Pages imported: 526
- Static assets copied: 130
- Config written: `mpress.yaml`

## Manual Review

The importer found patterns that may need manual cleanup:

- `community/showcase/_template.md` — private-content: Skipped underscore-prefixed Starlight content that is not normally routed.
- `contributing.mdx` — route-collision: Preserved at "/contributing-migrated" because "/contributing" is already provided by contributing/index.mdx.
- `de/concepts/lifecycle.mdx` — stale-fragment: Removed missing fragment #services-lifecycle from a link to /de/concepts/lifecycle; the page link was preserved.
- `fr/concepts/lifecycle.mdx` — stale-fragment: Removed missing fragment #services-lifecycle from a link to /fr/concepts/lifecycle; the page link was preserved.
- `id/blog/2021-09-27-v2-beta1-release-notes.md` — stale-fragment: Removed missing fragment #sponsors from a link to /id/credits; the page link was preserved.
- `id/blog/2021-11-08-v2-beta2-release-notes.md` — stale-fragment: Removed missing fragment #sponsors from a link to /id/credits; the page link was preserved.
- `id/blog/2022-02-22-v2-beta3-release-notes.md` — stale-fragment: Removed missing fragment #sponsors from a link to /id/credits; the page link was preserved.
- `id/community/showcase/_template.md` — private-content: Skipped underscore-prefixed Starlight content that is not normally routed.
- `id/concepts/lifecycle.mdx` — stale-fragment: Removed missing fragment #window-hooks-event-yang-dapat-dibatalkan from a link to /id/concepts/lifecycle; the page link was preserved.
- `id/contributing.mdx` — route-collision: Preserved at "/id/contributing-migrated" because "/id/contributing" is already provided by id/contributing/index.mdx.
- `id/getting-started/installation.mdx` — stale-fragment: Removed missing fragment #legacy-gtk3-support from a link to /id/guides/build/linux; the page link was preserved.
- `id/getting-started/setup.mdx` — stale-fragment: Removed missing fragment #platform-specific-dependencies from a link to /id/getting-started/installation; the page link was preserved.
- `id/getting-started/setup.mdx` — stale-fragment: Removed missing fragment #platform-specific-dependencies from a link to /id/getting-started/installation; the page link was preserved.
- `id/guides/build/building.mdx` — stale-fragment: Removed missing fragment #legacy-gtk3-support from a link to /id/guides/build/linux; the page link was preserved.
- `id/guides/build/linux.mdx` — stale-fragment: Removed missing fragment #linux-dialog-behavior from a link to /id/reference/dialogs; the page link was preserved.
- `id/quick-start/installation.mdx` — stale-fragment: Removed missing fragment #legacy-gtk3-support from a link to /id/guides/build/linux; the page link was preserved.
- `id/reference/application.mdx` — stale-fragment: Removed missing fragment #application-level-windows-options from a link to /id/features/windows/options; the page link was preserved.
- `id/reference/dialogs.mdx` — stale-fragment: Removed missing fragment #legacy-gtk3-support from a link to /id/guides/build/linux; the page link was preserved.
- `id/tutorials/04-self-update-a-wails-app.mdx` — stale-fragment: Removed missing fragment #replace-the-template from a link to /id/guides/updater; the page link was preserved.
- `id/tutorials/04-self-update-a-wails-app.mdx` — stale-fragment: Removed missing fragment #writing-your-own-provider from a link to /id/guides/updater; the page link was preserved.
- `ja/concepts/lifecycle.mdx` — stale-fragment: Removed missing fragment #services-lifecycle from a link to /ja/concepts/lifecycle; the page link was preserved.
- `ko/concepts/lifecycle.mdx` — stale-fragment: Removed missing fragment #서비스-수명-주기 from a link to /ko/concepts/lifecycle; the page link was preserved.
- `pt/concepts/lifecycle.mdx` — stale-fragment: Removed missing fragment #services-lifecycle from a link to /pt/concepts/lifecycle; the page link was preserved.
- `ru/concepts/lifecycle.mdx` — stale-fragment: Removed missing fragment #services-lifecycle from a link to /ru/concepts/lifecycle; the page link was preserved.
- `zh-cn/concepts/lifecycle.mdx` — stale-fragment: Removed missing fragment #services-lifecycle from a link to /zh-cn/concepts/lifecycle; the page link was preserved.
- `zh-tw/concepts/lifecycle.mdx` — stale-fragment: Removed missing fragment #services-lifecycle from a link to /zh-tw/concepts/lifecycle; the page link was preserved.
