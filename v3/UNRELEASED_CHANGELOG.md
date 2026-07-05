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
- Generate the docs changelog page in a site-friendly format: the page now starts with the latest release and category headers nest under version headers in the table of contents
- Fold the `webview2` binding into the v3 module as `v3/internal/webview2`, removing the standalone module, its nightly release/sync workflows, and the go.mod version dance (v3 is its only consumer) in [PR](https://github.com/wailsapp/wails/pull/5711) by @taliesin-ai

## Fixed
<!-- Bug fixes -->
- Move WebView2 monitor-scale detection and DPI-change host resync fix to Unreleased section in [PR](https://github.com/wailsapp/wails/pull/5750) by @taliesin-ai
- Update WebView2 COM marshaling for float64 and BOOL parameters in [PR](https://github.com/wailsapp/wails/pull/5741) by @wayneforrest
- Prevent panic and nil dereference in Windows system tray icon updates and destruction in [PR](https://github.com/wailsapp/wails/pull/5703) by @wayneforrest
- Fixes hidden windows not re-hiding correctly on Windows in [PR](https://github.com/wailsapp/wails/pull/5743) by @wayneforrest
- Synchronize WebView2 controller visibility with window minimize/maximize/restore in [PR](https://github.com/wailsapp/wails/pull/5742) by @wayneforrest

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
