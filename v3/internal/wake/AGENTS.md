# Wake - Go-native Build System for Wails v3

## Overview

Wake is a Go-native Taskfile executor embedded in `wailsapp/wails/v3/internal/wake/`. It replaces the external `task` CLI dependency, executing `Taskfile.yml` files directly in Go. Enabled via `WAILS_USE_WAKE=true` env var; falls back to `task` CLI when unset or unsupported features are used.

## Architecture

```
wake.go              Entry point: Parse -> Resolve -> DAG -> Execute (serial or parallel)
ast/                 Taskfile AST types + deep Clone() for include isolation
parse/               YAML parsing, include resolution, var/shell expansion, template expansion
resolve/             DAG builder (topological sort, cycle detection), platform filtering
exec/                Task executor, cache (SHA-256 hash + mtime), glob matching with ** support
fallback/            Falls back to external `task` CLI when WAKE disabled or unsupported
override/            Taskfile.local.yml / Taskfile.override.yml merge support
cmds/                Command routing: shell, native Go (go build, npm, etc), frontend ops
platform/            OS/arch detection
```

Build output is rendered through `internal/report` (a leaf contract) and
`internal/report/termui` (the lipgloss renderer). See "Build Reporting" below.

## Execution Flow

1. **Discover** `Taskfile.yml` / `Taskfile.yaml` / `build/Taskfile.yml`
2. **Parse** YAML into AST, recursively resolve includes with cycle detection
3. **Clone** included tasks (`task.Clone()`) to prevent pointer aliasing across namespaces
4. **Populate** built-in vars (`OS`, `ARCH`, `ROOT_DIR`, `TASKFILE_DIR`, etc.)
5. **Layer** local overrides (`override.LoadLocal` + `resolve.MergeTaskfile`): two stacked layers in increasing precedence — `Taskfile.override.*` (committed team layer) then `Taskfile.local.*` (personal layer, wins last); within a layer the first existing extension wins. Local-wins precedence — a same-named task replaces the base task (list fields replace, `env`/`vars` maps merge per-key); local-only tasks are added; top-level vars layer over base.
6. **Filter** platforms (remove non-matching `darwin:`/`linux:`/`windows:` namespaces)
7. **Resolve** vars (shell execution, ref chains, cycle detection)
8. **Expand** templates (`{{.VAR}}`) in task names, dirs, deps, cmds, env
9. **Resolve** dep namespaces (short deps -> fully qualified `prefix:task`)
10. **Build DAG** from target task's dependency tree
11. **Execute** serially (`ex.Execute`) or parallel (`executeParallel` via DAG in-degree)
12. **Cache** results in `.wake/cache.json` (hash + last_run timestamp)

## Parallel Execution

DAG-based: tasks with in-degree 0 grouped into waves, each wave executed via goroutines with `sync.WaitGroup`. Deadlock detection if no ready tasks remain but not all completed. Enabled via `ExecuteOptions.Parallel` or `WAKE_PARALLEL` env var.

## Cache System

- **Location**: `.wake/cache.json` per project
- **Key**: Fully namespaced task name (e.g., `darwin:common:generate:icons`)
- **Hash**: SHA-256 of task name + cmds + env + sources + generates patterns
- **Skip logic** (`ShouldSkip`):
  - Returns `false` if no sources/generates/status (task always runs)
  - Status commands: skip if all exit 0
  - Generates: skip if all files exist
  - Hash: skip if matches cache entry
  - Sources: skip if no source file mtime > cache `last_run`
- **Task Dir**: Sources/generates resolved relative to `task.Dir` (joined with `baseDir`), not `baseDir` directly
- **Glob**: `**` patterns handled via `recursiveMatch`; `exclude:` prefix filters sources

## Native Go Command Auto-Cache (wake-exclusive)

`exec/gocache.go` gives wake caching the external `task` CLI cannot match. For a single-command task that declares no explicit sources/generates/status and whose command is `go build` or `go mod tidy` (post template-expansion), wake derives the cache inputs itself and skips the subprocess when nothing changed:

- **`go mod tidy`**: inputs = `go.mod`, `go.sum`.
- **`go build`**: inputs = module-local non-test `.go` files + `go.mod`/`go.sum` + the resolved `generates` of the task's transitive dependency closure (this captures embedded assets like `frontend/dist`, so a frontend change still forces a re-link). Output existence (`-o` path) is also required for a hit.

Cache entries use the `#go` key suffix so they don't collide with `RecordTask`. The entry hash is the SHA-256 of the *expanded* command, so changing build flags (e.g. dev vs `-tags production`) invalidates correctly. Set `WAKE_DEBUG=1` to log cache misses.

Result (badge example, no-op rebuild): wake ~20ms vs `task` ~316ms (~94% faster) — wake skips both `go build` and `go mod tidy`; `task` always runs them.

**Limitation**: a file embedded directly by the main package that is neither a `.go` file nor produced by a dependency's `generates` won't be tracked. The Wails template routes embedded assets through `frontend/dist` (a dep generate), so the common case is covered.

## Known Gotchas

- **Var resolution is fixed-point**: `mergeVars` expands templated vars in a loop until stable (max 10 passes). A single pass was Go-map-order-dependent, so self-referential vars like `OUTPUT: '{{.OUTPUT | default .DEFAULT_OUTPUT}}'` resolved to empty ~50% of runs (causing a stray root binary + cache thrash). Don't revert to a single pass.
- **Glob prefix**: `recursiveGlob`/`recursiveMatch` only walk/match the literal path prefix before `**`. Previously `frontend/dist/**/*` matched the entire tree (including `.wake/cache.json`), which is why `sources`/`generates` were historically unreliable.
- **Empty sources/generates**: Tasks with no sources/generates/status still ALWAYS run via the normal path — *unless* they match the native Go auto-cache above.
- **Task cloning**: `ResolveIncludes` MUST clone tasks via `task.Clone()` before adding to `tf.Tasks` under namespaced keys. Without cloning, `expandTemplates` mutates shared pointers, corrupting all namespace entries.
- **Task Dir in cache**: `ShouldSkip` computes `taskDir` from `task.Dir` before globbing sources/generates. Tasks with `dir: build` resolve patterns relative to `build/`, not project root.
- **Glob exclude**: `recursiveMatch` treats `**/*` suffix as matching everything under the prefix. Pattern `frontend/**/*` matches any file under `frontend/`.
- **Namespace filtering**: `filterTaskNamespaces` removes `common:` tasks and non-matching platform prefixes when target is `darwin:*` / `linux:*` / `windows:*`.
- **Dep namespace resolution**: Short dep names resolved to full namespace via `resolveDepNamespaces` (pre-execution) and `resolveTaskName` (runtime). Tries `prefix:task`, then `prefix:common:task`, then `include:task`.
- **Fallback**: When `WAILS_USE_WAKE` not set or unsupported features detected (dotenv, output modes, defer, interval, short), falls back to `task` CLI if available.

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
captured (via `exec/capture.go`) and shown **only on failure**, so a clean
`[k/N] task  1.2s ✓` is the dominant signal. `Verbose` additionally streams each
command (`$ …`) and its output live; `Debug` adds resolver internals
(`[wake-debug] …`, dep wiring, DAG order) on stderr. Set via `WAKE_SILENT` /
`WAKE_VERBOSE` / `WAKE_DEBUG` (mapped to `ExecuteOptions` in `task_wrapper.go`).

**Step counting:** `countSteps` (wake.go) walks deps **and** task-ref commands
from the target and counts tasks with a real command — this is the `N`. The
executor reports each task exactly once through the `beginStep`/`reported`
chokepoint, so `k ≤ N`. A cache hit that prunes a dependency subtree reports the
pruned tasks as `cached` (`reportPrunedCached`) so an incremental build still
reaches `[N/N]` truthfully.

**Producer feedback (cross-process):** the executor sets `WAKE_REPORT=1` for
child processes. `wails3 generate bindings` (a subprocess) detects it and routes
its `config.Logger` through `commands.wakeLogger`, emitting wire events on stdout
instead of drawing its own spinner. The parent's `captureWriter` decodes those
events and forwards them to the live step as `StepInfo` sub-lines. In-process
producers can instead use `report.Active()`.

**Failures:** a failed command's captured output, command string and exit code
are rendered in a bordered panel; `wake.Execute` wraps the error as
`errReported` and `cmd/wails3/main.go` skips re-printing it (no double output).

## Unsupported Features (triggers fallback)

- `dotenv` at taskfile level
- `output` modes other than `interleaved`
- `requires` block
- `interval` (taskfile or task level)
- `run` modes other than `always`
- `short` in tasks
- `defer` in tasks

## Testing

```bash
go test ./internal/wake/...
```

All packages have tests. Run from `wails-v3/v3/` directory.

## Speed Test

```bash
/tmp/wake_speed_test.sh          # Cold build comparison
/tmp/wake_speed_test_cached.sh   # Cached build comparison
```

Current results (badge example, no-op cached build): wake **~20ms** vs task CLI **~316ms** (~94% faster). Wake skips every task including `go build`/`go mod tidy` via the native Go auto-cache; the task CLI re-runs both because the Taskfile declares no sources/generates for them. Cold builds are ~parity (dominated by `npm install` + vite + Go compile, all external).

## File Reference

| File | Purpose |
|------|---------|
| `wake.go` | Entry point, orchestration, parallel execution, platform filtering |
| `ast/ast.go` | Taskfile AST types, `Task.Clone()` deep copy |
| `ast/walk.go` | AST visitor pattern |
| `parse/parse.go` | YAML parsing, include resolution, var resolution, builtins |
| `parse/expr.go` | Template expression expansion (`{{.VAR}}`) |
| `resolve/dag.go` | DAG builder, topological sort, cycle detection |
| `resolve/filter.go` | Platform filtering |
| `resolve/merge.go` | Override merging |
| `exec/exec.go` | Task executor, dep resolution, var merging, command execution |
| `exec/cache.go` | Cache load/save, hash computation, ShouldSkip, glob matching |
| `exec/precond.go` | Precondition checking |
| `exec/runner.go` | Run confirmation |
| `exec/capture.go` | Per-command output capture + wire-event routing |
| `report/report.go` | Reporter contract, verbosity, wire protocol (leaf pkg) |
| `report/termui/termui.go` | Terminal renderer (spinner, [k/N], panels) |
| `fallback/task_cli.go` | External task CLI fallback |
| `override/override.go` | Taskfile.local.yml / override.yml loading |
| `cmds/shell.go` | Shell command execution |
| `cmds/native.go` | Native Go command routing |
| `cmds/frontend.go` | Frontend build ops |
| `cmds/task_ref.go` | Task reference resolution |
| `platform/detect.go` | OS/arch detection |

## Development Rules

- Experimental code stays on feature branches, never pushed to `master`
- Use `bd` for all issue tracking (see root AGENTS.md)
- Run `go test ./internal/wake/...` and `go vet ./internal/wake/...` before committing
- Store planning docs in `history/` directory
