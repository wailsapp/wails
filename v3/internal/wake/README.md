# Wake

Wake is the build engine inside Wails v3. It now contains two deliberately
separate systems:

- the manifest-native build system used automatically by projects with an
  active root `wails.toml`; and
- the experimental Go-native Taskfile runner retained for legacy Taskfile
  compatibility and migration.

`WAILS_USE_WAKE` controls only the legacy Taskfile runner. It does not enable or
disable the manifest-native pipeline.

## Manifest-native build system

New and fully migrated projects keep one sparse `wails.toml` at the project
root. Wails supplies versioned defaults from the installed CLI, so the file
normally contains only project identity and intentional customizations.

The normal commands are built in:

```bash
wails3 build
wails3 package
wails3 sign
wails3 dev
```

The planner turns each command into one immutable typed graph. Project work is
shared across Targets, Target and package outputs have one owner, and the
executor schedules the critical path with CPU, memory, and exclusive-tool
claims. Safe package branches and independent Target compiles overlap. Legacy
tools that mutate process state, such as AppImage generation, run in an
isolated subprocess.

### Caching

Inputs are discovered from typed Node semantics rather than user-maintained
`sources` lists. Action identities include:

- direct content snapshots;
- semantic Go binding API snapshots;
- consumed Artifact digests;
- resolved configuration;
- tool and handler identity; and
- relevant environment values.

Reproducible outputs are stored in a machine-local content-addressed Artifact
store and restored if a generated output is missing. Dependency installation
uses a Receipt rather than archiving `node_modules`. Signing and undeclared
side-effect hooks are never reusable.

### Customization

Stable customization belongs in `wails.toml`. File-script hooks are available
at six fixed barriers:

- `before_build`, `after_build`;
- `before_package`, `after_package`; and
- `before_sign`, `after_sign`.

Hooks call one project-owned script file and receive stable `WAILS_*`
environment values. Package formats may use strict, versioned user-owned file
or directory templates. Generated platform state lives under ignored
`.wails/` directories and is not the primary customization API.

`wails3 eject` freezes the complete resolved default configuration into the
manifest. `wails3 eject <profile>` freezes only that named profile. Re-ejection
never silently changes active values, and `--backup` creates a backup only when
requested.

### Dev Session

The manifest Dev Session owns persistent frontend and backend processes plus a
transactional watch set. File bursts request ordinary finite development Plans;
new generations cancel stale work without reporting a failed build. A healthy
application remains running until a changed replacement is built and ready.
No-op and restored binaries do not restart the backend.

## Migration and routing

`wails3 migrate` parses legacy v3 Taskfiles with Wake's AST parser. It compares
them with current embedded canonical templates and retained historical
fingerprints, classifying each file as current default, historical default,
customised, or custom.

Known configuration and script-file hooks are translated. Unsupported inline
shell receives a stable diagnostic and remains a manual migration. Migration
state, source digests, classifications, and diagnostics live in
`.wails/migration-report.json`, never in `wails.toml`.

An incomplete migration keeps the Taskfiles and routes through the legacy
system. A complete migration digest-checks and retires represented legacy
files; `--backup` first preserves them under `.wails/migration-backup`. A
project containing both systems without a migration report fails as ambiguous
instead of choosing silently.

## Legacy Taskfile runner

For projects that still use Taskfiles, `WAILS_USE_WAKE=true` selects the
experimental Go-native runner instead of the external `task` CLI:

```bash
WAILS_USE_WAKE=true wails3 build
WAILS_USE_WAKE=true wails3 task <task-name>
```

It supports v3 Taskfile parsing, includes, platform namespaces, variables,
local override layers, parallel dependencies, structured Wails reporting, and
its original `.wake/cache.json` task cache. Unsupported Taskfile features fall
back to the external Task CLI when available. This path is compatibility code,
not a second long-term customization language for manifest projects.

Legacy-runner environment variables:

| Variable | Effect |
| --- | --- |
| `WAILS_USE_WAKE` | Select the Go-native runner for a Taskfile project |
| `WAILS_NO_OVERRIDES` | Ignore `Taskfile.local.*` and `Taskfile.override.*` |
| `WAKE_VERBOSE` | Stream commands and subprocess output |
| `WAKE_SILENT` | Suppress task output |
| `WAKE_SERIAL` | Disable legacy dependency fan-out |
| `WAKE_FORCE` | Bypass legacy task cache entries |
| `WAKE_DEBUG` | Show resolver and execution diagnostics |

## Verification

Run the focused suite from `v3/`:

```bash
go test ./internal/wake/... ./internal/commands ./cmd/wails3
go test -race ./internal/wake/... ./internal/commands ./cmd/wails3
go vet ./internal/wake/... ./internal/commands ./cmd/wails3
```

Complete-command performance and native-host verification live in
`scripts/benchmark-manifest-build.go` and
`scripts/verify-manifest-build-system.go`. The permanent matrix and measured
baseline are recorded under `history/wayfinder/wails-build-system/`.

See [`AGENTS.md`](./AGENTS.md) for package boundaries, invariants, and detailed
execution flow.
