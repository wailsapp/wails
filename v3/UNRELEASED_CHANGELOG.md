# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3 release.
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

## Fixed
<!-- Bug fixes -->
- Fix `WebviewWindow.SetMenu(nil)` panicking on Windows instead of being a no-op like on macOS, fixing [#6104](https://github.com/wailsapp/wails/issues/6104)
- Preserve changelog entries across competing pushes with retries in [PR](https://github.com/wailsapp/wails/pull/6094) by @leaanthony
- Fix `go mod vendor` failing with `pattern arm64/WebView2Loader.dll: no matching files found` on every platform, by removing the embeds that referenced binaries never shipped in the module, fixing [#5782](https://github.com/wailsapp/wails/issues/5782) and [#5376](https://github.com/wailsapp/wails/issues/5376), in [PR](https://github.com/wailsapp/wails/pull/6031) by @Grantmartin2002
- Prevent duplicate window teardown on macOS in [PR](https://github.com/wailsapp/wails/pull/6159) by @charlesbalderson
- Handle WM_CONTEXTMENU for systray menu to open on right-click in [PR](https://github.com/wailsapp/wails/pull/6162) by @wzk2025
- Fix all five `InvokeSync*` variants blocking their caller forever when the dispatched callback panics, fixing [#6107](https://github.com/wailsapp/wails/issues/6107), in [PR](https://github.com/wailsapp/wails/pull/6108) by @RALIST

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
