# Wayfinder map: Wails manifest build system

Status: resolved
Type: map

## Destination

An implemented first release for replacing generated v3 Taskfiles and
`build/config.yml` with one config-only `wails.hcl`, using a fixed Wake pipeline
for graph execution, caching, parallelism, and reporting. Pipeline extension,
hooks, and typed `wails3 tool` calls are deliberately deferred.

## Notes

- Domain: Wails v3 build tooling and migration.
- The issue files in this directory preserve the design exploration. References
  there to TOML or first-release hooks are superseded by the accepted WEP.
- Standing decisions from the conversation: one `wails.hcl`; version `3` stays
  in lockstep with CLI releases; profiles select targets, formats, signing, and
  destinations; `wails3 eject` writes `wails.ejected.hcl`; the presence of
  `wails.hcl` opts in and makes Taskfiles completely ignored; migration never
  silently translates arbitrary commands or removes Taskfiles.
- The narrow cache slice has now been expanded into a complete implementation
  surface for hands-on backwards review. Matching-host release verification
  remains a release operation governed by the acceptance contract, not an
  unresolved build-system decision.
- The authoritative public proposal is
  [`v3/wep/proposals/manifest-build-system/proposal.md`](../../../v3/wep/proposals/manifest-build-system/proposal.md).

## Decisions so far

- [Manifest vocabulary and fully frozen eject semantics](issues/01-manifest-and-eject.md)
  — `wails.hcl` stays minimal until customized, then resolves through typed
  Platform/Target/Profile overlays; `wails3 eject [profile]` creates an
  independently frozen, host-neutral `wails.ejected.hcl` snapshot.
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
  Taskfiles; activation revalidates the reviewed draft and leaves all legacy
  sources untouched.
- [CLI migration cutover](issues/07-cli-compatibility.md)
  — one routing seam selects native Manifest, deliberate legacy fallback, or
  an actionable ambiguity error; `wails.hcl` itself is the opt-in flag.
- [Script hook contract and safe extension boundary](issues/05-script-hook-contract.md)
  — preserved as extension research only. Hooks and arbitrary calls are not in
  the config-only first release and require a later pipeline proposal.
- [Dev Session invalidation and process lifecycle](issues/10-dev-session-lifecycle.md)
  — finite generations drive a transactional, readiness-gated frontend,
  backend, and watch lifecycle with no-op preservation and clean cancellation.
- [Build performance acceptance matrix and regression budget](issues/09-performance-acceptance.md)
  — permanent semantic, variance, latency, concurrency, native-host, mobile,
  and signing gates make speed and platform correctness verifiable.

## Implementation evidence

- Minimal config, profile ejection, migration/cutover, generated assets, target
  overlays, build/package/sign plans, and the manifest dev lifecycle are
  implemented behind the normal CLI commands.
- Follow-up two-axis review closed cache identity for local Go sources and
  relevant environment, separated Receipts from Artifacts, made signing
  non-reusable, added missing-profile/re-eject semantics, combined multiple
  Targets in one Plan, rejected deferred script migration, and made Dev rebuilds
  generation-cancelable with process/watch reconfiguration.
- The final local audit measured a 99.092ms badge median with 5.05% MAD and all
  six Nodes cached. Dirty published binaries are regenerated byte-identically;
  missing binaries are restored byte-identically without executing a handler.
- Native Linux build, DEB/RPM/Arch/AppImage packaging, all nine built-in web
  templates, and the real Dev lifecycle pass locally. Foreign native packaging,
  signing, and SDK-backed mobile acceptance remain matching-host release work.
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
