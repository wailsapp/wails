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
- Add Klustr to community showcase documentation in [PR](https://github.com/wailsapp/wails/pull/5536) by @SametKUM
- Add Kira to community showcase with new pages and changelog entry in [PR](https://github.com/wailsapp/wails/pull/5685) by @thiennguyen93
- Add feedback section to MCP service guide in [PR](https://github.com/wailsapp/wails/pull/5694) by @taliesin-ai

## Changed
<!-- Changes in existing functionality -->

## Fixed
<!-- Bug fixes -->
- Ensure WebKit request metadata, response completion, and body stream handling run on the GTK main thread in [PR](https://github.com/wailsapp/wails/pull/5668) by @taliesin-ai
- Fix `Menu.Update()` not rebuilding the native menu on GTK4 Linux (#5659, independently diagnosed and fixed by @puneetdixit200 in #5539)
- Fix crash enumerating macOS screens on display change by copying screen id/name strings and snapshotting the count (#5565, independently diagnosed and fixed by @x-haose in #5584)
- Fix WebView2 content shrinking then disappearing after dragging a window across mixed-DPI monitors on Windows by re-asserting the controller bounds in the `WM_DPICHANGED` handler, mirroring the un-minimise DPI resync (#5677)

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
