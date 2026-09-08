# Additional platform fixes and native handoff

8 September 2026, follow-up to [the first fixes](fixes.md), on `codex/hcl-ga-fixes`.

## Additional fixes

1. **Windows development process ownership.** A non-inherited Job Object owns the wrapper and its descendants. Start the wrapper suspended, assign the job, then resume its primary thread. Closing the job after wrapper exit cleans descendants independently of the expired parent PID. Startup failures kill/reap the suspended process and close handles. Process-tree regressions now compile on Windows as well as Unix. This implementation requires the native Windows tests below before release acceptance is closed.
2. **Windows signing credentials.** Resolve the explicitly named password environment variable at signing time and pass its value to the signer. Missing credentials fail before invoking signing, instead of falling back to unrelated keychain state. The real HCL-to-signing regression failed before the fix and passes afterward. This test now runs on Windows too; real protected PFX verification remains outstanding.
3. **Android launch identity.** Both fresh-build and existing-APK deployment now use `android.application_id`, falling back to `project.identifier`. The emulator previously installed the APK correctly and then failed to launch the desktop ID. The same failure was reproduced in a manifest-loading regression; both paths now launch successfully.
4. **Windows test portability.** Plan-tool fixtures previously lived in a Unix-only file, preventing the entire commands test package from compiling on Windows. Shared native executable fixtures remove that compilation failure.

The Windows ownership approach follows Microsoft's [Job Object lifetime rules](https://learn.microsoft.com/en-us/windows/win32/procthread/job-objects) and [suspended thread startup](https://learn.microsoft.com/en-us/windows/win32/procthread/suspending-thread-execution). Go closes the initial thread handle, so the implementation obtains it through the documented thread-enumeration API before resuming.

## Android verification completed on Linux

Toolchain: Linux amd64, Android SDK build-tools 36.0.0, NDK 29.0.14206865, API 36 x86_64 emulator, bundletool 1.18.3. The disposable project is `/tmp/wails-hcl-ga-badge-0qn7nnnk`; it uses this checkout through a Go module replacement.

- `verify-manifest-build-system.go -android`: amd64, arm64 and universal AAB generation passed, with artifact receipts and zero-work unsigned reruns.
- Actual AAB manifest inspection confirmed application ID `com.wails.gareview.badge`, version name `3.7.2`, version code `372`, minimum SDK `28`, and target SDK `35`. See [manifest evidence](android-manifest.xml). Installed APK metadata reported the same settings.
- `jarsigner -verify` confirmed the unsigned profile produces an unsigned bundle. A disposable password-protected PKCS12 key then signed the bundle successfully; verification reported `jar verified` and the intended disposable certificate. This tests the release-signing mechanism, not a production certificate or Play Store upload.
- The final repeated signed build reused seven stages and reran signing, publication and artifact verification (1.3 seconds total). Incorrect and missing passwords each failed with exit 1; neither test password appeared in CLI output.
- `android run --device emulator-5554` built, installed and launched the development APK with the non-default ID.
- Chrome DevTools connected to the live Android WebView; it rendered the badge UI and received Go time events. A request through the Android runtime bridge successfully called `DockService.GetBadge`: [response](android-live-binding.json).
- `android run --apk ... --device emulator-5554 --logs` launched the existing APK, streamed app logs and exited cleanly on SIGINT (exit 0).

The emulator crashed during startup with software graphics. It booted with `-gpu host -feature -Vulkan`; that is an emulator-host issue, not a Wails application failure. The test-created emulator and forwarding port were cleaned up afterward. Physical arm64 device testing and Play release acceptance remain outstanding.

Go race tests and vet passed for `./internal/wake/... ./internal/commands ./cmd/wails3`. Windows amd64 commands tests, Windows arm64 CLI and Darwin arm64 CLI cross-compiled. Native Windows and Apple execution is not claimed. CodeRabbit reviewed the changes; its test-cleanup finding was fixed and the historical Windows ticket wording was clarified.

## Manual native checks

Check out `codex/hcl-ga-fixes` and record `git rev-parse HEAD`, OS/architecture and native tool versions with each result. Run from `v3/` unless a command explicitly changes into the disposable application. Use a migrated disposable app with its Go module pointing to this checkout; do not test against a previously installed Wails release.

### Windows

PowerShell:

```powershell
go build -o "$env:TEMP\wails3-ga.exe" ./cmd/wails3
go test ./internal/commands -run 'TestManifestProcess|TestHCLNativeSigningSettingsReachWindowsAndDarwinArtifacts|TestAndroidDeployUsesPlatformApplicationID' -count=20
go test ./internal/wake/... ./internal/commands ./cmd/wails3
# Replace the project path with your disposable migrated application:
go run scripts/verify-manifest-build-system.go -wails "$env:TEMP\wails3-ga.exe" -project C:\path\to\app
```

The process tests must release the listening port both on normal stop and after the wrapper exits first. Also run `wails3 dev`, restart after a frontend configuration edit, then Ctrl+C; confirm there are no leftover frontend/backend processes and the same port is reusable. Repeat from a terminal and the IDE you ordinarily use to exercise nested job ownership.

For signing, configure the app's `windows.signing` certificate, identity/publisher and credential variable name, and select `sign = true` in a release profile. Set the password variable only in the local shell. Build that profile twice; unsigned packaging should be reused while signing runs again. Verify the final EXE/MSIX signature using `signtool verify /pa /v <artifact>`. Repeat with a wrong and unset password: both must fail. Use a `sign = false` profile and confirm packaging invokes no signer. Install/uninstall NSIS and MSIX artifacts and make a frontend-to-Go call after launch. Repeat on native arm64 if GA claims that target.

### macOS

```sh
go build -o /tmp/wails3-ga ./cmd/wails3
go test ./internal/commands -run 'TestManifestProcess|TestHCLNativeSigningSettingsReachWindowsAndDarwinArtifacts|TestHCLIOS|TestHCLDevelopment' -count=3
go test ./internal/wake/... ./internal/commands ./cmd/wails3
go run scripts/verify-manifest-build-system.go -wails /tmp/wails3-ga -project /path/to/app -targets darwin/arm64,darwin/amd64,darwin/universal
```

Run a production build, save its app hash, then start and stop development. The production app must remain unchanged; the development app belongs under `.wails/dev/`. Test a complete custom Info.plist and inspect the assembled file. Mount the DMG, launch the app, and exercise a Go binding. With local credentials configured, verify Developer ID signing, notarization and stapling; rerun with invalid credentials and confirm failure. Repeat launch acceptance on native Intel where supported.

### iOS, from macOS with Xcode

```sh
go run scripts/verify-manifest-build-system.go -wails /tmp/wails3-ga -project /path/to/app -ios
# Separately, with a complete device/IPA/signing profile:
go run scripts/verify-manifest-build-system.go -wails /tmp/wails3-ga -project /path/to/app -ios-device -ios-device-profile ios-device
```

Use a visibly different `ios.icon`, non-default background modes and version/build metadata. Inspect the final app's Info.plist and icon assets, install on the simulator and launch. For a device IPA, verify provisioning and signer identity, install on a physical device and call a Go binding. The verifier builds artifacts; it does not replace these install/launch checks.

Return failing command output plus the profile (redact credentials), target, tool versions and commit. The open matching-host acceptance ticket remains the authority for GA completion; these instructions do not imply native tests have passed.
