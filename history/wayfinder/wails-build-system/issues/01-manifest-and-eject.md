# Manifest vocabulary and fully frozen eject semantics

Type: grilling
Status: resolved
Blocked by: none

## Question

What is the smallest coherent `wails.toml` vocabulary for project metadata,
frontend settings, targets, packaging, profiles, hooks, and user-owned
templates? What exactly does `wails3 eject` materialize, and how are implicit
defaults versioned without forcing a schema or profile block into every new
project?

## Comments

### 2026-08-14 — Grilling round 1

- Canonical terms are Manifest, Pipeline, Profile, and Target; user-facing
  configuration does not expose tasks or graph nodes.
- New projects contain a minimal manifest with stable project identity; build
  behavior remains implicit until customized or ejected.
- User-authored profiles are sparse overlays. Ejection materializes a complete,
  fully frozen profile.
- Unknown fields are errors except inside explicit `[extensions.<name>]`
  namespaces.

### 2026-08-14 — Grilling round 2

- Every new project writes the required identity configuration in
  `wails.toml`.
- Ejecting the effective base expands normal top-level configuration rather
  than creating an artificial `[profiles.default]`.
- Ejection can select a specific named profile. The target breadth within that
  selected profile remains open.
- Ejected snapshots record `ejected_by`; `ejected_from` is unnecessary.
- Sparse named profiles overlay the project's effective top-level base. Once
  that base is ejected, profiles overlay its frozen values rather than newer
  built-in defaults.

### 2026-08-14 — Grilling round 3, partial

- Profile selection is positional: `wails3 eject` ejects the effective base;
  `wails3 eject release` ejects the named `release` profile. A `--profile`
  flag and a `profile eject` subcommand are unnecessary.
- Ejecting a profile freezes every target reachable by that profile.
- Ejection metadata belongs under a `[wake]` namespace; its exact shape for
  independently ejected profiles remains open.
- The required project identity remains `name`, `product_name`, `identifier`,
  and `version`. The derived `binary_name` is shown as a commented discovery
  aid.
- Profiles cannot override project identity.
- Re-ejection may write candidate changes as comments, but the intended
  behavior needs clarification because commented TOML cannot provide a full
  freeze.

### 2026-08-14 — Grilling round 4

- Initial ejection writes the complete frozen configuration as active TOML.
  Re-ejection preserves that active snapshot and may add newer candidate
  defaults as commented upgrade suggestions rather than silently changing the
  build.
- Base ejection provenance is recorded as `ejected_by` under `[wake]`.
  Independently ejected named profiles are recorded by version under
  `[wake.ejected_profiles]`.
- Profiles use the same nested configuration shape as the top-level manifest.
  They are sparse overrides before ejection and complete snapshots afterward;
  there is no `extends` mechanism.
- Freshly ejected manifests remain clean. Per-value provenance comments are
  reserved for candidate changes proposed by a later re-ejection.

### 2026-08-15 — Grilling round 5, partial

- The schema may support the complete root vocabulary, but a new project's
  manifest emits only the minimal required `[project]` identity. Other
  sections remain implicit until customized or ejected.
- Shared application metadata is declared once and translated into each
  platform's generated assets. File associations and protocols remain
  platform-neutral manifest concepts with platform-specific overrides where
  needed.
- Target configuration should support platform-wide common values such as
  `[targets.darwin]` and a nested architecture-specific layer; the exact
  canonical name for the nested layer remains open.
- Build, packaging, and signing are distinct configuration areas.
- Frontend package manager, lockfile, conventional scripts, and output are
  inferred. `[frontend]` contains overrides, and ejection materializes the
  resolved values.
- User-owned templates are configured beside the feature that consumes them.
- Hooks are keyed by phase under `[hooks]`, using a concise form such as
  `before_build = "scripts/generate-version.sh"`, rather than an ordered array
  of records. Multiplicity and cached hook declarations remain open.

### 2026-08-15 — Grilling round 6

- A Platform is an operating-system family such as `darwin`; a Target is one
  concrete Platform and architecture pair such as `darwin/arm64`.
- `[targets.<platform>]` provides common values and
  `[targets.<platform>.<architecture>]` provides exact Target overrides.
  Resolution proceeds from built-in defaults to Platform values to exact
  Target values, followed by equivalent Profile overrides.
- Each hook phase invokes at most one user-owned script. That script may
  orchestrate other scripts without exposing another task model in the
  manifest.
- Hooks use `[hooks] before_build = "script"` shorthand. A phase that declares
  cache inputs and outputs instead uses a `[hooks.before_build]` table with a
  `script` field. Ejection preserves the user's chosen form and does not infer
  hook cache declarations.
- A newly generated manifest contains only the required `[project]` identity
  and the commented derived `binary_name`. It does not include commented
  examples for the remaining schema.

### 2026-08-15 — Grilling round 7

- Optional cross-platform project metadata uses `company_name`, `description`,
  `copyright`, `comments`, and `icon` under `[project]`. It remains absent in a
  new project and is materialized by ejection.
- `[project].version` is the human-facing application version. An optional
  platform-neutral integer `build_number` maps to platform build revisions and
  is derived when absent; exact Target overrides remain possible.
- Binding generation is configured under `[frontend.bindings]`, including
  language, model representation, time representation, and output directory.
- `[build]` expresses semantic policy such as output directory, production,
  obfuscation, path trimming, and symbol stripping. `[build.go]` contains Go
  tags, linker flags, and compiler flags. It does not contain a raw build
  command.
- `[dev]` contains Wails-owned development settings such as port, debounce,
  logging, watch/include patterns, exclusions, and gitignore handling. The
  Pipeline derives its build, frontend-server, and launch operations; the
  manifest does not reproduce the current arbitrary `executes` list.
- All manifest paths resolve from the directory containing `wails.toml`.
  User-owned scripts use that directory as their default working directory;
  long-form hook configuration may override it.

### 2026-08-15 — Grilling round 8

- Packaging formats are selected per Platform. Format-specific outcome
  configuration is nested under `[package.<platform>.<format>]`; Wails
  translates it into generated tool inputs rather than exposing packager CLI
  arguments.
- `[targets.<platform>]` accepts a limited set of platform identity mappings,
  including identifier, product name, minimum platform version, build number,
  and platform capabilities. These do not permit Profiles to alter identity.
- File associations use a platform-neutral list of extensions and metadata;
  protocols use a platform-neutral scheme and description. Either may
  optionally restrict registration with `platforms`.
- `wails.toml` stores signing policy and non-secret credential references, but
  never passwords or private tokens. Secrets come from supported keychain,
  user-default, or environment mechanisms.
- Merely configuring a signing identity does not enable signing. Automatic
  signing requires explicit `enabled = true` policy.
- `[wake]` is reserved for manifest provenance, including `ejected_by` and
  named-profile ejection metadata. Cache paths, parallelism, verbosity, and
  diagnostics remain CLI execution controls and are not emitted by ejection.

### 2026-08-15 — Grilling round 9, partial

- An unejected project resolves the defaults compiled into the running
  `wails3` CLI. Reproducibility comes from pinning that CLI in automation or
  ejecting; a normal manifest does not carry an explicit default version.
- Wails retains released default snapshots so re-ejection can compare the
  original defaults, the active frozen manifest, and current defaults using
  `ejected_by` as the baseline identifier.
- Ejection is host-independent: it materializes every Target supported by the
  selected Profile even if the current host cannot execute those Targets.
- Ejection writes every meaningful resolved value but omits semantically empty
  tables and lists. In an ejected scope, omission means frozen empty rather
  than inheritance from future defaults.
- Re-ejection never changes active TOML. It adds machine-recognizable candidate
  comments only where an active value still equals its historical default;
  user-customized values are not annotated, stale generated suggestions are
  replaced, and `ejected_by` remains unchanged.
- Ejection edits TOML with comment-preserving, canonically ordered, atomic
  replacement. It creates no backup by default, but supports an explicit
  request to retain one; the backup naming behavior remains open.

### 2026-08-15 — Grilling round 10

- `wails3 eject --backup` writes a timestamped sibling of `wails.toml` before
  atomically replacing the manifest. Ejection creates no backup unless the
  option is explicitly supplied.
- Extension names are lowercase slugs. Core Wails preserves opaque values
  under `[extensions.<name>]` without validating their internal fields or
  inventing defaults; Profiles may overlay those values.
- Ordinary manifest strings do not perform generic environment-variable
  interpolation. Exceptional environment-backed values use documented
  variables or typed references in the configuration area that owns them.
- `default` is a reserved Profile name because the effective base is expressed
  by top-level configuration. Other Profile names are lowercase slugs.
- The initial parser omits and rejects a speculative top-level `schema` field.
  Explicit schema versioning is introduced only if a future incompatible
  format requires it.
- `wails3 config show` is the non-mutating way to inspect fully resolved
  configuration. `wails3 eject` does not need a separate dry-run mode.

### 2026-08-15 — Grilling round 11

- First ejection preserves existing explicit customizations and materializes
  every remaining value from the effective current defaults.
- `wails3 eject <name>` creates a complete named Profile from the effective
  base when the Profile does not exist, or resolves its sparse overrides when
  it does.
- Base and named Profile snapshots are independently frozen. Ejecting one does
  not rewrite another; an ejected Profile no longer changes with its former
  base.
- `ejected_by` records the complete Wails CLI version, including prerelease or
  development identifiers. When a development build's historical defaults
  cannot be recovered, re-ejection leaves the snapshot untouched and reports
  that suggestions are unavailable.
- There is no `uneject` command. Returning to implicit defaults is the manual,
  reviewable operation of reducing the manifest to project identity plus
  desired sparse overrides and removing the relevant ejection metadata.

## Answer

`wails.toml` is a declarative Manifest containing project identity and user
build intent. It never exposes Taskfile tasks, Wake graph nodes, raw build
commands, or generated platform files.

### Minimal manifest

Every new project writes only its stable identity:

```toml
[project]
name = "my-app"
product_name = "My App"
identifier = "com.example.my-app"
version = "0.1.0"
# binary_name = "my-app"
```

`binary_name` is derived and shown as a commented discovery aid. Optional
project metadata includes `build_number`, `company_name`, `description`,
`copyright`, `comments`, and `icon`; it is not emitted until customized or
ejected.

### Vocabulary and structure

- `[project]` owns cross-platform identity and metadata.
- `[frontend]` owns inferred package-manager, script, directory, and output
  overrides; `[frontend.bindings]` owns generated frontend bindings.
- `[build]` owns semantic compilation policy and output; `[build.go]` contains
  Go tags, linker flags, and compiler flags.
- `[dev]` owns Wails-native watch, debounce, logging, frontend-server, and
  launch policy. It has no arbitrary execution list.
- `[targets.<platform>]` contains Platform-wide values and limited identity
  mappings. `[targets.<platform>.<architecture>]` contains exact Target
  overrides. A Platform is an OS family; a Target is a Platform/architecture
  pair.
- `[package.<platform>]` selects output formats. Format settings and optional
  user-owned templates live under `[package.<platform>.<format>]`.
- `[signing.<platform>]` owns explicit signing policy and non-secret credential
  references. An identity alone does not enable signing; `enabled = true` is
  required. Secrets never belong in the Manifest.
- `[[associations]]` and `[[protocols]]` describe platform-neutral operating
  system registrations, optionally restricted with `platforms`.
- `[hooks]` keys stable pipeline phases to one user-owned script each. The
  shorthand is `before_build = "scripts/example.sh"`; a cache-aware hook uses
  `[hooks.before_build]` with `script`, `inputs`, `outputs`, and optional
  `directory` fields. A hook script may call other scripts.
- `[profiles.<name>]` repeats the same nested shape as top-level configuration,
  except it cannot override project identity or identity mappings nested under
  Targets (`identifier`, `product_name`, `version`, or `build_number`). Schema
  validation rejects those fields in Profile scopes. Profiles are sparse
  overlays until independently ejected. There is no `extends`; `default` is
  reserved.
- `[extensions.<lowercase-slug>]` contains opaque extension-owned values. This
  is the only namespace in which core Wails accepts unknown fields.
- `[wake]` contains only provenance: base `ejected_by` and
  `[wake.ejected_profiles]` version entries. Execution controls do not belong
  in the Manifest.

There is no `[tasks]`, `[pipeline]`, global `[templates]`, or speculative
`schema` field. Unknown fields are errors outside extension namespaces.
Ordinary strings do not interpolate environment variables. Every path is
relative to the Manifest directory, which is also the default hook working
directory.

Target resolution proceeds from versioned built-in defaults to Platform values
to exact Target values, followed by equivalent Profile overrides. Shared
metadata is translated into generated platform assets; only a limited set of
identity mappings may vary by Platform or Target.

### Defaults and ejection

An unejected Manifest follows defaults compiled into the running `wails3` CLI.
CI can pin the CLI; users who need a source-level snapshot eject. Wails retains
released default snapshots for later three-way comparison.

```text
wails3 eject             # freeze the effective top-level base
wails3 eject release     # freeze or create the named release Profile
wails3 eject --backup    # also preserve the previous Manifest
```

First ejection preserves explicit user values, materializes every meaningful
resolved value, and freezes every Target supported by the selected Profile.
It is host-independent. Semantically empty tables and lists may be omitted;
within an ejected scope, omission means frozen empty rather than future
inheritance.

Base and named Profile snapshots are independent. Base provenance is
`[wake].ejected_by`; named provenance is stored under
`[wake.ejected_profiles]`. Values contain the complete CLI version. Ejecting a
missing Profile creates a complete copy of the effective base; ejecting an
existing sparse Profile resolves its overrides first.

Re-ejection never changes active TOML. It performs a three-way comparison
between historical defaults, the frozen Manifest, and current defaults, then
adds replaceable, machine-marked comments only beside values that still equal
their historical default. Customized values receive no suggestion and
provenance remains unchanged. If historical defaults cannot be identified,
Wails reports that suggestions are unavailable and leaves the file untouched.

Writes preserve user comments and formatting where practical, insert new
fields in canonical order, and replace the Manifest atomically. No backup is
created by default; `--backup` first writes a timestamped sibling. `wails3
config show` prints resolved configuration without mutation. Returning to
implicit defaults is a manual edit; there is no `uneject` command.
