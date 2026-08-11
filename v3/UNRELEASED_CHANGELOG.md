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
- Add macOS autoplay preference to disable user action requirement for media playback in [PR](https://github.com/wailsapp/wails/pull/5512) by @Eyalm321
- Move mailbox changelog entry to Unreleased in [PR](https://github.com/wailsapp/wails/pull/5935) by @leaanthony

## Changed
<!-- Changes in existing functionality -->
- macOS zoom animation uses CADisplayLink or NSTimer for smoother performance in [PR](https://github.com/wailsapp/wails/pull/5945) by @savely-krasovsky

## Fixed
<!-- Bug fixes -->
- Configure iOS Xcode project to retain inherited linker flags and add -ObjC in [PR](https://github.com/wailsapp/wails/pull/5915) by @mortenolsrud
- Fix excessive TCP connection churn in the `wails3 dev` asset proxy on large frontends, which could exhaust the host's ephemeral ports and make unrelated processes fail with `EADDRNOTAVAIL`
- Queue per-window event JavaScript for ordered dispatch and backpressure in [PR](https://github.com/wailsapp/wails/pull/5934) by @leaanthony

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
