# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3 alpha release.
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
- iOS: native message dialogs (UIAlertController) and open file/files/directory dialogs (UIDocumentPickerViewController); save dialogs return an explicit error
- iOS: clipboard support via UIPasteboard
- iOS: real screen metrics via UIScreen (points, pixels, scale, safe-area work area)
- iOS: device builds (`IOS_PLATFORM=device`), code-signing identity / provisioning profile / entitlements support, `.ipa` packaging, and `deploy-device` via devicectl
- iOS: configurable minimum iOS version (`ios.minIOSVersion` in build/config.yml)
- iOS: `wails3 doctor` reports Xcode and iOS SDK availability on macOS
- iOS: documentation (IOS.md and a docs-site guide)

## Changed
<!-- Changes in existing functionality -->

## Fixed
<!-- Bug fixes -->
- iOS: `GOOS=ios` compiles again (exported `events.IOS`, mobile method-name stubs) and production-tagged builds compile (build-tag fixes in pkg/application and several services)
- iOS: Go→JS events and ExecJS now work — the page no longer loads twice at startup and the `wails:runtime:ready` handshake can no longer be lost
- iOS: `ApplicationDidFinishLaunching`/`ApplicationStarted` no longer race app startup; removed the fixed 2-second startup sleep
- iOS: fixed a C-string leak on every Go→JS JavaScript execution
- iOS: `hasListeners` now reflects real listener registration
- iOS: framework debug logging is compiled out of production builds

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
