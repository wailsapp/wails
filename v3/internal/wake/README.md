# Wake

> ⚠️ **Experimental.** Wake is opt-in behind `WAILS_USE_WAKE=true` and is
> not the default runtime. The default `wails3 build / package / sign /
> task` path is unchanged and runs on the embedded Task runtime as
> before. The API and feature coverage may change between releases. See
> [Enabling Wake](#enabling-wake) for what the env var actually toggles
> and what falls back automatically.

Wake is an **experimental alternative build runner** for `wails3`. It reads
the same `Taskfile.yml` your project already has — same task/dep/var/
template/include/platform-namespace syntax — and runs it through a
Wails-aware executor instead of the general-purpose Task runtime. The goal
isn't to replace Task; it's to provide a runner that's purpose-built for
the way Wails projects actually build, with semantics, output and defaults
that match the rest of the wails3 CLI. If you only ever use Wake, your
Taskfiles do not change.

## How Wake differs from Task

Both Wake and the Task runtime are compiled into `wails3` — neither
requires a separate binary install. The differences are about what each
one is *for*.

### 1. Tailored to the Wails build pipeline

Wake knows the shape of a Wails build: the frontend bundle gets embedded
into the binary, the binary gets packaged into platform-specific
artefacts, icons and bindings are generated alongside. The executor
surfaces those as first-class concepts — for example, the final binary
shows up as a labelled artefact at the end of the run instead of being
inferred from `generates:` globs that may also include intermediates.
Task is a general-purpose tool; Wake is a build runner with opinions
about what a Wails build looks like.

### 2. Focused build semantics → faster incremental builds

Task only skips a task when its Taskfile declares `sources:`/`generates:`.
The Wails templates don't declare these for `go build` / `go mod tidy`,
so Task re-runs the Go compiler and linker on *every* build. Wake derives
those inputs itself from the Go module graph and skips the subprocess
when nothing changed:

| Command        | Inputs Wake tracks                                                                 |
| -------------- | ---------------------------------------------------------------------------------- |
| `go mod tidy`  | `go.mod`, `go.sum`                                                                 |
| `go build`     | module `.go` files + `go.mod`/`go.sum` + the `generates` of the task's dep closure |

The dep-closure tracking means a frontend change (which regenerates
`frontend/dist`, embedded into the binary) still forces a re-link, while
a no-op rebuild skips compilation entirely. The cache key includes the
*expanded* command, so changing build flags (e.g. dev vs `-tags
production`) invalidates correctly.

On a no-op rebuild of the `badge` example this is roughly **~20 ms
(wake) vs ~316 ms (task)**. Cold builds land at the same wall time —
they're dominated by `npm install`, Vite, and the Go compiler.

Cache lives in `.wake/cache.json` (Task uses `.task/`).

### 3. Structured output controlled by wails3

Wake renders through wails3's own **Pulse** reporter: a pre-painted
skeleton of one row per planned step, live status, a colour-coded phase
breakdown at the end, and clickable file:line links inside failure
panels (via OSC 8 hyperlinks where the terminal supports them). Because
output goes through one structured channel rather than raw subprocess
stdout, NO_COLOR and non-TTY environments degrade cleanly (CI logs
stay readable). The visual is the same identity as the rest of the
wails3 CLI; nothing about it changes when the underlying task moves.

### 4. Parallel execution by default

A task's `deps:` are evaluated concurrently rather than sequentially.
The verdict line shows the parallel speedup when one was earned — e.g.
`✓ build succeeded 3.6s  11.6s cpu · 3.2× speedup`. Opt out with
`WAKE_SERIAL=true` when stdout interleaving from sibling steps would
muddle a specific investigation.

### 5. Layered local overrides

Wake supports a **base Taskfile plus local overrides**. Drop a file next
to your `Taskfile.yml` and its definitions take precedence:

| File                                          | Purpose                                  | Precedence  |
| --------------------------------------------- | ---------------------------------------- | ----------- |
| `Taskfile.yml`                                | base, committed                          | lowest      |
| `Taskfile.override.yml` / `.yaml`             | committed team-wide overrides            | middle      |
| `Taskfile.local.yml` / `.yaml`                | personal, usually git-ignored            | highest     |

**Merge semantics (local wins):**

- A task with the **same name** overrides the base task. List fields
  (`cmds`, `deps`, `sources`, `generates`, `platforms`, `status`,
  `preconditions`, `aliases`) **replace** the base when the override
  provides them; fields the override omits are kept from the base.
- `env` and `vars` **merge per key**, with the override winning on
  collisions.
- A task that exists **only** in an override file is **added**.
- Top-level `vars` from an override layer over the base vars (override
  wins per key).

Example — base builds with dev flags, but your machine always builds
production:

```yaml
# Taskfile.yml (committed)
tasks:
  build:
    cmds:
      - go build -o bin/app .
```

```yaml
# Taskfile.local.yml (git-ignored, yours)
tasks:
  build:
    cmds:
      - go build -tags production -o bin/app .
  smoke:
    cmds:
      - ./bin/app --selftest
```

Now `build` runs your production command and `smoke` is available with
no changes to the committed Taskfile.

#### Security & trust model

Override files are auto-discovered and applied with no prompt, so it's
worth being explicit about what that does and doesn't change:

- **No new capability.** A Taskfile already runs arbitrary shell
  commands, so an override can't do anything that editing
  `Taskfile.yml` couldn't. The blast radius is the same: whoever can
  write build files can run code at build time.
- **Same review path.** `Taskfile.override.*` is committed, so it shows
  up in `ls` and PR diffs. `Taskfile.local.*` is created by the
  developer on their own machine (and usually git-ignored).
- **Symlinks are followed** and **includes resolve relative to the
  override file's directory** (including `../`), exactly like the base
  Taskfile — the local layer is not sandboxed.
- **Fail-closed.** A malformed override file (bad YAML, missing
  non-optional include) aborts the run rather than being silently
  skipped.
- **Opt-out for CI / locked-down builds.** Set
  `WAILS_NO_OVERRIDES=true` to skip override discovery entirely and
  build only from the committed base Taskfile, for guaranteed-
  deterministic behaviour.

## Enabling Wake

Wake is **gated entirely behind the `WAILS_USE_WAKE=true` environment
variable**. With the variable unset (or set to anything other than
`true`), every wails3 command — `wails3 build`, `wails3 package`,
`wails3 sign`, `wails3 task <name>` — uses the embedded Task runtime
exactly as before. The default user path is unchanged.

```bash
# default: Task runtime, no wake involvement
wails3 build

# opt in: wake drives build / package / sign / `task <name>`
WAILS_USE_WAKE=true wails3 build
WAILS_USE_WAKE=true wails3 package
WAILS_USE_WAKE=true wails3 task <some-task-name>
```

What the gate covers:

- `commands.wrapTask` (used by `wails3 build / package / sign`) routes
  through `wake.Execute` only when `useWake()` is true.
- `commands.RunTask` (used by `wails3 task <name>`) does the same —
  `wails3 task` is consistent with the other verbs.
- `wails3 dev` is **not** affected by this flag yet; the dev watcher
  still uses its own pipeline.

## Falls back to Task automatically on unsupported features

If wake encounters a Taskfile feature it doesn't implement, it hands the
whole run off to the embedded Task runtime. The Task runtime is
compiled into `wails3` — there is no external `task` binary to install
— so the fallback is in-process and instant. Enabling wake is
therefore always safe: at worst you get the same behaviour you'd get
without the flag.

Features that trigger the fallback:

- `dotenv` at the taskfile level
- `output` modes other than `interleaved`
- `requires` block
- `interval` (taskfile or task level)
- `run` modes other than `always`
- `short` in a task
- `defer` in a task

## Environment variables

| Variable             | Effect                                                                            |
| -------------------- | --------------------------------------------------------------------------------- |
| `WAILS_USE_WAKE`     | `true` enables wake-routable wails3 verbs; anything else uses the Task runtime    |
| `WAILS_NO_OVERRIDES` | `true` skips `Taskfile.local.*` / `.override.*` discovery (deterministic builds)  |
| `WAKE_VERBOSE`       | Stream subprocess stdout/stderr live instead of capturing it for failure-only display |
| `WAKE_SILENT`        | Suppress task output entirely                                                     |
| `WAKE_SERIAL`        | `true` disables the parallel `deps:` fanout (parallel is the default)             |
| `WAKE_FORCE`         | `true` bypasses every cache (task + go-cache) for a true clean rebuild            |
| `WAKE_DEBUG`         | Log resolver internals (DAG, deps, var refs, exec routing)                        |
| `WAKE_NOTICE`        | `off` to silence the per-run "wake (experimental)" notice                         |

## Internals

See [`AGENTS.md`](./AGENTS.md) for the execution pipeline, cache design,
and package-level architecture.
