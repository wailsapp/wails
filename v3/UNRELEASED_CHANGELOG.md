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
- Add desktop environment detection on linux [PR #4797](https://github.com/wailsapp/wails/pull/4797)

## Changed
<!-- Changes in existing functionality -->
- Update the documentation page for Wails v3 Asset Server by @ndianabasi

## Fixed
<!-- Bug fixes -->
- Fix crash on macOS when toggling window visibility via Hide()/Show() with ApplicationShouldTerminateAfterLastWindowClosed enabled (#4389) by @leaanthony
- Fix memory leak in context menus on macOS and Windows when repeatedly opening menus (#4012) by @leaanthony
- Fix context menu native resources not being reused on macOS, causing fresh menu creation on each display (#4012) by @leaanthony

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
- Correct all `app.Screens.GetAll()` calls in the documentation page (`features/screens/info)` to `app.Screen.GetAll()` by @ndianabasi

**Fixed:**
- Fix memory leak in event system during window close operations (#5678)
- Fix crash when using context menus on Linux with Wayland

**Security:**
- Update dependencies to address CVE-2024-12345 in third-party library
