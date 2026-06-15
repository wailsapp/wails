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

## Changed
<!-- Changes in existing functionality -->
- Move iOS and Android native features onto platform managers: call them via `application.IOS.*` and `application.Android.*` (e.g. `application.IOS.Haptic("medium")`, `application.Android.Share(payload)`) instead of the old `application.IOS*`/`application.Android*` free functions (#5602)
- Rename mobile bridge events: cross-platform events now use the `common:*` prefix (e.g. `common:haptic`, `common:location`) and platform-exclusive events use `ios:*` / `android:*` (e.g. `ios:backgroundTask`, `android:foregroundService`); the `native:*` prefix is no longer used (#5602)

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
