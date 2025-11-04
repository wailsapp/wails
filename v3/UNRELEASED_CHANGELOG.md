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
- Added NSIS Protocol template for Windows by @Tolfx in #4510
- Added tests for build-assets by @Tolfx in #4510

## Changed
<!-- Changes in existing functionality -->

## Fixed
- Fixed linux desktop.tmpl protocol range, by removing `<.Info.Protocol>` to `<.Protocol>` by @Tolfx in #4510
- Fixed redefinition error for liquid glass demo in [#4542](https://github.com/wailsapp/wails/pull/4542) by @Etesam913
- Fixed DnD dropzone detection failing at non-100% display scaling on Windows [#4632](https://github.com/wailsapp/wails/pull/4632) by @yulesxoxo

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
