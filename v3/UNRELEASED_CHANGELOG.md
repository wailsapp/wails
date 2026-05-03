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
- Add option to disable Escape key exiting fullscreen on macOS in [PR](https://github.com/wailsapp/wails/pull/5307) by @leaanthony
- Add option to disable Escape key exiting fullscreen on macOS in [PR](https://github.com/wailsapp/wails/pull/5310) by @leaanthony
- Adds Pausa community showcase documentation in [PR](https://github.com/wailsapp/wails/pull/5288) by @yuseferi

## Changed
<!-- Changes in existing functionality -->
- Update sponsors SVG in [PR](https://github.com/wailsapp/wails/pull/5308) by @github-actions[bot]
- Update icon generation command to handle unsupported platforms in [PR](https://github.com/wailsapp/wails/pull/5309) by @leaanthony
- Replace boolean fullscreen API with tri-state ButtonState, implement platform bindings in [PR](https://github.com/wailsapp/wails/pull/5224) by @leaanthony

## Fixed
<!-- Bug fixes -->
- Guard WebView2 focus operations against nil controller state in [PR](https://github.com/wailsapp/wails/pull/5315) by @leaanthony
- Update GitHub Actions workflow to correctly reference PR base branch in [PR](https://github.com/wailsapp/wails/pull/5313) by @leaanthony
- Ignore *_test.go files in dev mode to prevent unnecessary rebuilds in [PR](https://github.com/wailsapp/wails/pull/5203) by @leaanthony
- Prevent Menu.Update() segfault when app is not running in [PR](https://github.com/wailsapp/wails/pull/5291) by @wucm667

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
