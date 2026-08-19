# Wails build engine

`internal/wake` contains two separate systems:

- the native HCL build engine selected by an active project-root `wails.hcl`;
- the temporary Go-native Taskfile runner for legacy compatibility.

`WAILS_USE_WAKE` selects only the legacy Taskfile runner. It never enables or
disables the HCL engine.

## Native HCL projects

The presence of `wails.hcl` is the explicit cutover flag. Wails then ignores
all Taskfiles, including local overrides, and an invalid manifest fails without
falling back to legacy execution.

The canonical workflows are:

```bash
wails3 build
wails3 build release
wails3 build --plan
wails3 build --plan --json
wails3 dev
```

`wails3 package` and `wails3 sign` are deprecated compatibility aliases over
the same planner. New project documentation should use `wails3 build` with a
profile or anonymous `--targets` and `--formats` options.

The manifest is literal configuration. The first release has no expressions,
includes, inheritance, hooks, user-defined stages, or custom graph edges. The
planner turns resolved configuration into an immutable typed Plan; the
executor schedules its dependency and resource constraints. Independent
targets and safe package branches may run concurrently.

### Ownership and caching

Referenced assets and package templates are user-owned and read-only. Wails
stages them into disposable `.wails/` workspaces before generation. A failed
compile, signing, assembly, or packaging operation preserves the last complete
generated workspace and final artifact.

Action identities include direct content, semantic Go binding inputs,
consumed-artifact digests, resolved configuration, tool identity, and relevant
environment values. Reproducible outputs use the machine-local content store;
stateful dependency installation uses receipts. Credentials are direct
read-only inputs and are never staged or cached.

### Development sessions

The dev session owns persistent frontend and backend processes plus a
transactional watch set. File bursts request finite development Plans and
cancel stale generations. A failed candidate does not replace the last healthy
application. Development builds ignore production packaging, stripping,
obfuscation, signing, and notarisation policy.

## Migration and ejection

`wails3 migrate` analyses root and included Taskfiles, conventional overrides,
`build/config.yml`, frontend metadata, lockfiles, TypeScript configuration, and
conventional assets. It prints current diagnostics and exclusively creates an
inactive `wails.migrated.hcl` when that path does not exist.

- `--dry-run` writes nothing.
- `--json` prints the versioned analysis to stdout.
- `--output <path>` selects a new inactive project-relative `.hcl` path.
- `--activate` reruns analysis, validates the selected draft, and atomically
  activates it as `wails.hcl` only when no blockers remain.

Migration never overwrites a reviewed draft and never changes Taskfiles,
scripts, assets, or `build/config.yml`. It has no persistent report, backup,
retirement, rollback, or force-activation state. Reachable unrepresented build
behaviour blocks cutover; unrelated utility tasks are reported without
blocking.

`wails3 eject` exclusively writes the complete resolved reference manifest to
inactive `wails.ejected.hcl`. `--force` atomically replaces only that inactive
file.

## Legacy Taskfile projects

Without active `wails.hcl`, existing projects continue through the legacy
Taskfile path. `WAILS_USE_WAKE=true` selects the experimental Go-native runner
instead of the external `task` CLI:

```bash
WAILS_USE_WAKE=true wails3 build
WAILS_USE_WAKE=true wails3 task <task-name>
```

The runner supports v3 parsing, static includes, namespaces, variables, local
override layers, parallel dependencies, reporting, and `.wake/cache.json`.
Unsupported Taskfile features fall back to the external Task CLI when
available. This path is compatibility code, not a customisation language for
HCL projects.

## Verification

Run from `v3/`:

```bash
go test ./internal/wake/... ./internal/commands ./cmd/wails3
go test -race ./internal/wake/... ./internal/commands ./cmd/wails3
go vet ./internal/wake/... ./internal/commands ./cmd/wails3
```

Build and dev performance harnesses are
`scripts/benchmark-manifest-build.go` and
`scripts/benchmark-manifest-dev.go`. Native acceptance is driven by
`scripts/verify-manifest-build-system.go`.

See [`AGENTS.md`](./AGENTS.md) for package boundaries and implementation
invariants.
