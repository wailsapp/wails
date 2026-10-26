# Wails v3 Agent Guide

Read the [repository-wide agent instructions](../AGENTS.md) first. This file scopes work under `v3/`.

## Intent

`v3/` is the Wails v3 Go module. It contains the public application APIs, `wails3` CLI, platform implementations, generators, frontend runtime sources and bundles, templates, examples, tests, and v3-specific release tooling. Its module is defined by `v3/go.mod`.

## How to work here

- Keep v3 implementation and release changes in this tree. Do not assume a matching v2 change is required.
- Read `README.md`, `TESTING.md`, and the relevant platform or subsystem documentation before substantial changes.
- Use the `Taskfile.yaml` tasks for the touched subsystem and consult CI workflows for platform-specific expectations.
- Changes under `internal/runtime/desktop/@wailsio/runtime/src/` require both runtime outputs to be rebuilt: `task runtime:build:assets` and `task runtime:build:package`. The webview asset bundle and npm package are consumed independently; commit both generated outputs when sources change.
- Generated templates, API bindings, event definitions, and bundled assets must be updated through their source generators/tasks, not by hand where a generator owns them.
- Android and iOS support is experimental; do not describe it as part of the stable desktop compatibility contract.
- For subsystem-specific agent notes, follow the nearest `AGENTS.md` (for example, `internal/wake/AGENTS.md`).

## Release and issue scope

- Apply the root release policy: v3.1.0 is feature-only; bug fixes are severity-prioritised and released as nightly patch updates outside that milestone/project.
- Triage the actual issue or PR content. Do not infer priority or release placement from absent labels.
- Use GitHub issues and PRs as the authoritative tracker; do not create local task lists.
- Do not add entries to `UNRELEASED_CHANGELOG.md` or the published v3 changelog in ordinary PRs. They are maintained by the merge/release automation; provide the change's user-visible impact in the PR title and description instead. See the root changelog ownership policy for maintenance exceptions.
- For agent-assisted builds and app testing, follow the root Wails MCP guidance: `wails3 mcp` manages project jobs; the separate `mcp`-tagged application server controls a running app.

## In scope

- Wails v3 framework code, CLI, runtime, platform integrations, examples, tests, generated artifacts, and v3 release tooling.
- v3 documentation only where it is maintained alongside code for a specific subsystem (the main documentation source is `docs/`).

## Out of scope

- Wails v2 code or release decisions; use `v2/`.
- The main M-Press site source in `docs/` or Docusaurus source in `website/`, except when the task explicitly includes documentation or website changes.
- Unrelated changes to the root repository, shared automation, or assets.
