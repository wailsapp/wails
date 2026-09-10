# Development application arguments verification

Date: 2026-09-10. Branch: `codex/hcl-build-system`.

## Implemented contract

`dev.args` and `target "<os>[/<arch>]" { dev { args = [...] } }` provide literal
application argument vectors. Architecture overrides replace OS-wide defaults;
explicit empty lists clear them. These settings are independent of `run.args`.

Both `dev` and `run` accept arguments after `--` and the `--appargs` compatibility
spelling. The latter accepts one POSIX-quoted string on every host and performs no
shell execution or variable expansion. Mixed sources, repeated compatibility flags
and invalid quoting fail before launch. Explicit empty values clear defaults.
Development arguments require an active HCL project; legacy Taskfiles are not
rewritten.

Development resolution and restart decisions are shared across platforms. Desktop
launchers pass the vector directly. Windows rollback retains the previous process
arguments. iOS simulator/device launchers forward arguments, and mobile rollback
retains them alongside the previous staged bundle. Physical iOS readiness values
now use devicectl's JSON environment option. Android rejects nonempty arguments
before device discovery or build effects instead of dropping them.

## Automated checks

All passed:

```sh
# From v3/
go test ./internal/dev ./internal/commands ./internal/wake/... ./cmd/wails3
go test -race ./internal/dev ./internal/commands ./internal/wake/manifest ./cmd/wails3
go vet ./internal/dev ./internal/commands ./internal/wake/manifest ./cmd/wails3
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/wails3-dev-args.exe ./cmd/wails3
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go test -c -o /tmp/dev-commands-windows.test.exe ./internal/commands
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o /tmp/wails3-dev-args-darwin ./cmd/wails3
```

The session test covers both overlapping process replacement and the stop/restore
lifecycle used by Windows and mobile. It checks source rebuilds, argument-only
changes, failed replacement and recovery. The real Go lifecycle test verifies the
same arguments after rebuilds and an unchanged binary timestamp after an
argument-only manifest edit; the frontend process is unchanged.

CLI subprocess tests cover defaults, overrides, `--help`, paths containing spaces,
Windows paths, empty values, single/double-dash aliases, mixed sources and malformed
quoting. Schema tests cover OS/architecture precedence, Android clearing,
independent run defaults and ejection. Android rejection is also verified through
the actual CLI using the Android example.

`cd docs && npm run build` passes, building 197 pages and validating internal links.
CodeRabbit reviewed the change and found one documentation issue: the generated
field descriptions omitted Android's restriction. The schema source now describes
that restriction, and the generated reference has been refreshed.

## Native Linux Wails acceptance

Build the local runtime package first (`npm ci && npm run build:code` from
`v3/internal/runtime/desktop/@wailsio/runtime`), then run from `v3/`:

```sh
go build -o /tmp/wails3-dev-args ./cmd/wails3
python3 internal/dev/testdata/native_acceptance.py \
  --cli /tmp/wails3-dev-args \
  --output /tmp/wails-dev-args-acceptance \
  --expected-app-args '["--config-path","profile with spaces.yaml","--help"]' \
  -- --port 19347 --appargs '--config-path "profile with spaces.yaml" --help'
```

The native GTK4 / WebKitGTK 6.0 run passed all checks in
`dev-arguments-native-results.json`: runtime readiness, initial and restarted
argument delivery, Go/frontend binding round trips, stderr logging, HMR without a
backend restart, survival of a failed build, rebuild/restart, graceful shutdown,
unchanged production output and port reuse. The first attempt stopped because the
local runtime package was not built; the successful run followed the documented
runtime build prerequisite.

## Remaining native checks

Cross-compilation and platform-neutral adapter tests do not establish native
Windows/macOS/iOS execution. On those hosts, use the same argument vector with
`wails3 dev -- --config-path "profile with spaces.yaml" --help`, then trigger a Go
source rebuild. For the manifest-edit check, stop the CLI, put the arguments in
`dev.args`, and restart without CLI arguments; a CLI override intentionally wins
over later manifest edits. Check exact argument delivery,
frontend continuity, backend restart and rollback after a rejected candidate.
Windows can run `go test ./internal/commands -run TestManifestWindowsBackendStartupFailureRestoresPreviousImage`.

For iOS, add `--target ios/arm64 --device <udid>` before `--`; physical devices
also require `--destination device --host <Mac LAN IP>` and matching signing
configuration. Confirm readiness as well as argument delivery because physical
iOS now uses the JSON environment option. Android nonempty arguments should fail
before any device effects; `target "android"` can clear shared defaults with an
empty `dev.args` list. These checks remain under the existing native acceptance
ticket and do not block the verified Linux implementation.
