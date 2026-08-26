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

The Android first slice is implemented and accepted on the Linux audit host.
The branch provides `wails3 android devices` and
`wails3 android run [profile]` with:

- deterministic `--device`, `--emulator`, `--apk`, and `--no-launch` paths;
- opt-in `--logs` application log streaming and `--stop-emulator` cleanup;
- a Huh terminal picker for connected devices or configured AVDs;
- offline/unauthorised-state diagnostics and ABI-aware arm64/x86_64 selection;
- profile-configured development APK builds kept separate from production AAB;
- install/reinstall, explicit activity launch, reusable running emulators,
  bounded startup readiness that waits for Android boot completion and Package
  Manager, cancellation, failed-start cleanup, and `.wails/android` emulator
  logs;
- PID-filtered logcat attachment that follows application restarts, exits
  cleanly on Ctrl+C or app shutdown, and reports disconnect/offline transitions;
- injected adapters and focused tests for parsing, selection, build/install/
  launch sequencing, unavailable devices, and noninteractive ambiguity.

Generated Gradle state now applies manifest application ID, version, build
number, and minimum SDK settings so deployment launches the configured app.
The accepted host now has API 36, Build Tools 36.0.0, NDK r29, Emulator
37.1.11, platform tools 37.0.1, and the API 36 Google APIs x86_64 image. A real
`wails-hcl-api36` AVD run selected x86_64 ahead of its advertised ARM
translation ABI, built a 15,149,530-byte debug APK, installed
`com.mycompany.myproduct`, and resumed `com.wails.app.MainActivity`. The
installed native library is an x86-64 Android 21 ELF built by NDK r29; a
captured screen showed the live Wails badge application and logcat contained
no fatal exception.

On 2026-08-27, a second real AVD acceptance streamed the Wails application's
PID-filtered logs. Killing the emulator ended the command promptly with
`Android target emulator-5554 became offline`. A separate attached run ended
cleanly on Ctrl+C and `--stop-emulator` removed the ADB target. New-emulator
startup cancellation also has process-level cleanup and reaping coverage.

Remaining work is physical-device acceptance and the iOS simulator/device
extension on native macOS.
