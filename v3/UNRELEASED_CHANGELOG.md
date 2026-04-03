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

## Changed
<!-- Changes in existing functionality -->

## Fixed
<!-- Bug fixes -->
- Fix `wails3 doctor` reporting incorrect WebKitGTK packages on Fedora, openSUSE, Arch, and NixOS — 4.0 fallback entries have been removed since v3 requires the 4.1 API at compile time (#5071)
- Fix openSUSE webkit2gtk doctor package name (`webkit2gtk4_1-devel` → `webkit2gtk3-devel`, the correct openSUSE package name) (#5071)
- Fix `Unexpected token '<'` error when `/wails/custom.js` is missing in desktop dev mode. Added explicit 404 handler for `/wails/custom.js` and case-insensitive `Content-Type` validation in `loadOptionalScript` to prevent HTML SPA fallbacks from being injected as JavaScript. ([#5068](https://github.com/wailsapp/wails/issues/5068))

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
- Fix deadlock EventIPCTransport.DispatchWailsEvent holding RLock during InvokeSync (#5106)

**Security:**
- Update dependencies to address CVE-2024-12345 in third-party library
