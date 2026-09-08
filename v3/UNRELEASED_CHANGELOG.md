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
- Add the `appstore` build tag for macOS. Every private WebKit and AppKit API
  Wails calls now lives in `mac_private_api_darwin.go`; building with
  `-tags appstore` selects public equivalents or documented no-ops instead, with
  no change to the Go API. [PR #6060](https://github.com/wailsapp/wails/pull/6060) by @taliesin-ai
- Add `MacWebviewPreferences.PreferPageRenderingUpdatesNear60FPS`. Setting it to
  `application.Disabled` lets the webview render at the display's native refresh
  rate instead of WebKit's 60fps default, so `requestAnimationFrame` runs at
  120Hz on a ProMotion display. Uses a private WebKit API and is ignored under
  `-tags appstore`. Fixes [#6056](https://github.com/wailsapp/wails/issues/6056) in
  [PR #6060](https://github.com/wailsapp/wails/pull/6060) by @taliesin-ai

## Changed
<!-- Changes in existing functionality -->
- Enable the Web Inspector through the public `WKWebView.inspectable` property
  on macOS 13.3+, keeping the private `developerExtrasEnabled` preference only
  as a fallback for older systems. [PR #6060](https://github.com/wailsapp/wails/pull/6060) by @taliesin-ai

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
