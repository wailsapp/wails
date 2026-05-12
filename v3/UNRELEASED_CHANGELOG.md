# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3 alpha release.
Add your changes under the appropriate sections below.

Guidelines:
- Follow the "Keep a Changelog" format (https://keepachangelog.com/)
- Write clear, concise descriptions of changes
- Include the impact on users when relevant
- Use present tense ("Add feature" not "Added feature")
- Reference issue/PR numbers when applicable

This file is automatically processed by the nightly release workflow.
After processing, the content will be moved to the main changelog and this file will be reset.
-->

## Added
<!-- New features, capabilities, or enhancements -->

## Changed
<!-- Changes in existing functionality -->
- **BREAKING (macOS):** Normalise the macOS coordinate system so `GetScreens`, `Position`, and `SetPosition` all use the same space — logical points, Y-down, with `(0,0)` at the top-left of the primary screen. This matches Windows, GTK and the public APIs of Electron and the web. Screens physically above the primary now report negative `Bounds.Y` (previously positive), and `Position()`/`SetPosition()` values are now in logical points instead of `points × primaryScale`. Round-tripping `Position()` → `SetPosition()` is preserved; absolute values logged from earlier alpha builds or hand-computed workarounds (e.g. multiplying by `primaryScale` or flipping Y against a screen height) will need to be updated. Resolves [#5117](https://github.com/wailsapp/wails/issues/5117).

## Fixed
<!-- Bug fixes -->
- Fix memory safety issue in GTK menu handling on Linux in [PR](https://github.com/wailsapp/wails/pull/5363) by @leaanthony
- Detect NVIDIA GPUs and disable DMA-BUF renderer on Linux in [PR](https://github.com/wailsapp/wails/pull/5295) by @leaanthony
- Fix `SetPosition` cross-screen Y conversion on macOS: use primary screen height as global reference so windows land at the correct position on monitors that are vertically offset from the primary display in [#5117](https://github.com/wailsapp/wails/issues/5117)
- Fix git PR template to point to the correct feedback URL in [PR](https://github.com/wailsapp/wails/pull/5109) by @wayneforrest
- Fix a family of Windows systray `SetMenu` crashes caused by a broken `DestroyMenu` syscall that was passing four arguments instead of one, so every call returned FALSE and freed nothing. Also release HMENU and HBITMAP handles (including those allocated at runtime via `MenuItem.SetBitmap`) on menu rebuilds, reset stale checkbox/radio maps in `Win32Menu.Update`, and drop a redundant `Update()` call in `systemtray.updateMenu` that doubled allocations. Long-running systray apps no longer leak GDI/USER objects on each menu rebuild.


## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->

## Security
<!-- Security-related changes -->

---

### Example Entries:

**Added:**
- Add support for custom window icons in application options
- Add new `SetWindowIcon()` method to runtime API (#1234)

**Changed:**
- Update minimum Go version requirement to 1.21
- Improve error messages for invalid configuration files

**Fixed:**
- Fix memory leak in event system during window close operations (#5678)
- Fix crash when using context menus on Linux with Wayland

**Security:**
- Update dependencies to address CVE-2024-12345 in third-party library
