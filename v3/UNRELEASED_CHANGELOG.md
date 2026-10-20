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
<!-- Bug fixes -->
- Fix Windows webview not regaining focus on window activation, and Tab being trapped inside the webview instead of moving to the native window and back (#3783)
- Fix `internal/webview2/pkg/webview2` COM vtable structs (`ICoreWebView2_2`..`_27` and other derived interfaces) omitting the inherited method chain, which made calls land on the wrong vtable slot (#6129)
- Fix Windows question dialogs ignoring custom button labels and never firing button callbacks; custom buttons are now shown with a TaskDialog (falling back to MessageBox with positional callback mapping when comctl32 v6 is unavailable) (#5794)
- Fix `wails3 dev` hanging for 60s against Vite: the automatic frontend readiness gate now uses a TCP probe instead of an HTTP GET that sent no `Accept` header and received a 404 ([#6209](https://github.com/wailsapp/wails/issues/6209))
- Keep systray menus working when updated while open in [PR](https://github.com/wailsapp/wails/pull/6301) by @taliesin-ai
- Handle transient WebView2 COM errors gracefully in runtime operations in [PR](https://github.com/wailsapp/wails/pull/6302) by @leaanthony

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
