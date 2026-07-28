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

## Fixed
<!-- Bug fixes -->
- Fix a panic on startup on 32-bit Windows. The `IDropTarget` callbacks took a `POINT` by value, which is wider than a `uintptr` on `GOARCH=386`, so `syscall.NewCallback` rejected them while `pkg/w32` was still initialising - crashing every app before `main`, whether or not it used drag-and-drop.
- Fix a panic when styling a window on 32-bit Windows. `SetWindowLongPtr`/`GetWindowLongPtr` resolved `SetWindowLongPtrW`/`GetWindowLongPtrW`, which are compiler macros rather than exported symbols in 32-bit `user32.dll`. They now select the non-`Ptr` procedures on 32-bit, matching the existing `setClassLongPtr` behaviour.

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
