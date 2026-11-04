# Supporting Tools Agent Guide

Read the [repository-wide agent instructions](../AGENTS.md) first. This file scopes work under `tools/`.

## Intent

`tools/` contains small supporting utilities used for repository presentation or maintenance, such as SponsorKit and star-history tooling. These are auxiliary projects, not Wails runtime packages.

## How to work here

- Inspect the specific tool's README, manifests, and scripts before changing it; each tool may have its own runtime and dependency workflow.
- Keep changes isolated to the tool being maintained and avoid introducing dependencies into the Go modules unless required.
- Do not upload, publish, or alter third-party service state without explicit user direction.

## In scope

- Supporting utility source, configuration, and documentation under this directory.

## Out of scope

- Wails framework or website implementation.
- Generated outputs, secrets, credentials, and external publishing actions unless explicitly requested.
