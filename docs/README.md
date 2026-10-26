# Wails v3.0.0 documentation

This branch prepares the documentation for the Wails v3.0.0 release. Installation
examples target `v3.0.0`; those commands require the release tag to be published.
To work on the release before publication, follow the [source setup guide](mpress/content/contributing/setup.md)
using the `releases/v3.0.0` branch.

The M-Press documentation lives in `docs/mpress/`, with English as the default
language at `/` and existing translations under their language codes.
[Correction PRs are welcome](mpress/CONTRIBUTING.md), including typos, links,
examples, and translations. No issue or failing code test is required for a
documentation-only correction.

## Download, install and contribute

Follow [Write documentation with M-Press](mpress/content/contributing/documentation.md)
for executable downloads, checksum verification, PATH setup on macOS/Linux/Windows,
live preview, adding pages and submitting a PR against `master`. The public site
links to this guide through **Contributing → Documentation** in its sidebar.

## Preview and validate

Install the [latest stable M-Press release](https://github.com/leaanthony/mpress/releases/latest)
from a release archive, or use Go:

```sh
go install github.com/leaanthony/mpress/cmd/mpress@latest
```

From the repository root:

```sh
mpress version
mpress dev
```

Open `http://127.0.0.1:3000` and edit the Markdown sources. Stop the preview
with Ctrl+C before validation. With Python 3 also installed:

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Use `python` on Windows if that is the installed Python 3 command.

The root `mpress.yaml` points to `docs/mpress/content`, static assets, and CSS.
This entry point lets the contribution command work from a cloned repository.
`docs/mpress/mpress.yaml` also supports a standalone build from that directory;
its local contribution wizard is disabled because its paths are project-relative.
Keep both configurations aligned when changing site settings.

English source files live in `docs/mpress/content/`; translated files live under
folders such as `fr/`, `id/`, and `zh-cn/`. Output goes to `docs/mpress/site/`.
The build retains CSS because the imported home-page animation creates classes
at runtime. Contributors do not need deployment credentials, Node.js, private
services, or a translation provider to preview or edit the documentation.

See the [M-Press contribution guide](mpress/CONTRIBUTING.md) for the complete
correction and validation workflow.

M-Press is the only supported v3 documentation source. The nightly release
publisher writes `docs/mpress/content/changelog.md`, and generated release
notes use Markdown paths and metadata.
