# Repository Guidelines

## Project Structure & Module Organization
The repository hosts the current Wails tooling in `v3/`, the maintained v2 runtime in `v2/`, and shared assets under `assets/` and `pkg/`. CLI entrypoints live in `v3/cmd/wails3`, while reusable runtime code is grouped under `v3/internal/*` (platform integrations, packaging, templates) and `v3/pkg/*` (application, services, UI helpers). Example applications used for testing and documentation reside in `v3/examples`, with Docker build scripts and harnesses in `v3/test` and targeted regression suites in `v3/tests`. Documentation and site content are stored in `docs/` and `website/`; keep tutorial updates alongside related code changes.

## Build, Test, and Development Commands
Run `task v3:install` to compile and install the v3 CLI locally, and `task v3:precommit` before submitting changes to execute `go test ./...` and apply repository formatters. Use `task v3:test:examples` for a full platform-agnostic build of every example, or scope work with `task v3:test:example:windows DIR=badge` (replace `badge` as needed). Markdown formatting is standardized with `task format:md`, and changelog automation is driven via `task contributors:update` when contributor metadata changes.

## Coding Style & Naming Conventions
Go code must remain `gofmt`-clean with idiomatic package names (lowercase, no underscores) and exported API names following Go's PascalCase. Match existing directory naming by grouping platform-specific code under clearly named packages (e.g., `mac`, `w32`) and keep generator templates in `internal/templates`. Frontend or documentation updates should respect Prettier defaults; run `npx prettier --write "**/*.md"` or the Taskfile wrapper to stay consistent.

## Testing Guidelines
Unit and integration tests accompany Go sources inside the nearest package; add new suites under `v3/tests` for regression coverage. Cross-platform build validation relies on the Taskfile matrix—`task test:examples:all` exercises macOS, Windows, and Linux Docker pipelines, while `task test:example:linux:docker DIR=<name>` verifies a single sample with architecture auto-detection. Prefer deterministic tests and clean up generated binaries (`testbuild-*`) after manual runs.

## Commit & Pull Request Guidelines
Adopt conventional commits recognized by the release tooling (`feat:`, `fix:`, `chore:`, with `!` or `BREAKING CHANGE:` when necessary) so nightly pipelines classify versions correctly. Reference issues or discussion threads in the body, and update `v3/UNRELEASED_CHANGELOG.md` with concise entries under the appropriate heading. Pull requests should summarize scope, list validation commands executed (e.g., `task v3:precommit`), and include platform-specific notes or screenshots for UI-affecting work. Ensure CI passes and that new docs or samples link from the relevant navigation files in `docs/src/content` before requesting review.
