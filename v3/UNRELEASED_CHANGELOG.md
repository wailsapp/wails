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
- Add `-tags` flag to `wails3 build` command for passing custom Go build tags (e.g., `wails3 build -tags gtk4`) (#4957)
- Add documentation for automatic enum generation in binding generator, including dedicated Enums page and sidebar navigation (#4972)

## Changed
<!-- Changes in existing functionality -->

## Fixed
<!-- Bug fixes -->
- Fix `InvisibleTitleBarHeight` being applied to all macOS windows instead of only frameless or transparent title bar windows (#4960)
- Fix window shaking/jitter when resizing from top corners with `InvisibleTitleBarHeight` enabled, by skipping drag initiation near window edges (#4960)
- Fix generation of mapped types with enum keys in JS/TS bindings (#4437) by @fbbdev

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
