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
- Add `XDG_SESSION_TYPE` to `wails3 doctor` output on Linux by @leaanthony
<!-- New features, capabilities, or enhancements -->

## Changed
<!-- Changes in existing functionality -->

## Fixed
- Fix window menu crash on Wayland caused by appmenu-gtk-module accessing unrealized window (#4769) by @leaanthony
- Fix GTK application crash when app name contains invalid characters (spaces, parentheses, etc.) by @leaanthony
- Fix "not enough memory" error when initializing drag and drop on Windows (#4701) by @overlordtm
- Fix file explorer opening wrong directory on Linux due to incorrect URI escaping (#4397) by @leaanthony
- Fix AppImage build failure on modern Linux distributions (Arch, Fedora 39+, Ubuntu 24.04+) by auto-detecting `.relr.dyn` ELF sections and disabling stripping (#4642) by @leaanthony
- Fix `wails doctor` falsely reporting webkit packages as installed on Fedora/DNF-based systems (#4457) by @leaanthony
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
