# Wails Enhancement Proposal (WEP)

## HCL manifest build system for Wails v3

**WEP Number**: (assigned on acceptance)

**Status**: Draft

**Author**: Wails maintainers

**Created**: 2026-08-16

**Target**: Wails v3

## Summary

New Wails v3 projects use one root `wails.hcl` to describe project identity,
frontend commands, build and development policy, targets, packages, signing
references, file associations, protocols, and user-owned platform inputs.

Users request outcomes with `wails3 build`, a complete named build profile, or
bounded anonymous `--targets` and `--formats` options. Wails expands that intent
into an immutable multi-target Plan and executes a fixed Wails-owned pipeline.
Wake provides typed planning, ownership validation, caching, concurrency,
cancellation, and reporting without exposing a task graph as project
configuration.

The presence of `wails.hcl` is the explicit opt-in and cutover flag. Once the
file exists, Wails completely ignores Taskfiles. Migration is explicit,
reviewable, non-destructive, and inactive until the user activates it.

## Motivation

Generated Taskfiles became Wails' practical customization API. Users must
understand generated files and internal task names, while Wails cannot improve
generated orchestration without risking user edits. Production configuration is
also split across Taskfiles, variables, `build/config.yml`, and platform files.

Most projects need light customization, not a programmable build graph. The
project should record its identity and deviations from versioned Wails defaults;
the CLI should own the implementation that evolves with each Wails v3 release.

## Goals

- Centralize production and development intent in one readable file.
- Cover everything required for currently supported production binaries and
  packages on Windows, macOS, Linux, iOS, and Android.
- Preserve user-owned assets and templates byte-for-byte.
- Provide one documented build workflow with deterministic text and JSON Plans.
- Make generated state disposable, final artifacts exact, and cache behavior
  safe and inspectable.
- Keep existing projects on their Taskfile path until they explicitly cut over.

## Non-goals for the first release

- Hooks, arbitrary commands, custom stages, stage replacement, graph
  dependencies, or user-defined pipelines.
- HCL expressions, functions, includes, interpolation, inheritance, or HCL-JSON.
- Automatic translation of arbitrary shell logic or reachable customized
  Taskfile behavior without a typed equivalent.
- Device deployment, emulator management, IDE launch, or a generic `other`
  configuration bucket.
- Remote artifact caching.

Hooks and typed tool invocations may be designed later. This proposal reserves
no current stage name, command shape, or HCL syntax for that work.

## Manifest contract

The native manifest is exactly `wails.hcl`. Wails searches upward from the
working directory; the directory containing the nearest manifest is the project
root.

The first attribute is `version = 3`. It identifies the Wails major contract,
not an exact CLI version. The schema uses shallow semantic blocks, singular
labelled blocks for repeated concepts, and `snake_case` names.

The first release accepts only literal strings and heredocs, booleans, numbers,
tuples, and closed objects. References, interpolation, calls, operators,
conditions, comprehensions, and `null` are errors. Unknown or duplicate fields,
duplicate blocks, duplicate logical labels, and wrong primitive types are errors
with source ranges.

All manifest paths are project-relative. Absolute paths and references into
Wails-owned `.wails/` state are rejected. Normal commands never rewrite the
active manifest.

One typed schema and its default metadata drive decoding, exact type checking,
semantic validation, reference documentation, examples, Plan resolution, and
ejected output.

## Resolution

Wails resolves a request in this order:

1. Discover and parse the nearest manifest.
2. Validate literal shapes and exact primitive types.
3. Seed safe Wails-owned implementation defaults.
4. Apply explicit project, platform, and target values through typed rules.
5. Apply one complete build profile, or allowlisted anonymous CLI options.
6. Validate paths, ownership, formats, credentials, and host/tool capabilities.
7. Create a Plan only when no error remains.

Omitted fields inherit. Explicit collections replace inherited collections; an
explicit empty collection clears them. `tags` is additive: project, target, and
CLI tags are concatenated and deduplicated in order.

Normal builds do not infer package managers, lockfiles, binding shape, asset
paths, or production outputs from mutable project contents. Initialization and
migration may inspect conventional files, but write the selected values into
the manifest.

## Profiles and command line

A profile is a complete named production request, not a sparse overlay. It
selects concrete targets, formats, signing or notarization intent, and bounded
compiler policy.

The documented interface is:

| Command | Meaning |
| --- | --- |
| `wails3 build` | Anonymous host application build using finite build policy. |
| `wails3 build release` | Execute the complete `release` profile. |
| `wails3 build --targets … --formats …` | Bounded anonymous multi-target request; each option accepts one comma-separated list. |
| `wails3 build … --plan` | Validate and print the resolved Plan as text without changing files. |
| `wails3 build … --plan --json` | Emit the same versioned Plan as JSON. |
| `wails3 dev` | Start the long-lived development session. |
| `wails3 dev --plan` | Print the finite startup Plan and exit. |
| `wails3 eject` | Write a complete inactive reference manifest. |
| `wails3 migrate` | Analyze legacy inputs and write or report an inactive draft. |
| `wails3 migrate --activate` | Validate and atomically activate a reviewed, complete draft. |

Canonical native selection options are plural `--targets` and `--formats`.
Anonymous compiler overrides include `--tags`, `--obfuscated`, and
`--garble-args`; they appear in the Plan and action keys. A named profile rejects
artifact-affecting overrides and accepts only operational options such as Plan
inspection, force/cache bypass, and verbosity.

`package` and `sign` are temporary compatibility aliases. New documentation
teaches `build`. There is no `config show`, config export, or graph-navigation
interface: ordinary commands validate and `--plan` is the inspection surface.
Legacy-only tool commands and Taskfile flags remain on the legacy route and are
not part of the native HCL interface.

### Experimental follow-up commands

The experiment branch also prototypes two post-first-release surfaces so their
shape can be evaluated before inclusion in the accepted proposal:

- `wails3 config check [profile]` validates the base manifest and either every
  named profile or one requested profile without executing a build. HCL schema
  diagnostics include the source line, range and an actionable hint.
- `wails3 android devices` lists ADB targets. `wails3 android run [profile]`
  selects a connected device or AVD, builds an ABI-specific development APK,
  installs it, and optionally launches it. `--device` and `--emulator` provide
  deterministic automation; a terminal picker provides the initial interactive
  experience. APK remains unavailable as a production profile format.

These commands are experiments, not accepted first-release scope. Device log
streaming, long-lived lifecycle management, iOS deployment, and full source
ranges for planner or host-capability diagnostics remain follow-up work. This
section does not reserve or approve any hook syntax.

## Fixed pipeline and Plan

Wails owns one fixed pipeline. The planner expands requested final artifacts
into prerequisites, shares target-independent work, checks cycles and output
ownership, and schedules independent work concurrently. Deterministic work may
be restored from a content-addressed cache. Signing, notarization, credentials,
and externally stateful operations are never reusable cache entries.

Stable reporting stages are `resolve`, `prepare`, `generate`, `frontend`,
`compile`, `assemble`, `package`, `sign`, and `collect`. Stages are reporting
categories, not execution APIs. Users cannot select, replace, reorder, or attach
hooks to them.

The default Plan is a compact text table. `--plan --json` emits the same Plan as
a versioned machine-readable document. Plans include selected profile and
targets, nodes and dependencies, inputs, outputs, cache decisions, provenance,
and exact final artifacts. Plan inspection performs no writes.

## Assets, generated state, and artifacts

Ownership has three zones:

| Zone | Contract |
| --- | --- |
| User-owned inputs | `wails.hcl`, source, referenced assets/templates, existing `build/` content, and otherwise unclassified files. Wails may read, hash, and copy them, but never edit or delete them. |
| Generated inputs | Disposable Wails-owned workspaces, indexes, reports, and platform/package inputs beneath `.wails/`. |
| Final artifacts | Exact Plan-declared binaries and packages under `build.output`. Wails may atomically replace only those exact paths. |

An exact final-artifact path is the sole exception to otherwise user-owned,
unclassified content: once selected by a validated Plan, final-artifact
ownership takes precedence and Wails may replace a pre-existing file or
directory at that exact path. Siblings and ancestors remain user-owned. The
Plan must reject two Nodes claiming the same final path.

A custom file or directory completely replaces its corresponding generated
input. Wails copies it byte-for-byte into `.wails/`, preserving executable
permissions, with no interpolation, merge, sanitization, or line-ending change.
Structured settings that conflict with a complete replacement are errors.

`wails3 clean` may always remove `.wails/`. It removes a recorded final artifact
only when the current digest still matches the Wails-produced digest. Modified
and unknown files are preserved and reported.

## Targets, packaging, and signing

Supported targets and formats need no empty enablement blocks. Platform and
package blocks exist only for customization; profiles or anonymous flags select
outputs.

| Platform | Targets | Production formats |
| --- | --- | --- |
| Windows | `amd64`, `arm64` | binary, NSIS, MSIX |
| macOS | `amd64`, `arm64`, synthetic `universal` | app, DMG |
| Linux | `amd64`, `arm64` | binary, AppImage, DEB, RPM, Arch Linux |
| iOS | `arm64` with simulator/device destination | app, IPA |
| Android | `amd64`, `arm64`, synthetic `universal` | AAB only |

Platform blocks expose native identity, assets, SDK policy, and named signing
references. Target blocks are limited to architecture-sensitive compiler,
environment, minimum-version, and bounded toolchain policy. Toolchain strategies
are `auto`, `native`, `zig`, and `docker`.

Android Dev may assemble an APK internally for installation on an emulator or
device. APK is not a selectable production format in `wails.hcl`; profiles and
production `wails3 build` requests accept AAB only.

Package blocks are closed typed schemas. Credentials, passwords, and private
keys never appear in HCL; the manifest contains named references. Android
production never silently uses a debug keystore. A signed AAB requires a valid
credential, while unsigned release output must be deliberate.

File associations and protocols use explicit labelled blocks. Optional platform
filters default only to platforms where Wails supports the feature. Unsupported
combinations fail validation.

## Development

`wails3 dev` owns watchers and persistent frontend/application processes. It
requests a finite Plan for startup and each rebuild while retaining the same
ownership, cache, cancellation, and reporting rules as production.

The `dev` block contains log level, debounce, watch and exclusion patterns,
Git-ignore handling, shutdown grace period, and development tags. Legacy
`root_path` is unnecessary because the manifest directory is the project root.
Legacy `executes` is not migrated because arbitrary process orchestration is not
configuration.

Finite production policy such as `build.output`, `trim_path`, `strip`, and
obfuscation does not leak into development. Dev retains debug information, uses
the frontend dev command, and stores transient binaries beneath `.wails/dev/`.

Manifest reload is transactional: a failed replacement watcher, frontend,
backend, or build leaves the last healthy generation running. A newer file-event
generation cancels stale finite work; cancellation is not reported as a failed
build.

## Migration and cutover

Routing is determined only by project state:

| State | Behavior |
| --- | --- |
| No `wails.hcl` | Commands stay on the legacy Taskfile compatibility path. |
| `wails.migrated.hcl` | Inactive migration draft; normal commands ignore it. |
| `wails.hcl` | Native cutover; Taskfiles are completely ignored and invalid HCL never falls back. |

`wails3 migrate` analyzes root and included Taskfiles, conventional overrides,
`build/config.yml`, project metadata, lockfiles, TypeScript configuration, and
conventional assets. It extracts only behavior representable by typed HCL and
writes a fresh inactive draft. It never overwrites a draft, modifies an asset,
or deletes a Taskfile.

Reachable custom behavior that cannot be proved equivalent is a blocker. Stock
generated orchestration is replaced by the fixed pipeline. Utility tasks outside
Wails build/package/sign/dev reachability are reported but do not block.

`wails3 migrate --activate` reruns analysis, validates the reviewed draft,
refuses unresolved blockers, and atomically renames it to `wails.hcl`.
Activation changes no other project file. There is no automatic backup,
deletion, rollback, or stale-Taskfile fallback.

## Ejection

`wails3 eject` writes one complete canonical inactive sibling,
`wails.ejected.hcl`. It includes current Wails defaults and records the exact
generating CLI version in a comment. It excludes invocation overrides, host
facts, credentials, and tool versions, and never overwrites without explicit
force. Ejection does not create another way to build; the file is inert until a
user deliberately makes it the active `wails.hcl`.

## Security and correctness

- Manifest paths are contained within the project and cannot target `.wails/`.
- Exact output ownership is validated before execution.
- Action keys include typed intent, direct input snapshots, dependency artifact
  digests, tool/handler identity, and relevant explicit environment.
- Cache restore verifies content and publishes atomically; corrupt entries are
  discarded as misses.
- Secrets are referenced by name, redacted, and never serialized into Plans or
  receipts.
- Signing and other non-reproducible operations are not cached.
- Independent Plan branches continue after unrelated failures while descendants
  of a failed node are blocked.

## Examples

Minimal project:

```hcl
version = 3

project {
  name         = "hello"
  product_name = "Hello"
  identifier   = "com.example.hello"
  version      = "0.1.0"
  binary_name  = "hello"
  icon         = "assets/appicon.png"
}

frontend {
  directory = "frontend"
  install   = ["npm", "install"]
  build     = ["npm", "run", "build"]
  dev       = ["npm", "run", "dev"]
  output    = "frontend/dist"
}

build {
  output = "bin"
}
```

Light customization with one complete release profile:

```hcl
build {
  output = "dist"
  tags   = ["sqlite_fts5"]
}

dev {
  debounce_ms    = 500
  watch          = ["**/*.go", "frontend/src/**"]
  exclude        = ["dist/**", "frontend/node_modules/**"]
  use_git_ignore = true
}

target "windows/amd64" {
  toolchain = "zig"
  tags      = ["enterprise"]
}

package "nsis" {
  install_scope = "user"
}

profile "release" {
  target "windows/amd64" {
    formats = ["nsis"]
  }
}
```

User-owned inputs:

```hcl
windows {
  icon     = "assets/windows/app.ico"
  manifest = "assets/windows/app.manifest"
}

darwin {
  icon       = "assets/macos/app.icns"
  assets_car = "assets/macos/Assets.car"
  info_plist = "assets/macos/Info.plist"

  signing {
    entitlements = "assets/macos/entitlements.plist"
  }
}

package "nsis" {
  template = "packaging/windows/installer.nsi"
}

file_association "studio_project" {
  extensions  = ["studio"]
  name        = "Studio Project"
  description = "Studio project file"
  icon        = "assets/filetypes/studio-project.png"
  role        = "editor"
  platforms   = ["windows", "darwin", "linux"]
}

protocol "studio" {
  description = "Open Studio links"
  platforms   = ["windows", "darwin", "linux", "ios", "android"]
}
```

## Deferred extension work

A later proposal may define hooks, typed Wails tool invocations such as
`wails.generate.icons`, or other pipeline customization. That work must specify
lifecycle semantics, declared inputs and outputs, environment, caching,
cancellation, failures, and external process containment. It must not infer an
API from current stage names or legacy `wails3 tool` commands.
