# sponsorkit

A dependency-free Go replacement for the Node [sponsorkit](https://github.com/antfu-collective/sponsorkit)
package. It fetches live GitHub data and renders two self-contained SVGs:

- `website/static/img/sponsors.svg` — the sponsors image embedded in the project
  READMEs and both docs sites.
- `website/static/img/contributors.svg` — the contributors mosaic embedded on the
  credits pages of both docs sites.

## Sponsors mode (default)

1. Queries the GitHub GraphQL API for all active sponsors of the configured login.
2. Buckets them into tiers by monthly amount (see `config.go`). Bigger sponsors get
   bigger avatars and fancier treatments: animated gradient rings, glow halos,
   light sweeps, orbiting sparkles at the top; simple circles at the bottom.
3. Downloads each avatar at 2x resolution and re-encodes it as JPEG (Go stdlib only),
   embedding everything as data URIs so the SVG is fully self-contained.
4. Renders a dark-card SVG. Every avatar is a link to the sponsor's profile with a
   tooltip, animations are SMIL (they run even inside GitHub's image proxy), and a
   "Become a Sponsor" button links to the sponsors page.
5. Gold+ sponsors get a hover "power-up" (lift, glare sweep across the avatar, glow
   bloom, a fast light ring, an emanating pulse and a name underline). Hover and the
   per-sponsor links only work where the SVG loads as a document, which is why the
   website embeds it with `object` rather than `img`; inside `img` embeds the hover
   layers stay hidden.

## Contributors mode (`-mode contributors`)

1. Fetches the repository's contributors from the REST API (bots excluded) and
   counts each author's merged pull requests via GraphQL — one merged PR is one
   unit of work however it landed, so squash-merged and merge-committed work
   rank alike (`-metric commits` switches back to raw default-branch commit
   counts, which reward granular unsquashed histories).
2. Scans the v2 and v3 changelogs for `@login` credits. Squash-merged or
   hand-applied patches are often credited only there, so changelog-only
   contributors still appear; markdown profile links (`[@x](https://github.com/y)`)
   trust the URL rather than the link text, and every changelog-only login is
   validated against the API so typos and organisations are dropped. A
   contributor's credit is the larger of their metric count and mention count;
   contributors with commits but no recorded PRs (early direct pushes) stay
   visible in the tail bands.
3. Renders a mosaic of superellipse "squircles" on the same dark card, graded by
   credit into bands (see `bands` in `config.go`): the most prolific contributors
   get large named squircles with animated gradient rings and travelling light
   arcs, the long tail gets small plain squircles. Every squircle links to the
   contributor's profile.
4. Hover on the bigger bands lifts the squircle, blooms the ring and pops up a
   chip showing the commit (or changelog credit) count.
5. With hundreds of avatars the embedded images dominate the file size, so small
   bands use lower JPEG quality, flat-colour identicons keep their original PNG
   when it is smaller, and the squircle geometry is shared via `defs`/`use`.

## Usage

```sh
cd tools/sponsorkit
SPONSORKIT_GITHUB_TOKEN=<token> GOWORK=off go run . -out ../../website/static/img/sponsors.svg
SPONSORKIT_GITHUB_TOKEN=<token> GOWORK=off go run . -mode contributors -metric prs \
  -changelogs ../../docs/src/content/docs/changelog.mdx,../../website/src/pages/changelog.mdx \
  -out ../../website/static/img/contributors.svg
```

Flags:

| Flag          | Default          | Purpose                                          |
|---------------|------------------|--------------------------------------------------|
| `-mode`       | `sponsors`       | `sponsors` or `contributors`                     |
| `-login`      | `leaanthony`     | GitHub account whose sponsors to render          |
| `-repo`       | `wailsapp/wails` | Repository whose contributors to render          |
| `-metric`     | `prs`            | Contributor ranking: `prs` or `commits`          |
| `-changelogs` | (empty)          | Comma-separated changelogs scanned for `@login` credits |
| `-out`        | `sponsors.svg`   | Output path                                      |
| `-width`      | `800`            | SVG width in CSS pixels                          |
| `-scale`      | `2`              | Avatar oversampling for hi-dpi                   |
| `-quality`    | `80`             | JPEG quality for embedded avatars                |

For sponsors, the token must have sponsor-tier visibility for the account (the
maintainer's own token, e.g. the `SPONSORS_TOKEN` secret in CI); without it every
sponsor lands in the catch-all "Helpers" tier. For contributors any token works.
`GITHUB_TOKEN` is used as a fallback env var.

Tier and band thresholds and styling live in `config.go`; the visual language
(gradients, filters, animations) lives in `render.go` and `render_contributors.go`.

Both images are regenerated daily by
`.github/workflows/generate-sponsor-image.yml`, which commits them when they
change.

## Where the images are used

The single source of truth is `website/static/img/` (deployed to
`https://wails.io/img/`). The images are referenced from:

- `README.md` and all `README.*.md` translations — `img` with the repo-relative
  path to `sponsors.svg` (GitHub sanitises README HTML, so `object` embeds and
  therefore hover and in-image links are not possible there).
- v2 docs (`website/src/pages/credits.mdx` and the 11
  `website/i18n/*/docusaurus-plugin-content-pages/credits.mdx` copies) — `object`
  embeds of `/img/sponsors.svg` and `/img/contributors.svg`.
- v3 docs (`docs/src/content/docs/credits.mdx` and the 9 locale copies) — `object`
  embeds of `https://wails.io/img/sponsors.svg` and
  `https://wails.io/img/contributors.svg`, so they always show the latest
  deployed images.

If you add a new place that shows either image, reference that same file and
prefer an `object` embed so the per-person links and hover effects work.
