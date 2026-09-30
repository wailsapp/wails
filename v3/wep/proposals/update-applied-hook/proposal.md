# Wails Enhancement Proposal (WEP)

## Post-update hook for the v3 updater

**WEP Number**: (leave blank, assigned on acceptance)  
**Status**: Draft  
**Author**: Damon Blais (@Albinogeek)  
**Created**: 2026-09-30  
**Discussion**: [#6202](https://github.com/wailsapp/wails/pull/6202) (draft implementation opened before this WEP)  
**Implementor**: Damon Blais (@Albinogeek)
**Target**: Wails v3

## Summary

Add an optional `Config.OnUpdateApplied func(previousVersion string)` to `pkg/updater`. `Updater.Init` calls it once, on the first launch after `Restart` replaced the application, with the version that was replaced. Apps use it for bookkeeping that lives outside the binary and that the updater's file swap does not touch.

## Motivation

The v3 updater replaces the executable (or the `.app` bundle on macOS) and relaunches it. Anything that recorded the old version somewhere else is now stale, and the app has no signal that it just went through an update.

The case that prompted this: on Windows, the NSIS installer that `wails3` generates writes `DisplayVersion` under the uninstall key (`wails_tools.nsh`, `WriteRegStr ... "DisplayVersion" "${INFO_PRODUCTVERSION}"`). A self-update swaps the executable but leaves that value alone, so Settings > Apps keeps showing the version the user first installed. The app can fix the value itself, but only if it knows an update was just applied.

Other uses of the same signal:

- running one-off data or settings migrations that belong to a specific upgrade;
- showing "what's new" notes once after an update, and not after a fresh install;
- reporting update success to the app's own telemetry or logs.

## Detailed Design

### API

```go
type Config struct {
	// ...existing fields...

	// OnUpdateApplied, when set, is called by Init on the first launch after
	// Restart replaced the application, with the version that was replaced.
	OnUpdateApplied func(previousVersion string)
}
```

`previousVersion` is the `Config.CurrentVersion` of the process that called `Restart`. The new version is the app's own `CurrentVersion`, so it is not passed.

### Mechanism

1. `Restart` already starts the running executable in helper mode with `WAILS_UPDATER_HELPER_*` environment variables. It adds one more, `WAILS_UPDATER_HELPER_FROM`, set to the current version. This hop is a direct child process, so the environment reaches the helper on every platform.
2. The helper reads that value before it clears the helper variables. After the swap succeeds and immediately before it relaunches the target, it writes the version to a marker file: `os.TempDir()/wails-updated-<hex>`, where `<hex>` is the first 8 bytes of the SHA-256 of the target path (the executable, or the `.app` bundle on macOS). The file is written with mode `0600`. No marker is written when the version is empty.
3. If the relaunch fails, the helper deletes the marker before restoring the backup, so a rolled-back app is not told it was updated. If the swap itself fails, the marker is never written.
4. `Init` resolves the same target path, reads the marker, deletes it, and then calls `OnUpdateApplied` if one is set. It does this after releasing the updater's lock, so the callback may call other `Updater` methods. The marker is consumed even when no callback is set, so a later build that adds a callback is not told about an old update.

The reference implementation adds about 65 lines of non-test code across `config.go`, `helper.go`, `updater.go` and `spawn.go`.

### Why a file and not an environment variable

The obvious way to carry the version into the relaunched app is another environment variable. On macOS the helper relaunches the bundle with `open -n`, which asks LaunchServices to start it. The new process inherits launchd's environment, not the helper's, so the variable is lost. A file in the user's temp directory works the same way on all three desktop platforms.

### Impact on existing functionality

The field is optional and defaults to nil. Without it, the only change in behaviour is that the helper writes and `Init` deletes a small file in the temp directory around each applied update.

## Non-Goals

- Detecting version changes made by anything other than the Wails updater (a package manager, an installer run, a manual copy). Apps that need that should compare a stored version themselves; see Alternatives.
- Updating OS registrations such as `DisplayVersion` automatically. What to rewrite, and where, is installer-specific and stays in the app.
- A "before update" hook. The old process already knows when it calls `Restart`.
- Guaranteed exactly-once delivery across crashes. The design aims for at-most-once per update in normal operation; see Failure modes.

## Platform Considerations

- **Windows**: the helper starts the new executable directly. The marker lives in `%TEMP%`. The NSIS case applies here. Per-user installs write `DisplayVersion` under `HKCU`, which the app can update. Per-machine installs write it under `HKLM`, which a non-elevated app cannot change, so the hook reports the update but cannot fix that value by itself.
- **macOS**: the target is the `.app` bundle and the relaunch goes through `open -n`, which is the reason for the file. The helper and the relaunched app normally share the per-user `$TMPDIR`. If the app was started from a shell that set a different `TMPDIR`, the helper writes the marker where the relaunched app does not look, and the callback does not fire.
- **Linux**: the helper starts the new executable directly. For AppImages the marker key depends on the target resolving to `$APPIMAGE` rather than the per-run mount path; that resolution is proposed separately in [#6201](https://github.com/wailsapp/wails/pull/6201), and self-update of AppImages does not work without it anyway.
- **Mobile**: the updater does not apply there; no change.

## Failure modes

- **Marker left behind.** If the relaunched app exits before `Init` runs, or the new version does not use the updater at all, the marker stays in the temp directory. The next `Init` for the same target path reports the update late, possibly many launches later. Temp directory cleanup (reboot on a tmpfs `/tmp`, `systemd-tmpfiles`, macOS periodic cleanup, Windows Storage Sense) removes it eventually. A possible refinement is to also store the target's modification time in the marker and ignore a marker that does not match; the reference implementation does not do this yet.
- **Multiple instances.** Only the helper-launched instance is expected to run `Init` first. If the user starts a second instance at the same moment, both may read the marker before either deletes it, and the callback runs in both. Apps that need strict once-only behaviour should make the callback idempotent, which bookkeeping such as rewriting a registry value already is.
- **Several installs.** Two copies of the same app at different paths get different markers, so updating one does not notify the other.
- **Callback panics.** The callback runs synchronously on the goroutine that calls `Init`, and a panic propagates to that caller like a panic anywhere else in startup. The marker is already deleted, so the next launch starts normally. The updater does not recover the panic, because hiding a failure in the app's own bookkeeping would be worse than surfacing it. Maintainers may prefer a recover-and-log wrapper; that is a small change.
- **Marker write fails.** The helper logs it as non-fatal and the update proceeds; the callback simply does not fire.
- **Slow callback.** It delays the return of `Init`. The documentation should say to start long work in a goroutine.

## Pros/Cons

**Pros**

- One field; no new types, no change to existing behaviour.
- Fires only for updates the Wails updater applied, and never after a fresh install or a rollback.
- The app gets the previous version without keeping its own state.
- Works through the macOS `open -n` relaunch.

**Cons**

- Relies on a file in the temp directory and on the helper and the new process agreeing on that directory.
- At-most-once delivery is best effort (see Failure modes).
- It adds a small piece of state to a code path whose main job is replacing files safely.

## Alternatives Considered

1. **Environment variable on the relaunch.** Simplest, and it would work on Windows and Linux, but LaunchServices drops it on macOS. Rejected because the hook would silently not fire on one of the three desktop platforms.
2. **Command-line flag on the relaunch** (for example `open -n App.app --args --wails-updated-from=1.2.3`). This survives `open`, but it puts an argument the app did not define into its `os.Args`. Apps that parse flags strictly with `flag.Parse` would exit with an unknown-flag error on the first launch after every update. Rejected as breaking for existing apps.
3. **The app compares a stored last-run version itself.** At startup the app reads the version it saved on the previous run, compares it with `CurrentVersion`, runs its bookkeeping if they differ, and saves the new value. This is viable today with no framework change. It also catches version changes the updater did not make, such as a reinstall or a package manager upgrade, which for some uses (migrations, what's-new notes) is the better behaviour.

   The case for a framework hook anyway is narrower than it first looks. It is worth having when the app must act only on self-updates: the NSIS installer already writes `DisplayVersion` when it runs, so only the updater's path leaves it stale. It also spares every app from choosing a place and format for its own version state, and it puts the post-update step next to the rest of the updater configuration where users will look for it. If the maintainers judge that this does not justify a file in the temp directory, a reasonable outcome is to reject this WEP and instead document the stored-version pattern, with the NSIS `DisplayVersion` example, in the updater guide. The author is happy to write that instead.
4. **Marker next to the target instead of in the temp directory.** This removes the shared-`TMPDIR` assumption on macOS, but the install directory may not be writable at relaunch (for example a read-only AppImage mount without #6201, or a per-machine Windows install), and it leaves a stray file in the app's install directory. Kept as an option if reviewers prefer it.

## Backwards Compatibility

Fully compatible. The new field is optional, `Init`'s signature is unchanged, and an old app that never sets it only sees a temp file created by the helper and deleted by `Init`. The new environment variable is set only in the helper process and is cleared, with the other helper variables, before the app is relaunched.

## Security and Privacy

The marker contains only the previous version string. It is created with mode `0600` in the user's own temp directory and its name reveals only a truncated hash of the install path. Another process running as the same user could plant a marker and make the app believe it was updated from an arbitrary version string; that user can already modify the app itself, so this adds no new capability. Apps must still treat `previousVersion` as untrusted input and not, for example, interpolate it into a shell command. Updating per-machine registry values would need elevation, which this proposal does not provide.

## Test Plan

Unit tests in the reference implementation:

- `TestRunHelperSwap_RecordsReplacedVersion`: a successful swap writes the marker with the replaced version.
- `TestRunHelperSwap_LaunchFails_NoReplacedVersion`: when the relaunch fails and the backup is restored, no marker remains.
- `TestInit_OnUpdateApplied`: `Init` calls the callback once with the recorded version, deletes the marker, and a second `Init` in a new updater does not call it again.

Not yet done, and needed before the implementation leaves draft: a manual end-to-end update on Windows (NSIS per-user install, confirming a sample callback rewrites `DisplayVersion`), macOS (`.app` relaunched through `open -n`) and Linux (plain binary, and AppImage with #6201). Only the unit tests above have been run so far, on Linux.

## Reference Implementation

[#6202](https://github.com/wailsapp/wails/pull/6202) (draft), branch `feat/v3-updater-on-update-applied`. It is kept as a draft until this WEP is decided and will be revised to match the accepted design.

## Maintenance Plan

The author will maintain the feature and respond to issues about it. The code is small and sits inside the helper swap path, so it will be reviewed along with any future change to how the helper relaunches the app; if the relaunch mechanism changes (for example if macOS relaunches stop going through `open`), the marker could be replaced by a simpler channel without changing the public API.

## Conclusion

`OnUpdateApplied` gives apps a small, explicit signal that the updater has just replaced them, which is what they need to fix installer registrations and run upgrade-specific work. The temp-file mechanism is the least intrusive way to carry that signal through every platform's relaunch. The userland alternative is real and covers a wider set of version changes, so this WEP asks the maintainers to choose between a framework hook and a documented pattern; either would close the gap that the NSIS case exposes.

This proposal and the reference implementation were written with an AI assistant (Claude). I reviewed the text and the code and ran the tests listed above.
