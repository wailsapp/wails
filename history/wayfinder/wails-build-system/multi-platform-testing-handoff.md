# Multi-platform build-system testing handoff

This guide takes the manifest-driven build system from the completed local
Linux audit to release acceptance on matching native hosts.

## Handoff point

- Branch: `codex/hcl-build-system`
- Remote: `origin/codex/hcl-build-system`
- Minimum implementation commit: `ae193130c`
- Open acceptance ticket:
  [`issues/12-matching-host-release-verification.md`](issues/12-matching-host-release-verification.md)
- Acceptance contract:
  [`performance-acceptance.md`](performance-acceptance.md)

Always test the latest remote branch. On an existing clone:

```bash
git fetch origin
git switch codex/hcl-build-system
git pull --ff-only
git status --short --branch
```

The status must show a clean branch tracking
`origin/codex/hcl-build-system`. Keep fixes and test evidence on this branch and
push them before handing off again.

## Prepare each test host

Use a disposable native machine or VM for each host operating system. Record
the operating-system version, CPU architecture, CPU model, memory, Go version,
Node/package-manager versions, and native SDK/tool versions in the evidence.

From the repository's `v3` directory:

```bash
go test ./internal/commands ./internal/wake/... -count=1
go test -race ./internal/commands ./internal/wake/... -count=1
go vet ./internal/commands ./internal/wake/...
go build -o /tmp/wails3-hcl ./cmd/wails3
```

On Windows, build the CLI to a native temporary path instead:

```powershell
$WailsCLI = Join-Path $env:TEMP "wails3-hcl.exe"
go build -o $WailsCLI .\cmd\wails3
```

Run `/tmp/wails3-hcl doctor` (or `& $WailsCLI doctor` on Windows) and install
every tool required by that host's formats before collecting acceptance
evidence. Use the branch-built CLI throughout; an installed release CLI is not
valid evidence.

## Prepare the disposable project

Copy `v3/examples/badge` to a directory outside the checkout. Never collect
native acceptance evidence from the repository's example in place.

In the copied project:

```bash
/tmp/wails3-hcl migrate
```

Review `wails.migrated.hcl` and resolve every migration blocker. Then activate
it:

```bash
/tmp/wails3-hcl migrate --activate
```

The Windows equivalents are `& $WailsCLI migrate` and
`& $WailsCLI migrate --activate`.

Acceptance starts only after `wails.hcl` exists and the project builds without
consulting its Taskfiles. Preserve the reviewed `wails.hcl` with the evidence.

## Run the native matrix

Run the verifier from the repository's `v3` directory. `-project` points to
the disposable badge project and `-wails` points to the branch-built CLI. The
verifier builds the binary and native packages, verifies the artifact receipt,
then reruns every reusable command and requires zero executed Nodes.

| Host | Required targets | Required formats and checks |
| --- | --- | --- |
| Linux | `linux/amd64`, `linux/arm64` | Binary, DEB, RPM, Arch and AppImage |
| Windows | `windows/amd64`, `windows/arm64` | Binary, NSIS and MSIX |
| macOS | `darwin/amd64`, `darwin/arm64`, `darwin/universal` | Binary, app and DMG |
| macOS with Xcode | `ios/arm64` simulator | Simulator app |
| Credentialed macOS with Xcode | `ios/arm64` device profile | Signed device IPA |
| Android SDK/NDK host | `android/arm64`, `android/amd64`, `android/universal` | Production AAB only |

Examples:

```bash
# Linux; run each architecture on a suitable Linux host/toolchain.
go run ./scripts/verify-manifest-build-system.go \
  -wails /tmp/wails3-hcl \
  -project /path/to/disposable/badge \
  -targets linux/amd64 \
  -appimage

# Windows PowerShell.
go run .\scripts\verify-manifest-build-system.go `
  -wails $WailsCLI `
  -project C:\path\to\disposable\badge `
  -targets windows/amd64,windows/arm64

# macOS desktop matrix and iOS simulator acceptance.
go run ./scripts/verify-manifest-build-system.go \
  -wails /tmp/wails3-hcl \
  -project /path/to/disposable/badge \
  -targets darwin/amd64,darwin/arm64,darwin/universal \
  -ios

# Run separately with release credentials and a device profile.
# The -ios-device switch appends the ios/arm64 device row; -targets controls
# only the matching macOS desktop rows and therefore defaults to this Mac.
go run ./scripts/verify-manifest-build-system.go \
  -wails /tmp/wails3-hcl \
  -project /path/to/disposable/badge-device \
  -ios-device \
  -sign

# Android universal AAB on a host with its SDK and NDK configured.
go run ./scripts/verify-manifest-build-system.go \
  -wails /tmp/wails3-hcl \
  -project /path/to/disposable/badge-android \
  -android
```

The verifier's Android switch exercises the universal AAB. Run the two
architecture-specific rows directly as supplemental acceptance until the
runner supports them:

```bash
(
  cd /path/to/disposable/badge-android
  /tmp/wails3-hcl build --targets android/arm64 --formats aab
  /tmp/wails3-hcl build --targets android/arm64 --formats aab
  /tmp/wails3-hcl build --targets android/amd64 --formats aab
  /tmp/wails3-hcl build --targets android/amd64 --formats aab
)
```

For each second run, require `0 ran`, a positive cached count, a valid
`.wails/artifacts/receipt.json`, and byte-identical output. Inspect the
universal AAB and confirm it contains both required native ABIs. An APK is a
development artefact and does not satisfy production acceptance.

## Credentialed signing pass

Keep credentials outside the repository and declare only credential
references in the disposable project's `wails.hcl`. Run the applicable host
matrix again with `-sign`, plus the separate iOS device command above.

Signing and publication downstream of signing are intentionally non-reusable.
Their second execution must run again; a cache hit is a failure. Verify the
result with the native platform tool as well as the Wails receipt:

- Windows: validate Authenticode on the executable, NSIS installer and MSIX.
- macOS: run `codesign` verification; validate the DMG and app, and notarisation
  when the selected release profile requests it.
- iOS: validate the IPA signature and provisioning profile on a real device
  build.
- Android: verify the AAB signing certificate and bundle structure.
- Linux: verify the configured package signature for every signed format.

Redact credentials, passwords, private keys and signing-session tokens from
all logs before committing or sharing evidence.

## Controlled Linux performance gate

The latest uncontrolled desktop result was 103.913 ms and missed the 100 ms
absolute ceiling by 3.913 ms. It is diagnostic evidence, not the release
decision. Run the gate on the controlled Linux runner from a warmed disposable
badge project:

```bash
go build -o /tmp/wails-build-benchmark ./scripts/benchmark-manifest-build.go

/tmp/wails-build-benchmark \
  -name badge-noop-linux-amd64 \
  -samples 7 \
  -warmups 2 \
  -dir /path/to/disposable/badge \
  -artifacts bin/badge \
  -baseline ./internal/wake/benchmark/testdata/badge-noop-linux-amd64.json \
  -max-ms 100 \
  -max-regression 20 \
  -max-mad-percent 15 \
  -expect-ran 0 \
  -expect-cached 6 \
  -require-stable-artifacts \
  -output /path/to/evidence/badge-noop-linux-amd64.json \
  -- /tmp/wails3-hcl build --targets linux/amd64
```

The gate passes only when the seven-sample median is at most 100 ms, is no
more than 20% slower than the checked baseline, MAD is at most 15%, all samples
report `0 ran / 6 cached`, and artifact identities remain stable. A noisy run
above the MAD limit is rerun once on a clean controlled runner.

## Evidence for every row

Create one evidence directory per host, target and format. Preserve:

- host and toolchain inventory;
- the exact Git commit and clean status;
- reviewed `wails.hcl` with secret references redacted if necessary;
- JSON plan from the branch-built CLI, for example
  `/tmp/wails3-hcl build --plan --json`;
- complete cold and warm command logs with exit codes and timings;
- `.wails/artifacts/receipt.json` and independent artifact hashes;
- native signature-verification output where applicable;
- proof that the produced application launches and its frontend-to-Go binding
  call succeeds;
- installer/package smoke-test results on a disposable target machine;
- the controlled benchmark JSON on Linux.

Append a concise result table and links to the evidence under ticket 12. If a
row fails, preserve the failed evidence and create a narrowly scoped follow-up
ticket. Fixes stay on `codex/hcl-build-system`, receive focused regression
tests, and are pushed before that row is rerun.

## Definition of done

Multi-platform acceptance is done only when all of the following are true:

- The latest `origin/codex/hcl-build-system` passes the command and Wake tests,
  race detector and scoped vet on every available native host.
- Every target and format in the native matrix completes on its matching host
  operating system with the required target compiler or SDK; no row is
  represented only by a structural Plan.
- Every reusable build/package rerun executes zero Nodes, reports at least one
  cached Node, reproduces identical artifact bytes, and verifies its receipt.
- Linux packages install and launch; Windows NSIS/MSIX install and launch;
  macOS apps and DMGs mount or install and launch; the iOS simulator app runs;
  the device IPA validates on a device build; and the Android AAB validates
  with the required ABIs. The frontend-to-Go binding smoke test passes in each
  launched desktop/mobile application.
- Every credentialed format carries a valid native signature. Signing and its
  non-reusable downstream Nodes execute on every signing run.
- The controlled Linux seven-sample gate satisfies every latency, variance,
  work-count and artifact-stability budget above.
- Every row has complete, redacted evidence tied to the exact commit tested.
- Every discovered defect is fixed and covered by a regression test; there are
  no skipped, waived or merely documented failing rows.
- Ticket 12 records the final matrix and is resolved, all changes are committed
  and pushed, and the working tree is clean and up to date with
  `origin/codex/hcl-build-system`.
