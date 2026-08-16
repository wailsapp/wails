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
- Add Redis Viewer showcase page with screenshots and project link — see [documentation](https://v3.wails.io/community/showcase) and [documentation](https://v3.wails.io/community/showcase/redisviewer) in [PR](https://github.com/wailsapp/wails/pull/5984) by @redisviewer

## Changed
<!-- Changes in existing functionality -->
- Isolate Wails' private macOS WebKit and AppKit usage and add the
  `noprivateapis` build tag. Default builds preserve existing transparency,
  Liquid Glass grouping and inspector behaviour; public-only builds keep the
  same Go API and use documented fallbacks or no-ops where no public equivalent
  exists. [PR #5955](https://github.com/wailsapp/wails/pull/5955) by @taliesin-ai

## Fixed
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
