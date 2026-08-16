# Wayfinder map: Wails manifest build system

Status: open
Type: map

## Destination

An implementation-ready specification for replacing generated v3 Taskfiles and
`build/config.yml` with a single manifest-driven Wails build system, using Wake
for graph execution, caching, parallelism, and reporting. The specification is
complete when the implementation can begin without unresolved product or
architecture decisions.

## Notes

- Domain: Wails v3 build tooling and migration.
- Tracker: local Markdown files in this directory, following
  `docs/agents/issue-tracker.md`.
- Standing decisions from the conversation: one extendable `wails.toml`; no
  long-term `wails3 task`; profiles are opt-in; `wails3 eject` fully freezes
  resolved defaults; shell customization calls user-owned scripts; platform
  configuration is generated at build time; inline Taskfile shell migration is
  manual; default and modified Taskfiles should be migrated where possible.
- Prefer a narrow, measurable vertical-slice MVP over resolving every edge in
  advance. This map may carry that prototype far enough for hands-on review;
  production implementation remains outside the planning tickets.
- Consult the existing proposal at
  [`history/wails-toml-build-system.md`](../../wails-toml-build-system.md).
- Use the `wayfinder`, `grilling`, and `domain-modeling` skills while resolving
  decision tickets.

## Decisions so far

- [Manifest vocabulary and fully frozen eject semantics](issues/01-manifest-and-eject.md)
  — `wails.toml` stays minimal until customized, then resolves through typed
  Platform/Target/Profile overlays; `wails3 eject [profile]` creates an
  independently frozen, host-neutral snapshot with safe three-way upgrade
  suggestions.
- [Wake's typed build graph and execution boundary](issues/02-wake-domain-graph.md)
  — Wake produces one immutable multi-Target Plan of typed Nodes, deduplicates
  shared work, runs in-process handlers through a critical-path scheduler, and
  keeps Taskfile semantics isolated in legacy migration code.
- [Automatic input discovery and cache identity](issues/03-cache-and-input-discovery.md)
  — typed handlers derive portable content-addressed Action Keys from narrow
  automatic inputs, propagate Artifact content rather than transitive sources,
  and use persistent snapshots plus local cross-project storage for fast hits.

## Not yet specified

- The migration report format and the set of recognizable modifications that
  can be translated without preserving Taskfile semantics.
- The compatibility and release policy for projects before and after
  migration.
- The acceptance matrix across desktop, server, iOS, Android, packaging,
  signing, cross-compilation, and development mode.

## Out of scope

- Designing a general-purpose replacement for Task.
- Preserving arbitrary Taskfile syntax as a second supported build language.
- Rewriting unrelated v3 application/runtime APIs.
