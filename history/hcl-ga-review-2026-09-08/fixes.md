# HCL build-system fixes and Linux verification

8 September 2026. Fix branch: `codex/hcl-ga-fixes`, based on `codex/hcl-build-system` at `3452def9defa5ad7f711efe1aa14a34ca93cbe41`. Changes were made in an isolated worktree; the original working directory was left untouched.

The [original review](review.md) identified 16 defects. Fifteen are fixed, and the sixteenth (descendants surviving parent exit) is fixed for Unix. Windows process ownership remains an [open follow-up](../wayfinder/hcl-ga-review-2026-09-08/issues/17-windows-process-cleanup.md). This is not full v3 GA certification.

## Changes

- Cache correctness: discover embedded Go resources, preserve binding registrations inside methods, fingerprint compiler CPU/environment settings and persistent GOENV contents. The binding cache key is versioned to invalidate older fingerprints.
- Frontends: accept, probe and fingerprint explicit custom executables; propagate development environment and restart when it changes; use the discovered project root from nested directories.
- Migration: block unsupported Taskfile environment/execution policy and non-default VITE_PORT with actionable diagnostics. These settings are no longer silently discarded; arbitrary automatic translation was not added.
- Platform assets: apply Android versions and SDK policy, honour selected iOS icons and background modes, preserve complete replacement plist bytes, and isolate development app bundles from production outputs.
- Signing: MSIX packaging does not invoke signing; Android release Gradle configuration is unsigned and cannot fall back to a debug signer. Explicit signing stages retain responsibility for signing.
- Processes: on Unix, immediately clean the owned process group when its parent exits, including descendants that otherwise retain listening ports.

Per-finding implementation and acceptance notes are in the [ticket map](../wayfinder/hcl-ga-review-2026-09-08/map.md). Permanent regression tests are included alongside the implementation.

## Verification

All commands below ran from `v3/` on Linux amd64 unless stated otherwise:

```sh
go test ./internal/wake/... ./internal/commands ./cmd/wails3 ./internal/templates ./internal/setupwizard ./internal/report/...
go test -race ./internal/wake/... ./internal/commands ./cmd/wails3 ./internal/templates ./internal/setupwizard ./internal/report/...
go vet ./internal/wake/... ./internal/commands ./cmd/wails3 ./internal/templates ./internal/setupwizard ./internal/report/...
go build -o /tmp/wails-hcl-ga-wails3 ./cmd/wails3
```

All 18 packages passed both scoped and race-enabled suites. The final race suite and vet include the GOENV-file regression fix. Windows amd64 and Darwin arm64 CLI cross-compilation also passed; this verifies compilation, not native behavior.

Real CLI checks passed:

- A temporary project used a custom Python frontend and embedded JSON. A warm build executed zero stages; changing only the embedded JSON changed the executable's output. GOAMD64 changes from v3 to v1 were reflected by `go version -m` on the resulting binary.
- The badge example was migrated and activated, then built with the native Linux stack.
- `scripts/verify-manifest-build-system.go` passed on the final CLI for Linux binary, DEB, RPM and Arch packaging. Artifact receipts were verified. The warm package run reused 12 stages, completing in 131.3 ms.
- The seven-sample no-op benchmark passed the 150 ms gate with zero executed stages and stable binary bytes. Exact samples and timing are retained in [the benchmark result](badge-noop-fixes.json).

`coderabbit --plain` was attempted as required, but the installed CLI rejects that obsolete flag. Its supported replacement, `coderabbit review --uncommitted`, completed. The final review reported no new findings and retained the known Windows process-cleanup finding.

Raw `.log` and `.jsonl` command output remains locally available in this directory and is ignored to avoid committing megabytes of subprocess output. The original regression evidence and compact benchmark result are tracked.

## Remaining GA work

- Implement and validate Windows descendant ownership after wrapper exit.
- Complete the existing [matching-host release verification](../wayfinder/wails-build-system/issues/12-matching-host-release-verification.md) matrix: native Windows/macOS/iOS, Linux arm64, package installation/launch and binding calls, credentialed signing and failure paths, and physical Android device acceptance. Linux package production here did not install or launch packages, exercise real signing credentials, or retest AppImage.
- Verify final Android APK/AAB version and SDK metadata, release signer identity, and native Apple metadata/assets. Structural generation tests are not a substitute for native tools and credentials. Complete replacement plists are preserved; arbitrary user plist semantic correctness is not established by the byte-preservation test.
- Resolve scope and remaining work for [device/emulator deployment](../wayfinder/wails-build-system/issues/14-device-emulator-deployment.md), especially iOS, and experimental hooks/config/deployment surfaces.
- Integrate these fixes into the release candidate and retain acceptance evidence tied to that commit. The original review's outstanding-work table provides the full matrix.
