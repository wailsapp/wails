# Wake

Wake is a Go-native build runner embedded in the Wails v3 CLI. It executes
`Taskfile.yml` files directly, replacing the external [Task](https://taskfile.dev)
(`task`) binary for the build/package/sign commands.

It is designed as a **drop-in** for the subset of Taskfile syntax that Wails
projects use: the same `Taskfile.yml`, the same task/dep/var/template syntax,
the same includes and platform namespaces. If you only ever use Wake, your
Taskfiles do not change.

## Enabling Wake

Wake is **gated entirely behind the `WAILS_USE_WAKE=true` environment
variable**. With the variable unset (or set to anything other than `true`),
every wails3 command — `wails3 build`, `wails3 package`, `wails3 sign`,
`wails3 task <name>` — uses the embedded Task runtime exactly as before. The
default user path is unchanged.

```bash
# default: Task runtime, no wake involvement
wails3 build

# opt in: wake drives build / package / sign / `task <name>`
WAILS_USE_WAKE=true wails3 build
WAILS_USE_WAKE=true wails3 package
WAILS_USE_WAKE=true wails3 task <some-task-name>
```

What the gate covers:

- `commands.wrapTask` (used by `wails3 build / package / sign`) routes through
  `wake.Execute` only when `useWake()` is true.
- `commands.RunTask` (used by `wails3 task <name>`) does the same — `wails3
  task` is consistent with the other verbs.
- `wails3 dev` is **not** affected by this flag yet; the dev watcher still
  uses its own pipeline.

If wake encounters a Taskfile feature it does not implement (`dotenv:`,
`requires:`, `interval:`, `run:` modes other than `always`, `short:`,
`defer:`, `output:` modes other than `interleaved`), it **transparently
falls back to the embedded Task runtime** for the whole run. No external
`task` binary is required — the Task runtime is compiled into wails3.
Enabling wake is therefore always safe: at worst you get the same behaviour
you'd get without the flag.

### Environment variables wake responds to

| Variable           | Effect (only when `WAILS_USE_WAKE=true`)             |
| ------------------ | ---------------------------------------------------- |
| `WAILS_USE_WAKE`   | `true` enables every wake-routable wails3 verb       |
| `WAILS_NO_OVERRIDES` | `true` skips `Taskfile.local.*` / `.override.*`    |
| `WAKE_VERBOSE`     | Stream subprocess stdout/stderr live instead of capturing |
| `WAKE_SILENT`      | Suppress task output                                 |
| `WAKE_SERIAL`      | `true` disables parallel `deps:` fanout (parallel is the default) |
| `WAKE_FORCE`       | `true` bypasses every cache (task + go-cache) for a clean rebuild |
| `WAKE_DEBUG`       | Log resolver internals (DAG, deps, var refs, exec)   |

## How Wake differs from Taskfile

### 1. No external binary

Task is a separate executable. Wake is compiled into `wails3`, so there is no
`task` dependency to install or version-match, and task execution happens
in-process.

### 2. Incremental build cache for Go commands

This is the main reason Wake exists. Task only skips a task when its Taskfile
declares `sources:`/`generates:`. The Wails templates don't declare these for
`go build` / `go mod tidy`, so Task re-runs the Go compiler and linker on
*every* build.

Wake derives the inputs for those commands itself and skips the subprocess when
nothing changed:

| Command        | Inputs Wake tracks                                                                 |
| -------------- | ---------------------------------------------------------------------------------- |
| `go mod tidy`  | `go.mod`, `go.sum`                                                                 |
| `go build`     | module `.go` files + `go.mod`/`go.sum` + the `generates` of the task's dep closure |

The dependency-closure tracking means a frontend change (which regenerates
`frontend/dist`, embedded into the binary) still forces a re-link, while a
no-op rebuild skips compilation entirely. The cache key includes the *expanded*
command, so changing build flags (e.g. dev vs `-tags production`) invalidates
correctly.

On a no-op rebuild of the `badge` example this is roughly **~20 ms (wake) vs
~316 ms (task)**. Cold builds are about the same, since they're dominated by
`npm install`, Vite, and the Go compiler — all external.

Cache lives in `.wake/cache.json` (Task uses `.task/`).

### 3. Layered local overrides

Wake supports the concept of a **base Taskfile plus local overrides**. Drop a
file next to your `Taskfile.yml` and its definitions take precedence:

| File                                          | Purpose                                  | Precedence  |
| --------------------------------------------- | ---------------------------------------- | ----------- |
| `Taskfile.yml`                                | base, committed                          | lowest      |
| `Taskfile.override.yml` / `.yaml`             | committed team-wide overrides            | middle      |
| `Taskfile.local.yml` / `.yaml`                | personal, usually git-ignored            | highest     |

Layers are applied in that order, so a personal `Taskfile.local.yml` wins over a
committed `Taskfile.override.yml`, which wins over the base.

**Merge semantics (local wins):**

- A task with the **same name** overrides the base task. List fields (`cmds`,
  `deps`, `sources`, `generates`, `platforms`, `status`, `preconditions`,
  `aliases`) **replace** the base when the override provides them; fields the
  override omits are kept from the base.
- `env` and `vars` **merge per key**, with the override winning on collisions.
- A task that exists **only** in an override file is **added**.
- Top-level `vars` from an override layer over the base vars (override wins per
  key).

Example — base builds with dev flags, but your machine always builds production:

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

Now `build` runs your production command and `smoke` is available, with no
changes to the committed Taskfile.

#### Security & trust model

Override files are auto-discovered and applied with no prompt, so it's worth
being explicit about what that does and doesn't change:

- **No new capability.** A Taskfile already runs arbitrary shell commands, so an
  override can't do anything that editing `Taskfile.yml` couldn't. The blast
  radius is the same: whoever can write build files can run code at build time.
- **Same review path.** `Taskfile.override.*` is committed, so it shows up in
  `ls` and PR diffs. `Taskfile.local.*` is created by the developer on their own
  machine (and usually git-ignored).
- **Symlinks are followed** and **includes resolve relative to the override
  file's directory** (including `../`), exactly like the base Taskfile — the
  local layer is not sandboxed.
- **Fail-closed.** A malformed override file (bad YAML, missing non-optional
  include) aborts the run rather than being silently skipped.
- **Opt-out for CI / locked-down builds.** Set `WAILS_NO_OVERRIDES=true` to skip
  override discovery entirely and build only from the committed base Taskfile,
  for guaranteed-deterministic behavior.

## Unsupported features (automatic fallback to `task`)

If a Taskfile uses any of these, Wake hands the whole run to the `task` runtime:

- `dotenv` at the taskfile level
- `output` modes other than `interleaved`
- `requires` block
- `interval` (taskfile or task level)
- `run` modes other than `always`
- `short` in a task
- `defer` in a task

## Environment variables

| Variable           | Effect                                            |
| ------------------ | ------------------------------------------------- |
| `WAILS_USE_WAKE`   | `true` enables Wake; otherwise uses `task`        |
| `WAILS_NO_OVERRIDES` | `true` skips `Taskfile.local.*`/`.override.*` (deterministic builds) |
| `WAKE_VERBOSE`     | Log the resolved DAG and per-task detail          |
| `WAKE_SILENT`      | Suppress task output                              |
| `WAKE_PARALLEL`    | Execute independent tasks concurrently (opt-in)   |
| `WAKE_DEBUG`       | Log native Go cache misses                        |

## Internals

See [`AGENTS.md`](./AGENTS.md) for the execution pipeline, cache design, and
package-level architecture.
