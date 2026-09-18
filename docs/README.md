# Wails v3 documentation

The M-Press documentation lives in `docs/mpress/`, with English as the default
language at `/` and existing translations under their language codes.
[Correction PRs are welcome](mpress/CONTRIBUTING.md), including typos, links,
examples, and translations. No issue or failing code test is required for a
documentation-only correction.

## Preview and validate

Install [M-Press v1.0.17](https://github.com/leaanthony/mpress/releases/tag/v1.0.17)
from a release archive, or use Go:

```sh
go install github.com/leaanthony/mpress/cmd/mpress@v1.0.17
```

From the repository root:

```sh
mpress dev
mpress build --strict --no-purge-css
mpress check
```

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
