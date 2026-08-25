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

## Changed
<!-- Changes in existing functionality -->

## Fixed
- Make mobile secure storage fallible and fail-closed in [PR](https://github.com/wailsapp/wails/pull/5923) by @mortenolsrud
- Run application event hooks even when no listener is registered in [PR](https://github.com/wailsapp/wails/pull/5999) by @archy-rock3t-cloud
- Correct typos in comments and localized documentation in [PR](https://github.com/wailsapp/wails/pull/6023) by @haoku123
- Remove the precompiled macOS binaries committed under `v3/examples` in [PR](https://github.com/wailsapp/wails/pull/6025) by @4RH1T3CT0R7
<!-- Bug fixes -->

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
