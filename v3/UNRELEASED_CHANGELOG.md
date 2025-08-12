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
- Types Events by @fbbdev in [#4161](https://github.com/wailsapp/wails/pull/4161)
- Browser URL Sanitisation by @leaanthony in [#4500](https://github.dev/wailsapp/wails/pull/4500). Based on [#4484](https://github.com/wailsapp/wails/pull/4484) by @APShenkin.

## Changed
<!-- Changes in existing functionality -->
- When emitting a custom event with zero or one data argument, the data value will be assigned directly to the Data field without wrapping it in a slice by [@fbbdev](https://github.com/fbbdev) in [#4161](https://github.com/wailsapp/wails/pull/4161)

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
