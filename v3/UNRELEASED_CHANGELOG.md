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

feat(v3/mobile): add application.Mobile.Timezone()

Go's Android port hard-codes time.Local to UTC (time/zoneinfo_android.go
initLocal() is a stub with a decade-old TODO to read persist.sys.timezone),
so every wails3 Android app currently shows UTC timestamps regardless of
the device's configured timezone.

Add Timezone() to the MobileManager interface, following the existing
StoragePath()/PowerJSON() pattern: a Java bridge method backed by the
public TimeZone.getDefault() API on Android, an NSTimeZone-backed
implementation on iOS (which is unaffected by the bug but kept at parity
with the interface), and a "" no-op on desktop.

The method only reports the device's IANA timezone ID; it does not touch
time.Local itself. Applying it is left to the caller:

    if tz := application.Mobile.Timezone(); tz != "" {
        if loc, err := time.LoadLocation(tz); err == nil {
            time.Local = loc
        }
    }

This mirrors how the rest of Mobile.* is designed (query, don't mutate
global state), and avoids silently changing time.Local for apps that
don't opt in.

Touches:
- v3/pkg/application/mobile.go
- v3/pkg/application/mobile_features_android.go
- v3/pkg/application/mobile_features_ios.go
- v3/pkg/application/mobile_features_ios.h
- v3/pkg/application/mobile_features_ios.m
- v3/pkg/application/mobile_stub.go
- v3/internal/commands/build_assets/android/app/src/main/java/com/wails/app/WailsBridge.java
- docs/mpress/content/guides/mobile/mobile-api.md
- v3/UNRELEASED_CHANGELOG.md

Fixes #6189

## Changed
<!-- Changes in existing functionality -->

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
