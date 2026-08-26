# Matching-host release verification

Type: task
Status: open
Blocked by: none
Label: ready-for-human

## Question

Run the manifest build-system release matrix on the native hosts and with the
credentials that are unavailable on the Linux audit box. This is verification
of the implemented contract, not additional first-release design work.

## Acceptance criteria

- Linux arm64 on native hardware: binary plus every supported Linux package
  format, receipt verification, and zero-work reruns for reusable operations.
- Windows amd64 and arm64: build, NSIS, MSIX, receipt verification, and a
  zero-work rerun for every reusable operation.
- macOS amd64, arm64, and universal: binary, app, DMG, receipt verification,
  and zero-work reruns for reusable operations.
- iOS simulator and device: native build/package through Xcode; exercise IPA
  signing with release credentials and prove signing and downstream operations
  are never reported as cache hits.
- Android arm64, amd64, and universal: SDK/NDK build and production AAB output;
  verify the development APK path separately and confirm it is not selectable
  as a production format.
- Credentialed Windows, macOS, Linux package, iOS, and Android signing: verify
  the configured identity reaches the native tool and that unsigned or invalid
  credentials fail without publishing a stale artifact.
- On every host, archive the JSON plan, command log, artifact digests, receipt,
  cold timing, and warm timing. Run the seven-sample performance gate on the
  controlled Linux release runner.

Use `v3/scripts/verify-manifest-build-system.go` as the primary runner and the
matrix in `performance-acceptance.md` as the authoritative contract. Any host
failure becomes a narrowly scoped implementation ticket with the captured
evidence attached here. Follow
[`../multi-platform-testing-handoff.md`](../multi-platform-testing-handoff.md)
for host preparation, commands, evidence requirements, and the definition of
done.

## Comments

### 2026-08-22 — Linux audit hand-off

Linux amd64 binary, DEB, RPM, Arch, and AppImage production acceptance passed,
including receipt/digest checks and zero-work warm reruns. Structural planning
and cross-compilation passed for Windows, macOS, Linux arm64, and Android, but
those checks are not substitutes for matching-host native packaging, SDK, and
credentialed signing runs.

### 2026-08-23 — Local handoff preparation

The affected command and Wake suites pass normally and under the race
detector, and scoped vet is clean. Ten Plan/Dev tests were made independent of
a host-installed npm executable after the current Codex runtime exposed the
hidden dependency. The whole v3 sweep reaches only pre-existing
desktop-environment failures in X11 shortcut synthesis under Wayland and
file-manager launching in the test context.

The final combined seven-sample badge run measured 103.913ms median with 1.66%
MAD, zero executed and six cached Nodes in every sample, and byte-identical
artifacts. It passes semantic, artifact, variance, and relative-baseline gates.

### 2026-08-26 — Local work resumed and gate revised

The accepted warm no-op ceiling is now 150ms. Complete every Linux and Android
task available on the current Linux machine before handing off Windows, macOS,
and iOS acceptance. Android work now includes emulator/device deployment and a
polished interactive selection experience. Build/stage hooks and improved
configuration diagnostics are separately tracked follow-up capabilities and
must not be inferred from illustrative syntax alone.

### 2026-08-26 — Linux matrix completed on the audit host

The copied badge project passed native linux/amd64 binary, DEB, RPM, Arch and
AppImage acceptance. After the Podman and Android-command changes, the latest
seven-sample no-op result is 99.870ms median with 4.43% MAD, zero executed and
six cached Nodes in every sample, and stable artifact bytes, passing the
revised 150ms ceiling.

Linux arm64 cross acceptance also passed with Podman and the arm64 variant of
`ghcr.io/wailsapp/wails-cross:latest`: the produced ELF is AArch64 and DEB,
RPM and Arch packages pass receipt and zero-work verification. The cold QEMU
compile took 5m03s; warm binary and package reruns took 1.28s and 225ms. The
implementation now resolves Docker or Podman behind the container toolchain,
uses SELinux-safe Podman mounts, selects the target image platform, and uses
the image-native compiler for Linux. Native arm64 hardware remains a separate
matching-host launch/install row.

### 2026-08-26 — Android host preparation

The checksum-verified official Android command-line tools are installed
user-locally at `/home/lea/.local/share/android-sdk`, alongside Java 21, system
`adb`, and working KVM access. `wails3 android devices` runs against the real
ADB adapter and reports the empty device set correctly. Installing platform 35,
build tools, NDK r29, the emulator and an x86_64 system image—and therefore
running the real APK/AAB and deployment matrix—requires the user to accept the
Google Android SDK licence first.

### 2026-08-26 — Android SDK and emulator acceptance complete

After explicit licence acceptance, API 36, Build Tools 36.0.0, NDK r29,
Emulator 37.1.11, platform tools 37.0.1, and an API 36 Google APIs x86_64 image
were installed. The embedded Android project now uses AGP 9.0.1 with Gradle
9.2.1, a supported pairing for API 36, and its Groovy DSL uses assignment syntax
compatible with Gradle 10. The real verifier passed production AAB builds, artifact
receipts and zero-work reruns for android/amd64, android/arm64 and
android/universal. Final post-upgrade builds were 7.0s, 6.8s and 6.6s;
zero-work reruns were 111.8ms, 112.4ms and 102.7ms respectively.

The first real run exposed stale `.assets-stage-*` paths in generated Android
and iOS Go overlays. A regression test now requires every replacement path to
refer to a published file, and manifest generation rebases staged paths before
atomic publication. A second real run exposed ABI preference inversion on an
x86_64 emulator advertising ARM translation; ABI resolution now follows the
device-provided preference order.

The AGP migration also exposed that generated metadata substitution recognised
only legacy Groovy spacing. It now accepts both legacy and assignment syntax,
with regression coverage for application ID, version, build number and minimum
SDK. The final APK reports `com.mycompany.myproduct`, version `0.0.1`, API 36,
and minimum API 21.

`wails3 android run --emulator wails-hcl-api36` then built an x86_64 development
APK, installed it, and launched the configured application. ADB confirmed the
package, top-resumed Wails activity, running process, API 36 target, and x86-64
NDK r29 library. A clean regeneration after the AGP upgrade repeated both the
AAB build and APK deployment successfully. Physical-device deployment remains
outstanding.
