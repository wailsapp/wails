# Device and emulator deployment experience

Type: task
Status: open
Blocked by: none
Label: ready-for-agent

## Question

Allow Wails development builds to discover, select, launch and target an
emulator or connected device, deploy the application, start it, stream useful
status, and report actionable failures through a polished TUI or GUI.

Start with Android on the current Linux machine. Inventory the installed SDK,
NDK, adb, emulator and AVD state; define a testable deployment model beneath
the interactive surface; support deterministic non-interactive automation as
well as human selection; and keep production AAB builds distinct from internal
development APK deployment. The implementation must cover device loss,
multiple devices, offline/unauthorised devices, emulator startup readiness,
installation failure, process cancellation, logs, cache interactions and safe
cleanup. Extend the same domain model to iOS simulator/device workflows during
the later native-macOS pass.

## Comments

### 2026-08-26 — Added to the resumed goal

Android implementation and acceptance are prioritised before the remaining
Windows, macOS and iOS native-host matrix.

## Answer

Android implementation is in progress on the Linux audit host. The branch now
provides `wails3 android devices` and `wails3 android run [profile]` with:

- deterministic `--device`, `--emulator`, `--apk`, and `--no-launch` paths;
- a Huh terminal picker for connected devices or configured AVDs;
- offline/unauthorised-state diagnostics and ABI-aware arm64/x86_64 selection;
- profile-configured development APK builds kept separate from production AAB;
- install/reinstall, explicit activity launch, reusable running emulators,
  bounded startup readiness, cancellation, and `.wails/android` emulator logs;
- injected adapters and focused tests for parsing, selection, build/install/
  launch sequencing, unavailable devices, and noninteractive ambiguity.

Generated Gradle state now applies manifest application ID, version, build
number, and minimum SDK settings so deployment launches the configured app.
Real SDK/NDK/AVD acceptance is pending installation of the official Android
toolchain and explicit SDK licence acceptance on this host. The current official
command-line manager is installed at
`/home/lea/.local/share/android-sdk/cmdline-tools/latest`; Java 21, the system
`adb`, and `/dev/kvm` are available. A real CLI smoke test correctly reported no
connected devices, and explicit selection of a missing serial produced an
actionable error. SDK platform 35, build tools, NDK, emulator, system image, AVD
creation, APK installation, and launch acceptance remain licence-gated. iOS
remains for the native macOS pass.
