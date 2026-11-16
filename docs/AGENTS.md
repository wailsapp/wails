# Wails Documentation Agent Guide

Read the [repository-wide agent instructions](../AGENTS.md) first. This file scopes work under `docs/`.

## Intent

`docs/mpress/` is the supported source for the Wails v3 documentation site. English source pages live in `content/`; translations live in locale directories; shared static files and styles live alongside them. Generated site output is written to `docs/mpress/site/`.

The separate Docusaurus project in `website/` contains the v2 documentation and versioned historical content. Do not move or duplicate the M-Press v3 documentation there.

## How to work here

- Read [`mpress/CONTRIBUTING.md`](mpress/CONTRIBUTING.md) before substantial documentation changes.
- Edit source files under `mpress/content/` and translation sources, not generated output under `mpress/site/`.
- Keep root `mpress.yaml` and `mpress/mpress.yaml` aligned when changing site-wide configuration; the root configuration supports repository-level tasks and the nested one supports standalone builds.
- Preserve the site's language layout and metadata conventions. Update translations only when requested or when the task's scope includes them; do not invent translated content.
- Preview and validate with the pinned M-Press toolchain documented in `README.md` or `Taskfile.yml`. The standard repository task is `task docs:build`.
- Release publishing owns the generated changelog source and release notes; coordinate changes to those files with the release workflow.

## In scope

- Wails v3 guides, tutorials, API references, changelog content, translations, M-Press configuration, and documentation-specific scripts.
- Documentation corrections and examples that explain existing framework behavior.

## Out of scope

- Framework implementation changes; make those in `v3/` and update docs as needed.
- Docusaurus v2 or historical documentation, which belongs in `website/`.
- Editing generated site output, deployment infrastructure, or release artifacts by hand.
