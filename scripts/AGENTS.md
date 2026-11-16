# Repository Scripts Agent Guide

Read the [repository-wide agent instructions](../AGENTS.md) first. This file scopes work under `scripts/`.

## Intent

`scripts/` contains repository-level automation helpers, including issue/PR triage and labeling utilities. These scripts support maintainers; they are not part of the Wails application runtime.

## How to work here

- Inspect each script's help and current behavior before changing it. Preserve safe defaults and make side effects clear.
- Keep scripts consistent across supported shells/platforms when parallel `.sh` and `.ps1` implementations exist.
- For scripts that read or mutate GitHub issues, PRs, labels, or milestones, follow the root issue-tracking and release policy and use narrow, explicit selection criteria.
- Do not run a script that mutates external GitHub state unless the user explicitly asked for that action.
- Document new required tools and invocation details in the relevant script README or usage text.

## In scope

- Repository maintenance scripts, their usage documentation, and focused validation for those scripts.

## Out of scope

- Runtime, CLI, framework, and website behavior; those belong to their owning directories.
- One-off planning notes or parallel issue databases.
- External state changes not explicitly requested by the user.
