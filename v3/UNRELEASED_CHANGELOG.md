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
- Migrate Vite dev server port configuration to environment variables in [PR](https://github.com/wailsapp/wails/pull/5365) by @leaanthony
- Configure Vite dev server to bind to 127.0.0.1 in all templates in [PR](https://github.com/wailsapp/wails/pull/5361) by @leaanthony
- Update sponsors SVG in [PR](https://github.com/wailsapp/wails/pull/5384) by @github-actions[bot]

## Fixed
<!-- Bug fixes -->
- Fix stale state in macOS menus by applying the menu item mutators (`setMenuItemChecked()`, `setMenuItemLabel()`, `setMenuItemDisabled()`, `setMenuItemHidden()`, `setMenuItemTooltip()`) synchronously on the main thread, eliminating the `dispatch_async` race that left menus rendering the previous state when reopened quickly (#5002)
- Ignore *_test.go files in dev mode to prevent unnecessary rebuilds in [PR](https://github.com/wailsapp/wails/pull/5203) by @leaanthony
- Prevent Menu.Update() segfault when app is not running in [PR](https://github.com/wailsapp/wails/pull/5291) by @wucm667
- Use lastSizeWParam to gate menubar redraws on Windows in [PR](https://github.com/wailsapp/wails/pull/5382) by @taliesin-ai

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
