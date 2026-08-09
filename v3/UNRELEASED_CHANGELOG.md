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
- Implement macOS Dock bounce for window flashing in [PR](https://github.com/wailsapp/wails/pull/5921) by @julianstorer

## Changed
<!-- Changes in existing functionality -->

## Fixed
<!-- Bug fixes -->
- Asset server preserves content-type sniffer errors and unwritten prefixes during flushing in [PR](https://github.com/wailsapp/wails/pull/5931) by @leaanthony
- Prevent macOS applications from crashing when replacing the application menu from a Wails callback
- Fix unreadable native menus on Windows 10 1809 / Windows Server 2019 (build 17763). The dark-mode uxtheme exports were gated on build 18334, so the app-level dark-mode opt-in never ran on those hosts: the menu background was painted dark but Windows kept drawing menu text in the light theme, leaving dark text on a dark background. The ordinals exist from 17763, so the gate now matches.
- Fix `w32.GetStockObject` calling `GetDeviceCaps` instead of `GetStockObject`, which made it return 0 for every stock object.
- Improve WebView2 bootstrapper download error handling and reporting in [PR](https://github.com/wailsapp/wails/pull/5924) by @jannskiee

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
