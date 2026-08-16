# Wails v3 Manifest Build System

Status: MVP implemented for backwards review

## Summary

Wails v3 should move from generated Taskfiles to a built-in, Wails-aware build
pipeline configured by one root-level `wails.toml` file. New projects emit only
the required project identity; build behavior stays implicit until customized
or ejected.

The normal project should not contain a generated build graph. Wails owns the
default graph, its platform behavior, its cache, and the temporary files needed
by platform packagers. The project declares intent and only records deviations
from the defaults.

The existing Taskfile system is migrated into this model. It is not retained as
the long-term user-facing build API.

## Implementation snapshot (2026-08-16)

The MVP now exercises the proposal as one product surface rather than as an
isolated cache prototype:

- `internal/wake/manifest` owns sparse loading, compiled defaults, strict
  validation, named profiles, and full/profile ejection;
  `internal/wake/migration` owns private report state and migration cutover.
- `internal/wake/pipeline` resolves typed build/package/sign DAGs and executes
  them through a bounded critical-path scheduler with receipts and restorable
  content-addressed artifacts.
- One invocation may request comma-separated Targets. Equivalent project Nodes
  are shared, Target outputs are isolated, and incompatible project-level
  projections fail during planning instead of racing over an output.
- `wails3 build`, `package`, `sign`, and `dev` select the manifest pipeline for
  native or completed-migration projects; incomplete migrations retain the
  legacy Taskfile path.
- Platform inputs are generated under `.wails/`; desktop, iOS, Android,
  packaging, signing, hooks, target overlays, and Darwin universal builds all
  use typed specs rather than Task variables.
- `wails3 migrate` recognizes shipped historical Taskfiles, translates project
  metadata and safe script-file hook tasks, records source digests and
  diagnostics, and retires unchanged sources after automatic or explicitly
  reviewed completion; `--backup` preserves the legacy tree when requested.
- New templates contain a minimal `wails.toml` and no generated Taskfile or
  build directory. Older community templates are analysed and remain on legacy
  execution when their customizations are not completely represented.
- Compile cache identity includes local module/workspace sources and relevant
  execution environment. Receipts are state evidence rather than Artifacts,
  and signing-bearing iOS assembly never enters the reusable Artifact Store.
- The Dev Session cancels stale generations without presenting them as failed
  builds, preserves the healthy app through failures, skips replacement when
  the binary is unchanged, reloads watch policy without rebuilding, and stages
  frontend/port changes before committing the new session.
- The repository-required draft WEP is recorded at
  `v3/wep/proposals/manifest-build-system/proposal.md` for review without any
  GitHub activity.

Measured on the disposable badge fixture after the production integration:

| Scenario | Result |
|---|---:|
| cold desktop build with warm language/tool caches | 2.88 s |
| warm no-op desktop build | 70–80 ms wall; 4/4 cached |
| service method-body edit | 0.91 s; only Go compile ran |
| first DEB package after CLI rebuild | 2.76–2.84 s |
| warm no-op DEB package | 70–90 ms wall; 6/6 cached |

Linux build and DEB packaging have been exercised end to end. Windows, macOS,
iOS, Android, and their native signing tools are represented in deterministic
plans and covered by host-independent tests, but still require native-host
release verification before this moves beyond MVP status.

The final repeatability run also verified that packaging over an existing DEB
uses an isolated staging directory and replaces the expected artifact instead
of misidentifying the previous output as another packager candidate.

## Goals

- Make a new project's build system understandable from one file.
- Keep the default configuration close to zero lines.
- Make the exact default set observable, with reproducibility provided by a
  pinned CLI or an ejected snapshot.
- Allow advanced users to eject and fully own a default profile.
- Generate platform configuration at build time instead of committing it.
- Reuse Wake's graph, cache, parallel executor, and structured reporting.
- Support common customizations without requiring users to understand the
  internal pipeline.
- Provide a safe migration path from default and modified v3 Taskfiles.

## Non-goals

- Preserve arbitrary Taskfile syntax as a second build language.
- Support inline shell as a general-purpose extension mechanism.
- Require users to learn internal task names such as
  `common:generate:bindings`.
- Make generated platform files the primary customization surface.

## Project layout

A default project should look like:

```text
my-app/
  wails.toml              # required identity; build defaults remain implicit
  main.go
  go.mod
  frontend/
  assets/
  scripts/                # optional user-owned build hooks
```

The `build/` directory is no longer generated for orchestration. Generated
plists, manifests, NSIS files, nfpm files, Xcode projects, and similar files
are written to an ignored `.wails/` working directory or passed through memory
where the toolchain permits it. They are build products, not project source.

User-owned resources remain ordinary project files and are referenced from the
manifest, for example `assets/appicon.png` or
`packaging/windows/installer.nsi`.

## Manifest principles

The manifest is declarative. It describes the application, selected profile,
frontend, targets, packaging preferences, and extensions. It does not describe
the complete implementation graph.

The initial parser has no `schema` field and rejects it like any other unknown
top-level field. Explicit schema versioning should be introduced only if a
future incompatible format actually requires it.

The built-in defaults are keyed by the complete version of the running Wails
CLI. An ordinary project follows the defaults compiled into that CLI. An
ejected scope records the same complete version in `ejected_by`, freezes the
resolved values, and does not inherit later default changes.

Illustrative configuration:

```toml
[project]
name = "badge"
product_name = "Badge"
company_name = "My Company"
identifier = "com.mycompany.badge"
version = "0.1.0"
description = "A Wails application"
icon = "assets/appicon.png"

[frontend]
directory = "frontend"
package_manager = "auto"
build_command = "build"
dev_command = "dev"

[dev]
port = 9245

[targets.darwin]
minimum_version = "12.0"

[targets.darwin.arm64]

[targets.darwin.amd64]

[package.linux]
formats = ["appimage", "deb", "rpm"]

[package.windows]
formats = ["nsis", "msix"]

[package.darwin]
formats = ["app", "dmg"]
```

The exact field names are subject to schema review. The important distinction
is between user concepts (`frontend.package_manager`) and implementation
details (the individual dependency, binding, and compile tasks).

## Profiles

Profiles are available but are not emitted by default. The implicit default
profile covers normal development, production builds, packaging, and the
supported targets.

Users materialize the effective top-level base with:

```bash
wails3 eject
```

This expands normal top-level configuration in `wails.toml`. The result is
fully frozen:

- It contains every resolved value, including values that were previously
  implicit.
- It has no inheritance relationship with the built-in default.
- Future Wails releases do not silently modify it.
- Re-running ejection leaves active values unchanged and may add commented
  upgrade suggestions; adopting them is a manual edit.

Named profiles may be added later:

```toml
[profiles.debug.build]
production = false

[profiles.release.build]
production = true
```

Profile selection should be exposed by the primary commands, for example
`wails3 build --profile release`, but profiles should not be required to use
Wails.

## Customization

Customization has three layers, in order of preference.

### Structured configuration

Common needs are first-class manifest fields: package manager, output path,
Go tags, architectures, signing options, package formats, metadata, file
associations, protocols, and development settings.

### Script hooks

Shell is intentionally limited to invoking a user-owned script:

```toml
[hooks]
before_build = "scripts/generate-version.sh"
after_package = "scripts/publish-artifact.sh"
```

A cache-aware hook uses the long form for the same phase:

```toml
[hooks.before_build]
script = "scripts/generate-version.sh"
cache = true
inputs = ["version.txt"]
outputs = ["generated/version.go"]
```

Wails provides a stable environment to hooks:

```text
WAILS_PROJECT_DIR
WAILS_TARGET_OS
WAILS_TARGET_ARCH
WAILS_PROFILE
WAILS_OUTPUT
WAILS_PIPELINE_VERSION
```

Hook scope and `WAILS_OUTPUT` are fixed by phase:

| Phase | Invocation scope | `WAILS_OUTPUT` |
| --- | --- | --- |
| `before_build` | once for the shared Project | empty; target variables are also empty |
| `after_build` | once per Target | compiled binary |
| `before_package` | once per Target before its requested formats | future package path, or their common parent |
| `after_package` | once after the Target's requested formats | package path, or their common parent |
| `before_sign` | once before the Target's requested signing work | unsigned package path, or their common parent |
| `after_sign` | once after the Target's requested signing work | signed package path, or their common parent |

Package and signing hooks deliberately form barriers around the requested
format set. The package Nodes inside those barriers remain independent and may
run concurrently.

Hooks without `cache = true` are side-effectful and run every time. A cached
hook must declare complete inputs and outputs; declaring either without the
explicit cache opt-in is an error. The referenced script and relevant
executable metadata are always cache inputs, so users do not repeat the script
path in `inputs`.

Unix hooks are executable files with their own shebang. Windows `.cmd`, `.bat`,
and `.ps1` hooks are invoked through the corresponding platform interpreter,
with the script path passed as one argument rather than interpolated shell.
For cached hooks, multiple outputs must share one non-root directory; that
directory is the bounded Artifact Wails records and restores. The output root
cannot contain the script or any declared input. Scripts and working
directories are resolved through symlinks and must remain inside the Project.
Cancellation terminates the hook's process group so a stale build cannot leave
detached children.

### User-owned templates and extensions

Structured configuration should cover the common case. When it cannot, the
manifest can point to a stable user-owned template:

```toml
[package.windows.nsis]
template = "packaging/windows/installer.nsi"
```

Wails must never overwrite referenced user-owned files. Extension namespaces
can provide configuration for tooling that Wails does not understand without
exposing the full internal pipeline.

Package templates use strict Go `text/template` rendering against a versioned
model with `.Project`, `.Target`, `.Package`, `.Paths`, `.Associations`,
`.Protocols`, and `.Options`. A single source file renders to the input owned by
that package format. A directory source is copied atomically; `*.tmpl` files
are rendered with the suffix removed, ordinary files retain their mode, and
symlinks are rejected. The current destinations are:

| Format | User-owned template output |
| --- | --- |
| NSIS | `project.nsi` |
| MSIX | `AppxManifest.xml` |
| macOS app / iOS app or IPA | `Info.plist` |
| DMG | JSON package options |
| AppImage | desktop entry |
| DEB / RPM / Arch Linux | nfpm YAML |
| APK / AAB | complete Android Gradle directory |

The source remains beside the project while Wails renders into a disposable
`.wails/package/<format>/<target>/` workspace or the final package Artifact.
Package Nodes never modify the generated platform-assets Artifact they consume.

## Built-in pipeline

The built-in Pipeline is resolved by a read-only Planner into one immutable,
multi-Target Plan of typed Nodes rather than Taskfile tasks:

```text
install frontend dependencies ─┐
generate bindings ─────────────┼─→ build frontend ─┐
generate target assets ────────┘                  ├─→ compile / bundle
                                                  └─→ package / sign per format
```

Project inspection, toolchain checks, Target expansion, and input discovery
happen during planning rather than as dynamic Nodes. The actual graph is
Target- and Profile-dependent. Wails derives inputs from:

- `go.mod`, `go.sum`, and the Go package graph;
- Go source and embedded files;
- frontend source, `package.json`, and the selected lockfile;
- generated bindings and frontend output;
- manifest values and target environment;
- declared hook inputs and outputs.

Wake executes this graph through a single non-recursive, resource-aware worker
scheduler. Structural Node Keys deduplicate shared work; separate cache
fingerprints cover resolved configuration, tools, environment, source inputs,
and dependency results. Ready work is prioritized by estimated critical-path
duration. Wails operations run in-process where possible and unavoidable tools
are launched directly with typed arguments rather than shell strings.

The reporting contract, Pulse renderer, output capture, process cancellation,
and generalized DAG algorithms are reusable. Taskfile parsing, templates,
variables, namespaces, overrides, fallback, recursive execution, and the mtime
task cache remain isolated in the legacy migration/compatibility path.

### Cache and automatic input discovery

Each typed Node handler discovers the narrow input set implied by its resolved
Spec. Built-in Nodes do not require users to maintain source globs. Wake forms
an Action Key from the handler version, effective typed config, relevant
Target, recognized environment and tool identities, direct input Snapshots,
and consumed upstream Artifact digests. BLAKE3 and canonical typed encodings
provide fast, portable content identity.

Snapshots use a persistent filesystem fast path: unchanged file identity,
size, modification time, and change stamp reuse the previous content digest.
Directory inventories detect added and removed inputs, while Dev Session
watchers invalidate touched paths in memory. A full-verification path remains
available for unusual metadata-preserving workflows.

Go discovery caches the package inventory per Target, tags, module/workspace
state, and toolchain. It tracks selected source, Cgo, assembly, embed, and local
replacement files without scanning the module cache. Binding analysis produces
a semantic Binding Model, and downstream frontend work depends on generated
binding content, so implementation-only Go changes do not cascade when the
generated files remain identical.

Frontend installation uses a Receipt over package/workspace manifests,
lockfile, package-manager config, and tool identity; dependency directories are
not scanned or archived. Frontend builds exclude dependencies, previous output,
VCS data, and Wake state while including generated bindings and recognized
build configuration.

The project keeps a local Action Index backed by a machine-global
content-addressed Artifact Store. Missing outputs may be restored from it;
modified generated outputs are dirty misses and are atomically replaced by a
successful rebuild. Signing, notarization, publication, credential access, and
network side effects always run. Hooks run by default and become cacheable only
with an explicit opt-in and complete declared inputs and outputs.

### Development lifecycle

`wails3 dev` owns a long-lived Dev Session rather than representing watchers or
persistent processes as graph Nodes. A coalesced file burst requests the normal
finite development Plan. A newer generation cancels stale work and removes its
live report without rendering a failed-build verdict. A compile cache hit or
restoration leaves the current backend running; an executed `after_build` hook
or changed binary stages a replacement and swaps only after startup readiness.
A failed Plan or replacement leaves the healthy app intact.

The frontend server and backend receive explicit per-process
`FRONTEND_DEVSERVER_URL` and `WAILS_VITE_PORT` values; the CLI does not leak
candidate configuration through process-global environment. Both sides use
the same pinned IPv4 loopback address, avoiding `localhost` IPv4/IPv6
resolution disagreement. Port changes stage the new frontend first, then the
backend that consumes it, before terminating the old pair.

Watch sets are replaced transactionally. Root and nested `.gitignore` changes
reload policy immediately without a build, newly created directories are
registered and scanned for already-created inputs, and a failed replacement
retains the old watches. Shutdown cancels and waits for in-flight Plans, then
terminates and reaps the backend, frontend, and their complete process groups.

## Generated platform assets

Platform assets are generated at build time from the single manifest:

- macOS `Info.plist`, entitlements, app bundles, and Xcode support;
- Windows manifests, NSIS inputs, MSIX metadata, and `.syso` resources;
- Linux desktop metadata, AppImage inputs, and nfpm configuration;
- iOS and Android generated project/configuration files.

There is no general-purpose `update build-assets` operation in the new model.
Changing `wails.toml` changes the next generated build inputs. User-owned
templates and resources are explicit exceptions and are never regenerated.

## Migration

The migration command is:

```bash
wails3 migrate
wails3 migrate --dry-run
wails3 migrate --complete
wails3 migrate --complete --backup
```

Migration discovers the root Taskfile, included common/platform Taskfiles, and
`build/config.yml`. It parses them through the Task/Wake AST, normalizes the
result, and compares it with the known generated Wails templates.

Migration tiers:

1. **Unmodified generated projects** — migrate automatically and completely.
2. **Modified generated projects** — translate recognizable differences into
   manifest fields, profiles, hooks, or user-owned template references.
3. **Arbitrary Taskfile logic** — report the unsupported portions and require
   manual conversion to structured config or scripts.

Inline shell blocks are deliberately not converted into generated scripts.
They produce a migration warning and require the user to create a script. This
avoids creating opaque generated files while keeping the new system's shell
boundary clear.

Migration writes `wails.toml` plus a private machine-readable report and emits
a human-readable summary. Automatically complete migrations digest-check and
retire represented sources. Incomplete migrations retain Taskfiles; after
manually representing reported logic in config or scripts, the user confirms
cutover with `--complete`, optionally preserving the source tree with
`--backup`. The machine-readable completion state,
behavior when `wails.toml` already exists, and provenance requirements for
removing only fully represented legacy files are resolved by the
[Taskfile migration decision](wayfinder/wails-build-system/issues/06-taskfile-migration.md).
Command precedence while Taskfiles and a Manifest coexist is resolved by the
[CLI cutover decision](wayfinder/wails-build-system/issues/07-cli-compatibility.md).

## CLI surface

Initial commands:

```bash
wails3 build
wails3 build --profile release
wails3 package
wails3 sign
wails3 dev
wails3 eject
wails3 migrate
wails3 config check
wails3 config show
```

`wails3 task` is not part of the new public build API. Existing Taskfile-based
projects continue to work during the migration period, but newly initialized
projects do not receive Taskfiles.

## Implementation phases

### Phase 1: manifest and defaults

- Define typed manifest structs and validation.
- Implement implicit defaults and profile resolution.
- Add `wails3 config check` and `wails3 config show`.
- Add `wails3 eject` with full materialization and no inheritance.

### Phase 2: built-in graph

- Extract Wake's DAG, cache, execution, and reporting layers from Taskfile
  assumptions.
- Implement project inspection and frontend/Go input discovery.
- Implement the default build, dev, package, and sign nodes.

### Phase 3: generated platform assets

- Replace project-generated Taskfile/config templates with runtime generation.
- Move packaging inputs to typed manifest structures.
- Add explicit user-owned template references.
- Keep generated material under `.wails/` or in memory.

### Phase 4: migration

- Add Taskfile/config discovery and AST normalization.
- Recognize the shipped v3 templates.
- Translate common modifications.
- Produce warnings and dry-run reports for unsupported behavior.

### Phase 5: new project templates and documentation

- Stop generating Taskfiles from `wails init`.
- Update examples and documentation.
- Retain legacy Taskfile execution only for projects that have not migrated.

## Open design work

- Exact pipeline/default version compatibility policy.
- Structured signing configuration for each platform.
- Which existing packaging customizations deserve first-class fields.
- Migration mappings for the most common modified Taskfiles.
