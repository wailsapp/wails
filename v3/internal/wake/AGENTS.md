# Wake - Go-native Build Systems for Wails v3

Follow the [repository-wide instructions](../../../AGENTS.md) and the [v3 project guide](../../AGENTS.md). This file adds Wake-specific architecture and development rules for `v3/internal/wake/`.

## Overview

Wake is a native YAML build system. It never executes Taskfiles or Taskfile
commands. Migration is one-way: read and classify legacy declarations, write
an inactive YAML draft, validate it, then activate `wails.yaml`. There is no
reverse conversion or fallback from an active manifest to Taskfiles.

`WAILS_EXP_USE_WAKE` is a presence-based opt-in (empty and false still enable
it). `WAILS_USE_WAKE` is retired and ignored. Unmigrated projects and explicit
`wails3 task` commands use the separate embedded Task runtime.

## Architecture

- `manifest/`: strict YAML defaults, validation, profiles and ejection to YAML.
- `pipeline/`: typed immutable plans, resource-aware execution and receipts.
- `cache/`: snapshots, receipts, action index and artifact restoration.
- `migration/`: ephemeral migration analysis; `taskfile/` only parses legacy
  declarations. It must not expand templates, evaluate shell variables or run
  commands. Unknown behaviour blocks activation.
- `packagetemplate/`: atomic copying of user-owned package files.
- `benchmark/`: standard-library measurement and artifact verification.
- `buildinfo/`: shared target and package capability registry.

## Execution flow

### Manifest pipeline

1. Load and strictly validate sparse `wails.yaml`; resolve defaults and Profile.
2. Plan one immutable typed graph for every requested Target and package format.
3. Snapshot direct inputs and calculate tool/environment-aware Action Keys.
4. Prune or restore reusable Artifacts; validate stateful Receipts.
5. Run the remaining critical path with bounded CPU, memory, and exclusive-tool
   claims. The default memory capacity is one logical GiB per worker; constrained
   callers may provide a tighter `MemoryLimitMB`.
6. Generate platform state under `.wails/` and report user-visible Artifacts.

Package templates are user-owned files or directory trees, copied byte-for-byte
into an isolated workspace. They are not a template language: `.tmpl` files
receive no interpolation. Reject symlinks and atomically replace the generated
destination. Adapters must never mutate the platform-assets Artifact they consume.

The Plan has six bounded `RunHook` phases: `before_build`, `after_build`,
`before_package`, `after_package`, `before_sign`, and `after_sign`. Hooks invoke
one project-owned script directly and must remain separate from built-in stage
implementation. Preserve project/target/package barrier scope, versioned JSON
context through `WAILS_HOOK_CONTEXT_FILE`, process-group
cancellation, path containment, and default `CacheNever`. The generated context
path is ephemeral and must not affect cache identity; its semantic schema
version and resolved scope must. Remove context after success or cancellation,
retain it with a reported path after failure, and never copy inherited secrets
into it. A cacheable hook requires complete inputs and outputs; its bounded
output root is the only reusable Artifact. Typed tool calls and arbitrary
user-authored graph edges remain outside the public schema.

The Dev Session is outside the Plan: it owns persistent frontend/backend
processes and replaceable watch sets while requesting ordinary finite
development Plans. Treat a generation update as cancellation, not failure.
Use compile and `after_build` results to avoid replacing an unchanged backend;
stage changed frontend/port sessions before the backend that consumes them.
Keep candidate `FRONTEND_DEVSERVER_URL`/`WAILS_VITE_PORT` values scoped to the
build and child process, pin both sides to the same loopback address, replace
watches transactionally, and terminate/reap complete process groups on every
exit path.

### Migration cutover

`wails3 migrate` classifies Taskfiles against current embedded ASTs and known
historical fingerprints, computes reachability from Wails entry points, and
prints an ephemeral analysis. It exclusively creates inactive
`wails.migrated.yaml`; reruns never replace a reviewed draft. `--activate`
reruns current analysis, validates the selected draft, and atomically creates
active `wails.yaml` only when no reachable unrepresented behaviour remains.
Migration never modifies, retires, renames, or backs up legacy sources. There
is no persistent migration report. Active `wails.yaml` is the sole cutover and
routing signal; invalid YAML never falls back to Taskfiles.

## Build Reporting (UI)

Wake renders builds through a small reporting layer, kept separate from the
executor so it can be reused and so producers in other processes can feed it.

- **`internal/report`** — leaf contract (no lipgloss, no executor deps):
  - `Reporter` interface: `BuildStart` / `StepStart` / `StepInfo` / `StepCommand`
    / `StepOutput` / `StepEnd` / `StepFailed` / `BuildEnd`. `Nop` is the safe
    default (used by tests and a reporter-less `Executor`).
  - `Verbosity`: `Silent` < `Normal` < `Verbose` < `Debug`.
  - Wire protocol: `Encode`/`Decode`/`Emit`/`Route` move `Event`s across a
    process boundary on a single stdout line prefixed with a `\x1f`-delimited
    sentinel that never occurs in ordinary output.
- **`internal/report/termui`** — the renderer. On a TTY at `Normal` it draws one
  live spinner line per step, rewritten in place to a final `✓`/`·`/`✗` line
  with `[k/N]` and a duration. With no TTY / `NO_COLOR` it degrades to plain
  one-line-per-step output. No emojis (✓ U+2713, ✗ U+2717, `·`, Braille spinner).

**Verbosity behaviour (capture-by-default):** subprocess stdout/stderr is always
captured by the native command adapter and shown **only on failure**, so a clean
`[k/N] task  1.2s ✓` is the dominant signal. `Verbose` additionally streams each
command (`$ …`) and its output live; `Debug` adds resolver internals
(`[wake-debug] …`, dep wiring, DAG order) on stderr. Set via `WAKE_SILENT` /
`WAKE_VERBOSE` / `WAKE_DEBUG` (mapped to `ExecuteOptions` in `task_wrapper.go`).

**Step counting:** the native Plan defines the operation count, including
project, target, package and artifact collection nodes.

**Producer feedback (cross-process):** the executor sets `WAKE_REPORT=1` for
child processes. `wails3 generate bindings` (a subprocess) detects it and routes
its `config.Logger` through `commands.wakeLogger`, emitting wire events on stdout
instead of drawing its own spinner. The parent's `captureWriter` decodes those
events and forwards them to the live step as `StepInfo` sub-lines. In-process
producers can instead use `report.Active()`.

**Failures:** a failed command's captured output, command string and exit code
are rendered in a bordered panel; the native command adapter wraps the error as
`errReported` and `cmd/wails3/main.go` skips re-printing it (no double output).

## Testing

```bash
go test ./internal/wake/...
go test ./internal/commands ./cmd/wails3
go test -race ./internal/wake/... ./internal/commands ./cmd/wails3
go vet ./internal/wake/... ./internal/commands ./cmd/wails3
```

All packages have tests. Run from `wails-v3/v3/` directory.

## Development Rules

- Experimental code stays on feature branches, never pushed to `master`
- Run `go test ./internal/wake/...` and `go vet ./internal/wake/...` before committing
- Store planning docs in `history/` directory
