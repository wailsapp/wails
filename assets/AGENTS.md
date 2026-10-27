# Shared Assets Agent Guide

Read the [repository-wide agent instructions](../AGENTS.md) first. This file scopes work under `assets/`.

## Intent

`assets/` contains static artwork shared across the repository, such as logos and images referenced by the READMEs and project websites.

## How to work here

- Check references before renaming, replacing, or removing an asset; one file may be used by multiple versions or sites.
- Prefer updating the source asset and its references together. Keep licensing and attribution requirements intact.
- Do not place generated site output, binaries, or unrelated project-specific resources here.

## In scope

- Shared repository artwork and static files that are intentionally reused by multiple project areas.

## Out of scope

- Website-specific assets, which belong under `website/` or `docs/mpress/` as appropriate.
- Runtime binaries, generated output, and assets for a single unrelated project.
