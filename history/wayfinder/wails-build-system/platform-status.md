# HCL build system: platform status and next steps

Updated: 29 August 2026 (Australia/Sydney)

Branch: `codex/hcl-build-system`

Status reviewed at: `acdd11858`

This is the implementation handoff for the experimental HCL build system. It
records what is implemented, what has been proved on real platforms, and the
remaining work required before broader integration.

The project owner has confirmed that a WEP is **not** required for this
experiment. Do not block implementation or platform acceptance on a WEP.

## System-wide implementation

### Done

- A single root `wails.hcl` is the explicit opt-in. When it exists, Wails build
  and development commands ignore Taskfiles completely.
- Typed HCL configuration covers project identity, frontend commands,
  development settings, build settings, targets, profiles, package formats,
  signing references, file associations, protocols and user-owned assets.
- The setup wizard builds typed internal state and writes a validated HCL
  manifest. Combination tests cover the configurable wizard state.
- `wails3 build [profile]`, `wails3 dev`, compatibility package/sign commands,
  multi-target planning, text plans and JSON plans use the native pipeline.
- Wake provides an immutable typed graph, shared work, parallel execution,
  cancellation, local content-addressed caching, restoration and artifact
  receipts.
- `wails3 config check` and build-time validation preserve HCL source ownership
  through semantic, planning, signing and host-capability errors. Diagnostics
  include line/column carets, hints and environment-value redaction.
- `wails3 migrate` produces an inactive, reviewable draft and reports custom
  Taskfile logic instead of guessing. Activation is explicit and non-destructive.
- `wails3 eject` writes a complete standalone `wails.ejected.hcl`.
- Generated platform state stays under `.wails/`; project assets and package
  templates remain user-owned and are preserved byte-for-byte.
- Six bounded lifecycle hooks are implemented: `before_build`, `after_build`,
  `before_package`, `after_package`, `before_sign` and `after_sign`.
- Hooks execute project-owned scripts directly. They receive a protected JSON
  context file through `WAILS_HOOK_CONTEXT_FILE`, support cancellation and
  process-group cleanup, and are non-cacheable unless complete inputs and
  outputs are declared.
- The internal stages and exact hook boundaries are documented in
  [`support-report.html`](support-report.html).
- The accepted warm no-op Linux results are 99.870 ms before hooks and
  105.388 ms with the unchanged no-hook graph after hook support, both below
  the 150 ms gate.

### Still needed across the system

- Run the native-host acceptance matrix below and fix every defect it exposes.
- Verify credentialed signing, notarisation and provisioning with real secrets
  held outside the repository.
- Preserve evidence for each platform row: host/tool versions, exact commit,
  plan JSON, cold and warm logs, artifact hashes, receipt, native signature
  verification, installation and launch results.
- Rebase or merge the experiment onto the intended integration branch after
  platform defects are closed, then rerun affected acceptance rows.

## Platform summary

| Platform | Implementation | Real acceptance | Primary remaining work |
| --- | --- | --- | --- |
| Linux | Implemented | amd64 complete; arm64 cross-build proved | Native arm64 install/launch and credentialed signing |
| Android | Implemented | AAB matrix and x86_64 emulator complete | Physical device and release signing acceptance |
| Windows | Implemented, cross-compiled | Not run on native Windows | Native packages, install/launch, cache and signing |
| macOS | Implemented, cross-compiled where possible | Not run on native macOS | Native app/DMG, universal, launch, signing/notarisation |
| iOS | Build/package path implemented | Not run through Xcode | Deployment workflow plus simulator/device acceptance |

## Linux

### Done

- Native `linux/amd64` production binary builds pass.
- DEB, RPM, Arch Linux and AppImage production packaging passes on the Linux
  audit host with valid receipts and zero-work warm reruns.
- The default template and all nine built-in web templates build successfully.
- The real development lifecycle, frontend/backend restart behaviour and
  cancellation pass locally.
- Podman-backed `linux/arm64` cross-builds produce an AArch64 ELF. DEB, RPM and
  Arch packages pass receipt and zero-work verification using the arm64
  variant of the Wails cross image.
- Missing and modified published outputs are restored or regenerated with
  byte-stable results.

### Next

- Run the arm64 matrix on native Linux arm64 hardware.
- Install and launch the native arm64 binary and every supported package,
  including AppImage where supported by the host.
- Smoke-test installation and launch of the amd64 package formats on disposable
  target machines, including a frontend-to-Go binding call.
- Exercise configured package signing for every supported signed Linux format
  and verify signatures with the native package tools.
- Preserve native arm64 cold/warm timings and artifact evidence. The 150 ms
  Linux no-op benchmark only needs rerunning if a later fix changes planning,
  caching or execution.

## Android

### Done

- Real production AAB builds pass for `android/amd64`, `android/arm64` and
  `android/universal`, including receipts and zero-work reruns.
- Universal AAB generation includes both required native ABIs.
- Generated projects use API 36, AGP 9.0.1, Gradle 9.2.1 and NDK r29 on the
  accepted Linux host.
- `wails3 android devices` discovers ADB targets and reports unavailable,
  offline and unauthorised states.
- `wails3 android run [profile]` supports deterministic `--device`,
  `--emulator`, `--apk` and `--no-launch` operation plus an interactive terminal
  picker.
- The command starts or reuses an AVD, waits for Android and Package Manager,
  selects a compatible ABI, builds a development APK, installs it and launches
  the configured activity.
- Attached sessions support PID-filtered logs, process restarts, Ctrl+C,
  application exit, emulator loss, failed-start cleanup and explicit
  `--stop-emulator` cleanup.
- A real API 36 x86_64 AVD run installed and launched the default Wails app with
  the correct application ID, metadata and x86_64 native library.
- Production APK is intentionally not a supported output. APK remains an
  internal development/deployment artifact; AAB is the production format.

### Next

- Connect and authorise a physical Android device and run
  `wails3 android run --device SERIAL --logs`.
- Confirm ABI selection, install, launch, frontend-to-Go binding behaviour,
  application-scoped logs, clean Ctrl+C, app-exit detection and a useful error
  when the cable is disconnected.
- Confirm `--stop-emulator` is rejected for a physical device.
- Exercise real offline and unauthorised physical-device states where safe.
- Build and verify a credentialed release AAB, including its signing certificate
  and bundle structure.

## Windows

### Done

- Windows target resolution, planning, asset generation, compilation,
  packaging, publication and receipt code are implemented.
- Focused tests cover Windows paths and hook-context file permissions.
- The production code cross-compiles for Windows amd64 from the Linux audit
  host.

### Next

- Build the branch CLI on native Windows and run command/Wake tests, race tests
  where supported, and scoped vet.
- Build `windows/amd64` and `windows/arm64` outputs on suitable native hosts.
- Produce NSIS and MSIX packages for both targets.
- Verify cold builds, zero-work reusable reruns, receipts and byte-identical
  reusable outputs.
- Install and launch the executable, NSIS installer and MSIX in disposable
  Windows environments; exercise a frontend-to-Go binding call.
- Run credentialed Authenticode signing and validate the executable, NSIS and
  MSIX signatures with native Windows tools.
- Confirm invalid credentials fail without publishing a stale artifact.

## macOS

### Done

- macOS target resolution, universal target planning, asset generation, app
  assembly, DMG packaging, signing/notarisation nodes and receipts are
  implemented.
- Production Go code cross-compiles for macOS arm64 from the Linux audit host.
- Universal-binary merge and cache behaviour have focused automated coverage.

### Next

- Build the branch CLI on native macOS and run command/Wake tests, race tests
  and scoped vet.
- Build `darwin/amd64`, `darwin/arm64` and `darwin/universal` binaries and app
  bundles.
- Produce and mount DMGs for every applicable target.
- Verify cold builds, zero-work reusable reruns, receipts and byte-identical
  reusable outputs.
- Launch each app and smoke-test a frontend-to-Go binding call.
- Exercise real code signing and notarisation, validate with `codesign` and the
  native notarisation tools, and confirm signing is never reported as cached.
- Confirm invalid or unavailable credentials cannot publish stale artifacts.

## iOS

### Done

- iOS target/profile resolution, generated Xcode state, app assembly, device
  IPA packaging, signing nodes, publication and receipts are implemented.
- Generated-path publication was fixed so Android and iOS overlays no longer
  retain stale staging-directory references.
- Structural and focused tests cover the pipeline without claiming native
  Xcode acceptance.

### Next

- Implement the iOS simulator/device deployment experience. Reuse the tested
  deployment-domain shape from Android where appropriate: discovery, explicit
  non-interactive selection, interactive selection, launch, logs, cancellation,
  target loss and safe cleanup.
- Run the simulator build through native Xcode, install it, launch it and
  exercise a frontend-to-Go binding call.
- Build a real device IPA with a provisioning profile and release credentials.
- Install or validate the device build, inspect its signature and provisioning
  profile, and prove signing and downstream publication are never cache hits.
- Cover unavailable simulators/devices, multiple targets, untrusted devices,
  startup failures and target loss with focused tests and native evidence.

## Recommended execution order

1. **macOS/iOS host:** implement and accept iOS deployment while running the
   macOS desktop, universal, DMG, signing and notarisation matrix.
2. **Windows hosts:** run amd64 and arm64 binary/package/install/signing rows.
3. **Physical Android device:** complete the existing Android deployment row
   and release-signing evidence.
4. **Native Linux arm64:** complete package installation, launch and signing.
5. Fix each discovered defect on `codex/hcl-build-system`, add a focused
   regression test, push it, and rerun every affected row.
6. Update this file and
   [`multi-platform-testing-handoff.md`](multi-platform-testing-handoff.md) with
   the final evidence, then integrate the experiment into the intended release
   branch.

## Platform definition of done

A platform is complete only when:

- tests and scoped static analysis pass using a CLI built from the exact branch
  commit under test;
- every supported target and package format is built with its native compiler
  or SDK;
- a second reusable run executes zero nodes, reports cache hits, reproduces
  stable artifact bytes and verifies the receipt;
- packages install or mount, applications launch, and a frontend-to-Go binding
  call succeeds;
- configured signing is validated with native tools and signing work executes
  on every run;
- failures do not publish stale outputs;
- complete redacted evidence is attached to the exact commit; and
- every discovered defect has a regression test and a successful rerun.

Detailed host preparation, commands and evidence requirements are in
[`multi-platform-testing-handoff.md`](multi-platform-testing-handoff.md).

## Deliberately deferred

These are not blockers for the current experiment:

- custom stages or user-defined dependency graphs;
- inline shell commands in HCL;
- typed `wails3 tool` calls from HCL;
- HCL expressions, functions or includes;
- remote caching; and
- automatic translation of arbitrary customised Taskfile logic.
