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
- Add why-wails.mdx documentation in multiple languages in [PR](https://github.com/wailsapp/wails/pull/5739) by @taliesin-ai
- Support mapping Go time.Time to JS Date or string in bindings in [PR](https://github.com/wailsapp/wails/pull/5398) by @fbbdev
- Add Android App Bundle (AAB) packaging tasks (`bundle`, `bundle:fat`, `assemble:aab`, `assemble:aab:release`) for Play Store submission — APK tasks remain for local/emulator testing in [PR](https://github.com/wailsapp/wails/pull/5728) by @mortenolsrud (fixes [#5726](https://github.com/wailsapp/wails/issues/5726))
- Add Android physical device task targets and resume camera/location permissions in [PR](https://github.com/wailsapp/wails/pull/5735) by @taliesin-ai

## Changed
<!-- Changes in existing functionality -->

## Fixed
<!-- Bug fixes -->
- Update French documentation for SvelteKit and options in [PR](https://github.com/wailsapp/wails/pull/5744) by @leaanthony
- Re-enable WebView2 monitor-scale detection and gate host resync on DPI change in [PR](https://github.com/wailsapp/wails/pull/5734) by @taliesin-ai
- Fix SIGSEGV on macOS screen enumeration during display changes in [PR](https://github.com/wailsapp/wails/pull/5516) by @flofreud

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
