# Wails Website Agent Guide

Read the [repository-wide agent instructions](../AGENTS.md) first. This file scopes work under `website/`.

## Intent

`website/` is the Docusaurus 2 website source. It contains the Wails v2 documentation, site pages and components, localized content, and versioned historical documentation. Wails v3's current documentation source is the M-Press project under `docs/`.

## How to work here

- Read this directory's `README.md`, `package.json`, and `Taskfile.yml` before changing build or translation workflows.
- Edit Docusaurus source (`docs/`, `src/`, `static/`, configuration, or locale sources) appropriate to the requested change. Avoid editing generated `build/` output.
- Preserve existing versioned docs and sidebars; do not rewrite archived versions to match current behavior unless requested.
- Use the package manager and lockfile expected by the existing scripts and CI. Check configured Node version in `.nvmrc`.
- For English source changes that need translation, follow the Crowdin workflow in `README.md`; do not claim translations are updated unless locale files were actually updated.
- Verify the change with the relevant website task or npm script when requested or needed by the task.

## In scope

- Docusaurus site content, v2 docs, historical versioned docs, page components, styling, localization, and site configuration.

## Out of scope

- Current v3 documentation source, which belongs in `docs/`.
- Framework implementation, which belongs in `v2/` or `v3/`.
- Deployment credentials, generated site output, or unrelated release automation.
