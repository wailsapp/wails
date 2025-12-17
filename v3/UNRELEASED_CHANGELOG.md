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
- Add `Window.Print()` method to JavaScript runtime for triggering print dialog from frontend (#4290) by @leaanthony
- Add `DisableMenu` option to `WindowsWindow` to disable the menu bar on Windows (#4273)
<!-- New features, capabilities, or enhancements -->

## Changed
<!-- Changes in existing functionality -->

## Fixed
- Fix macOS dock icon click not showing hidden windows when app started with `Hidden: true` (#4583) by @leaanthony
- Fix macOS print dialog not opening due to incorrect window pointer type in CGO call (#4290) by @leaanthony
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
