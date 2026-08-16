# Wails Enhancement Proposal (WEP)

## Title

**Manifest-Driven Wails v3 Build System**

**WEP Number**: (leave blank, assigned on acceptance)

**Status**: Draft

**Author**: Wails maintainers

**Created**: 2026-08-16

**Discussion**: Local design and decision tickets under `history/wayfinder/wails-build-system/`

**Implementor**: Wails maintainers

**Target**: Wails v3

## Summary

Replace generated project Taskfiles and `build/config.yml` orchestration with
one sparse root `wails.toml` and a built-in Wails-aware Wake pipeline. Wails
owns versioned defaults, planning, caching, generated platform inputs, and
normal build/package/sign/dev execution. Users customize declarative fields or
user-owned script hooks and may eject the full base or one named profile when
they need a source-level frozen snapshot.

## Motivation

The generated v3 Taskfile system was intended to explain the build but became
the practical customization API. Users must understand several generated
files and internal task names, and Wails cannot safely update generated logic
without overwriting user changes. Projects that customize generated files also
stop receiving build-system improvements.

The build implementation should evolve with the pinned Wails CLI while the
project records only identity and deviations from defaults. The normal path
must also be substantially faster than repeatedly interpreting and invoking a
general-purpose task graph.

## Detailed Design

Every new project contains a minimal `wails.toml` with required project
identity. The manifest has no speculative schema field. It declares frontend,
build, dev, target, package, signing, association, protocol, hook, profile, and
extension intent; unknown core fields are errors. Defaults are compiled into
the exact CLI version and package-manager/icon/binding conventions may be
inferred where user values are absent.

`wails3 build`, `package`, `sign`, and `dev` resolve the manifest through a
read-only Planner. The Planner emits one immutable typed Plan for all requested
targets. Project nodes are shared only when their specs are equivalent;
target/package outputs have one owner. The executor schedules the critical
path with bounded CPU and exclusive-tool claims, continues independent
branches after failure, and reports each node through the existing Wake UI.

Direct inputs are content-snapshotted with BLAKE3. Action keys include typed
intent, direct snapshots, consumed artifact digests, tool/handler identity, and
relevant environment. Successful reproducible outputs enter a machine-local
content-addressed artifact store and may be restored when missing. Dependency
installation uses a Receipt rather than storing dependency directories.
Signing and undeclared-side-effect hooks never create reusable cache entries.

Generated platform resources live under ignored `.wails/` state. Stable
customization is expressed through structured manifest fields or optional
user-owned templates. Hooks call one project-relative user-owned script;
cacheable hooks must declare inputs and outputs.

`wails3 eject` materializes and freezes the resolved base. `wails3 eject
<profile>` freezes an existing sparse profile or creates one from the effective
base. Provenance records the complete CLI version. Re-ejection never changes
active values without a safe historical-default comparison. `--backup` creates
a timestamped sibling only when requested.

`wails3 migrate` parses legacy Taskfiles with Wake's AST parser, translates
known identity/configuration/script-hook patterns, records source digests and
stable diagnostics, and marks whether cutover is complete. Unsupported inline
shell remains manual. Legacy files are removed only on explicit request after
complete translation and digest verification. Incomplete migrations continue
using the legacy Taskfile path.

The Dev Session owns persistent frontend/backend processes and watchers. File
bursts request finite production Plans; a newer generation cancels stale work.
The healthy app remains alive until a replacement builds and starts. Manifest
changes reload watch policy and restart only affected persistent processes.

## Non-Goals

- Preserve Taskfile as a second long-term user-facing build language.
- Add arbitrary inline shell or user-defined graph nodes to `wails.toml`.
- Store remote cache artifacts in the first release.
- Treat generated platform files as the primary customization surface.
- Automatically translate arbitrary or inline legacy shell logic.

## Platform Considerations

Windows, macOS, Linux, iOS, and Android use the same manifest and Plan model.
Target overlays supply platform and architecture policy. Platform assets are
generated into `.wails/`, and package/sign adapters continue to call native
toolchains such as NSIS, Xcode/codesign, Android SDK tools, nfpm, and AppImage
tools. Unsupported target/format combinations fail during planning.

Linux build and DEB packaging are exercised end to end by the reference
implementation. Other plans are host-independent and tested structurally but
must pass native-host release verification before rollout.

## Pros/Cons

Pros:

- New projects expose one short, stable file instead of generated build logic.
- Build-system updates ship with Wails without overwriting user files.
- Typed planning and content caching provide fast no-op and incremental builds.
- Ejection supports users who require explicit, frozen configuration.
- Migration is conservative, inspectable, and provenance-aware.

Cons:

- Arbitrary Taskfile customizations require a manifest field, script hook, or
  manual migration.
- Wails owns more cross-platform build implementation and native-tool testing.
- Fully frozen manifests require explicit user review to adopt later defaults.
- The first release has no remote cache and needs native-host acceptance runs.

## Alternatives Considered

Keeping generated Taskfiles preserves flexibility but cannot reconcile user
edits with Wails updates and retains the current verbosity. Generating one
larger Taskfile reduces file count but leaves the same ownership problem.
Treating Taskfile as an override layer creates two supported build languages
and forces users to understand internal task names. A general plugin graph was
rejected in favor of typed first-party nodes plus stable script/template seams.

## Backwards Compatibility

Existing projects without `wails.toml` continue through Taskfile execution.
Migration writes an inactive manifest when unsupported customizations remain;
normal commands cut over only after `wake.migration.complete = true`. Stock
historical templates migrate automatically. Modified sources are never
deleted unless their digest still matches the analyzed source and migration is
complete. Source digests and detailed diagnostics live in the hidden
`.wails/migration-report.json`, not the user-owned manifest. Taskfile is not
retained after an explicit completed cutover.

The public CLI verbs remain, with manifest-native profile, target, format,
force, config, eject, and migration options. Arbitrary Task variables are
rejected on the manifest path with guidance to declarative fields or hooks.

## Security and Privacy

Manifest paths and hook paths must remain inside the project. Migration cannot
delete external includes or changed files. Secrets are not interpolated into
the manifest; signing stores environment-variable or keychain references.
Signing and default hooks are not reusable cache actions. Cache state is local,
disposable, and stores content digests/artifacts rather than credentials.

## Test Plan

- Strict manifest/default/profile/eject/path validation tests.
- Deterministic graph, target overlay, hook barrier, failure isolation, and
  multi-target tests.
- Snapshot, semantic Go API, local replacement, Receipt, artifact restoration,
  corruption, and coarse-timestamp cache tests.
- Stock/modified/external Taskfile migration and digest-guarded removal tests.
- Dev watch, exclusion, process-reconfiguration, and cancellation tests.
- Race detection, vet, and focused command/Wake suites on every change.
- Native build/package/sign acceptance on every supported release host.
- Controlled warm/cold/incremental benchmarks with a sub-100ms badge no-op
  target and explicit regression budget.

## Reference Implementation

The local reference implementation is on branch `codex/wake-cache-mvp`, from
commit `1cce9e558` onward. The detailed design and measured evidence are in
`history/wails-toml-build-system.md` and its local Wayfinder tickets. It has not
been pushed or submitted for external review.

## Maintenance Plan

The manifest package owns defaulting, validation, overlays, and ejection; the
pipeline package owns typed planning/execution; the cache package owns content
identity and storage; command adapters own native tools. Released default
snapshots should be retained when upgrade suggestions are introduced. New
public fields or node behavior require manifest, planner, migration,
documentation, and platform acceptance coverage. Performance regressions are
tracked against controlled benchmark medians.

The legacy Taskfile parser remains only for migration/compatibility during the
announced transition and can be removed after that policy expires.

## Conclusion

A sparse manifest plus a built-in typed Wake pipeline makes the normal Wails
build simple to understand, safe to update, and fast, while retaining explicit
escape hatches and a conservative path from existing v3 projects.
