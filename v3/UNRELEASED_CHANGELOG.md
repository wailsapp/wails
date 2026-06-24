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
- Remove go vet from webview2 release workflow cross-compilation in [PR](https://github.com/wailsapp/wails/pull/5672) by @taliesin-ai
- Update auto-changelog OpenRouter model to google/gemini-2.5-flash-lite in [PR](https://github.com/wailsapp/wails/pull/5670) by @taliesin-ai
- Bump `webview2` to v1.0.26.
  ### Fixes
  - **Recover from transient runtime COM errors instead of exiting** (#5658, #5580). `Chromium.errorCallback` previously called `os.Exit(1)` for *any* COM error, so a recoverable hiccup after startup killed the whole application. Runtime paths (`Resize`/`GetClientRect`, `Navigate`/`NavigateToString`, `Init`, `MessageReceived`, `PutZoomFactor`, `OpenDevToolsWindow`) now log and recover. In particular, a malformed/untrusted web message in `MessageReceived` is now dropped rather than taking the process down. This addresses the mixed-DPI monitor-crossing crash class (#5544, #5650). Environment/controller-creation paths remain fatal.
  **Full diff:** https://github.com/wailsapp/wails/compare/webview2/v1.0.25...webview2/v1.0.26
<!-- Changes in existing functionality -->

## Fixed
<!-- Bug fixes -->
- Fix release-webview2 workflow to correctly handle go.sum files in [PR](https://github.com/wailsapp/wails/pull/5671) by @taliesin-ai
- Fix Linux GTK4 menu updates by clearing and rebuilding the native menu in [PR](https://github.com/wailsapp/wails/pull/5659) by @taliesin-ai

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
