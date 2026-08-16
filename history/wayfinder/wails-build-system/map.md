# Wayfinder map: Wails manifest build system

Status: resolved
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
- The narrow cache slice has now been expanded into a complete implementation
  surface for hands-on backwards review. Matching-host release verification
  remains a release operation governed by the acceptance contract, not an
  unresolved build-system decision.
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
- [Manifest-to-binary cache fast-path MVP](issues/11-cache-fast-path-mvp.md)
  — the executable slice validated a 60–70ms no-op path, content-based
  Artifact propagation, and missing-output restoration; it now graduates into
  the complete production implementation.
- [Runtime-generated platform assets and customization escape hatches](issues/04-generated-platform-assets.md)
  — typed metadata generates target-local disposable assets, while every
  package format supports a strict, versioned user-owned template adapter.
- [Current build performance baseline](issues/08-performance-baseline.md)
  — reproducible complete-command and graph microbenchmarks establish no-op,
  incremental, packaging, restoration, Dev, Task, and larger-app timings.
- [Taskfile migration and customization classification](issues/06-taskfile-migration.md)
  — private migration reports classify current, historical, and customised
  Taskfiles; complete migrations digest-check and retire represented sources.
- [CLI migration cutover](issues/07-cli-compatibility.md)
  — one routing seam selects native Manifest, deliberate legacy fallback, or
  an actionable ambiguity error without storing workflow state in `wails.toml`.
- [Script hook contract and safe extension boundary](issues/05-script-hook-contract.md)
  — six typed file-script phases have fixed Project/Target/package barriers,
  stable environment and failure semantics, safe cancellation, and explicit
  bounded cache opt-in.
- [Dev Session invalidation and process lifecycle](issues/10-dev-session-lifecycle.md)
  — finite generations drive a transactional, readiness-gated frontend,
  backend, and watch lifecycle with no-op preservation and clean cancellation.
- [Build performance acceptance matrix and regression budget](issues/09-performance-acceptance.md)
  — permanent semantic, variance, latency, concurrency, native-host, mobile,
  and signing gates make speed and platform correctness verifiable.

## Implementation evidence

- Minimal config, base/profile ejection, migration/cutover, generated assets,
  hooks, target overlays, build/package/sign plans, and the manifest dev
  lifecycle are implemented behind the normal CLI commands.
- Follow-up two-axis review closed cache identity for local Go sources and
  relevant environment, separated Receipts from Artifacts, made signing
  non-reusable, added missing-profile/re-eject semantics, combined multiple
  Targets in one Plan, translated safe script-file hooks, and made Dev rebuilds
  generation-cancelable with process/watch reconfiguration.
- Badge no-op builds measure 70–80ms and DEB package no-ops measure 70–90ms
  wall time with all planned Nodes cached; a service method-body edit skips
  bindings and the frontend and executes only Go compilation.
- The final controlled run measured a 93.503ms badge median with 8.19% MAD,
  3.22% cold orchestration overhead, concurrent native Linux package adapters,
  and a 0/9 cached rerun after AppImage packaging.
- The full Go suite reaches unrelated desktop-environment failures only;
  focused manifest/Wake/command tests, race detection, vet, and real Linux
  build/package runs pass.
- A local draft WEP now records the public behavior and rollout proposal under
  `v3/wep/proposals/manifest-build-system/`; no GitHub issue, PR, or push was
  created.

## Not yet specified

None. The route to the destination is fully specified.

## Out of scope

- Designing a general-purpose replacement for Task.
- Preserving arbitrary Taskfile syntax as a second supported build language.
- Rewriting unrelated v3 application/runtime APIs.
