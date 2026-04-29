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
- Grant actions: write permission to trigger-release job in [PR](https://github.com/wailsapp/wails/pull/5270) by @leaanthony

## Changed
<!-- Changes in existing functionality -->
- Improve workflow efficiency by adding path filters and removing dead workflows in [PR](https://github.com/wailsapp/wails/pull/5280) by @leaanthony
- Update documentation to reference master branch for example links in [PR](https://github.com/wailsapp/wails/pull/5274) by @leaanthony
- Update documentation and examples for v3 in [PR](https://github.com/wailsapp/wails/pull/5272) by @leaanthony

## Fixed
<!-- Bug fixes -->
- Rewrite unreleased changelog trigger workflow in [PR](https://github.com/wailsapp/wails/pull/5281) by @leaanthony

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->
- Remove shell test scripts for various testing purposes in [PR](https://github.com/wailsapp/wails/pull/5267) by @leaanthony
- Delete v3-alpha documentation deployment workflow and CNAME record in [PR](https://github.com/wailsapp/wails/pull/5266) by @leaanthony

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
