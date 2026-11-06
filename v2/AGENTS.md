# Wails v2 Agent Guide

Read the [repository-wide agent instructions](../AGENTS.md) first. This file scopes work under `v2/`.

## Intent

`v2/` is the Wails v2 framework: its public Go APIs, CLI, platform backends, frontend/runtime support, templates, examples, and tests. It is a distinct maintained release line with its own module (`v2/go.mod`) and conventions.

## How to work here

- Make v2 changes in this tree and keep its module, dependencies, and public API compatibility in mind.
- Check `README.md`, `go.mod`, and `Taskfile.yaml` before changing build or release behavior.
- Common tasks include `task v2:lint`, `task v2:format`, and `go test ./...` from this directory, subject to platform dependencies.
- Update examples or templates when an API change affects how applications are created or used.
- Coordinate public API and behavior changes with the relevant GitHub issue or PR under the root tracking policy.

## In scope

- Wails v2 code, tests, templates, examples, module metadata, and v2-specific release tooling under this directory.
- Fixes to v2 compatibility or platform behavior that can be validated against this module.

## Out of scope

- Wails v3 implementation or release policy; use `v3/` for those changes.
- The current M-Press v3 documentation source in `docs/` or the Docusaurus source in `website/`, unless the request explicitly asks for v2 website/documentation work.
- Repository-wide CI, policy, or shared assets unless the task specifically requires changing them.
