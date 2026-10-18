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
- Add native macOS menu section headers for macOS 14+ — see the [documentation site](https://v3.wails.io) (20 pages updated) in [PR](https://github.com/wailsapp/wails/pull/6297) by @leaanthony

## Changed
<!-- Changes in existing functionality -->

## Fixed
<!-- Bug fixes -->
- Make SetRelativePosition relative to screen work area on macOS in [PR](https://github.com/wailsapp/wails/pull/6299) by @leaanthony
- Updater helper relaunches unchanged target after parent timeout or failed replacement in [PR](https://github.com/wailsapp/wails/pull/6298) by @taliesin-ai
- Pin pnpm to version 10 for versioned docs in [PR](https://github.com/wailsapp/wails/pull/6293) by @taliesin-ai

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
