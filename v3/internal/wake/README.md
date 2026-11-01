# Wails build engine

`internal/wake` is the experimental native YAML build engine, enabled by
`WAILS_EXP_USE_WAKE` and selected by an active project-root `wails.yaml`.
It never executes Taskfiles. Migration is one-way from legacy declarations
into reviewed YAML; the migration parser cannot run shell commands or expand
templates. Unmigrated projects use the separate embedded Task runtime.

Set `WAILS_EXP_USE_WAKE` to opt into the YAML CLI, project generation and setup
writer. Presence enables it, including an empty value, `0` or `false`. When
unset, unmigrated projects retain their embedded Task builds and legacy
project generation. Migrated projects fail with an opt-in diagnostic; they
cannot return to Taskfile execution by removing the flag.

The former `WAILS_USE_WAKE` variable is ignored. Only `WAILS_EXP_USE_WAKE`
enables the experimental YAML CLI.

## Native YAML projects

With the experiment enabled, the presence of `wails.yaml` is the cutover marker. Wails then ignores
all Taskfiles, including local overrides. Invalid manifests and an unset
experimental flag fail without falling back to legacy execution.

The canonical workflows are:

```bash
export WAILS_EXP_USE_WAKE=1
wails3 build
wails3 build release
wails3 build --plan
wails3 build --plan --json
wails3 dev
```

`wails3 package` and `wails3 sign` are deprecated compatibility aliases over
the same planner. New project documentation should use `wails3 build` with a
profile or anonymous `--targets` and `--formats` options.

The manifest is literal configuration. It has no expressions, includes,
inheritance, user-defined stages, or custom graph edges. The experiment supports
six bounded lifecycle hooks that invoke one project-owned script directly; they
do not expose the fixed stage graph as a programmable API. The planner turns
resolved configuration into an immutable typed Plan; the
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

### Desktop compiler selection

Set `CGO_ENABLED=0` in the caller environment or a manifest build environment
to cross-compile desktop targets with Go alone. Wake adds the experimental
`purego` backend tag for macOS and Linux GUI executables. Selecting that tag
explicitly also defaults compilation to `CGO_ENABLED=0`. Windows executables
and CLI/server programs default to CGo disabled; macOS and Linux GUI programs
retain the CGo default unless opted into pure Go. Per-executable environment
settings override the caller's choice. Development preserves the CGo choice
and declared `CC`/`CXX` settings.

Pure-Go desktop builds select the native Go toolchain across operating systems
and architectures. Binding generation uses an isolated Go-only environment
when the primary GUI uses pure Go, so analysis does not require a host C
compiler either. The resolved CGo choice appears in plan output and compiler
cache identities. Target GTK/WebKit libraries remain runtime requirements on
Linux.

With `CGO_ENABLED=1`, Wake uses a compatible native compiler or an explicitly
configured cross compiler (`CC`), then available Zig for Windows cross builds,
then the existing container fallback. Linux CGo cross builds still require
the target desktop libraries and headers; Zig alone does not provide those.

macOS application assembly, universal binary merging and DMG creation run on
all desktop hosts. DMGs use the existing Go image writer when macOS tools are
unavailable. Ad-hoc and exported-PKCS#12 Apple signatures use Go; keychain identities
retain macOS codesign. Team API keys use the Go Notary client, and authenticated
tickets are attached in Go. Notarytool profiles and DMG-editing hooks retain
their macOS tool requirements. Developer ID distribution still needs a valid
identity and a compatible secure timestamp. AppImage deployment, SquashFS writing and zsync generation run in
Go on all desktop hosts, using a pinned target AppImage runtime. Cross-host GUI
packaging requires a supplied Linux runtime tree with GTK/WebKit libraries,
compiled schemas and module caches. Pure-Go GUI runtime profiles follow each
executable's GTK tags, including companions. Selected runtime content is hashed
for cache invalidation; the entire Linux root is never snapshotted.

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
inactive `wails.migrated.yaml` when that path does not exist.

- `--dry-run` writes nothing.
- `--json` prints the versioned analysis to stdout.
- `--output <path>` selects a new inactive project-relative `.yaml` path.
- `--activate` reruns analysis, validates the selected draft, and atomically
  activates it as `wails.yaml` only when no blockers remain.

Migration never overwrites a reviewed draft and never changes Taskfiles,
scripts, assets, or `build/config.yml`. It has no persistent report, backup,
retirement, rollback, or force-activation state. Reachable unrepresented build
behaviour blocks cutover; unrelated utility tasks are reported without
blocking. Conservative, argument-free lifecycle tasks that invoke a contained
script file migrate to the equivalent typed hook.

`wails3 eject` exclusively writes the complete resolved reference manifest to
inactive `wails.ejected.yaml`. `--force` atomically replaces only that inactive
file.

## Legacy Taskfile projects

When the YAML experiment is disabled or no active `wails.yaml` exists,
`build`, `package`, `sign` and `task` use the embedded Task runtime. Root task
customisations, CLI variables, task selection and Task flags retain their
existing behavior. No environment variable selects the old native Taskfile
executor; it has been removed.

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
