# sponsorkit

A dependency-free Go replacement for the Node [sponsorkit](https://github.com/antfu-collective/sponsorkit)
package. It fetches live GitHub Sponsors data and renders `website/static/img/sponsors.svg`,
the image embedded in the project READMEs.

## What it does

1. Queries the GitHub GraphQL API for all active sponsors of the configured login.
2. Buckets them into tiers by monthly amount (see `config.go`). Bigger sponsors get
   bigger avatars and fancier treatments: animated gradient rings, glow halos,
   light sweeps, orbiting sparkles and tier badges at the top; simple circles at the bottom.
3. Downloads each avatar at 2x resolution and re-encodes it as JPEG (Go stdlib only),
   embedding everything as data URIs so the SVG is fully self-contained.
4. Renders a dark-card SVG. Every avatar is a link to the sponsor's profile with a
   tooltip, animations are SMIL (they run even inside GitHub's image proxy), and a
   "Become a Sponsor" button links to the sponsors page.

## Usage

```sh
cd tools/sponsorkit
SPONSORKIT_GITHUB_TOKEN=<token> GOWORK=off go run . -out ../../website/static/img/sponsors.svg
```

Flags:

| Flag       | Default        | Purpose                                  |
|------------|----------------|------------------------------------------|
| `-login`   | `leaanthony`   | GitHub account whose sponsors to render   |
| `-out`     | `sponsors.svg` | Output path                               |
| `-width`   | `800`          | SVG width in CSS pixels                   |
| `-scale`   | `2`            | Avatar oversampling for hi-dpi            |
| `-quality` | `80`           | JPEG quality for embedded avatars         |

The token must have sponsor-tier visibility for the account (the maintainer's own
token, e.g. the `SPONSORS_TOKEN` secret in CI). `GITHUB_TOKEN` is used as a fallback
env var. Without tier visibility every sponsor lands in the catch-all "Helpers" tier.

Tier thresholds and styling live in `config.go`; the visual language (gradients,
filters, animations) lives in `render.go`.

This tool is run daily by `.github/workflows/generate-sponsor-image.yml`, which
commits the regenerated SVG when it changes.
