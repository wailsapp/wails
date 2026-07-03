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
- Support mapping Go time.Time to JS Date or string in bindings in [PR](https://github.com/wailsapp/wails/pull/5398) by @fbbdev

## Changed
<!-- Changes in existing functionality -->
- Replace Node-based sponsor image pipeline with Go generator in [PR](https://github.com/wailsapp/wails/pull/5719) by @taliesin-ai

## Fixed
<!-- Bug fixes -->
- Reject the U+0085 (NEXT LINE) control character in `ValidateAndSanitizeURL`, completing the whitespace coverage of the URL validator
- Recalculate DWM frame on DPI change for frameless windows in [PR](https://github.com/wailsapp/wails/pull/4785) by @leaanthony
- Fix DnD dropzone detection failing at non-100% scaling on Windows in [PR](https://github.com/wailsapp/wails/pull/4632) by @yulesxoxo
- Add explicit Objective-C memory management for Cocoa objects across Darwin dialogs, menus, tray, and notifications in [PR](https://github.com/wailsapp/wails/pull/5714) by @taliesin-ai
- Fixes Linux CGO backend bugs and system tray issues in [PR](https://github.com/wailsapp/wails/pull/5718) by @taliesin-ai

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
