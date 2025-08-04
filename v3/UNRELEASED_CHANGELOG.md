# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3-alpha release.
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
- Support for dropzones with event sourcing dropped element data [@atterpac](https://github.com/atterpac) in [#4318](https://github.com/wailsapp/wails/pull/4318)
- Added `AdditionalLaunchArgs` to `WindowsWindow` options to allow for additional command line arguments to be passed to the WebView2 browser. in [PR](https://github.com/wailsapp/wails/pull/4467)
- Added Run go mod tidy automatically after wails init [@triadmoko](https://github.com/triadmoko) in [PR](https://github.com/wailsapp/wails/pull/4286)
- Windows Snapassist feature by @leaanthony in [PR](https://github.dev/wailsapp/wails/pull/4463)

## Changed
<!-- Changes in existing functionality -->

## Fixed
<!-- Bug fixes -->
- Fix Windows nil pointer dereference bug reported in [#4456](https://github.com/wailsapp/wails/issues/4456) by @leaanthony in [#4460](https://github.com/wailsapp/wails/pull/4460)
- Add support for `allowsBackForwardNavigationGestures` in macOS WKWebView to enable two-finger swipe navigation gestures (#1857)
- Fixes issue where onClick didn't work for menu items initially set as disabled by @leaanthony in [PR #4469](https://github.com/wailsapp/wails/pull/4469). Thanks to @IanVS for the initial investigation.
- Fix Vite server not being cleaned up when build fails (#4403)
- Fixed panic when closing or cancelling a `SaveFileDialog` on windows. Fixed in [PR](https://github.com/wailsapp/wails/pull/4284) by @hkhere
- Fixed HTML level drag and drop on Windows by [@mbaklor](https://github.com/mbaklor) in [#4259](https://github.com/wailsapp/wails/pull/4259)

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
