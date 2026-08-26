# Wayfinder map: Wails manifest build system

Status: active
Type: map

## Destination

An implemented first release for replacing generated v3 Taskfiles and
`build/config.yml` with one config-first `wails.hcl`, using a fixed Wake pipeline
for graph execution, caching, parallelism, and reporting. Pipeline extension
and typed `wails3 tool` calls were deliberately deferred from that first slice.
The bounded lifecycle-hook follow-up is now implemented experimentally;
follow-up work still covers
device/emulator deployment with an interactive experience, and stronger
line-referenced configuration diagnostics.

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
  — six project/target/package barrier phases invoke project-owned scripts
  directly, default to non-cacheable execution, and support bounded explicit
  input/output caching without exposing arbitrary graph edges.
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
- The checked local audit measured a 99.092ms badge median with 5.05% MAD and
  all six Nodes cached. A post-rebase profile raised filesystem traversal
  concurrency from 8 to 16, cutting a 21-sample median from 111.284ms to
  103.367ms; every sample still had zero executed and six cached Nodes. Dirty
  published binaries are regenerated byte-identically;
  missing binaries are restored byte-identically without executing a handler.
- The 2026-08-23 handoff audit removed hidden npm dependencies from ten
  command-level Plan/Dev tests and fixed artifact-stability checks against
  historical baselines without artifact metadata. Normal and race-tested
  command/Wake suites pass. The final local combined performance gate measured
  103.913ms with 1.66% MAD, zero executed and six cached Nodes, and stable
  bytes. After the ceiling was revised to 150ms, the final post-Podman/mobile
  seven-sample run passed at 99.870ms with 4.43% MAD and stable bytes.
- After adding bounded lifecycle hooks, the unchanged no-hook graph passed the
  same seven-sample gate at 105.388ms median with 3.45% MAD, zero executed and
  six cached Nodes in every sample, and byte-identical artifacts.
- Native Linux build, DEB/RPM/Arch/AppImage packaging, all nine built-in web
  templates, and the real Dev lifecycle pass locally. Podman-backed arm64
  cross-build and DEB/RPM/Arch packaging also pass with an AArch64 ELF and
  zero-work reruns. Real Android API 36/NDK r29 amd64, arm64 and universal AABs
  also pass with receipts and zero-work reruns. `wails3 android run` starts an
  API 36 AVD, builds its preferred x86_64 APK, installs it and resumes the Wails
  activity. Attached runs stream PID-filtered application logs, follow process
  restarts, detect real emulator loss, and support explicit post-session cleanup;
  failed startup is cleaned up automatically. Generated projects use the API
  36-compatible AGP 9.0.1/Gradle 9.2.1 pairing and avoid deprecated Groovy
  property assignment syntax.
  Physical-device, foreign native packaging, signing and native-arm64
  launch acceptance remain matching-host release work.
- The full Go suite reaches unrelated desktop-environment failures only;
  focused manifest/Wake/command tests, race detection, vet, and real Linux
  build/package runs pass.
- A draft WEP records the public behavior and rollout proposal under
  `v3/wep/proposals/manifest-build-system/` on the experiment branch.
- The experimental hook slice adds strict HCL schema and path validation,
  typed `RunHook` Plan nodes, versioned JSON context files, package/sign
  barriers, process-group cancellation, opt-in artifact restoration, and
  conservative Taskfile lifecycle-script migration.
- Configuration diagnostics retain explicit HCL block and attribute origins
  through semantic validation, planning, signing and host-capability checks.
  The shared Build/Dev/config-check renderer provides source carets and hints,
  redacts environment values, has golden and multiple-error coverage, and
  benchmarks at approximately 1.3 microseconds without source-file I/O.

## Not yet specified

The config-first destination, bounded hook contract, Android deployment
interaction model and configuration-diagnostic surface are implemented.
Physical Android-device and native Apple deployment acceptance remain follow-up
work under `issues/`.

## Out of scope

- Designing a general-purpose replacement for Task.
- Preserving arbitrary Taskfile syntax as a second supported build language.
- Rewriting unrelated v3 application/runtime APIs.
