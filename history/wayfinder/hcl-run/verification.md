# HCL run implementation verification

Date: 2026-09-10. Branch: `codex/hcl-build-system`.

## Scope and evidence

`wails3 run` resolves the nearest manifest, builds with production assets and
resolved build/run/platform tags, then launches the resulting application.
Without a discoverable manifest it executes `go run .` in the working directory.
Invalid manifests fail instead of falling back. HCL execution remains experimental
and requires `WAILS_EXP_USE_WAKE=1`; the no-manifest fallback does not.

The Linux suite discovers 105 runnable example directories independently from
Go entry points. Each has a valid local manifest and a unique application ID.
The launch sweep checks compilation, launch and termination, not every button or
platform feature demonstrated by each example. Per-example outcomes are recorded
in `example-results.json`: **100 launch smoke passes, five expected unsupported-host
rejections, zero failures**. After the last compile-input correction, the final CLI
also passed Go-only smoke checks plus a warm-cache fixture that changed an imported
Go package and an embedded asset under a disabled frontend directory. The expected unsupported Linux examples are `print`,
`notch-notification`, `liquid-glass`, `systray-stress` and `ios-poc`.

Checks run from `v3/`:

```sh
go test ./internal/generator ./internal/flags ./internal/wake/... ./internal/commands ./cmd/wails3
go test -race ./internal/wake/manifest ./internal/wake/pipeline ./internal/commands ./cmd/wails3
go vet ./internal/generator ./internal/flags ./internal/wake/manifest ./internal/wake/pipeline ./internal/commands ./cmd/wails3
go test ./internal/wake/manifest -run TestExampleRunManifestsCoverRunnableDirectories -v
go build -o /tmp/wails3-run ./cmd/wails3
```

The tests cover OS/architecture precedence and ejection, empty argument overrides,
real CLI parsing, environment delivery, native exit status, signal cancellation,
production+devtools source selection, parent-module cache invalidation and binding
package loading from the manifest root. The CLI test exercises configured defaults,
replacement flags including `--help`, and an explicit empty `--`.

The timeout check used a sleeping CLI fixture with `--plan-timeout 0.05` and two
selected examples: both were reported as `plan-timeout`, the report contained both
rows, and the sweep exited nonzero rather than stopping after the first timeout.

From the repository root, reproduce the example sweep with:

```sh
python3 v3/scripts/test-examples-run.py --cli /tmp/wails3-run --output /tmp/wails-example-tests
```

The CLI also cross-compiles with `CGO_ENABLED=0` for `windows/amd64` and
`darwin/arm64`. This validates the CLI code, not native Wails application linking
or execution. `platform-plans.json` records planning checks for every example on
Linux, Windows, macOS, Android and iOS: 404 plans resolved, 20 targets were
explicitly unsupported and 101 required macOS for iOS planning. There were no
unexpected failures. iOS requires a macOS host even for planning.

Documentation validation: `cd docs && npm run build` builds 197 pages and checks
all internal links. CodeRabbit review findings and their fixes are summarised in
`review.md`.

## Android build and remaining deployment check

On this Linux host, SDK build-tools 36.0.0, NDK 29.0.14206865 and Java are installed.
The following builds the real Android example through the production run pipeline
and checks the APK manifest, DEX and x86_64 native library:

```sh
cd v3
export ANDROID_HOME="$HOME/.local/share/android-sdk"
export ANDROID_NDK_HOME="$ANDROID_HOME/ndk/29.0.14206865"
WAILS_TEST_ANDROID_RUN=1 go test ./internal/commands -run TestRunAndroidAPKAcceptance -count=1 -timeout 15m
"$ANDROID_HOME/build-tools/36.0.0/apksigner" verify --verbose examples/android/bin/android.apk
```

The APK builds and its v1/v2 signatures verify. The local emulator crashes with
exit 139 before stable ADB readiness in three configurations: SwiftShader, GPU off
with Vulkan disabled, and software CPU emulation. No Wails app was deployed to it.
This is an outstanding native acceptance check, not a passing launch result.

On an authorised connected device or working emulator:

```sh
cd v3/examples/android
export WAILS_EXP_USE_WAKE=1
adb devices -l
wails3 run --target android/arm64 --device YOUR_ADB_SERIAL
# An x86_64 emulator requires android/amd64:
wails3 run --target android/amd64 --emulator YOUR_AVD_NAME
```

Acceptance: the app installs, opens and renders its frontend; logs attach; Ctrl-C
stops the app and CLI. A mismatched architecture or unavailable device produces an
error. Application arguments and runtime environment overrides are rejected before
deployment. Different examples must install independently without replacing one
another's application ID.

## Windows and macOS

Build this branch's CLI on the native host. In PowerShell, from `v3/`:

```powershell
go build -o "$env:TEMP\wails3-run.exe" ./cmd/wails3
$env:WAILS_EXP_USE_WAKE = "1"
cd examples\plain
& "$env:TEMP\wails3-run.exe" run
```

On macOS, from `v3/`:

```sh
go build -o /tmp/wails3-run ./cmd/wails3
export WAILS_EXP_USE_WAKE=1
cd examples/single-instance-url-scheme
/tmp/wails3-run run
```

In another macOS terminal, use `open 'wails-single-url://test'`. Verify the running
example receives the URL and remains the single application instance. Also check
`custom-protocol-example` with `open 'wailsexample://test'` and the file-association
example with a `.wails` file. The launcher registers the generated bundle using
Launch Services before executing its binary. Verify icons, bundle identity,
permissions, signing/entitlements, arguments and runtime environment on macOS.

On both hosts, run the example sweep with the native CLI path and manually check
interactive features. Ctrl-C must terminate the CLI and child processes; an app
that exits normally must return its own exit status. Run twice to check cache reuse.
The native acceptance ticket remains open until these results are supplied.

## iOS simulator and signed device

On macOS with Xcode and the appropriate SDK/simulator installed:

```sh
cd v3/examples/ios
export WAILS_EXP_USE_WAKE=1
xcrun simctl list devices available
wails3 run --target ios/arm64 --device YOUR_SIMULATOR_UDID
```

For a physical device, configure `ios.signing.identity`,
`ios.signing.provisioning_profile` and any required `ios.signing.entitlements` in
the local manifest. The profile must match the example's bundle identifier and
authorise the device. Keep credentials and provisioning files out of commits.

```sh
xcrun devicectl list devices
wails3 run --target ios/arm64 --destination device --device YOUR_DEVICE_UDID
```

Acceptance: installation, frontend rendering, arguments containing spaces/flags,
environment overrides and console attachment work. Test cancellation and app
termination independently; the CLI reports launcher status, not a portable mobile
app exit code. Simulators receive `SIMCTL_CHILD_` variables; devices receive the
JSON `--environment-variables` option. The device option is also used by
[Flutter's devicectl launcher](https://github.com/flutter/flutter/blob/master/packages/flutter_tools/lib/src/ios/core_devices.dart).
This source check does not replace native testing with the installed Xcode version.

## No-manifest fallback and platform limitations

In a separate temporary Go module with no discoverable `wails.hcl`, compare
`wails3 run -- --help` with `go run . --help`. No production or devtools tags are
added automatically. `wails3 run --tags example_tag` forwards the supplied tags.
Creating an invalid `wails.hcl` must produce an error without executing the app.

Linux native builds were tested with GTK4 / WebKitGTK 6.0. GTK3 / WebKit2GTK 4.1
development libraries are absent. Legacy `gtk3` native compilation and launch,
all other native hosts, Android deployment and interactive example behaviour
remain tracked in [native acceptance](issues/02-native-acceptance.md).
