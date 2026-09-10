# Wails development sessions

This is a Wails-specific **internal** subsystem for the experimental HCL CLI
(issue #6021), not a public library or a general process supervisor.

`WAILS_EXP_USE_WAKE` must be present and the project must have `wails.hcl`.
With the variable unset, the CLI keeps its existing Taskfile/Refresh path and
none of the new flags are registered. `WAILS_USE_WAKE` remains ignored.

The CLI owns flags and signals. `Run` owns the session context, build
generations, transactional watch replacement, readiness, restarts, rollback,
and cleanup. Its operations are adapters to Wails' existing manifest pipeline,
watcher and platform process tools. Build plans remain finite.

`Output` serialises lifecycle events, build reporting and attributed child
stdout/stderr. Incomplete lines and diagnostic tails are bounded to 64 KiB.
`--verbose` includes build commands and live build output; `--quiet` shows
failures. Errors retain component, phase, generation, command, exit status,
output and cause where available. The shared `internal/report` rendered-error
marker prevents duplicate CLI diagnostics. HCL builds do not install a global
reporter. Canceled/superseded work is not a build failure.

Frontend readiness requires an HTTP response below 500, with normal certificate
validation for HTTPS. Backend readiness requires a per-launch authenticated
callback from Wails' runtime-ready message; a previous process cannot satisfy
its replacement's handshake. Windowless apps use the application-started event.
The protocol is inert in ordinary launches and absent from production readiness
hooks. Android debug launches forward Go stdout/stderr into the attached logcat.

## Commands

```sh
export WAILS_EXP_USE_WAKE=1
wails3 dev
wails3 dev --verbose
wails3 dev --target ios/arm64 --device SIMULATOR_UDID
wails3 dev --target ios/arm64 --destination device --device DEVICE_ID --host MAC_LAN_IP
wails3 dev --target android/arm64 --emulator AVD_NAME
wails3 dev --target android/arm64 --device ADB_SERIAL
```

An iOS simulator must already be booted when `--device` is omitted. An explicitly
selected simulator is booted if needed. Physical iOS uses the project's Apple
development signing configuration and a Mac LAN address reachable from the
phone/tablet. The generated development bundle includes local-network metadata;
allow the device's local-network prompt when requested. Android uses loopback
and owns only its newly created `adb reverse` mappings. The emulator remains
running after the dev session; the app, log attachment and owned reverse mappings
are stopped. Mobile replacement retains only the current artifact and the
candidate needed for rollback. Package workspaces and outputs stay under
`.wails/dev/` and production artifacts remain untouched.

`--plan` prints the finite development build and effective application arguments
without connecting to devices. `dev.args` and target-specific `dev.args` provide
defaults; CLI arguments after `--` replace them, with a bare `--` clearing them.
`--appargs` is a compatibility spelling accepting one POSIX-quoted string.
Production `run.args` never supplies development defaults.

The session owns the active argument vector. Every rebuilt candidate gets the
newly resolved vector; a failed candidate retains the previous one. A runtime-only
argument edit restarts the backend while reusing its compile result and frontend.
Windows rollback reads the previous process's argument vector; mobile rollback
retains the arguments with the staged artifact. Android rejects nonempty vectors
before device discovery. iOS forwards them through its native launcher. See the
[argument reference](../../../docs/src/content/docs/reference/wails-hcl-run.md#development-application-arguments).
For HTTPS, configure the frontend server and trust its certificate; readiness
never silently disables certificate validation.

## Acceptance

Fast regression suites run with:

```sh
go test -race ./internal/dev ./internal/devruntime ./internal/report ./internal/commands ./cmd/wails3 ./internal/wake/...
go vet ./internal/dev ./internal/devruntime ./internal/report ./internal/commands ./cmd/wails3 ./internal/wake/...
```

The existing commands lifecycle/adapter/process suites exercise failed builds,
rapid edits, superseded generations, watcher configuration failures, unexpected
frontend/backend exits, startup failure, cancellation and descendant cleanup.
The mobile plan tests check development-only outputs and network configuration.

`testdata/native_acceptance.py` creates a disposable Wails fixture, using the
runtime from this checkout. First build the CLI and runtime package:

```sh
go build -o /tmp/wails3-dev ./cmd/wails3
npm --prefix internal/runtime/desktop/@wailsio/runtime ci
npm --prefix internal/runtime/desktop/@wailsio/runtime run build:code
python3 internal/dev/testdata/native_acceptance.py --cli /tmp/wails3-dev --output /tmp/wails-dev-acceptance -- --port 9351
```

On headless Linux, run the last command through `xvfb-run -a dbus-run-session --`.
On Windows, run it from an interactive console so Ctrl+Break reaches the child
process group. Mobile flags can follow `--`; use a separate, correctly signed
fixture for physical iOS. Native acceptance checks actual runtime readiness,
frontend-to-Go calls, frontend hot reload without a backend restart, a failed
build that keeps the session alive, recovery/relaunch and clean CLI shutdown.
It returns nonzero if any check fails and keeps JSON results and logs for review.

The platform acceptance matrix must cover macOS arm64, Windows amd64, Linux
amd64, iOS simulator and physical device, and an Android emulator/device.
Native GUI testing is required: compilation alone is not a passing platform
acceptance result. macOS Intel remains community-supported.
